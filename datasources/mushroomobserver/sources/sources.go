package main

import (
	"fmt"
	"strings"
	//"bitbucket.org/heindl/processors/store"
	//"golang.org/x/net/context"
	"time"
)


type MushroomObserverQueryResult struct {
	Version         float64   `json:"version"`
	RunDate         time.Time `json:"run_date"`
	Query           string    `json:"query"`
	NumberOfRecords int       `json:"number_of_records"`
	NumberOfPages   int       `json:"number_of_pages"`
	PageNumber      int       `json:"page_number"`
	Results         []MushroomObserverTaxonResult `json:"results"`
	RunTime float64 `json:"run_time"`
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
	Synonyms      []interface{} `json:"synonyms,omitempty"`
	Parents       []interface{} `json:"parents"`
}

func main() {

	// has_classification=true
	// is_deprecated=false
	// ok_for_export=true
	// updated_at=

	// Sort data sources by kind and last updated date.
	//taxa_store, err := store.NewTaxaStore()
	//if err != nil {
	//	panic(err)
	//}
	//
	//lastCreated, err := taxa_store.GetSourceLastCreated(context.Background(), store.DataSourceKindOccurrence, store.DataSourceIDMushroomObserver)
	//if err != nil {
	//	panic(err)
	//}

	parameters := []string{
		"updated_at=%s-%s",
		"format=json",
		"is_deprecated=false",
		//"ok_for_export=true",
		//"has_classification=true",
		"classification_has=morchella",
		//"has_synonyms=true",
		"detail=high",
		"rank=%s",
	}

	startDateStr := time.Now().Add(time.Hour * 24 * 30 * -120).Format("20060102")
	endDateStr := time.Now().Format("20060102")

	// Use the last updated date as the starting point for the query.
	// Not sure how to enumerate multiple so will need a separate query for each Rank.
	for _, rank := range []string{
		//"Subgenus",
		//"Section",
		//"Subsection",
		//"Series",
		"Species",
		"Subspecies",
		"Stirps",
		"Variety",
		} {
		path := "http://mushroomobserver.org/api/names?" + fmt.Sprintf(strings.Join(parameters, "&"), startDateStr, endDateStr, rank)
		fmt.Println(path)
	}

	// Collect all elements into a list.
	// If no existing data source exists with DataSource, try to match to existing taxon or create a new one.

	// If no match, create one.



}
