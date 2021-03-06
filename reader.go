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
	"log"
	"os"
)

type TapeReader struct {
	binName string
	wavName string
}

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
		tone1, more := <-toneChannel
		if !more {
			break
		}
		tone2, more := <-toneChannel
		if !more {
			break
		}
		cnt++
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

func bytesToFile(byteChannel chan byte, binName string) {
	f, err := os.Create(binName)
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
	log.Printf("Wrote %d bytes to file %s.", cnt, binName)
}

func reorderBytes(byteChannel chan byte, reorderedByteChannel chan byte) {
	// Kosmos CP1 writes bytes in this order: 0, 254, 253, ..., 1
	buf := make([]byte, 256)
	pos := 0
	block := 0
	for {
		b, more := <-byteChannel
		if !more {
			log.Printf("No more bytes to read")
			break
		}
		buf[pos] = b
		pos++
		if pos == 256 {
			log.Printf("Reordering Block %d", block)
			block++
			reorderedByteChannel <- buf[0]
			for i := 0; i < 255; i++ {
				reorderedByteChannel <- buf[255-i]
			}
			pos = 0
		}
	}
	log.Printf("Reordered %d blocks", block)
	close(reorderedByteChannel)
}

func (self *TapeReader) convert() {
	toneChannel := make(chan Tone)
	bitChannel := make(chan uint)
	byteChannel := make(chan byte)
	reorderedByteChannel := make(chan byte)
	wavReader := NewWavReader(self.wavName, toneChannel)
	go wavReader.Read()
	go tonesToBits(toneChannel, bitChannel)
	go bitsToBytes(bitChannel, byteChannel)
	go reorderBytes(byteChannel, reorderedByteChannel)
	bytesToFile(reorderedByteChannel, self.binName)
}
