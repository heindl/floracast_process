package store

import (
	"github.com/saleswise/errors/errors"
	"cloud.google.com/go/datastore"
	"context"
	"bitbucket.org/heindl/utils"
	. "github.com/saleswise/malias"
)

const EntityKindPhoto = "Photo"

type Photo struct {
	Key           *datastore.Key `datastore:"__key__"`
	Type          PhotoType `datastore:",omitempty,noindex" json:"type,omitempty" bson:"type,omitempty"`
	URL           string `datastore:",omitempty,noindex" json:"url,omitempty" bson:"url,omitempty"`
	NativePhotoID string `datastore:",omitempty,noindex" json:"nativePhotoId,omitempty" bson:"nativePhotoId,omitempty"`
	SquareURL     string `datastore:",omitempty,noindex" json:"squareUrl,omitempty" bson:"squareUrl,omitempty"`
	SmallURL      string `datastore:",omitempty,noindex" json:"smallUrl,omitempty" bson:"smallUrl,omitempty"`
	MediumURL     string `datastore:",omitempty,noindex" json:"mediumUrl,omitempty" bson:"mediumUrl,omitempty"`
	LargeURL      string `datastore:",omitempty,noindex" json:"largeUrl,omitempty" bson:"largeUrl,omitempty"`
	Attribution   string `datastore:",omitempty,noindex" json:"attribution,omitempty" bson:"attribution,omitempty"`
	LicenseCode   string `datastore:",omitempty,noindex" json:"licenseCode,omitempty" bson:"licenseCode,omitempty"`
	Flags         []interface{} `datastore:",omitempty,noindex" json:"flags,omitempty" bson:"flags,omitempty"`
}

type Photos []*Photo

func (Ω Photos) AddToSet(p *Photo) (Photos, error) {
	if err := SchemeKeyName(p.Key.Parent.Name).Validate(); err != nil {
		return nil, err
	}
	if p.Key.Parent.Kind != EntityKindMetaScheme {
		return nil, errors.New("photo does not have a meta scheme parent").SetState(M{utils.LogkeyPhoto: p, utils.LogkeyDatastoreKey: p.Key})
	}
	if p.Key.Incomplete() {
		return nil, errors.New("a photo id is required")
	}
	if p.Key.Kind != EntityKindPhoto {
		return nil, errors.New("invalid photo entity kind").SetState(M{utils.LogkeyPhoto: p, utils.LogkeyDatastoreKey: p.Key})
	}

	// TODO: Additional validation
	for i := range Ω {
		if Ω[i].Key.ID != p.Key.ID || Ω[i].Key.Parent.Name != p.Key.Parent.Name {
			continue
		}
		Ω[i] = Ω[i].Combine(p)
		return Ω, nil
	}
	return append(Ω, p), nil
}

func (Ω Photo) Combine(p *Photo) *Photo {
	if p.URL != "" {
		Ω.URL = p.URL
	}
	if p.NativePhotoID != "" {
		Ω.NativePhotoID = p.NativePhotoID
	}
	if p.SquareURL != "" {
		Ω.SquareURL = p.SquareURL
	}
	if p.SmallURL != "" {
		Ω.SmallURL = p.SmallURL
	}
	if p.MediumURL != "" {
		Ω.MediumURL = p.MediumURL
	}
	if p.LargeURL != "" {
		Ω.LargeURL = p.LargeURL
	}
	if p.Attribution != "" {
		Ω.Attribution = p.Attribution
	}
	if p.LicenseCode != "" {
		Ω.LicenseCode = p.LicenseCode
	}
	if len(p.Flags) > 0 {
		Ω.Flags = p.Flags
	}
	return &Ω
}

type PhotoType string
const (
	PhotoTypeFlickr = PhotoType("FlickrPhoto")
	PhotoTypeINaturalist = PhotoType("INaturalist")
)

func (Ω *store) SetPhotos(photos Photos) error {
	keys := []*datastore.Key{}
	for _, s := range photos {
		keys = append(keys, s.Key)
	}

	if _, err := Ω.DatastoreClient.PutMulti(context.Background(), keys, photos); err != nil {
		return errors.Wrap(err, "could not batch scheme puts to datastore")
	}
	return nil
}