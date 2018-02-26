package geofeatures

import (
	"google.golang.org/genproto/googleapis/type/latlng"
	"bitbucket.org/heindl/process/ecoregions"
	"github.com/dropbox/godropbox/errors"
	"bitbucket.org/heindl/process/terra"
	"encoding/json"
	"bitbucket.org/heindl/process/utils"
)

func (Ω *GeoFeatureSet) UnmarshalJSON(b []byte) error {
	m := map[string]interface{}{}
	if err := json.Unmarshal(b, &m); err != nil {
		return errors.Wrap(err, "Could not unmarshal FeatureSet")
	}

	gp, ok := m[keyGeoPoint]
	if !ok {
		return errors.New("geopoint missing from GeoFeatureSet")
	}

	geopoint := gp.(latlng.LatLng)

	if err := validateCoordinates(&geopoint); err != nil {
		return err
	}

	gs := GeoFeatureSet{
		coordinatesEstimated: m[keyCoordinatesEstimated].(bool),
		biome:                m[keyEcoBiome].(ecoregions.Biome),
		realm:                m[keyEcoRealm].(ecoregions.Realm),
		ecoNum:               m[keyEcoNum].(ecoregions.EcoNum),
		elevation:            utils.IntPtr(m[keyElevation].(int)),
		geopoint:             &geopoint,
	}

	if !gs.biome.Valid() {
		return errors.New("Invalid biome in GeoFeatureSet")
	}

	if !gs.ecoNum.Valid() {
		return errors.New("Invalid ecoNum in GeoFeatureSet")
	}

	if gs.elevation == nil {
		return errors.New("Invalid elevation in GeoFeatureSet")
	}

	Ω = &gs

	return nil
}

func (Ω *GeoFeatureSet) MarshalJSON() ([]byte, error) {

	if err := validateCoordinates(Ω.geopoint); err != nil {
		return nil, err
	}

	elevation, err := liveProcessor.getElevation(Ω.geopoint.GetLatitude(), Ω.geopoint.GetLongitude())
	if err != nil {
		return nil, err
	}

	// Flush elevations. Assumes our entry has been queued.
	if elevation == nil && len(liveProcessor.elevationsQueued) > 0 {
		if err := liveProcessor.flushElevations(); err != nil {
			return nil, err
		}
	}

	elevation, err = liveProcessor.getElevation(Ω.geopoint.GetLatitude(), Ω.geopoint.GetLongitude())
	if err != nil {
		return nil, err
	}
	if elevation == nil {
		return nil, errors.Newf("Elevation still not generated after flush [%s]", elevationKey(Ω.geopoint.GetLatitude(), Ω.geopoint.GetLongitude()))
	}

	if !Ω.biome.Valid() || !Ω.ecoNum.Valid() {
		return nil, errors.New("Invalid EcoID")
	}

	terraPoint, err := terra.NewPoint(Ω.geopoint.GetLatitude(), Ω.geopoint.GetLongitude())
	if err != nil {
		return nil, err
	}

	coordKey, err := Ω.CoordinateKey()
	if err != nil {
		return nil, err
	}

	return json.Marshal(map[string]interface{}{
		keyGeoPoint:             Ω.geopoint,
		keyCoordinatesEstimated: Ω.coordinatesEstimated,
		keyEcoBiome:             Ω.biome,
		keyEcoRealm:             Ω.realm,
		keyEcoNum:               Ω.ecoNum,
		keyS2Tokens:             terraPoint.S2TokenMap(),
		keyElevation:            elevation,
		keyCoordinate:           coordKey,
	})
}
