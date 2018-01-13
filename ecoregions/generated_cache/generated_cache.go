package cache

		
	type CachedEcoRegion struct {
		Realm  string
		//EcoName string
		EcoCode string
		EcoNum int64
		EcoID int64
		//GlobalStatus int64
		Biome int64
		//Area float64
		//AreaKm2 float64
		GeoHashes []string
		GeoJsonString string
	}


		var EcoRegionCache = []CachedEcoRegion{}