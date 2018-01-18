package terra

import (
	"sort"
	"github.com/gonum/stat"
	"github.com/montanaflynn/stats"
)

func(Ω FeatureCollection) AreaStats() *BasicStats {
	areas := []float64{}
	for _, feature := range Ω.Features(){
		areas = append(areas, feature.Area())
	}
	if len(areas) == 0 {
		return nil
	}

	sort.Float64s(areas)

	quartiles, err := stats.Quartile(areas)
	if err != nil {
		panic(err)
	}

	mode, err := stats.Mode(areas)
	if err != nil {
		panic(err)
	}

	mean, std := stat.MeanStdDev(areas, nil)

	variance, err := stats.Variance(stats.Float64Data(areas))
	if err != nil {
		panic(err)
	}
	aboveZero := 0
	for i := range areas {
		if areas[i] > 0 {
			aboveZero += 1
		}
	}

	return &BasicStats{
		AboveZero: aboveZero,
		Max: areas[len(areas)-1],
		Min: areas[0],
		Mean: mean,
		StandardDeviation: std,
		Median: variance,
		Mode: mode,
		Count: float64(len(areas)),
		Quantiles: quartiles,
	}
}

type BasicStats struct{
	Mean, Median, Max, Min, StandardDeviation, Variance, Count  float64
	Quantiles stats.Quartiles
	Mode []float64
	AboveZero int
}