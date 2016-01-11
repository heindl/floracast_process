package fetcher

import (
	"bitbucket.org/heindl/cxt"
	"encoding/json"
	"github.com/bitly/go-nsq"
	"github.com/dropbox/godropbox/errors"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
	. "bitbucket.org/heindl/malias"
	"bitbucket.org/heindl/occurrences"
	"bitbucket.org/heindl/taxon"
)

const AllTaxons = "allTaxons"

type TaxonFetcher struct {
	*cxt.Context
}

func init() {
	cxt.RegisterCollection(cxt.CollectionInitiator{
		Name: cxt.TaxonColl,
	})
}

func (s *TaxonFetcher) HandleMessage(m *nsq.Message) (err error) {

	// TODO Does nsq have a deterministic failure?
	// TODO: Set up logging package.

	taxons, err := s.parseTaxonRequest(string(m.Body))
	if err != nil {
		return err
	}

	for _, taxon := range taxons {

		if err := s.queueOccurrenceFetch(taxon); err != nil {
			return err
		}

	}

	return nil

}

func (s *TaxonFetcher) queueOccurrenceFetch(taxon taxon.Taxon) (err error) {

	start := taxon.LastModified
	if start.IsZero() {
		// start, err = time.Parse("20060102", "19840101") // Why not?
		start, err = time.Parse("20060102", "19600101")
		if err != nil {
			return err
		}
	}

	b, err := json.Marshal(occurrences.OccurrenceFetchQuery{
		StartDate: start,
		EndDate:   time.Now(),
		Key:       taxon.ID,
	})
	if err != nil {
		return errors.Wrap(err, "could not marshal taxon query")
	}

	if err := s.Producer.Publish(cxt.NSQOccurrencesFetch, b); err != nil {
		return errors.Wrapf(err, "could not publish message[%s]", cxt.NSQOccurrencesFetch)
	}

	// Should be a separate channel that listens from other boundaries.
	u := bson.M{
		"$set": bson.M{
			"lastModified": time.Now(),
		},
	}

	if _, err := s.Mongo.Coll(cxt.TaxonColl).UpsertId(taxon.ID, u); err != nil {
		return errors.Wrapf(err, "could not set last update time for taxon[%s]", taxon)
	}
	return nil
}

func (s *TaxonFetcher) parseTaxonRequest(request string) ([]taxon.Taxon, error) {

	request = strings.Trim(request, `"`)

	if string(request) == AllTaxons {
		return s.fetchTaxons(nil)
	}

	var ids []string
	for _, str := range strings.Split(string(request), ",") {
		if str == "" {
			continue
		}
		ids = append(ids, str)
	}

	if len(ids) == 0 {
		return nil, nil
	}

	return s.fetchTaxons(ids)

}

func (s *TaxonFetcher) fetchTaxons(ids []string) ([]taxon.Taxon, error) {

	q := M{}

	if len(ids) > 0 {
		q = M{
			"_id": M{
				"$in": ids,
			},
		}
	}

	var txns []taxon.Taxon
	if err := s.Mongo.Coll(cxt.TaxonColl).Find(q).All(&txns); err != nil {
		return nil, errors.Wrap(err, "could not get taxa from mongo")
	}

	if len(txns) == 0 {
		for _, id := range ids {
			k, err := taxon.NewKeyFromString(id)
			if err != nil {
				return nil, err
			}
			txns = append(txns, taxon.Taxon{ID: k})
		}
	}

	return txns, nil

}