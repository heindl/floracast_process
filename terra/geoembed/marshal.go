package geoembed

import (
	"bitbucket.org/heindl/process/terra/ecoregions"
	"bitbucket.org/heindl/process/terra/elevation"
	"bitbucket.org/heindl/process/terra/geo"
	"bitbucket.org/heindl/process/utils"
	"encoding/json"
	"github.com/dropbox/godropbox/errors"
	"google.golang.org/genproto/googleapis/type/latlng"
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

	if err := geo.ValidateCoordinates(geopoint.GetLatitude(), geopoint.GetLongitude()); err != nil {
		return err
	}

	Ω = &GeoFeatureSet{
		coordinatesEstimated: m[keyCoordinatesEstimated].(bool),
		biome:                m[keyEcoBiome].(ecoregions.Biome),
		realm:                m[keyEcoRealm].(ecoregions.Realm),
		ecoNum:               m[keyEcoNum].(ecoregions.EcoNum),
		elevation:            utils.IntPtr(m[keyElevation].(int)),
		geopoint:             &geopoint,
	}

	if !Ω.biome.Valid() {
		return errors.New("Invalid biome in GeoFeatureSet")
	}

	if !Ω.ecoNum.Valid() {
		return errors.New("Invalid ecoNum in GeoFeatureSet")
	}

	if Ω.elevation == nil {
		return errors.New("Invalid elevation in GeoFeatureSet")
	}

	return nil
}

func (Ω *GeoFeatureSet) MarshalJSON() ([]byte, error) {

	if err := geo.ValidateCoordinates(Ω.Lat(), Ω.Lng()); err != nil {
		return nil, err
	}

	elev, err := elevation.Get(Ω.Lat(), Ω.Lng())
	if err != nil {
		return nil, err
	}
	if elev == nil {
		return nil, errors.Newf("Elevation not fetched [%f, %f]", Ω.Lat(), Ω.Lng())
	}

	if !Ω.biome.Valid() || !Ω.ecoNum.Valid() {
		return nil, errors.Newf("Invalid Biome [%d] and EcoNum [%d]", Ω.biome, Ω.ecoNum)
	}

	terraPoint, err := geo.NewPoint(Ω.Lat(), Ω.Lng())
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
		keyElevation:            elev,
		keyCoordinate:           coordKey,
	})
}
