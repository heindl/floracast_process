package nameusage

import (
	"bitbucket.org/heindl/taxa/datasources"
)

type nameUsageSourceMap map[datasources.DataSourceType]map[datasources.DataSourceTargetID]*NameUsageSource

func (a nameUsageSourceMap) intersects(b nameUsageSourceMap) bool {

	for _, targets := range b {
		for _, src := range targets {
			if a.contains(src) {
				return true
			}
		}
	}

	return false
}

func (Ω nameUsageSourceMap) contains(src *NameUsageSource) bool {
	if _, ok := Ω[src.sourceType]; !ok {
		return false
	}

	if _, ok := Ω[src.sourceType][src.targetID]; !ok {
		return false
	}

	return true
}

func (Ω nameUsageSourceMap) targetIDCount() int {
	total := 0
	for _, targets := range Ω {
		total += len(targets)
	}
	return total
}

func (Ω nameUsageSourceMap) totalOccurrenceCount() int {
	total := 0
	for _, targets := range Ω {
		for _, src := range targets {
			total += src.occurrenceCount
		}
	}
	return total
}

func (Ω nameUsageSourceMap) targetIDs(srcType datasources.DataSourceType) datasources.DataSourceTargetIDs {

	res := datasources.DataSourceTargetIDs{}
	if _, ok := Ω[srcType]; !ok {
		return res
	}

	for _id, _ := range Ω[srcType] {
		res = res.AddToSet(_id)
	}

	return res
}

func (Ω nameUsageSourceMap) set(src *NameUsageSource) {
	if _, ok := Ω[src.sourceType]; !ok {
		Ω[src.sourceType] = map[datasources.DataSourceTargetID]*NameUsageSource{}
	}
	Ω[src.sourceType][src.targetID] = src
}

func (Ω nameUsageSourceMap) merge(b nameUsageSourceMap) {
	for _, srcs := range b {
		for _, src := range srcs {
			Ω.set(src)
		}
	}
}
