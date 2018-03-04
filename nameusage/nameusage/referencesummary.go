package nameusage

import (
	"sort"
	"strings"
)

type nameReferenceSummary struct {
	OccurrenceCount int
	ReferenceCount  int
	Name            string
}

type nameReferenceLedger []*nameReferenceSummary

func (Ω nameReferenceLedger) Len() int {
	return len(Ω)
}

func (Ω nameReferenceLedger) Swap(i, j int) {
	Ω[i], Ω[j] = Ω[j], Ω[i]
}

func (Ω nameReferenceLedger) Less(i, j int) bool {
	return Ω[i].ReferenceCount > Ω[j].ReferenceCount
}

func (Ω nameReferenceLedger) IncrementName(name string, occurrences int) nameReferenceLedger {

	name = strings.ToLower(name)

	exists := false
	for i := range Ω {
		if Ω[i].Name == name {
			Ω[i].OccurrenceCount += occurrences
			Ω[i].ReferenceCount++
			exists = true
			break
		}
	}

	if !exists {
		Ω = append(Ω, &nameReferenceSummary{
			OccurrenceCount: occurrences,
			ReferenceCount:  1,
			Name:            name,
		})
	}

	sort.Sort(Ω)

	return Ω

}
