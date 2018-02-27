package occurrences

import (
	"context"
	"fmt"
	"strings"

	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/geofeatures"
	"bitbucket.org/heindl/process/store"
	"cloud.google.com/go/firestore"
	"github.com/dropbox/godropbox/errors"
	"github.com/mongodb/mongo-tools/common/json"
	"google.golang.org/api/iterator"
)

// TODO: Should periodically check all occurrences for consistency.

func (Ω *OccurrenceAggregation) Upload(cxt context.Context, florastore store.FloraStore) error {
	for _, _o := range Ω.list {
		o := _o
		transactionFunc, err := o.UpsertTransactionFunc(florastore)
		if err != nil {
			return err
		}
		if err := florastore.FirestoreTransaction(cxt, transactionFunc); err != nil {
			return err
		}
	}
	return nil
}

func ClearRandomOccurrences(cxt context.Context, florastore store.FloraStore) error {

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
			return errors.Wrap(err, "Could not delete random occurrence")
		}
	}

	return nil
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

func (Ω *occurrence) docRef(florastore store.FloraStore) (*firestore.DocumentRef, error) {
	id, err := Ω.ID()
	if err != nil {
		return nil, err
	}
	col, err := Ω.Collection(florastore)
	if err != nil {
		return nil, err
	}
	return col.Doc(id), nil
}

func (Ω *occurrence) UpsertTransactionFunc(florastore store.FloraStore) (store.FirestoreTransactionFunc, error) {

	if !Ω.SrcType.Valid() || !Ω.TgtID.Valid(Ω.SrcType) || strings.TrimSpace(Ω.SrcOccurrenceID) == "" {
		return nil, errors.Newf("Invalid firestore reference ID: %s, %s, %s", Ω.SourceType(), Ω.TargetID(), Ω.SourceOccurrenceID())
	}

	// Design Note: Anything that can be checked and failed early, should be handled before the transaction.

	newDocRef, err := Ω.docRef(florastore)
	if err != nil {
		return nil, err
	}

	col, err := Ω.Collection(florastore)
	if err != nil {
		return nil, err
	}

	q, err := geofeatures.CoordinateQuery(col, Ω.GeoFeatureSet.Lat(), Ω.GeoFeatureSet.Lng())
	if err != nil {
		return nil, err
	}
	locationQuery := q.Where("FormattedDate", "==", Ω.FormattedDate)

	b, err := json.Marshal(Ω)
	if err != nil {
		return nil, errors.Wrap(err, "Could not marshal occurrence")
	}

	newOccurrenceMapDoc := map[string]interface{}{}
	if err := json.Unmarshal(b, &newOccurrenceMapDoc); err != nil {
		return nil, errors.Wrap(err, "Could not unmarshal occurrence doc into map")
	}

	return func(cxt context.Context, tx *firestore.Transaction) error {

		_, err := tx.Get(newDocRef)
		notFound := (err != nil && strings.Contains(err.Error(), "not found"))
		if !notFound && err != nil {
			return errors.Wrapf(err, "Could not get firestore occurrence doc [%s]", newDocRef.ID)
		}

		idAlreadyExists := !notFound

		imbricates, err := tx.Documents(locationQuery).GetAll()
		if err != nil {
			return errors.Wrap(err, "Error searching for a list of possibly overlapping occurrences")
		}

		if len(imbricates) > 1 {
			return errors.Newf("Unexpected: multiple imbricates found for occurrence with location [%f, %f, %s]", Ω.GeoFeatureSet.Lat(), Ω.GeoFeatureSet.Lng(), Ω.FormattedDate)
		}

		isImbricative := len(imbricates) > 0

		if idAlreadyExists && isImbricative {
			// This suggests the location has changed somewhere. Update code if we see this.
			fmt.Println(fmt.Sprintf("Unexpected: occurrence with id [%s] idAlreadyExists and is imbricative to another doc [%s]", newDocRef.ID, imbricates[0].Ref.ID))
		}

		if isImbricative && !idAlreadyExists {

			// TODO: Be wary of cases in which there are occurrence of two different species in the same spot. Not sure if this will come up.

			originalID, err := Ω.ID()
			if err != nil {
				return err
			}

			fmt.Println(fmt.Sprintf("Warning: Imbricative Occurrence Locations [%s, %s]", originalID, imbricates[0].Ref.ID))

			m := map[string]interface{}{}
			if err := imbricates[0].DataTo(&m); err != nil {
				return errors.Wrap(err, "Could not cast occurrence")
			}

			b, err := json.Marshal(m)
			if err != nil {
				return errors.Wrap(err, "Could not marshal firebase occurrence response")
			}

			imbricate := occurrence{}
			if err := json.Unmarshal(b, &imbricate); err != nil {
				return errors.Wrap(err, "Could not unmarshal occurrence")
			}

			if Ω.SourceType() != imbricate.SourceType() && imbricate.SourceType() == datasources.TypeGBIF {
				// So we have something other than GBIF, and the GBIF record is already in the database.
				// No opt to prefer the existing GBIF record.
				return nil
			}

			// Condition 1: The two are the same source, but one of the locations has changed, so delete the old to be safe.
			fmt.Println("Warning: Source type for imbricating locations are the same. Deleting the old one.")
			// Condition 2: So this is a GBIF source, and that is not, which means need to delete the old one.

			if err := tx.Delete(imbricates[0].Ref); err != nil {
				return errors.Wrapf(err, "Unable to delete occurrence [%s]", imbricates[0].Ref.ID)
			}
		}

		// Should be safe to override with new record
		if err := tx.Set(newDocRef, newOccurrenceMapDoc); err != nil {
			return errors.Wrap(err, "Could not set occurrence")
		}

		return nil

	}, nil
}
