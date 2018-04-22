package geoembed

import (
	"bitbucket.org/heindl/process/terra/ecoregions"
	"bitbucket.org/heindl/process/terra/elevation"
	"bitbucket.org/heindl/process/terra/geo"
	"encoding/json"
	"github.com/dropbox/godropbox/errors"
	"google.golang.org/genproto/googleapis/type/latlng"
)

type localStructForJSON struct {
	GeoPoint             *latlng.LatLng `json:""`
	CoordinatesEstimated bool
	EcoRealm             ecoregions.Realm  `json:""`
	EcoBiome             ecoregions.Biome  `json:""`
	EcoNum               ecoregions.EcoNum `json:""`
	Elevation            *int              `json:",omitempty"`
	S2Tokens             map[string]string `json:""`
}

// UnmarshalJSON is an interface method for converting a FeatureSet to JSON
func (Ω *GeoFeatureSet) UnmarshalJSON(b []byte) error {

	local := localStructForJSON{}
	if err := json.Unmarshal(b, &local); err != nil {
		return errors.Wrap(err, "Could not unmarshal GeoFeatureSet")
	}

	if local.Elevation == nil {
		return errors.New("Invalid Elevation")
	}

	Ω.coordinatesEstimated = local.CoordinatesEstimated
	Ω.biome = local.EcoBiome
	Ω.realm = local.EcoRealm
	Ω.ecoNum = local.EcoNum
	Ω.elevation = local.Elevation
	Ω.geoPoint = local.GeoPoint

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

// MarshalJSON is an interface method for converting JSON to a FeatureSet
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

	return json.Marshal(localStructForJSON{
		GeoPoint:             Ω.geoPoint,
		CoordinatesEstimated: Ω.coordinatesEstimated,
		EcoBiome:             Ω.biome,
		EcoRealm:             Ω.realm,
		EcoNum:               Ω.ecoNum,
		S2Tokens:             terraPoint.S2TokenMap(),
		Elevation:            elev,
	})
}
