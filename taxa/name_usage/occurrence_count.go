package name_usage

import "bitbucket.org/heindl/taxa/store"

type SourceTargetOccurrenceCount map[store.DataSourceType]map[store.DataSourceTargetID]int

type CanonicalNameUsage struct {
	CanonicalName               string                                                  `json:",omitempty"`
	Synonyms                    []string                                                `json:",omitempty"`
	Ranks                       []string                                                `json:",omitempty"`
	SourceTargetOccurrenceCount SourceTargetOccurrenceCount `json:",omitempty"`
}

func (a SourceTargetOccurrenceCount) Intersects(b SourceTargetOccurrenceCount) bool {

	for srcType, counts := range b {
		for targetID, _ := range counts {
			if a.Contains(srcType, targetID) {
				return true
			}
		}
	}

	return false
}

func (Ω SourceTargetOccurrenceCount) Contains(srcType store.DataSourceType, targetId store.DataSourceTargetID) bool {
	if _, ok := Ω[srcType]; !ok {
		return false
	}

	if _, ok := Ω[srcType][targetId]; !ok {
		return false
	}

	return true
}

func (Ω SourceTargetOccurrenceCount) TargetIDCount() int {
	total := 0
	for _, counts := range Ω {
		total += len(counts)
	}
	return total
}

func (Ω SourceTargetOccurrenceCount) TotalOccurrenceCount() int {
	total := 0
	for _, counts := range Ω {
		for _, count := range counts {
			total += count
		}
	}
	return total
}

func (Ω SourceTargetOccurrenceCount) TargetIDs(srcType store.DataSourceType) store.DataSourceTargetIDs {

	res := store.DataSourceTargetIDs{}
	if _, ok := Ω[srcType]; !ok {
		return res
	}

	for _id, _ := range Ω[srcType] {
		res = res.AddToSet(_id)
	}

	return res
}

func (Ω SourceTargetOccurrenceCount) Add(srcType store.DataSourceType, targetId store.DataSourceTargetID, count int) {

	if _, ok := Ω[srcType]; !ok {
		Ω[srcType] = map[store.DataSourceTargetID]int{}
	}

	if _, ok := Ω[srcType][targetId]; ok {
		if Ω[srcType][targetId] < count {
			Ω[srcType][targetId] = count
		}
	} else {
		Ω[srcType][targetId] = count
	}
}

func (Ω SourceTargetOccurrenceCount) Set(srcType store.DataSourceType, targetId store.DataSourceTargetID, count int) {

	if _, ok := Ω[srcType]; !ok {
		Ω[srcType] = map[store.DataSourceTargetID]int{}
	}

	Ω[srcType][targetId] = count
}

func (Ω SourceTargetOccurrenceCount) SetAll(b SourceTargetOccurrenceCount) {
	for srcID, counts := range b {
		for targetID, count := range counts {
			Ω.Set(srcID, targetID, count)
		}
	}
}

func (Ω SourceTargetOccurrenceCount) AddAll(b SourceTargetOccurrenceCount) {
	for srcID, counts := range b {
		for targetID, count := range counts {
			Ω.Add(srcID, targetID, count)
		}
	}
}
