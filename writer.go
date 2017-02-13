package main

import (
	"log"
	"os"
)

type TapeWriter struct {
	binName string
	wavName string
}

func fileToBytes(binName string, byteChannel chan byte) {
	f, err := os.Open(binName)
	checkErr(err)
	defer f.Close()
	stat, err := f.Stat()
	checkErr(err)
	buf := make([]byte, stat.Size())
	f.Read(buf)
	log.Printf("Read %d bytes from file %s.", stat.Size(), binName)
	for _, b := range buf {
		byteChannel <- b
	}
	close(byteChannel)
}

func bytesToBits(byteChannel chan byte, bitChannel chan uint) {
	for {
		b, more := <-byteChannel
		if !more {
			log.Printf("No more bytes to read")
			break
		}
		var i uint
		for i = 8; i > 0; i-- {
			if (b & (1 << 7)) > 0 {
				bitChannel <- 1
			} else {
				bitChannel <- 0
			}
		}
	}
	close(bitChannel)
}

func bitsToTones(bitChannel chan uint, toneChannel chan Tone) {
	// Send 16 sec low freq intro
	// wait for the signal tone (16 secs low frequency)
	toneChannel <- Tone{Freq: FreqLow, Duration: 16}

	// Now, read bits
	cnt := 0
	for {
		cnt++
		bit, more := <-bitChannel
		if !more {
			break
		}
		if bit == 0 {
			// 0.066 secs high, 0.033 secs low
			toneChannel <- Tone{Freq: FreqLow, Duration: 0.066}
			toneChannel <- Tone{Freq: FreqHigh, Duration: 0.033}
		} else {
			// 0.033 secs high, 0.066 secs low
			toneChannel <- Tone{Freq: FreqLow, Duration: 0.033}
			toneChannel <- Tone{Freq: FreqHigh, Duration: 0.066}
		}
	}

	// Finally, a few secs high freq (0V)
	toneChannel <- Tone{Freq: FreqHigh, Duration: 5}

	close(toneChannel)

	log.Printf("Read %d bits", cnt)
}

func tonesToWav(toneChannel chan Tone, wavName string) {

}

func (self *TapeWriter) convert() {
	toneChannel := make(chan Tone)
	bitChannel := make(chan uint)
	byteChannel := make(chan byte)
	go fileToBytes(self.binName, byteChannel)
	go bytesToBits(byteChannel, bitChannel)
	go bitsToTones(bitChannel, toneChannel)
	tonesToWav(toneChannel, self.wavName)
}
