package nameusage

import "strings"

type NameReferenceSummary struct {
	OccurrenceCount int
	ReferenceCount  int
	Name            string
}

type NameReferenceLedger []*NameReferenceSummary

func (s NameReferenceLedger) Len() int {return len(s)}
func (s NameReferenceLedger) Swap(i, j int) {s[i], s[j] = s[j], s[i]}
func (s NameReferenceLedger) Less(i, j int) bool {return s[i].ReferenceCount < s[j].ReferenceCount }

func (Ω NameReferenceLedger) IncrementName(name string, occurrences int) NameReferenceLedger {

	name = strings.ToLower(name)

	for _, ref := range Ω {
		if ref.Name == name {
			ref.OccurrenceCount += occurrences
			ref.ReferenceCount += 1
		}
		return Ω
	}

	Ω = append(Ω, &NameReferenceSummary{
		OccurrenceCount: occurrences,
		ReferenceCount:  1,
		Name:            name,
	})

	return Ω

}
