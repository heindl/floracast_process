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

type queryResult struct {
	NumberOfPages   int            `json:"number_of_pages"`
	NumberOfRecords int            `json:"number_of_records"`
	PageNumber      int            `json:"page_number"`
	Query           string         `json:"query"`
	Results         []*taxonResult `json:"results"`
	RunDate         time.Time      `json:"run_date"`
	RunTime         float64        `json:"run_time"`
	Version         float64        `json:"version"`
}

type taxonResult struct {
	Author        string    `json:"author"`
	Citation      string    `json:"citation"`
	CreatedAt     time.Time `json:"created_at"`
	Deprecated    bool      `json:"deprecated"`
	ID            int       `json:"id"`
	LastViewed    time.Time `json:"last_viewed"`
	Misspelled    bool      `json:"misspelled"`
	Name          string    `json:"name"`
	Notes         string    `json:"notes"`
	NumberOfViews int       `json:"number_of_views"`
	OkForExport   bool      `json:"ok_for_export"`
	Rank          string    `json:"rank"`
	Type          string    `json:"type"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Could not find a case in which either of these were ever rendered,
// even when "has_synonyms" flag is used exclusively.
//Synonyms      []interface{} `json:"synonyms,omitempty"`
//Parents       []interface{} `json:"parents"`

// FetchNameUsages impliments the NameUsage fetch interface.
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
	qResult := queryResult{}
	if err := utils.RequestJSON(nameURL, &qResult); err != nil {
		if strings.Contains(err.Error(), "StatusCode: 503") {
			fmt.Println(fmt.Sprintf("Warning: MushroomObserver [%s] NameUsage temporarily unavailable", nameURL))
			return nil, nil
		}
		return nil, errors.Wrapf(err, "Could not fetch name from MushroomObserver [%s]", name)
	}

	// Since we are always fetching by a specific name, it will never be more than one page.
	// But just to be sure ...
	if qResult.NumberOfPages > 1 {
		return nil, errors.Newf("Unexpected: Multiple pages returned from MushroomObserver name query [%s]", name)
	}

	if qResult.NumberOfRecords == 0 {
		return nil, nil
	}

	res := []nameusage.NameUsage{}

	for _, result := range qResult.Results {
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
