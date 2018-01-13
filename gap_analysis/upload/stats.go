package main

import (
	"bitbucket.org/heindl/taxa/store"
	"bitbucket.org/heindl/taxa/utils"
	"fmt"
	"github.com/montanaflynn/stats"
	"github.com/saleswise/errors/errors"
	"sync"
)

type StatsContainer struct {
	sync.Mutex
	States map[string]Stats
}

type Stats struct {
	Areas            stats.Float64Data
	Acres            stats.Float64Data
	Invalid          map[string]int
	CentroidDistance stats.Float64Data
}

func NewStatsContainer() *StatsContainer {
	states := map[string]Stats{}
	for _, name := range store.ValidProtectedAreaStates {
		states[name] = Stats{
			Areas:            stats.Float64Data{},
			Acres:            stats.Float64Data{},
			Invalid:          map[string]int{},
			CentroidDistance: stats.Float64Data{},
		}
	}

	return &StatsContainer{
		States: states,
	}
}

func (Ω *Stats) Print(label string) error {

	invalidCount := 0
	for _, v := range Ω.Invalid {
		invalidCount += v
	}

	fmt.Println("------", label, "-", len(Ω.Areas)-invalidCount, "-", invalidCount, "------")

	if len(Ω.Areas) == 0 && invalidCount == 0 {
		fmt.Println("")
		return nil
	}

	if err := print_stats("Areas:", Ω.Areas); err != nil {
		return err
	}
	if err := print_stats("Acres:", Ω.Acres); err != nil {
		return err
	}

	if err := print_stats("CentroidDistance:", Ω.CentroidDistance); err != nil {
		return err
	}
	fmt.Println("Invalid: ", utils.JsonOrSpew(Ω.Invalid))
	fmt.Println("")
	return nil
}

func print_stats(label string, s stats.Float64Data) error {
	if len(s) == 0 {
		fmt.Println(label, 0, 0, 0, 0)
		return nil
	}
	min, err := s.Min()
	if err != nil {
		return errors.Wrap(err, "couldn't print min")
	}
	max, err := s.Max()
	if err != nil {
		return errors.Wrap(err, "couldn't print max")
	}
	median, err := s.Median()
	if err != nil {
		return errors.Wrap(err, "couldn't print median")
	}
	mean, err := s.Mean()
	if err != nil {
		return errors.Wrap(err, "couldn't print mean")
	}
	fmt.Println(label, fmt.Sprintf("%6.2f, %6.2f, %6.2f, %6.2f", min, max, median, mean))
	return nil
}

func (Ω *StatsContainer) Print() error {
	// Compile stats.
	total := Stats{
		Areas:            stats.Float64Data{},
		Acres:            stats.Float64Data{},
		Invalid:          map[string]int{},
		CentroidDistance: stats.Float64Data{},
	}
	for state, sts := range Ω.States {
		if err := sts.Print(state); err != nil {
			return err
		}
		total.Areas = append(total.Areas, sts.Areas...)
		total.Acres = append(total.Acres, sts.Acres...)
		total.CentroidDistance = append(total.CentroidDistance, sts.CentroidDistance...)
		for reason, count := range sts.Invalid {
			if _, ok := total.Invalid[reason]; !ok {
				total.Invalid[reason] = count
			} else {
				total.Invalid[reason] += count
			}
		}
	}
	return total.Print("total")
}

func (p *Parser) UpdateStats(pa *store.ProtectedArea, valid bool, reason string, minDistanceFromNearestPoint float64) error {
	p.Stats.Lock()
	defer p.Stats.Unlock()

	// Only working within a specific range of states.
	if reason == "state" {
		return nil
	}

	sts := p.Stats.States[pa.StateAbbr]

	sts.Areas = append(sts.Areas, pa.Height*pa.Width)
	sts.Acres = append(sts.Acres, pa.GISAcres)
	sts.CentroidDistance = append(sts.CentroidDistance, minDistanceFromNearestPoint)

	if !valid {
		if _, ok := sts.Invalid[reason]; !ok {
			sts.Invalid[reason] = 1
		} else {
			sts.Invalid[reason] += 1
		}
	}

	p.Stats.States[pa.StateAbbr] = sts

	return nil

}
