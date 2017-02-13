/*
 * Copyright (c) 2016 - 2017 Andreas Signer <asigner@gmail.com>
 *
 * This file is part of kosmos_tape_tool.
 *
 * kosmos_tape_tool is free software: you can redistribute it and/or
 * modify it under the terms of the GNU General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * kosmos_tape_tool is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kosmos_tape_tool.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"encoding/binary"
	"math"
	"os"
)

const maxAmplitude = 20000

func writeUint32(f *os.File, val uint32) {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, val)
	_, err := f.Write(buf)
	checkErr(err)
}

func writeUint16(f *os.File, val uint16) {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, val)
	_, err := f.Write(buf)
	checkErr(err)
}

type WavChunk struct {
	next *WavChunk
	data []byte
}

type WavWriter struct {
	sampleRate    uint32
	byteRate      uint32
	blockAlign    uint16
	bitsPerSample uint16

	firstChunk *WavChunk
	lastChunk  *WavChunk
}

func NewWavWriter() *WavWriter {
	writer := WavWriter{}
	writer.sampleRate = 44100
	writer.blockAlign = 2
	writer.byteRate = writer.sampleRate * uint32(writer.blockAlign)
	writer.bitsPerSample = 16

	writer.firstChunk = &WavChunk{next: nil, data: make([]byte, 0)}
	writer.lastChunk = writer.firstChunk

	return &writer
}

func (self *WavWriter) emit(tone Tone) {
	freq := 0.0
	if tone.Freq == FreqLow {
		freq = 1000
	} else {
		freq = 2000
	}

	samples := int(tone.Duration * float64(self.sampleRate))
	chunk := &WavChunk{next: nil, data: make([]byte, samples*2)}
	step := float64(freq) * 2 * math.Pi / float64(self.sampleRate)

	for i := 0; i < samples; i++ {
		v := uint16(int16(math.Sin(float64(i)*step) * maxAmplitude))
		binary.LittleEndian.PutUint16(chunk.data[2*i:], v)
	}

	self.lastChunk.next = chunk
	self.lastChunk = chunk
}

func (self *WavWriter) dataSize() uint32 {
	var size uint32 = 0
	cur := self.firstChunk
	for cur != nil {
		size += uint32(len(cur.data))
		cur = cur.next
	}
	return size
}

func (self *WavWriter) Write(filename string) {
	f, err := os.Create(filename)
	checkErr(err)

	dataSize := self.dataSize()

	writeUint32(f, 0x46464952) // "RIFF"
	writeUint32(f, uint32(4+24+8+dataSize))
	writeUint32(f, 0x45564157) // "WAVE"

	writeUint32(f, 0x20746d66) // "fmt "
	writeUint32(f, 16)         // size remaining header
	writeUint16(f, 1)          // PCM format
	writeUint16(f, 1)          // # of channels
	writeUint32(f, self.sampleRate)
	writeUint32(f, self.byteRate)
	writeUint16(f, self.blockAlign)
	writeUint16(f, self.bitsPerSample)

	writeUint32(f, 0x61746164) // "data"
	writeUint32(f, self.dataSize())
	cur := self.firstChunk
	for cur != nil {
		f.Write(cur.data)
		cur = cur.next
	}

	f.Close()
}
