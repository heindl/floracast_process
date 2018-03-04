package occurrence

import (
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/terra/geo"
	"bitbucket.org/heindl/process/utils"
	"encoding/json"
	"errors"
	"fmt"
	dropboxErrors "github.com/dropbox/godropbox/errors"
	"sync"
)

// Aggregation handles collecting occurrences, validating them, avoiding collisions,
// and uploading them to FireStore.
type Aggregation struct {
	collisions int
	sync.Mutex
	list []Occurrence
}

// NewOccurrenceAggregation creates one.
func NewOccurrenceAggregation() *Aggregation {
	oa := Aggregation{
		list: []Occurrence{},
	}
	return &oa
}

// Count returns the length of the aggregation.
func (Ω *Aggregation) Count() int {
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

// Merge combines aggregations selects between collisions.
func (Ω *Aggregation) Merge(æ *Aggregation) error {
	if æ == nil {
		return nil
	}
	for _, o := range æ.list {
		if err := Ω.AddOccurrence(o); err != nil && !utils.ContainsError(err, ErrCollision) {
			return err
		}
	}
	Ω.Lock()
	defer Ω.Unlock()
	Ω.collisions += æ.collisions
	return nil
}

func (Ω *Aggregation) indexOf(qKey string) (int, error) {
	for _i := range Ω.list {
		i := _i
		iKey, err := Ω.list[i].LocationKey()
		if err != nil {
			return 0, err
		}

		if qKey == iKey {
			return i, nil
		}
	}
	return -1, nil
}

// ErrCollision warns of a collision.
var ErrCollision = errors.New("Occurrence Collision")

var counter = 0

// AddOccurrence adds a new record to the aggregation and returns error if it's
// an unselected record in collision.
func (Ω *Aggregation) AddOccurrence(q Occurrence) error {

	if q == nil {
		return nil
	}

	qKey, err := q.LocationKey()
	if err != nil {
		return err
	}

	counter++
	fmt.Println("ADD OCCURRENCE", qKey, counter)
	defer func() {
		fmt.Println("FINISHED ADDING OCCURRENCE", qKey, counter)
	}()

	Ω.Lock()
	defer Ω.Unlock()

	if Ω.list == nil {
		Ω.list = []Occurrence{q}
		return nil
	}

	i, err := Ω.indexOf(qKey)
	if err != nil {
		return err
	}
	if i == -1 {
		Ω.list = append(Ω.list, q)
		return nil
	}

	if Ω.list[i].SourceType() == datasources.TypeGBIF && q.SourceType() != datasources.TypeGBIF {
		return dropboxErrors.Wrapf(
			ErrCollision,
			"Key [%s] - [%s] [%s]",
			qKey,
			fmt.Sprint(Ω.list[i].SourceType(), ",", Ω.list[i].TargetID(), ",", Ω.list[i].SourceOccurrenceID()),
			fmt.Sprint(q.SourceType(), ",", q.TargetID(), ",", q.SourceOccurrenceID()),
		)
	}

	Ω.list[i] = q

	return nil

}

// GeoJSON creates a GeoJSON point collection
func (Ω *Aggregation) GeoJSON() ([]byte, error) {
	points := geo.Points{}
	for _, o := range Ω.list {
		lat, lng, err := o.Coordinates()
		if err != nil {
			return nil, err
		}
		p, err := geo.NewPoint(lat, lng)
		if err != nil {
			return nil, err
		}
		date, err := o.Date()
		if err != nil {
			return nil, err
		}
		if err := p.SetProperty("Date", date); err != nil {
			return nil, err
		}
		points = append(points, p)
	}
	return points.GeoJSON()
}

// MarshalJSON will convert record list to JSON
func (Ω *Aggregation) MarshalJSON() ([]byte, error) {
	return json.Marshal(Ω.list)
}

// UnmarshalJSON takes a list of occurrences and creates an aggregation
func (Ω *Aggregation) UnmarshalJSON(b []byte) error {
	list := []*record{}
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
