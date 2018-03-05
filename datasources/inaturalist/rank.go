package inaturalist

import "strings"

type rank string

const (
	// Originating from INaturalist:
	rankKingdom     = rank("Kingdom")
	rankPhylum      = rank("Phylum")
	rankSubPhylum   = rank("SubPhylum")
	rankClass       = rank("Class")
	rankSubClass    = rank("SubClass")
	rankOrder       = rank("Order")
	rankSuperFamily = rank("SuperFamily")
	rankFamily      = rank("Family")
	rankSubFamily   = rank("SubFamily")
	rankTribe       = rank("Tribe")
	rankSubTribe    = rank("SubTribe")
	rankGenus       = rank("Genus")
	rankSpecies     = rank("Species")
	rankSubSpecies  = rank("SubSpecies")
	rankForm        = rank("Form")
	rankVariety     = rank("Variety")
)

var taxonRankMap = map[string]rank{
	"kingdom":     rankKingdom,
	"phylum":      rankPhylum,
	"subphylum":   rankSubPhylum,
	"class":       rankClass,
	"subclass":    rankSubClass,
	"order":       rankOrder,
	"superfamily": rankSuperFamily,
	"family":      rankFamily,
	"subfamily":   rankSubFamily,
	"tribe":       rankTribe,
	"subtribe":    rankSubTribe,
	"genus":       rankGenus,
	"species":     rankSpecies,
	"subspecies":  rankSubSpecies,
	"form":        rankForm,
	"variety":     rankVariety,
}

func (Ω rank) Valid() bool {
	if _, ok := taxonRankMap[strings.ToLower(string(Ω))]; !ok {
		return false
	}
	return true
}

type rankLevel int

const (
	// Originating from INaturalist:
	//rankLevelKingdom     = rankLevel(70)
	//rankLevelPhylum      = rankLevel(60)
	//rankLevelSubPhylum   = rankLevel(57)
	//rankLevelClass       = rankLevel(50)
	//rankLevelSubClass    = rankLevel(47)
	//rankLevelOrder       = rankLevel(40)
	//rankLevelSuperFamily = rankLevel(33)
	//rankLevelFamily      = rankLevel(30)
	//rankLevelSubFamily   = rankLevel(27)
	//rankLevelTribe       = rankLevel(25)
	//rankLevelSubTribe    = rankLevel(24)
	//rankLevelGenus       = rankLevel(20)
	rankLevelSpecies = rankLevel(10)
	//rankLevelSubSpecies  = rankLevel(5)
	//rankLevelVariety     = rankLevel(5)
)
