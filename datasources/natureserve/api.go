package natureserve

import (
	"encoding/xml"
	"github.com/kennygrant/sanitize"
)

/*****
			SPECIES SEARCH REPORT
*****/

type commonName struct {
	Attrlanguage string   `xml:"language,attr"  json:",omitempty"` // maxLength=2
	Text         string   `xml:",chardata" json:",omitempty"`      // maxLength=20
	XMLName      xml.Name `xml:"commonName,omitempty" json:"commonName,omitempty"`
}

type conditions struct {
	Match   *match   `xml:"match,omitempty" json:"match,omitempty"` // ZZmaxLength=0
	XMLName xml.Name `xml:"conditions,omitempty" json:"conditions,omitempty"`
}

type globalSpeciesUID struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=23
	XMLName xml.Name `xml:"globalSpeciesUID,omitempty" json:"globalSpeciesUID,omitempty"`
}

type ii struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=25
	XMLName xml.Name `xml:"i,omitempty" json:"i,omitempty"`
}

type jurisdictionScientificName struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=35
	XMLName xml.Name `xml:"jurisdictionScientificName,omitempty" json:"jurisdictionScientificName,omitempty"`
}

type match struct {
	SearchTerm *searchTerm `xml:"searchTerm,omitempty" json:"searchTerm,omitempty"` // ZZmaxLength=0
	XMLName    xml.Name    `xml:"match,omitempty" json:"match,omitempty"`
}

type natureServeExplorerURI struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=94
	XMLName xml.Name `xml:"natureServeExplorerURI,omitempty" json:"natureServeExplorerURI,omitempty"`
}

type searchTerm struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=10
	XMLName xml.Name `xml:"searchTerm,omitempty" json:"searchTerm,omitempty"`
}

type speciesSearchReport struct {
	AttrXsiSpaceschemaLocation string                   `xml:"http://www.w3.org/2001/XMLSchema-instance schemaLocation,attr"  json:",omitempty"` // maxLength=136
	AttrschemaVersion          string                   `xml:"schemaVersion,attr"  json:",omitempty"`                                            // maxLength=3
	Attrxmlns                  string                   `xml:"xmlns,attr"  json:",omitempty"`                                                    // maxLength=67
	AttrXmlnsxsi               string                   `xml:"xmlns xsi,attr"  json:",omitempty"`                                                // maxLength=41
	Conditions                 *conditions              `xml:"conditions,omitempty" json:"conditions,omitempty"`                                 // ZZmaxLength=0
	SpeciesSearchResultList    *speciesSearchResultList `xml:"speciesSearchResultList,omitempty" json:"speciesSearchResultList,omitempty"`       // ZZmaxLength=0
	XMLName                    xml.Name                 `xml:"speciesSearchReport,omitempty" json:"speciesSearchReport,omitempty"`
}

type speciesSearchResult struct {
	Attruid                    string                      `xml:"uid,attr"  json:",omitempty"`                                                      // maxLength=23
	CommonName                 *commonName                 `xml:"commonName,omitempty" json:"commonName,omitempty"`                                 // ZZmaxLength=0
	GlobalSpeciesUID           *globalSpeciesUID           `xml:"globalSpeciesUID,omitempty" json:"globalSpeciesUID,omitempty"`                     // ZZmaxLength=0
	JurisdictionScientificName *jurisdictionScientificName `xml:"jurisdictionScientificName,omitempty" json:"jurisdictionScientificName,omitempty"` // ZZmaxLength=0
	NatureServeExplorerURI     *natureServeExplorerURI     `xml:"natureServeExplorerURI,omitempty" json:"natureServeExplorerURI,omitempty"`         // ZZmaxLength=0
	TaxonomicComments          *taxonomicComments          `xml:"taxonomicComments,omitempty" json:"taxonomicComments,omitempty"`                   // ZZmaxLength=0
	XMLName                    xml.Name                    `xml:"speciesSearchResult,omitempty" json:"speciesSearchResult,omitempty"`
}

type speciesSearchResultList struct {
	SpeciesSearchResult []*speciesSearchResult `xml:"speciesSearchResult,omitempty" json:"speciesSearchResult,omitempty"` // ZZmaxLength=0
	XMLName             xml.Name               `xml:"speciesSearchResultList,omitempty" json:"speciesSearchResultList,omitempty"`
}

type taxonomicComments struct {
	I       []*ii    `xml:"i,omitempty" json:"i,omitempty"` // ZZmaxLength=0
	Text    string   `xml:",chardata" json:",omitempty"`    // maxLength=733
	XMLName xml.Name `xml:"taxonomicComments,omitempty" json:"taxonomicComments,omitempty"`
}

/*********************************************

SPECIES COMPREHENSIVE RESULTS

**********************************************/

type globalSpeciesList struct {
	AttrXsiSpaceschemaLocation string           `xml:"http://www.w3.org/2001/XMLSchema-instance schemaLocation,attr"  json:",omitempty"` // maxLength=136
	AttrschemaVersion          string           `xml:"schemaVersion,attr"  json:",omitempty"`                                            // maxLength=3
	Attrxmlns                  string           `xml:"xmlns,attr"  json:",omitempty"`                                                    // maxLength=67
	AttrXmlnsxsi               string           `xml:"xmlns xsi,attr"  json:",omitempty"`                                                // maxLength=41
	GlobalSpecies              []*globalSpecies `xml:"globalSpecies,omitempty" json:"globalSpecies,omitempty"`                           // ZZmaxLength=0
	XMLName                    xml.Name         `xml:"globalSpeciesList,omitempty" json:"globalSpeciesList,omitempty"`
}

type alternateSeparationProcedure struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=528
	XMLName xml.Name `xml:"alternateSeparationProcedure,omitempty" json:"alternateSeparationProcedure,omitempty"`
}

type barriers struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=209
	XMLName xml.Name `xml:"barriers,omitempty" json:"barriers,omitempty"`
}

type biologicalResearchNeeds struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=32
	XMLName xml.Name `xml:"biologicalResearchNeeds,omitempty" json:"biologicalResearchNeeds,omitempty"`
}

type citation struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=314
	XMLName xml.Name `xml:"citation,omitempty" json:"citation,omitempty"`
}

type class struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=7
	XMLName xml.Name `xml:"class,omitempty" json:"class,omitempty"`
}

type classification struct {
	Names    *names    `xml:"names,omitempty" json:"names,omitempty"`       // ZZmaxLength=0
	Taxonomy *taxonomy `xml:"taxonomy,omitempty" json:"taxonomy,omitempty"` // ZZmaxLength=0
	XMLName  xml.Name  `xml:"classification,omitempty" json:"classification,omitempty"`
}

type classificationStatus struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=8
	XMLName xml.Name `xml:"classificationStatus,omitempty" json:"classificationStatus,omitempty"`
}

type code struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=6
	XMLName xml.Name `xml:"code,omitempty" json:"code,omitempty"`
}

type comments struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=273
	XMLName xml.Name `xml:"comments,omitempty" json:"comments,omitempty"`
}

type conceptReference struct {
	Attrcode                   string                      `xml:"code,attr"  json:",omitempty"`                                                     // maxLength=12
	ClassificationStatus       *classificationStatus       `xml:"classificationStatus,omitempty" json:"classificationStatus,omitempty"`             // ZZmaxLength=0
	FormattedFullCitation      *formattedFullCitation      `xml:"formattedFullCitation,omitempty" json:"formattedFullCitation,omitempty"`           // ZZmaxLength=0
	NameUsedInConceptReference *nameUsedInConceptReference `xml:"nameUsedInConceptReference,omitempty" json:"nameUsedInConceptReference,omitempty"` // ZZmaxLength=0
	XMLName                    xml.Name                    `xml:"conceptReference,omitempty" json:"conceptReference,omitempty"`
}

type conservationStatus struct {
	NatureServeStatus *natureServeStatus `xml:"natureServeStatus,omitempty" json:"natureServeStatus,omitempty"` // ZZmaxLength=0
	OtherStatuses     *otherStatuses     `xml:"otherStatuses,omitempty" json:"otherStatuses,omitempty"`         // ZZmaxLength=0
	XMLName           xml.Name           `xml:"conservationStatus,omitempty" json:"conservationStatus,omitempty"`
}

type conservationStatusAuthors struct {
	AttrdisplayValue string   `xml:"displayValue,attr"  json:",omitempty"` // maxLength=27
	XMLName          xml.Name `xml:"conservationStatusAuthors,omitempty" json:"conservationStatusAuthors,omitempty"`
}

type conservationStatusFactors struct {
	ConservationStatusAuthors    *conservationStatusAuthors    `xml:"conservationStatusAuthors,omitempty" json:"conservationStatusAuthors,omitempty"`       // ZZmaxLength=0
	EstimatedNumberOfOccurrences *estimatedNumberOfOccurrences `xml:"estimatedNumberOfOccurrences,omitempty" json:"estimatedNumberOfOccurrences,omitempty"` // ZZmaxLength=0
	GlobalAbundance              *globalAbundance              `xml:"globalAbundance,omitempty" json:"globalAbundance,omitempty"`                           // ZZmaxLength=0
	GlobalInventoryNeeds         *globalInventoryNeeds         `xml:"globalInventoryNeeds,omitempty" json:"globalInventoryNeeds,omitempty"`                 // ZZmaxLength=0
	GlobalProtection             *globalProtection             `xml:"globalProtection,omitempty" json:"globalProtection,omitempty"`                         // ZZmaxLength=0
	GlobalShortTermTrend         *globalShortTermTrend         `xml:"globalShortTermTrend,omitempty" json:"globalShortTermTrend,omitempty"`                 // ZZmaxLength=0
	IntrinsicVulnerability       *intrinsicVulnerability       `xml:"intrinsicVulnerability,omitempty" json:"intrinsicVulnerability,omitempty"`             // ZZmaxLength=0
	OtherConsiderations          *otherConsiderations          `xml:"otherConsiderations,omitempty" json:"otherConsiderations,omitempty"`                   // ZZmaxLength=0
	StatusFactorsEditionDate     *statusFactorsEditionDate     `xml:"statusFactorsEditionDate,omitempty" json:"statusFactorsEditionDate,omitempty"`         // ZZmaxLength=0
	Threat                       *threat                       `xml:"threat,omitempty" json:"threat,omitempty"`                                             // ZZmaxLength=0
	XMLName                      xml.Name                      `xml:"conservationStatusFactors,omitempty" json:"conservationStatusFactors,omitempty"`
}

type conservationStatusMap struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=51
	XMLName xml.Name `xml:"conservationStatusMap,omitempty" json:"conservationStatusMap,omitempty"`
}

type countyCode struct {
	Attrtype string   `xml:"type,attr"  json:",omitempty"` // maxLength=4
	Text     string   `xml:",chardata" json:",omitempty"`  // maxLength=5
	XMLName  xml.Name `xml:"countyCode,omitempty" json:"countyCode,omitempty"`
}

type countyDistribution struct {
	OccurrenceNations *occurrenceNations `xml:"occurrenceNations,omitempty" json:"occurrenceNations,omitempty"` // ZZmaxLength=0
	XMLName           xml.Name           `xml:"countyDistribution,omitempty" json:"countyDistribution,omitempty"`
}

type countyName struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=10
	XMLName xml.Name `xml:"countyName,omitempty" json:"countyName,omitempty"`
}

type currentPresenceAbsence struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=7
	XMLName xml.Name `xml:"currentPresenceAbsence,omitempty" json:"currentPresenceAbsence,omitempty"`
}

type delineation struct {
	AttrgroupName                     string                             `xml:"groupName,attr"  json:",omitempty"`                                                              // maxLength=30
	AttrmigratoryUseType              string                             `xml:"migratoryUseType,attr"  json:",omitempty"`                                                       // maxLength=14
	DelineationAuthors                *delineationAuthors                `xml:"delineationAuthors,omitempty" json:"delineationAuthors,omitempty"`                               // ZZmaxLength=0
	InferredMinimumExtentOfHabitatUse *inferredMinimumExtentOfHabitatUse `xml:"inferredMinimumExtentOfHabitatUse,omitempty" json:"inferredMinimumExtentOfHabitatUse,omitempty"` // ZZmaxLength=0
	MappingGuidance                   *mappingGuidance                   `xml:"mappingGuidance,omitempty" json:"mappingGuidance,omitempty"`                                     // ZZmaxLength=0
	MinimumCriteriaForOccurrence      *minimumCriteriaForOccurrence      `xml:"minimumCriteriaForOccurrence,omitempty" json:"minimumCriteriaForOccurrence,omitempty"`           // ZZmaxLength=0
	Notes                             *notes                             `xml:"notes,omitempty" json:"notes,omitempty"`                                                         // ZZmaxLength=0
	Separation                        *separation                        `xml:"separation,omitempty" json:"separation,omitempty"`                                               // ZZmaxLength=0
	VersionDate                       *versionDate                       `xml:"versionDate,omitempty" json:"versionDate,omitempty"`                                             // ZZmaxLength=0
	XMLName                           xml.Name                           `xml:"delineation,omitempty" json:"delineation,omitempty"`
}

type delineationAuthors struct {
	AttrdisplayValue string   `xml:"displayValue,attr"  json:",omitempty"` // maxLength=16
	XMLName          xml.Name `xml:"delineationAuthors,omitempty" json:"delineationAuthors,omitempty"`
}

type delineations struct {
	Delineation *delineation `xml:"delineation,omitempty" json:"delineation,omitempty"` // ZZmaxLength=0
	XMLName     xml.Name     `xml:"delineations,omitempty" json:"delineations,omitempty"`
}

type description struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=79
	XMLName xml.Name `xml:"description,omitempty" json:"description,omitempty"`
}

type distance struct {
	Attrunits string   `xml:"units,attr"  json:",omitempty"` // maxLength=10
	Text      string   `xml:",chardata" json:",omitempty"`   // maxLength=1
	XMLName   xml.Name `xml:"distance,omitempty" json:"distance,omitempty"`
}

type distanceForSuitableHabitat struct {
	Attrunits string   `xml:"units,attr"  json:",omitempty"` // maxLength=10
	Text      string   `xml:",chardata" json:",omitempty"`   // maxLength=2
	XMLName   xml.Name `xml:"distanceForSuitableHabitat,omitempty" json:"distanceForSuitableHabitat,omitempty"`
}

type distanceForUnsuitableHabitat struct {
	Attrunits string   `xml:"units,attr"  json:",omitempty"` // maxLength=10
	Text      string   `xml:",chardata" json:",omitempty"`   // maxLength=1
	XMLName   xml.Name `xml:"distanceForUnsuitableHabitat,omitempty" json:"distanceForUnsuitableHabitat,omitempty"`
}

type distribution struct {
	AttrUSAndCanadianDistributionComplete string                 `xml:"USAndCanadianDistributionComplete,attr"  json:",omitempty"`              // maxLength=5
	ConservationStatusMap                 *conservationStatusMap `xml:"conservationStatusMap,omitempty" json:"conservationStatusMap,omitempty"` // ZZmaxLength=0
	CountyDistribution                    *countyDistribution    `xml:"countyDistribution,omitempty" json:"countyDistribution,omitempty"`       // ZZmaxLength=0
	Endemism                              *endemism              `xml:"endemism,omitempty" json:"endemism,omitempty"`                           // ZZmaxLength=0
	GlobalRange                           *globalRange           `xml:"globalRange,omitempty" json:"globalRange,omitempty"`                     // ZZmaxLength=0
	Nations                               *nations               `xml:"nations,omitempty" json:"nations,omitempty"`                             // ZZmaxLength=0
	Watersheds                            *watersheds            `xml:"watersheds,omitempty" json:"watersheds,omitempty"`                       // ZZmaxLength=0
	XMLName                               xml.Name               `xml:"distribution,omitempty" json:"distribution,omitempty"`
}

type distributionConfidence struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=9
	XMLName xml.Name `xml:"distributionConfidence,omitempty" json:"distributionConfidence,omitempty"`
}

type durations struct {
	XMLName xml.Name `xml:"durations,omitempty" json:"durations,omitempty"`
}

type ecologyAndLifeHistory struct {
	Durations                        *durations                        `xml:"durations,omitempty" json:"durations,omitempty"`                                               // ZZmaxLength=0
	EcologyAndLifeHistoryAuthors     *ecologyAndLifeHistoryAuthors     `xml:"ecologyAndLifeHistoryAuthors,omitempty" json:"ecologyAndLifeHistoryAuthors,omitempty"`         // ZZmaxLength=0
	EcologyAndLifeHistoryDescription *ecologyAndLifeHistoryDescription `xml:"ecologyAndLifeHistoryDescription,omitempty" json:"ecologyAndLifeHistoryDescription,omitempty"` // ZZmaxLength=0
	EcologyAndLifeHistoryEditionDate *ecologyAndLifeHistoryEditionDate `xml:"ecologyAndLifeHistoryEditionDate,omitempty" json:"ecologyAndLifeHistoryEditionDate,omitempty"` // ZZmaxLength=0
	FoodHabits                       *foodHabits                       `xml:"foodHabits,omitempty" json:"foodHabits,omitempty"`                                             // ZZmaxLength=0
	Habitats                         *habitats                         `xml:"habitats,omitempty" json:"habitats,omitempty"`                                                 // ZZmaxLength=0
	Migration                        *migration                        `xml:"migration,omitempty" json:"migration,omitempty"`                                               // ZZmaxLength=0
	Phenologies                      *phenologies                      `xml:"phenologies,omitempty" json:"phenologies,omitempty"`                                           // ZZmaxLength=0
	XMLName                          xml.Name                          `xml:"ecologyAndLifeHistory,omitempty" json:"ecologyAndLifeHistory,omitempty"`
}

type ecologyAndLifeHistoryAuthors struct {
	AttrdisplayValue string   `xml:"displayValue,attr"  json:",omitempty"` // maxLength=17
	XMLName          xml.Name `xml:"ecologyAndLifeHistoryAuthors,omitempty" json:"ecologyAndLifeHistoryAuthors,omitempty"`
}

type ecologyAndLifeHistoryDescription struct {
	ShortGeneralDescription *shortGeneralDescription `xml:"shortGeneralDescription,omitempty" json:"shortGeneralDescription,omitempty"` // ZZmaxLength=0
	XMLName                 xml.Name                 `xml:"ecologyAndLifeHistoryDescription,omitempty" json:"ecologyAndLifeHistoryDescription,omitempty"`
}

type ecologyAndLifeHistoryEditionDate struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=10
	XMLName xml.Name `xml:"ecologyAndLifeHistoryEditionDate,omitempty" json:"ecologyAndLifeHistoryEditionDate,omitempty"`
}

type endemism struct {
	Code        *code        `xml:"code,omitempty" json:"code,omitempty"`               // ZZmaxLength=0
	Description *description `xml:"description,omitempty" json:"description,omitempty"` // ZZmaxLength=0
	XMLName     xml.Name     `xml:"endemism,omitempty" json:"endemism,omitempty"`
}

type estimatedNumberOfOccurrences struct {
	Code        *code        `xml:"code,omitempty" json:"code,omitempty"`               // ZZmaxLength=0
	Comments    *comments    `xml:"comments,omitempty" json:"comments,omitempty"`       // ZZmaxLength=0
	Description *description `xml:"description,omitempty" json:"description,omitempty"` // ZZmaxLength=0
	SearchValue *searchValue `xml:"searchValue,omitempty" json:"searchValue,omitempty"` // ZZmaxLength=0
	XMLName     xml.Name     `xml:"estimatedNumberOfOccurrences,omitempty" json:"estimatedNumberOfOccurrences,omitempty"`
}

type family struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=11
	XMLName xml.Name `xml:"family,omitempty" json:"family,omitempty"`
}

type foodComments struct {
	I       []*ii    `xml:"i,omitempty" json:"i,omitempty"` // ZZmaxLength=0
	Text    string   `xml:",chardata" json:",omitempty"`    // maxLength=101
	XMLName xml.Name `xml:"foodComments,omitempty" json:"foodComments,omitempty"`
}

type foodHabits struct {
	FoodComments *foodComments `xml:"foodComments,omitempty" json:"foodComments,omitempty"` // ZZmaxLength=0
	XMLName      xml.Name      `xml:"foodHabits,omitempty" json:"foodHabits,omitempty"`
}

type formalTaxonomy struct {
	Class             *class             `xml:"class,omitempty" json:"class,omitempty"`                         // ZZmaxLength=0
	Family            *family            `xml:"family,omitempty" json:"family,omitempty"`                       // ZZmaxLength=0
	Genus             *genus             `xml:"genus,omitempty" json:"genus,omitempty"`                         // ZZmaxLength=0
	GenusSize         *genusSize         `xml:"genusSize,omitempty" json:"genusSize,omitempty"`                 // ZZmaxLength=0
	Kingdom           *kingdom           `xml:"kingdom,omitempty" json:"kingdom,omitempty"`                     // ZZmaxLength=0
	Order             *order             `xml:"order,omitempty" json:"order,omitempty"`                         // ZZmaxLength=0
	Phylum            *phylum            `xml:"phylum,omitempty" json:"phylum,omitempty"`                       // ZZmaxLength=0
	TaxonomicComments *taxonomicComments `xml:"taxonomicComments,omitempty" json:"taxonomicComments,omitempty"` // ZZmaxLength=0
	XMLName           xml.Name           `xml:"formalTaxonomy,omitempty" json:"formalTaxonomy,omitempty"`
}

type formattedFullCitation struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=314
	XMLName xml.Name `xml:"formattedFullCitation,omitempty" json:"formattedFullCitation,omitempty"`
}

type formattedName struct {
	I       []*ii    `xml:"i,omitempty" json:"i,omitempty"` // ZZmaxLength=0
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"formattedName,omitempty" json:"formattedName,omitempty"`
}

type genus struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=9
	XMLName xml.Name `xml:"genus,omitempty" json:"genus,omitempty"`
}

type genusSize struct {
	Code        *code        `xml:"code,omitempty" json:"code,omitempty"`               // ZZmaxLength=0
	Description *description `xml:"description,omitempty" json:"description,omitempty"` // ZZmaxLength=0
	XMLName     xml.Name     `xml:"genusSize,omitempty" json:"genusSize,omitempty"`
}

type globalAbundance struct {
	Code        *code        `xml:"code,omitempty" json:"code,omitempty"`               // ZZmaxLength=0
	Comments    *comments    `xml:"comments,omitempty" json:"comments,omitempty"`       // ZZmaxLength=0
	Description *description `xml:"description,omitempty" json:"description,omitempty"` // ZZmaxLength=0
	SearchValue *searchValue `xml:"searchValue,omitempty" json:"searchValue,omitempty"` // ZZmaxLength=0
	XMLName     xml.Name     `xml:"globalAbundance,omitempty" json:"globalAbundance,omitempty"`
}

type globalInventoryNeeds struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=113
	XMLName xml.Name `xml:"globalInventoryNeeds,omitempty" json:"globalInventoryNeeds,omitempty"`
}

type globalProtection struct {
	Code        *code        `xml:"code,omitempty" json:"code,omitempty"`               // ZZmaxLength=0
	Comments    *comments    `xml:"comments,omitempty" json:"comments,omitempty"`       // ZZmaxLength=0
	Description *description `xml:"description,omitempty" json:"description,omitempty"` // ZZmaxLength=0
	Needs       *needs       `xml:"needs,omitempty" json:"needs,omitempty"`             // ZZmaxLength=0
	SearchValue *searchValue `xml:"searchValue,omitempty" json:"searchValue,omitempty"` // ZZmaxLength=0
	XMLName     xml.Name     `xml:"globalProtection,omitempty" json:"globalProtection,omitempty"`
}

type globalRange struct {
	Code        *code        `xml:"code,omitempty" json:"code,omitempty"`               // ZZmaxLength=0
	Comments    *comments    `xml:"comments,omitempty" json:"comments,omitempty"`       // ZZmaxLength=0
	Description *description `xml:"description,omitempty" json:"description,omitempty"` // ZZmaxLength=0
	SearchValue *searchValue `xml:"searchValue,omitempty" json:"searchValue,omitempty"` // ZZmaxLength=0
	XMLName     xml.Name     `xml:"globalRange,omitempty" json:"globalRange,omitempty"`
}

type globalShortTermTrend struct {
	Code        *code        `xml:"code,omitempty" json:"code,omitempty"`               // ZZmaxLength=0
	Comments    *comments    `xml:"comments,omitempty" json:"comments,omitempty"`       // ZZmaxLength=0
	Description *description `xml:"description,omitempty" json:"description,omitempty"` // ZZmaxLength=0
	SearchValue *searchValue `xml:"searchValue,omitempty" json:"searchValue,omitempty"` // ZZmaxLength=0
	XMLName     xml.Name     `xml:"globalShortTermTrend,omitempty" json:"globalShortTermTrend,omitempty"`
}

type globalSpecies struct {
	AttrspeciesCode        string                  `xml:"speciesCode,attr"  json:",omitempty"`                                      // maxLength=10
	Attruid                string                  `xml:"uid,attr"  json:",omitempty"`                                              // maxLength=23
	Classification         *classification         `xml:"classification,omitempty" json:"classification,omitempty"`                 // ZZmaxLength=0
	ConservationStatus     *conservationStatus     `xml:"conservationStatus,omitempty" json:"conservationStatus,omitempty"`         // ZZmaxLength=0
	Distribution           *distribution           `xml:"distribution,omitempty" json:"distribution,omitempty"`                     // ZZmaxLength=0
	EcologyAndLifeHistory  *ecologyAndLifeHistory  `xml:"ecologyAndLifeHistory,omitempty" json:"ecologyAndLifeHistory,omitempty"`   // ZZmaxLength=0
	License                *license                `xml:"license,omitempty" json:"license,omitempty"`                               // ZZmaxLength=0
	ManagementSummary      *managementSummary      `xml:"managementSummary,omitempty" json:"managementSummary,omitempty"`           // ZZmaxLength=0
	NatureServeExplorerURI *natureServeExplorerURI `xml:"natureServeExplorerURI,omitempty" json:"natureServeExplorerURI,omitempty"` // ZZmaxLength=0
	PopulationOccurrence   *populationOccurrence   `xml:"populationOccurrence,omitempty" json:"populationOccurrence,omitempty"`     // ZZmaxLength=0
	References             *references             `xml:"references,omitempty" json:"references,omitempty"`                         // ZZmaxLength=0
	XMLName                xml.Name                `xml:"globalSpecies,omitempty" json:"globalSpecies,omitempty"`
}

func (Ω *globalSpecies) Synonyms() ([]*taxonScientificName, error) {
	names := Ω.Classification.Names

	res := []*taxonScientificName{}

	if names.Synonyms != nil && len(names.Synonyms.SynonymName) > 0 {
		for _, synonym := range names.Synonyms.SynonymName {
			sn, err := synonym.asTaxonScientificName()
			if err != nil {
				return nil, err
			}
			res = append(res, sn)
		}
	}
	return res, nil
}

type globalStatus struct {
	ConservationStatusFactors *conservationStatusFactors `xml:"conservationStatusFactors,omitempty" json:"conservationStatusFactors,omitempty"` // ZZmaxLength=0
	NationalStatuses          *nationalStatuses          `xml:"nationalStatuses,omitempty" json:"nationalStatuses,omitempty"`                   // ZZmaxLength=0
	Rank                      *rank                      `xml:"rank,omitempty" json:"rank,omitempty"`                                           // ZZmaxLength=0
	Reasons                   *reasons                   `xml:"reasons,omitempty" json:"reasons,omitempty"`                                     // ZZmaxLength=0
	RoundedRank               *roundedRank               `xml:"roundedRank,omitempty" json:"roundedRank,omitempty"`                             // ZZmaxLength=0
	StatusLastChanged         *statusLastChanged         `xml:"statusLastChanged,omitempty" json:"statusLastChanged,omitempty"`                 // ZZmaxLength=0
	StatusLastReviewed        *statusLastReviewed        `xml:"statusLastReviewed,omitempty" json:"statusLastReviewed,omitempty"`               // ZZmaxLength=0
	XMLName                   xml.Name                   `xml:"globalStatus,omitempty" json:"globalStatus,omitempty"`
}

type habitat struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=20
	XMLName xml.Name `xml:"habitat,omitempty" json:"habitat,omitempty"`
}

type habitatComments struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=266
	XMLName xml.Name `xml:"habitatComments,omitempty" json:"habitatComments,omitempty"`
}

type habitats struct {
	HabitatComments     *habitatComments     `xml:"habitatComments,omitempty" json:"habitatComments,omitempty"`         // ZZmaxLength=0
	PalustrineHabitats  *palustrineHabitats  `xml:"palustrineHabitats,omitempty" json:"palustrineHabitats,omitempty"`   // ZZmaxLength=0
	TerrestrialHabitats *terrestrialHabitats `xml:"terrestrialHabitats,omitempty" json:"terrestrialHabitats,omitempty"` // ZZmaxLength=0
	XMLName             xml.Name             `xml:"habitats,omitempty" json:"habitats,omitempty"`
}

type inferredExtentjustification struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=464
	XMLName xml.Name `xml:"inferredExtentjustification,omitempty" json:"inferredExtentjustification,omitempty"`
}

type inferredMinimumExtentOfHabitatUse struct {
	Distance                    *distance                    `xml:"distance,omitempty" json:"distance,omitempty"`                                       // ZZmaxLength=0
	InferredExtentjustification *inferredExtentjustification `xml:"inferredExtentjustification,omitempty" json:"inferredExtentjustification,omitempty"` // ZZmaxLength=0
	XMLName                     xml.Name                     `xml:"inferredMinimumExtentOfHabitatUse,omitempty" json:"inferredMinimumExtentOfHabitatUse,omitempty"`
}

type informalTaxonomy struct {
	InformalTaxonomyLevel1Name *informalTaxonomyLevel1Name `xml:"informalTaxonomyLevel1Name,omitempty" json:"informalTaxonomyLevel1Name,omitempty"` // ZZmaxLength=0
	InformalTaxonomyLevel2Name *informalTaxonomyLevel2Name `xml:"informalTaxonomyLevel2Name,omitempty" json:"informalTaxonomyLevel2Name,omitempty"` // ZZmaxLength=0
	InformalTaxonomyLevel3Name *informalTaxonomyLevel3Name `xml:"informalTaxonomyLevel3Name,omitempty" json:"informalTaxonomyLevel3Name,omitempty"` // ZZmaxLength=0
	InformalTaxonomyLevel4Name *informalTaxonomyLevel4Name `xml:"informalTaxonomyLevel4Name,omitempty" json:"informalTaxonomyLevel4Name,omitempty"` // ZZmaxLength=0
	XMLName                    xml.Name                    `xml:"informalTaxonomy,omitempty" json:"informalTaxonomy,omitempty"`
}

type informalTaxonomyLevel1Name struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=22
	XMLName xml.Name `xml:"informalTaxonomyLevel1Name,omitempty" json:"informalTaxonomyLevel1Name,omitempty"`
}

type informalTaxonomyLevel2Name struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=7
	XMLName xml.Name `xml:"informalTaxonomyLevel2Name,omitempty" json:"informalTaxonomyLevel2Name,omitempty"`
}

type informalTaxonomyLevel3Name struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=21
	XMLName xml.Name `xml:"informalTaxonomyLevel3Name,omitempty" json:"informalTaxonomyLevel3Name,omitempty"`
}

type informalTaxonomyLevel4Name struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=24
	XMLName xml.Name `xml:"informalTaxonomyLevel4Name,omitempty" json:"informalTaxonomyLevel4Name,omitempty"`
}

type intrinsicVulnerability struct {
	Comments *comments `xml:"comments,omitempty" json:"comments,omitempty"` // ZZmaxLength=0
	XMLName  xml.Name  `xml:"intrinsicVulnerability,omitempty" json:"intrinsicVulnerability,omitempty"`
}

type kingdom struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=8
	XMLName xml.Name `xml:"kingdom,omitempty" json:"kingdom,omitempty"`
}

type license struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=144
	XMLName xml.Name `xml:"license,omitempty" json:"license,omitempty"`
}

type locallyMigrant struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=5
	XMLName xml.Name `xml:"locallyMigrant,omitempty" json:"locallyMigrant,omitempty"`
}

type longDistanceMigrant struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=5
	XMLName xml.Name `xml:"longDistanceMigrant,omitempty" json:"longDistanceMigrant,omitempty"`
}

type managementSummary struct {
	BiologicalResearchNeeds *biologicalResearchNeeds `xml:"biologicalResearchNeeds,omitempty" json:"biologicalResearchNeeds,omitempty"` // ZZmaxLength=0
	XMLName                 xml.Name                 `xml:"managementSummary,omitempty" json:"managementSummary,omitempty"`
}

type mappingGuidance struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=889
	XMLName xml.Name `xml:"mappingGuidance,omitempty" json:"mappingGuidance,omitempty"`
}

type maximumLastObservedYear struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=4
	XMLName xml.Name `xml:"maximumLastObservedYear,omitempty" json:"maximumLastObservedYear,omitempty"`
}

type migration struct {
	LocallyMigrant      *locallyMigrant      `xml:"locallyMigrant,omitempty" json:"locallyMigrant,omitempty"`           // ZZmaxLength=0
	LongDistanceMigrant *longDistanceMigrant `xml:"longDistanceMigrant,omitempty" json:"longDistanceMigrant,omitempty"` // ZZmaxLength=0
	NonMigrant          *nonMigrant          `xml:"nonMigrant,omitempty" json:"nonMigrant,omitempty"`                   // ZZmaxLength=0
	XMLName             xml.Name             `xml:"migration,omitempty" json:"migration,omitempty"`
}

type minimumCriteriaForOccurrence struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=304
	XMLName xml.Name `xml:"minimumCriteriaForOccurrence,omitempty" json:"minimumCriteriaForOccurrence,omitempty"`
}

type nameUsedInConceptReference struct {
	FormattedName   *formattedName   `xml:"formattedName,omitempty" json:"formattedName,omitempty"`     // ZZmaxLength=0
	UnformattedName *unformattedName `xml:"unformattedName,omitempty" json:"unformattedName,omitempty"` // ZZmaxLength=0
	XMLName         xml.Name         `xml:"nameUsedInConceptReference,omitempty" json:"nameUsedInConceptReference,omitempty"`
}

func (Ω *nameUsedInConceptReference) asFormattedName() string {
	if Ω == nil {
		return ""
	}
	if Ω.UnformattedName != nil {
		return Ω.UnformattedName.Text
	}
	if Ω.FormattedName == nil {
		return ""
	}
	if len(Ω.FormattedName.I) > 0 {
		return Ω.FormattedName.I[0].Text
	}
	return sanitize.HTML(Ω.FormattedName.Text)
}

type names struct {
	NatureServePrimaryGlobalCommonName *natureServePrimaryGlobalCommonName `xml:"natureServePrimaryGlobalCommonName,omitempty" json:"natureServePrimaryGlobalCommonName,omitempty"` // ZZmaxLength=0
	OtherGlobalCommonNames             *otherGlobalCommonNames             `xml:"otherGlobalCommonNames,omitempty" json:"otherGlobalCommonNames,omitempty"`                         // ZZmaxLength=0
	ScientificName                     *scientificName                     `xml:"scientificName,omitempty" json:"scientificName,omitempty"`                                         // ZZmaxLength=0
	Synonyms                           *synonyms                           `xml:"synonyms,omitempty" json:"synonyms,omitempty"`                                                     // ZZmaxLength=0
	XMLName                            xml.Name                            `xml:"names,omitempty" json:"names,omitempty"`
}

func (Ω *names) taxonSynonyms() ([]*taxonScientificName, error) {
	if Ω == nil || Ω.Synonyms == nil || Ω.Synonyms.SynonymName == nil {
		return nil, nil
	}
	res := []*taxonScientificName{}
	for _, synonymName := range Ω.Synonyms.SynonymName {
		sn, err := synonymName.asTaxonScientificName()
		if err != nil {
			return nil, err
		}
		if sn != nil {
			res = append(res, sn)
		}
	}
	return res, nil
}

type nation struct {
	AttrnationCode        string                 `xml:"nationCode,attr"  json:",omitempty"`                                     // maxLength=2
	AttrnationName        string                 `xml:"nationName,attr"  json:",omitempty"`                                     // maxLength=13
	NationalDistributions *nationalDistributions `xml:"nationalDistributions,omitempty" json:"nationalDistributions,omitempty"` // ZZmaxLength=0
	Subnations            *subnations            `xml:"subnations,omitempty" json:"subnations,omitempty"`                       // ZZmaxLength=0
	XMLName               xml.Name               `xml:"nation,omitempty" json:"nation,omitempty"`
}

type nationalDistribution struct {
	CurrentPresenceAbsence *currentPresenceAbsence `xml:"currentPresenceAbsence,omitempty" json:"currentPresenceAbsence,omitempty"` // ZZmaxLength=0
	DistributionConfidence *distributionConfidence `xml:"distributionConfidence,omitempty" json:"distributionConfidence,omitempty"` // ZZmaxLength=0
	Origin                 *origin                 `xml:"origin,omitempty" json:"origin,omitempty"`                                 // ZZmaxLength=0
	Population             *population             `xml:"population,omitempty" json:"population,omitempty"`                         // ZZmaxLength=0
	Regularity             *regularity             `xml:"regularity,omitempty" json:"regularity,omitempty"`                         // ZZmaxLength=0
	XMLName                xml.Name                `xml:"nationalDistribution,omitempty" json:"nationalDistribution,omitempty"`
}

type nationalDistributions struct {
	NationalDistribution *nationalDistribution `xml:"nationalDistribution,omitempty" json:"nationalDistribution,omitempty"` // ZZmaxLength=0
	XMLName              xml.Name              `xml:"nationalDistributions,omitempty" json:"nationalDistributions,omitempty"`
}

type nationalStatus struct {
	AttrnationCode      string               `xml:"nationCode,attr"  json:",omitempty"`                                 // maxLength=2
	AttrnationName      string               `xml:"nationName,attr"  json:",omitempty"`                                 // maxLength=13
	Rank                *rank                `xml:"rank,omitempty" json:"rank,omitempty"`                               // ZZmaxLength=0
	RoundedRank         *roundedRank         `xml:"roundedRank,omitempty" json:"roundedRank,omitempty"`                 // ZZmaxLength=0
	StatusLastReviewed  *statusLastReviewed  `xml:"statusLastReviewed,omitempty" json:"statusLastReviewed,omitempty"`   // ZZmaxLength=0
	SubnationalStatuses *subnationalStatuses `xml:"subnationalStatuses,omitempty" json:"subnationalStatuses,omitempty"` // ZZmaxLength=0
	XMLName             xml.Name             `xml:"nationalStatus,omitempty" json:"nationalStatus,omitempty"`
}

type nationalStatuses struct {
	NationalStatus []*nationalStatus `xml:"nationalStatus,omitempty" json:"nationalStatus,omitempty"` // ZZmaxLength=0
	XMLName        xml.Name          `xml:"nationalStatuses,omitempty" json:"nationalStatuses,omitempty"`
}

type nations struct {
	Nation  []*nation `xml:"nation,omitempty" json:"nation,omitempty"` // ZZmaxLength=0
	XMLName xml.Name  `xml:"nations,omitempty" json:"nations,omitempty"`
}

type natureServePrimaryGlobalCommonName struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=20
	XMLName xml.Name `xml:"natureServePrimaryGlobalCommonName,omitempty" json:"natureServePrimaryGlobalCommonName,omitempty"`
}

type natureServeStatus struct {
	GlobalStatus *globalStatus `xml:"globalStatus,omitempty" json:"globalStatus,omitempty"` // ZZmaxLength=0
	XMLName      xml.Name      `xml:"natureServeStatus,omitempty" json:"natureServeStatus,omitempty"`
}

type needs struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=120
	XMLName xml.Name `xml:"needs,omitempty" json:"needs,omitempty"`
}

type nomenclaturalAuthor struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=31
	XMLName xml.Name `xml:"nomenclaturalAuthor,omitempty" json:"nomenclaturalAuthor,omitempty"`
}

type nonMigrant struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=5
	XMLName xml.Name `xml:"nonMigrant,omitempty" json:"nonMigrant,omitempty"`
}

type notes struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=967
	XMLName xml.Name `xml:"notes,omitempty" json:"notes,omitempty"`
}

type occurrenceCounties struct {
	OccurrenceCounty []*occurrenceCounty `xml:"occurrenceCounty,omitempty" json:"occurrenceCounty,omitempty"` // ZZmaxLength=0
	XMLName          xml.Name            `xml:"occurrenceCounties,omitempty" json:"occurrenceCounties,omitempty"`
}

type occurrenceCounty struct {
	CountyCode              *countyCode              `xml:"countyCode,omitempty" json:"countyCode,omitempty"`                           // ZZmaxLength=0
	CountyName              *countyName              `xml:"countyName,omitempty" json:"countyName,omitempty"`                           // ZZmaxLength=0
	MaximumLastObservedYear *maximumLastObservedYear `xml:"maximumLastObservedYear,omitempty" json:"maximumLastObservedYear,omitempty"` // ZZmaxLength=0
	SpeciesOccurrenceCount  *speciesOccurrenceCount  `xml:"speciesOccurrenceCount,omitempty" json:"speciesOccurrenceCount,omitempty"`   // ZZmaxLength=0
	XMLName                 xml.Name                 `xml:"occurrenceCounty,omitempty" json:"occurrenceCounty,omitempty"`
}

type occurrenceNation struct {
	Attrcode         string            `xml:"code,attr"  json:",omitempty"`                                 // maxLength=2
	OccurrenceStates *occurrenceStates `xml:"occurrenceStates,omitempty" json:"occurrenceStates,omitempty"` // ZZmaxLength=0
	XMLName          xml.Name          `xml:"occurrenceNation,omitempty" json:"occurrenceNation,omitempty"`
}

type occurrenceNations struct {
	OccurrenceNation *occurrenceNation `xml:"occurrenceNation,omitempty" json:"occurrenceNation,omitempty"` // ZZmaxLength=0
	XMLName          xml.Name          `xml:"occurrenceNations,omitempty" json:"occurrenceNations,omitempty"`
}

type occurrenceState struct {
	Attrcode           string              `xml:"code,attr"  json:",omitempty"`                                     // maxLength=2
	OccurrenceCounties *occurrenceCounties `xml:"occurrenceCounties,omitempty" json:"occurrenceCounties,omitempty"` // ZZmaxLength=0
	XMLName            xml.Name            `xml:"occurrenceState,omitempty" json:"occurrenceState,omitempty"`
}

type occurrenceStates struct {
	OccurrenceState []*occurrenceState `xml:"occurrenceState,omitempty" json:"occurrenceState,omitempty"` // ZZmaxLength=0
	XMLName         xml.Name           `xml:"occurrenceStates,omitempty" json:"occurrenceStates,omitempty"`
}

type order struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=11
	XMLName xml.Name `xml:"order,omitempty" json:"order,omitempty"`
}

type origin struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=20
	XMLName xml.Name `xml:"origin,omitempty" json:"origin,omitempty"`
}

type otherConsiderations struct {
	XMLName xml.Name `xml:"otherConsiderations,omitempty" json:"otherConsiderations,omitempty"`
}

type otherGlobalCommonNames struct {
	CommonName []*commonName `xml:"commonName,omitempty" json:"commonName,omitempty"` // ZZmaxLength=0
	XMLName    xml.Name      `xml:"otherGlobalCommonNames,omitempty" json:"otherGlobalCommonNames,omitempty"`
}

type otherStatuses struct {
	Status  []*status `xml:"status,omitempty" json:"status,omitempty"` // ZZmaxLength=0
	XMLName xml.Name  `xml:"otherStatuses,omitempty" json:"otherStatuses,omitempty"`
}

type palustrineHabitats struct {
	Habitat []*habitat `xml:"habitat,omitempty" json:"habitat,omitempty"` // ZZmaxLength=0
	XMLName xml.Name   `xml:"palustrineHabitats,omitempty" json:"palustrineHabitats,omitempty"`
}

type phenologies struct {
	PhenologyComments *phenologyComments `xml:"phenologyComments,omitempty" json:"phenologyComments,omitempty"` // ZZmaxLength=0
	XMLName           xml.Name           `xml:"phenologies,omitempty" json:"phenologies,omitempty"`
}

type phenologyComments struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=256
	XMLName xml.Name `xml:"phenologyComments,omitempty" json:"phenologyComments,omitempty"`
}

type phylum struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=11
	XMLName xml.Name `xml:"phylum,omitempty" json:"phylum,omitempty"`
}

type population struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=10
	XMLName xml.Name `xml:"population,omitempty" json:"population,omitempty"`
}

type populationOccurrence struct {
	Delineations *delineations `xml:"delineations,omitempty" json:"delineations,omitempty"` // ZZmaxLength=0
	Viabilities  *viabilities  `xml:"viabilities,omitempty" json:"viabilities,omitempty"`   // ZZmaxLength=0
	XMLName      xml.Name      `xml:"populationOccurrence,omitempty" json:"populationOccurrence,omitempty"`
}

type rank struct {
	Code    *code    `xml:"code,omitempty" json:"code,omitempty"` // ZZmaxLength=0
	XMLName xml.Name `xml:"rank,omitempty" json:"rank,omitempty"`
}

type reasons struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=168
	XMLName xml.Name `xml:"reasons,omitempty" json:"reasons,omitempty"`
}

type references struct {
	Citation []*citation `xml:"citation,omitempty" json:"citation,omitempty"` // ZZmaxLength=0
	XMLName  xml.Name    `xml:"references,omitempty" json:"references,omitempty"`
}

type regularity struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=19
	XMLName xml.Name `xml:"regularity,omitempty" json:"regularity,omitempty"`
}

type roundedRank struct {
	Code        *code        `xml:"code,omitempty" json:"code,omitempty"`               // ZZmaxLength=0
	Description *description `xml:"description,omitempty" json:"description,omitempty"` // ZZmaxLength=0
	XMLName     xml.Name     `xml:"roundedRank,omitempty" json:"roundedRank,omitempty"`
}

type scientificName struct {
	ConceptReference    *conceptReference    `xml:"conceptReference,omitempty" json:"conceptReference,omitempty"`       // ZZmaxLength=0
	FormattedName       *formattedName       `xml:"formattedName,omitempty" json:"formattedName,omitempty"`             // ZZmaxLength=0
	NomenclaturalAuthor *nomenclaturalAuthor `xml:"nomenclaturalAuthor,omitempty" json:"nomenclaturalAuthor,omitempty"` // ZZmaxLength=0
	UnformattedName     *unformattedName     `xml:"unformattedName,omitempty" json:"unformattedName,omitempty"`         // ZZmaxLength=0
	XMLName             xml.Name             `xml:"scientificName,omitempty" json:"scientificName,omitempty"`
}

func (Ω *scientificName) asTaxonScientificName() (*taxonScientificName, error) {

	name := Ω.CanonicalName()
	if name == "" {
		return nil, nil
	}

	txn := taxonScientificName{
		Name: name,
	}

	if Ω.NomenclaturalAuthor != nil {
		txn.Author = Ω.NomenclaturalAuthor.Text
	}

	if Ω.ConceptReference == nil {
		return &txn, nil
	}

	txn.ConceptReferenceCode = Ω.ConceptReference.Attrcode
	if Ω.ConceptReference.ClassificationStatus != nil {
		txn.ConceptReferenceClassificationStatus = Ω.ConceptReference.ClassificationStatus.Text
	}
	if Ω.ConceptReference.FormattedFullCitation != nil {
		txn.ConceptReferenceFullCitation = Ω.ConceptReference.FormattedFullCitation.Text
	}
	txn.ConceptReferenceNameUsed = Ω.ConceptReference.NameUsedInConceptReference.asFormattedName()
	txn.ConceptReferenceCode = Ω.ConceptReference.Attrcode

	return &txn, nil
}

func (Ω *scientificName) CanonicalName() string {
	name := ""
	if Ω.UnformattedName != nil {
		name = Ω.UnformattedName.Text
	} else if Ω.FormattedName != nil {
		if len(Ω.FormattedName.I) > 0 {
			name = Ω.FormattedName.I[0].Text
		} else {
			name = sanitize.HTML(Ω.FormattedName.Text)
		}
	}
	return name
}

type synonymName struct {
	FormattedName       *formattedName       `xml:"formattedName,omitempty" json:"formattedName,omitempty"`             // ZZmaxLength=0
	NomenclaturalAuthor *nomenclaturalAuthor `xml:"nomenclaturalAuthor,omitempty" json:"nomenclaturalAuthor,omitempty"` // ZZmaxLength=0
	UnformattedName     *unformattedName     `xml:"unformattedName,omitempty" json:"unformattedName,omitempty"`         // ZZmaxLength=0
	XMLName             xml.Name             `xml:"synonymName,omitempty" json:"synonymName,omitempty"`
}

func (Ω *synonymName) asTaxonScientificName() (*taxonScientificName, error) {
	name := ""
	if Ω.UnformattedName != nil {
		name = Ω.UnformattedName.Text
	} else if Ω.FormattedName != nil {
		if len(Ω.FormattedName.I) > 0 {
			name = Ω.FormattedName.I[0].Text
		} else {
			name = sanitize.HTML(Ω.FormattedName.Text)
		}
	} else {
		return nil, nil
	}

	txnName := taxonScientificName{
		Name: name,
	}

	if Ω.NomenclaturalAuthor != nil {
		txnName.Author = Ω.NomenclaturalAuthor.Text
	}

	return &txnName, nil
}

type searchValue struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=3
	XMLName xml.Name `xml:"searchValue,omitempty" json:"searchValue,omitempty"`
}

type separation struct {
	AlternateSeparationProcedure *alternateSeparationProcedure `xml:"alternateSeparationProcedure,omitempty" json:"alternateSeparationProcedure,omitempty"` // ZZmaxLength=0
	Barriers                     *barriers                     `xml:"barriers,omitempty" json:"barriers,omitempty"`                                         // ZZmaxLength=0
	DistanceForSuitableHabitat   *distanceForSuitableHabitat   `xml:"distanceForSuitableHabitat,omitempty" json:"distanceForSuitableHabitat,omitempty"`     // ZZmaxLength=0
	DistanceForUnsuitableHabitat *distanceForUnsuitableHabitat `xml:"distanceForUnsuitableHabitat,omitempty" json:"distanceForUnsuitableHabitat,omitempty"` // ZZmaxLength=0
	SeparationJustification      *separationJustification      `xml:"separationJustification,omitempty" json:"separationJustification,omitempty"`           // ZZmaxLength=0
	XMLName                      xml.Name                      `xml:"separation,omitempty" json:"separation,omitempty"`
}

type separationJustification struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=1191
	XMLName xml.Name `xml:"separationJustification,omitempty" json:"separationJustification,omitempty"`
}

type shortGeneralDescription struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=38
	XMLName xml.Name `xml:"shortGeneralDescription,omitempty" json:"shortGeneralDescription,omitempty"`
}

type speciesOccurrenceCount struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=2
	XMLName xml.Name `xml:"speciesOccurrenceCount,omitempty" json:"speciesOccurrenceCount,omitempty"`
}

type status struct {
	Attrname          string             `xml:"name,attr"  json:",omitempty"`                                   // maxLength=14
	StatusDate        *statusDate        `xml:"statusDate,omitempty" json:"statusDate,omitempty"`               // ZZmaxLength=0
	StatusDescription *statusDescription `xml:"statusDescription,omitempty" json:"statusDescription,omitempty"` // ZZmaxLength=0
	StatusValue       *statusValue       `xml:"statusValue,omitempty" json:"statusValue,omitempty"`             // ZZmaxLength=0
	XMLName           xml.Name           `xml:"status,omitempty" json:"status,omitempty"`
}

type statusDate struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=10
	XMLName xml.Name `xml:"statusDate,omitempty" json:"statusDate,omitempty"`
}

type statusDescription struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=15
	XMLName xml.Name `xml:"statusDescription,omitempty" json:"statusDescription,omitempty"`
}

type statusFactorsEditionDate struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=10
	XMLName xml.Name `xml:"statusFactorsEditionDate,omitempty" json:"statusFactorsEditionDate,omitempty"`
}

type statusLastChanged struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=10
	XMLName xml.Name `xml:"statusLastChanged,omitempty" json:"statusLastChanged,omitempty"`
}

type statusLastReviewed struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=10
	XMLName xml.Name `xml:"statusLastReviewed,omitempty" json:"statusLastReviewed,omitempty"`
}

type statusValue struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=2
	XMLName xml.Name `xml:"statusValue,omitempty" json:"statusValue,omitempty"`
}

type subnation struct {
	AttrsubnationCode        string                    `xml:"subnationCode,attr"  json:",omitempty"`                                        // maxLength=2
	AttrsubnationName        string                    `xml:"subnationName,attr"  json:",omitempty"`                                        // maxLength=21
	SubnationalDistributions *subnationalDistributions `xml:"subnationalDistributions,omitempty" json:"subnationalDistributions,omitempty"` // ZZmaxLength=0
	XMLName                  xml.Name                  `xml:"subnation,omitempty" json:"subnation,omitempty"`
}

type subnationalDistribution struct {
	CurrentPresenceAbsence *currentPresenceAbsence `xml:"currentPresenceAbsence,omitempty" json:"currentPresenceAbsence,omitempty"` // ZZmaxLength=0
	DistributionConfidence *distributionConfidence `xml:"distributionConfidence,omitempty" json:"distributionConfidence,omitempty"` // ZZmaxLength=0
	Origin                 *origin                 `xml:"origin,omitempty" json:"origin,omitempty"`                                 // ZZmaxLength=0
	Population             *population             `xml:"population,omitempty" json:"population,omitempty"`                         // ZZmaxLength=0
	Regularity             *regularity             `xml:"regularity,omitempty" json:"regularity,omitempty"`                         // ZZmaxLength=0
	XMLName                xml.Name                `xml:"subnationalDistribution,omitempty" json:"subnationalDistribution,omitempty"`
}

type subnationalDistributions struct {
	SubnationalDistribution *subnationalDistribution `xml:"subnationalDistribution,omitempty" json:"subnationalDistribution,omitempty"` // ZZmaxLength=0
	XMLName                 xml.Name                 `xml:"subnationalDistributions,omitempty" json:"subnationalDistributions,omitempty"`
}

type subnationalStatus struct {
	AttrsubnationCode string       `xml:"subnationCode,attr"  json:",omitempty"`              // maxLength=2
	AttrsubnationName string       `xml:"subnationName,attr"  json:",omitempty"`              // maxLength=21
	Rank              *rank        `xml:"rank,omitempty" json:"rank,omitempty"`               // ZZmaxLength=0
	RoundedRank       *roundedRank `xml:"roundedRank,omitempty" json:"roundedRank,omitempty"` // ZZmaxLength=0
	XMLName           xml.Name     `xml:"subnationalStatus,omitempty" json:"subnationalStatus,omitempty"`
}

type subnationalStatuses struct {
	SubnationalStatus []*subnationalStatus `xml:"subnationalStatus,omitempty" json:"subnationalStatus,omitempty"` // ZZmaxLength=0
	XMLName           xml.Name             `xml:"subnationalStatuses,omitempty" json:"subnationalStatuses,omitempty"`
}

type subnations struct {
	Subnation []*subnation `xml:"subnation,omitempty" json:"subnation,omitempty"` // ZZmaxLength=0
	XMLName   xml.Name     `xml:"subnations,omitempty" json:"subnations,omitempty"`
}

type synonyms struct {
	SynonymName []*synonymName `xml:"synonymName,omitempty" json:"synonymName,omitempty"` // ZZmaxLength=0
	XMLName     xml.Name       `xml:"synonyms,omitempty" json:"synonyms,omitempty"`
}

type taxonomy struct {
	FormalTaxonomy   *formalTaxonomy   `xml:"formalTaxonomy,omitempty" json:"formalTaxonomy,omitempty"`     // ZZmaxLength=0
	InformalTaxonomy *informalTaxonomy `xml:"informalTaxonomy,omitempty" json:"informalTaxonomy,omitempty"` // ZZmaxLength=0
	XMLName          xml.Name          `xml:"taxonomy,omitempty" json:"taxonomy,omitempty"`
}

type terrestrialHabitats struct {
	Habitat []*habitat `xml:"habitat,omitempty" json:"habitat,omitempty"` // ZZmaxLength=0
	XMLName xml.Name   `xml:"terrestrialHabitats,omitempty" json:"terrestrialHabitats,omitempty"`
}

type threat struct {
	Comments *comments `xml:"comments,omitempty" json:"comments,omitempty"` // ZZmaxLength=0
	XMLName  xml.Name  `xml:"threat,omitempty" json:"threat,omitempty"`
}

type unformattedName struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=35
	XMLName xml.Name `xml:"unformattedName,omitempty" json:"unformattedName,omitempty"`
}

type versionDate struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=10
	XMLName xml.Name `xml:"versionDate,omitempty" json:"versionDate,omitempty"`
}

type viabilities struct {
	Viability *viability `xml:"viability,omitempty" json:"viability,omitempty"` // ZZmaxLength=0
	XMLName   xml.Name   `xml:"viabilities,omitempty" json:"viabilities,omitempty"`
}

type viability struct {
	ViabilityJustification *viabilityJustification `xml:"viabilityJustification,omitempty" json:"viabilityJustification,omitempty"` // ZZmaxLength=0
	XMLName                xml.Name                `xml:"viability,omitempty" json:"viability,omitempty"`
}

type viabilityJustification struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=201
	XMLName xml.Name `xml:"viabilityJustification,omitempty" json:"viabilityJustification,omitempty"`
}

type watershed struct {
	Attrtype                string                   `xml:"type,attr"  json:",omitempty"`                                               // maxLength=5
	MaximumLastObservedYear *maximumLastObservedYear `xml:"maximumLastObservedYear,omitempty" json:"maximumLastObservedYear,omitempty"` // ZZmaxLength=0
	SpeciesOccurrenceCount  *speciesOccurrenceCount  `xml:"speciesOccurrenceCount,omitempty" json:"speciesOccurrenceCount,omitempty"`   // ZZmaxLength=0
	WatershedCodes          *watershedCodes          `xml:"watershedCodes,omitempty" json:"watershedCodes,omitempty"`                   // ZZmaxLength=0
	WatershedName           *watershedName           `xml:"watershedName,omitempty" json:"watershedName,omitempty"`                     // ZZmaxLength=0
	XMLName                 xml.Name                 `xml:"watershed,omitempty" json:"watershed,omitempty"`
}

type watershedCode struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=8
	XMLName xml.Name `xml:"watershedCode,omitempty" json:"watershedCode,omitempty"`
}

type watershedCodes struct {
	WatershedCode *watershedCode `xml:"watershedCode,omitempty" json:"watershedCode,omitempty"` // ZZmaxLength=0
	XMLName       xml.Name       `xml:"watershedCodes,omitempty" json:"watershedCodes,omitempty"`
}

type watershedList struct {
	Watershed []*watershed `xml:"watershed,omitempty" json:"watershed,omitempty"` // ZZmaxLength=0
	XMLName   xml.Name     `xml:"watershedList,omitempty" json:"watershedList,omitempty"`
}

type watershedName struct {
	Text    string   `xml:",chardata" json:",omitempty"` // maxLength=31
	XMLName xml.Name `xml:"watershedName,omitempty" json:"watershedName,omitempty"`
}

type watersheds struct {
	WatershedList *watershedList `xml:"watershedList,omitempty" json:"watershedList,omitempty"` // ZZmaxLength=0
	XMLName       xml.Name       `xml:"watersheds,omitempty" json:"watersheds,omitempty"`
}
