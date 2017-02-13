package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	wavname string
	binname string
)

func usage() {
	_, file := filepath.Split(os.Args[0])
	fmt.Printf("Usage: %s <wav-file> <target>\n", file)
	os.Exit(1)
}

func init() {
	if len(os.Args) != 3 {
		usage()
	}
	wavname = os.Args[1]
	binname = os.Args[2]
}

type converterState int

func tonesToBits(toneChannel chan Tone, bitChannel chan uint) {
	// wait for the signal tone (16 secs low frequency)
	for {
		tone := <-toneChannel
		if tone.Freq == FreqLow && tone.Duration > 15 {
			break
		}
	}

	// Now, read bits
	cnt := 0
	for {
		cnt++
		tone1, more := <-toneChannel
		if !more {
			break
		}
		tone2, more := <-toneChannel
		if !more {
			break
		}
		if tone1.Freq != FreqHigh || tone2.Freq != FreqLow {
			log.Printf("Unrecognized data, quitting")
			break
		}

		var bit uint
		if tone1.Duration > tone2.Duration {
			bit = 0
		} else {
			bit = 1
		}
		log.Printf("Bit %04d: %d", cnt, bit)
		bitChannel <- bit
	}
	close(bitChannel)

	log.Printf("DONE")
}

func bitsToBytes(bitChannel chan uint, byteChannel chan byte) {
	var cnt uint = 0
	var cur byte = 0
	for {
		bit, more := <-bitChannel
		if !more {
			log.Printf("No more bits to read")
			break
		}
		if bit > 0 {
			cur |= 1 << cnt
		}
		cnt++
		if cnt == 8 {
			byteChannel <- cur
			cur = 0
			cnt = 0
		}
	}
	close(byteChannel)
}

func bytesToFile(byteChannel chan byte) {
	f, err := os.Create(binname)
	checkErr(err)
	defer f.Close()
	cnt := 0
	for {
		b, more := <-byteChannel
		if !more {
			log.Printf("No more bytes to read")
			break
		}
		log.Printf("Byte %03d: %03d", cnt, b)
		f.Write([]byte{b})
		cnt++
	}
	log.Printf("Wrote %d bytes to file %s.", cnt, binname)
}

func main() {
	toneChannel := make(chan Tone)
	bitChannel := make(chan uint)
	byteChannel := make(chan byte)
	wavReader := NewWavReader(wavname, toneChannel)
	go wavReader.read()
	go tonesToBits(toneChannel, bitChannel)
	go bitsToBytes(bitChannel, byteChannel)
	bytesToFile(byteChannel)
}
