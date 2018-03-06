package gbif

import (
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/utils"
	"context"
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"gopkg.in/tomb.v2"
	"math"
	"sync"
	"time"
)

// FetchOccurrences gathers a list of OccurrenceProvider interfaces.
func FetchOccurrences(cxt context.Context, targetID datasources.TargetID, since *time.Time) ([]*occurrence, error) {

	txnID, err := taxonIDFromTargetID(targetID)
	if err != nil {
		return nil, err
	}

	lastInterpreted := ""
	if since != nil && !since.IsZero() {
		lastInterpreted = since.Format("20060102")
	}

	apiList, err := fetchOccurrences(occurrenceSearchQuery{
		TaxonKey:        int(txnID),
		LastInterpreted: lastInterpreted,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not fetch occurrences [%d] from the gbif", txnID)
	}

	res := []*occurrence{}

	for _, gbifO := range apiList {

		// Check geospatial issue.
		if gbifO.Issues.isUnacceptable() {
			continue
		}

		if gbifO.Issues.hasIssue(occurrenceIssuePresumedSwappedCoordinate) {
			fmt.Println("WARNING:", occurrenceIssuePresumedSwappedCoordinate, fmt.Sprintf("Latitude/Longitude [%f, %f]", gbifO.DecimalLatitude, gbifO.DecimalLongitude))
			continue
		}

		if gbifO.Issues.isUncertain() {
			continue
		}
		res = append(res, gbifO)

	}

	return res, nil

}

type occurrenceSearchQuery struct {
	// Continent, as defined in our Continent enum.
	// WARNING: that this is an unreliable filter, as some records have not saved the continent.
	Continent continentEnum `json:"continent"`
	// The maximum number of results to return. This can't be greater than 300, any value greater is set to 300.
	Limit int `json:"limit"`
	// The month of the year, starting with 1 for January. Supports range queries.
	Month int `json:"month"`
	// A taxon key from the GBIF backbone. All included and synonym taxa are included in the search,
	// so a search for aves with taxonKey=212 (i.e. /occurrence/search?taxonKey=212) will match all birds, no matter which species.
	TaxonKey int `json:"taxonKey"`
	// The 4 digit year. A year of 98 will be interpreted as AD 98.
	// Supports range queries: http://api.gbif.org/v1/occurrence/search?year=1800,1899
	Year int `json:"year"`
	// Limits searches to occurrence records which contain a value in both latitude and longitude
	// (i.e. hasCoordinate=true limits to occurrence records with coordinate values and hasCoordinate=false limits to occurrence records without coordinate values).
	HasCoordinate bool `json:"hasCoordinate"`
	// Includes/excludes occurrence records which contain spatial issues (as determined in our record interpretation), i.e. hasGeospatialIssue=true returns only those records with spatial issues while hasGeospatialIssue=false includes only records without spatial issues. The absence of this parameter returns any record with or without spatial issues.
	HasGeospatialIssue bool `json:"hasGeospatialIssue"`
	// This date the record was last modified in GBIF, in ISO 8601 format: yyyy, yyyy-MM, yyyy-MM-dd, or MM-dd. Supports range queries.
	LastInterpreted string `json:"lastInterpreted"`
}

type continentEnum string

func (c continentEnum) IsZero() bool {
	return string(c) == ""
}

//const (
//	ContinentEnumAfrica       = continentEnum("AFRICA")
//	ContinentEnumAntartica    = continentEnum("ANTARCTICA")
//	ContinentEnumAsia         = continentEnum("ASIA")
//	ContinentEnumEurope       = continentEnum("EUROPE")
//	ContinentEnumNorthAmerica = continentEnum("NORTH_AMERICA")
//	ContinentEnumNorthOceania = continentEnum("OCEANIA")
//	ContinentEnumSouthAmerica = continentEnum("SOUTH_AMERICA")
//)

// datasetKey, year, month, eventDate, lastInterpreted, decimalLatitude, decimalLongitude, country, continent, publishingCountry, elevation, depth, institutionCode, collectionCode, catalogNumber, recordedBy, recordNumber, basisOfRecord, taxonKey, scientificName, hasCoordinate, geometry, hasGeospatialIssue, issue, mediaType

func (q *occurrenceSearchQuery) url(offset int) string {

	if q.Limit == 0 {
		q.Limit = 300
	}

	u := fmt.Sprintf("http://api.gbif.org/v1/occurrence/search?taxonKey=%d&offset=%d&limit=%d",
		q.TaxonKey,
		offset*q.Limit,
		q.Limit,
	)

	if !q.Continent.IsZero() {
		u += fmt.Sprintf("&continent=%s", q.Continent)
	}

	if q.HasCoordinate {
		u += "&hasCoordinate=true"
	}

	if q.HasGeospatialIssue {
		u += "&hasGeospatialIssue=true"
	}

	if q.LastInterpreted != "" {
		u += fmt.Sprintf("&lastInterpreted=%s", q.LastInterpreted)
	}

	return u
}

// A note about data found: http://tools.gbif.org/dwca-validator/extension.do?id=gbif:TypesAndSpecimen

type occurrence struct {
	BasisOfRecord                   string           `json:"basisOfRecord"`
	CatalogNumber                   string           `json:"catalogNumber"`
	Class                           string           `json:"class"`
	ClassKey                        int              `json:"classKey"`
	CollectionCode                  string           `json:"collectionCode"`
	Country                         string           `json:"country"`
	CountryCode                     string           `json:"countryCode"`
	DatasetKey                      string           `json:"datasetKey"`
	DatasetName                     string           `json:"datasetName"`
	DateIdentified                  interface{}      `json:"dateIdentified"`
	Day                             int              `json:"day"`
	DecimalLatitude                 float64          `json:"decimalLatitude"`
	DecimalLongitude                float64          `json:"decimalLongitude"`
	EventDate                       gbifTime         `json:"eventDate"`
	EventTime                       string           `json:"eventTime"`
	Extensions                      struct{}         `json:"extensions"`
	Facts                           []interface{}    `json:"facts"`
	Family                          string           `json:"family"`
	FamilyKey                       int              `json:"familyKey"`
	GbifID                          string           `json:"gbifID"`
	GenericName                     string           `json:"genericName"`
	Genus                           string           `json:"genus"`
	GenusKey                        int              `json:"genusKey"`
	GeodeticDatum                   string           `json:"geodeticDatum"`
	HTTPUnknownOrgOccurrenceDetails string           `json:"http://unknown.org/occurrenceDetails"`
	IdentificationID                string           `json:"identificationID"`
	Identifier                      string           `json:"identifier"`
	Identifiers                     []interface{}    `json:"identifiers"`
	InstitutionCode                 string           `json:"institutionCode"`
	Issues                          occurrenceIssues `json:"issues"`
	Key                             int              `json:"key"`
	Kingdom                         string           `json:"kingdom"`
	KingdomKey                      int              `json:"kingdomKey"`
	LastCrawled                     interface{}      `json:"lastCrawled"`
	LastInterpreted                 interface{}      `json:"lastInterpreted"`
	LastParsed                      interface{}      `json:"lastParsed"`
	Media                           []media          `json:"media"`
	Modified                        interface{}      `json:"modified"`
	Month                           int              `json:"month"`
	OccurrenceID                    string           `json:"occurrenceID"`
	OccurrenceRemarks               string           `json:"occurrenceRemarks"`
	Order                           string           `json:"order"`
	OrderKey                        int              `json:"orderKey"`
	Phylum                          string           `json:"phylum"`
	PhylumKey                       int              `json:"phylumKey"`
	Protocol                        string           `json:"protocol"`
	PublishingCountry               string           `json:"publishingCountry"`
	PublishingOrgKey                string           `json:"publishingOrgKey"`
	RecordedBy                      string           `json:"recordedBy"`
	References                      string           `json:"references"`
	Relations                       []interface{}    `json:"relations"`
	Rights                          string           `json:"rights"`
	RightsHolder                    string           `json:"rightsHolder"`
	OccurrenceStatus                string           `json:"occurrenceStatus"`
	ScientificName                  string           `json:"scientificName"`
	Species                         string           `json:"species"`
	SpeciesKey                      int              `json:"speciesKey"`
	SpecificEpithet                 string           `json:"specificEpithet"`
	TaxonID                         string           `json:"taxonID"`
	TaxonKey                        int              `json:"taxonKey"`
	TaxonRank                       string           `json:"taxonRank"`
	VerbatimEventDate               string           `json:"verbatimEventDate"`
	Year                            int              `json:"year"`
}

func (Ω *occurrence) Lat() (float64, error) {
	return Ω.DecimalLatitude, nil
}

func (Ω *occurrence) Lng() (float64, error) {
	return Ω.DecimalLongitude, nil
}

func (Ω *occurrence) DateString() string {
	return Ω.EventDate.Time.Format("20060102")
}

func (Ω *occurrence) CoordinatesEstimated() bool {
	// Rounded to 5 decimal place. Not what I expected.
	// isEstimated := s.Issues.hasIssue(ogbif.OCCURRENCE_ISSUE_COORDINATE_ROUNDED)
	return false
}

func (Ω *occurrence) SourceOccurrenceID() string {
	return Ω.GbifID
}

type response struct {
	page
	Results []*occurrence `json:"results"`
}

// TotalOccurrenceCount returns a full search across all occurrences. Results are ordered by relevance.
func fetchOccurrences(q occurrenceSearchQuery) ([]*occurrence, error) {

	var syncList struct {
		sync.Mutex
		list []*occurrence
	}

	response, err := q.request(0)
	if err != nil {
		return nil, err
	}
	syncList.list = append(syncList.list, response.Results...)

	offsets := math.Ceil(float64(response.Count) / float64(response.Limit))
	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _i := 1; _i < int(offsets); _i++ {
			i := _i
			tmb.Go(func() error {
				response, err := q.request(i)
				if err != nil {
					return err
				}
				syncList.Lock()
				defer syncList.Unlock()
				syncList.list = append(syncList.list, response.Results...)
				return nil
			})
		}
		return nil
	})

	if err := tmb.Wait(); err != nil {
		return nil, err
	}

	return syncList.list, nil

}

func (Ω *occurrenceSearchQuery) request(offset int) (*response, error) {
	res := response{}
	if err := utils.RequestJSON(Ω.url(offset), &res); err != nil {
		return nil, err
	}
	return &res, nil
}
