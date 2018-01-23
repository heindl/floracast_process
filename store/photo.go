package store

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/saleswise/errors/errors"
)

const EntityKindPhoto = "Photo"

type Photo struct {
	ID            string             `datastore:",omitempty"`
	DataSourceID  DataSourceID       `datastore:",omitempty"`
	TaxonID       INaturalistTaxonID `datastore:",omitempty"`
	PhotoType     PhotoType          `datastore:",omitempty,noindex" json:"type,omitempty" bson:"type,omitempty"`
	URL           string             `datastore:",omitempty,noindex" json:"url,omitempty" bson:"url,omitempty"`
	NativePhotoID string             `datastore:",omitempty,noindex" json:"nativePhotoId,omitempty" bson:"nativePhotoId,omitempty"`
	SquareURL     string             `datastore:",omitempty,noindex" json:"squareUrl,omitempty" bson:"squareUrl,omitempty"`
	SmallURL      string             `datastore:",omitempty,noindex" json:"smallUrl,omitempty" bson:"smallUrl,omitempty"`
	MediumURL     string        `datastore:",omitempty,noindex" json:"mediumUrl,omitempty" bson:"mediumUrl,omitempty"`
	LargeURL      string        `datastore:",omitempty,noindex" json:"largeUrl,omitempty" bson:"largeUrl,omitempty"`
	Attribution   string        `datastore:",omitempty,noindex" json:"attribution,omitempty" bson:"attribution,omitempty"`
	LicenseCode   string        `datastore:",omitempty,noindex" json:"licenseCode,omitempty" bson:"licenseCode,omitempty"`
	Flags         []interface{} `datastore:",omitempty,noindex" json:"flags,omitempty" bson:"flags,omitempty"`
}

type Photos []*Photo

//func (Ω Photos) AddToSet(p *Photo) (Photos, error) {
//	if err := SchemeKeyName(p.Key.Parent.Name).Validate(); err != nil {
//		return nil, err
//	}
//	if p.Key.Parent.Kind != EntityKindMetaScheme {
//		return nil, errors.New("photo does not have a meta scheme parent").SetState(M{utils.LogkeyPhoto: p, utils.LogkeyDatastoreKey: p.Key})
//	}
//	if p.Key.Incomplete() {
//		return nil, errors.New("a photo id is required")
//	}
//	if p.Key.Kind != EntityKindPhoto {
//		return nil, errors.New("invalid photo entity kind").SetState(M{utils.LogkeyPhoto: p, utils.LogkeyDatastoreKey: p.Key})
//	}
//
//	// TODO: Additional validation
//	for i := range Ω {
//		if Ω[i].Key.ID != p.Key.ID || Ω[i].Key.Parent.Name != p.Key.Parent.Name {
//			continue
//		}
//		Ω[i] = Ω[i].PushSynonym(p)
//		return Ω, nil
//	}
//	return append(Ω, p), nil
//}

//func (Ω Photo) PushSynonym(p *Photo) *Photo {
//	if p.URL != "" {
//		Ω.URL = p.URL
//	}
//	if p.NativePhotoID != "" {
//		Ω.NativePhotoID = p.NativePhotoID
//	}
//	if p.SquareURL != "" {
//		Ω.SquareURL = p.SquareURL
//	}
//	if p.SmallURL != "" {
//		Ω.SmallURL = p.SmallURL
//	}
//	if p.MediumURL != "" {
//		Ω.MediumURL = p.MediumURL
//	}
//	if p.LargeURL != "" {
//		Ω.LargeURL = p.LargeURL
//	}
//	if p.Attribution != "" {
//		Ω.Attribution = p.Attribution
//	}
//	if p.LicenseCode != "" {
//		Ω.LicenseCode = p.LicenseCode
//	}
//	if len(p.Flags) > 0 {
//		Ω.Flags = p.Flags
//	}
//	return &Ω
//}

type PhotoType string

const (
	PhotoTypeFlickr      = PhotoType("FlickrPhoto")
	PhotoTypeINaturalist = PhotoType("INaturalist")
)

func (Ω PhotoType) Valid() bool {
	return Ω == PhotoTypeFlickr || Ω == PhotoTypeINaturalist
}

func (Ω *store) NewPhotoDocumentRef(taxonID INaturalistTaxonID, dataSourceID DataSourceID, photoID string) (*firestore.DocumentRef, error) {

	if !taxonID.Valid() {
		return nil, errors.New("invalid photo taxon id")
	}
	if !dataSourceID.Valid() {
		return nil, errors.New("invalid data source id")
	}
	if photoID == "" {
		return nil, errors.New("invalid photo id")
	}

	return Ω.FirestoreClient.Collection(CollectionTypePhotos).
		Doc(fmt.Sprintf("%s|%s|%s", string(taxonID), dataSourceID, photoID)), nil

}

func (Ω *store) SetPhoto(cxt context.Context, p Photo) error {

	ref, err := Ω.NewPhotoDocumentRef(p.TaxonID, p.DataSourceID, p.ID)
	if err != nil {
		return err
	}

	if _, err := ref.Set(cxt, p); err != nil {
		return errors.Wrap(err, "could not set photo")
	}

	return nil
}
