package materialtaxa

func (Ω *NameUsage) materialize() (map[string]interface{}, error) {

	name, err := Ω.CommonName()
	if err != nil {
		return nil, err
	}

	m := map[string]interface{}{
		"ScientificName": Ω.canonicalName.ScientificName(),
		"CommonName": name,
	}

	return m, nil
}