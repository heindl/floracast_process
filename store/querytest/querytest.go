package main

import (
	"bitbucket.org/heindl/taxa/store"
	"golang.org/x/net/context"
	"google.golang.org/genproto/googleapis/type/latlng"
	"time"
)

func main() {
	client, err := store.NewLiveFirestore()
	if err != nil {
		panic(err)
	}

	//p := store.Prediction{
	//	FormattedDate: "20170806",
	//	Location: latlng.LatLng {
	//		Latitude: 43.788655,
	//		Longitude: -75.097508,
	//	},
	//	Month: time.August,
	//	PredictionValue: 0.06885252892971039,
	//	PercentileOverAllTaxonPredictions: 0.7,
	//	PercentileOverAllTaxaPredictionsForDay: 0.6,
	//	INaturalistTaxonID: store.INaturalistTaxonID("473935"),
	//}

	o := store.Occurrence{
		FormattedDate: "20170806",
		TargetID:      "473935",
		DataSourceID:  store.DataSourceTypeGBIF,
		Location: latlng.LatLng{
			Latitude:  43.788655,
			Longitude: -75.097508,
		},
		Month:   time.August,
		TaxonID: store.INaturalistTaxonID("473935"),
	}

	if _, _, err := client.Collection("OccurrenceCount").Add(context.Background(), o); err != nil {
		panic(err)
	}

	//docs, err := client.Collection("OccurrenceCount").Where("INaturalistTaxonID", "==", "143393").
	////Select("Location", ).
	////docs, err := client.Collection("OccurrenceCount").
	////Where("S2CellIDs.1_10001", "==", true).
	//Documents(context.Background()).GetAll()
	//if err != nil {
	//	panic(err)
	//}
	//res := store.OccurrenceCount{}
	//for _, doc := range docs {
	//	o := store.Occurrence{}
	//	if err := doc.DataTo(&o); err != nil {
	//		panic(err)
	//	}
	//	res = append(res, o)
	//}
	//fmt.Println("length", len(res))
	//fmt.Println(utils.JsonOrSpew(res))
}
