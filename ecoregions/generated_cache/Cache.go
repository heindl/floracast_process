package cache

type CachedEcoRegion struct {
	Realm      string
	EcoName    string
	EcoCode    string
	EcoNum     int64
	EcoID      int64
	Biome      int64
	Geometries []string
}

var EcoRegionCache = []CachedEcoRegion{}
