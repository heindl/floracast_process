package geoembed

import (
	"bitbucket.org/heindl/process/terra/ecoregions"
	"bitbucket.org/heindl/process/terra/ecoregions/cache"
	"bitbucket.org/heindl/process/terra/elevation"
	"bitbucket.org/heindl/process/terra/geo"
	"cloud.google.com/go/firestore"
	"google.golang.org/genproto/googleapis/type/latlng"
)

// Description of Precision
// https://gis.stackexchange.com/questions/8650/measuring-accuracy-of-latitude-and-longitude
// The sign tells us whether we are north or south, east or west on the globe.
// A nonzero hundreds digit tells us we're using longitude, not latitude!
// The tens digit gives a position to about 1,000 kilometers. It gives us useful information about what continent or ocean we are on.
// The units digit (one decimal degree) gives a position up to 111 kilometers (60 nautical miles, about 69 miles). It can tell us roughly what large state or country we are in.
// The first decimal place is worth up to 11.1 km: it can distinguish the position of one large city from a neighboring large city.
// The second decimal place is worth up to 1.1 km: it can separate one village from the next.
// The third decimal place is worth up to 110 m: it can identify a large agricultural field or institutional campus.
// The fourth decimal place is worth up to 11 m: it can identify a parcel of land. It is comparable to the typical accuracy of an uncorrected GPS unit with no interference.
// The fifth decimal place is worth up to 1.1 m: it distinguish trees from each other. Accuracy to this level with commercial GPS units can only be achieved with differential correction.
// The sixth decimal place is worth up to 0.11 m: you can use this for laying out structures in detail, for designing landscapes, building roads. It should be more than good enough for tracking movements of glaciers and rivers. This can be achieved by taking painstaking measures with GPS, such as differentially corrected GPS.

type GeoFeatureSet struct {
	coordinatesEstimated bool
	biome                ecoregions.Biome
	realm                ecoregions.Realm
	ecoNum               ecoregions.EcoNum
	elevation            *int
	geopoint             *latlng.LatLng
}

const keyGeoPoint = "GeoPoint"
const keyCoordinatesEstimated = "CoordinatesEstimated"
const keyEcoRealm = "EcoRealm"
const keyEcoBiome = "EcoBiome"
const keyEcoNum = "EcoNum"
const keyS2Tokens = "S2Tokens"
const keyElevation = "Elevation"
const keyCoordinate = "CoordinateKey"

func (Ω *GeoFeatureSet) Lat() float64 {
	return Ω.geopoint.GetLatitude()
}

func (Ω *GeoFeatureSet) Lng() float64 {
	return Ω.geopoint.GetLongitude()
}

func NewGeoFeatureSet(lat, lng float64, coordinatesEstimated bool) (*GeoFeatureSet, error) {

	if err := geo.ValidateCoordinates(lat, lng); err != nil {
		return nil, err
	}

	region, err := cache.FetchEcologicalRegion(lat, lng)
	if err != nil {
		return nil, err
	}

	if err := elevation.Queue(lat, lng); err != nil {
		return nil, err
	}

	return &GeoFeatureSet{
		geopoint:             &latlng.LatLng{Latitude: lat, Longitude: lng},
		coordinatesEstimated: coordinatesEstimated,
		biome:                region.Biome(),
		realm:                region.Realm(),
		ecoNum:               region.EcoNum(),
	}, nil
}

func CoordinateQuery(collection *firestore.CollectionRef, lat, lng float64) (*firestore.Query, error) {
	k, err := NewCoordinateKey(lat, lng)
	if err != nil {
		return nil, err
	}
	q := collection.Where(keyCoordinate, "==", k)
	return &q, nil
}
