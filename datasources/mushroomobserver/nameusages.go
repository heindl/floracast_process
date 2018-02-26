package mushroomobserver

import (
	"fmt"
	"strings"
	"time"
	"net/url"
	"bitbucket.org/heindl/process/utils"
	"github.com/dropbox/godropbox/errors"
	"gopkg.in/tomb.v2"
	"sync"
	"context"
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/nameusage/canonicalname"
	"bitbucket.org/heindl/process/nameusage/nameusage"
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

var lmtr = utils.NewLimiter(10)

func FetchNameUsages(cxt context.Context, names []string, sources datasources.TargetIDs) ([]nameusage.NameUsage, error) {

	//TODO: If names are three, consider adding var. "Cantharellus cibarius var. cibarius"
	// Only if missing in parent.

	usageResponse := []nameusage.NameUsage{}
	locker := sync.Mutex{}

	tmb := tomb.Tomb{}
	tmb.Go(func() error{
		for _, _name := range names {
			name := _name
			tmb.Go(func() error {
				release := lmtr.Go()
				defer release()

				nameURL := getMatchNameURL(name)
				var queryResult MushroomObserverQueryResult
				if err := utils.RequestJSON(nameURL, &queryResult); err != nil {
					if strings.Contains(err.Error(), "StatusCode: 503") {
						fmt.Println(fmt.Sprintf("Warning: Could not fetch name from MushroomObserver [%s]", nameURL))
						return nil
					}
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

				for _, _result := range queryResult.Results {
					result := _result
					tmb.Go(func() error {
						targetID, err := datasources.NewDataSourceTargetIDFromInt(result.ID)
						if err != nil {
							return err
						}

						cn, err := canonicalname.NewCanonicalName(result.Name, strings.ToLower(result.Rank))
						if err != nil {
							return err
						}

						src, err := nameusage.NewSource(datasources.TypeMushroomObserver, targetID, cn)
						if err != nil {
							return err
						}

						usage, err := nameusage.NewNameUsage(src)
						if err != nil {
							return err
						}

						locker.Lock()
						defer locker.Unlock()
						usageResponse = append(usageResponse, usage)
						return nil
					})
				}
				return nil
			})
		}
		return nil
	})

	if err := tmb.Wait(); err != nil {
		return nil, err
	}

	return usageResponse, nil

}

func getMatchNameURL(name string) string {
	parameters := strings.Join([]string{
		//fmt.Sprintf("updated_at=%s-%s", time.Now().AddUsage(time.Hour * 24 * 30 * -120).Format("20060102"), time.Now().Format("20060102")),
		"format=json",
		"is_deprecated=false",
		//"ok_for_export=true",
		//"has_classification=true",
		fmt.Sprintf("name=%s", url.QueryEscape(name)),
		//"has_synonyms=true",
		"detail=high",
		//fmt.Sprintf("rank=%s", rank),
	}, "&")
	return "http://mushroomobserver.org/api/names?" + parameters
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