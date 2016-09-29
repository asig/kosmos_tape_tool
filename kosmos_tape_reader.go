package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
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

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func check(condition bool, msg string) {
	if !condition {
		log.Printf(msg)
		os.Exit(1)
	}
}

func readUint32(f *os.File) uint32 {
	buf := make([]byte, 4)
	_, err := f.Read(buf)
	checkErr(err)
	return binary.LittleEndian.Uint32(buf)
}

func readUint16(f *os.File) uint16 {
	buf := make([]byte, 2)
	_, err := f.Read(buf)
	checkErr(err)
	return binary.LittleEndian.Uint16(buf)
}

func init() {
	if len(os.Args) != 2 {
		usage()
	}
	filename = os.Args[1]
}

func normalize(val uint32) uint32 {
	// The thresholds are only correct for 44100 sample rate!
	if val < 35 {
		return 20
	} else {
		return 50
	}
}

func main() {
	log.Println("Reading Kosmos Tape from " + filename)

	f, err := os.Open(filename)
	checkErr(err)

	check(readUint32(f) == 0x46464952, "Not a RIFF file")
	chunkSize := readUint32(f)
	log.Printf("Chunksize is %d", chunkSize)
	check(readUint32(f) == 0x45564157, "Not a WAVE file")

	check(readUint32(f) == 0x20746d66, "Not a 'fmt ' chunk")
	subchunkSize := readUint32(f)
	check(subchunkSize == 16, fmt.Sprintf("fmt Subchunk size is %d", subchunkSize))
	audioFormat := readUint16(f)
	numChannels := readUint16(f)
	sampleRate := readUint32(f)
	byteRate := readUint32(f)
	blockAlign := readUint16(f)
	bitsPerSample := readUint16(f)
	check(bitsPerSample == 16, "bitsPerSample is not 16")

	log.Printf("audioFormat   = %d", audioFormat)
	log.Printf("numChannels   = %d", numChannels)
	log.Printf("sampleRate    = %d", sampleRate)
	log.Printf("byteRate      = %d", byteRate)
	log.Printf("blockAlign    = %d", blockAlign)
	log.Printf("bitsPerSample = %d", bitsPerSample)

	check(readUint32(f) == 0x61746164, "Not a 'data' chunk")
	subchunkSize = readUint32(f)

	// Read data, drop everything but the first channel
	log.Print("Loading data")
	numSamples := subchunkSize / uint32(numChannels) / uint32(bitsPerSample/8)
	data := make([]int16, numSamples)
	var i uint32
	for i = 0; i < numSamples; i++ {
		data[i] = int16(readUint16(f))
		f.Seek(int64((numChannels-1)*(bitsPerSample/8)), 1) // Skip other channels
	}
	log.Print("Data loaded")

	var lastTransition uint32 = 0
	var lastΔt uint32 = 0
	var startPos uint32 = 0
	for i = 0; i < numSamples-1; i++ {
		// find zero transitions
		if data[i] >= 0 && data[i-1] < 0 {
			// compute time since last zero transition, normalizing values
			Δt := normalize(i - lastTransition)
			if math.Abs(float64(int64(Δt)-int64(lastΔt))) > 5 {
				// frequency changed. Figure out how long we kept the same frequency.
				l := i - startPos
				log.Printf("Δt = %d for %d cycles (%.4f secs)", lastΔt, l, float64(l)/float64(sampleRate))
				lastΔt = Δt
				startPos = i
			}
			lastTransition = i
		}
	}
	l := i - startPos
	log.Printf("Δt = %d for %d cycles (%.4f secs)", lastΔt, l, float64(l)/float64(sampleRate))

	log.Printf("numSamples    = %d", numSamples)

	f.Close()

}
