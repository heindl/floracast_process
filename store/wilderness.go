package store

import (
	"google.golang.org/genproto/googleapis/type/latlng"
	"github.com/fatih/structs"
	"github.com/saleswise/errors/errors"
	"cloud.google.com/go/firestore"
	"context"
	"strings"
)

type WildernessArea struct {
	ID string `datastore:",omitempty"`
	State       WildernessAreaState                `datastore:",omitempty"`
	Acres       float64               `datastore:",omitempty"`
	Name        string                `datastore:",omitempty"`
	Centre latlng.LatLng `datastore:",omitempty"`
	RadiusKilometers float64  `datastore:",omitempty"`
	ManagerType string `datastore:",omitempty"`
	ManagerName string `datastore:",omitempty"`
	ManagementDesignation string `datastore:",omitempty"`
	OwnerType string `datastore:",omitempty"`
	OwnerName string `datastore:",omitempty"`
	Category string `datastore:",omitempty"`
	YearEstablished int `datastore:",omitempty"`
	PublicAccess string `datastore:",omitempty"`

}

type WildernessAreaState string

type WildernessAreas []WildernessArea

func (a WildernessAreas) Len() int           { return len(a) }
func (a WildernessAreas) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a WildernessAreas) Less(i, j int) bool { return a[i].Acres > a[j].Acres }

func (Ω *store) ReadWildernessArea(cxt context.Context, lat, lng float64) (*WildernessArea, error) {

	// Validate
	if lat == 0 || lng == 0 {
		return nil, errors.New("invalid wilderness area id")
	}

	docs, err := Ω.FirestoreClient.Collection(CollectionTypeWildernessAreas).
		Where("Centre.Latitude", "==", lat).
		Where("Centre.Longitude", "==", lng).
		Documents(cxt).
		GetAll()

	if err != nil {
		return nil, err
	}
	if len(docs) > 0 {
		return nil, errors.New("more than one wilderness area found")
	}

	if len(docs) == 0 {
		return nil, errors.New("no wilderness area found")
	}

	w := WildernessArea{}
	if err := docs[0].DataTo(&w); err != nil {
		return nil, errors.Wrap(err, "could not type cast WildernessArea")
	}

	return &w, nil
}


func (Ω *store) UpsertWildernessArea(cxt context.Context, wa WildernessArea) error {

	// Validate
	if wa.ID == "" {
		return errors.New("invalid wilderness area id")
	}

	ref := Ω.FirestoreClient.Collection(CollectionTypeWildernessAreas).Doc(wa.ID)

	if err := Ω.FirestoreClient.RunTransaction(cxt, func(cxt context.Context, tx *firestore.Transaction) error {
		if _, err := tx.Get(ref); err != nil {
			if strings.Contains(err.Error(), "not found") {
				return tx.Set(ref, wa)
			} else {
				return err
			}
		}
		return tx.UpdateMap(ref, structs.Map(wa))
	}); err != nil {
		return errors.Wrap(err, "could not update occurrence")
	}

	return nil
}
