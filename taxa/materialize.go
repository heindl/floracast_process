package taxa

import (
	"bitbucket.org/heindl/processors/nameusage/nameusage"
	"context"
	"time"
	"bitbucket.org/heindl/processors/store"
	"github.com/dropbox/godropbox/errors"
)

type materializedTaxon struct {
	ScientificName string `json:""`
	CommonName string `json:""`
	Photos []photo `json:""`
	Description description `json:""`
	CreatedAt time.Time `json:""`
	ModifiedAt time.Time `json:""`
}

func UploadMaterializedTaxa(ctx context.Context, florastore store.FloraStore, usage nameusage.NameUsage, deletedUsageIDs nameusage.NameUsageIDs) error {
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

func materialize(ctx context.Context, usage nameusage.NameUsage) (map[string]interface{}, error) {

	name, err := usage.CommonName()
	if err != nil {
		return nil, err
	}

	photos, err := photos(ctx, usage)
	if err != nil {
		return nil, err
	}

	descriptions, err := descriptions(ctx, usage)
	if err != nil {
		return nil, err
	}

	m := map[string]interface{}{
		"ScientificName": usage.CanonicalName().ScientificName(),
		"CommonName": name,
		"Photos": photos,
		"Descriptions": descriptions,
	}

	return m, nil
}