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
