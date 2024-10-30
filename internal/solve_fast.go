package internal

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"runtime"
	"sort"
	"strings"
)

func SolveFast(filename string) string {
	workers := int64(runtime.NumCPU())

	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	info, err := file.Stat()

	if err != nil {
		log.Fatal(err)
	}

	fileSize := info.Size()

	offset := fileSize / int64(workers)

	results := make(map[string]*StationResult)
	send := make(chan map[string]*StationResult)

	for i := int64(0); i < workers; i++ {
		limit := offset * (i + 1)

		if i == workers-1 {
			limit = fileSize
		}

		go readSection(filename, offset*i, limit, send)
	}

	for i := int64(0); i < workers; i++ {
		workerResults := <-send

		for workerName := range workerResults {
			if station, ok := results[workerName]; ok {
				workerStation := workerResults[workerName]

				station.Total += workerStation.Total
				station.Count += workerStation.Count

				if workerStation.Maximum > station.Maximum {
					station.Maximum = workerStation.Maximum
				}

				if workerStation.Minimum < station.Minimum {
					station.Minimum = workerStation.Minimum
				}
			} else {
				results[workerName] = workerResults[workerName]
			}
		}
	}

	var buffer bytes.Buffer

	ordered := orderResults(results)
	count := len(ordered)

	for i, station := range ordered {
		buffer.WriteString(station.String())

		if i < count-1 {
			buffer.WriteString(";")
		}
	}

	return buffer.String()
}

func orderResults(results map[string]*StationResult) []*StationResult {
	ordered := make([]*StationResult, 0, len(results))

	for name := range results {
		ordered = append(ordered, results[name])
	}

	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].Name < ordered[j].Name
	})

	return ordered
}

func readSection(filename string, offset int64, limit int64, send chan map[string]*StationResult) {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	read := int64(0)
	results := make(map[string]*StationResult)

	file.Seek(offset, 0)
	reader := bufio.NewReader(file)

	if offset != 0 {
		b, _ := reader.ReadBytes('\n')
		read += int64(len(b))
	}

	for {
		if offset+read > limit {
			break
		}

		b, err := reader.ReadBytes('\n')

		if err != nil && err != io.EOF {
			log.Panic(err)
		}

		line := strings.TrimSpace(string(b))
		infos := strings.Split(line, ";")

		if len(infos) == 2 {
			name := infos[0]
			temperature := temperatureToInt(infos[1])

			if station, ok := results[name]; ok {
				station.Total += temperature

				station.Count += 1

				if temperature > station.Maximum {
					station.Maximum = temperature
				}

				if temperature < station.Minimum {
					station.Minimum = temperature
				}
			} else {
				results[name] = &StationResult{
					Name:    name,
					Total:   temperature,
					Count:   1,
					Minimum: temperature,
					Maximum: temperature,
				}
			}
		}

		if err == io.EOF {
			break
		}

		read += int64(len(b))
	}

	send <- results
}

func temperatureToInt(temperature string) int {
	result := 0
	negative := false

	for _, char := range temperature {
		if char == 45 {
			negative = true
			continue
		}

		if char == 46 {
			continue
		}

		result = result*10 + int(char-48)
	}

	if negative {
		return -result
	}

	return result
}
