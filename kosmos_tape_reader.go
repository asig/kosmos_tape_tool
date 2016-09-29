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
	before_sync_tone converterState = iota
	in_data          converterState = iota
	after_data                      = iota
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
	for {
		tone := <-toneChannel
		if tone.Duration == 0 {
			break
		}
		log.Printf("CONSUMER: freq = %d, duration = %f", tone.Freq, tone.Duration)
	}
	log.Printf("DONE")
}

func main() {
	toneChannel := make(chan Tone)
	wavReader := NewWavReader(filename, toneChannel)
	go wavReader.read()
	go readBits(toneChannel)
}
