package taxa

import (
	"github.com/heindl/floracast_process/nameusage/nameusage"
	"github.com/heindl/floracast_process/store"
	"github.com/heindl/floracast_process/utils"
	"context"
	"github.com/dropbox/godropbox/errors"
	"strings"
)

type MaterializedTaxon struct {
	ScientificName string       `json:",omitempty" firestore:",omitempty"`
	CommonName     string       `json:",omitempty" firestore:",omitempty"`
	Photo          *photo       `json:",omitempty" firestore:",omitempty"`
	Description    *description `json:",omitempty" firestore:",omitempty"`
}

func Fetch(ctx context.Context, floraStore store.FloraStore, id nameusage.ID) (*MaterializedTaxon, error) {
	col, err := floraStore.FirestoreCollection(store.CollectionTaxa)
	if err != nil {
		return nil, err
	}
	snap, err := col.Doc(id.String()).Get(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not get MaterializedTaxon from FireStore [%s]", id)
	}
	txn := MaterializedTaxon{}
	if err := snap.DataTo(&txn); err != nil {
		return nil, errors.Wrapf(err, "Could not cast MaterializedTaxon from FireStore [%s]", id)
	}
	return &txn, nil
}

// UploadMaterializedTaxon takes a NameUsage, materializes it, clears old references, and saves to FireStore.
func UploadMaterializedTaxon(ctx context.Context, florastore store.FloraStore, usage nameusage.NameUsage, deletedUsageIDs ...nameusage.ID) error {
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

func materialize(ctx context.Context, usage nameusage.NameUsage) (*MaterializedTaxon, error) {

	description, err := fetchDescription(ctx, usage)
	if err != nil {
		return nil, err
	}

	commonName, err := usage.CommonName()
	if err != nil {
		return nil, err
	}

	// Remove 'Mushroom' from Name
	commonNameSlice := []string{}
	for _, f := range strings.Fields(commonName) {
		if !strings.Contains(strings.ToLower(f), "mushroom") {
			commonNameSlice = append(commonNameSlice, f)
		}
	}

	commonName, err = utils.FormatTitle(strings.Join(commonNameSlice, " "))
	if err != nil {
		return nil, err
	}

	photo, err := fetchPhoto(ctx, usage)
	if err != nil {
		return nil, err
	}

	mt := MaterializedTaxon{
		ScientificName: utils.CapitalizeString(usage.CanonicalName().ScientificName()),
		CommonName:     commonName,
		Photo:          photo,
		Description:    description,
	}

	return &mt, nil
}
