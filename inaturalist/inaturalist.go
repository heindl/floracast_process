package inaturalist

import (
	"fmt"
	"strings"
	"io/ioutil"
	"net/http"
	"github.com/saleswise/errors/errors"
	"encoding/json"
	"gopkg.in/tomb.v2"
	"time"
	"bitbucket.org/heindl/utils"
	. "github.com/saleswise/malias"
	"bitbucket.org/heindl/species/store"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strconv"
	"sync"
	"cloud.google.com/go/datastore"
	"github.com/jonboulle/clockwork"
)


type fetcher struct {
	// Reference to google data store.
	Store       store.TaxaStore
	Taxa store.Taxa
	Schema store.Schema
	Photos store.Photos
	sync.Mutex
	Clock       clockwork.Clock
	Tomb *tomb.Tomb
}

func (Ω *fetcher) FetchProcessTaxa(parent_taxa int) error {

	if parent_taxa == 0 {
		return errors.New("invalid taxa")
	}

	Ω.Tomb = &tomb.Tomb{}
	Ω.Tomb.Go(func() error {

		type TaxonCount struct {
			Count int `json:"count"`
			Taxon Taxon `json:"taxon"`
		}

		var response struct {
			Counter
			Results []*TaxonCount `json:"results"`
		}

		// Right now pagination - ?page=2 - appears to not work, so having to jack up it up to 10,000. Be sure to check below that the
		// total returned is less than the s.
		url := fmt.Sprintf("http://api.inaturalist.org/v1/observations/species_counts?taxon_id=%d&per_page=10000", parent_taxa)
		if err := request(url, &response); err != nil {
			return err
		}

		if response.TotalResults > response.PerPage {
			return errors.New("The total number of taxa [%d] is higher than the number I could fit on the page [%d]. This is problematic becuase pagination, to my knowledge, does not work in the API now.")
		}

		limiter := make(chan struct{}, 20)
		for i := 0; i < 20; i++ {
			limiter <- struct{}{}
		}

		for _, _taxon := range response.Results {
			taxon := _taxon
			// No opt if there are no observations. Likely irrelevant to the project for now because it is to rare or difficult to find.
			if taxon.Taxon.ObservationsCount == 0 {
				fmt.Printf("Zero observations for taxon %s [%d].", taxon.Taxon.Name, taxon.Taxon.ID)
				continue
			}
			<-limiter
			Ω.Tomb.Go(func() error {
				defer func() {
					limiter <- struct{}{}
				}()
				if err := Ω.fetchProcessTaxon(store.TaxonID(taxon.Taxon.ID)); err != nil {
					return err
				}
				return nil
			})
		}

		return nil
	})
	if err := Ω.Tomb.Wait(); err != nil {
		return err
	}

	fmt.Println("TAXA", len(Ω.Taxa))

	if err := Ω.Store.SetTaxa(Ω.Taxa); err != nil {
		return err
	}

	fmt.Println("SCHEMA", len(Ω.Schema))

	if err := Ω.Store.SetSchema(Ω.Schema); err != nil {
		return err
	}

	fmt.Println("PHOTOS", len(Ω.Photos))

	return Ω.Store.SetPhotos(Ω.Photos)

}


type Counter struct {
	TotalResults int `json:"total_results"`
	Page int `json:"page"`
	PerPage int `json:"per_page"`
}

type Taxon struct {
	ObservationsCount int `json:"observations_count"`
	TaxonSchemesCount int `json:"taxon_schemes_count"`
	IsActive bool `json:"is_active"`
	Ancestry string `json:"ancestry"`
	IconicTaxonID int `json:"iconic_taxon_id"`
	TaxonPhotos []struct {
		Photo Photo `json:"photo"`
		Taxon Taxon `json:"taxon"`
	} `json:"taxon_photos"`
	RankLevel int `json:"rank_level"`
	TaxonChangesCount int `json:"taxon_changes_count"`
	AtlasID int `json:"atlas_id"`
	ParentID int `json:"parent_id"`
	Name string `json:"name"`
	Rank string `json:"rank"`
	ID int64 `json:"id"`
	DefaultPhoto Photo `json:"default_photo"`
	AncestorIds []int `json:"ancestor_ids"`
	IconicTaxonName string `json:"iconic_taxon_name"`
	PreferredCommonName string `json:"preferred_common_name"`
	Ancestors []*Taxon `json:"ancestors"`
	Children []*Taxon `json:"children"`
	ListedTaxa []struct {
		ID int `json:"id"`
		TaxonID int `json:"taxon_id"`
		EstablishmentMeans string `json:"establishment_means"`
		Place struct {
			   ID int `json:"id"`
			   Name string `json:"name"`
			   DisplayName string `json:"display_name"`
			   AdminLevel int `json:"admin_level"`
			   AncestorPlaceIds []int `json:"ancestor_place_ids"`
		   } `json:"place"`
		List struct {
			   ID int `json:"id"`
			   Title string `json:"title"`
		   } `json:"list"`
	} `json:"listed_taxa"`
	WikipediaSummary string `json:"wikipedia_summary"`
	MinSpeciesAncestry string `json:"min_species_ancestry"`
	CreatedAt time.Time `json:"created_at"`
}

type Photo struct {
	Flags []interface{} `json:"flags"`
	Type string `json:"type"`
	URL string `json:"url"`
	SquareURL string `json:"square_url"`
	NativePageURL string `json:"native_page_url"`
	NativePhotoID string `json:"native_photo_id"`
	SmallURL string `json:"small_url"`
	Attribution string `json:"attribution"`
	MediumURL string `json:"medium_url"`
	ID int `json:"id"`
	LicenseCode string `json:"license_code"`
	OriginalDimensions interface{} `json:"original_dimensions"`
	LargeURL string `json:"large_url"`
}

func (Ω Photo) Format(parentKey *datastore.Key) *store.Photo {
	return &store.Photo{
		Key: datastore.NameKey(store.EntityKindPhoto, strconv.Itoa(Ω.ID), parentKey),
		Type: store.PhotoType(Ω.Type),
		URL: Ω.URL,
		SquareURL: Ω.SquareURL,
		SmallURL: Ω.SmallURL,
		MediumURL: Ω.MediumURL,
		LargeURL: Ω.LargeURL,
		NativePhotoID: Ω.NativePageURL,
		Attribution: Ω.Attribution,
		LicenseCode: Ω.LicenseCode,
		Flags: Ω.Flags,
	}
}

func (Ω *fetcher) fetchProcessTaxon(taxonID store.TaxonID) error {

	if !taxonID.Valid() {
		return errors.New("invalid taxon id")
	}

	// Check to see if we've already processed it. Having states likely means we've grabbed the full page already.
	if i := Ω.Taxa.Index(taxonID); i != -1 && len(Ω.Taxa[i].States) > 0 {
		return nil
	}

	var response struct {
		Counter
		Results []*Taxon `json:"results"`
	}

	url := fmt.Sprintf("http://api.inaturalist.org/v1/taxa/%d", taxonID)

	if err := request(url, &response); err != nil {
		return err
	}

	if response.TotalResults == 0 {
		return errors.New("no taxon returned from ").SetState(M{utils.LogkeyIdentifier: taxonID})
	}

	if response.TotalResults > 1 {
		return errors.New("taxon request has more than one result").SetState(M{utils.LogkeyDatastoreKey: taxonID})
	}

	// Should only be one result.
	taxa := response.Results[0]

	var lastAncestor *datastore.Key
	for i, _a := range taxa.Ancestors {
		a := _a
		key, err := Ω.processTaxon(a, lastAncestor, (i != 0))
		if err != nil {
			return err
		}
		lastAncestor = key
	}
	Ω.Tomb.Go(func() error {
		_, err := Ω.processTaxon(taxa, lastAncestor, true)
		return err
	})

	Ω.Tomb.Go(func() error {
		for _, _a := range taxa.Children {
			a := _a
			rank := store.RankLevel(a.RankLevel)
			if rank == store.RankLevelSpecies || rank == store.RankLevelSubSpecies {
				Ω.Tomb.Go(func() error {
					if err := Ω.fetchProcessTaxon(store.TaxonID(a.ID)); err != nil {
						return err
					}
					return nil
				})
			}
		}
		return nil
	})

	return nil
}

func (Ω *fetcher) processTaxon(txn *Taxon, parent *datastore.Key, shouldHaveParent bool) (*datastore.Key, error) {

	// No need to reprocess a parent key we've already processed it before. Noopt and return the key.
	//taxonKey := Ω.Fetched.Find(txn.ID, store.EntityKindTaxon)
	//if taxonKey != nil {
	//	return taxonKey, nil
	//}

	taxonKey := store.NewTaxonKey(txn.ID)
	// Add parent as last ancestory. Again, these should be in order.
	if parent != nil {
		taxonKey.Parent = parent
	}

	t := &store.Taxon{
		Key: taxonKey,
		CanonicalName: store.CanonicalName(txn.Name),
		Rank: store.TaxonRank(txn.Rank),
		RankLevel: store.RankLevel(txn.RankLevel),
		CommonName: txn.PreferredCommonName,
		CreatedAt: Ω.Clock.Now(),
		ModifiedAt: Ω.Clock.Now(),
		WikipediaSummary: txn.WikipediaSummary,
	}

	for _, lt := range txn.ListedTaxa {
		if lt.Place.AdminLevel == 1 {
			t.States = append(t.States, store.State{
				EstablishmentMeans: lt.EstablishmentMeans,
				Name: lt.Place.Name,
			})
		}
	}

	// Create this taxon
	//var err error
	//taxonKey, err = Ω.Store.CreateTaxon(t, shouldHaveParent)
	//if err != nil {
	//	return nil, errors.Wrap(err, "could not save photo")
	//}

	// Lock and attach to list.
	Ω.Tomb.Go(func() error {
		Ω.Lock()
		defer Ω.Unlock()
		var err error
		Ω.Taxa, err = Ω.Taxa.AddToSet(t)
		if err != nil {
			return err
		}
		return nil
	})

	Ω.Tomb.Go(func() error {
		Ω.Lock()
		defer Ω.Unlock()
		var err error
		// Create a inaturalist scheme, a subobject of the taxon.
		iNatScheme := store.NewMetaScheme(
			store.SchemeSourceIDINaturalist,
			store.SchemeTargetID(strconv.FormatInt(taxonKey.ID, 10)),
			taxonKey,
		)
		Ω.Schema, err = Ω.Schema.AddToSet(iNatScheme)
		if err != nil {
			return err
		}

		// Note that the photos store sub-species, and so far the only place i can find them.
		for _, p := range txn.TaxonPhotos {
			if p.Taxon.ID == t.Key.ID {
				Ω.Photos, err = Ω.Photos.AddToSet(p.Photo.Format(iNatScheme.Key))
				if err != nil {
					return err
				}
			}
		}
		// Add default photo, if present
		Ω.Photos, err = Ω.Photos.AddToSet(txn.DefaultPhoto.Format(iNatScheme.Key))
		if err != nil {
			return err
		}
		return nil
	})

	Ω.Tomb.Go(func() error {
		// Fetch additional schemes from the iNaturalist, like the GBIF
		if schemes, err := Ω.fetchSchemes(taxonKey, (t.RankLevel == store.RankLevelSubSpecies || t.RankLevel == store.RankLevelSpecies)); err != nil {
			return err
		} else if len(schemes) > 0 {
			Ω.Lock()
			defer Ω.Unlock()
			for _, s := range schemes {
				Ω.Schema, err = Ω.Schema.AddToSet(s)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})

	return taxonKey, nil
}

var schemeRegex = regexp.MustCompile(`\(([^\)]+)\)`)

func (Ω *fetcher) fetchSchemes(txn *datastore.Key, isSpecies bool) ([]*store.Scheme, error) {

	if !store.ValidTaxonKey(txn) {
		return nil, errors.New("invalid taxonID").SetState(M{utils.LogkeyTaxon: txn})
	}

	url := fmt.Sprintf("http://www.inaturalist.org/taxa/%d/schemes", txn.ID)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, errors.Wrap(err, "could parse site for goquery").SetState(M{utils.LogkeyURL: url})
	}

	res := []*store.Scheme{}
	// Find the review items
	pairs := []struct{
		OriginID store.SchemeSourceID
		TargetID store.SchemeTargetID
	}{}
	doc.Find(`a[href*="/taxon_schemes/"]`).Each(func(i int, s *goquery.Selection) {
		v, _ := s.Attr("href")
		originID := store.SchemeSourceID(strings.TrimLeft(v, "/taxon_schemes/"))
		if string(originID) == "" {
			return
		}
		dataID := schemeRegex.FindString(s.Parent().Text())
		if dataID == "" {
			return
		}
		targetID := store.SchemeTargetID(strings.TrimRight(strings.TrimLeft(dataID, "("), ")"))
		pairs = append(pairs, struct{
			OriginID store.SchemeSourceID
			TargetID store.SchemeTargetID
		}{originID, targetID})
	})
	for _, pair := range pairs {
		if pair.OriginID == store.SchemeSourceIDGBIF && isSpecies {
			res = append(res, store.NewOccurrenceScheme(store.SchemeSourceIDGBIF, pair.TargetID, txn))
		}
		res = append(res, store.NewMetaScheme(pair.OriginID, pair.TargetID, txn))
	}

	return res, nil

}

func request(url string, response interface{}) error {

	resp, err := http.Get(url)
	if err != nil {
		return errors.Wrapf(err, "could not get gbif http response from request: %s", url)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.Wrapf(errors.New(resp.Status), "code: %d; request: %s", resp.StatusCode, url)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "could not read http response body")
	}

	if err := json.Unmarshal(b, &response); err != nil {
		return errors.Wrap(err, "could not unmarshal http response")
	}

	return nil
}