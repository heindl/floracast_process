package main

import (
	"bitbucket.org/heindl/logkeys"
	. "bitbucket.org/heindl/malias"
	"bitbucket.org/heindl/nsqeco"
	"bitbucket.org/heindl/species"
	"bitbucket.org/heindl/species/store"
	"bitbucket.org/heindl/utils"
	"encoding/json"
	"github.com/bitly/go-nsq"
	"github.com/heindl/gbif"
	"github.com/saleswise/errors/errors"
	"time"
)

func main() {
	producer, err := nsqeco.NewProducer()
	if err != nil {
		panic(err)
	}
	defer producer.Stop()
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
	list, err := this.SpeciesStore.ReadFromCanonicalNames(name)
	if err != nil {
		return err
	}
	if len(list) != 0 {
		return this.queueOccurrenceFetch(list[0])
	}

	subspecies, err := gatherSubspecies(name)
	if err != nil {
		return err
	}

	for _, s := range subspecies {
		if err := this.SpeciesStore.AddSources(s.CanonicalName, s.Sources...); err != nil {
			return err
		}
		if err := this.NSQProducer.Publish(nsqeco.NSQSpeciesMetaFetch, []byte(s.CanonicalName)); err != nil {
			return errors.Wrap(err, "could not publish message").SetState(M{logkeys.Topic: nsqeco.NSQSpeciesMetaFetch, logkeys.CanonicalName: s.CanonicalName})
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
			Since:  &time.Time{}, // Use all of time for now. Eventually use last sync on source data.
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

func gatherSubspecies(name species.CanonicalName) ([]species.Species, error) {

	subspecies, err := gbif.Search(gbif.SearchQuery{
		Q:    string(name),
		Rank: "SUBSPECIES",
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not search gbif")
	}

	m := make(map[string][]int)

	addtoset := func(name string, n gbif.NameUsage) error {
		if n.Key == 0 {
			return errors.Wrapf(err, "no key found for subspecies:\n %v", utils.JsonOrSpew(n))
		}
		if _, ok := m[name]; ok {
			m[name] = utils.AddIntToSet(m[name], n.Key)
		} else {
			m[name] = []int{n.Key}
		}
		return nil
	}

	for _, sub := range subspecies {

		if err := addtoset(sub.CanonicalName, sub); err != nil {
			return nil, err
		}
		s := gbif.Species(sub.Key)
		synonyms, err := s.Synonyms()
		if err != nil {
			return nil, err
		}
		for _, synonym := range synonyms {
			if err := addtoset(sub.CanonicalName, synonym); err != nil {
				return nil, err
			}
		}
	}

	var response []species.Species

	for k, v := range m {
		s := species.Species{
			CanonicalName: species.CanonicalName(k),
		}
		for _, i := range v {
			s.Sources = append(s.Sources, species.Source{
				Type:     species.SourceTypeGBIF,
				IndexKey: species.IndexKey(i),
			})
		}
		response = append(response, s)
	}

	return response, nil

}
