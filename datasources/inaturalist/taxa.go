package inaturalist

import (
	"fmt"
	"context"
	"github.com/dropbox/godropbox/errors"
	"gopkg.in/tomb.v2"
	"bitbucket.org/heindl/process/utils"
	"sync"
	"time"
	"regexp"
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"bitbucket.org/heindl/process/datasources"
)

type page struct {
	TotalResults int `json:"total_results"`
	Page         int `json:"page"`
	PerPage      int `json:"per_page"`
}

type Taxon struct {
	CompleteSpeciesCount      int `json:"complete_species_count"`
	Extinct                   bool        `json:"extinct"`
	ObservationsCount         int         `json:"observations_count"`
	WikipediaURL string `json:"wikipedia_url"`
	TaxonSchemesCount         int         `json:"taxon_schemes_count"`
	Ancestry                  string      `json:"ancestry"`
	IsActive                  bool        `json:"is_active"`
	// TODO: Must have way to sync synonyms when they change.
	CurrentSynonymousTaxonIds []TaxonID `json:"current_synonymous_taxon_ids"`
	IconicTaxonID             TaxonID         `json:"iconic_taxon_id"`
	TaxonPhotos       []struct {
		Photo Photo `json:"photo"`
		Taxon Taxon `json:"taxon"`
	} `json:"taxon_photos"`
	RankLevel            RankLevel      `json:"rank_level"`
	TaxonChangesCount    int                  `json:"taxon_changes_count"`
	AtlasID              int                  `json:"atlas_id"`
	ParentID             TaxonID              `json:"parent_id"`
	Name                 string               `json:"name"`
	Rank                 string               `json:"rank"`
	ID                   TaxonID              `json:"id"`
	DefaultPhoto         Photo                `json:"default_photo"`
	AncestorIds          []TaxonID            `json:"ancestor_ids"`
	IconicTaxonName      string               `json:"iconic_taxon_name"`
	PreferredCommonName  string               `json:"preferred_common_name"`
	Ancestors            []*Taxon             `json:"ancestors"`
	Children             []*Taxon             `json:"children"`
	WikipediaSummary     string               `json:"wikipedia_summary"`
	MinSpeciesAncestry   string               `json:"min_species_ancestry"`
	CreatedAt            time.Time            `json:"created_at"`
	ConservationStatuses []ConservationStatus `json:"conservation_statuses"`
	TaxonSchemes         []*TaxonScheme
}

type ConservationStatus struct {
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

type TaxaFetcher struct {
	list []*Taxon
	sync.Mutex
	ctx context.Context
	includeChildren bool
	includeSchemes bool
}

func NewTaxaFetcher(ctx context.Context, includeChildren, includeSchemes bool) *TaxaFetcher {
	return &TaxaFetcher{
		includeChildren: includeChildren,
		includeSchemes: includeSchemes,
		list: []*Taxon{},
	}
}

func (Ω *TaxaFetcher) IndexOf(æ TaxonID) int {
	Ω.Lock()
	defer Ω.Unlock()
	for i := range Ω.list {
		if Ω.list[i].ID == æ {
			return i
		}
	}
	return -1
}

func (Ω *TaxaFetcher) Set(æ *Taxon) {

	if æ == nil {
		return
	}
	i := Ω.IndexOf(æ.ID)
	if i == -1 {
		Ω.list = append(Ω.list, æ)
	} else {
		Ω.list[i] = æ
	}
	return
}


func (Ω *TaxaFetcher) FetchTaxa(parent_taxa ...TaxonID) ([]*Taxon, error) {

	tmb := tomb.Tomb{}

	tmb.Go(func()error {
		for _, _taxonID := range parent_taxa {
			taxonID := _taxonID
			tmb.Go(func()error {
				return Ω.fetchTaxon(taxonID)
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

func (Ω *TaxaFetcher) fetchTaxon(taxonID TaxonID) error {

	if Ω.IndexOf(taxonID) != -1 {
		return nil
	}

	var response struct {
		page
		Results []*Taxon `json:"results"`
	}

	done := globalTaxonLimiter.Go()
	url := fmt.Sprintf("http://api.inaturalist.org/v1/taxa/%d", taxonID)
	if err := utils.RequestJSON(url, &response); err != nil {
		done()
		return err
	}
	done()

	if response.TotalResults == 0 {
		return errors.Newf("no taxon returned from ID: %s", taxonID)
	}

	if response.TotalResults > 1 {
		return errors.Newf("taxon request has more than one result: %s", taxonID)
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
		return errors.Newf("Sanity check failed. Have synonymous taxon ids from taxon[%s] with no way to handle.", taxonID)
	}

	//for _, _txn := range taxon.Ancestors{
	//	txn := _txn
	//	Ω.Tmb.Go(func() error {
	//		return Ω.fetchTaxon(cxt, store.INaturalistTaxonID(txn.ID))
	//	})
	//}
	return nil

}


func (Ω *TaxaFetcher) parseTaxon(txn *Taxon, isFromFullPageRequest bool) error {

	// Exit early if not a species.
	if txn.RankLevel > RankLevelSpecies {
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
		schemes, err := fetchTaxonSchemes(txn.ID)
		if err != nil {
			return err
		}
		txn.TaxonSchemes = schemes
	}

	Ω.Set(txn)
	return nil
}

type TaxonScheme struct {
	SourceType datasources.SourceType
	TargetID   datasources.TargetID
}

var taxonSchemeRegex = regexp.MustCompile(`\(([^\)]+)\)`)

func fetchTaxonSchemes(taxonID TaxonID) ([]*TaxonScheme, error) {

	url := fmt.Sprintf("http://www.inaturalist.org/taxa/%d/schemes", taxonID)

	r := bytes.NewReader([]byte{})
	if err := utils.RequestJSON(url, r); err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, errors.Wrap(err, "could parse site for goquery")
	}

	res := []*TaxonScheme{}
	doc.Find(`a[href*="/taxon_schemes/"]`).Each(func(i int, s *goquery.Selection) {
		v, _ := s.Attr("href")
		originID := datasources.SourceType(strings.TrimLeft(v, "/taxon_schemes/"))
		if string(originID) == "" {
			return
		}
		dataID := taxonSchemeRegex.FindString(s.Parent().Text())
		if dataID == "" {
			return
		}
		targetID := datasources.TargetID(strings.TrimRight(strings.TrimLeft(dataID, "("), ")"))
		res = append(res, &TaxonScheme{
			SourceType: originID,
			TargetID:   targetID,
		})
	})
	return res, nil

}

