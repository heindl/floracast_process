package geofeatures

import (
	"bitbucket.org/heindl/process/ecoregions"
	"github.com/dropbox/godropbox/errors"
	"os"
	"math"
	"sync"
	"fmt"
	"cloud.google.com/go/firestore"
	"google.golang.org/genproto/googleapis/type/latlng"
	"strings"
)

var liveProcessor = &geoFeaturesProcessor{
	elevationsQueued:    map[string]*latlng.LatLng{},
	elevationsFetched: map[string]*int{},
}

type geoFeaturesProcessor struct {
	ecoRegionCache    *ecoregions.EcoRegionsCache
	sync.Mutex
	elevationsQueued    map[string]*latlng.LatLng
	elevationsFetched map[string]*int
}

func init() {

	apiKey := os.Getenv("FLORACAST_GCLOUD_API_KEY")
	if apiKey == "" {
		panic(errors.New("Missing API Key"))
	}

	var err error
	liveProcessor.ecoRegionCache, err = ecoregions.NewEcoRegionsCache()
	if err != nil {
		panic(err)
	}
}

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

type CoordinateKey string

func (Ω CoordinateKey) Valid() bool {
	if len(Ω) != 11 {
		return false
	}
	return true
}

func NewCoordinateKey(lat, lng float64) (CoordinateKey, error) {
	ll := latlng.LatLng{
		Latitude: lat,
		Longitude: lng,
	}
	if err := validateCoordinates(&ll); err != nil {
		return "", err
	}
	// Intentionally reduce the precision of the coordinates to ensure we're not duplicating occurrences.
	k := fmt.Sprintf("%.3f_%.3f", lat, lng)
	// Replace decimals in order to use as firestore key.
	k = strings.Replace(k, ".", "|", -1)
	return CoordinateKey(k), nil
}

func (Ω *GeoFeatureSet) CoordinateKey() (CoordinateKey, error) {
	return NewCoordinateKey(Ω.Lat(), Ω.Lng())

}

func CoordinateQuery(collection *firestore.CollectionRef, lat, lng float64) (*firestore.Query, error) {
	k, err := NewCoordinateKey(lat, lng)
	if err != nil {
		return nil, err
	}
	q := collection.Where(keyCoordinate, "==", k)
	return &q, nil
}


var ErrInvalidCoordinate = errors.New("Invalid Coordinate")

func validateCoordinates(geopoint *latlng.LatLng) error {
	//
	//if !geopoint.Valid() {
	//	return errors.Wrapf(ErrInvalidCoordinate,"Invalid GeoPoint [%f, %f]", geopoint.Lat, geopoint.Lng)
	//}

	lat := geopoint.GetLatitude()
	lng := geopoint.GetLongitude()

	if lat < 6.6 || lat > 83.3 {
		return errors.Wrapf(ErrInvalidCoordinate,"latitude [%f] is out of bounds", lat)
	}
	if lng < -178.2 || lng > -49.0 {
		return errors.Wrapf(ErrInvalidCoordinate,"longitude [%f] is out of bounds", lng)
	}
	// We need the decimal precision to be at least a football field, so require at least three decimal places (110m).
	if hasDecimalPlaces(2, lat) || hasDecimalPlaces(1, lat) {
		return errors.Wrapf(ErrInvalidCoordinate,"latitude [%f] has insufficient precision", lat)
	}
	if hasDecimalPlaces(2, lng) || hasDecimalPlaces(1, lng) {
		return errors.Wrapf(ErrInvalidCoordinate,"longitude [%f] has insufficient precision", lng)
	}
	return nil
}

func NewGeoFeatureSet(lat, lng float64, coordinatesEstimated bool) (*GeoFeatureSet, error) {

	geopoint := latlng.LatLng{Latitude: lat, Longitude: lng}

	if err := validateCoordinates(&geopoint); err != nil {
		return nil, err
	}

	ecoID, err := liveProcessor.ecoRegionCache.EcoID(lat, lng)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not find EcoRegion [%.3f, %.3f]", lat, lng)
	}

	if err := liveProcessor.queueElevation(lat, lng); err != nil {
		return nil, err
	}
	return &GeoFeatureSet{
		geopoint:             &geopoint,
		coordinatesEstimated: coordinatesEstimated,
		biome:                ecoID.Biome(),
		realm:                ecoID.Realm(),
		ecoNum:               ecoID.EcoNum(),
		}, nil
}

func hasDecimalPlaces(i int, v float64) bool {
	vf := v * math.Pow(10.0, float64(i))
	extra := vf - float64(int(vf))
	return extra == 0
}

