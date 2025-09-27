package internal

import (
	"bytes"
	"log"
	"os"
	"runtime"
	"sort"

	"github.com/edsrzf/mmap-go"
)

func SolveFast(filename string) string {
	workers := runtime.NumCPU()

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	fileSize := int(info.Size())

	mmap, err := mmap.Map(file, mmap.RDONLY, 0)

	if err != nil {
		log.Fatal(err)
	}

	defer mmap.Unmap()

	offset := fileSize / workers

	if offset == 0 {
		workers = 1
		offset = fileSize
	}

	results := make(map[string]StationResult)
	send := make(chan map[string]StationResult, workers)

	for i := 0; i < workers; i++ {
		start := offset * i
		end := offset * (i + 1)

		if i == workers-1 {
			end = fileSize
		}

		go processMemorySection(mmap, start, end, send)
	}

	for i := 0; i < workers; i++ {
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

				results[workerName] = station
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

func orderResults(results map[string]StationResult) []StationResult {
	ordered := make([]StationResult, 0, len(results))

	for name := range results {
		ordered = append(ordered, results[name])
	}

	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].Name < ordered[j].Name
	})

	return ordered
}

func processMemorySection(data mmap.MMap, start, end int, send chan map[string]StationResult) {
	results := make(map[string]StationResult)

	position := int(start)

	if start > 0 && data[start-1] != '\n' {
		next_line_position := bytes.IndexByte(data[position:], '\n')
		if next_line_position != -1 {
			position = position + next_line_position + 1
		}
	}

	for position < end {
		nextlinePost := bytes.IndexByte(data[position:], '\n')
		if nextlinePost == -1 {
			break
		}
		newlinePos := position + nextlinePost

		line := data[position:newlinePos]

		if len(line) == 0 {
			position = newlinePos + 1
			continue
		}

		semi := bytes.IndexByte(line, ';')
		if semi <= 0 || semi == len(line)-1 {
			position = newlinePos + 1
			continue
		}

		name := string(line[:semi])
		temperature_bytes := line[semi+1:]

		if len(temperature_bytes) > 0 && temperature_bytes[len(temperature_bytes)-1] == '\r' {
			temperature_bytes = temperature_bytes[:len(temperature_bytes)-1]
		}

		temperature := temperatureToInt(temperature_bytes)

		if station, ok := results[name]; ok {
			station.Total += temperature

			station.Count++

			if temperature > station.Maximum {
				station.Maximum = temperature
			}

			if temperature < station.Minimum {
				station.Minimum = temperature
			}

			results[name] = station
		} else {
			results[name] = StationResult{
				Name:    name,
				Total:   temperature,
				Count:   1,
				Minimum: temperature,
				Maximum: temperature,
			}
		}

		position = newlinePos + 1
	}

	send <- results
}

func temperatureToInt(temperature []byte) int {
	result := 0
	negative := false

	for i := 0; i < len(temperature); i++ {
		char := temperature[i]
		if char == '-' {
			negative = true
			continue
		}
		if char == '.' {
			continue
		}
		result = result*10 + int(char-'0')
	}

	if negative {
		return -result
	}

	return result
}
