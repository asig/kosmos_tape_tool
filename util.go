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
