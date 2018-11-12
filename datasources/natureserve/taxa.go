package natureserve

import (
	"github.com/heindl/floracast_process/utils"
	"context"
	"fmt"
	"github.com/kennygrant/sanitize"
	"gopkg.in/tomb.v2"
	url2 "net/url"
	"os"
	"strings"
	"sync"
)

var natureServeAPIKey = os.Getenv("FLORACAST_NATURESERVE_API_KEY")

func init() {
	if natureServeAPIKey == "" {
		panic("FLORACAST_NATURESERVE_API_KEY environment variable required")
	}
}

func fetchTaxaFromSearch(cxt context.Context, names ...string) ([]*taxon, error) {

	locker := sync.Mutex{}
	taxa := []*taxon{}

	limit := utils.NewLimiter(20)

	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _, _name := range names {
			name := _name
			done := limit.Go()
			tmb.Go(func() error {
				defer done()
				localTaxa, err := searchName(cxt, name)
				if err != nil {
					return err
				}
				locker.Lock()
				defer locker.Unlock()
				taxa = append(taxa, localTaxa...)
				return nil
			})
		}
		return nil
	})

	if err := tmb.Wait(); err != nil {
		return nil, err
	}

	return taxa, nil
}

func searchName(cxt context.Context, name string) ([]*taxon, error) {

	url := fmt.Sprintf(
		"https://services.natureserve.org/idd/rest/ns/v1/globalSpecies/list/nameSearch?NSAccessKeyId=%s&name=%s",
		natureServeAPIKey,
		url2.QueryEscape(name),
	)

	searchResults := speciesSearchReport{}
	if err := utils.RequestXML(url, &searchResults); err != nil {
		return nil, err
	}

	if searchResults.SpeciesSearchResultList == nil {
		return nil, nil
	}

	uids := []string{}

	for _, sr := range searchResults.SpeciesSearchResultList.SpeciesSearchResult {
		//if sr.GlobalSpeciesUID == nil {
		//	fmt.Println("NIL UID", utils.JsonOrSpew(sr))
		//	continue
		//}
		uids = append(uids, sr.Attruid)
	}

	uids = utils.RemoveStringDuplicates(uids)

	if len(uids) == 0 {
		return nil, nil
	}

	return fetchTaxaWithUID(cxt, uids...)

}

//func FetchTaxonWithUID(cxt context.Context, uid string, referencedCanonicalName string) (*taxon, error) {
//	taxa, err := fetchTaxaWithUID(cxt, uid)
//	if err != nil {
//		return nil, err
//	}
//	txn := taxa[0]
//	if !utils.ContainsString(txn.ScientificNameStrings(), referencedCanonicalName) {
//		return nil, errors.Newf("Expected nature serve taxon [%s] to contain reference canonical name [%s]", uid, referencedCanonicalName)
//	}
//	return txn, nil
//}

func fetchTaxaWithUID(_ context.Context, uids ...string) ([]*taxon, error) {

	if len(uids) == 0 {
		return nil, nil
	}

	url := fmt.Sprintf(
		"https://services.natureserve.org/idd/rest/ns/v1.1/globalSpecies/comprehensive?NSAccessKeyId=%s&uid=%s",
		natureServeAPIKey,
		strings.Join(uids, ","),
	)

	speciesList := globalSpeciesList{}
	if err := utils.RequestXML(url, &speciesList); err != nil {
		return nil, err
	}

	taxa := []*taxon{}

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

type taxon struct {
	//kingdom string `json:",omitempty"`
	//phylum string `json:",omitempty"`
	//class string `json:",omitempty"`
	//order string `json:",omitempty"`
	//family string `json:",omitempty"`
	//genus string `json:",omitempty"`
	ID             string                 `json:",omitempty"`
	ScientificName *taxonScientificName   `json:",omitempty"`
	Synonyms       []*taxonScientificName `json:",omitempty"`
	CommonNames    []*taxonCommonName     `json:",omitempty"`
}

func (Ω *taxon) ScientificNameStrings() []string {
	res := []string{strings.ToLower(Ω.ScientificName.Name)}
	for _, synonym := range Ω.Synonyms {
		res = append(res, strings.ToLower(synonym.Name))
	}
	return res
}

// TODO: Convert commonName to struct with language field.

type taxonScientificName struct {
	Name                                 string `json:",omitempty"`
	Author                               string `json:",omitempty"`
	ConceptReferenceCode                 string `json:",omitempty"`
	ConceptReferenceFullCitation         string `json:",omitempty"`
	ConceptReferenceNameUsed             string `json:",omitempty"`
	ConceptReferenceClassificationStatus string `json:",omitempty"`
}

type taxonCommonName struct {
	LanguageCode string `json:",omitempty"`
	IsPrimary    bool   `json:",omitempty"`
	Name         string `json:",omitempty"`
}

// TODO: Really dig more into various available information. informalTaxonomy is useful for choosing broad categories
// that can be filterable by the user.
// TODO: Dig more into the conservation status. Very good for specifying how careful to be with each species.
// For some species, great information about migration patterns.
// https://services.natureserve.org/idd/rest/ns/v1.1/globalSpecies/comprehensive?NSAccessKeyId=b2374ab2-275c-48eb-b3c1-8f7afe9af5c4&uid=ELEMENT_GLOBAL.2.116078,ELEMENT_GLOBAL.2.121086,ELEMENT_GLOBAL.2.735443,ELEMENT_GLOBAL.2.735442,ELEMENT_GLOBAL.9.24619,ELEMENT_GLOBAL.9.24616,ELEMENT_GLOBAL.2.108328,ELEMENT_GLOBAL.2.114107,ELEMENT_GLOBAL.2.121010,ELEMENT_GLOBAL.2.107284,ELEMENT_GLOBAL.2.111490,ELEMENT_GLOBAL.2.108561,ELEMENT_GLOBAL.2.107412,ELEMENT_GLOBAL.2.115920,ELEMENT_GLOBAL.2.108251,ELEMENT_GLOBAL.2.116121,ELEMENT_GLOBAL.2.841062,ELEMENT_GLOBAL.2.841061,ELEMENT_GLOBAL.9.24619

func parseGlobalSpecies(spcs *globalSpecies) (*taxon, error) {

	if spcs.Classification == nil {
		fmt.Println(fmt.Sprintf("Warning: Invalid/missing NatureServe species [%s]", spcs.Attruid))
		return nil, nil
	}

	txn := taxon{
		ID: spcs.Attruid,
	}

	//if species.classification.taxonomy != nil && species.classification.taxonomy.formalTaxonomy != nil {
	//	ft := species.classification.taxonomy.formalTaxonomy
	//	txn.kingdom = ft.kingdom.Text
	//	txn.phylum = ft.phylum.Text
	//	txn.class = ft.class.Text
	//	txn.order = ft.order.Text
	//	txn.family = ft.family.Text
	//	txn.genus = ft.genus.Text
	//}

	if err := setNames(spcs, &txn); err != nil {
		return nil, err
	}

	return &txn, nil

}

func setNames(species *globalSpecies, txn *taxon) (err error) {
	if species.Classification == nil || species.Classification.Names == nil {
		return nil
	}

	names := species.Classification.Names

	txn.ScientificName, err = names.ScientificName.asTaxonScientificName()
	if err != nil {
		return err
	}

	if sn := parseXMLScientificName(names.ScientificName); sn != nil {
		txn.ScientificName = sn
	}

	txn.Synonyms, err = names.taxonSynonyms()
	if err != nil {
		return err
	}

	txn.CommonNames = []*taxonCommonName{}

	if names.NatureServePrimaryGlobalCommonName != nil {
		txn.CommonNames = append(txn.CommonNames, &taxonCommonName{
			IsPrimary: true,
			Name:      names.NatureServePrimaryGlobalCommonName.Text,
		})
	}

	if names.OtherGlobalCommonNames != nil && len(names.OtherGlobalCommonNames.CommonName) > 0 {
		for _, cn := range names.OtherGlobalCommonNames.CommonName {
			txn.CommonNames = append(txn.CommonNames, &taxonCommonName{
				IsPrimary:    false,
				Name:         cn.Text,
				LanguageCode: cn.Attrlanguage,
			})
		}
	}
	return nil
}

func parseXMLScientificName(given *scientificName) *taxonScientificName {

	name := given.CanonicalName()
	if name == "" {
		return nil
	}

	txn := taxonScientificName{
		Name: name,
	}

	if given.NomenclaturalAuthor != nil {
		txn.Author = given.NomenclaturalAuthor.Text
	}

	if given.ConceptReference == nil {
		return &txn
	}

	cr := given.ConceptReference
	if cr.ClassificationStatus != nil {
		txn.ConceptReferenceClassificationStatus = cr.ClassificationStatus.Text
	}
	if crnu := cr.NameUsedInConceptReference; crnu != nil {
		if crnu.UnformattedName != nil {
			txn.ConceptReferenceNameUsed = crnu.UnformattedName.Text
		} else if crnu.FormattedName != nil {
			if len(crnu.FormattedName.I) > 0 {
				txn.ConceptReferenceNameUsed = crnu.FormattedName.I[0].Text
			} else {
				txn.ConceptReferenceNameUsed = sanitize.HTML(crnu.FormattedName.Text)
			}

		}
	}
	if cr.FormattedFullCitation != nil {
		txn.ConceptReferenceFullCitation = cr.FormattedFullCitation.Text
	}
	txn.ConceptReferenceCode = cr.Attrcode

	return &txn
}
