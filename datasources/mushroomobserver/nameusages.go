package mushroomobserver

import (
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/nameusage/canonicalname"
	"bitbucket.org/heindl/process/nameusage/nameusage"
	"bitbucket.org/heindl/process/utils"
	"context"
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"gopkg.in/tomb.v2"
	"net/url"
	"strings"
	"sync"
	"time"
)

type MushroomObserverQueryResult struct {
	Version         float64                        `json:"version"`
	RunDate         time.Time                      `json:"run_date"`
	Query           string                         `json:"query"`
	NumberOfRecords int                            `json:"number_of_records"`
	NumberOfPages   int                            `json:"number_of_pages"`
	PageNumber      int                            `json:"page_number"`
	Results         []*MushroomObserverTaxonResult `json:"results"`
	RunTime         float64                        `json:"run_time"`
}

func FetchNameUsages(cxt context.Context, names []string, _ datasources.TargetIDs) ([]nameusage.NameUsage, error) {

	//TODO: If names are three, consider adding var. "Cantharellus cibarius var. cibarius"
	// Only if missing in parent.

	res := []nameusage.NameUsage{}
	lock := sync.Mutex{}

	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _, _name := range names {
			releaseLmtr := fetchLmtr.Go()
			name := _name
			tmb.Go(func() error {
				defer releaseLmtr()
				usages, err := fetchNameUsages(name)
				if err != nil {
					return err
				}
				lock.Lock()
				defer lock.Unlock()
				res = append(res, usages...)
				return nil
			})
		}
		return nil
	})

	if err := tmb.Wait(); err != nil {
		return nil, err
	}

	return res, nil

}

func fetchNameUsages(name string) ([]nameusage.NameUsage, error) {
	nameURL := getMatchNameURL(name)
	var queryResult MushroomObserverQueryResult
	if err := utils.RequestJSON(nameURL, &queryResult); err != nil {
		if strings.Contains(err.Error(), "StatusCode: 503") {
			fmt.Println(fmt.Sprintf("Warning: MushroomObserver [%s] NameUsage temporarily unavailable", nameURL))
			return nil, nil
		}
		return nil, errors.Wrapf(err, "Could not fetch name from MushroomObserver [%s]", name)
	}

	// Since we are always fetching by a specific name, it will never be more than one page.
	// But just to be sure ...
	if queryResult.NumberOfPages > 1 {
		return nil, errors.Newf("Unexpected: Multiple pages returned from MushroomObserver name query [%s]", name)
	}

	if queryResult.NumberOfRecords == 0 {
		return nil, nil
	}

	res := []nameusage.NameUsage{}

	for _, result := range queryResult.Results {
		targetID, err := datasources.NewDataSourceTargetIDFromInt(datasources.TypeMushroomObserver, result.ID)
		if err != nil {
			return nil, err
		}

		cn, err := canonicalname.NewCanonicalName(result.Name, strings.ToLower(result.Rank))
		if err != nil {
			return nil, err
		}

		src, err := nameusage.NewSource(datasources.TypeMushroomObserver, targetID, cn)
		if err != nil {
			return nil, err
		}

		usage, err := nameusage.NewNameUsage(src)
		if err != nil {
			return nil, err
		}
		res = append(res, usage)
	}

	return res, nil
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
	ID            int       `json:"id"`
	Type          string    `json:"type"`
	Name          string    `json:"name"`
	Author        string    `json:"author"`
	Rank          string    `json:"rank"`
	Deprecated    bool      `json:"deprecated"`
	Misspelled    bool      `json:"misspelled"`
	Citation      string    `json:"citation"`
	Notes         string    `json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	NumberOfViews int       `json:"number_of_views"`
	LastViewed    time.Time `json:"last_viewed"`
	OkForExport   bool      `json:"ok_for_export"`
	// Could not find a case in which either of these were ever rendered,
	// even when "has_synonyms" flag is used exclusively.
	//Synonyms      []interface{} `json:"synonyms,omitempty"`
	//Parents       []interface{} `json:"parents"`
}
