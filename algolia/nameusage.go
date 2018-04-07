package algolia

import (
	"context"
	"strings"

	"bitbucket.org/heindl/process/nameusage/nameusage"
	"bitbucket.org/heindl/process/store"
	"bitbucket.org/heindl/process/taxa"
	"bitbucket.org/heindl/process/utils"
	"gopkg.in/tomb.v2"
	"sync"
)

type StandardNameUsageRecord struct {
	NameUsageID    nameusage.ID `json:""`
	ScientificName string       `json:""`
	CommonName     string       `json:""`
	Thumbnail      string       `json:""`
	Occurrences    int          `json:""`
	Predictions    int          `json:""`
	floraStore     store.FloraStore
	sync.Mutex
}

func (Ω *StandardNameUsageRecord) hydrate() error {
	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		tmb.Go(Ω.fetchUsageData)
		tmb.Go(Ω.fetchThumbnail)
		return nil
	})
	return tmb.Wait()
}

//func nameUsageIDFilter(nameUsageIDs ...nameusage.ID) algoliasearch.Map {
//	if len(nameUsageIDs) == 0 {
//		return nil
//	}
//	filterKeys := []string{}
//	for _, id := range nameUsageIDs {
//		filterKeys = append(filterKeys, fmt.Sprintf("%s:%s", keyNameUsageID, id))
//	}
//	return algoliasearch.Map{
//		"filters": strings.Join(filterKeys, " OR "),
//	}
//}

func (Ω *StandardNameUsageRecord) fetchThumbnail() error {
	txn, err := taxa.Fetch(context.Background(), Ω.floraStore, Ω.NameUsageID)
	if err != nil {
		return err
	}
	// TODO: Generate thumbnail from image
	if txn.Photo != nil {
		Ω.Thumbnail = txn.Photo.Thumbnail
	}
	return nil
}

func (Ω *StandardNameUsageRecord) fetchUsageData() error {

	usage, err := nameusage.Fetch(context.Background(), Ω.floraStore, Ω.NameUsageID)
	if err != nil {
		return err
	}

	Ω.CommonName, err = usage.CommonName()
	if err != nil {
		return err
	}
	Ω.CommonName = strings.Title(Ω.CommonName)

	Ω.ScientificName = utils.CapitalizeString(usage.CanonicalName().ScientificName())

	Ω.Occurrences, err = usage.Occurrences()
	if err != nil {
		return err
	}

	return nil
}
