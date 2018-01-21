package nature_serve

import (
	"context"
	"os"
	"fmt"
	"bitbucket.org/heindl/taxa/utils"
	"sync"
	"gopkg.in/tomb.v2"
	"strings"
)

var natureServeAPIKey = os.Getenv("FLORACAST_NATURESERVE_API_KEY")

func init() {
	if natureServeAPIKey == "" {
		panic("FLORACAST_NATURESERVE_API_KEY environment variable required")
	}
}

func FetchTaxaFromSearch(cxt context.Context, names ...string) ([]*Taxon, error) {

	locker := sync.Mutex{}
	taxa := []*Taxon{}

	tmb := tomb.Tomb{}
	tmb.Go(func()error {
		for _, _name := range names {
			name := _name
			tmb.Go(func() error {
				local_taxa, err := searchName(cxt, name)
				if err != nil {
					return err
				}
				locker.Lock()
				defer locker.Unlock()
				taxa = append(taxa, local_taxa...)
				return nil
			})
		}
		return nil
	})

	if err := tmb.Wait(); err != nil {
		return nil, err
	}

	return nil, nil
}

func searchName(cxt context.Context, name string) ([]*Taxon, error) {
	url := fmt.Sprintf(
		"https://services.natureserve.org/idd/rest/ns/v1/globalSpecies/list/nameSearch?&NSAccessKeyId=%s&name=%s",
		natureServeAPIKey,
		name,
	)

	searchResults := SpeciesSearchReport{}
	if err := utils.RequestXML(url, &searchResults); err != nil {
		return nil, err
	}

	fmt.Println(utils.JsonOrSpew(searchResults))

	if searchResults.SpeciesSearchResultList == nil {
		return nil, nil
	}

	uids := []string{}

	for _, sr := range searchResults.SpeciesSearchResultList.SpeciesSearchResult {
		uids = append(uids, sr.GlobalSpeciesUid.Text)
	}

	return FetchTaxaWithUID(cxt, uids...)

}

func FetchTaxaWithUID(cxt context.Context, uids ...string) ([]*Taxon, error) {
	
	if len(uids) == 0 {
		return nil, nil
	}

	url := fmt.Sprintf(
		"https://services.natureserve.org/idd/rest/ns/v1.1/globalSpecies/comprehensive?NSAccessKeyId=%s&uid=%s",
		natureServeAPIKey,
		strings.Join(uids, ","),
		)
	
	speciesList := GlobalSpeciesList{}
	if err := utils.RequestXML(url, &speciesList); err != nil {
		return nil, err
	}

	taxa := []*Taxon{}

	for _, globalSpecies := range speciesList.GlobalSpecies {
		txn, err := parseGlobalSpecies(globalSpecies)
		if err != nil {
			return nil, err
		}
		if txn != nil {
			taxa = append(taxa, txn)
		}
	}

	return taxa, nil

}

type Taxon struct {
	Kingdom string
	Phylum string
	Class string
	Order string
	Family string
	Genus string
	ScientificName *TaxonScientificName
	Synonyms []*TaxonScientificName
	CommonNames []*TaxonCommonName
}

type TaxonScientificName struct {
	Name string
	Author string
	ConceptReferenceCode string
	ConceptReferenceFullCitation string
	ConceptReferenceNameUsed string
	ConceptReferenceClassificationStatus string
}

type TaxonCommonName struct {
	LanguageCode string
	IsPrimary bool
	Name string
}

// TODO: Really dig more into various available information. InformalTaxonomy is useful for choosing broad categories
// that can be filterable by the user.
// TODO: Dig more into the conservation status. Very good for specifying how careful to be with each species.
// For some species, great information about migration patterns.
// https://services.natureserve.org/idd/rest/ns/v1.1/globalSpecies/comprehensive?NSAccessKeyId=b2374ab2-275c-48eb-b3c1-8f7afe9af5c4&uid=ELEMENT_GLOBAL.2.116078,ELEMENT_GLOBAL.2.121086,ELEMENT_GLOBAL.2.735443,ELEMENT_GLOBAL.2.735442,ELEMENT_GLOBAL.9.24619,ELEMENT_GLOBAL.9.24616,ELEMENT_GLOBAL.2.108328,ELEMENT_GLOBAL.2.114107,ELEMENT_GLOBAL.2.121010,ELEMENT_GLOBAL.2.107284,ELEMENT_GLOBAL.2.111490,ELEMENT_GLOBAL.2.108561,ELEMENT_GLOBAL.2.107412,ELEMENT_GLOBAL.2.115920,ELEMENT_GLOBAL.2.108251,ELEMENT_GLOBAL.2.116121,ELEMENT_GLOBAL.2.841062,ELEMENT_GLOBAL.2.841061,ELEMENT_GLOBAL.9.24619

func parseGlobalSpecies(species *GlobalSpecies) (*Taxon, error) {

	if species.Classification == nil {
		fmt.Println("Invalid or missing classification", species.Attruid)
		return nil, nil
	}

	txn := Taxon{}

	if species.Classification.Taxonomy != nil && species.Classification.Taxonomy.FormalTaxonomy != nil {
		ft := species.Classification.Taxonomy.FormalTaxonomy
		txn.Kingdom = ft.Kingdom.Text
		txn.Phylum = ft.Phylum.Text
		txn.Class = ft.Class.Text
		txn.Order = ft.Order.Text
		txn.Family = ft.Family.Text
		txn.Genus = ft.Genus.Text
	}

	if species.Classification.Names != nil {

		names := species.Classification.Names

		if sn := parseXMLScientificName(names.ScientificName); sn != nil {
			txn.ScientificName = sn
		}

		if len(names.Synonyms.SynonymName) > 0 {
			txn.Synonyms = []*TaxonScientificName{}
			for _, synonym := range names.Synonyms.SynonymName {
				if sn := parseXMLScientificName(synonym); sn != nil {
					txn.Synonyms = append(txn.Synonyms, sn)
				}
			}
		}

		txn.CommonNames = []*TaxonCommonName{}

		if names.NatureServePrimaryGlobalCommonName != nil {
			txn.CommonNames = append(txn.CommonNames, &TaxonCommonName{
				IsPrimary: true,
				Name: names.NatureServePrimaryGlobalCommonName.Text,
			})
		}

		if len(names.OtherGlobalCommonNames.CommonName) > 0 {
			for _, cn := range names.OtherGlobalCommonNames.CommonName {
				txn.CommonNames = append(txn.CommonNames, &TaxonCommonName{
					IsPrimary: false,
					Name: cn.Text,
					LanguageCode:cn.Attrlanguage,
				})
			}
		}
	}

	return &txn, nil

}


func parseXMLScientificName(scientific *ScientificName, synonym *SynonymName) *TaxonScientificName {

	if scientific == nil && synonym == nil {
		return nil
	}

	g

	if given == nil {
		return nil
	}

	name := ""
	if given.UnformattedName != nil {
		name = given.UnformattedName.Text
	} else if given.FormattedName != nil {
		name = given.FormattedName.I[0].Text
	} else {
		return nil
	}

	txn := TaxonScientificName{
		Name: name,
	}

	if given.NomenclaturalAuthor != nil {
		txn.Author = given.NomenclaturalAuthor.Text
	}

	if cr := given.ConceptReference; cr != nil {
		if cr.ClassificationStatus != nil {
			txn.ConceptReferenceClassificationStatus = cr.ClassificationStatus.Text
		}
		if crnu := cr.NameUsedInConceptReference; crnu != nil {
			if crnu.UnformattedName != nil {
				txn.ConceptReferenceNameUsed = crnu.UnformattedName.Text
			} else if crnu.FormattedName != nil {
				txn.ConceptReferenceNameUsed = crnu.FormattedName.I[0].Text
			}
		}
		if cr.FormattedFullCitation != nil {
			txn.ConceptReferenceFullCitation = cr.FormattedFullCitation.Text
		}
		txn.ConceptReferenceCode = cr.Attrcode
	}

	return &txn
}