package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"os"
)

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

func normalize(val uint32) uint32 {
	// The thresholds are only correct for 44100 sample rate!
	if val < 35 {
		return 20
	}
	return 50
}

type Frequency int

const (
	FreqLow  Frequency = iota
	FreqHigh Frequency = iota
)

type Tone struct {
	Freq     Frequency
	Duration float64
}

type WavReader struct {
	audioFormat   uint16
	numChannels   uint16
	sampleRate    uint32
	byteRate      uint32
	blockAlign    uint16
	bitsPerSample uint16

	data    []int16
	channel chan Tone
}

func NewWavReader(filename string, channel chan Tone) *WavReader {
	reader := WavReader{}
	reader.init(filename)
	reader.channel = channel
	return &reader
}

func (self *WavReader) init(filename string) {
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
	self.audioFormat = readUint16(f)
	self.numChannels = readUint16(f)
	self.sampleRate = readUint32(f)
	self.byteRate = readUint32(f)
	self.blockAlign = readUint16(f)
	self.bitsPerSample = readUint16(f)
	check(self.bitsPerSample == 16, "bitsPerSample is not 16")

	log.Printf("audioFormat   = %d", self.audioFormat)
	log.Printf("numChannels   = %d", self.numChannels)
	log.Printf("sampleRate    = %d", self.sampleRate)
	log.Printf("byteRate      = %d", self.byteRate)
	log.Printf("blockAlign    = %d", self.blockAlign)
	log.Printf("bitsPerSample = %d", self.bitsPerSample)

	check(readUint32(f) == 0x61746164, "Not a 'data' chunk")
	subchunkSize = readUint32(f)

	// Read data, drop everything but the first channel
	log.Print("Loading data")
	numSamples := subchunkSize / uint32(self.numChannels) / uint32(self.bitsPerSample/8)
	self.data = make([]int16, numSamples)
	var i uint32
	for i = 0; i < numSamples; i++ {
		self.data[i] = int16(readUint16(f))
		f.Seek(int64((self.numChannels-1)*(self.bitsPerSample/8)), 1) // Skip other channels
	}
	log.Print("Data loaded")
	f.Close()
}

func (self *WavReader) emit(Δt, duration uint32) {
	var freq Frequency
	if Δt < 30 {
		freq = FreqHigh
	} else {
		freq = FreqLow
	}

	log.Printf("Δt = %d for %d cycles (%.4f secs)", Δt, duration, float64(duration)/float64(self.sampleRate))
	self.channel <- Tone{freq, float64(duration) / float64(self.sampleRate)}
}

func (self *WavReader) read() {
	var lastTransition uint32 = 0
	var lastΔt uint32 = 0
	var startPos uint32 = 0
	var i uint32
	for i = 0; i < uint32(len(self.data))-1; i++ {
		// find zero transitions
		if self.data[i] >= 0 && self.data[i-1] < 0 {
			// compute time since last zero transition, normalizing values
			Δt := normalize(i - lastTransition)
			if math.Abs(float64(int64(Δt)-int64(lastΔt))) > 5 {
				// frequency changed. Figure out how long we kept the same frequency.
				self.emit(lastΔt, i-startPos)
				lastΔt = Δt
				startPos = i
			}
			lastTransition = i
		}
	}
	self.emit(lastΔt, i-startPos)
	close(self.channel)
}
