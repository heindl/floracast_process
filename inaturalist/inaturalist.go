package inaturalist
//
//import (
//	"bitbucket.org/heindl/taxa/store"
//	"bitbucket.org/heindl/taxa/utils"
//	"bytes"
//	"context"
//	"flag"
//	"fmt"
//	"github.com/PuerkitoBio/goquery"
//	"github.com/heindl/gbif"
//	"github.com/jonboulle/clockwork"
//	"github.com/saleswise/errors/errors"
//	. "github.com/saleswise/malias"
//	"gopkg.in/tomb.v2"
//	"regexp"
//	"strconv"
//	"strings"
//	"sync"
//	"time"
//)
//
//func main() {
//
//	taxa := flag.String("taxa", "", "parent taxa for query, string separated")
//	flag.Parse()
//
//	if *taxa == "" {
//		return
//	}
//
//	toProcess := []store.INaturalistTaxonID{}
//	for _, s := range strings.Split(*taxa, ",") {
//		toProcess = append(toProcess, store.INaturalistTaxonID(s))
//	}
//
//	ts, err := store.NewTaxaStore()
//	if err != nil {
//		panic(err)
//	}
//	f := fetcher{
//		Store: ts,
//		Clock: clockwork.NewRealClock(),
//	}
//
//	if err := f.FetchProcessTaxa(context.Background(), toProcess); err != nil {
//		panic(err)
//	}
//}
//
//type fetcher struct {
//	// Reference to google data store.
//	Store         store.TaxaStore
//	ProcessedTaxa []string
//	//DataSources store.DataSources
//	//Photos store.Photos
//	sync.Mutex
//	Clock clockwork.Clock
//	Tomb  *tomb.Tomb
//}
//
//func (Ω *fetcher) FetchProcessTaxa(cxt context.Context, parent_taxa []store.INaturalistTaxonID) error {
//	for _, t := range parent_taxa {
//		if err := Ω._fetchProcessTaxa(cxt, t); err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//func (Ω *fetcher) _fetchProcessTaxa(cxt context.Context, parent_taxon store.INaturalistTaxonID) error {
//
//	if !parent_taxon.Valid() {
//		return errors.New("invalid taxa")
//	}
//
//	Ω.Tomb = &tomb.Tomb{}
//	Ω.Tomb.Go(func() error {
//
//		type TaxonCount struct {
//			Count int   `json:"count"`
//			Taxon Taxon `json:"taxon"`
//		}
//
//		var response struct {
//			Counter
//			Results []*TaxonCount `json:"results"`
//		}
//
//		// Right now pagination - ?page=2 - appears to not work, so having to jack up it up to 10,000. Be sure to check below that the
//		// total returned is less than the s.
//		url := fmt.Sprintf("http://api.inaturalist.org/v1/observations/species_counts?taxon_id=%s&per_page=10000", string(parent_taxon))
//
//		if err := utils.RequestJSON(url, &response); err != nil {
//			return err
//		}
//
//		if response.TotalResults > response.PerPage {
//			return errors.New("The total number of taxa [%d] is higher than the number I could fit on the page [%d]. This is problematic becuase pagination, to my knowledge, does not work in the API now.")
//		}
//
//		limiter := make(chan struct{}, 20)
//		for i := 0; i < 20; i++ {
//			limiter <- struct{}{}
//		}
//
//		for _, _taxon := range response.Results {
//			taxon := _taxon
//			// No opt if there are no observations. Likely irrelevant to the project for now because it is to rare or difficult to find.
//			//if taxon.INaturalistTaxon.ObservationsCount == 0 {
//			//	fmt.Printf("Zero observations for taxon %s [%d].", taxon.INaturalistTaxon.Name, taxon.INaturalistTaxon.ID)
//			//	continue
//			//}
//			<-limiter
//			Ω.Tomb.Go(func() error {
//				defer func() {
//					limiter <- struct{}{}
//				}()
//				if err := Ω.fetchProcessTaxon(cxt, store.INaturalistTaxonID(strconv.Itoa(int(taxon.Taxon.ID)))); err != nil {
//					return err
//				}
//				return nil
//			})
//		}
//
//		return nil
//	})
//	return Ω.Tomb.Wait()
//
//}
//
//type Counter struct {
//	TotalResults int `json:"total_results"`
//	Page         int `json:"page"`
//	PerPage      int `json:"per_page"`
//}
//
//type Taxon struct {
//	ObservationsCount int    `json:"observations_count"`
//	TaxonSchemesCount int    `json:"taxon_schemes_count"`
//	IsActive          bool   `json:"is_active"`
//	Ancestry          string `json:"ancestry"`
//	IconicTaxonID     int    `json:"iconic_taxon_id"`
//	TaxonPhotos       []struct {
//		Photo Photo `json:"photo"`
//		Taxon Taxon `json:"taxon"`
//	} `json:"taxon_photos"`
//	RankLevel           int     `json:"rank_level"`
//	TaxonChangesCount   int     `json:"taxon_changes_count"`
//	AtlasID             int     `json:"atlas_id"`
//	ParentID            int     `json:"parent_id"`
//	Name                string  `json:"name"`
//	Rank                string  `json:"rank"`
//	ID                  int64   `json:"id"`
//	DefaultPhoto        Photo   `json:"default_photo"`
//	AncestorIds         []int   `json:"ancestor_ids"`
//	IconicTaxonName     string  `json:"iconic_taxon_name"`
//	PreferredCommonName string  `json:"preferred_common_name"`
//	Ancestors           []Taxon `json:"ancestors"`
//	Children            []Taxon `json:"children"`
//	ListedTaxa          []struct {
//		ID                 int    `json:"id"`
//		TaxonID            int    `json:"taxon_id"`
//		EstablishmentMeans string `json:"establishment_means"`
//		Place              struct {
//			ID               int    `json:"id"`
//			Name             string `json:"name"`
//			DisplayName      string `json:"display_name"`
//			AdminLevel       int    `json:"admin_level"`
//			AncestorPlaceIds []int  `json:"ancestor_place_ids"`
//		} `json:"place"`
//		List struct {
//			ID    int    `json:"id"`
//			Title string `json:"title"`
//		} `json:"list"`
//	} `json:"listed_taxa"`
//	WikipediaSummary   string    `json:"wikipedia_summary"`
//	MinSpeciesAncestry string    `json:"min_species_ancestry"`
//	CreatedAt          time.Time `json:"created_at"`
//}
//
//type Photo struct {
//	Flags              []interface{} `json:"flags"`
//	Type               string        `json:"type"`
//	URL                string        `json:"url"`
//	SquareURL          string        `json:"square_url"`
//	NativePageURL      string        `json:"native_page_url"`
//	NativePhotoID      string        `json:"native_photo_id"`
//	SmallURL           string        `json:"small_url"`
//	Attribution        string        `json:"attribution"`
//	MediumURL          string        `json:"medium_url"`
//	ID                 int           `json:"id"`
//	LicenseCode        string        `json:"license_code"`
//	OriginalDimensions interface{}   `json:"original_dimensions"`
//	LargeURL           string        `json:"large_url"`
//}
//
//func (Ω Photo) Format(taxonID store.INaturalistTaxonID, sourceID store.DataSourceID) store.Photo {
//	return store.Photo{
//		ID:            strconv.Itoa(Ω.ID),
//		DataSourceID:  sourceID,
//		TaxonID:       taxonID,
//		PhotoType:     store.PhotoType(Ω.Type),
//		URL:           Ω.URL,
//		SquareURL:     Ω.SquareURL,
//		SmallURL:      Ω.SmallURL,
//		MediumURL:     Ω.MediumURL,
//		LargeURL:      Ω.LargeURL,
//		NativePhotoID: Ω.NativePageURL,
//		Attribution:   Ω.Attribution,
//		LicenseCode:   Ω.LicenseCode,
//		Flags:         Ω.Flags,
//	}
//}
//
//func (Ω *fetcher) fetchProcessTaxon(cxt context.Context, taxonID store.INaturalistTaxonID) error {
//
//	if !taxonID.Valid() {
//		return errors.New("invalid taxon id")
//	}
//
//	// Check to see if we've already processed it. Having states likely means we've grabbed the full page already.
//	//if i := Ω.Taxa.Index(taxonID); i != -1 && len(Ω.Taxa[i].States) > 0 {
//	//	return nil
//	//}
//
//	var response struct {
//		Counter
//		Results []Taxon `json:"results"`
//	}
//
//	url := fmt.Sprintf("http://api.inaturalist.org/v1/taxa/%s", string(taxonID))
//
//	if err := utils.RequestJSON(url, &response); err != nil {
//		return err
//	}
//
//	if response.TotalResults == 0 {
//		return errors.New("no taxon returned from ").SetState(M{utils.LogkeyIdentifier: taxonID})
//	}
//
//	if response.TotalResults > 1 {
//		return errors.New("taxon request has more than one result").SetState(M{utils.LogkeyDatastoreKey: taxonID})
//	}
//
//	// Should only be one result.
//	taxa := response.Results[0]
//
//	Ω.Tomb.Go(func() (err error) {
//		var lastAncestor store.INaturalistTaxonID
//		for i, _a := range taxa.Ancestors {
//			a := _a
//			la, err := Ω.processTaxon(cxt, a, lastAncestor, (i != 0))
//			if err != nil {
//				return err
//			}
//			lastAncestor = la
//		}
//		if _, err := Ω.processTaxon(cxt, taxa, lastAncestor, true); err != nil {
//			return err
//		}
//		return nil
//	})
//
//	Ω.Tomb.Go(func() error {
//		for _, _a := range taxa.Children {
//			a := _a
//			rank := store.RankLevel(a.RankLevel)
//			if rank == store.RankLevelSpecies || rank == store.RankLevelSubSpecies || rank == store.RankLevelVariety {
//				Ω.Tomb.Go(func() error {
//					if err := Ω.fetchProcessTaxon(cxt, store.INaturalistTaxonID(strconv.Itoa(int(a.ID)))); err != nil {
//						return err
//					}
//					return nil
//				})
//			}
//		}
//		return nil
//	})
//
//	return nil
//}
//
//func (Ω *fetcher) processTaxon(cxt context.Context, txn Taxon, parent store.INaturalistTaxonID, shouldHaveParent bool) (store.INaturalistTaxonID, error) {
//
//	// No need to reprocess a parent key we've already processed it before. Noopt and return the key.
//	//taxonKey := Ω.Fetched.Find(txn.ID, store.EntityKindTaxon)
//	//if taxonKey != nil {
//	//	return taxonKey, nil
//	//}
//
//	taxonID := store.INaturalistTaxonID(strconv.Itoa(int(txn.ID)))
//
//	if utils.ContainsString(Ω.ProcessedTaxa, string(taxonID)) {
//		return taxonID, nil
//	}
//
//	Ω.Lock()
//	Ω.ProcessedTaxa = append(Ω.ProcessedTaxa, string(taxonID))
//	Ω.Unlock()
//
//
//	rank, ok := store.TaxonRankMap[txn.Rank]
//	if !ok {
//		return store.INaturalistTaxonID(""), errors.Newf("unsupported rank: %s", txn.Rank)
//	}
//
//	if !parent.Valid() && shouldHaveParent {
//		return store.INaturalistTaxonID(""), errors.New("parent taxon id expected but invalid")
//	}
//
//	t := store.Taxon{
//		ParentID:         parent,
//		CanonicalName:    store.CanonicalName(txn.Name),
//		Rank:             rank,
//		ID:               taxonID,
//		RankLevel:        store.RankLevel(txn.RankLevel),
//		CommonName:       txn.PreferredCommonName,
//		CreatedAt:        Ω.Clock.Now(),
//		ModifiedAt:       Ω.Clock.Now(),
//		WikipediaSummary: txn.WikipediaSummary,
//	}
//
//	for _, lt := range txn.ListedTaxa {
//		if lt.Place.AdminLevel == 1 {
//			t.States = append(t.States, store.State{
//				EstablishmentMeans: lt.EstablishmentMeans,
//				Name:               lt.Place.Name,
//			})
//		}
//	}
//
//	if err := Ω.Store.CreateTaxonIfNotExists(cxt, t); err != nil {
//		return store.INaturalistTaxonID(""), err
//	}
//
//	taxonPhoto := ""
//	if txn.DefaultPhoto.ID != 0 {
//		if err := Ω.Store.SetPhoto(cxt, txn.DefaultPhoto.Format(t.ID, store.DataSourceIDINaturalist)); err != nil {
//			return store.INaturalistTaxonID(""), err
//		}
//		taxonPhoto = txn.DefaultPhoto.MediumURL
//	}
//
//	// Note that the photos store sub-species, and so far the only place i can find them.
//	for _, p := range txn.TaxonPhotos {
//		if strconv.Itoa(int(p.Taxon.ID)) == string(t.ID) {
//			if err := Ω.Store.SetPhoto(cxt, p.Photo.Format(t.ID, store.DataSourceIDINaturalist)); err != nil {
//				return store.INaturalistTaxonID(""), err
//			}
//			if taxonPhoto == "" {
//				taxonPhoto = p.Photo.MediumURL
//			}
//		}
//	}
//
//	if err := Ω.Store.SetTaxonPhoto(cxt, t.ID, taxonPhoto); err != nil {
//		return store.INaturalistTaxonID(""), err
//	}
//
//	dataSources, err := Ω.fetchDataSources(t.ID, t.CanonicalName, (t.RankLevel == store.RankLevelSubSpecies || t.RankLevel == store.RankLevelSpecies || t.RankLevel == store.RankLevelVariety))
//	if err != nil {
//		return store.INaturalistTaxonID(""), err
//	}
//	if len(dataSources) > 0 {
//		for _, s := range dataSources {
//			if err := Ω.Store.UpsertDataSource(cxt, s); err != nil {
//				return store.INaturalistTaxonID(""), err
//			}
//		}
//	}
//
//	return t.ID, nil
//}
//
//var schemeRegex = regexp.MustCompile(`\(([^\)]+)\)`)
//
//var schemeFetchRateLimiter = time.Tick(time.Second / 2)
//
//func (Ω *fetcher) fetchDataSources(taxonID store.INaturalistTaxonID, canonicalName store.CanonicalName, isSpecies bool) ([]store.DataSource, error) {
//
//	if !taxonID.Valid() {
//		return nil, errors.New("invalid taxon id")
//	}
//
//	<-schemeFetchRateLimiter
//
//	url := fmt.Sprintf("http://www.inaturalist.org/taxa/%s/schemes", string(taxonID))
//
//	r := bytes.NewReader([]byte{})
//	if err := utils.RequestJSON(url, r); err != nil {
//		return nil, err
//	}
//
//	doc, err := goquery.NewDocumentFromReader(r)
//	if err != nil {
//		return nil, errors.Wrap(err, "could parse site for goquery").SetState(M{utils.LogkeyURL: url})
//	}
//
//	res := []store.DataSource{}
//	// Find the review items
//	pairs := []struct {
//		OriginID store.DataSourceID
//		TargetID store.DataSourceTargetID
//	}{}
//	doc.Find(`a[href*="/taxon_schemes/"]`).Each(func(i int, s *goquery.Selection) {
//		v, _ := s.Attr("href")
//		originID := store.DataSourceID(strings.TrimLeft(v, "/taxon_schemes/"))
//		if string(originID) == "" {
//			return
//		}
//		dataID := schemeRegex.FindString(s.Parent().Text())
//		if dataID == "" {
//			return
//		}
//		targetID := store.DataSourceTargetID(strings.TrimRight(strings.TrimLeft(dataID, "("), ")"))
//		pairs = append(pairs, struct {
//			OriginID store.DataSourceID
//			TargetID store.DataSourceTargetID
//		}{originID, targetID})
//	})
//	for _, pair := range pairs {
//		if pair.OriginID == store.DataSourceIDGBIF && isSpecies {
//			res = append(res, store.DataSource{
//				Kind:     store.DataSourceKindOccurrence,
//				SourceID: store.DataSourceIDGBIF,
//				TargetID: pair.TargetID,
//				TaxonID:  taxonID,
//			})
//			res = append(res, store.DataSource{
//				Kind:     store.DataSourceKindDescription,
//				SourceID: store.DataSourceIDGBIF,
//				TargetID: pair.TargetID,
//				TaxonID:  taxonID,
//			})
//			res = append(res, store.DataSource{
//				Kind:     store.DataSourceKindPhoto,
//				SourceID: store.DataSourceIDGBIF,
//				TargetID: pair.TargetID,
//				TaxonID:  taxonID,
//			})
//
//			additionalIDs, err := Ω.fetchAdditionalGBIFTaxonIDs(string(canonicalName), pair.TargetID)
//			if err != nil {
//				return nil, err
//			}
//
//			for _, a := range additionalIDs {
//				res = append(res, store.DataSource{
//					Kind:     store.DataSourceKindOccurrence,
//					SourceID: store.DataSourceIDGBIF,
//					TargetID: a,
//					TaxonID:  taxonID,
//				})
//				res = append(res, store.DataSource{
//					Kind:     store.DataSourceKindDescription,
//					SourceID: store.DataSourceIDGBIF,
//					TargetID: a,
//					TaxonID:  taxonID,
//				})
//				res = append(res, store.DataSource{
//					Kind:     store.DataSourceKindPhoto,
//					SourceID: store.DataSourceIDGBIF,
//					TargetID: a,
//					TaxonID:  taxonID,
//				})
//			}
//		}
//	}
//
//	return res, nil
//
//}
//
//func (Ω *fetcher) fetchAdditionalGBIFTaxonIDs(iNaturalistCanonicalName string, gbifTaxonID store.DataSourceTargetID) ([]store.DataSourceTargetID, error) {
//	names, err := gbif.Search(gbif.SearchQuery{Q: iNaturalistCanonicalName})
//	if err != nil {
//		return nil, errors.Wrap(err, "could not get name usage")
//	}
//	nubs := []int{}
//	for _, n := range names {
//		if n.CanonicalName == iNaturalistCanonicalName &&
//			n.NubKey != 0 &&
//			(n.TaxonomicStatus == gbif.TaxonomicStatusACCEPTED || n.TaxonomicStatus == gbif.TaxonomicStatusSYNONYM) {
//			nubs = utils.AddIntToSet(nubs, n.NubKey)
//		}
//	}
//
//	res := []store.DataSourceTargetID{}
//	for _, n := range nubs {
//		t := store.DataSourceTargetID(strconv.Itoa(n))
//		if t == gbifTaxonID {
//			continue
//		}
//		res = append(res, t)
//	}
//	return res, nil
//}