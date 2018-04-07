package nameusage

import (
	"sort"
	"strings"
)

type NameReferenceSummary struct {
	Occurrences int
	References  int
	Name        string
}

func (Ω NameReferenceLedger) Names() (res []string) {
	for _, s := range Ω {
		res = append(res, s.Name)
	}
	return
}

type NameReferenceLedger []*NameReferenceSummary

func (Ω NameReferenceLedger) Len() int {
	return len(Ω)
}

func (Ω NameReferenceLedger) Swap(i, j int) {
	Ω[i], Ω[j] = Ω[j], Ω[i]
}

func (Ω NameReferenceLedger) Less(i, j int) bool {
	return Ω[i].References > Ω[j].References
}

func (Ω NameReferenceLedger) IncrementName(name string, occurrences int) NameReferenceLedger {

	name = strings.ToLower(name)

	exists := false
	for i := range Ω {
		if Ω[i].Name == name {
			Ω[i].Occurrences += occurrences
			Ω[i].References++
			exists = true
			break
		}
	}

	if !exists {
		Ω = append(Ω, &NameReferenceSummary{
			Occurrences: occurrences,
			References:  1,
			Name:        name,
		})
	}

	sort.Sort(Ω)

	return Ω

}
