package mushroom_observer

import (
	"fmt"
	"strings"
	//"bitbucket.org/heindl/taxa/store"
	//"golang.org/x/net/context"
	"time"
	"bitbucket.org/heindl/taxa/taxa/name_usage"
	"net/url"
	"bitbucket.org/heindl/taxa/utils"
	"github.com/dropbox/godropbox/errors"
	"bitbucket.org/heindl/taxa/store"
	"gopkg.in/tomb.v2"
	"sync"
)

type MushroomObserverQueryResult struct {
	Version         float64   `json:"version"`
	RunDate         time.Time `json:"run_date"`
	Query           string    `json:"query"`
	NumberOfRecords int       `json:"number_of_records"`
	NumberOfPages   int       `json:"number_of_pages"`
	PageNumber      int       `json:"page_number"`
	Results         []*MushroomObserverTaxonResult `json:"results"`
	RunTime float64 `json:"run_time"`
}

func MatchCanonicalNames(names ...string) ([]*name_usage.NameUsageSource, error) {

	//TODO: If names are three, consider adding var. "Cantharellus cibarius var. cibarius"
	// Only if missing in parent.

	nameResponse := []*name_usage.NameUsageSource{}
	locker := sync.Mutex{}

	tmb := tomb.Tomb{}
	tmb.Go(func() error{
		for _, _name := range names {
			name := _name
			tmb.Go(func() error {

				parameters := strings.Join([]string{
					//fmt.Sprintf("updated_at=%s-%s", time.Now().Add(time.Hour * 24 * 30 * -120).Format("20060102"), time.Now().Format("20060102")),
					"format=json",
					"is_deprecated=false",
					//"ok_for_export=true",
					//"has_classification=true",
					fmt.Sprintf("name=%s", url.QueryEscape(name)),
					//"has_synonyms=true",
					"detail=high",
					//fmt.Sprintf("rank=%s", rank),
				}, "&")

				apiURL := "http://mushroomobserver.org/api/names?" + parameters

				fmt.Println(apiURL)

				var queryResult MushroomObserverQueryResult
				if err := utils.RequestJSON(apiURL, &queryResult); err != nil {
					return errors.Wrapf(err, "could not fetch name from mushroom observer [%s]", name)
				}

				// Since we are always fetching by a specific name, it will never be more than one page.
				// But just to be sure ...
				if queryResult.NumberOfPages > 1 {
					return errors.Newf("Unexpected: multiple pages returned from MushroomObserver name query [%s]", name)
				}

				if queryResult.NumberOfRecords == 0 {
					return nil
				}

				for _, r := range queryResult.Results {
					usage, err := parseTaxonResult(r)
					if err != nil {
						return err
					}
					if usage == nil {
						continue
					}
					locker.Lock()
					nameResponse = append(nameResponse, usage)
					locker.Unlock()
				}
				return nil
			})
		}
		return nil
	})

	if err := tmb.Wait(); err != nil {
		return nil, err
	}

	return nameResponse, nil

}

type MushroomObserverTaxonResult struct {
	ID            int           `json:"id"`
	Type          string        `json:"type"`
	Name          string        `json:"name"`
	Author        string        `json:"author"`
	Rank          string        `json:"rank"`
	Deprecated    bool          `json:"deprecated"`
	Misspelled    bool          `json:"misspelled"`
	Citation      string        `json:"citation"`
	Notes         string        `json:"notes"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	NumberOfViews int           `json:"number_of_views"`
	LastViewed    time.Time     `json:"last_viewed"`
	OkForExport   bool          `json:"ok_for_export"`
	// Could not find a case in which either of these were ever rendered,
	// even when "has_synonyms" flag is used exclusively.
	//Synonyms      []interface{} `json:"synonyms,omitempty"`
	//Parents       []interface{} `json:"parents"`
}

func parseTaxonResult(r *MushroomObserverTaxonResult) (*name_usage.NameUsageSource, error) {

	targetID, err := store.NewDataSourceTargetIDFromInt(r.ID)
	if err != nil {
		return nil, err
	}

	cn, err := store.NewCanonicalName(r.Name)
	if err != nil {
		return nil, err
	}

	return name_usage.NewNameUsageSource(store.DataSourceTypeMushroomObserver, *targetID, *cn, false)
}

