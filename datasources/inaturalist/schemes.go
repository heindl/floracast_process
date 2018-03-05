package inaturalist

import (
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/utils"
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/dropbox/godropbox/errors"
	"regexp"
	"strings"
	"sync"
)

type taxonScheme struct {
	SourceType datasources.SourceType
	TargetID   datasources.TargetID
}

var taxonSchemeRegex = regexp.MustCompile(`\(([^\)]+)\)`)

func (Ω taxonID) fetchTaxonSchemes() ([]*taxonScheme, error) {

	if !Ω.Valid() {
		return nil, errors.Newf("Invalid taxonID [%d]", Ω)
	}

	url := fmt.Sprintf("http://www.inaturalist.org/taxa/%d/schemes", Ω)

	r := bytes.NewReader([]byte{})
	if err := utils.RequestJSON(url, r); err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, errors.Wrap(err, "could parse site for goquery")
	}

	parser := schemePageParser{
		schemes:      []*taxonScheme{},
		hrefSelector: "/taxon_schemes/",
	}

	doc.Find(fmt.Sprintf(`a[href*="%s"]`, parser.hrefSelector)).Each(parser.parseHREF)

	if parser.error != nil {
		return nil, parser.error
	}

	if len(parser.schemes) == 0 {
		return nil, nil
	}

	return parser.schemes, nil

}

type schemePageParser struct {
	sync.Mutex
	schemes      []*taxonScheme
	hrefSelector string
	error        error
}

func (Ω *schemePageParser) parseHREF(i int, s *goquery.Selection) {

	srcStr, _ := s.Attr("href")
	srcStr = strings.TrimPrefix(strings.TrimSpace(srcStr), Ω.hrefSelector)
	if srcStr == "" {
		return
	}

	srcType, err := datasources.NewSourceType(srcStr)
	if err != nil {
		Ω.error = err
		return
	}

	targetStr := taxonSchemeRegex.FindString(s.Parent().Text())
	if targetStr == "" {
		return
	}
	targetStr = strings.TrimRight(strings.TrimLeft(targetStr, "("), ")")

	targetID, err := datasources.NewTargetID(targetStr, srcType)
	if err != nil {
		Ω.error = err
		return
	}

	Ω.Lock()
	defer Ω.Unlock()
	Ω.schemes = append(Ω.schemes, &taxonScheme{
		SourceType: srcType,
		TargetID:   targetID,
	})
}
