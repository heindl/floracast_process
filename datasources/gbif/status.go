package gbif

type taxonomicStatus string

type taxonomicStatuses []taxonomicStatus

func (Ω taxonomicStatuses) Contains(b taxonomicStatus) bool {
	for _, status := range Ω {
		if status == b {
			return true
		}
	}
	return false
}

const (
	//taxonomicStatusAccepted           = taxonomicStatus("ACCEPTED")
	//taxonomicStatusDoubtful           = taxonomicStatus("DOUBTFUL")
	taxonomicStatusHeterotypicSynonym = taxonomicStatus("HETEROTYPIC_SYNONYM")
	taxonomicStatusHomotypicSynonym   = taxonomicStatus("HOMOTYPIC_SYNONYM")
	//taxonomicStatusMisapplied         = taxonomicStatus("MISAPPLIED")
	taxonomicStatusProparteSynonym = taxonomicStatus("PROPARTE_SYNONYM")
	taxonomicStatusSynonym         = taxonomicStatus("SYNONYM")
)
