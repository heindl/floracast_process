package inaturalist

import (
	"fmt"
	"strconv"
	"bitbucket.org/heindl/taxa/store"
	"context"
	"github.com/dropbox/godropbox/errors"
	"gopkg.in/tomb.v2"
	"bitbucket.org/heindl/taxa/utils"
	"sync"
	"time"
	"regexp"
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type page struct {
	TotalResults int `json:"total_results"`
	Page         int `json:"page"`
	PerPage      int `json:"per_page"`
}

type TaxonID int64


type INaturalistTaxon struct {
	CompleteSpeciesCount      int `json:"complete_species_count"`
	Extinct                   bool        `json:"extinct"`
	ObservationsCount         int         `json:"observations_count"`
	TaxonSchemesCount         int         `json:"taxon_schemes_count"`
	Ancestry                  string      `json:"ancestry"`
	IsActive                  bool        `json:"is_active"`
	// TODO: Must have way to sync synonyms when they change.
	CurrentSynonymousTaxonIds []TaxonID `json:"current_synonymous_taxon_ids"`
	IconicTaxonID             TaxonID         `json:"iconic_taxon_id"`
	TaxonPhotos       []struct {
		Photo INaturalistPhoto            `json:"photo"`
		Taxon INaturalistTaxon `json:"taxon"`
	} `json:"taxon_photos"`
	RankLevel           store.RankLevel                `json:"rank_level"`
	TaxonChangesCount   int                `json:"taxon_changes_count"`
	AtlasID             int                `json:"atlas_id"`
	ParentID            TaxonID                `json:"parent_id"`
	Name                string             `json:"name"`
	Rank                string             `json:"rank"`
	ID                  TaxonID              `json:"id"`
	DefaultPhoto        INaturalistPhoto              `json:"default_photo"`
	AncestorIds         []TaxonID              `json:"ancestor_ids"`
	IconicTaxonName     string             `json:"iconic_taxon_name"`
	PreferredCommonName string             `json:"preferred_common_name"`
	Ancestors           []INaturalistTaxon `json:"ancestors"`
	Children            []INaturalistTaxon `json:"children"`
	WikipediaSummary   string    `json:"wikipedia_summary"`
	MinSpeciesAncestry string    `json:"min_species_ancestry"`
	CreatedAt          time.Time `json:"created_at"`
	ConservationStatuses []ConservationStatus `json:"conservation_statuses"`
	TaxonSchemes []INaturalistTaxonScheme
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

type INaturalistPhoto struct {
	OriginalURL        string        `json:"original_url"`
	Flags              []interface{} `json:"flags"`
	Type               string        `json:"type"`
	URL                string        `json:"url"`
	SquareURL          string        `json:"square_url"`
	NativePageURL      string        `json:"native_page_url"`
	NativePhotoID      string        `json:"native_photo_id"`
	SmallURL           string        `json:"small_url"`
	Attribution        string        `json:"attribution"`
	MediumURL          string        `json:"medium_url"`
	ID                 int           `json:"id"`
	LicenseCode        string        `json:"license_code"`
	OriginalDimensions interface{}   `json:"original_dimensions"`
	LargeURL           string        `json:"large_url"`
}

func (Ω INaturalistPhoto) ToStorePhoto(taxonID store.INaturalistTaxonID, sourceID store.DataSourceID) store.Photo {
	return store.Photo{
		ID:            strconv.Itoa(Ω.ID),
		DataSourceID:  sourceID,
		TaxonID:       taxonID,
		PhotoType:     store.PhotoType(Ω.Type),
		URL:           Ω.URL,
		SquareURL:     Ω.SquareURL,
		SmallURL:      Ω.SmallURL,
		MediumURL:     Ω.MediumURL,
		LargeURL:      Ω.LargeURL,
		NativePhotoID: Ω.NativePageURL,
		Attribution:   Ω.Attribution,
		LicenseCode:   Ω.LicenseCode,
		Flags:         Ω.Flags,
	}
}

type fetchOrchestrator struct {
	Tmb tomb.Tomb
	Limiter chan struct{}
	Taxa map[store.INaturalistTaxonID]*INaturalistTaxon // Use map to avoid duplicates in recursive search.
	sync.Mutex
}

func FetchTaxaAndChildren(cxt context.Context, parent_taxa ...TaxonID) ([]*INaturalistTaxon, error) {

	orch := fetchOrchestrator{
		Tmb: tomb.Tomb{},
		Taxa: map[store.INaturalistTaxonID]*INaturalistTaxon{},
	}

	orch.Limiter = make(chan struct{}, 20)
	for i := 0; i < 20; i++ {
		orch.Limiter <- struct{}{}
	}

	orch.Tmb.Go(func() error {
		for _, t := range parent_taxa {
			orch.Tmb.Go(func() error {
				return orch.fetchTaxon(cxt, t);
			})
		}
		return nil
	})

	if err := orch.Tmb.Wait(); err != nil {
		return nil, err
	}

	res := []*INaturalistTaxon{}
	for _, v := range orch.Taxa {
		res = append(res, v)
	}
	return res, nil
}

func (Ω *fetchOrchestrator) fetchTaxon(cxt context.Context, taxonID store.INaturalistTaxonID) error {

	if !taxonID.Valid() {
		return errors.New("invalid taxon id")
	}

	// Check to see if we've already processed the full page
	if _, ok := Ω.Taxa[taxonID]; ok {
		return nil
	}

	var response struct {
		page
		Results []INaturalistTaxon `json:"results"`
	}

	url := fmt.Sprintf("http://api.inaturalist.org/v1/taxa/%d", taxonID)

	<- Ω.Limiter
	if err := utils.RequestJSON(url, &response); err != nil {
		Ω.Limiter <- struct{}{}
		return err
	}
	Ω.Limiter <- struct{}{}

	if response.TotalResults == 0 {
		return errors.Newf("no taxon returned from ID: %s", taxonID)
	}

	if response.TotalResults > 1 {
		return errors.Newf("taxon request has more than one result: %s", taxonID)
	}

	// Should only be one result.
	taxon := response.Results[0]

	for _, _txn := range taxon.Children {
		txn := _txn
		Ω.Tmb.Go(func() error {
			return Ω.fetchTaxon(cxt, store.INaturalistTaxonID(txn.ID))
		})
	}

	// Fetch synonyms? Maybe be good to have a source connection for each of them,
	// as well as a connection to the gbif.
	if len(taxon.CurrentSynonymousTaxonIds) > 0 {
		fmt.Println("Have Synonymous Taxon Ids", taxonID, taxon.CurrentSynonymousTaxonIds)
		panic("ANOMOLY DETECTED")
		for _, _txnID := range taxon.CurrentSynonymousTaxonIds {
			txnID := _txnID
			Ω.Tmb.Go(func() error {
				return Ω.fetchTaxon(cxt, store.INaturalistTaxonID(txnID))
			})
		}
	}

	//for _, _txn := range taxon.Ancestors{
	//	txn := _txn
	//	Ω.Tmb.Go(func() error {
	//		return Ω.fetchTaxon(cxt, store.INaturalistTaxonID(txn.ID))
	//	})
	//}

	rank := store.RankLevel(taxon.RankLevel)

	// Exit early if not a species.
	if rank != store.RankLevelSpecies && rank != store.RankLevelSubSpecies {
		return nil
	}

	if taxon.Extinct {
		return nil
	}

	if !taxon.IsActive {
		return nil
	}

	// Fetch Schemes
	if taxon.TaxonSchemesCount > 0 {
		var err error
		taxon.TaxonSchemes, err = Ω.fetchTaxonSchemes(taxonID)
		if err != nil {
			return err
		}
	}

	Ω.Lock()
	defer Ω.Unlock()
	Ω.Taxa[taxonID] = &taxon

	return nil
}

type INaturalistTaxonScheme struct {
	DataSourceID store.DataSourceID
	TargetID store.DataSourceTargetID
}

var taxonSchemeRegex = regexp.MustCompile(`\(([^\)]+)\)`)

func (Ω *fetchOrchestrator) fetchTaxonSchemes(taxonID store.INaturalistTaxonID) ([]INaturalistTaxonScheme, error) {

	<- Ω.Limiter
	defer func() {
		Ω.Limiter <- struct{}{}
	}()

	url := fmt.Sprintf("http://www.inaturalist.org/taxa/%d/schemes", taxonID)

	r := bytes.NewReader([]byte{})
	if err := utils.RequestJSON(url, r); err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, errors.Wrap(err, "could parse site for goquery")
	}

	res := []INaturalistTaxonScheme{}
	doc.Find(`a[href*="/taxon_schemes/"]`).Each(func(i int, s *goquery.Selection) {
		v, _ := s.Attr("href")
		originID := store.DataSourceID(strings.TrimLeft(v, "/taxon_schemes/"))
		if string(originID) == "" {
			return
		}
		dataID := taxonSchemeRegex.FindString(s.Parent().Text())
		if dataID == "" {
			return
		}
		targetID := store.DataSourceTargetID(strings.TrimRight(strings.TrimLeft(dataID, "("), ")"))
		res = append(res, INaturalistTaxonScheme {
			DataSourceID: originID,
			TargetID: targetID,
		})
	})

	return res, nil

}

