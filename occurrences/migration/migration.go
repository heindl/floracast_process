package main

import (
	"bitbucket.org/heindl/taxa/store"
	"context"
	"gopkg.in/tomb.v2"
)

func main() {


	//0 13343.39119734705 [90.0000000, 90.0000000] 2.2214414690791835
	//1 6671.695598673532 [45.0000000, 45.0000000] 0.555360367269797
	//2 3373.639861740696 [23.9624890, 22.6198649] 0.137437725222794
	//3 1517.9232350633238 [10.9521370, 10.6196553] 0.027287207440910972
	//4 764.837489493313 [5.5770774, 5.0598688] 0.00684499781515081
	//5 369.35866309005206 [2.7094999, 2.4604145] 0.0015863233255728505
	//6 181.362198552915 [1.3345051, 1.2121322] 0.0003812234884628829
	//7 90.77020381049476 [0.6691799, 0.6014717] 9.530126567965143e-05
	//8 45.59232713177376 [0.3358469, 0.3018932] 2.4063379259794675e-05
	//9 22.801382443667002 [0.1680417, 0.1506583] 6.01560690771692e-06
	//10 11.401984368274213 [0.0840504, 0.0752570] 1.503866974949068e-06
	//11 5.701314012680624 [0.0420326, 0.0376104] 3.759620843102548e-07
	//12 2.8507372805005597 [0.0210181, 0.0188007] 9.398991886986995e-08
	//13 1.4255701565429497 [0.0105103, 0.0094015] 2.3504752478269158e-08
	//14 0.7127900840274252 [0.0052553, 0.0047005] 5.876178539415875e-09


	//for i:=0;i<30;i++ {
	//	cellID := s2.CellIDFromLatLng(s2.LatLngFromDegrees(38.6530169, -90.3135084)).Parent(i)
	//	cell := s2.CellFromCellID(cellID)
	//	lo := cell.RectBound().Lo()
	//	high := cell.RectBound().Hi()
	//	distance := geo.NewPoint(lo.Lat.Degrees(), lo.Lng.Degrees()).GreatCircleDistance(geo.NewPoint(high.Lat.Degrees(), high.Lng.Degrees()))
	//	fmt.Println(i, distance, cell.RectBound().Size(), cell.RectBound().Area())
	//}
	//for i:=0;i<30;i++ {
	//	fmt.Println(i, s2.CellFromCellID(cellID.Parent(i)).ExactArea())
	//}
	ts, err := store.NewTaxaStore()
	if err != nil {
		panic(err)
	}

	cxt := context.Background()

	taxa, err := ts.ReadTaxa(cxt)
	if err != nil {
		panic(err)
	}
	tmb := tomb.Tomb{}
	tmb.Go(func() error{
		for _, _taxon := range taxa {
			taxon := _taxon
			if taxon.ID != store.TaxonID("58682") {
				continue
			}
			tmb.Go(func() error {
				occurrences, err := ts.GetOccurrences(cxt, taxon.ID)
				if err != nil {
					return err
				}
				for _, o := range occurrences {
					if _, err := ts.UpsertOccurrence(cxt, o); err != nil {
						return err
					}
					//if err := ts.IncrementTaxonEcoRegion(cxt, o.TaxonID, o.EcoRegion); err != nil {
					//	return err
					//}
				}
				return nil
			})
		}
		return nil
	})
	if err := tmb.Wait(); err != nil {
		panic(err)
	}

}