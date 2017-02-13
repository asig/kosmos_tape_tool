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
	if stat.Size()%256 != 0 {
		panic("File size is not a multiple of 256!")
	}
	buf := make([]byte, 256)
	for blocks := stat.Size() / 256; blocks > 0; blocks-- {
		f.Read(buf)
		byteChannel <- buf[0]
		for j := 255; j > 0; j-- {
			byteChannel <- buf[j]
		}
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
		for i = 0; i < 8; i++ {
			if (b & (1 << i)) > 0 {
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
		bit, more := <-bitChannel
		if !more {
			break
		}
		cnt++
		// The Kosmos CP1 uses 30ms and 60ms slices, but there's some
		// overhead, so let's instead use 35ms and 65ms
		if bit == 0 {
			// 0.065 secs high, 0.035 secs low
			toneChannel <- Tone{Freq: FreqLow, Duration: 0.065}
			toneChannel <- Tone{Freq: FreqHigh, Duration: 0.035}
		} else {
			// 0.035 secs high, 0.065 secs low
			toneChannel <- Tone{Freq: FreqLow, Duration: 0.035}
			toneChannel <- Tone{Freq: FreqHigh, Duration: 0.065}
		}
	}

	// Finally, a few secs high freq (0V)
	toneChannel <- Tone{Freq: FreqHigh, Duration: 5}

	close(toneChannel)

	log.Printf("Read %d bits", cnt)
}

func tonesToWav(toneChannel chan Tone, wavName string) {
	writer := NewWavWriter()
	for {
		tone, more := <-toneChannel
		if !more {
			break
		}
		writer.emit(tone)
	}
	writer.Write(wavName)
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
