package fetcher

import (
	"bitbucket.org/heindl/cxt"
	. "bitbucket.org/heindl/malias"
	"bitbucket.org/heindl/occurrences/fetcher"
	"bitbucket.org/heindl/species"
	"encoding/json"
	"github.com/bitly/go-nsq"
	"github.com/dropbox/godropbox/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type SpeciesFetcher struct {
	*cxt.Context
}

func init() {
	cxt.RegisterCollection(cxt.CollectionInitiator{
		Name: cxt.SpeciesColl,
	})
}

func (this *SpeciesFetcher) HandleMessage(m *nsq.Message) error {

	s := species.Species{
		CanonicalName: species.CanonicalName(m.Body),
	}
	err := this.Mongo.Coll(cxt.SpeciesColl).Find(s).One(&s)
	if err != nil && err != mgo.ErrNotFound {
		return errors.Wrap(err, "could not get taxon from mongo")
	}
	if err == nil {
		return this.queueOccurrenceFetch(s)
	}

	subspecies, err := gatherSubspecies(s.CanonicalName)
	if err != nil {
		return err
	}

	for _, s := range subspecies {

		if _, err := this.Mongo.Coll(cxt.SpeciesColl).Upsert(M{"canonicalName": s.CanonicalName}, M{
			"$addToSet": M{
				"sources": M{"$each": s.Sources},
			},
			"$set": bson.M{
				"lastModified": time.Now(),
			},
		}); err != nil {
			errors.Wrap(err, "could not upsert taxon")
		}

		// Queue data fetch.
		if err := this.Producer.Publish(cxt.NSQSpeciesDataFetch, []byte(s.CanonicalName)); err != nil {
			return errors.Wrapf(err, "could not publish message[%s]", cxt.NSQSpeciesDataFetch)
		}

		if err := this.queueOccurrenceFetch(s); err != nil {
			return err
		}

	}

	return nil

}

func (this *SpeciesFetcher) queueOccurrenceFetch(s species.Species) error {

	for _, src := range s.Sources {

		b, err := json.Marshal(fetcher.OccurrenceFetchQuery{
			Since:  s.LastModified,
			Source: src,
		})
		if err != nil {
			return errors.Wrap(err, "could not marshal taxon query")
		}

		if err := this.Producer.Publish(cxt.NSQOccurrencesFetch, b); err != nil {
			return errors.Wrapf(err, "could not publish message[%s]", cxt.NSQOccurrencesFetch)
		}

	}

	return nil
}
