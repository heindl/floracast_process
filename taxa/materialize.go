package taxa

import (
	"bitbucket.org/heindl/process/nameusage/nameusage"
	"bitbucket.org/heindl/process/store"
	"bitbucket.org/heindl/process/utils"
	"context"
	"github.com/dropbox/godropbox/errors"
	"strings"
)

type materializedTaxon struct {
	ScientificName string       `json:",omitempty" firestore:",omitempty"`
	CommonName     string       `json:",omitempty" firestore:",omitempty"`
	Photo          *photo       `json:",omitempty" firestore:",omitempty"`
	Description    *description `json:",omitempty" firestore:",omitempty"`
}

func UploadMaterializedTaxa(ctx context.Context, florastore store.FloraStore, usage nameusage.NameUsage, deletedUsageIDs ...nameusage.ID) error {
	if err := clearMaterializedTaxa(ctx, florastore, deletedUsageIDs); err != nil {
		return err
	}

	id, err := usage.ID()
	if err != nil {
		return nil
	}

	col, err := florastore.FirestoreCollection(store.CollectionTaxa)
	if err != nil {
		return err
	}

	docRef := col.Doc(id.String())
	materialized, err := materialize(ctx, usage)
	if err != nil {
		return err
	}
	if _, err := docRef.Set(ctx, materialized); err != nil {
		return err
	}
	return nil
}

func clearMaterializedTaxa(ctx context.Context, florastore store.FloraStore, allUsageIDs nameusage.IDs) error {

	for _, usageIDs := range allUsageIDs.Batch(500) {
		if len(usageIDs) == 0 {
			return nil
		}
		batch := florastore.FirestoreBatch()
		for _, id := range usageIDs {
			col, err := florastore.FirestoreCollection(store.CollectionTaxa)
			if err != nil {
				return err
			}
			batch = batch.Delete(col.Doc(id.String()))
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
		CommonName:     strings.Title(commonName),
		Photo:          photo,
		Description:    description,
	}

	return &mt, nil
}
