package inaturalist

import "strings"

type TaxonRank string

const (
	// Originating from INaturalist:
	RankKingdom     = TaxonRank("Kingdom")
	RankPhylum      = TaxonRank("Phylum")
	RankSubPhylum   = TaxonRank("SubPhylum")
	RankClass       = TaxonRank("Class")
	RankSubClass    = TaxonRank("SubClass")
	RankOrder       = TaxonRank("Order")
	RankSuperFamily = TaxonRank("SuperFamily")
	RankFamily      = TaxonRank("Family")
	RankSubFamily   = TaxonRank("SubFamily")
	RankTribe       = TaxonRank("Tribe")
	RankSubTribe    = TaxonRank("SubTribe")
	RankGenus       = TaxonRank("Genus")
	RankSpecies     = TaxonRank("Species")
	RankSubSpecies  = TaxonRank("SubSpecies")
	RankForm        = TaxonRank("Form")
	RankVariety     = TaxonRank("Variety")
)

var TaxonRankMap = map[string]TaxonRank{
	"kingdom":     RankKingdom,
	"phylum":      RankPhylum,
	"subphylum":   RankSubPhylum,
	"class":       RankClass,
	"subclass":    RankSubClass,
	"order":       RankOrder,
	"superfamily": RankSuperFamily,
	"family":      RankFamily,
	"subfamily":   RankSubFamily,
	"tribe":       RankTribe,
	"subtribe":    RankSubTribe,
	"genus":       RankGenus,
	"species":     RankSpecies,
	"subspecies":  RankSubSpecies,
	"form":        RankForm,
	"variety":     RankVariety,
}

func (Ω TaxonRank) Valid() bool {
	if _, ok := TaxonRankMap[strings.ToLower(string(Ω))]; !ok {
		return false
	}
	return true
}

type RankLevel int

const (
	// Originating from INaturalist:
	RankLevelKingdom     = RankLevel(70)
	RankLevelPhylum      = RankLevel(60)
	RankLevelSubPhylum   = RankLevel(57)
	RankLevelClass       = RankLevel(50)
	RankLevelSubClass    = RankLevel(47)
	RankLevelOrder       = RankLevel(40)
	RankLevelSuperFamily = RankLevel(33)
	RankLevelFamily      = RankLevel(30)
	RankLevelSubFamily   = RankLevel(27)
	RankLevelTribe       = RankLevel(25)
	RankLevelSubTribe    = RankLevel(24)
	RankLevelGenus       = RankLevel(20)
	RankLevelSpecies     = RankLevel(10)
	RankLevelSubSpecies  = RankLevel(5)
	RankLevelVariety     = RankLevel(5)
)
