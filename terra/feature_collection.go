package terra

import (
	"github.com/dropbox/godropbox/errors"
	"github.com/paulmach/go.geojson"
)


type FeatureCollection struct{
	area float64
	polyLabel Point
	features []*Feature
}

func (Ω *FeatureCollection) Count() int {
	return len(Ω.features)
}

func (Ω *FeatureCollection) Features() []*Feature {
	return Ω.features
}

func (Ω *FeatureCollection) Append(features ...*Feature) error {
	for i := range features {
		features[i].Normalize()
		if !features[i].Valid() {
			return errors.New("invalid feature")
		}
		Ω.features = append(Ω.features, features[i])
	}
	largestArea := 0.0
	for i := range Ω.features {
		Ω.area += Ω.features[i].Area()
		if Ω.features[i].Area() > largestArea {
			largestArea = Ω.features[i].Area()
			Ω.polyLabel = Ω.features[i].PolyLabel()
		}
	}
	return nil
}

func (Ω *FeatureCollection) GeoJSON() ([]byte, error) {
	fc := geojson.NewFeatureCollection()
	for _, feature := range Ω.features {
		encoded_feature, err := feature.GeoJSON()
		if err != nil {
			return nil, err
		}
		geojson_feature, err := geojson.UnmarshalFeature(encoded_feature)
		if err != nil {
			return nil, errors.Wrap(err, "could not unmarshal geojson feature")
		}
		fc = fc.AddFeature(geojson_feature)
	}

	// PolyLabel
	//fc = fc.AddFeature(geojson.NewPointFeature())

	encoded_feature_collection, err := fc.MarshalJSON()
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal geojson feature collection")
	}
	return encoded_feature_collection, nil
}

func (Ω *FeatureCollection) Contains(lat, lng float64) bool {
	for _, feature := range Ω.features {
		if feature.Contains(lat, lng) {
			return true
		}
	}
	return false
}

func (Ω *FeatureCollection) PolyLabel() Point {
	return Ω.polyLabel
}

func (Ω *FeatureCollection) Area() float64 {
	return Ω.area
}
func (Ω *FeatureCollection) FilterByProperty(should_filter func(interface{}) bool, property_key string, ) *FeatureCollection {
	if property_key == "" {
		return Ω
	}
	// This instead of Feature collection to reduce normalization calculations.
	output_holder := []*Feature{}
	for _, feature := range Ω.features {
		i, err := feature.GetProperty(property_key)
		if err != nil {
			panic(err)
			// Hard break because should function may be looking for nil or something
			return nil
		}
		if should_filter(i) {
			continue
		}
		output_holder = append(output_holder, feature)
	}
	fc := FeatureCollection{}
	// Ignore internal validate error because features must have been valid to be created.
	_ = fc.Append(output_holder...)
	return &fc
}

// Note that this function ignores missing strings.
// max_distance_from_centroid
func (Ω *FeatureCollection) GroupByProperties(property_keys ...string) FeatureCollections {
	if len(property_keys) == 0 {
		return nil
	}
	output_holder := map[string][]*Feature{}
	for _, feature := range Ω.features {
		a := ""
		for _, k := range property_keys {
			i, err := feature.GetProperty(k)
			if err != nil {
				return nil
			}
			b := i.([]byte)
			a += string(b)
		}
		if _, ok := output_holder[a]; !ok {
			output_holder[a] = []*Feature{}
		}
		output_holder[a] = append(output_holder[a], feature)
	}

	output := FeatureCollections{}

	for _, features := range output_holder {
		fc := FeatureCollection{}
		fc.Append(features...)
		output = append(output, &fc)
	}

	return output
}

func (Ω *FeatureCollection) FilterByMinimumArea(minimum_area_kilometers float64) *FeatureCollection {
	output_holder := []*Feature{}
	for _, ic := range Ω.features {
		if ic.Area() < minimum_area_kilometers {
			continue
		}
		output_holder = append(output_holder, ic)
	}
	fc := FeatureCollection{}
	// Ignore internal validate error because features must have been valid to be created.
	_ = fc.Append(output_holder...)
	return &fc
}

func (Ω FeatureCollection) MaxDistanceFromCentroid() float64 {

	polylabels := Points{}
	for _, fc := range Ω.features {
		polylabels = append(polylabels, fc.PolyLabel())
	}

	centroid := polylabels.Centroid()

	max := 0.0
	for _, p := range polylabels {
		distance := centroid.DistanceKilometers(p)
		if distance > max {
			max = distance
		}
	}

	return max
}

type CondenseMergePropertiesFunc func(a, b []byte) []byte

func (Ω FeatureCollection) Condense(merge_properties CondenseMergePropertiesFunc) *Feature {

	multipolygon := MultiPolygon{}
	properties := []byte{}
	for _, f := range Ω.features {
		properties = merge_properties(properties, f.properties)
		multipolygon = multipolygon.PushMultiPolygon(f.multiPolygon)
	}
	f := Feature{
		multiPolygon: multipolygon,
		properties: properties,
	}

	f.Normalize()

	return &f
}


/////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////


type FeatureCollections []*FeatureCollection

func (Ω FeatureCollections) FilterByMinimumArea(minimum_area_kilometers float64) (FeatureCollections) {
	out := FeatureCollections{}
	for _, ic := range Ω {
		if ic.Area() < minimum_area_kilometers {
			continue
		}
		out = append(out, ic)
	}
	return out
}

func (Ω FeatureCollections) PolyLabels() Points {
	labels := Points{}
	for _, ic := range Ω {
		labels = append(labels, ic.PolyLabel())
	}
	return labels
}

func (Ω FeatureCollections) Condense(merge_properties CondenseMergePropertiesFunc) (*FeatureCollection, error) {
	features := []*Feature{}
	for _, fc := range Ω {
		features = append(features, fc.Condense(merge_properties))
	}
	fc := FeatureCollection{}
	if err := fc.Append(features...); err != nil {
		return nil, err
	}
	return &fc, nil
}

// Sorter accepts two property maps and returns true if a is greater than b. If sorter is nil, the area will be used.
func (Ω FeatureCollections) DecimateClusters(minKm float64) FeatureCollections {

	a := FeatureCollections{}
	a = append(a, Ω...)
	// Slow process. Would be more efficient to cluster. Could grab bounding box seperate into quadrants.
	complete := false;

Restart:
	for complete == false {
		for i := range a {
			// Find nearest Point.
			nearest_distance := 0.0
			array_position := 0
			for k := range a {
				if k == i {
					continue
				}
				distance := a[i].PolyLabel().DistanceKilometers(a[k].PolyLabel())
				// If the distance between the two is less than the required minimum.
				if distance > minKm {
					continue
				}

				if nearest_distance == 0 || distance < nearest_distance {
					array_position = k
					nearest_distance = distance
				}
			}
			if nearest_distance > 0 {
				if a[i].Area() > a[array_position].Area() {
					a = append(a[:array_position], a[array_position+1:]...)
				} else {
					a = append(a[:i], a[i+1:]...)
				}
				continue Restart
			}
		}
		complete = true
	}


	return a
}

