package occurrences

import (
	"bitbucket.org/heindl/taxa/utils"
	"cloud.google.com/go/firestore"
	"golang.org/x/net/context"
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"strings"
	"time"
	"bitbucket.org/heindl/taxa/store"
	"bitbucket.org/heindl/taxa/geofeatures"
	"go.uber.org/ratelimit"
	"gopkg.in/tomb.v2"
	"strconv"
	"sync"
	"bitbucket.org/heindl/taxa/datasources"
	"encoding/json"
)

func NewOccurrence(srcType datasources.SourceType, targetID datasources.TargetID, occurrenceID string) (*Occurrence, error) {
	if !srcType.Valid() {
		return nil, errors.Newf("Invalid source type [%s]", srcType)
	}
	if !targetID.Valid(srcType) {
		return nil, errors.Newf("Invalid target id [%s]", targetID)
	}
	if occurrenceID == "" {
		return nil, errors.Newf("Invalid occurrence id")
	}

	return &Occurrence{
		sourceType: srcType,
		targetID: targetID,
		sourceOccurrenceID: occurrenceID,
		createdAt: utils.TimePtr(time.Now()),
		modifiedAt: utils.TimePtr(time.Now()),
	}, nil
}

type Occurrence struct {
	sourceType          datasources.SourceType
	targetID            datasources.TargetID
	sourceOccurrenceID  string
	formattedDate       string
	createdAt           *time.Time
	modifiedAt          *time.Time
	*geofeatures.GeoFeatureSet
	fsDocumentReference *firestore.DocumentRef
	fsTimeLocationQuery firestore.Query
}

func (Ω *Occurrence) SourceType() datasources.SourceType {
	return Ω.sourceType
}

func (Ω *Occurrence) TargetID() datasources.TargetID {
	return Ω.targetID
}

func (Ω *Occurrence) locationKey() (string, error) {
	if Ω == nil || Ω.GeoFeatureSet == nil {
		return "", errors.New("Nil Occurrence")
	}
	if Ω.GeoFeatureSet == nil {
		return "", errors.New("Nil FeatureSet")
	}
	return Ω.GeoFeatureSet.CoordinateKey() + "|" + Ω.formattedDate, nil
}

const keySourceType = "SourceType"
const keyTargetID = "TargetID"
const keySourceOccurrenceID = "SourceOccurrenceID"
const keyFormattedDate = "FormattedDate"
const keyCreatedAt = "CreatedAt"
const keyModifiedAt = "ModifiedAt"


func (Ω *Occurrence) MarshalJSON() ([]byte, error) {
	m, err := Ω.toMap()
	if err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

func (Ω *Occurrence) toMap() (map[string]interface{}, error) {

	oc := map[string]interface{}{
		keySourceType:    Ω.sourceType,
		keyTargetID:   Ω.targetID,
		keySourceOccurrenceID:       Ω.sourceOccurrenceID,
		keyFormattedDate:      Ω.formattedDate,
		keyCreatedAt:    Ω.createdAt,
		keyModifiedAt:   Ω.modifiedAt,
	}

	gfs, err := Ω.GeoFeatureSet.ToMap()
	if err != nil {
		return nil, err
	}

	for k, v := range gfs {
		if _, ok := oc[k]; ok {
			return nil, errors.Newf("Occurrence field collides with GeoFeatureSet [%s]", k)
		}
		oc[k] = v
	}

	return oc, nil
}

func fromMap(m map[string]interface{}) (*Occurrence, error) {

	o, err := NewOccurrence(
		m[keySourceType].(datasources.SourceType),
		m[keyTargetID].(datasources.TargetID),
		m[keySourceOccurrenceID].(string),
	)
	if err != nil {
		return nil, err
	}

	o.GeoFeatureSet, err = geofeatures.NewGeoFeatureSetFromMap(m)
	if err != nil {
		return nil, err
	}

	o.formattedDate = m[keyFormattedDate].(string)
	if o.formattedDate == "" || len(o.formattedDate) != 8 {
		return nil, errors.New("Invalid formatted date")
	}

	o.createdAt = utils.TimePtr(m[keyCreatedAt].(time.Time))
	o.modifiedAt = utils.TimePtr(m[keyModifiedAt].(time.Time))

	return o, nil
}


var ErrInvalidDate = errors.New("Invalid Date")

func (Ω *Occurrence) SetGeospatial(lat, lng float64, date string, coordinatesEstimated bool) error {

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

	Ω.formattedDate = date

	return nil
}

type OccurrenceAggregation struct {
	collisions int
	iterator_position int
	sync.Mutex
	list []*Occurrence
}

func NewOccurrenceAggregation() *OccurrenceAggregation {
	oa := OccurrenceAggregation{
		list: []*Occurrence{},
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
func (Ω *OccurrenceAggregation) AddOccurrence(b *Occurrence) error {

	if b == nil {
		return nil
	}

	bKey, err := b.locationKey()
	if err != nil {
		return err
	}

	bSourceType := b.sourceType

	Ω.Lock()
	defer Ω.Unlock()

	if Ω.list == nil {
		Ω.list = []*Occurrence{}
	}

	for i := range Ω.list {
		aKey, err := Ω.list[i].locationKey()
		if err != nil {
			return err
		}
		if aKey != bKey {
			continue
		}
		aSourceType := Ω.list[i].sourceType

		fmt.Println("Warning: Collision",
			aKey,
			"[" + fmt.Sprint(Ω.list[i].sourceType, ",", Ω.list[i].targetID, ",", Ω.list[i].sourceOccurrenceID) + "]",
			"[" + fmt.Sprint(b.sourceType, ",", b.targetID, ",", b.sourceOccurrenceID) + "]")

		if aSourceType != bSourceType && bSourceType == datasources.DataSourceTypeGBIF {
			Ω.list[i] = b
		}

		return ErrCollision
	}

	Ω.list = append(Ω.list, b)

	return nil


}


// TODO: Should periodically check all occurrences for consistency.

func (Ω OccurrenceAggregation) Upsert(cxt context.Context, firestoreClient *firestore.Client) error {

	col := firestoreClient.Collection(store.CollectionTypeOccurrences)
	limiter := ratelimit.New(100)

	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _, _o := range Ω.list {
			o := _o
			limiter.Take()
			tmb.Go(func() error{
				if err := o.reference(col); err != nil {
					return err
				}
				return firestoreClient.RunTransaction(cxt, o.upsert)
			})
		}
		return nil
	})
	return tmb.Wait()
}

func (Ω *Occurrence) reference(col *firestore.CollectionRef) error {
	if !Ω.sourceType.Valid() || !Ω.targetID.Valid(Ω.sourceType) || strings.TrimSpace(Ω.sourceOccurrenceID) == "" {
		return errors.Newf("Invalid firestore reference ID: %s, %s, %s", Ω.sourceType, Ω.targetID, Ω.sourceOccurrenceID)
	}
	id := fmt.Sprintf("%s-%s-%s", Ω.sourceType, Ω.targetID, Ω.sourceOccurrenceID)
	Ω.fsDocumentReference = col.Doc(id)
	q, err := Ω.GeoFeatureSet.CoordinateQuery(col)
	if err != nil {
		return err
	}
	Ω.fsTimeLocationQuery = q.Where("FormattedDate", "==", Ω.formattedDate)
	return nil
}

func (Ω *Occurrence) upsert(cxt context.Context, tx *firestore.Transaction) error {
	if Ω.fsDocumentReference == nil {
		return errors.New("Firebase document reference not set")
	}
	existingDoc, err := tx.Get(Ω.fsDocumentReference)
	notFound := (err != nil && strings.Contains(err.Error(), "not found"))
	if !notFound && err != nil {
		return err
	}
	exists := !notFound

	imbricates, err := tx.Documents(Ω.fsTimeLocationQuery).GetAll()
	if err != nil {
		return errors.Wrap(err, "Error searching for a list of possibly overlapping occurrences")
	}

	if len(imbricates) > 1 {
		return errors.Newf("Unexpected: multiple imbricates found for occurrence with location [%f, %f, %s]", Ω.Lat(), Ω.Lng(), Ω.formattedDate)
	}

	isImbricative := len(imbricates) > 0

	if exists && isImbricative {
		// This suggests the location has changed somewhere. Update code if we see this.
		return errors.Newf("Unexpected: occurrence with id [%s] exists and is imbricative to another doc [%s]", existingDoc.Ref.ID, imbricates[0].Ref.ID)
	}

	if isImbricative {

		// TODO: Be wary of cases in which there are occurrence of two different species in the same spot. Not sure if this will come up.

		fmt.Println(fmt.Sprintf("Warning: Imbricative Occurrence Locations [%s, %s]", existingDoc.Ref.ID, imbricates[0].Ref.ID))

		imbricate, err := fromMap(imbricates[0].Data())
		if err != nil {
			return err
		}

		if Ω.sourceType != imbricate.sourceType && imbricate.sourceType == datasources.DataSourceTypeGBIF {
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
	
	updateDoc, err := Ω.toMap()
	if err != nil {
		return err
	}

	if exists {
		delete(updateDoc, keyCreatedAt)
	}

	// Should be safe to override with new record
	if err := tx.Set(Ω.fsDocumentReference, updateDoc); err != nil {
		return errors.Wrap(err, "Could not set")
	}

	return nil

}


//func (Ω *store) UpsertOccurrence(cxt context.Context, o Occurrence) (isNewOccurrence bool, err error) {
//
//	isNewOccurrence = false
//
//	ref, err := Ω.NewOccurrenceDocumentRef(o.TaxonID, o.SourceType, o.TargetID)
//	if err != nil {
//		return false, err
//	}
//
//	bkf := backoff.NewExponentialBackOff()
//	bkf.InitialInterval = time.Second * 1
//	ticker := backoff.NewTicker(bkf)
//	for _ = range ticker.C {
//
//		//o.S2CellIDs, err = s2Cells(o.Location.Latitude, o.Location.Longitude)
//		//if err != nil {
//		//	return false, err
//		//}
//
//		err := Ω.FirestoreClient.RunTransaction(cxt, func(cxt context.Context, tx *firestore.Transaction) error {
//			if _, err := tx.Get(ref); err != nil {
//				if strings.Contains(err.Error(), "not found") {
//					isNewOccurrence = true
//					return tx.Set(ref, o)
//				} else {
//					return err
//				}
//			}
//			fields := o.mergeFields()
//			return tx.Set(ref, o, firestore.Merge(fields...))
//		})
//
//		if err != nil && strings.Contains(err.Error(), "rpc error") {
//			fmt.Println("RPC ERROR", err, utils.JsonOrSpew(o))
//			continue
//		}
//
//		if err != nil {
//			ticker.Stop()
//			return false, errors.Wrap(err, "could not update occurrence")
//		}
//
//		ticker.Stop()
//		break
//	}
//	return isNewOccurrence, nil
//}

//func s2Cells(lat, lng float64) (map[string]bool, error) {
//
//	if lat == 0 || lng == 0 {
//		return nil, errors.New("invalid lat/lng")
//	}
//
//	cell := s2.CellIDFromLatLng(s2.LatLngFromDegrees(lat, lng))
//	cells := map[string]bool{strings.Replace(cell.String(), "/", "_", -1): true}
//	for i:=1; i<14; i++ {
//		cells[strings.Replace(cell.Parent(i).String(), "/", "_", -1)] = true
//	}
//	return cells, nil
//}

//func (Ω *store) GetOccurrences(cxt context.Context, taxonID INaturalistTaxonID) (res OccurrenceCollection, err error) {
//
//	if !taxonID.Valid() {
//		return nil, errors.Newf("invalid taxon id [%s]", taxonID)
//	}
//
//	docs, err := Ω.FirestoreClient.Collection(CollectionTypeOccurrences).
//		Where("INaturalistTaxonID", "==", taxonID).
//		Documents(cxt).
//		GetAll()
//
//	if err != nil {
//		return nil, errors.Wrapf(err, "could not get occurrences with taxon id [%s]", taxonID)
//	}
//
//	for _, doc := range docs {
//		o := Occurrence{}
//		if err := doc.DataTo(&o); err != nil {
//			return nil, errors.Wrap(err, "could not type cast occurrence")
//		}
//		res = append(res, o)
//	}
//
//	return
//}
