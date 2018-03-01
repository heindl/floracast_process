package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"

	"bitbucket.org/heindl/process/protectedarea/padus"
	"bitbucket.org/heindl/process/terra/geo"
	"bitbucket.org/heindl/process/utils"
	"github.com/dropbox/godropbox/errors"
	"gopkg.in/tomb.v2"
)

// Processor parses, groups, filters and writes PAD-US GeoJSON.
type Processor interface {
	ProcessFeatureCollections() (geo.FeatureCollections, *metrics, error)
	WriteCollections(collections geo.FeatureCollections) error
}

// NewProcessor validates and initiates a processor to parse, group, filter and write PAD-US GeoJSON.
// Input and Output directories must be specified.
func NewProcessor(inpath, outpath string) (Processor, error) {
	inpath = strings.TrimSpace(inpath)
	outpath = strings.TrimSpace(outpath)
	if inpath == "" || outpath == "" {
		return nil, errors.Newf("Invalid directories [%s, %s]", inpath, outpath)
	}

	return &orchestrator{
		inPath:  inpath,
		outPath: outpath,
		metrics: &metrics{
			Stats: map[string]int{},
		},
	}, nil
}

type orchestrator struct {
	sync.Mutex
	inPath     string
	outPath    string
	aggregated geo.FeatureCollection
	*metrics
}

type metrics struct {
	Stats                  map[string]int `json:""`
	PublicAccessClosed     int            `json:""`
	PublicAccessRestricted int            `json:""`
	//golfCourse             int
	PublicAccessUnknown int `json:""`
	//unassignedIUCNCategory int
	Total      int `json:""`
	EmptyAreas int `json:""`
	//localParks          int
	//marineProtectedArea int
}

// ProcessFeatureCollections reads, parses, groups, and filters PAD-US GeoJSON
func (Ω *orchestrator) ProcessFeatureCollections() (geo.FeatureCollections, *metrics, error) {

	nameGrouped, err := Ω.readGroupAndFilter()
	if err != nil {
		return nil, nil, err
	}

	// Filter by distance from cluster centroid to avoid widely dispersed parks.
	centroidDistanceFiltered, err := Ω.filterByCentroidDistance(nameGrouped)
	if err != nil {
		return nil, nil, err
	}

	Ω.Stats["After Centroid Distance Filter"] = len(centroidDistanceFiltered)

	minimumAreaFiltered := centroidDistanceFiltered.FilterByMinimumArea(0.50)

	Ω.Stats["After Minimum Area Filter"] = len(minimumAreaFiltered)

	// TODO: Sort based on additional fields, particularly protected or access status.
	decimatedClusters, err := minimumAreaFiltered.DecimateClusters(clusterDecimationKm)
	if err != nil {
		return nil, nil, err
	}

	Ω.Stats["After Cluster Decimation"] = len(decimatedClusters)

	return decimatedClusters, Ω.metrics, nil
}

// WriteCollections writes grouped feature collections to a geojson file
// The collection polylabel is used for the filename
func (Ω *orchestrator) WriteCollections(collections geo.FeatureCollections) error {
	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _, 𝝨 := range collections {
			col := 𝝨
			tmb.Go(func() error {
				polyLabel, err := col.PolyLabel()
				if err != nil {
					return err
				}
				filePath := path.Join(Ω.outPath, fmt.Sprintf("%.6f_%.6f.geojson", polyLabel.Latitude(), polyLabel.Longitude()))
				gj, err := col.GeoJSON()
				if err != nil {
					return err
				}
				if err := ioutil.WriteFile(filePath, gj, os.ModePerm); err != nil {
					return errors.Wrapf(err, "Could not write GeoJSON [%s]", filePath)
				}
				return nil
			})
		}
		return nil
	})
	return tmb.Wait()
}

func (Ω *orchestrator) readGroupAndFilter() (geo.FeatureCollections, error) {

	if err := geo.ReadFeaturesFromGeoJSONFeatureCollectionFile(Ω.inPath, Ω.receiveFeature); err != nil {
		return nil, err
	}

	Ω.Stats["Initial Filtered Total"] = Ω.aggregated.Count()

	filteredUnitNames, err := Ω.aggregated.FilterByProperty(func(i interface{}) bool {
		s := strings.ToLower(string(i.([]byte)))
		return utils.WordInArrayIsASubstring(s, flagsToFilter) || utils.StringContainsOnlyNumbers(s)
	}, "Unit_Nm")
	if err != nil {
		return nil, err
	}

	Ω.Stats["After Name Filter"] = filteredUnitNames.Count()

	nameGrouped, err := filteredUnitNames.GroupByProperties("Unit_Nm", "Loc_Nm")
	if err != nil {
		return nil, err
	}

	Ω.Stats["After Name Group"] = len(nameGrouped)

	return nameGrouped, nil
}

func (Ω *orchestrator) filterByCentroidDistance(collections geo.FeatureCollections) (geo.FeatureCollections, error) {
	// Filter by distance from cluster centroid to avoid widely dispersed parks.
	res := geo.FeatureCollections{}
	tmb := tomb.Tomb{}
	lock := sync.Mutex{}
	tmb.Go(func() error {
		for _, 𝝨 := range collections {
			col := 𝝨
			tmb.Go(func() error {
				if col.Count() == 1 {
					lock.Lock()
					defer lock.Unlock()
					res = append(res, col)
					return nil
				}
				// TODO: Explode and regroup those that are too large by Unit_Nm & Loc_Nm
				maxDistance, err := col.MaxDistanceFromCentroid()
				if err != nil {
					return err
				}
				lock.Lock()
				defer lock.Unlock()
				if maxDistance <= maxCentroidDistance {
					res = append(res, col)
				} else {
					cols, err := col.Explode()
					if err != nil {
						return err
					}
					res = append(res, cols...)
				}
				return nil
			})
		}
		return nil
	})
	if err := tmb.Wait(); err != nil {
		return nil, err
	}
	return res, nil
}

type fieldValidator interface {
	Valid() bool
}

func (Ω *orchestrator) isFeatureValid(feature *geo.Feature) (bool, error) {

	if !feature.Valid() {
		Ω.EmptyAreas++
		return false, nil
	}

	pa := padus.ProtectedArea{}
	if err := feature.GetProperties(&pa); err != nil {
		return false, err
	}

	switch pa.Access {
	case padus.PublicAccessRestricted: // In Alabama, this includes Wildlife Refuges and WMAs.
		Ω.PublicAccessRestricted++
	case padus.PublicAccessUnknown: // There are so many of these that look valid, we should ignore this.
		Ω.PublicAccessUnknown++
	case padus.PublicAccessClosed:
		Ω.PublicAccessClosed++
		return false, nil
	}

	for _, field := range []struct {
		v fieldValidator
		s string
	}{
		{v: pa.GAPStatusCode, s: fmt.Sprintf(`GapStatus("%s"): "%s", | %s`, pa.GAPStatusCode, pa.DGAPSts, pa.UnitNm)},
		{v: pa.IUCNCategory, s: fmt.Sprintf(`IUCNCategory("%s"): "%s", | %s`, pa.IUCNCategory, pa.DIUCNCat, pa.UnitNm)},
		{v: pa.Category, s: fmt.Sprintf(`Category("%s"): "%s", | %s`, pa.Category, pa.DCategory, pa.UnitNm)},
		{v: pa.OwnerType, s: fmt.Sprintf(`OwnerType("%s"): "%s", | %s`, pa.OwnerType, pa.DOwnType, pa.UnitNm)},
		{v: pa.OwnerName, s: fmt.Sprintf(`OwnerName("%s"): "%s", | %s`, pa.OwnerName, pa.DOwnName, pa.UnitNm)},
		{v: pa.ManagerName, s: fmt.Sprintf(`ManagerName("%s"): "%s", | %s`, pa.ManagerName, pa.DMangNam, pa.UnitNm)},
		{v: pa.ManagerType, s: fmt.Sprintf(`ManagerType("%s"): "%s", | %s`, pa.ManagerType, pa.DMangTyp, pa.UnitNm)},
		{v: pa.Designation, s: fmt.Sprintf(`Designation("%s"): "%s", | %s`, pa.Designation, pa.DDesTp, pa.UnitNm)},
	} {
		if !field.v.Valid() {
			//if pa.Designation == pad_us.Designation("LP") {
			//	Ω.localParks += 1
			//	return false
			//}
			//if pa.Designation == pad_us.Designation("MPA") {
			//	Ω.marineProtectedArea += 1
			//	return false
			//}
			//fmt.Println(field.s)
			return false, nil
		}
	}

	return true, nil
}

//func (Ω *orchestrator) getID(nf *geo.Feature) string {
//	pa := padus.ProtectedArea{}
//	if err := nf.GetProperties(&pa); err != nil {
//		panic(err)
//	}
//
//	if pa.WDPACd != 0 {
//		return "wdpa_" + strconv.Itoa(int(pa.WDPACd))
//	}
//
//	if pa.SourcePAI != "" {
//		return "pai_" + pa.SourcePAI
//	}
//
//	randID, err := gostrgen.RandGen(20, gostrgen.Lower|gostrgen.Digit, "", "")
//	if err != nil {
//		panic(err)
//	}
//	return "unidentified_" + randID
//}

func (Ω *orchestrator) receiveFeature(nf *geo.Feature) error {

	Ω.Total++

	shouldSave, err := Ω.isFeatureValid(nf)
	if err != nil {
		return err
	}
	if !shouldSave {
		return nil
	}

	Ω.Lock()
	defer Ω.Unlock()

	return Ω.aggregated.Append(nf)
}
