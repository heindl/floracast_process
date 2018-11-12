package inaturalist

import (
	"github.com/heindl/floracast_process/utils"
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"github.com/lytics/cache"
	"gopkg.in/tomb.v2"
	"sync"
	"unsafe"
)

var globalTaxaCache *cache.Cache

var taxaCacheFetchLimiter = utils.NewLimiter(20)

func loader(txnID string) (interface{}, error) {

	var response struct {
		page
		Results []*taxon `json:"results"`
	}

	done := globalTaxonLimiter.Go()
	url := fmt.Sprintf("http://api.inaturalist.org/v1/taxa/%s", txnID)
	if err := utils.RequestJSON(url, &response); err != nil {
		done()
		return nil, err
	}
	done()

	if response.TotalResults == 0 {
		return nil, errors.Newf("no taxon returned from ID: %s", txnID)
	}

	if response.TotalResults > 1 {
		return nil, errors.Newf("taxon request has more than one result: %s", txnID)
	}

	txn := response.Results[0]

	if txn.Extinct {
		return nil, nil
	}

	//if !txn.IsActive {
	//	return nil, nil
	//}

	return txn, nil

}

func sizer(v interface{}) int64 {
	if v == nil {
		return 0
	}
	txn := v.(*taxon)
	size := unsafe.Sizeof(txn)
	return int64(size)
}

func init() {
	globalTaxaCache = cache.NewCache(16, loader, sizer)
}

// TODO: Rewrite this to merge fetching logic, and properly store and wait for values that are in progress.

func GetCachedTaxon(id taxonID) (*taxon, error) {

	v, err := globalTaxaCache.GetOrLoad(string(id.TargetID()))
	if err != nil {
		return nil, err
	}

	if v == nil {
		return nil, err
	}

	return v.(*taxon), nil

}

func GetCachedNamesAndSynonyms(ids ...taxonID) ([]string, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	tmb := tomb.Tomb{}
	res := []string{}
	locker := sync.Mutex{}
	tmb.Go(func() error {
		for _, _id := range ids {
			id := _id
			tmb.Go(func() error {
				txn, err := GetCachedTaxon(id)
				if err != nil {
					return err
				}
				synonyms, err := GetCachedNamesAndSynonyms(txn.CurrentSynonymousTaxonIds...)
				if err != nil {
					return err
				}
				locker.Lock()
				defer locker.Unlock()
				res = append(res, txn.Name)
				res = append(res, synonyms...)
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
