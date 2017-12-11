package store

import (
	"google.golang.org/genproto/googleapis/type/latlng"
	"github.com/saleswise/errors/errors"
	"context"
	"math"
	"bitbucket.org/heindl/taxa/utils"
	"time"
	"github.com/cenkalti/backoff"
	"strings"
	"fmt"
)

type ProtectedArea struct {
	ID                    string             `datastore:",omitempty"`
	State                 ProtectedAreaState `datastore:",omitempty"`
	Acres                 float64            `datastore:",omitempty"`
	Name                  string             `datastore:",omitempty"`
	Centre                latlng.LatLng      `datastore:",omitempty"`
	RadiusKilometers      float64            `datastore:",omitempty"`
	ManagerType           string             `datastore:",omitempty"`
	ManagerName           string             `datastore:",omitempty"`
	ManagementDesignation string             `datastore:",omitempty"`
	OwnerType             string             `datastore:",omitempty"`
	OwnerName             string             `datastore:",omitempty"`
	Category              string             `datastore:",omitempty"`
	YearEstablished       int                `datastore:",omitempty"`
	PublicAccess          string             `datastore:",omitempty"`

}

type ProtectedAreaState string

type ProtectedAreas []ProtectedArea

func (a ProtectedAreas) Len() int           { return len(a) }
func (a ProtectedAreas) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ProtectedAreas) Less(i, j int) bool { return a[i].Acres > a[j].Acres }

func (Ω *store) ReadProtectedArea(cxt context.Context, lat, lng float64) (*ProtectedArea, error) {

	// Validate
	if lat == 0 || lng == 0 {
		return nil, errors.New("invalid protected area id")
	}

	docs, err := Ω.FirestoreClient.Collection(CollectionTypeProtectedAreas).
		// TODO: Would be great to use a geo query here or at least an approximation.
		Where("Centre.Longitude", ">", math.Floor(lng)).
		Where("Centre.Longitude", "<=", math.Ceil(lng)).
		Documents(cxt).
		GetAll()

	if err != nil {
		return nil, errors.Wrap(err, "could not find wilderness area")
	}

	for _, d := range docs {
		w := ProtectedArea{}
		if err := d.DataTo(&w); err != nil {
			return nil, errors.Wrap(err, "could not type cast ProtectedArea")
		}
		if utils.CoordinatesEqual(lat, w.Centre.Latitude) && utils.CoordinatesEqual(lat, w.Centre.Latitude) {
			return &w, nil
		}
	}

	return nil, errors.New("no wilderness area found")
}

var counter = 0;

func (Ω *store) SetProtectedArea(cxt context.Context, wa ProtectedArea) error {

	counter = counter + 1;

	fmt.Println(counter, utils.JsonOrSpew(wa))

	// Validate
	if wa.ID == "" {
		return errors.New("invalid wilderness area id")
	}

	bkf := backoff.NewExponentialBackOff()
	bkf.InitialInterval = time.Second * 1
	ticker := backoff.NewTicker(bkf)
	for _ = range ticker.C {
		_, err := Ω.FirestoreClient.Collection(CollectionTypeProtectedAreas).Doc(wa.ID).Set(cxt, wa)
		if err != nil && strings.Contains(err.Error(), "Internal error encountered") {
			fmt.Println("Internal error encountered", err)
			continue
		}
		if err != nil {
			ticker.Stop()
			return errors.Wrap(err, "could nto set protected area")
		}
		ticker.Stop()
		break
	}

	return nil
}
