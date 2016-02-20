package metafetcher

import (
	"bitbucket.org/heindl/cxt"
	"bitbucket.org/heindl/logkeys"
	. "bitbucket.org/heindl/malias"
	"bitbucket.org/heindl/species"
	"github.com/bitly/go-nsq"
	"github.com/dropbox/godropbox/errors"
	"github.com/heindl/eol"
	"strings"
)

type SpeciesMetaFetchHandler struct {
	*cxt.Context
}

func (s *SpeciesMetaFetchHandler) HandleMessage(m *nsq.Message) (err error) {

	scientificName := strings.Trim(string(m.Body), `"`)

	if scientificName == "" {
		return nil
	}

	results, err := eol.Search(eol.SearchQuery{
		Query: scientificName,
		Limit: 10,
	})
	if err != nil {
		s.Log.WithFields(M{logkeys.CanonicalName: scientificName}.Fields()).Error("could not search encyclopedia of life")
		return err
	}

	if len(results) == 0 {
		s.Log.WithFields(M{logkeys.CanonicalName: scientificName}.Fields()).Warn("no search results found from the encyclopedia of life")
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
		s.Log.WithFields(M{
			logkeys.CanonicalName: scientificName,
		}.Fields()).Warn("no page identifier found from the encyclopedia of life")
		return nil
	}

	q := species.Species{
		CanonicalName: species.CanonicalName(scientificName),
	}

	if _, err := s.Mongo.Coll(cxt.SpeciesColl).Upsert(q, s.genUpdateFromPage(highest)); err != nil {
		return errors.Wrap(err, "could not update taxon data")
	}

	return nil

}

func (s *SpeciesMetaFetchHandler) genUpdateFromPage(p eol.PageResponse) M {

	u := M{
		"$addToSet": M{
			"sources": species.Source{
				Type:     species.SourceTypeEOL,
				IndexKey: species.IndexKey(p.Identifier),
			},
		},
	}

	if len(p.Texts()) == 0 || len(p.Images()) == 0 {
		return u
	}

	set := species.Species{}

	if len(p.Texts()) > 0 {
		set.Description = &species.Media{
			Source: p.Texts()[0].Source,
			Value:  p.Texts()[0].Value,
		}
	}

	if len(p.Images()) > 0 {
		set.Image = &species.Media{
			Source: p.Images()[0].Source,
			Value:  p.Images()[0].Value,
		}
	}

	u["$set"] = set
	return u
}
