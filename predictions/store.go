package predictions

import (
	"bitbucket.org/heindl/process/nameusage/nameusage"
	"bitbucket.org/heindl/process/store"
	"context"
	"github.com/dropbox/godropbox/errors"
	"google.golang.org/api/iterator"
)

func FetchFromFireStore(cxt context.Context, floraStore store.FloraStore, id nameusage.ID) (Predictions, error) {
	col, err := floraStore.FirestoreCollection(store.CollectionOccurrences)
	if err != nil {
		return nil, err
	}
	res := Predictions{}
	iter := col.Where("NameUsageID", "==", id).Documents(cxt)
	for {
		snap, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, errors.Wrap(err, "Could not get Occurrence from Firestore")
		}
		p := prediction{}
		if err := snap.DataTo(&p); err != nil {
			return nil, errors.Wrap(err, "Could not cast Firestore data to Prediction")
		}
		res = append(res, &p)
	}
	return res, nil
}

//func (Ω *store) PredictionDocumentID(p Prediction) (string, error) {
//	if !p.TaxonID.Valid() {
//		return "", errors.New("invalid taxon id")
//	}
//	if p.Date == nil || p.Date.IsZero() {
//		return "", errors.New("invalid date")
//	}
//	if p.Location.GetLatitude() == 0 {
//		return "", errors.New("invalid latitude")
//	}
//	if p.Location.GetLongitude() == 0 {
//		return "", errors.New("invalid longitude")
//	}
//	return fmt.Sprintf("%s|%s|%.6f|%.6f", string(p.TaxonID), p.Date.Format("20060102"), p.Location.GetLatitude(), p.Location.GetLongitude()), nil
//}
