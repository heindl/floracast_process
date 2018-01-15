package main

import (
	"encoding/json"
	"flag"
	"bitbucket.org/heindl/taxa/terra"
	"sync"
	"bitbucket.org/heindl/taxa/pad_us"
	"github.com/dropbox/godropbox/errors"
	"fmt"
	"bitbucket.org/heindl/taxa/utils"
)

func main() {
	in := flag.String("in", "/tmp/gap_analysis/CA/state.geojson", "Input json file")
	out := flag.String("out", "/tmp/gap_analysis/CA/areas", "Combined json directory")

	flag.Parse()

	if *in == "" || *out == "" {
		panic("input file and output directory required.")
	}

	processor := Processor{
		OutputDirectory: *out,
		NameCount: map[string]map[int64][]*Feature{},
		Cache: map[int64]*Feature{},
		Unidentified: []*Feature{},
	}

	if err := terra.ReadGeoJSONFeatureCollectionFile(*in, processor.ReceiveFeature); err != nil {
		panic(err)
	}

	fmt.Println("NAME COUNTS")
	for name, ids := range processor.NameCount {
		if len(ids) > 1 {
			for id, features := range ids {
				for _, f := range features {
					fmt.Println(name, id, f.MultiPolygon.ReferencePoints())
				}
			}
		}
		//for id, count := range ids {
		//	if count > 2 {
		//		fmt.Println(id, name, count)
		//	}
		//}
	}

	return

	//fmt.Println("MULTIPOLYGONS")
	//for id, feature := range processor.Cache {
	//	fmt.Println(id, len(feature.MultiPolygon))
	//}

	for _, u := range processor.Unidentified {
		if u.Parameters.Access == "XA" {
			continue
		}
		fmt.Println(utils.JsonOrSpew(u))
	}

	//fmt.Println("UNIDENTIFIED: ", len(processor.Unidentified))

	//for k, ofc := range outputFeatureCollections {
	//	b, err := json.Marshal(ofc)
	//	if err != nil {
	//		panic(err)
	//	}
	//	fname := k + ".feature_collection.geojson"
	//	if err := ioutil.WriteFile(path.Join(*out, fname), b, 0700); err != nil {
	//		panic(err)
	//	}
	//}

}

type Feature struct {
	Parameters *pad_us.ProtectedArea
	MultiPolygon terra.MultiPolygon
}

type Processor struct {
	sync.Mutex
	OutputDirectory string
	NameCount map[string]map[int64][]*Feature
	Cache map[int64]*Feature
	Unidentified []*Feature
}

func (Ω *Processor) ReceiveFeature(encoded_properties []byte, multipolygon terra.MultiPolygon) error {

	if multipolygon.Empty() {
		return nil
	}



	pa := pad_us.ProtectedArea{}
	if err := json.Unmarshal(encoded_properties, &pa); err != nil {
		return errors.Wrap(err, "could not unmarshal properties")
	}

	nf := Feature{
			Parameters:   &pa,
			MultiPolygon: multipolygon,
		}

	Ω.Lock()
	defer Ω.Unlock()

	if pa.WDPACd == 0 {
		Ω.Unidentified = append(Ω.Unidentified, &nf)
		return nil
	}
	if nf.Parameters.UnitNm == "" {
		fmt.Println("no unit name", nf.Parameters.LocalName)
		return nil
	}

	if _, ok := Ω.Cache[pa.WDPACd]; !ok {
		Ω.Cache[pa.WDPACd] = &nf
	} else {
		Ω.Cache[pa.WDPACd].MultiPolygon = Ω.Cache[pa.WDPACd].MultiPolygon.PushMultiPolygon(nf.MultiPolygon)
	}

	if _, ok := Ω.NameCount[nf.Parameters.UnitNm]; !ok {
		Ω.NameCount[nf.Parameters.UnitNm] = map[int64][]*Feature{
			pa.WDPACd: {&nf},
		}
	} else {
		if _, ok := Ω.NameCount[nf.Parameters.UnitNm][pa.WDPACd]; !ok {
			Ω.NameCount[nf.Parameters.UnitNm][pa.WDPACd] = []*Feature{&nf}
		} else {
			Ω.NameCount[nf.Parameters.UnitNm][pa.WDPACd] = append(Ω.NameCount[nf.Parameters.UnitNm][pa.WDPACd], &nf)
		}
	}
	return nil
}
