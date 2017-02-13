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
