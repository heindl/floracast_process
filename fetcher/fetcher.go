package main

import (
	"bitbucket.org/heindl/nsqeco"
	"bitbucket.org/heindl/species"
	"bitbucket.org/heindl/species/store"
	"encoding/json"
	"github.com/bitly/go-nsq"
	"github.com/dropbox/godropbox/errors"
)

func main() {
	producer, err := nsqeco.NewProducer()
	if err != nil {
		panic(err)
	}
	store, err := store.NewSpeciesStore()
	if err != nil {
		panic(err)
	}
	defer store.Close()

	if err := nsqeco.Listen(nsqeco.NSQSpeciesFetch, &SpeciesFetchHandler{
		NSQProducer:  producer,
		SpeciesStore: store,
	}, 10); err != nil {
		panic(err)
	}
	<-make(chan bool)
}

type SpeciesFetchHandler struct {
	NSQProducer  nsqeco.Producer
	SpeciesStore store.SpeciesStore
}

func (this *SpeciesFetchHandler) HandleMessage(m *nsq.Message) error {

	name := species.CanonicalName(m.Body)

	// Index: SpeciesColl.0
	s, err := this.SpeciesStore.ReadFromCanonicalName(name)
	if err != nil {
		return err
	}
	if s != nil {
		return this.queueOccurrenceFetch(*s)
	}

	subspecies, err := gatherSubspecies(name)
	if err != nil {
		return err
	}

	for _, s := range subspecies {
		if err := this.SpeciesStore.AddSources(s.CanonicalName, s.Sources...); err != nil {
			return err
		}
		// Queue data fetch.
		if err := this.NSQProducer.Publish(nsqeco.NSQSpeciesMetaFetch, []byte(s.CanonicalName)); err != nil {
			return errors.Wrapf(err, "could not publish message[%s]", nsqeco.NSQSpeciesMetaFetch)
		}
		if err := this.queueOccurrenceFetch(s); err != nil {
			return err
		}
	}

	return nil

}

func (this *SpeciesFetchHandler) queueOccurrenceFetch(s species.Species) error {

	for _, src := range s.Sources {
		b, err := json.Marshal(nsqeco.OccurrenceFetchQuery{
			Since:  s.ModifiedAt,
			Source: src,
		})
		if err != nil {
			return errors.Wrap(err, "could not marshal taxon query")
		}
		if err := this.NSQProducer.Publish(nsqeco.NSQOccurrenceFetch, b); err != nil {
			return errors.Wrapf(err, "could not publish message[%s]", nsqeco.NSQOccurrenceFetch)
		}
	}

	return nil
}
