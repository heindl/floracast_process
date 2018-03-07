package gbif

import (
	"bitbucket.org/heindl/process/utils"
	"fmt"
	"net/url"
	"strings"
)

type matchResult struct {
	Alternatives []nameUsage `json:"alternatives"`
	nameUsage
}

type matchQuery struct {
	// Optional class classification accepting a canonical name.
	Class string
	// Optional family classification accepting a canonical name.
	Family string
	// Optional genus classification accepting a canonical name.
	Genus string
	// Optional kingdom classification accepting a canonical name.
	Kingdom string
	// A scientific name which can be either a case insensitive filter for a canonical namestring, e.g. 'Puma concolor', or an input to the name parser
	Name string
	// Optional order classification accepting a canonical name.
	Order string
	// Optional phylum classification accepting a canonical name.
	Phylum string
	// Filters by taxonomic rank as given in our rank enum
	Rank string
	// If true it (fuzzy) matches only the given name, but never a taxon in the upper classification
	Strict bool
	// If true it shows alternative matches which were considered but then rejected
	Verbose bool
}

func (q matchQuery) url() string {

	u := fmt.Sprintf("http://api.gbif.org/v1/species/match?verbose=%v&strict=%v", q.Verbose, q.Strict)

	if q.Class != "" {
		u += fmt.Sprintf("&class=%s", url.QueryEscape(q.Class))
	}

	if q.Family != "" {
		u += fmt.Sprintf("&family=%s", url.QueryEscape(q.Family))
	}

	if q.Genus != "" {
		u += fmt.Sprintf("&genus=%s", url.QueryEscape(q.Genus))
	}

	if q.Kingdom != "" {
		u += fmt.Sprintf("&kingdom=%s", url.QueryEscape(q.Kingdom))
	}

	if q.Name != "" {
		name := utils.CapitalizeString(url.QueryEscape(strings.ToLower(q.Name)))
		u += fmt.Sprintf("&name=%s", name)
	}

	if q.Order != "" {
		u += fmt.Sprintf("&order=%s", url.QueryEscape(q.Order))
	}

	if q.Phylum != "" {
		u += fmt.Sprintf("&phylum=%s", url.QueryEscape(q.Phylum))
	}

	if q.Rank != "" {
		u += fmt.Sprintf("&rank=%s", url.QueryEscape(q.Rank))
	}

	return u
}

// Fuzzy matches scientific names against the GBIF Backbone Taxonomy with the optional classification provided.
// If a classification is provided and strict is not set to true,
// the default matching will also try to match against these if no direct match is found for the name parameter alone.
func match(q matchQuery) (response matchResult, err error) {
	if err := utils.RequestJSON(q.url(), &response); err != nil {
		return matchResult{}, err
	}
	return
}

//
//type searchQuery struct {
//	// Filters by the checklist dataset key (a uuid)
//	DatasetKey int `json:"datasetKey"`
//	// A list of facet names used to retrieve the 100 most frequent values for a field.
//	// Allowed facets are: datasetKey, higherTaxonKey, rank, status, isExtinct, habitat
//	// and nameType. Additionally threat and nomenclaturalStatus are legal values but
//	// not yet implemented, so data will not yet be returned for them.
//	Facet string `json:"facet"`
//	// Used in combination with the facet parameter. Set facetMincount={#} to exclude
//	// facets with a count less than {#}, e.g. /search?facet=type&limit=0&facetMincount=10000
//	// only shows the type value 'OCCURRENCE' because 'CHECKLIST' and 'METADATA'
//	// have counts less than 10000.
//	FacetMinCount int `json:"facetMincount"`
//	// Used in combination with the facet parameter. Set facetMultiselect=true to still
//	// return counts for values that are not currently filtered, e.g.
//	// /search?facet=type&limit=0&type=CHECKLIST&facetMultiselect=true still shows type
//	// values 'OCCURRENCE' and 'METADATA' even though type is being filtered by type=CHECKLIST
//	FacetMultiSelect bool `json:"facetMultiselect"`
//	// 	Filters by the habitat. Currently only 3 major biomes are accepted in our Habitat enum
//	Habitat string `json:"habitat"`
//	// Filters by any of the higher Linnean rank keys. Note this is within the respective
//	// checklist and not searching nub keys across all checklists.
//	HigherTaxonKey int `json:"highertaxonKey"`
//	// Set hl=true to highlight terms matching the query when in fulltext search fields.
//	// The highlight will be an emphasis tag of class 'gbifH1' e.g. /search?q=plant&hl=true.
//	// Fulltext search fields include: title, keyword, country, publishing country,
//	// publishing organization title, hosting organization title, and description.
//	// One additional full text field is searched which includes information from
//	// metadata documents, but the text of this field is not returned in the response.
//	HL bool `json:"hl"`
//	// Filters by extinction status (a boolean, e.g. isExtinct=true)
//	IsExtinct bool `json:"isExtinct"`
//	// Simple full text search parameter. The value for this parameter can be a simple
//	// word or a phrase. Wildcards are not supported.
//	Q string
//	// 	Filters by taxonomic rank as given in our rank enum
//	Rank []rank `json:"rank"`
//	// Filters by the taxonomic status as given in our taxonomicStatus enum
//	//Status []taxonomicStatus `json:"status"`
//	// The maximum number of results to return. This can't be greater than 300, any
//	// value greater is set to 300.
//	Limit int
//}
//
//func (q searchQuery) url(offset int) string {
//
//	if q.Limit == 0 {
//		q.Limit = 300
//	}
//
//	q.Q = strings.TrimSpace(q.Q)
//	q.Q = strings.Trim(q.Q, `"`)
//
//	u := fmt.Sprintf("http://api.gbif.org/v1/species/search?q=%s", url.QueryEscape(q.Q))
//
//	params := []string{
//		fmt.Sprintf("offset=%d", offset),
//		fmt.Sprintf("limit=%d", q.Limit),
//		fmt.Sprintf("datasetKey=%d", q.DatasetKey),
//		fmt.Sprintf("facet=%s", q.Facet),
//		fmt.Sprintf("facetMincount=%d", q.FacetMinCount),
//		fmt.Sprintf("facetMultiselect=%t", q.FacetMultiSelect),
//		fmt.Sprintf("habitat=%s", q.Habitat),
//		fmt.Sprintf("highertaxonKey=%d", q.HigherTaxonKey),
//		fmt.Sprintf("hl=%t", q.HL),
//		fmt.Sprintf("isExtinct=%t", q.IsExtinct),
//	}
//
//	filtered := []string{}
//
//	for _, p := range params {
//		if !strings.HasSuffix(p, "=0") && !strings.HasSuffix(p, "=") {
//			filtered = append(filtered, p)
//		}
//	}
//
//	for _, r := range q.Rank {
//		if string(r) != "" {
//			filtered = append(filtered, fmt.Sprintf("rank=%s", r))
//		}
//	}
//
//	//for _, s := range q.Status {
//	//	if string(s) != "" {
//	//		filtered = append(filtered, fmt.Sprintf("status=%s", s))
//	//	}
//	//}
//
//	return u + "&" + strings.Join(filtered, "&")
//}
//
//type searchResponse struct {
//	page
//	Results []nameUsage `json:"results"`
//}
//
//func search(q searchQuery) (names []nameUsage, err error) {
//
//	if q.Q == "" {
//		return nil, errors.New("a query value is required to search")
//	}
//
//	qResponse := searchResponse{}
//
//	if err := utils.RequestJSON(q.url(0), &qResponse); err != nil {
//		return nil, err
//	}
//
//	res := qResponse.Results
//	locker := sync.Mutex{}
//	totalRequests := math.Ceil(float64(qResponse.Count) / float64(qResponse.Limit))
//
//	tmb := tomb.Tomb{}
//	tmb.Go(func() error {
//		for _i := 1; _i <= int(totalRequests); _i++ {
//			i := _i
//			tmb.Go(func() error {
//				nqResponse := searchResponse{}
//				if err := utils.RequestJSON(q.url(i), &qResponse); err != nil {
//					return err
//				}
//				locker.Lock()
//				defer locker.Unlock()
//				res = append(res, nqResponse.Results...)
//				return nil
//			})
//		}
//		return nil
//	})
//	if err := tmb.Wait(); err != nil {
//		return nil, err
//	}
//
//	return res, nil
//
//}
