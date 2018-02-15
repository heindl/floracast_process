package occurrences

import (
	"bitbucket.org/heindl/processors/utils"
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"strings"
	"time"
	"bitbucket.org/heindl/processors/store"
	"bitbucket.org/heindl/processors/geofeatures"
	"go.uber.org/ratelimit"
	"gopkg.in/tomb.v2"
	"strconv"
	"sync"
	"bitbucket.org/heindl/processors/datasources"
	"encoding/json"
)

type Occurrence interface{
	SourceType() datasources.SourceType
	TargetID() datasources.TargetID
	LocationKey() (string, error)
	SourceOccurrenceID() string
	UpsertTransaction(florastore store.FloraStore) (store.FirestoreTransactionFunc, error)
	SetGeospatial(lat, lng float64, date string, coordinatesEstimated bool) error
}

func NewOccurrence(srcType datasources.SourceType, targetID datasources.TargetID, occurrenceID string) (Occurrence, error) {
	if !srcType.Valid() {
		return nil, errors.Newf("Invalid source type [%s]", srcType)
	}
	if !targetID.Valid(srcType) {
		return nil, errors.Newf("Invalid target id [%s]", targetID)
	}
	if occurrenceID == "" {
		return nil, errors.Newf("Invalid occurrence id")
	}

	return &occurrence{
		SrcType: srcType,
		TgtID: targetID,
		SrcOccurrenceID: occurrenceID,
		CreatedAt: utils.TimePtr(time.Now()),
		ModifiedAt: utils.TimePtr(time.Now()),
	}, nil
}

type occurrence struct {
	SrcType          datasources.SourceType `json:"SourceType" firestore:"SourceType"`
	TgtID            datasources.TargetID `json:"TargetID" firestore:"TargetID"`
	SrcOccurrenceID  string `json:"SourceOccurrenceID" firestore:"SourceOccurrenceID"`
	FormattedDate       string `json:"" firestore:""`
	CreatedAt           *time.Time `json:"" firestore:""`
	ModifiedAt          *time.Time `json:"" firestore:""`
	*geofeatures.GeoFeatureSet
	fsDocumentReference *firestore.DocumentRef
}

func (Ω *occurrence) SourceType() datasources.SourceType {
	return Ω.SrcType
}

func (Ω *occurrence) TargetID() datasources.TargetID {
	return Ω.TgtID
}

func (Ω *occurrence) SourceOccurrenceID() string {
	return Ω.SrcOccurrenceID
}

func (Ω *occurrence) LocationKey() (string, error) {
	if Ω == nil || Ω.GeoFeatureSet == nil {
		return "", errors.New("Nil Occurrence")
	}
	if Ω.GeoFeatureSet == nil {
		return "", errors.New("Nil FeatureSet")
	}
	return Ω.GeoFeatureSet.CoordinateKey() + "|" + Ω.FormattedDate, nil
}


var ErrInvalidDate = errors.New("Invalid Date")

func (Ω *occurrence) SetGeospatial(lat, lng float64, date string, coordinatesEstimated bool) error {

	var err error
	// GeoFeatureSet placeholder should validate for decimal places.
	Ω.GeoFeatureSet, err = geofeatures.NewGeoFeatureSet(lat, lng, coordinatesEstimated)
	if err != nil {
		return err
	}

	if len(date) != 8 {
		return errors.Wrapf(ErrInvalidDate, "Date [%s] must be in format YYYYMMDD", date)
	}

	intDate, err := strconv.Atoi(date)
	if err != nil || intDate == 0 {
		return errors.Wrapf(ErrInvalidDate, "Date [%s] must be in format YYYYMMDD", date)
	}

	if intDate < 19600101 {
		return errors.Wrapf(ErrInvalidDate, "Date [%s] must be after 1960", date)
	}

	// TODO: Reconsider the time zone of each source.
	// Not going to store time because it's useless for the date model.
	// But do need to cast time to the local for that coordinate, to be certain we have the correct day.
	//if Ω.Lat() == 0 || Ω.Lng() == 0 {
	//	return errors.New("Could not calculate timezone because occurrence location is invalid")
	//}
	//tz, err := time.LoadLocation(latlong.LookupZoneName(Ω.Lat(), Ω.Lng()))
	//if err != nil {
	//	return errors.Wrap(err, "could not load location")
	//}
	//
	//loc := date.Location().String()
	//
	//if tz.String() != loc && (loc != "" && loc != "Local") {
	//	return errors.Newf("Locations [%s, %s] are not equal: %s", tz.String(), date.Location().String(), date)
	//}

	Ω.FormattedDate = date

	return nil
}

type OccurrenceAggregation struct {
	collisions int
	iterator_position int
	sync.Mutex
	list []Occurrence
}

func NewOccurrenceAggregation() *OccurrenceAggregation {
	oa := OccurrenceAggregation{
		list: []Occurrence{},
	}
	return &oa
}

func (Ω *OccurrenceAggregation) Collisions() int {
	Ω.Lock()
	defer Ω.Unlock()
	return Ω.collisions
}

func (Ω *OccurrenceAggregation) Count() int {
	if Ω == nil {
		return 0
	}
	Ω.Lock()
	defer Ω.Unlock()
	if Ω.list == nil {
		return 0
	}
	return len(Ω.list)
}

func (Ω *OccurrenceAggregation) Merge(æ *OccurrenceAggregation) error {
	for _, o := range æ.list {
		if err := Ω.AddOccurrence(o); err != nil && !utils.ContainsError(err, ErrCollision) {
			return err
		}
	}
	Ω.collisions += æ.collisions
	return nil
}

var ErrCollision = errors.New("Occurrence Collision")
func (Ω *OccurrenceAggregation) AddOccurrence(b Occurrence) error {

	if b == nil {
		return nil
	}

	bKey, err := b.LocationKey()
	if err != nil {
		return err
	}

	bSourceType := b.SourceType()

	Ω.Lock()
	defer Ω.Unlock()

	if Ω.list == nil {
		Ω.list = []Occurrence{}
	}

	for i := range Ω.list {
		aKey, err := Ω.list[i].LocationKey()
		if err != nil {
			return err
		}
		if aKey != bKey {
			continue
		}
		aSourceType := Ω.list[i].SourceType()

		fmt.Println("Warning: Collision",
			aKey,
			"[" + fmt.Sprint(Ω.list[i].SourceType(), ",", Ω.list[i].TargetID(), ",", Ω.list[i].SourceOccurrenceID()) + "]",
			"[" + fmt.Sprint(b.SourceType(), ",", b.TargetID(), ",", b.SourceOccurrenceID()) + "]")

		if aSourceType != bSourceType && bSourceType == datasources.TypeGBIF {
			Ω.list[i] = b
		}

		return ErrCollision
	}

	Ω.list = append(Ω.list, b)

	return nil


}


// TODO: Should periodically check all occurrences for consistency.

func (Ω *OccurrenceAggregation) Upload(cxt context.Context, florastore store.FloraStore) error {

	limiter := ratelimit.New(100)

	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _, _o := range Ω.list {
			o := _o
			limiter.Take()
			tmb.Go(func() error{
				transaction, err := o.UpsertTransaction(florastore)
				if err != nil {
					return err
				}
				return florastore.FirestoreTransaction(cxt, transaction)
			})
		}
		return nil
	})
	return tmb.Wait()
}

func (Ω *occurrence) docRef(florastore store.FloraStore) *firestore.DocumentRef {
	id := fmt.Sprintf("%s-%s-%s", Ω.SourceType(), Ω.TargetID(), Ω.SourceOccurrenceID())
	return florastore.FirestoreCollection(store.CollectionOccurrences).Doc(id)
}

func (Ω *occurrence) UpsertTransaction(florastore store.FloraStore) (store.FirestoreTransactionFunc, error) {

	if !Ω.SrcType.Valid() || !Ω.TgtID.Valid(Ω.SrcType) || strings.TrimSpace(Ω.SrcOccurrenceID) == "" {
		return nil, errors.Newf("Invalid firestore reference ID: %s, %s, %s", Ω.SourceType(), Ω.TargetID(), Ω.SourceOccurrenceID())
	}

	return func(cxt context.Context, tx *firestore.Transaction) error {

		existingDoc, err := tx.Get(Ω.docRef(florastore))
		notFound := (err != nil && strings.Contains(err.Error(), "not found"))
		if !notFound && err != nil {
			return err
		}

		exists := !notFound

		q, err := Ω.CoordinateQuery(florastore.FirestoreCollection(store.CollectionOccurrences))
		if err != nil {
			return err
		}
		locationQuery := q.Where("FormattedDate", "==", Ω.FormattedDate)

		imbricates, err := tx.Documents(locationQuery).GetAll()
		if err != nil {
			return errors.Wrap(err, "Error searching for a list of possibly overlapping occurrences")
		}

		if len(imbricates) > 1 {
			return errors.Newf("Unexpected: multiple imbricates found for occurrence with location [%f, %f, %s]", Ω.Lat(), Ω.Lng(), Ω.FormattedDate)
		}

		isImbricative := len(imbricates) > 0

		if exists && isImbricative {
			// This suggests the location has changed somewhere. Update code if we see this.
			return errors.Newf("Unexpected: occurrence with id [%s] exists and is imbricative to another doc [%s]", existingDoc.Ref.ID, imbricates[0].Ref.ID)
		}

		if isImbricative {

			// TODO: Be wary of cases in which there are occurrence of two different species in the same spot. Not sure if this will come up.

			fmt.Println(fmt.Sprintf("Warning: Imbricative Occurrence Locations [%s, %s]", existingDoc.Ref.ID, imbricates[0].Ref.ID))

			imbricate := occurrence{}
			if err := imbricates[0].DataTo(&imbricate); err != nil {
				return errors.Wrap(err, "Could not cast occurrence")
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

		if exists {
			Ω.CreatedAt = nil
		}

		// Should be safe to override with new record
		if err := tx.Set(Ω.docRef(florastore), Ω); err != nil {
			return errors.Wrap(err, "Could not set")
		}

		return nil

	}, nil
}

func (Ω *OccurrenceAggregation) MarshalJSON() ([]byte, error) {
	return json.Marshal(Ω.list)
}

func (Ω *OccurrenceAggregation) UnmarshalJSON(b []byte) error {
	list := []*occurrence{}
	if err := json.Unmarshal(b, &list); err != nil {
		return err
	}
	res := []Occurrence{}
	for _, o := range list {
		res = append(res, Occurrence(o))
	}
	Ω.list = res
	return nil
}