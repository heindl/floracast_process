package occurrences

import (
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/terra/geo"
	"bitbucket.org/heindl/process/utils"
	"encoding/json"
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"sync"
)

type OccurrenceAggregation struct {
	collisions int
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
	if æ == nil {
		return nil
	}
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
			"["+fmt.Sprint(Ω.list[i].SourceType(), ",", Ω.list[i].TargetID(), ",", Ω.list[i].SourceOccurrenceID())+"]",
			"["+fmt.Sprint(b.SourceType(), ",", b.TargetID(), ",", b.SourceOccurrenceID())+"]")

		if aSourceType != bSourceType && bSourceType == datasources.TypeGBIF {
			Ω.list[i] = b
		}

		return ErrCollision
	}

	Ω.list = append(Ω.list, b)

	return nil

}

func (Ω *OccurrenceAggregation) GeoJSON() ([]byte, error) {
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
