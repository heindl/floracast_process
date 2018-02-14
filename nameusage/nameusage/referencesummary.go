package nameusage

import (
	"strings"
	"sort"
)

type NameReferenceSummary struct {
	OccurrenceCount int
	ReferenceCount  int
	Name            string
}

type NameReferenceLedger []*NameReferenceSummary

func (s NameReferenceLedger) Len() int {
	return len(s)
}

func (s NameReferenceLedger) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s NameReferenceLedger) Less(i, j int) bool {
	return s[i].ReferenceCount > s[j].ReferenceCount
}

func (Ω NameReferenceLedger) IncrementName(name string, occurrences int) NameReferenceLedger {

	name = strings.ToLower(name)

	exists := false
	for i := range Ω {
		if Ω[i].Name == name {
			Ω[i].OccurrenceCount += occurrences
			Ω[i].ReferenceCount += 1
			exists = true
			break
		}
	}

	if !exists {
		Ω = append(Ω, &NameReferenceSummary{
			OccurrenceCount: occurrences,
			ReferenceCount:  1,
			Name:            name,
		})
	}

	sort.Sort(Ω)

	return Ω

}
