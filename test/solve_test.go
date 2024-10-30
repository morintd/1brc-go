package test

import (
	"1brc/internal"
	"bufio"
	"bytes"
	"log"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("POST /barcode", Ordered, func() {
	expected := getExpected("../assets/expected-medium.txt")
	measurements_path := "../assets/medium.txt"

	It("(Slow) Should return ordered stations with min/max/average", func() {
		result := internal.SolveSlow(measurements_path)
		Expect(result).To(Equal(expected))
	})

	It("(Fast) Should return ordered stations with min/max/average", func() {
		result := internal.SolveFast(measurements_path)
		Expect(result).To(Equal(expected))
	})
})

func getExpected(path string) string {
	file, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)

	var buffer bytes.Buffer

	for scanner.Scan() {
		buffer.WriteString(scanner.Text())
	}

	file.Close()

	return buffer.String()
}
