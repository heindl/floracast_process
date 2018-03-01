package geo

import (
	"github.com/dropbox/godropbox/errors"
	"github.com/smira/go-point-clustering"
)

type FeatureCollections []*FeatureCollection

func (Ω FeatureCollections) FilterByMinimumArea(minimum_area_kilometers float64) FeatureCollections {
	out := FeatureCollections{}
	for _, ic := range Ω {
		if ic.Area() < minimum_area_kilometers {
			continue
		}
		out = append(out, ic)
	}
	return out
}

func (Ω FeatureCollections) PolyLabels() (Points, error) {
	labels := Points{}
	for _, ic := range Ω {
		p, err := ic.PolyLabel()
		if err != nil {
			return nil, err
		}
		labels = append(labels, p)
	}
	return labels, nil
}

func (Ω FeatureCollections) Condense(mergeProperties CondenseMergePropertiesFunc) (*FeatureCollection, error) {
	features := []*Feature{}
	for _, fc := range Ω {
		condensed, err := fc.Condense(mergeProperties)
		if err != nil {
			return nil, err
		}
		features = append(features, condensed)
	}
	fc := FeatureCollection{}
	if err := fc.Append(features...); err != nil {
		return nil, err
	}
	return &fc, nil
}

func (Ω FeatureCollections) DecimateClusters(radiusKm float64) (FeatureCollections, error) {

	polyLabels, err := Ω.PolyLabels()
	if err != nil {
		return nil, err
	}

	clusterPointList := cluster.PointList{}
	for _, p := range polyLabels {
		clusterPointList = append(clusterPointList, cluster.Point{p.Longitude(), p.Latitude()})
	}
	clusters, _ := cluster.DBScan(clusterPointList, radiusKm, 1)

	res := FeatureCollections{}
	for _, clstr := range clusters {
		col, err := largestWithinCluster(Ω, clstr.Points)
		if err != nil {
			return nil, err
		}
		res = append(res, col)
	}

	return res, nil

}

func largestWithinCluster(fcs FeatureCollections, positions []int) (*FeatureCollection, error) {
	res := &FeatureCollection{}
	for _, position := range positions {
		if res == nil || fcs[position].Area() > res.Area() {
			res = fcs[position]
		}
	}
	if res == nil {
		return nil, errors.New("Expected Largest Cluster to not be nil")
	}
	return res, nil
}
