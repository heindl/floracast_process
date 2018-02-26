package taxa

import (
	"bitbucket.org/heindl/process/datasources/sourcefetchers"
	"golang.org/x/net/context"
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/nameusage/nameusage"
	"gopkg.in/tomb.v2"
	"sync"
	"github.com/grokify/html-strip-tags-go"
	"strings"
)

type photo struct {
	Citation string `json:",omitempty" firestore:",omitempty"`
	Thumbnail string `json:",omitempty" firestore:",omitempty"`
	Large string `json:",omitempty" firestore:",omitempty"`
	Rank int `json:"-" firestore:"-"`
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
			tmb.Go(func() error {
				if src.OccurrenceCount() < 1 {
					return nil
				}

				srcType, err := src.SourceType()
				if err != nil {
					return err
				}

				targetID, err := src.TargetID()
				if err != nil {
					return err
				}

				fetchedPhotos, err := sourcefetchers.FetchPhotos(ctx, srcType, targetID)
				if err != nil {
					return err
				}

				for _, p := range fetchedPhotos {
					rank := 0
					if p.Source() == datasources.TypeINaturalist {
						rank += 1
					}
					if p.Thumbnail() != "" {
						rank += 1
					}
					if p.Citation() != "" {
						rank += 1
					}
					lock.Lock()
					if rank > highestRank {
						highestRank = rank
						leadPhoto = photo{
							Citation:  p.Citation(),
							Large:     p.Large(),
							Thumbnail: p.Thumbnail(),
							Rank:      rank,
						}
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
	Text string `json:",omitempty" firestore:",omitempty"`
	Rank string `json:"-" firestore:"-"`
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
		for _, _src := range srcs {
			src := _src
			tmb.Go(func() error {
				if src.OccurrenceCount() < 1 {
					return nil
				}

				srcType, err := src.SourceType()
				if err != nil {
					return err
				}

				targetID, err := src.TargetID()
				if err != nil {
					return err
				}

				fetchedDescriptions, err := sourcefetchers.FetchDescriptions(ctx, srcType, targetID)
				if err != nil {
					return err
				}
				for _, p := range fetchedDescriptions {
					rank := 0
					if p.Source() == datasources.TypeINaturalist {
						rank += 1
					}
					citation, err := p.Citation()
					if err != nil {
						return err
					}
					citation = strip.StripTags(citation)
					if citation != "" {
						rank += 1
					}
					text, err := p.Text()
					if err != nil {
						return err
					}
					text = strings.TrimSpace(strip.StripTags(text))
					if text == "" {
						continue
					}
					lock.Lock()
					if rank > highestRank || highestRank == 0 {
						highestRank = rank
						leadDescription = description{
							Citation: citation,
							Text: text,
						}
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