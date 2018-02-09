package materialtaxa

type description struct {
	Citation string `json:""`
	Text string `json:""`
}

func (Ω *NameUsage) descriptions() ([]description, error) {

	for _, src := range Ω.Sources() {
		if src.occurrenceCount > 0 {

		}
	}

	return nil, nil

}