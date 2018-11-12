package predictions

import (
	"github.com/heindl/floracast_process/nameusage/nameusage"
	"github.com/heindl/floracast_process/store"
	"github.com/heindl/floracast_process/terra/geoembed"
	"github.com/heindl/floracast_process/utils"
	"context"
	"github.com/dropbox/godropbox/errors"
	"google.golang.org/api/iterator"
	"strings"
)

func FetchFromFireStore(cxt context.Context, floraStore store.FloraStore, id nameusage.ID) (Predictions, error) {
	col, err := floraStore.FirestoreCollection(store.CollectionPredictionIndex)
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

		r := record{}
		if err := snap.DataTo(&r); err != nil {
			return nil, errors.Wrap(err, "Could not cast Firestore data to Prediction")
		}

		lat, lng, err := geoembed.S2Key(strings.Split(snap.Ref.ID, "-")[1]).Parse()
		if err != nil {
			return nil, err
		}

		for date, m := range r.Timeline {
			p, err := NewPrediction(r.NameUsageID, utils.FormattedDate(date), lat, lng, m.Value)
			if err != nil {
				return nil, err
			}
			res = append(res, p)
		}
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
