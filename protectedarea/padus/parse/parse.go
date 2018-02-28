package main

import (
	"bitbucket.org/heindl/process/protectedarea/padus"
	"bitbucket.org/heindl/process/terra/geo"
	"bitbucket.org/heindl/process/utils"
	"flag"
	"fmt"
	"github.com/elgs/gostrgen"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
)

const LogStatus = 1

var flagsToFilter = []string{
	"golf", "soccer", "recreation", "athletic", "softball",
	"baseball", "horse", "arts", "gym", "cemetery", "museum",
	"community center", "sports", "city high", "high school", "tennis", "pavilion",
	"skate park", "unknown", "elementary", "library",
}

const maxCentroidDistance = 20.0

func main() {
	in := flag.String("in", "/tmp/gap_analysis/ID/state.geojson", "Input json file")
	out := flag.String("out", "/tmp/gap_analysis/ID/areas", "Combined json directory")

	flag.Parse()

	if *in == "" || *out == "" {
		panic("input file and output directory required.")
	}

	processor := Processor{
		OutputDirectory: *out,
		Aggregated:      geo.FeatureCollection{},
		Stats:           map[string]int{},
	}

	grouped, err := processor.ReadGroupAndFilter(*in)
	if err != nil {
		panic(err)
	}

	// Filter by distance from cluster centroid to avoid widely dispersed parks.
	validCentroidDistance := geo.FeatureCollections{}
	for _, v := range grouped {
		if v.Count() == 1 {
			validCentroidDistance = append(validCentroidDistance, v)
			continue
		}
		// TODO: Explode and regroup those that are too large by Unit_Nm & Loc_Nm
		maxDistance, err := v.MaxDistanceFromCentroid()
		if err != nil {
			panic(err)
		}
		if maxDistance <= maxCentroidDistance {
			validCentroidDistance = append(validCentroidDistance, v)
		}
	}

	processor.Stats["After Centroid Distance Filter"] = len(validCentroidDistance)

	minimum_area_filtered := validCentroidDistance.FilterByMinimumArea(0.50)

	processor.Stats["After Minimum Area Filter"] = len(minimum_area_filtered)

	// TODO: Sort based on additional fields, particularly protected or access status.
	decimatedClusters, err := minimum_area_filtered.DecimateClusters(15)
	if err != nil {
		panic(err)
	}

	processor.Stats["After Cluster Decimation"] = len(decimatedClusters)

	processor.PrintStats()

	for _, v := range decimatedClusters {
		polyLabel, err := v.PolyLabel()
		if err != nil {
			panic(err)
		}

		fname := fmt.Sprintf("%s/%.6f_%.6f.geojson", *out, polyLabel.Latitude(), polyLabel.Longitude())

		gj, err := v.GeoJSON()
		if err != nil {
			panic(err)
		}

		if err := ioutil.WriteFile(fname, gj, os.ModePerm); err != nil {
			panic(err)
		}
	}

}

type Processor struct {
	sync.Mutex
	OutputDirectory string
	Aggregated      geo.FeatureCollection
	*counts
	Stats map[string]int
}

type counts struct {
	PublicAccessClosed     int
	PublicAccessRestricted int
	GolfCourse             int
	PublicAccessUnknown    int
	UnassignedIUCNCategory int
	Total                  int
	EmptyAreas             int
	LocalParks             int
	MarineProtectedArea    int
}

func (Ω *Processor) ReadGroupAndFilter(filepath string) (geo.FeatureCollections, error) {
	if err := geo.ReadFeaturesFromGeoJSONFeatureCollectionFile(filepath, Ω.ReceiveFeature); err != nil {
		return nil, err
	}

	Ω.Stats["Initial Filtered Total"] = Ω.Aggregated.Count()

	filteredUnitNames, err := Ω.Aggregated.FilterByProperty(func(i interface{}) bool {
		s := strings.ToLower(string(i.([]byte)))
		if utils.WordInArrayIsASubstring(s, flagsToFilter) || utils.StringContainsOnlyNumbers(s) {
			return true
		}
		return false
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

func (Ω *Processor) PrintStats() {

	if LogStatus == 0 {
		return
	}

	fmt.Println("Total", Ω.Total)
	fmt.Println("PublicAccessClosed", Ω.PublicAccessClosed)
	fmt.Println("PublicAccessRestricted", Ω.PublicAccessRestricted)
	fmt.Println("PublicAccessUnknown", Ω.PublicAccessUnknown)
	fmt.Println("GolfCourse", Ω.GolfCourse)
	fmt.Println("IsZero Areas", Ω.EmptyAreas)
	fmt.Println("UnassignedIUCNCategory", Ω.UnassignedIUCNCategory)
	fmt.Println("MarineProtectedArea", Ω.MarineProtectedArea)
	fmt.Println("LocalParks", Ω.LocalParks)

	for k, v := range Ω.Stats {
		fmt.Println(k, v)
	}
}

type fieldValidator interface {
	Valid() bool
}

func (Ω *Processor) ShouldSaveProtectedArea(feature *geo.Feature) (bool, error) {

	if !feature.Valid() {
		Ω.EmptyAreas += 1
		return false, nil
	}

	pa := padus.ProtectedArea{}
	if err := feature.GetProperties(&pa); err != nil {
		return false, err
	}

	switch pa.Access {
	case padus.PublicAccessClosed:
		Ω.PublicAccessClosed += 1
		return false, nil
	case padus.PublicAccessRestricted: // In Alabama, this includes Wildlife Refuges and WMAs.
		Ω.PublicAccessRestricted += 1
	case padus.PublicAccessUnknown: // There are so many of these that look valid, we should ignore this.
		Ω.PublicAccessUnknown += 1
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
			//	Ω.LocalParks += 1
			//	return false
			//}
			//if pa.Designation == pad_us.Designation("MPA") {
			//	Ω.MarineProtectedArea += 1
			//	return false
			//}
			fmt.Println(fmt.Sprintf(`GapStatus("%s"): "%s", | %s`, pa.GAPStatusCode, pa.DGAPSts, pa.UnitNm))
			return false, nil
		}
	}

	return true, nil
}

func (Ω *Processor) GetID(nf *geo.Feature) string {
	pa := padus.ProtectedArea{}
	if err := nf.GetProperties(&pa); err != nil {
		panic(err)
	}

	if pa.WDPACd != 0 {
		return "wdpa_" + strconv.Itoa(int(pa.WDPACd))
	}

	if pa.SourcePAI != "" {
		return "pai_" + pa.SourcePAI
	}

	rand_id, err := gostrgen.RandGen(20, gostrgen.Lower|gostrgen.Digit, "", "")
	if err != nil {
		panic(err)
	}
	return "unidentified_" + rand_id
}

func (Ω *Processor) ReceiveFeature(nf *geo.Feature) error {

	Ω.Total += 1

	shouldSave, err := Ω.ShouldSaveProtectedArea(nf)
	if err != nil {
		return err
	}
	if !shouldSave {
		return nil
	}

	Ω.Lock()
	defer Ω.Unlock()

	return Ω.Aggregated.Append(nf)
}
