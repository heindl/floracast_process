package occurrences

import (
	"fmt"
	"time"
	"bitbucket.org/heindl/species/store"
	"bitbucket.org/heindl/utils"
	"github.com/heindl/gbif"
	"github.com/saleswise/errors/errors"
	"cloud.google.com/go/datastore"
	"strconv"
)

type GBIF int

func (this GBIF) Fetch(begin *time.Time, end time.Time) (store.Occurrences, error) {

	if begin == nil || begin.IsZero() {
		begin = utils.TimePtr(time.Date(1900, time.January, 1, 0, 0, 0, 0, time.UTC))
	}

	q := gbif.OccurrenceSearchQuery{
		TaxonKey:          int(this),
		LastInterpreted:    fmt.Sprintf("%s,%s", begin.Format("2006-01-02"), end.Format("2006-01-02")),
		HasCoordinate:      true,
		HasGeospatialIssue: false,
	}

	results, err := gbif.Occurrences(q)
	if err != nil {
		return nil, errors.Wrap(err, "could not request occurrences")
	}

	if len(results) == 0 {
		return nil, nil
	}

	res := store.Occurrences{}
	for _, r := range results {
		o := this.parse(r)
		if o == nil {
			continue
		}
		res = append(res, o)
	}

	return res, nil
}

func (this GBIF) parse(o gbif.Occurrence) *store.Occurrence {
	// Note that the OccurrenceID, which I originally used, appears to be incomplete, duplicated, and missing in some cases.

	id, err := strconv.ParseInt(o.GbifID, 10, 64)
	if err != nil {
		fmt.Println("error parsing gbifId", err)
		return nil
	}
	if id == 0 {
		return nil
	}
	if o.EventDate == nil || o.EventDate.Time.IsZero() {
		// TODO: Consider reporting malformed occurrence error.
		return nil
	}
	p := datastore.GeoPoint{o.DecimalLatitude, o.DecimalLongitude}
	if !p.Valid() {
		return nil
	}
	return &store.Occurrence{
		Key: datastore.IDKey(store.EntityKindOccurrence, id, nil),
		OccurrenceID: o.OccurrenceID,
		Location: &p,
		Date:      o.EventDate.Time,
		RecordedBy: o.RecordedBy,
		References: o.References,
		CreatedAt: time.Now(),
	}
}
