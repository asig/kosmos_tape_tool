package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	filename string
)

func usage() {
	_, file := filepath.Split(os.Args[0])
	fmt.Printf("Usage: %s <wav-file>\n", file)
	os.Exit(1)
}

func init() {
	if len(os.Args) != 2 {
		usage()
	}
	filename = os.Args[1]
}

type converterState int

const ( // iota is reset to 0
	beforeSyncTone converterState = iota
	inData         converterState = iota
	afterData                     = iota
)

// func convertToBits(input chan Tone, output chan int) {

// 	state := before_sync_tone

// 	for {
// 		tone := <-input
// 		if tone.Duration == 0 {
// 			break
// 		}
// 		switch(state) {
// 		case before_sync_tone:
// 			if tone.Freq == low && tone.Duration > 15 {
// 				state = in_data
// 			}
// 			break
// 		case in_data:
// 			if tone.Freq =
// 			case
// 		}
// 		log.Printf("CONSUMER: freq = %d, duration = %f", tone.Freq, tone.Duration)
// 	}

// }

func readBits(toneChannel chan Tone) {

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
		tone1 := <-toneChannel
		tone2 := <-toneChannel

		if tone1.Freq != FreqHigh || tone2.Freq != FreqLow {
			log.Printf("Unrecognized data, quitting")
			break
		}

		var bit int
		if tone1.Duration > tone2.Duration {
			bit = 0
		} else {
			bit = 1
		}
		log.Printf("Bit %04d: %d", cnt, bit)
	}

	log.Printf("DONE")
}

func main() {
	toneChannel := make(chan Tone)
	wavReader := NewWavReader(filename, toneChannel)
	go wavReader.read()
	readBits(toneChannel)
}
