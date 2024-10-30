package internal

import (
	"fmt"
	"math"
)

type StationResult struct {
	Name    string
	Total   int
	Count   int
	Minimum int
	Maximum int
}

func (result *StationResult) String() string {
	return fmt.Sprintf("%v:%v/%v/%v", result.Name, float64(result.Minimum)/10, math.Round(float64(result.Total)/float64(result.Count))/10, float64(result.Maximum)/10)
}
