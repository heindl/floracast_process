package datasources


type DataSourceType string

const (
	DataSourceTypeGBIF             = DataSourceType("27")
	DataSourceTypeINaturalist      = DataSourceType("INAT")
	DataSourceTypeMushroomObserver = DataSourceType("MUOB")
	DataSourceTypeNatureServe      = DataSourceType("11")
)

func HasDataSourceType(srcs []DataSourceType, sourceType DataSourceType) bool {
	for _, src := range srcs {
		if src == sourceType {
			return true
		}
	}
	return false
}

func (Ω DataSourceType) Valid() bool {
	_, ok := SchemeSourceIDMap[Ω]
	return ok
}

var SchemeSourceIDMap = map[DataSourceType]string{
	// Floracast
	DataSourceTypeINaturalist: "iNaturalist",
	// INaturalist
	DataSourceType("1"):       "IUCN Red List of Threatened Species. Version 2012.1",
	DataSourceType("3"):       "Amphibian Species of the World 5.6",
	DataSourceType("2"):       "Amphibiaweb. 2012",
	DataSourceType("5"):       "Amphibian Species of the World 5.5",
	DataSourceType("17"):      "New England Wild Flower Society's Flora Novae Angliae",
	DataSourceTypeNatureServe: "NatureServe Explorer: An online encyclopedia of life. Version 7.1",
	DataSourceType("12"):      "Calflora",
	DataSourceType("13"):      "Odonata Central",
	DataSourceType("14"):      "IUCN Red List of Threatened Species. Version 2012.2",
	DataSourceType("10"):      "eBird/Clements Checklist 6.7",
	DataSourceType("15"):      "CONABIO",
	DataSourceType("6"):       "The Reptile Database",
	DataSourceType("16"):      "Afribats",
	DataSourceType("18"):      "Norma 059, 2010",
	DataSourceType("4"):          "Draft IUCN/SSC, 2013.1",
	DataSourceType("19"):         "Draft IUCN/SSC Amphibian Specialist Group, 2011",
	DataSourceType("20"):         "eBird/Clements Checklist 6.8",
	DataSourceType("21"):         "IUCN Red List of Threatened Species. Version 2013.2",
	DataSourceType("22"):         "eBird/Clements Checklist 6.9",
	DataSourceType("23"):         "NatureWatch NZ",
	DataSourceType("24"):           "The world spider catalog, version 15.5",
	DataSourceType("25"):           "Carabidae of the World",
	DataSourceType("26"):           "IUCN Red List of Threatened Species. Version 2014.3",
	DataSourceTypeGBIF:             "GBIF",
	DataSourceType("28"):           "NPSpecies",
	DataSourceType("29"):           "Esslinger&#39;s North American Lichens",
	DataSourceType("30"):           "Amphibian Species of the World 6.0",
	DataSourceType("31"):           "Esslinger&#39;s North American Lichens, Version 21",
	DataSourceTypeMushroomObserver: "MushroomObserver.org",
}



