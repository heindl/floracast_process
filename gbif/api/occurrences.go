package api

import (
	"encoding/json"
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"gopkg.in/tomb.v2"
	"io/ioutil"
	"math"
	"sync"
	"github.com/sethgrid/pester"
	"time"
)

type OccurrenceSearchQuery struct {
	// Continent, as defined in our Continent enum.
	// WARNING: that this is an unreliable filter, as some records have not saved the continent.
	Continent ContinentEnum `json:"continent"`
	// The maximum number of results to return. This can't be greater than 300, any value greater is set to 300.
	Limit int `json:"limit"`
	// The month of the year, starting with 1 for January. Supports range queries.
	Month int `json:"month"`
	// A taxon key from the GBIF backbone. All included and synonym taxa are included in the search,
	// so a search for aves with taxonKey=212 (i.e. /occurrence/search?taxonKey=212) will match all birds, no matter which species.
	TaxonKey int `json:"taxonKey"`
	// The 4 digit year. A year of 98 will be interpreted as AD 98.
	// Supports range queries: http://api.gbif.org/v1//occurrence/search?year=1800,1899
	Year int `json:"year"`
	// Limits searches to occurrence records which contain a value in both latitude and longitude
	// (i.e. hasCoordinate=true limits to occurrence records with coordinate values and hasCoordinate=false limits to occurrence records without coordinate values).
	HasCoordinate bool `json:"hasCoordinate"`
	// Includes/excludes occurrence records which contain spatial issues (as determined in our record interpretation), i.e. hasGeospatialIssue=true returns only those records with spatial issues while hasGeospatialIssue=false includes only records without spatial issues. The absence of this parameter returns any record with or without spatial issues.
	HasGeospatialIssue bool `json:"hasGeospatialIssue"`
	// This date the record was last modified in GBIF, in ISO 8601 format: yyyy, yyyy-MM, yyyy-MM-dd, or MM-dd. Supports range queries.
	LastInterpreted string `json:"lastInterpreted"`
}

type ContinentEnum string

func (c ContinentEnum) IsZero() bool {
	return string(c) == ""
}

const (
	ContinentEnumAfrica       = ContinentEnum("AFRICA")
	ContinentEnumAntartica    = ContinentEnum("ANTARCTICA")
	ContinentEnumAsia         = ContinentEnum("ASIA")
	ContinentEnumEurope       = ContinentEnum("EUROPE")
	ContinentEnumNorthAmerica = ContinentEnum("NORTH_AMERICA")
	ContinentEnumNorthOceania = ContinentEnum("OCEANIA")
	ContinentEnumSouthAmerica = ContinentEnum("SOUTH_AMERICA")
)

// datasetKey, year, month, eventDate, lastInterpreted, decimalLatitude, decimalLongitude, country, continent, publishingCountry, elevation, depth, institutionCode, collectionCode, catalogNumber, recordedBy, recordNumber, basisOfRecord, taxonKey, scientificName, hasCoordinate, geometry, hasGeospatialIssue, issue, mediaType

func (q *OccurrenceSearchQuery) url(offset int) string {

	if q.Limit == 0 {
		q.Limit = 300
	}

	u := fmt.Sprintf("http://api.gbif.org/v1/occurrence/search?taxonKey=%d&offset=%d&limit=%d",
		q.TaxonKey,
		offset * q.Limit,
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

type Occurrence struct {
	BasisOfRecord                        string        `json:"basisOfRecord"`
	CatalogNumber                        string        `json:"catalogNumber"`
	Class                                string        `json:"class"`
	ClassKey                             int           `json:"classKey"`
	CollectionCode                       string        `json:"collectionCode"`
	Country                              string        `json:"country"`
	CountryCode                          string        `json:"countryCode"`
	DatasetKey                           string        `json:"datasetKey"`
	DatasetName                          string        `json:"datasetName"`
	DateIdentified                       interface{}         `json:"dateIdentified"`
	Day                                  int           `json:"day"`
	DecimalLatitude                      float64       `json:"decimalLatitude"`
	DecimalLongitude                     float64       `json:"decimalLongitude"`
	EventDate                            Time         `json:"eventDate"`
	EventTime                            string        `json:"eventTime"`
	Extensions                           struct{}      `json:"extensions"`
	Facts                                []interface{} `json:"facts"`
	Family                               string        `json:"family"`
	FamilyKey                            int           `json:"familyKey"`
	GbifID                               string        `json:"gbifID"`
	GenericName                          string        `json:"genericName"`
	Genus                                string        `json:"genus"`
	GenusKey                             int              `json:"genusKey"`
	GeodeticDatum                        string           `json:"geodeticDatum"`
	HTTP___unknown_org_occurrenceDetails string           `json:"http://unknown.org/occurrenceDetails"`
	IdentificationID                     string           `json:"identificationID"`
	Identifier                           string           `json:"identifier"`
	Identifiers                          []interface{}    `json:"identifiers"`
	InstitutionCode                      string           `json:"institutionCode"`
	Issues                               OccurrenceIssues `json:"issues"`
	Key                                  int              `json:"key"`
	Kingdom                              string           `json:"kingdom"`
	KingdomKey                           int              `json:"kingdomKey"`
	LastCrawled                          interface{}      `json:"lastCrawled"`
	LastInterpreted                      interface{}      `json:"lastInterpreted"`
	LastParsed                           interface{}      `json:"lastParsed"`
	Media                                []Media          `json:"media"`
	Modified                             interface{}           `json:"modified"`
	Month                                int           `json:"month"`
	OccurrenceID                         string        `json:"occurrenceID"`
	OccurrenceRemarks                    string        `json:"occurrenceRemarks"`
	Order                                string        `json:"order"`
	OrderKey                             int           `json:"orderKey"`
	Phylum                               string        `json:"phylum"`
	PhylumKey                            int           `json:"phylumKey"`
	Protocol                             string        `json:"protocol"`
	PublishingCountry                    string        `json:"publishingCountry"`
	PublishingOrgKey                     string        `json:"publishingOrgKey"`
	RecordedBy                           string        `json:"recordedBy"`
	References                           string        `json:"references"`
	Relations                            []interface{} `json:"relations"`
	Rights                               string        `json:"rights"`
	RightsHolder                         string        `json:"rightsHolder"`
	OccurrenceStatus                 string        `json:"occurrenceStatus"`
	ScientificName                       string        `json:"scientificName"`
	Species                              string        `json:"species"`
	SpeciesKey                           int           `json:"speciesKey"`
	SpecificEpithet                      string        `json:"specificEpithet"`
	TaxonID                              string        `json:"taxonID"`
	TaxonKey                             int           `json:"taxonKey"`
	TaxonRank                            string        `json:"taxonRank"`
	VerbatimEventDate                    string        `json:"verbatimEventDate"`
	Year                                 int           `json:"year"`
}

type response struct {
	page
	Results []Occurrence `json:"results"`
}

// OccurrenceCount returns a full search across all occurrences. Results are ordered by relevance.
func  Occurrences(q OccurrenceSearchQuery) ([]Occurrence, error) {

	var syncList struct{
		sync.Mutex
		list []Occurrence
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

var throttle = time.Tick(time.Second / 10)
func (s *OccurrenceSearchQuery) request(offset int) (response, error) {
	<-throttle
	client := pester.New()
	client.Concurrency = 1
	client.MaxRetries = 5
	client.Backoff = pester.ExponentialJitterBackoff
	client.KeepLog = true

	resp, err := client.Get(s.url(offset))
	if err != nil {
		return response{}, errors.Wrap(err, "could not get http response")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return response{}, errors.Wrapf(errors.New(resp.Status), "StatusCode: %d; URL: %s", resp.StatusCode, s.url(offset))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response{}, errors.Wrapf(err, "could not read http response body: %d", resp.StatusCode)
	}

	var res response
	if err := json.Unmarshal(body, &res); err != nil {
		return response{}, errors.Wrap(err, "could not unmarshal http response")
	}

	return res, nil
}

type OccurrenceIssue string

type OccurrenceIssues []OccurrenceIssue

func (Ω OccurrenceIssues) Intersects(query OccurrenceIssues) bool {
	for i := range Ω {
		if query.HasIssue(Ω[i]) {
			return true
		}
	}
	return false
}

func (Ω OccurrenceIssues) HasIssue(query OccurrenceIssue) bool {
	for i := range Ω {
		if Ω[i] == query {
			return true
		}
	}
	return false
}


const OCCURRENCE_ISSUE_BASIS_OF_RECORD_INVALID = OccurrenceIssue("BASIS_OF_RECORD_INVALID")
//The given basis of record is impossible to interpret or seriously different from the recommended vocabulary.
const OCCURRENCE_ISSUE_CONTINENT_COUNTRY_MISMATCH = OccurrenceIssue("CONTINENT_COUNTRY_MISMATCH")
//The interpreted continent and country do not match up.
const OCCURRENCE_ISSUE_CONTINENT_DERIVED_FROM_COORDINATES = OccurrenceIssue("DERIVED_FROM_COORDINATES")
//The interpreted continent is based on the coordinates, not the verbatim string information.
const OCCURRENCE_ISSUE_CONTINENT_INVALID = OccurrenceIssue("CONTINENT_INVALID")
//Uninterpretable continent values found.
const OCCURRENCE_ISSUE_COORDINATE_INVALID = OccurrenceIssue("COORDINATE_INVALID")
//Coordinate value given in some form but GBIF is unable to interpret it.
const OCCURRENCE_ISSUE_COORDINATE_OUT_OF_RANGE = OccurrenceIssue("COORDINATE_OUT_OF_RANGE")
//Coordinate has invalid lat/lon values out of their decimal max range.
const OCCURRENCE_ISSUE_COORDINATE_PRECISION_INVALID = OccurrenceIssue("COORDINATE_PRECISION_INVALID")
//Indicates an invalid or very unlikely coordinatePrecision
const OCCURRENCE_ISSUE_COORDINATE_REPROJECTED = OccurrenceIssue("COORDINATE_REPROJECTED")
//The original coordinate was successfully reprojected from a different geodetic datum to WGS84.
const OCCURRENCE_ISSUE_COORDINATE_REPROJECTION_FAILED = OccurrenceIssue("COORDINATE_REPROJECTION_FAILED")
//The given decimal latitude and longitude could not be reprojected to WGS84 based on the provided datum.
const OCCURRENCE_ISSUE_COORDINATE_REPROJECTION_SUSPICIOUS = OccurrenceIssue("COORDINATE_REPROJECTION_SUSPICIOUS")
//Indicates successful coordinate reprojection according to provided datum, but which results in a datum shift larger than 0.1 decimal degrees.
const OCCURRENCE_ISSUE_COORDINATE_ROUNDED = OccurrenceIssue("COORDINATE_ROUNDED")
//Original coordinate modified by rounding to 5 decimals.
const OCCURRENCE_ISSUE_COORDINATE_UNCERTAINTY_METERS_INVALID = OccurrenceIssue("COORDINATE_UNCERTAINTY_METERS_INVALID")
//Indicates an invalid or very unlikely dwc:uncertaintyInMeters.
const OCCURRENCE_ISSUE_COUNTRY_COORDINATE_MISMATCH = OccurrenceIssue("COUNTRY_COORDINATE_MISMATCH")
//The interpreted occurrence coordinates fall outside of the indicated country.
const OCCURRENCE_ISSUE_COUNTRY_DERIVED_FROM_COORDINATES = OccurrenceIssue("COUNTRY_DERIVED_FROM_COORDINATES")
//The interpreted country is based on the coordinates, not the verbatim string information.
const OCCURRENCE_ISSUE_COUNTRY_INVALID = OccurrenceIssue("COUNTRY_INVALID")
//Uninterpretable country values found.
const OCCURRENCE_ISSUE_COUNTRY_MISMATCH = OccurrenceIssue("COUNTRY_MISMATCH")
//Interpreted country for dwc:country and dwc:countryCode contradict each other.
const OCCURRENCE_ISSUE_DEPTH_MIN_MAX_SWAPPED = OccurrenceIssue("DEPTH_MIN_MAX_SWAPPED")
//Set if supplied min>max
const OCCURRENCE_ISSUE_DEPTH_NON_NUMERIC = OccurrenceIssue("DEPTH_NON_NUMERIC")
//Set if depth is a non numeric value
const OCCURRENCE_ISSUE_DEPTH_NOT_METRIC = OccurrenceIssue("DEPTH_NOT_METRIC")
//Set if supplied depth is not given in the metric system, for example using feet instead of meters
const OCCURRENCE_ISSUE_DEPTH_UNLIKELY = OccurrenceIssue("DEPTH_UNLIKELY")
//Set if depth is larger than 11.000m or negative.
const OCCURRENCE_ISSUE_ELEVATION_MIN_MAX_SWAPPED = OccurrenceIssue("ELEVATION_MIN_MAX_SWAPPED")
//Set if supplied min > max elevation
const OCCURRENCE_ISSUE_ELEVATION_NON_NUMERIC = OccurrenceIssue("ELEVATION_NON_NUMERIC")
//Set if elevation is a non numeric value
const OCCURRENCE_ISSUE_ELEVATION_NOT_METRIC = OccurrenceIssue("ELEVATION_NOT_METRIC")
//Set if supplied elevation is not given in the metric system, for example using feet instead of meters
const OCCURRENCE_ISSUE_ELEVATION_UNLIKELY = OccurrenceIssue("ELEVATION_UNLIKELY")
//Set if elevation is above the troposphere (17km) or below 11km (Mariana Trench).
const OCCURRENCE_ISSUE_GEODETIC_DATUM_ASSUMED_WGS84 = OccurrenceIssue("GEODETIC_DATUM_ASSUMED_WGS84")
//Indicating that the interpreted coordinates assume they are based on WGS84 datum as the datum was either not indicated or interpretable.
const OCCURRENCE_ISSUE_GEODETIC_DATUM_INVALID = OccurrenceIssue("GEODETIC_DATUM_INVALID")
//The geodetic datum given could not be interpreted.
const OCCURRENCE_ISSUE_IDENTIFIED_DATE_INVALID = OccurrenceIssue("IDENTIFIED_DATE_INVALID")
//The date given for dwc:dateIdentified is invalid and cant be interpreted at all.
const OCCURRENCE_ISSUE_IDENTIFIED_DATE_UNLIKELY = OccurrenceIssue("IDENTIFIED_DATE_UNLIKELY")
//The date given for dwc:dateIdentified is in the future or before Linnean times (1700).
const OCCURRENCE_ISSUE_INDIVIDUAL_COUNT_INVALID = OccurrenceIssue("INDIVIDUAL_COUNT_INVALID")
//Individual count value not parsable into an integer.
const OCCURRENCE_ISSUE_INTERPRETATION_ERROR = OccurrenceIssue("INTERPRETATION_ERROR")
//An error occurred during interpretation, leaving the record interpretation incomplete.
const OCCURRENCE_ISSUE_MODIFIED_DATE_INVALID = OccurrenceIssue("MODIFIED_DATE_INVALID")
//A (partial) invalid date is given for dc:modified, such as a non existing date, invalid zero month, etc.
const OCCURRENCE_ISSUE_MODIFIED_DATE_UNLIKELY = OccurrenceIssue("MODIFIED_DATE_UNLIKELY")
//The date given for dc:modified is in the future or predates unix time (1970).
const OCCURRENCE_ISSUE_MULTIMEDIA_DATE_INVALID = OccurrenceIssue("MULTIMEDIA_DATE_INVALID")
//An invalid date is given for dc:created of a multimedia object.
const OCCURRENCE_ISSUE_MULTIMEDIA_URI_INVALID = OccurrenceIssue("MULTIMEDIA_URI_INVALID")
//An invalid uri is given for a multimedia object.
const OCCURRENCE_ISSUE_PRESUMED_NEGATED_LATITUDE = OccurrenceIssue("PRESUMED_NEGATED_LATITUDE")
//Latitude appears to be negated, e.g.
const OCCURRENCE_ISSUE_PRESUMED_NEGATED_LONGITUDE = OccurrenceIssue("PRESUMED_NEGATED_LONGITUDE")
//Longitude appears to be negated, e.g.
const OCCURRENCE_ISSUE_PRESUMED_SWAPPED_COORDINATE = OccurrenceIssue("PRESUMED_SWAPPED_COORDINATE")
//Latitude and longitude appear to be swapped.
const OCCURRENCE_ISSUE_RECORDED_DATE_INVALID = OccurrenceIssue("RECORDED_DATE_INVALID")
//A (partial) invalid date is given, such as a non existing date, invalid zero month, etc.
const OCCURRENCE_ISSUE_RECORDED_DATE_MISMATCH = OccurrenceIssue("RECORDED_DATE_MISMATCH")
//The recording date specified as the eventDate string and the individual year, month, day are contradicting.
const OCCURRENCE_ISSUE_RECORDED_DATE_UNLIKELY = OccurrenceIssue("RECORDED_DATE_UNLIKELY")
//The recording date is highly unlikely, falling either into the future or represents a very old date before 1600 that predates modern taxonomy.
const OCCURRENCE_ISSUE_REFERENCES_URI_INVALID = OccurrenceIssue("REFERENCES_URI_INVALID")
//An invalid uri is given for dc:references.
const OCCURRENCE_ISSUE_TAXON_MATCH_FUZZY = OccurrenceIssue("TAXON_MATCH_FUZZY")
//Matching to the taxonomic backbone can only be done using a fuzzy, non exact match.
const OCCURRENCE_ISSUE_TAXON_MATCH_HIGHERRANK = OccurrenceIssue("TAXON_MATCH_HIGHERRANK")
//Matching to the taxonomic backbone can only be done on a higher rank and not the scientific name.
const OCCURRENCE_ISSUE_TAXON_MATCH_NONE = OccurrenceIssue("TAXON_MATCH_NONE")
//Matching to the taxonomic backbone cannot be done cause there was no match at all or several matches with too little information to keep them apart (homonyms).
const OCCURRENCE_ISSUE_TYPE_STATUS_INVALID = OccurrenceIssue("TYPE_STATUS_INVALID")
//The given type status is impossible to interpret or seriously different from the recommended vocabulary.
const OCCURRENCE_ISSUE_ZERO_COORDINATE = OccurrenceIssue("ZERO_COORDINATE")
//Coordinate is the exact 0/0 coordinate, often indicating a bad null coordinate.
