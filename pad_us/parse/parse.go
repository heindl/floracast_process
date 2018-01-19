package main

import (
	"bitbucket.org/heindl/taxa/pad_us"
	"bitbucket.org/heindl/taxa/terra"
	"flag"
	"fmt"
	"github.com/elgs/gostrgen"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

const LogStatus = 1

var filteredNames = []string{
	"golf", "soccer", "recreation", "athletic", "softball",
	"baseball", "horse", "arts", "gym", "cemetery", "museum",
	"community center", "sports", "city high", "high school", "tennis", "pavilion",
	"skate park", "unknown", "elementary", "library",
}

func main() {
	in := flag.String("in", "/tmp/gap_analysis/ID/state.geojson", "Input json file")
	out := flag.String("out", "/tmp/gap_analysis/ID/areas", "Combined json directory")

	flag.Parse()

	if *in == "" || *out == "" {
		panic("input file and output directory required.")
	}

	processor := Processor{
		OutputDirectory: *out,
		Aggregated:      terra.FeatureCollection{},
		Stats:           map[string]int{},
	}

	if err := terra.ReadFeaturesFromGeoJSONFeatureCollectionFile(*in, processor.ReceiveFeature); err != nil {
		panic(err)
	}

	processor.Stats["Initial Filtered Total"] = processor.Aggregated.Count()

	filtered_unit_names := processor.Aggregated.FilterByProperty(func(i interface{}) bool {
		s := strings.ToLower(string(i.([]byte)))
		for _, f := range filteredNames {
			if strings.Contains(s, f) {
				return true
			}
		}
		// Filter out places with only numbers
		for _, r := range s {
			if unicode.IsLetter(r) {
				return false
			}
		}
		return true
	}, "Unit_Nm")

	processor.Stats["After Name Filter"] = filtered_unit_names.Count()

	name_grouped := filtered_unit_names.GroupByProperties("Unit_Nm", "Loc_Nm")

	processor.Stats["After Name Group"] = len(name_grouped)

	max_centroid_distance := 20.0
	above_centroid_distance := terra.FeatureCollections{}
	below_centroid_distance := terra.FeatureCollections{}
	for _, v := range name_grouped {
		// TODO: Explode and regroup those that are too large by Unit_Nm & Loc_Nm
		if v.Count() > 1 && v.MaxDistanceFromCentroid() > max_centroid_distance {
			above_centroid_distance = append(above_centroid_distance, v)
		} else {
			below_centroid_distance = append(below_centroid_distance, v)
		}
	}

	processor.Stats["After Centroid Distance Filter"] = len(below_centroid_distance)

	minimum_area_filtered := below_centroid_distance.FilterByMinimumArea(0.50)

	processor.Stats["After Minimum Area Filter"] = len(minimum_area_filtered)

	// TODO: Sort based on additional fields, particularly protected or access status.
	decimated_cluster := minimum_area_filtered.DecimateClusters(15)

	processor.Stats["After Cluster Decimation"] = len(decimated_cluster)

	processor.PrintStats()

	for _, v := range decimated_cluster {

		fname := fmt.Sprintf("%s/%.6f_%.6f.geojson", *out, v.PolyLabel().Latitude(), v.PolyLabel().Longitude())

		gj, err := v.GeoJSON()
		if err != nil {
			panic(err)
		}

		if err := ioutil.WriteFile(fname, gj, os.ModePerm); err != nil {
			panic(err)
		}
	}

	//gj, err := decimated_cluster.PolyLabels().GeoJSON()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(string(gj))

}

func filter_exists(shouldExist bool) func(interface{}) bool {
	return func(i interface{}) bool {
		s := string(i.([]byte))
		//fmt.Println(s)
		exists := (s != "" && s != "null" && s != "0")

		if shouldExist != exists {
			fmt.Println(s, exists, shouldExist != exists)
		}

		return shouldExist != exists
	}
}

func cast_string(i interface{}) string {
	return string(i.([]byte))
}

type Processor struct {
	sync.Mutex
	OutputDirectory        string
	PublicAccessClosed     int
	Aggregated             terra.FeatureCollection
	PublicAccessRestricted int
	GolfCourse             int
	PublicAccessUnknown    int
	UnassignedIUCNCategory int
	Total                  int
	EmptyAreas             int
	LocalParks             int
	MarineProtectedArea    int
	Stats                  map[string]int
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
	fmt.Println("Empty Areas", Ω.EmptyAreas)
	fmt.Println("UnassignedIUCNCategory", Ω.UnassignedIUCNCategory)
	fmt.Println("MarineProtectedArea", Ω.MarineProtectedArea)
	fmt.Println("LocalParks", Ω.LocalParks)

	for k, v := range Ω.Stats {
		fmt.Println(k, v)
	}
}

func (Ω *Processor) ShouldSaveProtectedArea(feature *terra.Feature) bool {

	if !feature.Valid() {
		Ω.EmptyAreas += 1
		return false
	}

	pa := pad_us.ProtectedArea{}
	if err := feature.GetProperties(&pa); err != nil {
		panic(err)
	}

	switch pa.Access {
	case pad_us.PublicAccessClosed:
		Ω.PublicAccessClosed += 1
		return false
	case pad_us.PublicAccessRestricted:
		Ω.PublicAccessRestricted += 1
		// In Alabama, this includes Wildlife Refuges and WMAs.
		//return false
	case pad_us.PublicAccessUnknown:
		Ω.PublicAccessUnknown += 1
		// There are so many of these that look valid, we should ignore this.
		//return false
	}

	if _, ok := pad_us.GAPStatusDefinitions[pa.GAPStatusCode]; !ok {
		fmt.Println(fmt.Sprintf(`GapStatus("%s"): "%s", | %s`, pa.GAPStatusCode, pa.DGAPSts, pa.UnitNm))
		return false
	}

	if _, ok := pad_us.AreaIUCNCategoryDefinitions[pa.IUCNCategory]; !ok {
		fmt.Println(fmt.Sprintf(`IUCNCategory("%s"): "%s", | %s`, pa.IUCNCategory, pa.DIUCNCat, pa.UnitNm))
		return false
	}

	if _, ok := pad_us.CategoryDefinitions[pa.Category]; !ok {
		fmt.Println(fmt.Sprintf(`Category("%s"): "%s", | %s`, pa.Category, pa.DCategory, pa.UnitNm))
		return false
	}

	if _, ok := pad_us.OwnerTypeDefinitions[pa.OwnerType]; !ok {
		fmt.Println(fmt.Sprintf(`OwnerType("%s"): "%s", | %s`, pa.OwnerType, pa.DOwnType, pa.UnitNm))
		return false
	}

	if _, ok := pad_us.OwnerNameDefinitions[pa.OwnerName]; !ok {
		fmt.Println(fmt.Sprintf(`OwnerName("%s"): "%s", | %s`, pa.OwnerName, pa.DOwnName, pa.UnitNm))
		return false
	}

	if _, ok := pad_us.ManagerNameDefinitions[pa.ManagerName]; !ok {
		fmt.Println(fmt.Sprintf(`ManagerName("%s"): "%s", | %s`, pa.ManagerName, pa.DMangNam, pa.UnitNm))
		return false
	}

	if _, ok := pad_us.ManagerTypeDefinitions[pa.ManagerType]; !ok {
		fmt.Println(fmt.Sprintf(`ManagerType("%s"): "%s", | %s`, pa.ManagerType, pa.DMangTyp, pa.UnitNm))
		return false
	}

	if _, ok := pad_us.DesignationDefinitions[pa.Designation]; !ok {
		if pa.Designation == pad_us.Designation("LP") {
			Ω.LocalParks += 1
			return false
		}
		if pa.Designation == pad_us.Designation("MPA") {
			Ω.MarineProtectedArea += 1
			return false
		}
		fmt.Println(fmt.Sprintf(`Designation("%s"): "%s", | %s`, pa.Designation, pa.DDesTp, pa.UnitNm))
		return false
	}

	return true
}

func (Ω *Processor) GetID(nf *terra.Feature) string {
	pa := pad_us.ProtectedArea{}
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

func (Ω *Processor) ReceiveFeature(nf *terra.Feature) error {

	Ω.Total += 1

	if !Ω.ShouldSaveProtectedArea(nf) {
		return nil
	}

	Ω.Lock()
	defer Ω.Unlock()

	if err := Ω.Aggregated.Append(nf); err != nil {
		return err
	}

	return nil
}
