package nature_serve

import "encoding/xml"

/*****
			SPECIES SEARCH REPORT
*****/

type CommonName struct {
	Attrlanguage string`xml:"language,attr"  json:",omitempty"`  // maxLength=2
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=20
	XMLName  xml.Name `xml:"commonName,omitempty" json:"commonName,omitempty"`
}

type Conditions struct {
	Match *Match `xml:"match,omitempty" json:"match,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"conditions,omitempty" json:"conditions,omitempty"`
}

type GlobalSpeciesUid struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=23
	XMLName  xml.Name `xml:"globalSpeciesUid,omitempty" json:"globalSpeciesUid,omitempty"`
}

type I struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=25
	XMLName  xml.Name `xml:"i,omitempty" json:"i,omitempty"`
}

type JurisdictionScientificName struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=35
	XMLName  xml.Name `xml:"jurisdictionScientificName,omitempty" json:"jurisdictionScientificName,omitempty"`
}

type Match struct {
	SearchTerm *SearchTerm `xml:"searchTerm,omitempty" json:"searchTerm,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"match,omitempty" json:"match,omitempty"`
}

type NatureServeExplorerURI struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=94
	XMLName  xml.Name `xml:"natureServeExplorerURI,omitempty" json:"natureServeExplorerURI,omitempty"`
}

type SearchTerm struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=10
	XMLName  xml.Name `xml:"searchTerm,omitempty" json:"searchTerm,omitempty"`
}

type SpeciesSearchReport struct {
	AttrXsiSpaceschemaLocation string`xml:"http://www.w3.org/2001/XMLSchema-instance schemaLocation,attr"  json:",omitempty"`  // maxLength=136
	AttrschemaVersion string`xml:"schemaVersion,attr"  json:",omitempty"`  // maxLength=3
	Attrxmlns string`xml:"xmlns,attr"  json:",omitempty"`  // maxLength=67
	AttrXmlnsxsi string`xml:"xmlns xsi,attr"  json:",omitempty"`  // maxLength=41
	Conditions *Conditions `xml:"conditions,omitempty" json:"conditions,omitempty"`   // ZZmaxLength=0
	SpeciesSearchResultList *SpeciesSearchResultList `xml:"speciesSearchResultList,omitempty" json:"speciesSearchResultList,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"speciesSearchReport,omitempty" json:"speciesSearchReport,omitempty"`
}

type SpeciesSearchResult struct {
	Attruid string`xml:"uid,attr"  json:",omitempty"`  // maxLength=23
	CommonName *CommonName `xml:"commonName,omitempty" json:"commonName,omitempty"`   // ZZmaxLength=0
	GlobalSpeciesUid *GlobalSpeciesUid `xml:"globalSpeciesUid,omitempty" json:"globalSpeciesUid,omitempty"`   // ZZmaxLength=0
	JurisdictionScientificName *JurisdictionScientificName `xml:"jurisdictionScientificName,omitempty" json:"jurisdictionScientificName,omitempty"`   // ZZmaxLength=0
	NatureServeExplorerURI *NatureServeExplorerURI `xml:"natureServeExplorerURI,omitempty" json:"natureServeExplorerURI,omitempty"`   // ZZmaxLength=0
	TaxonomicComments *TaxonomicComments `xml:"taxonomicComments,omitempty" json:"taxonomicComments,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"speciesSearchResult,omitempty" json:"speciesSearchResult,omitempty"`
}

type SpeciesSearchResultList struct {
	SpeciesSearchResult []*SpeciesSearchResult `xml:"speciesSearchResult,omitempty" json:"speciesSearchResult,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"speciesSearchResultList,omitempty" json:"speciesSearchResultList,omitempty"`
}

type TaxonomicComments struct {
	I []*I `xml:"i,omitempty" json:"i,omitempty"`   // ZZmaxLength=0
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=733
	XMLName  xml.Name `xml:"taxonomicComments,omitempty" json:"taxonomicComments,omitempty"`
}






/*********************************************

SPECIES COMPREHENSIVE RESULTS

**********************************************/


type GlobalSpeciesList struct {
	AttrXsiSpaceschemaLocation string`xml:"http://www.w3.org/2001/XMLSchema-instance schemaLocation,attr"  json:",omitempty"`  // maxLength=136
	AttrschemaVersion string`xml:"schemaVersion,attr"  json:",omitempty"`  // maxLength=3
	Attrxmlns string`xml:"xmlns,attr"  json:",omitempty"`  // maxLength=67
	AttrXmlnsxsi string`xml:"xmlns xsi,attr"  json:",omitempty"`  // maxLength=41
	GlobalSpecies []*GlobalSpecies `xml:"globalSpecies,omitempty" json:"globalSpecies,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"globalSpeciesList,omitempty" json:"globalSpeciesList,omitempty"`
}

type AlternateSeparationProcedure struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=528
	XMLName  xml.Name `xml:"alternateSeparationProcedure,omitempty" json:"alternateSeparationProcedure,omitempty"`
}

type Barriers struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=209
	XMLName  xml.Name `xml:"barriers,omitempty" json:"barriers,omitempty"`
}

type BiologicalResearchNeeds struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=32
	XMLName  xml.Name `xml:"biologicalResearchNeeds,omitempty" json:"biologicalResearchNeeds,omitempty"`
}

type Citation struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=314
	XMLName  xml.Name `xml:"citation,omitempty" json:"citation,omitempty"`
}

type Class struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=7
	XMLName  xml.Name `xml:"class,omitempty" json:"class,omitempty"`
}

type Classification struct {
	Names *Names `xml:"names,omitempty" json:"names,omitempty"`   // ZZmaxLength=0
	Taxonomy *Taxonomy `xml:"taxonomy,omitempty" json:"taxonomy,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"classification,omitempty" json:"classification,omitempty"`
}

type ClassificationStatus struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=8
	XMLName  xml.Name `xml:"classificationStatus,omitempty" json:"classificationStatus,omitempty"`
}

type Code struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=6
	XMLName  xml.Name `xml:"code,omitempty" json:"code,omitempty"`
}

type Comments struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=273
	XMLName  xml.Name `xml:"comments,omitempty" json:"comments,omitempty"`
}

type ConceptReference struct {
	Attrcode string`xml:"code,attr"  json:",omitempty"`  // maxLength=12
	ClassificationStatus *ClassificationStatus `xml:"classificationStatus,omitempty" json:"classificationStatus,omitempty"`   // ZZmaxLength=0
	FormattedFullCitation *FormattedFullCitation `xml:"formattedFullCitation,omitempty" json:"formattedFullCitation,omitempty"`   // ZZmaxLength=0
	NameUsedInConceptReference *NameUsedInConceptReference `xml:"nameUsedInConceptReference,omitempty" json:"nameUsedInConceptReference,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"conceptReference,omitempty" json:"conceptReference,omitempty"`
}

type ConservationStatus struct {
	NatureServeStatus *NatureServeStatus `xml:"natureServeStatus,omitempty" json:"natureServeStatus,omitempty"`   // ZZmaxLength=0
	OtherStatuses *OtherStatuses `xml:"otherStatuses,omitempty" json:"otherStatuses,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"conservationStatus,omitempty" json:"conservationStatus,omitempty"`
}

type ConservationStatusAuthors struct {
	AttrdisplayValue string`xml:"displayValue,attr"  json:",omitempty"`  // maxLength=27
	XMLName  xml.Name `xml:"conservationStatusAuthors,omitempty" json:"conservationStatusAuthors,omitempty"`
}

type ConservationStatusFactors struct {
	ConservationStatusAuthors *ConservationStatusAuthors `xml:"conservationStatusAuthors,omitempty" json:"conservationStatusAuthors,omitempty"`   // ZZmaxLength=0
	EstimatedNumberOfOccurrences *EstimatedNumberOfOccurrences `xml:"estimatedNumberOfOccurrences,omitempty" json:"estimatedNumberOfOccurrences,omitempty"`   // ZZmaxLength=0
	GlobalAbundance *GlobalAbundance `xml:"globalAbundance,omitempty" json:"globalAbundance,omitempty"`   // ZZmaxLength=0
	GlobalInventoryNeeds *GlobalInventoryNeeds `xml:"globalInventoryNeeds,omitempty" json:"globalInventoryNeeds,omitempty"`   // ZZmaxLength=0
	GlobalProtection *GlobalProtection `xml:"globalProtection,omitempty" json:"globalProtection,omitempty"`   // ZZmaxLength=0
	GlobalShortTermTrend *GlobalShortTermTrend `xml:"globalShortTermTrend,omitempty" json:"globalShortTermTrend,omitempty"`   // ZZmaxLength=0
	IntrinsicVulnerability *IntrinsicVulnerability `xml:"intrinsicVulnerability,omitempty" json:"intrinsicVulnerability,omitempty"`   // ZZmaxLength=0
	OtherConsiderations *OtherConsiderations `xml:"otherConsiderations,omitempty" json:"otherConsiderations,omitempty"`   // ZZmaxLength=0
	StatusFactorsEditionDate *StatusFactorsEditionDate `xml:"statusFactorsEditionDate,omitempty" json:"statusFactorsEditionDate,omitempty"`   // ZZmaxLength=0
	Threat *Threat `xml:"threat,omitempty" json:"threat,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"conservationStatusFactors,omitempty" json:"conservationStatusFactors,omitempty"`
}

type ConservationStatusMap struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=51
	XMLName  xml.Name `xml:"conservationStatusMap,omitempty" json:"conservationStatusMap,omitempty"`
}

type CountyCode struct {
	Attrtype string`xml:"type,attr"  json:",omitempty"`  // maxLength=4
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=5
	XMLName  xml.Name `xml:"countyCode,omitempty" json:"countyCode,omitempty"`
}

type CountyDistribution struct {
	OccurrenceNations *OccurrenceNations `xml:"occurrenceNations,omitempty" json:"occurrenceNations,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"countyDistribution,omitempty" json:"countyDistribution,omitempty"`
}

type CountyName struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=10
	XMLName  xml.Name `xml:"countyName,omitempty" json:"countyName,omitempty"`
}

type CurrentPresenceAbsence struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=7
	XMLName  xml.Name `xml:"currentPresenceAbsence,omitempty" json:"currentPresenceAbsence,omitempty"`
}

type Delineation struct {
	AttrgroupName string`xml:"groupName,attr"  json:",omitempty"`  // maxLength=30
	AttrmigratoryUseType string`xml:"migratoryUseType,attr"  json:",omitempty"`  // maxLength=14
	DelineationAuthors *DelineationAuthors `xml:"delineationAuthors,omitempty" json:"delineationAuthors,omitempty"`   // ZZmaxLength=0
	InferredMinimumExtentOfHabitatUse *InferredMinimumExtentOfHabitatUse `xml:"inferredMinimumExtentOfHabitatUse,omitempty" json:"inferredMinimumExtentOfHabitatUse,omitempty"`   // ZZmaxLength=0
	MappingGuidance *MappingGuidance `xml:"mappingGuidance,omitempty" json:"mappingGuidance,omitempty"`   // ZZmaxLength=0
	MinimumCriteriaForOccurrence *MinimumCriteriaForOccurrence `xml:"minimumCriteriaForOccurrence,omitempty" json:"minimumCriteriaForOccurrence,omitempty"`   // ZZmaxLength=0
	Notes *Notes `xml:"notes,omitempty" json:"notes,omitempty"`   // ZZmaxLength=0
	Separation *Separation `xml:"separation,omitempty" json:"separation,omitempty"`   // ZZmaxLength=0
	VersionDate *VersionDate `xml:"versionDate,omitempty" json:"versionDate,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"delineation,omitempty" json:"delineation,omitempty"`
}

type DelineationAuthors struct {
	AttrdisplayValue string`xml:"displayValue,attr"  json:",omitempty"`  // maxLength=16
	XMLName  xml.Name `xml:"delineationAuthors,omitempty" json:"delineationAuthors,omitempty"`
}

type Delineations struct {
	Delineation *Delineation `xml:"delineation,omitempty" json:"delineation,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"delineations,omitempty" json:"delineations,omitempty"`
}

type Description struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=79
	XMLName  xml.Name `xml:"description,omitempty" json:"description,omitempty"`
}

type Distance struct {
	Attrunits string`xml:"units,attr"  json:",omitempty"`  // maxLength=10
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=1
	XMLName  xml.Name `xml:"distance,omitempty" json:"distance,omitempty"`
}

type DistanceForSuitableHabitat struct {
	Attrunits string`xml:"units,attr"  json:",omitempty"`  // maxLength=10
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=2
	XMLName  xml.Name `xml:"distanceForSuitableHabitat,omitempty" json:"distanceForSuitableHabitat,omitempty"`
}

type DistanceForUnsuitableHabitat struct {
	Attrunits string`xml:"units,attr"  json:",omitempty"`  // maxLength=10
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=1
	XMLName  xml.Name `xml:"distanceForUnsuitableHabitat,omitempty" json:"distanceForUnsuitableHabitat,omitempty"`
}

type Distribution struct {
	AttrUSAndCanadianDistributionComplete string`xml:"USAndCanadianDistributionComplete,attr"  json:",omitempty"`  // maxLength=5
	ConservationStatusMap *ConservationStatusMap `xml:"conservationStatusMap,omitempty" json:"conservationStatusMap,omitempty"`   // ZZmaxLength=0
	CountyDistribution *CountyDistribution `xml:"countyDistribution,omitempty" json:"countyDistribution,omitempty"`   // ZZmaxLength=0
	Endemism *Endemism `xml:"endemism,omitempty" json:"endemism,omitempty"`   // ZZmaxLength=0
	GlobalRange *GlobalRange `xml:"globalRange,omitempty" json:"globalRange,omitempty"`   // ZZmaxLength=0
	Nations *Nations `xml:"nations,omitempty" json:"nations,omitempty"`   // ZZmaxLength=0
	Watersheds *Watersheds `xml:"watersheds,omitempty" json:"watersheds,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"distribution,omitempty" json:"distribution,omitempty"`
}

type DistributionConfidence struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=9
	XMLName  xml.Name `xml:"distributionConfidence,omitempty" json:"distributionConfidence,omitempty"`
}

type Durations struct {
	XMLName  xml.Name `xml:"durations,omitempty" json:"durations,omitempty"`
}

type EcologyAndLifeHistory struct {
	Durations *Durations `xml:"durations,omitempty" json:"durations,omitempty"`   // ZZmaxLength=0
	EcologyAndLifeHistoryAuthors *EcologyAndLifeHistoryAuthors `xml:"ecologyAndLifeHistoryAuthors,omitempty" json:"ecologyAndLifeHistoryAuthors,omitempty"`   // ZZmaxLength=0
	EcologyAndLifeHistoryDescription *EcologyAndLifeHistoryDescription `xml:"ecologyAndLifeHistoryDescription,omitempty" json:"ecologyAndLifeHistoryDescription,omitempty"`   // ZZmaxLength=0
	EcologyAndLifeHistoryEditionDate *EcologyAndLifeHistoryEditionDate `xml:"ecologyAndLifeHistoryEditionDate,omitempty" json:"ecologyAndLifeHistoryEditionDate,omitempty"`   // ZZmaxLength=0
	FoodHabits *FoodHabits `xml:"foodHabits,omitempty" json:"foodHabits,omitempty"`   // ZZmaxLength=0
	Habitats *Habitats `xml:"habitats,omitempty" json:"habitats,omitempty"`   // ZZmaxLength=0
	Migration *Migration `xml:"migration,omitempty" json:"migration,omitempty"`   // ZZmaxLength=0
	Phenologies *Phenologies `xml:"phenologies,omitempty" json:"phenologies,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"ecologyAndLifeHistory,omitempty" json:"ecologyAndLifeHistory,omitempty"`
}

type EcologyAndLifeHistoryAuthors struct {
	AttrdisplayValue string`xml:"displayValue,attr"  json:",omitempty"`  // maxLength=17
	XMLName  xml.Name `xml:"ecologyAndLifeHistoryAuthors,omitempty" json:"ecologyAndLifeHistoryAuthors,omitempty"`
}

type EcologyAndLifeHistoryDescription struct {
	ShortGeneralDescription *ShortGeneralDescription `xml:"shortGeneralDescription,omitempty" json:"shortGeneralDescription,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"ecologyAndLifeHistoryDescription,omitempty" json:"ecologyAndLifeHistoryDescription,omitempty"`
}

type EcologyAndLifeHistoryEditionDate struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=10
	XMLName  xml.Name `xml:"ecologyAndLifeHistoryEditionDate,omitempty" json:"ecologyAndLifeHistoryEditionDate,omitempty"`
}

type Endemism struct {
	Code *Code `xml:"code,omitempty" json:"code,omitempty"`   // ZZmaxLength=0
	Description *Description `xml:"description,omitempty" json:"description,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"endemism,omitempty" json:"endemism,omitempty"`
}

type EstimatedNumberOfOccurrences struct {
	Code *Code `xml:"code,omitempty" json:"code,omitempty"`   // ZZmaxLength=0
	Comments *Comments `xml:"comments,omitempty" json:"comments,omitempty"`   // ZZmaxLength=0
	Description *Description `xml:"description,omitempty" json:"description,omitempty"`   // ZZmaxLength=0
	SearchValue *SearchValue `xml:"searchValue,omitempty" json:"searchValue,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"estimatedNumberOfOccurrences,omitempty" json:"estimatedNumberOfOccurrences,omitempty"`
}

type Family struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=11
	XMLName  xml.Name `xml:"family,omitempty" json:"family,omitempty"`
}

type FoodComments struct {
	I []*I `xml:"i,omitempty" json:"i,omitempty"`   // ZZmaxLength=0
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=101
	XMLName  xml.Name `xml:"foodComments,omitempty" json:"foodComments,omitempty"`
}

type FoodHabits struct {
	FoodComments *FoodComments `xml:"foodComments,omitempty" json:"foodComments,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"foodHabits,omitempty" json:"foodHabits,omitempty"`
}

type FormalTaxonomy struct {
	Class *Class `xml:"class,omitempty" json:"class,omitempty"`   // ZZmaxLength=0
	Family *Family `xml:"family,omitempty" json:"family,omitempty"`   // ZZmaxLength=0
	Genus *Genus `xml:"genus,omitempty" json:"genus,omitempty"`   // ZZmaxLength=0
	GenusSize *GenusSize `xml:"genusSize,omitempty" json:"genusSize,omitempty"`   // ZZmaxLength=0
	Kingdom *Kingdom `xml:"kingdom,omitempty" json:"kingdom,omitempty"`   // ZZmaxLength=0
	Order *Order `xml:"order,omitempty" json:"order,omitempty"`   // ZZmaxLength=0
	Phylum *Phylum `xml:"phylum,omitempty" json:"phylum,omitempty"`   // ZZmaxLength=0
	TaxonomicComments *TaxonomicComments `xml:"taxonomicComments,omitempty" json:"taxonomicComments,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"formalTaxonomy,omitempty" json:"formalTaxonomy,omitempty"`
}

type FormattedFullCitation struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=314
	XMLName  xml.Name `xml:"formattedFullCitation,omitempty" json:"formattedFullCitation,omitempty"`
}

type FormattedName struct {
	I []*I `xml:"i,omitempty" json:"i,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"formattedName,omitempty" json:"formattedName,omitempty"`
}

type Genus struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=9
	XMLName  xml.Name `xml:"genus,omitempty" json:"genus,omitempty"`
}

type GenusSize struct {
	Code *Code `xml:"code,omitempty" json:"code,omitempty"`   // ZZmaxLength=0
	Description *Description `xml:"description,omitempty" json:"description,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"genusSize,omitempty" json:"genusSize,omitempty"`
}

type GlobalAbundance struct {
	Code *Code `xml:"code,omitempty" json:"code,omitempty"`   // ZZmaxLength=0
	Comments *Comments `xml:"comments,omitempty" json:"comments,omitempty"`   // ZZmaxLength=0
	Description *Description `xml:"description,omitempty" json:"description,omitempty"`   // ZZmaxLength=0
	SearchValue *SearchValue `xml:"searchValue,omitempty" json:"searchValue,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"globalAbundance,omitempty" json:"globalAbundance,omitempty"`
}

type GlobalInventoryNeeds struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=113
	XMLName  xml.Name `xml:"globalInventoryNeeds,omitempty" json:"globalInventoryNeeds,omitempty"`
}

type GlobalProtection struct {
	Code *Code `xml:"code,omitempty" json:"code,omitempty"`   // ZZmaxLength=0
	Comments *Comments `xml:"comments,omitempty" json:"comments,omitempty"`   // ZZmaxLength=0
	Description *Description `xml:"description,omitempty" json:"description,omitempty"`   // ZZmaxLength=0
	Needs *Needs `xml:"needs,omitempty" json:"needs,omitempty"`   // ZZmaxLength=0
	SearchValue *SearchValue `xml:"searchValue,omitempty" json:"searchValue,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"globalProtection,omitempty" json:"globalProtection,omitempty"`
}

type GlobalRange struct {
	Code *Code `xml:"code,omitempty" json:"code,omitempty"`   // ZZmaxLength=0
	Comments *Comments `xml:"comments,omitempty" json:"comments,omitempty"`   // ZZmaxLength=0
	Description *Description `xml:"description,omitempty" json:"description,omitempty"`   // ZZmaxLength=0
	SearchValue *SearchValue `xml:"searchValue,omitempty" json:"searchValue,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"globalRange,omitempty" json:"globalRange,omitempty"`
}

type GlobalShortTermTrend struct {
	Code *Code `xml:"code,omitempty" json:"code,omitempty"`   // ZZmaxLength=0
	Comments *Comments `xml:"comments,omitempty" json:"comments,omitempty"`   // ZZmaxLength=0
	Description *Description `xml:"description,omitempty" json:"description,omitempty"`   // ZZmaxLength=0
	SearchValue *SearchValue `xml:"searchValue,omitempty" json:"searchValue,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"globalShortTermTrend,omitempty" json:"globalShortTermTrend,omitempty"`
}

type GlobalSpecies struct {
	AttrspeciesCode string`xml:"speciesCode,attr"  json:",omitempty"`  // maxLength=10
	Attruid string`xml:"uid,attr"  json:",omitempty"`  // maxLength=23
	Classification *Classification `xml:"classification,omitempty" json:"classification,omitempty"`   // ZZmaxLength=0
	ConservationStatus *ConservationStatus `xml:"conservationStatus,omitempty" json:"conservationStatus,omitempty"`   // ZZmaxLength=0
	Distribution *Distribution `xml:"distribution,omitempty" json:"distribution,omitempty"`   // ZZmaxLength=0
	EcologyAndLifeHistory *EcologyAndLifeHistory `xml:"ecologyAndLifeHistory,omitempty" json:"ecologyAndLifeHistory,omitempty"`   // ZZmaxLength=0
	License *License `xml:"license,omitempty" json:"license,omitempty"`   // ZZmaxLength=0
	ManagementSummary *ManagementSummary `xml:"managementSummary,omitempty" json:"managementSummary,omitempty"`   // ZZmaxLength=0
	NatureServeExplorerURI *NatureServeExplorerURI `xml:"natureServeExplorerURI,omitempty" json:"natureServeExplorerURI,omitempty"`   // ZZmaxLength=0
	PopulationOccurrence *PopulationOccurrence `xml:"populationOccurrence,omitempty" json:"populationOccurrence,omitempty"`   // ZZmaxLength=0
	References *References `xml:"references,omitempty" json:"references,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"globalSpecies,omitempty" json:"globalSpecies,omitempty"`
}

type GlobalStatus struct {
	ConservationStatusFactors *ConservationStatusFactors `xml:"conservationStatusFactors,omitempty" json:"conservationStatusFactors,omitempty"`   // ZZmaxLength=0
	NationalStatuses *NationalStatuses `xml:"nationalStatuses,omitempty" json:"nationalStatuses,omitempty"`   // ZZmaxLength=0
	Rank *Rank `xml:"rank,omitempty" json:"rank,omitempty"`   // ZZmaxLength=0
	Reasons *Reasons `xml:"reasons,omitempty" json:"reasons,omitempty"`   // ZZmaxLength=0
	RoundedRank *RoundedRank `xml:"roundedRank,omitempty" json:"roundedRank,omitempty"`   // ZZmaxLength=0
	StatusLastChanged *StatusLastChanged `xml:"statusLastChanged,omitempty" json:"statusLastChanged,omitempty"`   // ZZmaxLength=0
	StatusLastReviewed *StatusLastReviewed `xml:"statusLastReviewed,omitempty" json:"statusLastReviewed,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"globalStatus,omitempty" json:"globalStatus,omitempty"`
}

type Habitat struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=20
	XMLName  xml.Name `xml:"habitat,omitempty" json:"habitat,omitempty"`
}

type HabitatComments struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=266
	XMLName  xml.Name `xml:"habitatComments,omitempty" json:"habitatComments,omitempty"`
}

type Habitats struct {
	HabitatComments *HabitatComments `xml:"habitatComments,omitempty" json:"habitatComments,omitempty"`   // ZZmaxLength=0
	PalustrineHabitats *PalustrineHabitats `xml:"palustrineHabitats,omitempty" json:"palustrineHabitats,omitempty"`   // ZZmaxLength=0
	TerrestrialHabitats *TerrestrialHabitats `xml:"terrestrialHabitats,omitempty" json:"terrestrialHabitats,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"habitats,omitempty" json:"habitats,omitempty"`
}

type InferredExtentjustification struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=464
	XMLName  xml.Name `xml:"inferredExtentjustification,omitempty" json:"inferredExtentjustification,omitempty"`
}

type InferredMinimumExtentOfHabitatUse struct {
	Distance *Distance `xml:"distance,omitempty" json:"distance,omitempty"`   // ZZmaxLength=0
	InferredExtentjustification *InferredExtentjustification `xml:"inferredExtentjustification,omitempty" json:"inferredExtentjustification,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"inferredMinimumExtentOfHabitatUse,omitempty" json:"inferredMinimumExtentOfHabitatUse,omitempty"`
}

type InformalTaxonomy struct {
	InformalTaxonomyLevel1Name *InformalTaxonomyLevel1Name `xml:"informalTaxonomyLevel1Name,omitempty" json:"informalTaxonomyLevel1Name,omitempty"`   // ZZmaxLength=0
	InformalTaxonomyLevel2Name *InformalTaxonomyLevel2Name `xml:"informalTaxonomyLevel2Name,omitempty" json:"informalTaxonomyLevel2Name,omitempty"`   // ZZmaxLength=0
	InformalTaxonomyLevel3Name *InformalTaxonomyLevel3Name `xml:"informalTaxonomyLevel3Name,omitempty" json:"informalTaxonomyLevel3Name,omitempty"`   // ZZmaxLength=0
	InformalTaxonomyLevel4Name *InformalTaxonomyLevel4Name `xml:"informalTaxonomyLevel4Name,omitempty" json:"informalTaxonomyLevel4Name,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"informalTaxonomy,omitempty" json:"informalTaxonomy,omitempty"`
}

type InformalTaxonomyLevel1Name struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=22
	XMLName  xml.Name `xml:"informalTaxonomyLevel1Name,omitempty" json:"informalTaxonomyLevel1Name,omitempty"`
}

type InformalTaxonomyLevel2Name struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=7
	XMLName  xml.Name `xml:"informalTaxonomyLevel2Name,omitempty" json:"informalTaxonomyLevel2Name,omitempty"`
}

type InformalTaxonomyLevel3Name struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=21
	XMLName  xml.Name `xml:"informalTaxonomyLevel3Name,omitempty" json:"informalTaxonomyLevel3Name,omitempty"`
}

type InformalTaxonomyLevel4Name struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=24
	XMLName  xml.Name `xml:"informalTaxonomyLevel4Name,omitempty" json:"informalTaxonomyLevel4Name,omitempty"`
}

type IntrinsicVulnerability struct {
	Comments *Comments `xml:"comments,omitempty" json:"comments,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"intrinsicVulnerability,omitempty" json:"intrinsicVulnerability,omitempty"`
}

type Kingdom struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=8
	XMLName  xml.Name `xml:"kingdom,omitempty" json:"kingdom,omitempty"`
}

type License struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=144
	XMLName  xml.Name `xml:"license,omitempty" json:"license,omitempty"`
}

type LocallyMigrant struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=5
	XMLName  xml.Name `xml:"locallyMigrant,omitempty" json:"locallyMigrant,omitempty"`
}

type LongDistanceMigrant struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=5
	XMLName  xml.Name `xml:"longDistanceMigrant,omitempty" json:"longDistanceMigrant,omitempty"`
}

type ManagementSummary struct {
	BiologicalResearchNeeds *BiologicalResearchNeeds `xml:"biologicalResearchNeeds,omitempty" json:"biologicalResearchNeeds,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"managementSummary,omitempty" json:"managementSummary,omitempty"`
}

type MappingGuidance struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=889
	XMLName  xml.Name `xml:"mappingGuidance,omitempty" json:"mappingGuidance,omitempty"`
}

type MaximumLastObservedYear struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=4
	XMLName  xml.Name `xml:"maximumLastObservedYear,omitempty" json:"maximumLastObservedYear,omitempty"`
}

type Migration struct {
	LocallyMigrant *LocallyMigrant `xml:"locallyMigrant,omitempty" json:"locallyMigrant,omitempty"`   // ZZmaxLength=0
	LongDistanceMigrant *LongDistanceMigrant `xml:"longDistanceMigrant,omitempty" json:"longDistanceMigrant,omitempty"`   // ZZmaxLength=0
	NonMigrant *NonMigrant `xml:"nonMigrant,omitempty" json:"nonMigrant,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"migration,omitempty" json:"migration,omitempty"`
}

type MinimumCriteriaForOccurrence struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=304
	XMLName  xml.Name `xml:"minimumCriteriaForOccurrence,omitempty" json:"minimumCriteriaForOccurrence,omitempty"`
}

type NameUsedInConceptReference struct {
	FormattedName *FormattedName `xml:"formattedName,omitempty" json:"formattedName,omitempty"`   // ZZmaxLength=0
	UnformattedName *UnformattedName `xml:"unformattedName,omitempty" json:"unformattedName,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"nameUsedInConceptReference,omitempty" json:"nameUsedInConceptReference,omitempty"`
}

type Names struct {
	NatureServePrimaryGlobalCommonName *NatureServePrimaryGlobalCommonName `xml:"natureServePrimaryGlobalCommonName,omitempty" json:"natureServePrimaryGlobalCommonName,omitempty"`   // ZZmaxLength=0
	OtherGlobalCommonNames *OtherGlobalCommonNames `xml:"otherGlobalCommonNames,omitempty" json:"otherGlobalCommonNames,omitempty"`   // ZZmaxLength=0
	ScientificName *ScientificName `xml:"scientificName,omitempty" json:"scientificName,omitempty"`   // ZZmaxLength=0
	Synonyms *Synonyms `xml:"synonyms,omitempty" json:"synonyms,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"names,omitempty" json:"names,omitempty"`
}

type Nation struct {
	AttrnationCode string`xml:"nationCode,attr"  json:",omitempty"`  // maxLength=2
	AttrnationName string`xml:"nationName,attr"  json:",omitempty"`  // maxLength=13
	NationalDistributions *NationalDistributions `xml:"nationalDistributions,omitempty" json:"nationalDistributions,omitempty"`   // ZZmaxLength=0
	Subnations *Subnations `xml:"subnations,omitempty" json:"subnations,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"nation,omitempty" json:"nation,omitempty"`
}

type NationalDistribution struct {
	CurrentPresenceAbsence *CurrentPresenceAbsence `xml:"currentPresenceAbsence,omitempty" json:"currentPresenceAbsence,omitempty"`   // ZZmaxLength=0
	DistributionConfidence *DistributionConfidence `xml:"distributionConfidence,omitempty" json:"distributionConfidence,omitempty"`   // ZZmaxLength=0
	Origin *Origin `xml:"origin,omitempty" json:"origin,omitempty"`   // ZZmaxLength=0
	Population *Population `xml:"population,omitempty" json:"population,omitempty"`   // ZZmaxLength=0
	Regularity *Regularity `xml:"regularity,omitempty" json:"regularity,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"nationalDistribution,omitempty" json:"nationalDistribution,omitempty"`
}

type NationalDistributions struct {
	NationalDistribution *NationalDistribution `xml:"nationalDistribution,omitempty" json:"nationalDistribution,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"nationalDistributions,omitempty" json:"nationalDistributions,omitempty"`
}

type NationalStatus struct {
	AttrnationCode string`xml:"nationCode,attr"  json:",omitempty"`  // maxLength=2
	AttrnationName string`xml:"nationName,attr"  json:",omitempty"`  // maxLength=13
	Rank *Rank `xml:"rank,omitempty" json:"rank,omitempty"`   // ZZmaxLength=0
	RoundedRank *RoundedRank `xml:"roundedRank,omitempty" json:"roundedRank,omitempty"`   // ZZmaxLength=0
	StatusLastReviewed *StatusLastReviewed `xml:"statusLastReviewed,omitempty" json:"statusLastReviewed,omitempty"`   // ZZmaxLength=0
	SubnationalStatuses *SubnationalStatuses `xml:"subnationalStatuses,omitempty" json:"subnationalStatuses,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"nationalStatus,omitempty" json:"nationalStatus,omitempty"`
}

type NationalStatuses struct {
	NationalStatus []*NationalStatus `xml:"nationalStatus,omitempty" json:"nationalStatus,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"nationalStatuses,omitempty" json:"nationalStatuses,omitempty"`
}

type Nations struct {
	Nation []*Nation `xml:"nation,omitempty" json:"nation,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"nations,omitempty" json:"nations,omitempty"`
}

type NatureServePrimaryGlobalCommonName struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=20
	XMLName  xml.Name `xml:"natureServePrimaryGlobalCommonName,omitempty" json:"natureServePrimaryGlobalCommonName,omitempty"`
}

type NatureServeStatus struct {
	GlobalStatus *GlobalStatus `xml:"globalStatus,omitempty" json:"globalStatus,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"natureServeStatus,omitempty" json:"natureServeStatus,omitempty"`
}

type Needs struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=120
	XMLName  xml.Name `xml:"needs,omitempty" json:"needs,omitempty"`
}

type NomenclaturalAuthor struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=31
	XMLName  xml.Name `xml:"nomenclaturalAuthor,omitempty" json:"nomenclaturalAuthor,omitempty"`
}

type NonMigrant struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=5
	XMLName  xml.Name `xml:"nonMigrant,omitempty" json:"nonMigrant,omitempty"`
}

type Notes struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=967
	XMLName  xml.Name `xml:"notes,omitempty" json:"notes,omitempty"`
}

type OccurrenceCounties struct {
	OccurrenceCounty []*OccurrenceCounty `xml:"occurrenceCounty,omitempty" json:"occurrenceCounty,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"occurrenceCounties,omitempty" json:"occurrenceCounties,omitempty"`
}

type OccurrenceCounty struct {
	CountyCode *CountyCode `xml:"countyCode,omitempty" json:"countyCode,omitempty"`   // ZZmaxLength=0
	CountyName *CountyName `xml:"countyName,omitempty" json:"countyName,omitempty"`   // ZZmaxLength=0
	MaximumLastObservedYear *MaximumLastObservedYear `xml:"maximumLastObservedYear,omitempty" json:"maximumLastObservedYear,omitempty"`   // ZZmaxLength=0
	SpeciesOccurrenceCount *SpeciesOccurrenceCount `xml:"speciesOccurrenceCount,omitempty" json:"speciesOccurrenceCount,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"occurrenceCounty,omitempty" json:"occurrenceCounty,omitempty"`
}

type OccurrenceNation struct {
	Attrcode string`xml:"code,attr"  json:",omitempty"`  // maxLength=2
	OccurrenceStates *OccurrenceStates `xml:"occurrenceStates,omitempty" json:"occurrenceStates,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"occurrenceNation,omitempty" json:"occurrenceNation,omitempty"`
}

type OccurrenceNations struct {
	OccurrenceNation *OccurrenceNation `xml:"occurrenceNation,omitempty" json:"occurrenceNation,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"occurrenceNations,omitempty" json:"occurrenceNations,omitempty"`
}

type OccurrenceState struct {
	Attrcode string`xml:"code,attr"  json:",omitempty"`  // maxLength=2
	OccurrenceCounties *OccurrenceCounties `xml:"occurrenceCounties,omitempty" json:"occurrenceCounties,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"occurrenceState,omitempty" json:"occurrenceState,omitempty"`
}

type OccurrenceStates struct {
	OccurrenceState []*OccurrenceState `xml:"occurrenceState,omitempty" json:"occurrenceState,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"occurrenceStates,omitempty" json:"occurrenceStates,omitempty"`
}

type Order struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=11
	XMLName  xml.Name `xml:"order,omitempty" json:"order,omitempty"`
}

type Origin struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=20
	XMLName  xml.Name `xml:"origin,omitempty" json:"origin,omitempty"`
}

type OtherConsiderations struct {
	XMLName  xml.Name `xml:"otherConsiderations,omitempty" json:"otherConsiderations,omitempty"`
}

type OtherGlobalCommonNames struct {
	CommonName []*CommonName `xml:"commonName,omitempty" json:"commonName,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"otherGlobalCommonNames,omitempty" json:"otherGlobalCommonNames,omitempty"`
}

type OtherStatuses struct {
	Status []*Status `xml:"status,omitempty" json:"status,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"otherStatuses,omitempty" json:"otherStatuses,omitempty"`
}

type PalustrineHabitats struct {
	Habitat []*Habitat `xml:"habitat,omitempty" json:"habitat,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"palustrineHabitats,omitempty" json:"palustrineHabitats,omitempty"`
}

type Phenologies struct {
	PhenologyComments *PhenologyComments `xml:"phenologyComments,omitempty" json:"phenologyComments,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"phenologies,omitempty" json:"phenologies,omitempty"`
}

type PhenologyComments struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=256
	XMLName  xml.Name `xml:"phenologyComments,omitempty" json:"phenologyComments,omitempty"`
}

type Phylum struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=11
	XMLName  xml.Name `xml:"phylum,omitempty" json:"phylum,omitempty"`
}

type Population struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=10
	XMLName  xml.Name `xml:"population,omitempty" json:"population,omitempty"`
}

type PopulationOccurrence struct {
	Delineations *Delineations `xml:"delineations,omitempty" json:"delineations,omitempty"`   // ZZmaxLength=0
	Viabilities *Viabilities `xml:"viabilities,omitempty" json:"viabilities,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"populationOccurrence,omitempty" json:"populationOccurrence,omitempty"`
}

type Rank struct {
	Code *Code `xml:"code,omitempty" json:"code,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"rank,omitempty" json:"rank,omitempty"`
}

type Reasons struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=168
	XMLName  xml.Name `xml:"reasons,omitempty" json:"reasons,omitempty"`
}

type References struct {
	Citation []*Citation `xml:"citation,omitempty" json:"citation,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"references,omitempty" json:"references,omitempty"`
}

type Regularity struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=19
	XMLName  xml.Name `xml:"regularity,omitempty" json:"regularity,omitempty"`
}

type RoundedRank struct {
	Code *Code `xml:"code,omitempty" json:"code,omitempty"`   // ZZmaxLength=0
	Description *Description `xml:"description,omitempty" json:"description,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"roundedRank,omitempty" json:"roundedRank,omitempty"`
}

type ScientificName struct {
	ConceptReference *ConceptReference `xml:"conceptReference,omitempty" json:"conceptReference,omitempty"`   // ZZmaxLength=0
	FormattedName *FormattedName `xml:"formattedName,omitempty" json:"formattedName,omitempty"`   // ZZmaxLength=0
	NomenclaturalAuthor *NomenclaturalAuthor `xml:"nomenclaturalAuthor,omitempty" json:"nomenclaturalAuthor,omitempty"`   // ZZmaxLength=0
	UnformattedName *UnformattedName `xml:"unformattedName,omitempty" json:"unformattedName,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"scientificName,omitempty" json:"scientificName,omitempty"`
}

type SynonymName struct {
	ConceptReference *ConceptReference `xml:"conceptReference,omitempty" json:"conceptReference,omitempty"`   // ZZmaxLength=0
	FormattedName *FormattedName `xml:"formattedName,omitempty" json:"formattedName,omitempty"`   // ZZmaxLength=0
	NomenclaturalAuthor *NomenclaturalAuthor `xml:"nomenclaturalAuthor,omitempty" json:"nomenclaturalAuthor,omitempty"`   // ZZmaxLength=0
	UnformattedName *UnformattedName `xml:"unformattedName,omitempty" json:"unformattedName,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"synonymName,omitempty" json:"synonymName,omitempty"`
}

type SearchValue struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=3
	XMLName  xml.Name `xml:"searchValue,omitempty" json:"searchValue,omitempty"`
}

type Separation struct {
	AlternateSeparationProcedure *AlternateSeparationProcedure `xml:"alternateSeparationProcedure,omitempty" json:"alternateSeparationProcedure,omitempty"`   // ZZmaxLength=0
	Barriers *Barriers `xml:"barriers,omitempty" json:"barriers,omitempty"`   // ZZmaxLength=0
	DistanceForSuitableHabitat *DistanceForSuitableHabitat `xml:"distanceForSuitableHabitat,omitempty" json:"distanceForSuitableHabitat,omitempty"`   // ZZmaxLength=0
	DistanceForUnsuitableHabitat *DistanceForUnsuitableHabitat `xml:"distanceForUnsuitableHabitat,omitempty" json:"distanceForUnsuitableHabitat,omitempty"`   // ZZmaxLength=0
	SeparationJustification *SeparationJustification `xml:"separationJustification,omitempty" json:"separationJustification,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"separation,omitempty" json:"separation,omitempty"`
}

type SeparationJustification struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=1191
	XMLName  xml.Name `xml:"separationJustification,omitempty" json:"separationJustification,omitempty"`
}

type ShortGeneralDescription struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=38
	XMLName  xml.Name `xml:"shortGeneralDescription,omitempty" json:"shortGeneralDescription,omitempty"`
}

type SpeciesOccurrenceCount struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=2
	XMLName  xml.Name `xml:"speciesOccurrenceCount,omitempty" json:"speciesOccurrenceCount,omitempty"`
}

type Status struct {
	Attrname string`xml:"name,attr"  json:",omitempty"`  // maxLength=14
	StatusDate *StatusDate `xml:"statusDate,omitempty" json:"statusDate,omitempty"`   // ZZmaxLength=0
	StatusDescription *StatusDescription `xml:"statusDescription,omitempty" json:"statusDescription,omitempty"`   // ZZmaxLength=0
	StatusValue *StatusValue `xml:"statusValue,omitempty" json:"statusValue,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"status,omitempty" json:"status,omitempty"`
}

type StatusDate struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=10
	XMLName  xml.Name `xml:"statusDate,omitempty" json:"statusDate,omitempty"`
}

type StatusDescription struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=15
	XMLName  xml.Name `xml:"statusDescription,omitempty" json:"statusDescription,omitempty"`
}

type StatusFactorsEditionDate struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=10
	XMLName  xml.Name `xml:"statusFactorsEditionDate,omitempty" json:"statusFactorsEditionDate,omitempty"`
}

type StatusLastChanged struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=10
	XMLName  xml.Name `xml:"statusLastChanged,omitempty" json:"statusLastChanged,omitempty"`
}

type StatusLastReviewed struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=10
	XMLName  xml.Name `xml:"statusLastReviewed,omitempty" json:"statusLastReviewed,omitempty"`
}

type StatusValue struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=2
	XMLName  xml.Name `xml:"statusValue,omitempty" json:"statusValue,omitempty"`
}

type Subnation struct {
	AttrsubnationCode string`xml:"subnationCode,attr"  json:",omitempty"`  // maxLength=2
	AttrsubnationName string`xml:"subnationName,attr"  json:",omitempty"`  // maxLength=21
	SubnationalDistributions *SubnationalDistributions `xml:"subnationalDistributions,omitempty" json:"subnationalDistributions,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"subnation,omitempty" json:"subnation,omitempty"`
}

type SubnationalDistribution struct {
	CurrentPresenceAbsence *CurrentPresenceAbsence `xml:"currentPresenceAbsence,omitempty" json:"currentPresenceAbsence,omitempty"`   // ZZmaxLength=0
	DistributionConfidence *DistributionConfidence `xml:"distributionConfidence,omitempty" json:"distributionConfidence,omitempty"`   // ZZmaxLength=0
	Origin *Origin `xml:"origin,omitempty" json:"origin,omitempty"`   // ZZmaxLength=0
	Population *Population `xml:"population,omitempty" json:"population,omitempty"`   // ZZmaxLength=0
	Regularity *Regularity `xml:"regularity,omitempty" json:"regularity,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"subnationalDistribution,omitempty" json:"subnationalDistribution,omitempty"`
}

type SubnationalDistributions struct {
	SubnationalDistribution *SubnationalDistribution `xml:"subnationalDistribution,omitempty" json:"subnationalDistribution,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"subnationalDistributions,omitempty" json:"subnationalDistributions,omitempty"`
}

type SubnationalStatus struct {
	AttrsubnationCode string`xml:"subnationCode,attr"  json:",omitempty"`  // maxLength=2
	AttrsubnationName string`xml:"subnationName,attr"  json:",omitempty"`  // maxLength=21
	Rank *Rank `xml:"rank,omitempty" json:"rank,omitempty"`   // ZZmaxLength=0
	RoundedRank *RoundedRank `xml:"roundedRank,omitempty" json:"roundedRank,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"subnationalStatus,omitempty" json:"subnationalStatus,omitempty"`
}

type SubnationalStatuses struct {
	SubnationalStatus []*SubnationalStatus `xml:"subnationalStatus,omitempty" json:"subnationalStatus,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"subnationalStatuses,omitempty" json:"subnationalStatuses,omitempty"`
}

type Subnations struct {
	Subnation []*Subnation `xml:"subnation,omitempty" json:"subnation,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"subnations,omitempty" json:"subnations,omitempty"`
}

type Synonyms struct {
	SynonymName []*SynonymName `xml:"synonymName,omitempty" json:"synonymName,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"synonyms,omitempty" json:"synonyms,omitempty"`
}

type Taxonomy struct {
	FormalTaxonomy *FormalTaxonomy `xml:"formalTaxonomy,omitempty" json:"formalTaxonomy,omitempty"`   // ZZmaxLength=0
	InformalTaxonomy *InformalTaxonomy `xml:"informalTaxonomy,omitempty" json:"informalTaxonomy,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"taxonomy,omitempty" json:"taxonomy,omitempty"`
}

type TerrestrialHabitats struct {
	Habitat []*Habitat `xml:"habitat,omitempty" json:"habitat,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"terrestrialHabitats,omitempty" json:"terrestrialHabitats,omitempty"`
}

type Threat struct {
	Comments *Comments `xml:"comments,omitempty" json:"comments,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"threat,omitempty" json:"threat,omitempty"`
}

type UnformattedName struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=35
	XMLName  xml.Name `xml:"unformattedName,omitempty" json:"unformattedName,omitempty"`
}

type VersionDate struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=10
	XMLName  xml.Name `xml:"versionDate,omitempty" json:"versionDate,omitempty"`
}

type Viabilities struct {
	Viability *Viability `xml:"viability,omitempty" json:"viability,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"viabilities,omitempty" json:"viabilities,omitempty"`
}

type Viability struct {
	ViabilityJustification *ViabilityJustification `xml:"viabilityJustification,omitempty" json:"viabilityJustification,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"viability,omitempty" json:"viability,omitempty"`
}

type ViabilityJustification struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=201
	XMLName  xml.Name `xml:"viabilityJustification,omitempty" json:"viabilityJustification,omitempty"`
}

type Watershed struct {
	Attrtype string`xml:"type,attr"  json:",omitempty"`  // maxLength=5
	MaximumLastObservedYear *MaximumLastObservedYear `xml:"maximumLastObservedYear,omitempty" json:"maximumLastObservedYear,omitempty"`   // ZZmaxLength=0
	SpeciesOccurrenceCount *SpeciesOccurrenceCount `xml:"speciesOccurrenceCount,omitempty" json:"speciesOccurrenceCount,omitempty"`   // ZZmaxLength=0
	WatershedCodes *WatershedCodes `xml:"watershedCodes,omitempty" json:"watershedCodes,omitempty"`   // ZZmaxLength=0
	WatershedName *WatershedName `xml:"watershedName,omitempty" json:"watershedName,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"watershed,omitempty" json:"watershed,omitempty"`
}

type WatershedCode struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=8
	XMLName  xml.Name `xml:"watershedCode,omitempty" json:"watershedCode,omitempty"`
}

type WatershedCodes struct {
	WatershedCode *WatershedCode `xml:"watershedCode,omitempty" json:"watershedCode,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"watershedCodes,omitempty" json:"watershedCodes,omitempty"`
}

type WatershedList struct {
	Watershed []*Watershed `xml:"watershed,omitempty" json:"watershed,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"watershedList,omitempty" json:"watershedList,omitempty"`
}

type WatershedName struct {
	Text string `xml:",chardata" json:",omitempty"`   // maxLength=31
	XMLName  xml.Name `xml:"watershedName,omitempty" json:"watershedName,omitempty"`
}

type Watersheds struct {
	WatershedList *WatershedList `xml:"watershedList,omitempty" json:"watershedList,omitempty"`   // ZZmaxLength=0
	XMLName  xml.Name `xml:"watersheds,omitempty" json:"watersheds,omitempty"`
}