package internal

import (
	"bufio"
	"bytes"
	"log"
	"os"
)

func SolveSlow(filename string) string {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	results := make(map[string]StationResult)
	i := 0

	for scanner.Scan() {
		line := scanner.Bytes()
		semi := bytes.IndexByte(line, ';')
		if semi <= 0 || semi == len(line)-1 {
			continue
		}

		name := string(line[:semi])
		temperature := temperatureToInt(line[semi+1:])

		if err != nil {
			log.Panic(err)
		}

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

		i += 1

		if i%10000000 == 0 {
			log.Printf("\n%v percent", i/10000000)
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
