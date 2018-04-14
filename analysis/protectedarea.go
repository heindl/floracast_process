package main

//
//type protectedAreas struct {
//	floraStore store.FloraStore
//	cxt        context.Context
//}
//
//func (Ω *protectedAreas) PrintCoordinates() {
//	colRef, err := Ω.floraStore.FirestoreCollection(store.CollectionProtectedAreas)
//	if err != nil {
//		panic(err)
//	}
//
//	iter := colRef.Documents(Ω.cxt)
//
//	type Coords struct {
//		Lat float64 `json:"lat"`
//		Lng float64 `json:"lng"`
//	}
//
//	res := []Coords{}
//
//	for {
//		snap, err := iter.Next()
//		if err == iterator.Done {
//			break
//		}
//		if err != nil {
//			panic(err)
//		}
//		coords := strings.Split(strings.Replace(snap.Ref.ID, "|", ".", -1), "_")
//		lat, _ := strconv.ParseFloat(coords[0], 64)
//		lng, _ := strconv.ParseFloat(coords[1], 64)
//		res = append(res, Coords{lat, lng})
//	}
//
//	b, err := json.Marshal(res)
//	if err != nil {
//		panic(err)
//	}
//
//	fmt.Println(string(b))
//}
//
//func (Ω *protectedAreas) PrintGeohashPrecisionGroups() {
//
//	colRef, err := Ω.floraStore.FirestoreCollection(store.CollectionProtectedAreas)
//	if err != nil {
//		panic(err)
//	}
//
//	iter := colRef.Documents(Ω.cxt)
//
//	c, err := geohashindex.NewCollection("itdw574")
//	if err != nil {
//		panic(err)
//	}
//
//	//tokens := map[uint]map[uint64]int{}
//
//	//start := uint(1)
//	//end := uint(6)
//	//for i := start; i <= end; i++ {
//	//	tokens[i] = map[uint64]int{}
//	//}
//
//	for {
//		snap, err := iter.Next()
//		if err == iterator.Done {
//			break
//		}
//		if err != nil {
//			panic(err)
//		}
//		var loc struct {
//			Name          string `json:""`
//			GeoFeatureSet struct {
//				GeoPoint latlng.LatLng `json:""`
//			} `json:""`
//		}
//		if err := snap.DataTo(&loc); err != nil {
//			panic(err)
//		}
//
//		if geohash.EncodeWithPrecision(
//			loc.GeoFeatureSet.GeoPoint.GetLatitude(),
//			loc.GeoFeatureSet.GeoPoint.GetLongitude(),
//			3,
//		) == "9yp" {
//			fmt.Println(geohash.Encode(
//				loc.GeoFeatureSet.GeoPoint.GetLatitude(),
//				loc.GeoFeatureSet.GeoPoint.GetLongitude(),
//			))
//			fmt.Println(utils.JsonOrSpew(loc))
//		}
//
//		continue
//
//		if err := c.AddPoint(
//			loc.GeoFeatureSet.GeoPoint.GetLatitude(),
//			loc.GeoFeatureSet.GeoPoint.GetLongitude(),
//			"20180102",
//			loc.Name,
//		); err != nil {
//			panic(err)
//		}
//
//		if err := c.AddPoint(
//			loc.GeoFeatureSet.GeoPoint.GetLatitude(),
//			loc.GeoFeatureSet.GeoPoint.GetLongitude(),
//			"20180108",
//			loc.Name,
//		); err != nil {
//			panic(err)
//		}
//	}
//
//	//predictionRef, err := Ω.floraStore.FirestoreCollection(store.CollectionGeoIndex)
//	//if err != nil {
//	//	panic(err)
//	//}
//	//
//	//if err := c.Upload(context.Background(), predictionRef); err != nil {
//	//	panic(err)
//	//}
//
//	//pColRef, err := Ω.floraStore.FirestoreCollection(store.CollectionPredictions)
//	//if err != nil {
//	//	panic(err)
//	//}
//	//
//	////unwound := []nested{}
//	//for _, v := range res {
//	//if _, _, err := pColRef.Add(context.Background(), v); err != nil {
//	//	panic(err)
//	//}
//	//}
//	//for _, v := range res {
//	//	//fmt.Println(utils.JsonOrSpew(v))
//	//	//continue
//	//	b, err := json.Marshal(v)
//	//	if err != nil {
//	//		panic(err)
//	//	}
//	//	res := map[string]interface{}{}
//	//	if err := json.Unmarshal(b, &res); err != nil {
//	//		panic(err)
//	//	}
//	//	if _, _, err := pColRef.Add(context.Background(), res); err != nil {
//	//		panic(err)
//	//	}
//	//}
//	//fmt.Println("Count", len(res))
//
//	//for i := range tokens {
//	//	fmt.Println(fmt.Sprintf("Precision %d", i))
//	//	fmt.Println("Groups:", len(tokens[i]))
//	//	for k, v := range tokens[i] {
//	//		fmt.Println(fmt.Sprintf("%d: %d", k, v))
//	//	}
//	//}
//}
//
//func (Ω *protectedAreas) GroupByS2Tokens() []map[string]int {
//	colRef, err := Ω.floraStore.FirestoreCollection(store.CollectionProtectedAreas)
//	if err != nil {
//		panic(err)
//	}
//
//	snaps, err := colRef.Documents(Ω.cxt).GetAll()
//	if err != nil {
//		panic(err)
//	}
//
//	//total := len(snaps)
//	tokens := []map[string]int{}
//	for i := 0; i < 10; i++ {
//		tokens = append(tokens, map[string]int{})
//	}
//	for _, snap := range snaps {
//		data, err := snap.DataAt("GeoFeatureSet.S2Tokens")
//		if err != nil {
//			panic(err)
//		}
//		for k, v := range data.(map[string]interface{}) {
//			i, _ := strconv.Atoi(k)
//			token := v.(string)
//			if _, ok := tokens[i][token]; ok {
//				tokens[i][token] = tokens[i][token] + 1
//			} else {
//				tokens[i][token] = 1
//			}
//		}
//	}
//	return tokens
//}
