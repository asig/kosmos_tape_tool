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
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var (
	wavName = flag.String("wav", "", "WAV-file to read/write")
	binName = flag.String("bin", "", "bin-file to read/write")
)

func usage() {
	_, file := filepath.Split(os.Args[0])
	fmt.Printf("Usage: %s -wav=<wav-file> -bin=<bin-file> <read|write>\n", file)
	os.Exit(1)
}

func init() {
	flag.Parse()
	if len(flag.Args()) != 1 || *wavName == "" || *binName == "" {
		usage()
	}
}

func main() {
	op := flag.Arg(0)
	if op == "read" {
		reader := TapeReader{wavName: *wavName, binName: *binName}
		reader.convert()
	} else if op == "write" {
		writer := TapeWriter{wavName: *wavName, binName: *binName}
		writer.convert()
	} else {
		usage()
	}
}
