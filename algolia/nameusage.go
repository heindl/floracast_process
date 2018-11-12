package algolia

import (
	"github.com/heindl/floracast_process/nameusage/nameusage"
	"github.com/heindl/floracast_process/store"
	"github.com/heindl/floracast_process/taxa"
	"github.com/heindl/floracast_process/utils"
	"context"
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
	Ω.CommonName, err = utils.FormatTitle(Ω.CommonName)
	if err != nil {
		return err
	}

	Ω.ScientificName = utils.CapitalizeString(usage.CanonicalName().ScientificName())

	Ω.Occurrences, err = usage.Occurrences()
	if err != nil {
		return err
	}

	return nil
}
