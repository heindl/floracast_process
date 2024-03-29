package occurrence

import (
	"github.com/heindl/floracast_process/datasources"
	"github.com/heindl/floracast_process/nameusage/nameusage"
	"github.com/heindl/floracast_process/store"
	"github.com/heindl/floracast_process/terra/geoembed"
	"cloud.google.com/go/firestore"
	"encoding/json"
	"errors"
	"fmt"
	dropboxError "github.com/dropbox/godropbox/errors"
	"strconv"
)

// Occurrences represents both a species record and random point.
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
func NewOccurrence(nameUsageID *nameusage.ID, srcType datasources.SourceType, targetID datasources.TargetID, occurrenceID string) (Occurrence, error) {
	if !srcType.Valid() {
		return nil, dropboxError.Newf("Invalid source type [%s]", srcType)
	}
	if !targetID.Valid(srcType) {
		return nil, dropboxError.Newf("Invalid target id [%s]", targetID)
	}
	if occurrenceID == "" {
		return nil, dropboxError.Newf("Invalid record id")
	}

	return &record{
		NameUsageID:     nameUsageID,
		SrcType:         srcType,
		TgtID:           targetID,
		SrcOccurrenceID: occurrenceID,
	}, nil
}

type record struct {
	NameUsageID     *nameusage.ID           `json:",omitempty"` // Omitempty because will not appear if data is random.
	SrcType         datasources.SourceType  `json:"SourceType"`
	TgtID           datasources.TargetID    `json:"TargetID"`
	SrcOccurrenceID string                  `json:"SourceOccurrenceID"`
	FormattedDate   string                  `json:"FormattedDate"`
	FormattedMonth  string                  `json:"FormattedMonth"`
	GeoFeatureSet   *geoembed.GeoFeatureSet `json:"GeoFeatureSet"`
}

func (Ω *record) Collection(florastore store.FloraStore) (*firestore.CollectionRef, error) {
	if Ω.SourceType() == datasources.TypeRandom {
		return florastore.FirestoreCollection(store.CollectionRandom)
	}
	return florastore.FirestoreCollection(store.CollectionOccurrences)
}

func (Ω *record) MarshalJSON() ([]byte, error) {
	o := *Ω
	return json.Marshal(o)
}

func (Ω *record) Coordinates() (lat, lng float64, err error) {
	return Ω.GeoFeatureSet.Lat(), Ω.GeoFeatureSet.Lng(), nil
}

//func (Ω *record) UnmarshalJSON(b []byte) error {
//
//	gf := geoembed.GeoFeatureSet{}
//	if err := json.Unmarshal(b, &gf); err != nil {
//		return err
//	}
//
//	o := record{}
//	if err := json.Unmarshal(b, &o); err != nil {
//		return err
//	}
//
//	o.GeoFeatureSet = &gf
//
//	*Ω = o
//
//	return nil
//}

func (Ω *record) SourceType() datasources.SourceType {
	return Ω.SrcType
}

func (Ω *record) TargetID() datasources.TargetID {
	return Ω.TgtID
}

func (Ω *record) SourceOccurrenceID() string {
	return Ω.SrcOccurrenceID
}

func (Ω *record) Date() (string, error) {
	if len(Ω.FormattedDate) != 8 {
		return "", dropboxError.Newf("Invalid Occurrences Date [%s]", Ω.FormattedDate)
	}
	return Ω.FormattedDate, nil
}

func (Ω *record) LocationKey() (string, error) {
	if Ω == nil {
		return "", dropboxError.New("Occurrences is Invalid")
	}

	if Ω.GeoFeatureSet == nil {
		return "", dropboxError.New("Occurrences GeoFeatureSet is Invalid")
	}

	date, err := Ω.Date()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s|%s", Ω.GeoFeatureSet.CoordinateToken(), date), nil
}

// ErrInvalidDate flags a date that isn't in the format 20060101
var ErrInvalidDate = errors.New("invalid date")

// SetGeoSpatial creates and adds the record GeoFeatureSet.
func (Ω *record) SetGeoSpatial(lat, lng float64, date string, coordinatesEstimated bool) error {

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
	//	return errors.New("Could not calculate timezone because record location is invalid")
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
	Ω.FormattedMonth = date[4:6]

	return nil
}
