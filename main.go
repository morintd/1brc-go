package main

import (
	"1brc/internal"
	"errors"
	"log"
	"os"
	"time"

	_ "github.com/onsi/ginkgo/v2"
	_ "github.com/onsi/gomega"
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

	log.Printf("\nTook %s", elapsed)
	log.Println(result)
}
