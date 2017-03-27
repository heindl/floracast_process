package main

import (
	"bitbucket.org/heindl/logkeys"
	. "bitbucket.org/heindl/malias"
	"bitbucket.org/heindl/nsqeco"
	"bitbucket.org/heindl/species"
	"bitbucket.org/heindl/species/store"
	"github.com/bitly/go-nsq"
	"github.com/dropbox/godropbox/errors"
	"github.com/heindl/eol"
	"github.com/omidnikta/logrus"
	"strings"
)

func main() {
	store, err := store.NewSpeciesStore()
	if err != nil {
		panic(err)
	}
	defer store.Close()
	if err := nsqeco.Listen(nsqeco.NSQSpeciesMetaFetch, &SpeciesMetaFetchHandler{
		Log:          logrus.New(),
		SpeciesStore: store,
	}, 10, nsqeco.DefaultConfig()); err != nil {
		panic(err)
	}
	<-make(chan bool)

}

type SpeciesMetaFetchHandler struct {
	Log          *logrus.Logger
	SpeciesStore store.SpeciesStore
}

func (this *SpeciesMetaFetchHandler) HandleMessage(m *nsq.Message) (err error) {

	name := species.CanonicalName(strings.Trim(string(m.Body), `"`))

	if name == "" {
		return nil
	}

	results, err := eol.Search(eol.SearchQuery{
		Query: string(name),
		Limit: 10,
	})
	if err != nil {
		this.Log.WithFields(M{logkeys.CanonicalName: name}.Fields()).Error("could not search encyclopedia of life")
		return err
	}

	if len(results) == 0 {
		this.Log.WithFields(M{logkeys.CanonicalName: name}.Fields()).Warn("no search results found from the encyclopedia of life")
		return nil
	}

	// The first result should be the most relevant, but check the top ten for the highest score.

	var highest eol.PageResponse

	for _, r := range results {

		page, err := eol.Page(eol.PageQuery{
			ID:      r.ID,
			Images:  1,
			Text:    1,
			Details: true,
		})
		if err != nil {
			return errors.Wrapf(err, "could not find page query from id[%v]", r.ID)
		}
		if page.RichnessScore > highest.RichnessScore {
			highest = *page
		}
	}

	if highest.Identifier == 0 {
		this.Log.WithFields(M{
			logkeys.CanonicalName: name,
		}.Fields()).Warn("no page identifier found from the encyclopedia of life")
		return nil
	}

	if err := this.SpeciesStore.AddSources(species.CanonicalName(name), species.Source{
		Type:     species.SourceTypeEOL,
		IndexKey: species.IndexKey(highest.Identifier),
	}); err != nil {
		return err
	}

	if len(highest.Texts()) > 0 {
		if err := this.SpeciesStore.SetDescription(name, &species.Media{
			Source: highest.Texts()[0].Source,
			Value:  highest.Texts()[0].Value,
		}); err != nil {
			return err
		}
	}

	if len(highest.Images()) > 0 {
		if err := this.SpeciesStore.SetImage(name, &species.Media{
			Source: highest.Images()[0].Source,
			Value:  highest.Images()[0].Value,
		}); err != nil {
			return err
		}
	}

	return nil
}
