package wikipedia

import (
	"github.com/dropbox/godropbox/errors"
	"net/url"
	"fmt"
	"strings"
	"github.com/sadbox/mediawiki"
	"time"
	"github.com/grokify/html-strip-tags-go"
)

func Citation(wikipedia_url string) (string, error) {

	requestURI, err := url.ParseRequestURI(wikipedia_url)
	if err != nil {
		return "", errors.Wrapf(err, "Invalid URL [%s]", wikipedia_url)
	}

	wiki, err := mediawiki.New("https://en.wikipedia.org/w/api.php", "FloracastWikiBot")
	if err != nil {
		return "", errors.Wrapf(err, "Could not fetch MediaWiki [%s]", wikipedia_url)
	}

	pageName := strings.TrimLeft(requestURI.Path, "wiki/")
	wikiPage, err := wiki.Read(pageName)
	if err != nil {
		return "", errors.Wrapf(err, "Could not read wikipedia page [%s] from link [%s]", pageName, wikipedia_url)
	}

	modifiedAt := time.Time{}

	for _, revision := range wikiPage.Revisions{
		if revision.Revid == wikiPage.Lastrevid {
			modifiedAt = revision.Timestamp
			break
		}
	}

	if modifiedAt.IsZero() {
		return "", errors.Newf("Could not find modifiedAt [%s]", wikipedia_url)
	}

	c := strings.Join([]string{
		"Wikipedia contributors",
		fmt.Sprintf(`"%s"`, wikiPage.Title),
		"Wikipedia, The Free Encyclopedia",
		modifiedAt.Format("2 Jan. 2006"),
		"Web",
		time.Now().Format("2 Jan. 2006"),
		fmt.Sprintf("<%s>", wikipedia_url),
	}, ". ")

	return strip.StripTags(c), nil
}