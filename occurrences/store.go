package occurrences

import (
	"context"
	"fmt"
	"strings"

	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/store"
	"bitbucket.org/heindl/process/terra/geoembed"
	"cloud.google.com/go/firestore"
	"encoding/json"
	"github.com/dropbox/godropbox/errors"
	"google.golang.org/api/iterator"
)

// TODO: Should periodically check all occurrences for consistency.

// Upload saves Occurrences to either the Occurrence or Random Firestore Collections.
func (Ω *OccurrenceAggregation) Upload(cxt context.Context, floraStore store.FloraStore) error {
	for _, _o := range Ω.list {
		o := _o
		transactionFunc, err := o.UpsertTransactionFunc(floraStore)
		if err != nil {
			return err
		}
		if err := floraStore.FirestoreTransaction(cxt, transactionFunc); err != nil {
			return err
		}
	}
	return nil
}

// ClearRandomPoints clears all points from the Random collection.
func ClearRandomPoints(cxt context.Context, florastore store.FloraStore) error {

	col, err := florastore.FirestoreCollection(store.CollectionRandom)
	if err != nil {
		return err
	}
	docs := col.Where("SourceType", "==", datasources.TypeRandom).Documents(cxt)
	for {
		snap, err := docs.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return errors.Wrap(err, "Could not get random occurrence snapshot")
		}
		if _, err := snap.Ref.Delete(cxt); err != nil {
			return errors.Wrap(err, "Could not delete Random occurrence")
		}
	}

	return nil
}

// UpsertTransactionFunc returns a transaction function for adding an Occurrence to FireStore.
func (Ω *occurrence) UpsertTransactionFunc(florastore store.FloraStore) (store.FirestoreTransactionFunc, error) {
	if !Ω.SrcType.Valid() || !Ω.TgtID.Valid(Ω.SrcType) || strings.TrimSpace(Ω.SrcOccurrenceID) == "" {
		return nil, errors.Newf("Invalid FireStore reference ID: %s, %s, %s", Ω.SourceType(), Ω.TargetID(), Ω.SourceOccurrenceID())
	}
	// Design Note: Anything that can be checked and failed early, should be handled before the transaction.
	docRef, err := Ω.docRef(florastore)
	if err != nil {
		return nil, err
	}
	return Ω.returnTransaction(docRef), nil
}

func (Ω *occurrence) ID() (string, error) {
	if !Ω.SourceType().Valid() || !Ω.TargetID().Valid(Ω.SourceType()) {
		return "", errors.Newf("Invalid SourceType [%s] and TargetID [%s]", Ω.SourceType(), Ω.TargetID())
	}

	if Ω.SourceOccurrenceID() == "" || Ω.SourceOccurrenceID() == "0" {
		return "", errors.Newf("Invalid SourceOccurrenceID [%s]", Ω.SourceOccurrenceID())
	}
	return fmt.Sprintf("%s-%s-%s", Ω.SourceType(), Ω.TargetID(), Ω.SourceOccurrenceID()), nil
}

func (Ω *occurrence) toMap() (map[string]interface{}, error) {
	b, err := json.Marshal(Ω)
	if err != nil {
		return nil, errors.Wrap(err, "Could not marshal occurrence")
	}

	m := map[string]interface{}{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, errors.Wrap(err, "Could not unmarshal occurrence doc into map")
	}

	return m, nil
}

func (Ω *occurrence) docRef(floraStore store.FloraStore) (*firestore.DocumentRef, error) {
	id, err := Ω.ID()
	if err != nil {
		return nil, err
	}
	col, err := Ω.Collection(floraStore)
	if err != nil {
		return nil, err
	}
	return col.Doc(id), nil
}

func (Ω *occurrence) returnTransaction(docRef *firestore.DocumentRef) store.FirestoreTransactionFunc {
	return func(cxt context.Context, tx *firestore.Transaction) error {

		_, err := tx.Get(docRef)
		if err != nil && !strings.Contains(err.Error(), "not found") {
			return err
		}
		idAlreadyExists := err != nil

		imbricate, err := Ω.fetchImbricate(tx, docRef.Parent)
		if err != nil {
			return err
		}

		if idAlreadyExists && imbricate != nil {
			// This suggests the location has changed somewhere. Update code if we see this.
			return errors.Newf("Unexpected: occurrence with id [%s] idAlreadyExists and is imbricative to another doc [%s]", docRef.ID, imbricate.Ref.ID)
		}

		if imbricate != nil && !idAlreadyExists {
			// TODO: Be wary of cases in which there are occurrence of two different species
			// in the same spot. Not sure if this will come up.
			shouldOverride, err := Ω.handleExistingRecord(tx, imbricate)
			if err != nil || !shouldOverride {
				return err
			}
		}

		return Ω.setInFirestore(tx, docRef)
	}
}

func (Ω *occurrence) setInFirestore(tx *firestore.Transaction, docRef *firestore.DocumentRef) error {
	m, err := Ω.toMap()
	if err != nil {
		return err
	}
	// Should be safe to override with new record
	if err := tx.Set(docRef, m); err != nil {
		return errors.Wrap(err, "Could not set occurrence")
	}
	return nil
}

func (Ω *occurrence) fetchImbricate(tx *firestore.Transaction, collection *firestore.CollectionRef) (*firestore.DocumentSnapshot, error) {

	q, err := geoembed.CoordinateQuery(collection, Ω.GeoFeatureSet.Lat(), Ω.GeoFeatureSet.Lng())
	if err != nil {
		return nil, err
	}
	locationQuery := q.Where("FormattedDate", "==", Ω.FormattedDate)

	imbricates, err := tx.Documents(locationQuery).GetAll()
	if err != nil {
		return nil, errors.Wrap(err, "Error searching for a list of possibly overlapping occurrences")
	}

	if len(imbricates) > 1 {
		return nil, errors.Newf("Unexpected: multiple imbricates found for occurrence with location [%f, %f, %s]", Ω.GeoFeatureSet.Lat(), Ω.GeoFeatureSet.Lng(), Ω.FormattedDate)
	}

	if len(imbricates) == 1 {
		return imbricates[0], nil
	}

	return nil, nil
}

func newOccurrencefromFireStoreSnap(doc *firestore.DocumentSnapshot) (*occurrence, error) {

	m := doc.Data()

	b, err := json.Marshal(m)
	if err != nil {
		return nil, errors.Wrap(err, "Could not marshal firebase occurrence response")
	}

	o := occurrence{}
	if err := json.Unmarshal(b, &o); err != nil {
		return nil, errors.Wrap(err, "Could not unmarshal occurrence")
	}

	return &o, nil

}

func (Ω *occurrence) handleExistingRecord(tx *firestore.Transaction, imbDoc *firestore.DocumentSnapshot) (shouldOverride bool, err error) {
	originalID, err := Ω.ID()
	if err != nil {
		return false, err
	}

	fmt.Println(fmt.Sprintf("Warning: Imbricative Occurrence Locations [%s, %s]", originalID, imbDoc.Ref.ID))

	occurrence, err := newOccurrencefromFireStoreSnap(imbDoc)
	if err != nil {
		return false, err
	}

	if Ω.SourceType() != occurrence.SourceType() && occurrence.SourceType() == datasources.TypeGBIF {
		// So we have something other than GBIF, and the GBIF record is already in the database.
		// No opt to prefer the existing GBIF record.
		return false, nil
	}

	// Condition 1: The two are the same source, but one of the locations has changed, so delete the old to be safe.
	fmt.Println("Warning: Source type for imbricating locations are the same. Deleting the old one.")
	// Condition 2: So this is a GBIF source, and that is not, which means need to delete the old one.

	if err := tx.Delete(imbDoc.Ref); err != nil {
		return false, errors.Wrapf(err, "Unable to delete occurrence [%s]", imbDoc.Ref.ID)
	}

	return true, nil
}
