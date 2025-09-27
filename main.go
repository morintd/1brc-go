package main

import (
	"1brc/internal"
	"errors"
	"log"
	"os"
	"time"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		log.Panic(errors.New("missing file path"))
	}

	filePath := args[1]

	start := time.Now()
	result := internal.SolveFast(filePath)
	elapsed := time.Since(start)

	log.Printf("\n%s", result)
	log.Printf("\nTook %s", elapsed)
}
