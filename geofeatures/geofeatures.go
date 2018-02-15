package geofeatures

import (
	"bitbucket.org/heindl/processors/ecoregions"
	"bitbucket.org/heindl/processors/utils"
	"github.com/dropbox/godropbox/errors"
	"googlemaps.github.io/maps"
	"os"
	//"github.com/cenkalti/backoff"
	"math"
	"sync"
	"fmt"
	"cloud.google.com/go/firestore"
	"google.golang.org/appengine"
	"strings"
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
	CoordinatesEstimated bool `json:"" firestore:""`
	Biome                ecoregions.Biome `json:"" firestore:""`
	Realm                ecoregions.Realm `json:"" firestore:""`
	EcoNum               ecoregions.EcoNum `json:"" firestore:""`
	Elevation            *float64 `json:"" firestore:""`
	Geopoint             *appengine.GeoPoint `json:"" firestore:""`
}

func (Ω GeoFeatureSet) Lat() float64 {
	return Ω.Geopoint.Lat
}

func (Ω GeoFeatureSet) Lng() float64 {
	return Ω.Geopoint.Lng
}

const keyCoordinate = "CoordinateKey"
//const keyGeoPoint = "GeoPoint"
//const keyCoordinatesEstimated = "CoordinatesEstimated"
//const keyRealm = "Realm"
//const keyBiome = "Biome"
//const keyEcoNum = "EcoNum"
//const keyS2Tokens = "S2Tokens"
//const keyElevation = "Elevation"

func (Ω GeoFeatureSet) CoordinateKey() string {
	// Intentionally reduce the precision of the coordinates to ensure we're not duplicating occurrences.
	return fmt.Sprintf("%.3f|%.3f", Ω.Lat(), Ω.Lng())
}

func (Ω GeoFeatureSet) CoordinateQuery(collection *firestore.CollectionRef) (firestore.Query, error) {
	return collection.Where(keyCoordinate, "==", Ω.CoordinateKey()), nil
}

//func NewGeoFeatureSetFromMap(m map[string]interface{}) (*GeoFeatureSet, error) {
//
//	geopoint := m[keyGeoPoint].(appengine.GeoPoint)
//
//	if err := validateCoordinates(&geopoint); err != nil {
//		return nil, err
//	}
//
//	gs := GeoFeatureSet{
//		CoordinatesEstimated: m[keyCoordinatesEstimated].(bool),
//		Biome:                m[keyBiome].(ecoregions.Biome),
//		Realm:                m[keyRealm].(ecoregions.Realm),
//		EcoNum:               m[keyEcoNum].(ecoregions.EcoNum),
//		Elevation:            utils.FloatPtr(m[keyElevation].(float64)),
//		Geopoint:             &geopoint,
//	}
//
//	if !gs.Biome.Valid() {
//		return nil, errors.New("Invalid Biome")
//	}
//
//	if !gs.EcoNum.Valid() {
//		return nil, errors.New("Invalid EcoNum")
//	}
//
//	if gs.Elevation == nil {
//		return nil, errors.New("Invalid Elevation")
//	}
//
//	return &gs, nil
//
//}

var ErrInvalidCoordinate = errors.New("Invalid Coordinate")

func validateCoordinates(geopoint *appengine.GeoPoint) error {

	if !geopoint.Valid() {
		return errors.Wrapf(ErrInvalidCoordinate,"Invalid GeoPoint [%f, %f]", geopoint.Lat, geopoint.Lng)
	}

	if geopoint.Lat < 6.6 || geopoint.Lat > 83.3 {
		return errors.Wrapf(ErrInvalidCoordinate,"latitude [%f] is out of bounds", geopoint.Lat)
	}
	if geopoint.Lng < -178.2 || geopoint.Lng > -49.0 {
		return errors.Wrapf(ErrInvalidCoordinate,"longitude [%f] is out of bounds", geopoint.Lng)
	}
	// We need the decimal precision to be at least a football field, so require at least three decimal places (110m).
	if hasDecimalPlaces(2, geopoint.Lat) || hasDecimalPlaces(1, geopoint.Lat) {
		return errors.Wrapf(ErrInvalidCoordinate,"latitude [%f] has insufficient precision", geopoint.Lat)
	}
	if hasDecimalPlaces(2, geopoint.Lng) || hasDecimalPlaces(1, geopoint.Lng) {
		return errors.Wrapf(ErrInvalidCoordinate,"longitude [%f] has insufficient precision", geopoint.Lng)
	}
	return nil
}

func NewGeoFeatureSet(lat, lng float64, coordinatesEstimated bool) (*GeoFeatureSet, error) {

	geopoint := appengine.GeoPoint{Lat: lat, Lng: lng}

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
		Geopoint:             &geopoint,
		CoordinatesEstimated: coordinatesEstimated,
		Biome:                ecoID.Biome(),
		Realm:                ecoID.Realm(),
		EcoNum:               ecoID.EcoNum(),
		}, nil
}
//
//func (Ω *GeoFeatureSet) ToMap() (map[string]interface{}, error) {
//
//	if err := validateCoordinates(Ω.Geopoint); err != nil {
//		return nil, err
//	}
//
//	elevation, err := liveProcessor.getElevation(Ω.Geopoint.Lat, Ω.Geopoint.Lng)
//	if err != nil {
//		return nil, err
//	}
//
//	// Flush elevations. Assumes our entry has been queued.
//	if elevation == nil && len(liveProcessor.elevationsQueued) > 0 {
//		if err := liveProcessor.flushElevations(); err != nil {
//			return nil, err
//		}
//	}
//
//	elevation, err = liveProcessor.getElevation(Ω.Geopoint.Lat, Ω.Geopoint.Lng)
//	if err != nil {
//		return nil, err
//	}
//	if elevation == nil {
//		return nil, errors.Newf("Elevation still not generated after flush [%s]", elevationKey(Ω.Geopoint.Lat, Ω.Geopoint.Lng))
//	}
//
//	if !Ω.Biome.Valid() || !Ω.EcoNum.Valid() {
//		return nil, errors.New("Invalid EcoID")
//	}
//
//	return map[string]interface{}{
//		keyGeoPoint: Ω.Geopoint,
//		keyCoordinatesEstimated: Ω.CoordinatesEstimated,
//		keyBiome:       Ω.Biome,
//		keyRealm:       Ω.Realm,
//		keyEcoNum:      Ω.EcoNum,
//		keyS2Tokens:    terra.NewPoint(Ω.Geopoint.Lat, Ω.Geopoint.Lng).S2TokenMap(),
//		keyElevation:   elevation,
//		keyCoordinate: Ω.CoordinateKey(),
//	}, nil
//}


func hasDecimalPlaces(i int, v float64) bool {
	vf := v * math.Pow(10.0, float64(i))
	extra := vf - float64(int(vf))
	return extra == 0
}

var liveProcessor = &geoFeaturesProcessor{
	elevationsQueued:    map[string]maps.LatLng{},
	elevationsFetched: map[string]float64{},
}

type geoFeaturesProcessor struct {
	mapClient         *maps.Client
	ecoRegionCache    *ecoregions.EcoRegionsCache
	sync.Mutex
	elevationsQueued    map[string]maps.LatLng
	elevationsFetched map[string]float64
}

func init() {

	apiKey := os.Getenv("FLORACAST_GCLOUD_API_KEY")
	if apiKey == "" {
		panic(errors.New("Missing API Key"))
	}

	var err error
	liveProcessor.mapClient, err = maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		panic(errors.Wrap(err, "could not get google maps client"))
	}

	liveProcessor.ecoRegionCache, err = ecoregions.NewEcoRegionsCache()
	if err != nil {
		panic(err)
	}
}

func elevationKey(lat, lng float64) string {
	return fmt.Sprintf("%.4f,%.4f", lat, lng)
}

func (Ω *geoFeaturesProcessor) elevationQueueStatus(lat, lng float64) (queued, fetched bool) {
	Ω.Lock()
	defer Ω.Unlock()
	k := elevationKey(lat, lng)
	_, queued = Ω.elevationsQueued[k]
	_, fetched = Ω.elevationsFetched[k]
	return
}

func (Ω *geoFeaturesProcessor) queueElevation(lat, lng float64) error {
	queued, fetched := Ω.elevationQueueStatus(lat, lng)
	if queued || fetched {
		return nil
	}
	k := elevationKey(lat, lng)
	Ω.Lock()
	Ω.elevationsQueued[k] = maps.LatLng{lat, lng}
	Ω.Unlock()
	if len(Ω.elevationsQueued) >= 10 {
		return Ω.flushElevations()
	}
	return nil
}

func (Ω *geoFeaturesProcessor) getElevation(lat, lng float64) (*float64, error) {
	queued, fetched := Ω.elevationQueueStatus(lat, lng)

	if !fetched {
		return nil, nil
	}
	if !fetched && queued {
		return nil, nil
	}
	if !fetched && !queued {
		return nil, errors.Newf("Trying to get coordinates neither fetched or queued [%s] with %d previously fetched", elevationKey(lat, lng), len(Ω.elevationsFetched))
	}
	return utils.FloatPtr(Ω.elevationsFetched[elevationKey(lat, lng)]), nil
}

func (Ω *geoFeaturesProcessor) flushElevations() error {
	Ω.Lock()
	defer Ω.Unlock()

	locs := []string{}
	for _, k := range Ω.elevationsQueued {
		locs = append(locs, elevationKey(k.Lat, k.Lng))
	}
	if len(locs) == 0 {
		return nil
	}

	var res struct {
		Results []struct {
			Lat float64 `json:"latitude"`
			Lng float64 `json:"longitude"`
			Elevation float64 `json:"Elevation"` // Meters
		} `json:"results"`
	}

	if err := utils.RequestJSON("https://api.open-Elevation.com/api/v1/lookup?locations=" + strings.Join(locs, "|"), &res); err != nil {
		return errors.Wrap(err, "Could not fetch Elevation api")
	}

	//resolvedElevations, err := Ω.mapClient.Elevation(context.Background(), &eleReq)
	//if err != nil {
	//	return errors.Wrap(err, "could not fetch elevations")
	//}
	//for _, e := range resolvedElevations {
	//	k := elevationKey(e.Location.Lat, e.Location.Lng)
	//	Ω.elevationsFetched[k] = e.Elevation
	//}
	for _, r := range res.Results {
		k := elevationKey(r.Lat, r.Lng)
		if r.Elevation == 0 {
			fmt.Println("Elevation key [%s] not resolved", k)
		}
		Ω.elevationsFetched[k] = r.Elevation
	}

	//fmt.Println("Flushing Elevations", len(Ω.elevationsQueued), len(Ω.elevationsFetched), len(locs), len(resolvedElevations))
	Ω.elevationsQueued = map[string]maps.LatLng{}
	return nil
}
//
//
//func (Ω *geoFeaturesProcessor) processElevationBatch(cxt context.Context, locations ...PredictableLocation) error {
//
//
//	// First batch and fetch elevations
//	eleReq := maps.ElevationRequest{Locations: []maps.LatLng{}}
//	// Gather lat/lng pairs for Elevation fetch.
//	for _, o := range locations {
//		eleReq.Locations = append(eleReq.Locations, maps.LatLng{o.Lat(), o.Lng()})
//	}
//
//	//resolvedElevations := []maps.ElevationResult{}
//
//	//bkf := backoff.NewExponentialBackOff()
//	//bkf.InitialInterval = time.Second * 1
//	//ticker := backoff.NewTicker(bkf)
//	//for _ = range ticker.C {
//
//
//		resolvedElevations, err := Ω.mapClient.Elevation(cxt, &eleReq)
//		//if err != nil && strings.Contains(err.Error(), "TLS handshake timeout") {
//		//	fmt.Println("TLS handshake timeout encountered. Backing off ...")
//		//	continue
//		//}
//		//if err != nil && strings.Contains(err.Error(), "DATA_NOT_AVAILABLE") {
//		//	fmt.Println("DATA_NOT_AVAILABLE", len(resolvedElevations))
//		//	ticker.Stop()
//		//	break
//		//}
//		if err != nil {
//			//ticker.Stop()
//			return errors.Wrap(err, "could not fetch elevations")
//		}
//
//	//	ticker.Stop()
//	//	break
//	//}
//
//	tmb := tomb.Tomb{}
//	tmb.Go(func() error {
//		for _, _loc := range locations {
//			loc := _loc
//			tmb.Go(func() error {
//				for _, _r := range resolvedElevations {
//					r := _r
//					if !utils.CoordinatesEqual(loc.Lat(), r.Location.Lat) {
//						continue
//					}
//					if !utils.CoordinatesEqual(loc.Lng(), r.Location.Lng) {
//						continue
//					}
//					return loc.SetElevation(r.Elevation)
//				}
//				return errors.Newf("Elevation not found: %.5f, %.5f", loc.Lat(), loc.Lng())
//			})
//		}
//		return nil
//	})
//
//	return tmb.Wait()
//
//}
