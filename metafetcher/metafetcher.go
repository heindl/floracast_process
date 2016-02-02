package metafetcher

import (
	"bitbucket.org/heindl/cxt"
	"bitbucket.org/heindl/logkeys"
	. "bitbucket.org/heindl/malias"
	"bitbucket.org/heindl/species"
	"bitbucket.org/heindl/utils"
	"fmt"
	"github.com/bitly/go-nsq"
	"github.com/dropbox/godropbox/errors"
	"github.com/heindl/eol"
	"strings"
)

type MetaFetcher struct {
	*cxt.Context
}

func (s *MetaFetcher) HandleMessage(m *nsq.Message) (err error) {

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

	fmt.Println(utils.JsonOrSpew(highest))

	if highest.Identifier == 0 {
		s.Log.WithFields(M{
			logkeys.CanonicalName: scientificName,
		}.Fields()).Warn("no page identifier found from the encyclopedia of life")
		return nil
	}

	q := species.Species{
		CanonicalName: species.CanonicalName(scientificName),
	}

	fmt.Println(s.genUpdateFromPage(highest))

	if _, err := s.Mongo.Coll(cxt.SpeciesColl).Upsert(q, s.genUpdateFromPage(highest)); err != nil {
		return errors.Wrap(err, "could not update taxon data")
	}

	return nil

}

func (s *MetaFetcher) genUpdateFromPage(p eol.PageResponse) M {

	u := M{
		"$addToSet": M{
			"sources": species.Source{
				Type:     species.SourceTypeEOL,
				IndexKey: p.Identifier,
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
