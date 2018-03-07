package inaturalist

import (
	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/utils"
	"context"
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"google.golang.org/api/iterator"
	"strconv"
	"strings"
	"time"
)

type occurrence struct {
	OutOfRange        bool          `json:"out_of_range"`
	QualityGrade      string        `json:"quality_grade"`
	TimeObservedAt    time.Time     `json:"time_observed_at"` // 2012-03-18T17:31:53-04:00
	Annotations       []interface{} `json:"annotations"`
	UUID              string        `json:"uuid"`
	Photos            []interface{} `json:"photos"`
	ObservedOnDetails struct {
		Date  string `json:"date"`
		Week  int    `json:"week"`
		Month int    `json:"month"`
		Hour  int    `json:"hour"`
		Year  int    `json:"year"`
		Day   int    `json:"day"`
	} `json:"observed_on_details"`
	ID                       int  `json:"id"`
	CachedVotesTotal         int  `json:"cached_votes_total"`
	IdentificationsMostAgree bool `json:"identifications_most_agree"`
	CreatedAtDetails         struct {
		Date  string `json:"date"`
		Week  int    `json:"week"`
		Month int    `json:"month"`
		Hour  int    `json:"hour"`
		Year  int    `json:"year"`
		Day   int    `json:"day"`
	} `json:"created_at_details"`
	SpeciesGuess                string        `json:"species_guess"`
	IdentificationsMostDisagree bool          `json:"identifications_most_disagree"`
	Tags                        []interface{} `json:"tags"`
	PositionalAccuracy          int           `json:"positional_accuracy"`
	CommentsCount               int           `json:"comments_count"`
	SiteID                      int           `json:"site_id"`
	CreatedTimeZone             string        `json:"created_time_zone"`
	IDPlease                    bool          `json:"id_please"`
	LicenseCode                 string        `json:"license_code"`
	ObservedTimeZone            string        `json:"observed_time_zone"`
	QualityMetrics              []interface{} `json:"quality_metrics"`
	PublicPositionalAccuracy    int           `json:"public_positional_accuracy"`
	ReviewedBy                  []int         `json:"reviewed_by"`
	OauthApplicationID          int           `json:"oauth_application_id"`
	Flags                       []interface{} `json:"flags"`
	CreatedAt                   string        `json:"created_at"`
	Description                 string        `json:"description"`
	TimeZoneOffset              string        `json:"time_zone_offset"`
	ProjectIdsWithCuratorID     []int         `json:"project_ids_with_curator_id"`
	ObservedOn                  string        `json:"observed_on"`
	ObservedOnString            string        `json:"observed_on_string"`
	UpdatedAt                   string        `json:"updated_at"`
	Sounds                      []interface{} `json:"sounds"`
	PlaceIds                    []int         `json:"place_ids"`
	Captive                     bool          `json:"captive"`
	Taxon                       taxon         `json:"taxon"`
	Outlinks                    []interface{} `json:"outlinks"`
	FavesCount                  int           `json:"faves_count"`
	Ofvs                        []interface{} `json:"ofvs"`
	NumIdentificationAgreements int           `json:"num_identification_agreements"`
	Preferences                 struct {
		PrefersCommunityTaxon interface{} `json:"prefers_community_taxon"`
	} `json:"preferences"`
	Comments         []interface{}   `json:"comments"`
	MapScale         int             `json:"map_scale"`
	URI              string          `json:"uri"`
	ProjectIds       []int           `json:"project_ids"`
	Identifications  identifications `json:"identifications"`
	CommunityTaxonID interface{}     `json:"community_taxon_id"`
	Geojson          *struct {
		Coordinates []string `json:"coordinates"`
		Type        string   `json:"type"`
	} `json:"geojson"`
	OwnersIdentificationFromVision bool                 `json:"owners_identification_from_vision"`
	IdentificationsCount           int                  `json:"identifications_count"`
	Obscured                       bool                 `json:"obscured"`
	ProjectObservations            []projectObservation `json:"project_observations"`
	NumIdentificationDisagreements int                  `json:"num_identification_disagreements"`
	ObservationPhotos              []interface{}        `json:"observation_photos"`
	Geoprivacy                     interface{}          `json:"geoprivacy"`
	Location                       string               `json:"location"`
	Votes                          []interface{}        `json:"votes"`
	User                           user                 `json:"user"`
	Mappable                       bool                 `json:"mappable"`
	IdentificationsSomeAgree       bool                 `json:"identifications_some_agree"`
	ProjectIdsWithoutCuratorID     []int                `json:"project_ids_without_curator_id"`
	PlaceGuess                     string               `json:"place_guess"`
	Faves                          []interface{}        `json:"faves"`
	NonOwnerIds                    []interface{}        `json:"non_owner_ids"`
	Application                    struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		URL  string `json:"url"`
		Icon string `json:"icon"`
	} `json:"application"`
}

// Lat returns the latitude.
func (Ω *occurrence) Lat() (float64, error) {
	lat, err := strconv.ParseFloat(Ω.Geojson.Coordinates[1], 64)
	if err != nil {
		return 0, errors.Wrap(err, "Could not parse latitude from INaturalist occurrence")
	}
	return lat, nil
}

// Lng returns the longitude.
func (Ω *occurrence) Lng() (float64, error) {
	lng, err := strconv.ParseFloat(Ω.Geojson.Coordinates[0], 64)
	if err != nil {
		return 0, errors.Wrap(err, "Could not parse longitude from INaturalist occurrence")
	}
	return lng, nil
}

// DateString converts the date to YYYYMMDD.
func (Ω *occurrence) DateString() string {
	return strings.Replace(Ω.ObservedOnDetails.Date, "-", "", -1)
}
func (Ω *occurrence) CoordinatesEstimated() bool {
	// Rounded to 5 decimal place. Not what I expected.
	// isEstimated := s.Issues.hasIssue(ogbif.OCCURRENCE_ISSUE_COORDINATE_ROUNDED)
	return false
}

// SourceOccurrenceID upholds the occurrence interface.
func (Ω *occurrence) SourceOccurrenceID() string {
	return strconv.Itoa(Ω.ID)
}

type projectObservation struct {
	Preferences struct {
		AllowsCuratorCoordinateAccess bool `json:"allows_curator_coordinate_access"`
	} `json:"preferences"`
	UserID    interface{} `json:"user_id"`
	ID        int         `json:"id"`
	UUID      string      `json:"uuid"`
	Project   project     `json:"project"`
	ProjectID int         `json:"project_id"`
}

type project struct {
	ID                       int           `json:"id"`
	ProjectObservationFields []interface{} `json:"project_observation_fields"`
	Icon                     string        `json:"icon"`
	Description              string        `json:"description"`
	Location                 string        `json:"location"`
	Title                    string        `json:"title"`
	Slug                     string        `json:"slug"`
	Latitude                 string        `json:"latitude"`
	Longitude                string        `json:"longitude"`
}

type user struct {
	ID          int    `json:"id"`
	Login       string `json:"login"`
	Preferences struct {
	} `json:"preferences"`
	LoginAutocomplete    string        `json:"login_autocomplete"`
	Name                 string        `json:"name"`
	NameAutocomplete     string        `json:"name_autocomplete"`
	Icon                 string        `json:"icon"`
	ObservationsCount    int           `json:"observations_count"`
	IdentificationsCount int           `json:"identifications_count"`
	JournalPostsCount    int           `json:"journal_posts_count"`
	ActivityCount        int           `json:"activity_count"`
	Roles                []interface{} `json:"roles"`
	SiteID               interface{}   `json:"site_id"`
	IconURL              string        `json:"icon_url"`
}

type identifications []*identification

func (Ω identifications) LatestIdentification() time.Time {
	latest := time.Time{}
	for _, identification := range Ω {
		if identification.CreatedAt.After(latest) {
			latest = identification.CreatedAt
		}
	}
	return latest
}

type identification struct {
	Disagreement     interface{}   `json:"disagreement"`
	Flags            []interface{} `json:"flags"`
	CreatedAt        time.Time     `json:"created_at"`
	TaxonID          int           `json:"taxon_id"`
	Body             interface{}   `json:"body"`
	OwnObservation   bool          `json:"own_observation"`
	UUID             string        `json:"uuid"`
	TaxonChange      interface{}   `json:"taxon_change"`
	Vision           bool          `json:"vision"`
	Current          bool          `json:"current"`
	ID               int           `json:"id"`
	CreatedAtDetails struct {
		Date  string `json:"date"`
		Week  int    `json:"week"`
		Month int    `json:"month"`
		Hour  int    `json:"hour"`
		Year  int    `json:"year"`
		Day   int    `json:"day"`
	} `json:"created_at_details"`
	Category                   string      `json:"category"`
	User                       user        `json:"user"`
	PreviousObservationTaxonID interface{} `json:"previous_observation_taxon_id"`
	Taxon                      taxon       `json:"taxon"`
}

// FetchOccurrences returns returns a slice of OccurrenceProviders.
func FetchOccurrences(_ context.Context, targetID datasources.TargetID, since *time.Time) ([]*occurrence, error) {

	if !taxonIDFromTargetID(targetID).Valid() {
		return nil, errors.New("Invalid taxonID")
	}

	res := []*occurrence{}
	page := 1
	for {
		list, err := fetchOccurrences(page, taxonIDFromTargetID(targetID), since)
		if err != nil && err != iterator.Done {
			return nil, err
		}
		res = append(res, list...)
		if err != nil && err == iterator.Done {
			break
		}
		page++
	}

	return res, nil
}

var throttle = time.NewTicker(time.Second / 20)

func fetchOccurrences(page int, txnID taxonID, since *time.Time) ([]*occurrence, error) {
	var response struct {
		TotalResults int           `json:"total_results"`
		Page         int           `json:"page"`
		PerPage      int           `json:"per_page"`
		Results      []*occurrence `json:"results"`
	}

	u := "https://api.inaturalist.org/v1/observations?place_id=97394&quality_grade=research&captive=false&per_page=200&geoprivacy=open"
	u += fmt.Sprintf("&taxon_id=%d", txnID)
	u += fmt.Sprintf("&page=%d", page)
	if since != nil && !since.IsZero() {
		u += fmt.Sprintf("&updated_since=%s", since.Format("2006-01-02"))
	}

	<-throttle.C

	if err := utils.RequestJSON(u, &response); err != nil {
		return nil, err
	}

	res := []*occurrence{}

	for _, inatOccurrence := range response.Results {

		// https://www.inaturalist.org/pages/help#quality
		// All of these are covered in query, but just to be safe ...
		if inatOccurrence.QualityGrade != "research" ||
			inatOccurrence.Captive ||
			inatOccurrence.Taxon.ID != txnID { // Ignore descendents
			continue
		}

		if inatOccurrence.Geojson == nil {
			return nil, errors.Newf("Invalid Geojson [%d]. Should have been solved by changing privacy setting", inatOccurrence.ID)
		}

		res = append(res, inatOccurrence)
	}

	if (response.Page * response.PerPage) < response.TotalResults {
		return res, nil
	}

	return res, iterator.Done
}
