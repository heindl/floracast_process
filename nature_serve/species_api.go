package nature_serve
//
//import "encoding/xml"
//
//
//
//type Citation struct {
//	Text string `xml:",chardata" json:",omitempty"`   // maxLength=130
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 citation,omitempty" json:"citation,omitempty"`
//}
//
//type Class struct {
//	Text string `xml:",chardata" json:",omitempty"`   // maxLength=11
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 class,omitempty" json:"class,omitempty"`
//}
//
//type Classification struct {
//	Names *Names `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 names,omitempty" json:"names,omitempty"`   // ZZmaxLength=0
//	Taxonomy *Taxonomy `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 taxonomy,omitempty" json:"taxonomy,omitempty"`   // ZZmaxLength=0
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 classification,omitempty" json:"classification,omitempty"`
//}
//
//type ClassificationStatus struct {
//	Text string `xml:",chardata" json:",omitempty"`   // maxLength=8
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 classificationStatus,omitempty" json:"classificationStatus,omitempty"`
//}
//
//type Code struct {
//	Text string `xml:",chardata" json:",omitempty"`   // maxLength=4
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 code,omitempty" json:"code,omitempty"`
//}
//
//type CommonName struct {
//	Attrlanguage string`xml:"language,attr"  json:",omitempty"`  // maxLength=2
//	Text string `xml:",chardata" json:",omitempty"`   // maxLength=12
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 commonName,omitempty" json:"commonName,omitempty"`
//}
//
//type ConceptReference struct {
//	Attrcode string`xml:"code,attr"  json:",omitempty"`  // maxLength=12
//	ClassificationStatus *ClassificationStatus `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 classificationStatus,omitempty" json:"classificationStatus,omitempty"`   // ZZmaxLength=0
//	FormattedFullCitation *FormattedFullCitation `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 formattedFullCitation,omitempty" json:"formattedFullCitation,omitempty"`   // ZZmaxLength=0
//	NameUsedInConceptReference *NameUsedInConceptReference `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 nameUsedInConceptReference,omitempty" json:"nameUsedInConceptReference,omitempty"`   // ZZmaxLength=0
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 conceptReference,omitempty" json:"conceptReference,omitempty"`
//}
//
//type ConservationStatus struct {
//	NatureServeStatus *NatureServeStatus `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 natureServeStatus,omitempty" json:"natureServeStatus,omitempty"`   // ZZmaxLength=0
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 conservationStatus,omitempty" json:"conservationStatus,omitempty"`
//}
//
//type ConservationStatusMap struct {
//	Text string `xml:",chardata" json:",omitempty"`   // maxLength=5
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 conservationStatusMap,omitempty" json:"conservationStatusMap,omitempty"`
//}
//
//type CurrentPresenceAbsence struct {
//	Text string `xml:",chardata" json:",omitempty"`   // maxLength=20
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 currentPresenceAbsence,omitempty" json:"currentPresenceAbsence,omitempty"`
//}
//
//type Description struct {
//	Text string `xml:",chardata" json:",omitempty"`   // maxLength=6
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 description,omitempty" json:"description,omitempty"`
//}
//
//type Distribution struct {
//	AttrUSAndCanadianDistributionComplete string`xml:"USAndCanadianDistributionComplete,attr"  json:",omitempty"`  // maxLength=5
//	ConservationStatusMap *ConservationStatusMap `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 conservationStatusMap,omitempty" json:"conservationStatusMap,omitempty"`   // ZZmaxLength=0
//	Nations *nations `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 nations,omitempty" json:"nations,omitempty"`   // ZZmaxLength=0
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 distribution,omitempty" json:"distribution,omitempty"`
//}
//
//type DistributionConfidence struct {
//	Text string `xml:",chardata" json:",omitempty"`   // maxLength=9
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 distributionConfidence,omitempty" json:"distributionConfidence,omitempty"`
//}
//
//type Durations struct {
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 durations,omitempty" json:"durations,omitempty"`
//}
//
//type EcologyAndLifeHistory struct {
//	Durations *Durations `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 durations,omitempty" json:"durations,omitempty"`   // ZZmaxLength=0
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 ecologyAndLifeHistory,omitempty" json:"ecologyAndLifeHistory,omitempty"`
//}
//
//type Family struct {
//	Text string `xml:",chardata" json:",omitempty"`   // maxLength=13
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 family,omitempty" json:"family,omitempty"`
//}
//
//type FormalTaxonomy struct {
//	Class *Class `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 class,omitempty" json:"class,omitempty"`   // ZZmaxLength=0
//	Family *Family `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 family,omitempty" json:"family,omitempty"`   // ZZmaxLength=0
//	Genus *Genus `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 genus,omitempty" json:"genus,omitempty"`   // ZZmaxLength=0
//	Kingdom *Kingdom `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 kingdom,omitempty" json:"kingdom,omitempty"`   // ZZmaxLength=0
//	Order *Order `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 order,omitempty" json:"order,omitempty"`   // ZZmaxLength=0
//	Phylum *Phylum `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 phylum,omitempty" json:"phylum,omitempty"`   // ZZmaxLength=0
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 formalTaxonomy,omitempty" json:"formalTaxonomy,omitempty"`
//}
//
//type FormattedFullCitation struct {
//	Text string `xml:",chardata" json:",omitempty"`   // maxLength=91
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 formattedFullCitation,omitempty" json:"formattedFullCitation,omitempty"`
//}
//
//type FormattedName struct {
//	I *I `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 i,omitempty" json:"i,omitempty"`   // ZZmaxLength=0
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 formattedName,omitempty" json:"formattedName,omitempty"`
//}
//
//type Genus struct {
//	Text string `xml:",chardata" json:",omitempty"`   // maxLength=9
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 genus,omitempty" json:"genus,omitempty"`
//}
//
//type GlobalSpecies struct {
//	AttrspeciesCode string`xml:"speciesCode,attr"  json:",omitempty"`  // maxLength=10
//	Attruid string`xml:"uid,attr"  json:",omitempty"`  // maxLength=23
//	Classification *Classification `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 classification,omitempty" json:"classification,omitempty"`   // ZZmaxLength=0
//	ConservationStatus *ConservationStatus `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 conservationStatus,omitempty" json:"conservationStatus,omitempty"`   // ZZmaxLength=0
//	Distribution *Distribution `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 distribution,omitempty" json:"distribution,omitempty"`   // ZZmaxLength=0
//	EcologyAndLifeHistory *EcologyAndLifeHistory `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 ecologyAndLifeHistory,omitempty" json:"ecologyAndLifeHistory,omitempty"`   // ZZmaxLength=0
//	License *License `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 license,omitempty" json:"license,omitempty"`   // ZZmaxLength=0
//	NatureServeExplorerURI *NatureServeExplorerURI `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 natureServeExplorerURI,omitempty" json:"natureServeExplorerURI,omitempty"`   // ZZmaxLength=0
//	ReferenceCount *ReferenceCount `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 references,omitempty" json:"references,omitempty"`   // ZZmaxLength=0
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 globalSpecies,omitempty" json:"globalSpecies,omitempty"`
//}
//
//type GlobalSpeciesList struct {
//	AttrXsiSpaceschemaLocation string`xml:"http://www.w3.org/2001/XMLSchema-instance schemaLocation,attr"  json:",omitempty"`  // maxLength=136
//	AttrschemaVersion string`xml:"schemaVersion,attr"  json:",omitempty"`  // maxLength=3
//	Attrxmlns string`xml:"xmlns,attr"  json:",omitempty"`  // maxLength=67
//	AttrXmlnsxsi string`xml:"xmlns xsi,attr"  json:",omitempty"`  // maxLength=41
//	GlobalSpecies *GlobalSpecies `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 globalSpecies,omitempty" json:"globalSpecies,omitempty"`   // ZZmaxLength=0
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 globalSpeciesList,omitempty" json:"globalSpeciesList,omitempty"`
//}
//
//type GlobalStatus struct {
//	NationalStatuses *NationalStatuses `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 nationalStatuses,omitempty" json:"nationalStatuses,omitempty"`   // ZZmaxLength=0
//	Rank *Rank `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 rank,omitempty" json:"rank,omitempty"`   // ZZmaxLength=0
//	RoundedRank *RoundedRank `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 roundedRank,omitempty" json:"roundedRank,omitempty"`   // ZZmaxLength=0
//	StatusLastChanged *StatusLastChanged `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 statusLastChanged,omitempty" json:"statusLastChanged,omitempty"`   // ZZmaxLength=0
//	StatusLastReviewed *StatusLastReviewed `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 statusLastReviewed,omitempty" json:"statusLastReviewed,omitempty"`   // ZZmaxLength=0
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 globalStatus,omitempty" json:"globalStatus,omitempty"`
//}
//
//type I struct {
//	Text string `xml:",chardata" json:",omitempty"`   // maxLength=19
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 i,omitempty" json:"i,omitempty"`
//}
//
//type InformalTaxonomy struct {
//	InformalTaxonomyLevel1Name *InformalTaxonomyLevel1Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 informalTaxonomyLevel1Name,omitempty" json:"informalTaxonomyLevel1Name,omitempty"`   // ZZmaxLength=0
//	InformalTaxonomyLevel2Name *InformalTaxonomyLevel2Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 informalTaxonomyLevel2Name,omitempty" json:"informalTaxonomyLevel2Name,omitempty"`   // ZZmaxLength=0
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 informalTaxonomy,omitempty" json:"informalTaxonomy,omitempty"`
//}
//
//type InformalTaxonomyLevel1Name struct {
//	Text string `xml:",chardata" json:",omitempty"`   // maxLength=13
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 informalTaxonomyLevel1Name,omitempty" json:"informalTaxonomyLevel1Name,omitempty"`
//}
//
//type InformalTaxonomyLevel2Name struct {
//	Text string `xml:",chardata" json:",omitempty"`   // maxLength=22
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 informalTaxonomyLevel2Name,omitempty" json:"informalTaxonomyLevel2Name,omitempty"`
//}
//
//type Kingdom struct {
//	Text string `xml:",chardata" json:",omitempty"`   // maxLength=5
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 kingdom,omitempty" json:"kingdom,omitempty"`
//}
//
//type License struct {
//	Text string `xml:",chardata" json:",omitempty"`   // maxLength=144
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 license,omitempty" json:"license,omitempty"`
//}
//
//type NameUsedInConceptReference struct {
//	FormattedName *FormattedName `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 formattedName,omitempty" json:"formattedName,omitempty"`   // ZZmaxLength=0
//	UnformattedName *UnformattedName `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 unformattedName,omitempty" json:"unformattedName,omitempty"`   // ZZmaxLength=0
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 nameUsedInConceptReference,omitempty" json:"nameUsedInConceptReference,omitempty"`
//}
//
//type Names struct {
//	NatureServePrimaryGlobalCommonName *NatureServePrimaryGlobalCommonName `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 natureServePrimaryGlobalCommonName,omitempty" json:"natureServePrimaryGlobalCommonName,omitempty"`   // ZZmaxLength=0
//	OtherGlobalCommonNames *OtherGlobalCommonNames `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 otherGlobalCommonNames,omitempty" json:"otherGlobalCommonNames,omitempty"`   // ZZmaxLength=0
//	ScientificName *ScientificName `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 scientificName,omitempty" json:"scientificName,omitempty"`   // ZZmaxLength=0
//	Synonyms *Synonyms `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 synonyms,omitempty" json:"synonyms,omitempty"`   // ZZmaxLength=0
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 names,omitempty" json:"names,omitempty"`
//}
//
//type Nation struct {
//	AttrnationCode string`xml:"nationCode,attr"  json:",omitempty"`  // maxLength=2
//	AttrnationName string`xml:"nationName,attr"  json:",omitempty"`  // maxLength=13
//	NationalDistributions *NationalDistributions `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 nationalDistributions,omitempty" json:"nationalDistributions,omitempty"`   // ZZmaxLength=0
//	Subnations *Subnations `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 subnations,omitempty" json:"subnations,omitempty"`   // ZZmaxLength=0
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 nation,omitempty" json:"nation,omitempty"`
//}
//
//type NationalDistribution struct {
//	CurrentPresenceAbsence *CurrentPresenceAbsence `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 currentPresenceAbsence,omitempty" json:"currentPresenceAbsence,omitempty"`   // ZZmaxLength=0
//	DistributionConfidence *DistributionConfidence `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 distributionConfidence,omitempty" json:"distributionConfidence,omitempty"`   // ZZmaxLength=0
//	Origin *Origin `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 origin,omitempty" json:"origin,omitempty"`   // ZZmaxLength=0
//	Regularity *Regularity `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 regularity,omitempty" json:"regularity,omitempty"`   // ZZmaxLength=0
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 nationalDistribution,omitempty" json:"nationalDistribution,omitempty"`
//}
//
//type NationalDistributions struct {
//	NationalDistribution *NationalDistribution `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 nationalDistribution,omitempty" json:"nationalDistribution,omitempty"`   // ZZmaxLength=0
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 nationalDistributions,omitempty" json:"nationalDistributions,omitempty"`
//}
//
//type NationalStatus struct {
//	AttrnationCode string`xml:"nationCode,attr"  json:",omitempty"`  // maxLength=2
//	AttrnationName string`xml:"nationName,attr"  json:",omitempty"`  // maxLength=13
//	Rank *Rank `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 rank,omitempty" json:"rank,omitempty"`   // ZZmaxLength=0
//	RoundedRank *RoundedRank `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 roundedRank,omitempty" json:"roundedRank,omitempty"`   // ZZmaxLength=0
//	StatusLastReviewed *StatusLastReviewed `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 statusLastReviewed,omitempty" json:"statusLastReviewed,omitempty"`   // ZZmaxLength=0
//	SubnationalStatuses *SubnationalStatuses `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 subnationalStatuses,omitempty" json:"subnationalStatuses,omitempty"`   // ZZmaxLength=0
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 nationalStatus,omitempty" json:"nationalStatus,omitempty"`
//}
//
//type NationalStatuses struct {
//	NationalStatus []*NationalStatus `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 nationalStatus,omitempty" json:"nationalStatus,omitempty"`   // ZZmaxLength=0
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 nationalStatuses,omitempty" json:"nationalStatuses,omitempty"`
//}
//
//type nations struct {
//	Nation []*Nation `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 nation,omitempty" json:"nation,omitempty"`   // ZZmaxLength=0
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 nations,omitempty" json:"nations,omitempty"`
//}
//
//type NatureServeExplorerURI struct {
//	Text string `xml:",chardata" json:",omitempty"`   // maxLength=84
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 natureServeExplorerURI,omitempty" json:"natureServeExplorerURI,omitempty"`
//}
//
//type NatureServePrimaryGlobalCommonName struct {
//	Text string `xml:",chardata" json:",omitempty"`   // maxLength=12
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 natureServePrimaryGlobalCommonName,omitempty" json:"natureServePrimaryGlobalCommonName,omitempty"`
//}
//
//type NatureServeStatus struct {
//	GlobalStatus *GlobalStatus `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 globalStatus,omitempty" json:"globalStatus,omitempty"`   // ZZmaxLength=0
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 natureServeStatus,omitempty" json:"natureServeStatus,omitempty"`
//}
//
//type NomenclaturalAuthor struct {
//	Text string `xml:",chardata" json:",omitempty"`   // maxLength=5
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 nomenclaturalAuthor,omitempty" json:"nomenclaturalAuthor,omitempty"`
//}
//
//type Order struct {
//	Text string `xml:",chardata" json:",omitempty"`   // maxLength=9
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 order,omitempty" json:"order,omitempty"`
//}
//
//type Origin struct {
//	Text string `xml:",chardata" json:",omitempty"`   // maxLength=20
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 origin,omitempty" json:"origin,omitempty"`
//}
//
//type OtherGlobalCommonNames struct {
//	CommonName []*CommonName `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 commonName,omitempty" json:"commonName,omitempty"`   // ZZmaxLength=0
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 otherGlobalCommonNames,omitempty" json:"otherGlobalCommonNames,omitempty"`
//}
//
//type Phylum struct {
//	Text string `xml:",chardata" json:",omitempty"`   // maxLength=10
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 phylum,omitempty" json:"phylum,omitempty"`
//}
//
//type Rank struct {
//	Code *Code `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 code,omitempty" json:"code,omitempty"`   // ZZmaxLength=0
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 rank,omitempty" json:"rank,omitempty"`
//}
//
//type References struct {
//	Citation []*Citation `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 citation,omitempty" json:"citation,omitempty"`   // ZZmaxLength=0
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 references,omitempty" json:"references,omitempty"`
//}
//
//type Regularity struct {
//	Text string `xml:",chardata" json:",omitempty"`   // maxLength=19
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 regularity,omitempty" json:"regularity,omitempty"`
//}
//
//type RoundedRank struct {
//	Code *Code `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 code,omitempty" json:"code,omitempty"`   // ZZmaxLength=0
//	Description *Description `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 description,omitempty" json:"description,omitempty"`   // ZZmaxLength=0
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 roundedRank,omitempty" json:"roundedRank,omitempty"`
//}
//
//type ScientificName struct {
//	ConceptReference *ConceptReference `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 conceptReference,omitempty" json:"conceptReference,omitempty"`   // ZZmaxLength=0
//	FormattedName *FormattedName `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 formattedName,omitempty" json:"formattedName,omitempty"`   // ZZmaxLength=0
//	UnformattedName *UnformattedName `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 unformattedName,omitempty" json:"unformattedName,omitempty"`   // ZZmaxLength=0
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 scientificName,omitempty" json:"scientificName,omitempty"`
//}
//
//type StatusLastChanged struct {
//	Text string `xml:",chardata" json:",omitempty"`   // maxLength=10
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 statusLastChanged,omitempty" json:"statusLastChanged,omitempty"`
//}
//
//type StatusLastReviewed struct {
//	Text string `xml:",chardata" json:",omitempty"`   // maxLength=10
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 statusLastReviewed,omitempty" json:"statusLastReviewed,omitempty"`
//}
//
//type Subnation struct {
//	AttrsubnationCode string`xml:"subnationCode,attr"  json:",omitempty"`  // maxLength=2
//	AttrsubnationName string`xml:"subnationName,attr"  json:",omitempty"`  // maxLength=16
//	SubnationalDistributions *SubnationalDistributions `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 subnationalDistributions,omitempty" json:"subnationalDistributions,omitempty"`   // ZZmaxLength=0
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 subnation,omitempty" json:"subnation,omitempty"`
//}
//
//type SubnationalDistribution struct {
//	CurrentPresenceAbsence *CurrentPresenceAbsence `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 currentPresenceAbsence,omitempty" json:"currentPresenceAbsence,omitempty"`   // ZZmaxLength=0
//	DistributionConfidence *DistributionConfidence `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 distributionConfidence,omitempty" json:"distributionConfidence,omitempty"`   // ZZmaxLength=0
//	Origin *Origin `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 origin,omitempty" json:"origin,omitempty"`   // ZZmaxLength=0
//	Regularity *Regularity `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 regularity,omitempty" json:"regularity,omitempty"`   // ZZmaxLength=0
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 subnationalDistribution,omitempty" json:"subnationalDistribution,omitempty"`
//}
//
//type SubnationalDistributions struct {
//	SubnationalDistribution *SubnationalDistribution `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 subnationalDistribution,omitempty" json:"subnationalDistribution,omitempty"`   // ZZmaxLength=0
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 subnationalDistributions,omitempty" json:"subnationalDistributions,omitempty"`
//}
//
//type SubnationalStatus struct {
//	AttrsubnationCode string`xml:"subnationCode,attr"  json:",omitempty"`  // maxLength=2
//	AttrsubnationName string`xml:"subnationName,attr"  json:",omitempty"`  // maxLength=16
//	Rank *Rank `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 rank,omitempty" json:"rank,omitempty"`   // ZZmaxLength=0
//	RoundedRank *RoundedRank `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 roundedRank,omitempty" json:"roundedRank,omitempty"`   // ZZmaxLength=0
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 subnationalStatus,omitempty" json:"subnationalStatus,omitempty"`
//}
//
//type SubnationalStatuses struct {
//	SubnationalStatus []*SubnationalStatus `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 subnationalStatus,omitempty" json:"subnationalStatus,omitempty"`   // ZZmaxLength=0
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 subnationalStatuses,omitempty" json:"subnationalStatuses,omitempty"`
//}
//
//type Subnations struct {
//	Subnation []*Subnation `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 subnation,omitempty" json:"subnation,omitempty"`   // ZZmaxLength=0
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 subnations,omitempty" json:"subnations,omitempty"`
//}
//
//type SynonymName struct {
//	FormattedName *FormattedName `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 formattedName,omitempty" json:"formattedName,omitempty"`   // ZZmaxLength=0
//	NomenclaturalAuthor *NomenclaturalAuthor `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 nomenclaturalAuthor,omitempty" json:"nomenclaturalAuthor,omitempty"`   // ZZmaxLength=0
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 synonymName,omitempty" json:"synonymName,omitempty"`
//}
//
//type Synonyms struct {
//	SynonymName *SynonymName `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 synonymName,omitempty" json:"synonymName,omitempty"`   // ZZmaxLength=0
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 synonyms,omitempty" json:"synonyms,omitempty"`
//}
//
//type Taxonomy struct {
//	FormalTaxonomy *FormalTaxonomy `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 formalTaxonomy,omitempty" json:"formalTaxonomy,omitempty"`   // ZZmaxLength=0
//	InformalTaxonomy *InformalTaxonomy `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 informalTaxonomy,omitempty" json:"informalTaxonomy,omitempty"`   // ZZmaxLength=0
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 taxonomy,omitempty" json:"taxonomy,omitempty"`
//}
//
//type UnformattedName struct {
//	Text string `xml:",chardata" json:",omitempty"`   // maxLength=19
//	XMLName  xml.Name `xml:"http://services.natureserve.org/docs/schemas/biodiversityDataFlow/1 unformattedName,omitempty" json:"unformattedName,omitempty"`
//}