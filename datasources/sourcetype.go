package datasources


type SourceType string

const (
	DataSourceTypeGBIF             = SourceType("27")
	DataSourceTypeINaturalist      = SourceType("INAT")
	DataSourceTypeMushroomObserver = SourceType("MUOB")
	DataSourceTypeNatureServe      = SourceType("11")
)

func HasDataSourceType(srcs []SourceType, sourceType SourceType) bool {
	for _, src := range srcs {
		if src == sourceType {
			return true
		}
	}
	return false
}

func (Ω SourceType) Valid() bool {
	_, ok := SchemeSourceIDMap[Ω]
	return ok
}

var SchemeSourceIDMap = map[SourceType]string{
	// Floracast
	DataSourceTypeINaturalist: "iNaturalist",
	// INaturalist
	SourceType("1"):                "IUCN Red List of Threatened Species. Version 2012.1",
	SourceType("3"):                "Amphibian Species of the World 5.6",
	SourceType("2"):                "Amphibiaweb. 2012",
	SourceType("5"):                "Amphibian Species of the World 5.5",
	SourceType("17"):               "New England Wild Flower Society's Flora Novae Angliae",
	DataSourceTypeNatureServe:      "NatureServe Explorer: An online encyclopedia of life. Version 7.1",
	SourceType("12"):               "Calflora",
	SourceType("13"):               "Odonata Central",
	SourceType("14"):               "IUCN Red List of Threatened Species. Version 2012.2",
	SourceType("10"):               "eBird/Clements Checklist 6.7",
	SourceType("15"):               "CONABIO",
	SourceType("6"):                "The Reptile Database",
	SourceType("16"):               "Afribats",
	SourceType("18"):               "Norma 059, 2010",
	SourceType("4"):                "Draft IUCN/SSC, 2013.1",
	SourceType("19"):               "Draft IUCN/SSC Amphibian Specialist Group, 2011",
	SourceType("20"):               "eBird/Clements Checklist 6.8",
	SourceType("21"):               "IUCN Red List of Threatened Species. Version 2013.2",
	SourceType("22"):               "eBird/Clements Checklist 6.9",
	SourceType("23"):               "NatureWatch NZ",
	SourceType("24"):               "The world spider catalog, version 15.5",
	SourceType("25"):               "Carabidae of the World",
	SourceType("26"):               "IUCN Red List of Threatened Species. Version 2014.3",
	DataSourceTypeGBIF:             "GBIF",
	SourceType("28"):               "NPSpecies",
	SourceType("29"):               "Esslinger&#39;s North American Lichens",
	SourceType("30"):               "Amphibian Species of the World 6.0",
	SourceType("31"):               "Esslinger&#39;s North American Lichens, Version 21",
	DataSourceTypeMushroomObserver: "MushroomObserver.org",
}



