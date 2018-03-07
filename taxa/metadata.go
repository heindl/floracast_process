package taxa

import (
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/datasources/providers"
	"bitbucket.org/heindl/process/datasources/sourcefetchers"
	"bitbucket.org/heindl/process/nameusage/nameusage"
	"context"
	"github.com/grokify/html-strip-tags-go"
	"gopkg.in/tomb.v2"
	"strings"
	"sync"
)

type photo struct {
	Citation  string `json:",omitempty" firestore:",omitempty"`
	Thumbnail string `json:",omitempty" firestore:",omitempty"`
	Large     string `json:",omitempty" firestore:",omitempty"`
	rank      int
}

func parsePhoto(p providers.Photo) (*photo, error) {
	rank := 0
	if p.SourceType() == datasources.TypeINaturalist {
		rank++
	}
	if p.Thumbnail() != "" {
		rank++
	}
	if p.Citation() != "" {
		rank++
	}
	return &photo{
		Citation:  p.Citation(),
		Large:     p.Large(),
		Thumbnail: p.Thumbnail(),
		rank:      rank,
	}, nil
}

func fetchPhoto(ctx context.Context, usage nameusage.NameUsage) (*photo, error) {

	srcs, err := usage.Sources(datasources.TypeGBIF, datasources.TypeINaturalist)
	if err != nil {
		return nil, err
	}

	highestRank := 0
	leadPhoto := photo{}
	lock := sync.Mutex{}

	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _, _src := range srcs {
			src := _src
			if src.OccurrenceCount() < 1 {
				continue
			}
			tmb.Go(func() error {

				fetchedPhotos, err := sourcefetchers.FetchPhotos(ctx, src.SourceType, src.TargetID)
				if err != nil {
					return err
				}

				for _, p := range fetchedPhotos {
					pho, err := parsePhoto(p)
					if err != nil {
						return err
					}
					lock.Lock()
					if pho.rank > highestRank {
						highestRank = pho.rank
						leadPhoto = *pho
					}
					lock.Unlock()
				}
				return nil
			})
		}
		return nil
	})
	if err := tmb.Wait(); err != nil {
		return nil, err
	}

	// TODO: Sort and clip to five here.
	return &leadPhoto, nil
}

type description struct {
	Citation string `json:",omitempty" firestore:",omitempty"`
	Text     string `json:",omitempty" firestore:",omitempty"`
	rank     int
}

func parseDescription(p providers.Description) (*description, error) {
	rank := 0
	if p.SourceType() == datasources.TypeINaturalist {
		rank++
	}
	citation, err := p.Citation()
	if err != nil {
		return nil, err
	}
	if citation != "" {
		rank++
	}
	text, err := p.Text()
	if err != nil {
		return nil, err
	}
	text = strings.TrimSpace(strip.StripTags(text))
	if text == "" {
		return nil, nil
	}
	return &description{
		Text:     text,
		Citation: citation,
		rank:     rank,
	}, nil
}

func fetchDescription(ctx context.Context, usage nameusage.NameUsage) (*description, error) {

	srcs, err := usage.Sources(datasources.TypeGBIF, datasources.TypeINaturalist)
	if err != nil {
		return nil, err
	}

	highestRank := 0
	leadDescription := description{}
	lock := sync.Mutex{}

	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _, 𝝨 := range srcs {
			src := 𝝨
			if src.OccurrenceCount() < 1 {
				continue
			}
			tmb.Go(func() error {
				fetchedDescriptions, err := sourcefetchers.FetchDescriptions(ctx, src.SourceType, src.TargetID)
				if err != nil {
					return err
				}
				for _, p := range fetchedDescriptions {
					desc, err := parseDescription(p)
					if err != nil {
						return err
					}
					lock.Lock()
					if desc.rank > highestRank || highestRank == 0 {
						highestRank = desc.rank
						leadDescription = *desc
					}
					lock.Unlock()
				}
				return nil
			})
		}
		return nil
	})
	if err := tmb.Wait(); err != nil {
		return nil, err
	}

	return &leadDescription, nil
}
