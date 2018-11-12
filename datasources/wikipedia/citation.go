package wikipedia

import (
	"github.com/heindl/floracast_process/utils"
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"github.com/grokify/html-strip-tags-go"
	"net/url"
	"strings"
	"time"
)

type query struct {
	Query struct {
		Pages []page `json:"pages"`
	} `json:"query"`
}

type page struct {
	Pageid    int        `json:"pageid"`
	Ns        int        `json:"ns"`
	Title     string     `json:"title"`
	Revisions []revision `json:"revisions"`
}

type revision struct {
	Revid     int       `json:"revid"`
	Parentid  int       `json:"parentid"`
	Minor     bool      `json:"minor"`
	User      string    `json:"user"`
	Timestamp time.Time `json:"timestamp"`
	Comment   string    `json:"comment"`
}

// MLACitation fetches and validates a wikipedia page and returns a Modern Language Association (MLA) formatted citation string.
func MLACitation(wikipediaURL string) (string, error) {

	wikipediaURL = strings.TrimSpace(wikipediaURL)

	parsedWikipediaURL, err := url.Parse(wikipediaURL)
	if err != nil {
		return "", errors.Wrapf(err, "Could not parse URL [%s]", wikipediaURL)
	}

	pageTitle := strings.TrimPrefix(parsedWikipediaURL.Path, "/wiki/")
	if pageTitle == "" {
		return "", errors.Newf("Invalid Wikipedia Title [%s]", wikipediaURL)
	}

	path := []string{
		"action=query",
		fmt.Sprintf("titles=%s", pageTitle),
		"prop=revisions",
		"format=json",
		"formatversion=2",
	}

	u := "https://en.wikipedia.org/w/api.php?" + strings.Join(path, "&")

	q := query{}
	if err := utils.RequestJSON(u, &q); err != nil {
		return "", err
	}

	if len(q.Query.Pages) == 0 {
		return "", errors.Newf("Wikipedia page not found [%s]", wikipediaURL)
	}

	title := strip.StripTags(strings.TrimSpace(q.Query.Pages[0].Title))

	if title == "" {
		return "", errors.Newf("Wikipedia page [%s] missing title", wikipediaURL)
	}

	if len(q.Query.Pages[0].Revisions) == 0 {
		return "", errors.Newf("Wikipedia page [%s] missing revisions", wikipediaURL)
	}

	modifiedAt := q.Query.Pages[0].Revisions[0].Timestamp

	if modifiedAt.IsZero() {
		return "", errors.Newf("Could not find Wikipedia modifiedAt time [%s]", wikipediaURL)
	}

	citationComponents := []string{
		"Wikipedia contributors",
		fmt.Sprintf(`"%s"`, title),
		"Wikipedia, The Free Encyclopedia",
		modifiedAt.Format("2 Jan. 2006"),
		"Web",
		time.Now().Format("2 Jan. 2006"),
		fmt.Sprintf("<%s>", wikipediaURL),
	}

	return strings.Join(citationComponents, ". "), nil
}
