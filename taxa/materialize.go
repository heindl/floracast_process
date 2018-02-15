package taxa

import (
	"bitbucket.org/heindl/processors/nameusage/nameusage"
	"context"
	"time"
	"bitbucket.org/heindl/processors/store"
	"github.com/dropbox/godropbox/errors"
	"bitbucket.org/heindl/processors/utils"
	"strings"
)

type materializedTaxon struct {
	ScientificName string    `json:",omitempty" firestore:",omitempty"`
	CommonName string        `json:",omitempty" firestore:",omitempty"`
	Photo *photo             `json:",omitempty" firestore:",omitempty"`
	Description *description `json:",omitempty" firestore:",omitempty"`
	CreatedAt time.Time      `json:"" firestore:""`
	ModifiedAt time.Time     `json:"" firestore:""`
}

func UploadMaterializedTaxa(ctx context.Context, florastore store.FloraStore, usage nameusage.NameUsage, deletedUsageIDs ...nameusage.NameUsageID) error {
	if err := clearMaterializedTaxa(ctx, florastore, deletedUsageIDs); err != nil {
		return err
	}

	id, err := usage.ID()
	if err != nil {
		return nil
	}

	docRef := florastore.FirestoreCollection(store.CollectionTaxa).Doc(id.String())
	materialized, err := materialize(ctx, usage)
	if err != nil {
		return err
	}
	if _, err := docRef.Set(ctx, materialized); err != nil {
		return err
	}
	return nil
}

func clearMaterializedTaxa(ctx context.Context, florastore store.FloraStore, allUsageIDs nameusage.NameUsageIDs) error {

	for _, usageIDs := range allUsageIDs.Batch(500) {
		if len(usageIDs) == 0 {
			return nil
		}
		batch := florastore.FirestoreBatch()
		for _, id := range usageIDs {
			docRef := florastore.FirestoreCollection(store.CollectionTaxa).Doc(id.String())
			batch = batch.Delete(docRef)
		}
		if _, err := batch.Commit(ctx); err != nil {
			return errors.Wrap(err, "Could not commit materialized taxa")
		}
	}
	return nil
}

func materialize(ctx context.Context, usage nameusage.NameUsage) (*materializedTaxon, error) {

	description, err := fetchDescription(ctx, usage)
	if err != nil {
		return nil, err
	}

	commonName, err := usage.CommonName()
	if err != nil {
		return nil, err
	}

	photo, err := fetchPhoto(ctx, usage)
	if err != nil {
		return nil, err
	}

	mt := materializedTaxon{
		ScientificName: utils.CapitalizeString(usage.CanonicalName().ScientificName()),
		CommonName: strings.Title(commonName),
		Photo: photo,
		Description: description,
		CreatedAt: time.Now(),
		ModifiedAt: time.Now(),
	}

	return &mt, nil
}