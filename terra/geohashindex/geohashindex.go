package geohashindex

import (
	"bitbucket.org/heindl/process/nameusage/nameusage"
	"bitbucket.org/heindl/process/store"
	"bitbucket.org/heindl/process/terra/geo"
	"bitbucket.org/heindl/process/utils"
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/mmcloughlin/geohash"
	"github.com/saleswise/errors/errors"
	"gopkg.in/mgo.v2/bson"
)

type Collection interface {
	AddPoint(lat, lng float64, date string, value interface{}) error
	Print()
	Upload(context.Context, *firestore.CollectionRef) error
}

func NewCollection(id nameusage.ID) (Collection, error) {
	if !id.Valid() {
		return nil, errors.New("Invalid NameUsageID")
	}
	return &collection{nameUsageID: id, indexes: map[hash3]dateMap{}}, nil
}

type hash3 string

type dateMap map[string]*LetterField

type collection struct {
	nameUsageID nameusage.ID
	indexes     map[hash3]dateMap
}

func (Ω *collection) Print() {
	for k, v := range Ω.indexes {
		b, err := bson.Marshal(v)
		if err != nil {
			panic(err)
		}
		res := map[string]interface{}{}
		if err := bson.Unmarshal(b, &res); err != nil {
			panic(err)
		}
		fmt.Println("Key:", k)
		fmt.Println("NameUsageID:", Ω.nameUsageID)
		fmt.Println(utils.JsonOrSpew(res))
	}
}

func (Ω *collection) Upload(cxt context.Context, colRef *firestore.CollectionRef) error {

	for k, v := range Ω.indexes {

		docRef := colRef.Doc(string(k)).Collection(string(store.CollectionPredictions)).Doc(string(Ω.nameUsageID))
		//
		//b, err := bson.Marshal(v)
		//if err != nil {
		//	panic(err)
		//}
		//res := map[string]interface{}{}
		//if err := bson.Unmarshal(b, &res); err != nil {
		//	panic(err)
		//}

		if _, err := docRef.Set(cxt, v); err != nil {
			return errors.Newf("Could not create GeoHashIndex [%s]", docRef.Path)
		}
	}

	return nil
}

func (Ω *collection) AddPoint(lat, lng float64, date string, value interface{}) error {
	if err := geo.ValidateCoordinates(lat, lng); err != nil {
		return err
	}
	if len(date) != 8 {
		return errors.Newf("Invalid date [%s] for GeoHashIndex Point", date)
	}
	if value == nil {
		return errors.New("Value required for GeoHashIndex Point")
	}

	key := hash3(geohash.EncodeWithPrecision(lat, lng, 3))
	if len(key) != 3 {
		return errors.New("Could not encode GeoHash")
	}

	if _, ok := Ω.indexes[key]; !ok {
		Ω.indexes[key] = dateMap{}
	}
	if _, ok := Ω.indexes[key][date]; !ok {
		Ω.indexes[key][date] = &LetterField{}
	}
	Ω.indexes[key][date].Add(geohash.Encode(lat, lng), 3, value)

	return nil

}
