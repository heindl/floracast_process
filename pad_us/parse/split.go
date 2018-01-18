package main

import (
	"flag"
	"bitbucket.org/heindl/taxa/terra"
	"sync"
	"bitbucket.org/heindl/taxa/pad_us"
	"fmt"
	"strconv"
	"github.com/elgs/gostrgen"
	"strings"
)

const PrintName = ""


func main() {
	in := flag.String("in", "/tmp/gap_analysis/CA/state.geojson", "Input json file")
	out := flag.String("out", "/tmp/gap_analysis/CA/areas", "Combined json directory")

	flag.Parse()

	if *in == "" || *out == "" {
		panic("input file and output directory required.")
	}

	processor := Processor{
		OutputDirectory: *out,
		Aggregated: terra.FeatureCollection{},
	}

	if err := terra.ReadGeoJSONFeatureCollectionFile(*in, processor.ReceiveFeature); err != nil {
		panic(err)
	}

	processor.PrintStats()

	fmt.Println("Initial Filtered Total", processor.Aggregated.Count())

	grouped_by_wdpa := processor.Aggregated.FilterByProperty(filter_exists(true), "WDPA_Cd").GroupByProperty(cast_string, "Unit_Nm")
	grouped_by_pai := processor.Aggregated.FilterByProperty(filter_exists(true), "Source_PAI").GroupByProperty(cast_string,  "Unit_Nm")
	grouped_by_undefined := processor.Aggregated.
		FilterByProperty(filter_exists(false), "Source_PAI").
		FilterByProperty(filter_exists(false), "WDPA_Cd").
		GroupByProperty(cast_string, "Unit_Nm")

		name_grouped := terra.FeatureCollections{}
		name_grouped = append(name_grouped, grouped_by_wdpa...)
	name_grouped = append(name_grouped, grouped_by_pai...)
	name_grouped = append(name_grouped, grouped_by_undefined...)

	fmt.Println("After Name Group", len(name_grouped))

	max_centroid_distance := 20.0
	above_centroid_distance := terra.FeatureCollections{}
	below_centroid_distance := terra.FeatureCollections{}
	for _, v := range name_grouped {
		if v.Count() > 1 && v.MaxDistanceFromCentroid() > max_centroid_distance {
			above_centroid_distance = append(above_centroid_distance, v)
		} else {
			below_centroid_distance = append(below_centroid_distance, v)
		}
	}

	fmt.Println("After Centroid Distance", len(below_centroid_distance))

	minimum_area_filtered := below_centroid_distance.FilterByMinimumArea(1)

	fmt.Println("After Minimum Area", len(minimum_area_filtered))

	decimated := minimum_area_filtered.DecimateClusters(25)

	fmt.Println("After Decimation", len(decimated))

	//gj, err := decimated.PolyLabels().GeoJSON()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(string(gj))



	return

}

func filter_exists(shouldExist bool) (func(interface{}) bool) {
	return func(i interface{}) bool {
		s := string(i.([]byte))
		exists := s != "" && s != "null" && s != "0"
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
	Aggregated            terra.FeatureCollection
	PublicAccessRestricted int
	GolfCourse int
	PublicAccessUnknown    int
	UnassignedIUCNCategory int
	UnknownDesignation int
	Total                  int
	EmptyAreas             int
}

func (Ω *Processor) PrintStats() {
	fmt.Println("Total", Ω.Total)
	fmt.Println("PublicAccessClosed", Ω.PublicAccessClosed)
	fmt.Println("PublicAccessRestricted", Ω.PublicAccessRestricted)
	fmt.Println("PublicAccessUnknown", Ω.PublicAccessUnknown)
	fmt.Println("GolfCourse", Ω.GolfCourse)
	fmt.Println("Empty Areas", Ω.EmptyAreas)
	fmt.Println("UnassignedIUCNCategory", Ω.UnassignedIUCNCategory)
	fmt.Println("UnknownDesignation", Ω.UnknownDesignation)
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

	if pa.UnitNm == "" {
		panic("no unit name")
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

	if strings.Contains(strings.ToLower(pa.UnitNm), "golf") {
		Ω.GolfCourse += 1
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
		return "wdpa_"+strconv.Itoa(int(pa.WDPACd))
	}

	if pa.SourcePAI != "" {
		return "pai_"+pa.SourcePAI
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
