package occurrences

import (
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/store"
	"bitbucket.org/heindl/process/terra/geoembed"
	"cloud.google.com/go/firestore"
	"encoding/json"
	"errors"
	"fmt"
	dropboxError "github.com/dropbox/godropbox/errors"
	"strconv"
)

// Occurrence represents both a species occurrence and random point.
type Occurrence interface {
	ID() (string, error)
	Collection(florastore store.FloraStore) (*firestore.CollectionRef, error)
	SourceType() datasources.SourceType
	TargetID() datasources.TargetID
	LocationKey() (string, error)
	SourceOccurrenceID() string
	UpsertTransactionFunc(florastore store.FloraStore) (store.FirestoreTransactionFunc, error)
	SetGeoSpatial(lat, lng float64, date string, coordinatesEstimated bool) error
	MarshalJSON() ([]byte, error)
	Coordinates() (lat, lng float64, err error)
	Date() (string, error)
}

// NewOccurrence creates and validates a new one.
func NewOccurrence(srcType datasources.SourceType, targetID datasources.TargetID, occurrenceID string) (Occurrence, error) {
	if !srcType.Valid() {
		return nil, dropboxError.Newf("Invalid source type [%s]", srcType)
	}
	if !targetID.Valid(srcType) {
		return nil, dropboxError.Newf("Invalid target id [%s]", targetID)
	}
	if occurrenceID == "" {
		return nil, dropboxError.Newf("Invalid occurrence id")
	}

	return &occurrence{
		SrcType:         srcType,
		TgtID:           targetID,
		SrcOccurrenceID: occurrenceID,
	}, nil
}

type occurrence struct {
	SrcType         datasources.SourceType  `json:"SourceType"`
	TgtID           datasources.TargetID    `json:"TargetID"`
	SrcOccurrenceID string                  `json:"SourceOccurrenceID"`
	FormattedDate   string                  `json:""`
	GeoFeatureSet   *geoembed.GeoFeatureSet `json:""`
}

func (Ω *occurrence) Collection(florastore store.FloraStore) (*firestore.CollectionRef, error) {
	if Ω.SourceType() == datasources.TypeRandom {
		return florastore.FirestoreCollection(store.CollectionRandom)
	}
	return florastore.FirestoreCollection(store.CollectionOccurrences)
}

func (Ω *occurrence) MarshalJSON() ([]byte, error) {
	o := *Ω
	return json.Marshal(o)

	//gb, err := json.Marshal(o.GeoFeatureSet)
	//if err != nil {
	//	return nil, err
	//}
	//
	//gm := map[string]interface{}{}
	//if err := json.Unmarshal(gb, &gm); err != nil {
	//	return nil, err
	//}
	//
	//o.GeoFeatureSet = nil
	//
	//ob, err := json.Marshal(o)
	//if err != nil {
	//	return nil, err
	//}
	//
	//om := map[string]interface{}{}
	//if err := json.Unmarshal(ob, &om); err != nil {
	//	return nil, err
	//}
	//
	//for k, v := range gm {
	//	if _, ok := om[k]; ok {
	//		return nil, errors.Newf("Occurrence field [%s] collides with GeoFeatures field", k)
	//	}
	//	om[k] = v
	//}
	//
	//return json.Marshal(om)
}

func (Ω *occurrence) Coordinates() (lat, lng float64, err error) {
	return Ω.GeoFeatureSet.Lat(), Ω.GeoFeatureSet.Lng(), nil
}

func (Ω *occurrence) UnmarshalJSON(b []byte) error {

	gf := geoembed.GeoFeatureSet{}
	if err := json.Unmarshal(b, &gf); err != nil {
		return err
	}

	o := occurrence{}
	if err := json.Unmarshal(b, &o); err != nil {
		return err
	}

	o.GeoFeatureSet = &gf

	*Ω = o

	return nil
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

func (Ω *occurrence) Date() (string, error) {
	if len(Ω.FormattedDate) != 8 {
		return "", dropboxError.Newf("Invalid Occurrence Date [%s]", Ω.FormattedDate)
	}
	return Ω.FormattedDate, nil
}

func (Ω *occurrence) LocationKey() (string, error) {
	if Ω == nil {
		return "", dropboxError.New("Occurrence is Invalid")
	}

	if Ω.GeoFeatureSet == nil {
		return "", dropboxError.New("Occurrence GeoFeatureSet is Invalid")
	}

	coordKey, err := Ω.GeoFeatureSet.CoordinateKey()
	if err != nil {
		return "", err
	}

	date, err := Ω.Date()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s|%s", coordKey, date), nil
}

// ErrInvalidDate flags a date that isn't in the format 20060101
var ErrInvalidDate = errors.New("invalid date")

// SetGeoSpatial creates and adds the occurrence GeoFeatureSet.
func (Ω *occurrence) SetGeoSpatial(lat, lng float64, date string, coordinatesEstimated bool) error {

	var err error
	// GeoFeatureSet placeholder should validate for decimal places.
	Ω.GeoFeatureSet, err = geoembed.NewGeoFeatureSet(lat, lng, coordinatesEstimated)
	if err != nil {
		return err
	}

	if len(date) != 8 {
		return dropboxError.Wrapf(ErrInvalidDate, "Date [%s] must be in format YYYYMMDD", date)
	}

	intDate, err := strconv.Atoi(date)
	if err != nil || intDate == 0 {
		return dropboxError.Wrapf(ErrInvalidDate, "Date [%s] must be in format YYYYMMDD", date)
	}

	if intDate < 19600101 {
		return dropboxError.Wrapf(ErrInvalidDate, "Date [%s] must be after 1960", date)
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
