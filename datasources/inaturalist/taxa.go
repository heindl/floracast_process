package inaturalist

import (
	"bitbucket.org/heindl/process/utils"
	"context"
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"gopkg.in/tomb.v2"
	"sync"
	"time"
)

type page struct {
	TotalResults int `json:"total_results"`
	Page         int `json:"page"`
	PerPage      int `json:"per_page"`
}

type taxon struct {
	CompleteSpeciesCount int    `json:"complete_species_count"`
	Extinct              bool   `json:"extinct"`
	ObservationsCount    int    `json:"observations_count"`
	WikipediaURL         string `json:"wikipedia_url"`
	TaxonSchemesCount    int    `json:"taxon_schemes_count"`
	Ancestry             string `json:"ancestry"`
	IsActive             bool   `json:"is_active"`
	// TODO: Must have way to sync synonyms when they change.
	CurrentSynonymousTaxonIds []taxonID `json:"current_synonymous_taxon_ids"`
	IconicTaxonID             taxonID   `json:"iconic_taxon_id"`
	TaxonPhotos               []struct {
		Photo photo `json:"photo"`
		Taxon taxon `json:"taxon"`
	} `json:"taxon_photos"`
	RankLevel            rankLevel            `json:"rank_level"`
	TaxonChangesCount    int                  `json:"taxon_changes_count"`
	AtlasID              int                  `json:"atlas_id"`
	ParentID             taxonID              `json:"parent_id"`
	Name                 string               `json:"name"`
	Rank                 string               `json:"rank"`
	ID                   taxonID              `json:"id"`
	DefaultPhoto         photo                `json:"default_photo"`
	AncestorIds          []taxonID            `json:"ancestor_ids"`
	IconicTaxonName      string               `json:"iconic_taxon_name"`
	PreferredCommonName  string               `json:"preferred_common_name"`
	Ancestors            []*taxon             `json:"ancestors"`
	Children             []*taxon             `json:"children"`
	WikipediaSummary     string               `json:"wikipedia_summary"`
	MinSpeciesAncestry   string               `json:"min_species_ancestry"`
	CreatedAt            time.Time            `json:"created_at"`
	ConservationStatuses []conservationStatus `json:"conservation_statuses"`
	TaxonSchemes         []*taxonScheme
}

type conservationStatus struct {
	PlaceID    int    `json:"place_id"`
	SourceID   int    `json:"source_id"`
	Authority  string `json:"authority"`
	Status     string `json:"status"`
	Iucn       int    `json:"iucn"`
	Geoprivacy string `json:"geoprivacy"`
	Place      struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		DisplayName string `json:"display_name"`
	} `json:"place"`
}

type taxaFetcher struct {
	list []*taxon
	sync.Mutex
	includeChildren bool
	includeSchemes  bool
}

func newTaxaFetcher(ctx context.Context, includeChildren, includeSchemes bool) *taxaFetcher {
	return &taxaFetcher{
		includeChildren: includeChildren,
		includeSchemes:  includeSchemes,
		list:            []*taxon{},
	}
}

func (Ω *taxaFetcher) IndexOf(æ taxonID) int {
	Ω.Lock()
	defer Ω.Unlock()
	for i := range Ω.list {
		if Ω.list[i].ID == æ {
			return i
		}
	}
	return -1
}

func (Ω *taxaFetcher) Set(æ *taxon) {
	if æ == nil {
		return
	}
	i := Ω.IndexOf(æ.ID)
	Ω.Lock()
	defer Ω.Unlock()
	if i == -1 {
		Ω.list = append(Ω.list, æ)
	} else {
		Ω.list[i] = æ
	}
}

func (Ω *taxaFetcher) FetchTaxa(parentTaxa ...taxonID) ([]*taxon, error) {

	tmb := tomb.Tomb{}

	tmb.Go(func() error {
		for _, _txnID := range parentTaxa {
			txnID := _txnID
			tmb.Go(func() error {
				return Ω.fetchTaxon(txnID)
			})
		}
		return nil
	})

	if err := tmb.Wait(); err != nil {
		return nil, err
	}

	return Ω.list, nil
}

var globalTaxonLimiter = utils.NewLimiter(20)

func (Ω *taxaFetcher) fetchTaxon(txnID taxonID) error {

	if Ω.IndexOf(txnID) != -1 {
		return nil
	}

	var response struct {
		page
		Results []*taxon `json:"results"`
	}

	done := globalTaxonLimiter.Go()
	url := fmt.Sprintf("http://api.inaturalist.org/v1/taxa/%d", txnID)
	if err := utils.RequestJSON(url, &response); err != nil {
		done()
		return err
	}
	done()

	if response.TotalResults == 0 {
		return errors.Newf("no taxon returned from ID: %s", txnID)
	}

	if response.TotalResults > 1 {
		return errors.Newf("taxon request has more than one result: %s", txnID)
	}

	if Ω.includeChildren {
		tmb := tomb.Tomb{}
		tmb.Go(func() error {
			for _, _txn := range response.Results[0].Children {
				txn := _txn
				tmb.Go(func() error {
					return Ω.parseTaxon(txn, false)
				})
			}
			return nil
		})
		if err := tmb.Wait(); err != nil {
			return err
		}
	}

	if err := Ω.parseTaxon(response.Results[0], true); err != nil {
		return err
	}

	// Fetch synonyms? Maybe be good to have a source connection for each of them,
	// as well as a connection to the gbif.
	if len(response.Results[0].CurrentSynonymousTaxonIds) > 0 {
		return errors.Newf("Sanity check failed. Have synonymous taxon ids from taxon[%s] with no way to handle.", txnID)
	}

	//for _, _txn := range taxon.Ancestors{
	//	txn := _txn
	//	Ω.Tmb.Go(func() error {
	//		return Ω.fetchTaxon(cxt, store.INaturalistTaxonID(txn.ID))
	//	})
	//}
	return nil

}

func (Ω *taxaFetcher) parseTaxon(txn *taxon, isFromFullPageRequest bool) error {

	// Exit early if not a species.
	if txn.RankLevel > rankLevelSpecies {
		// Fetch children if this was the child of another request. Otherwise we're safe stopping with species.
		if !isFromFullPageRequest {
			return Ω.fetchTaxon(txn.ID)
		}
		// We expect children to be parsed in the caller of this function.
		return nil
	}

	if txn.Extinct {
		return nil
	}

	if !txn.IsActive {
		return nil
	}

	// Fetch Schemes
	if txn.TaxonSchemesCount > 0 && Ω.includeSchemes {
		schemes, err := txn.ID.fetchTaxonSchemes()
		if err != nil {
			return err
		}
		txn.TaxonSchemes = schemes
	}

	Ω.Set(txn)
	return nil
}
