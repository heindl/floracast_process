package inaturalist

import (
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"time"
	"bitbucket.org/heindl/taxa/utils"
	"strconv"
	"context"
	"bitbucket.org/heindl/taxa/datasources"
	"strings"
)

type Occurrence struct {
		OutOfRange        bool          `json:"out_of_range"`
		QualityGrade      string        `json:"quality_grade"`
		TimeObservedAt    time.Time       `json:"time_observed_at"` // 2012-03-18T17:31:53-04:00
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
		Taxon                       Taxon         `json:"taxon"`
		Outlinks                    []interface{} `json:"outlinks"`
		FavesCount                  int           `json:"faves_count"`
		Ofvs                        []interface{} `json:"ofvs"`
		NumIdentificationAgreements int           `json:"num_identification_agreements"`
		Preferences                 struct {
			PrefersCommunityTaxon interface{} `json:"prefers_community_taxon"`
		} `json:"preferences"`
		Comments        []interface{} `json:"comments"`
		MapScale        int           `json:"map_scale"`
		URI             string        `json:"uri"`
		ProjectIds      []int         `json:"project_ids"`
		Identifications Identifications  `json:"identifications"`
		CommunityTaxonID interface{} `json:"community_taxon_id"`
		Geojson          *struct {
			Coordinates []string `json:"coordinates"`
			Type        string   `json:"type"`
		} `json:"geojson"`
		OwnersIdentificationFromVision bool `json:"owners_identification_from_vision"`
		IdentificationsCount           int  `json:"identifications_count"`
		Obscured                       bool `json:"obscured"`
		ProjectObservations            []ProjectObservation `json:"project_observations"`
		NumIdentificationDisagreements int           `json:"num_identification_disagreements"`
		ObservationPhotos              []interface{} `json:"observation_photos"`
		Geoprivacy                     interface{}   `json:"geoprivacy"`
		Location                       string        `json:"location"`
		Votes                          []interface{} `json:"votes"`
		User                           User `json:"user"`
		Mappable                   bool          `json:"mappable"`
		IdentificationsSomeAgree   bool          `json:"identifications_some_agree"`
		ProjectIdsWithoutCuratorID []int         `json:"project_ids_without_curator_id"`
		PlaceGuess                 string        `json:"place_guess"`
		Faves                      []interface{} `json:"faves"`
		NonOwnerIds                []interface{} `json:"non_owner_ids"`
		Application                struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			URL  string `json:"url"`
			Icon string `json:"icon"`
		} `json:"application"`
}

func (Ω *Occurrence) Lat() (float64, error) {
	lat, err := strconv.ParseFloat(Ω.Geojson.Coordinates[1], 64)
	if err != nil {
		return 0, errors.Wrap(err, "Could not parse latitude from INaturalist occurrence")
	}
	return lat, nil
}
func (Ω *Occurrence) Lng() (float64, error) {
	lng, err := strconv.ParseFloat(Ω.Geojson.Coordinates[0], 64)
	if err != nil {
		return 0, errors.Wrap(err, "Could not parse longitude from INaturalist occurrence")
	}
	return lng, nil
}
func (Ω *Occurrence) DateString() string {
	return strings.Replace(Ω.ObservedOnDetails.Date, "-", "", -1)
}
func (Ω *Occurrence) CoordinatesEstimated() bool {
	// Rounded to 5 decimal place. Not what I expected.
	// isEstimated := s.Issues.HasIssue(ogbif.OCCURRENCE_ISSUE_COORDINATE_ROUNDED)
	return false
}
func (Ω *Occurrence) SourceOccurrenceID() string {
	return strconv.Itoa(Ω.ID)
}

type ProjectObservation struct {
	Preferences struct {
		AllowsCuratorCoordinateAccess bool `json:"allows_curator_coordinate_access"`
	} `json:"preferences"`
	UserID  interface{} `json:"user_id"`
	ID      int         `json:"id"`
	UUID    string      `json:"uuid"`
	Project Project `json:"project"`
	ProjectID int `json:"project_id"`
}

type Project struct {
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

type User struct {
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

type Identifications []*Identification

func (Ω Identifications) LatestIdentification() time.Time {
	latest := time.Time{}
	for _, identification := range Ω {
		if identification.CreatedAt.After(latest) {
			latest = identification.CreatedAt
		}
	}
	return latest
}

type Identification struct {
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
	User                       User        `json:"user"`
	PreviousObservationTaxonID interface{} `json:"previous_observation_taxon_id"`
	Taxon                      Taxon       `json:"taxon"`
}

var throttle = time.Tick(time.Second / 20)
func FetchOccurrences(cxt context.Context, targetID datasources.TargetID, since *time.Time) ([]*Occurrence, error) {

	taxonID := TaxonIDFromTargetID(targetID)

	if !taxonID.Valid() {
		return nil, errors.New("Invalid TaxonID")
	}

	output := []*Occurrence{}

	page := 1

	for {

		var response struct {
			TotalResults int `json:"total_results"`
			Page         int `json:"page"`
			PerPage      int `json:"per_page"`
			Results []*Occurrence `json:"results"`
		}

		u := "https://api.inaturalist.org/v1/observations?place_id=97394&quality_grade=research&captive=false&per_page=200&geoprivacy=open"
		u += fmt.Sprintf("&taxon_id=%d", taxonID)
		u += fmt.Sprintf("&page=%d", page)
		if since != nil && !since.IsZero() {
			u += fmt.Sprintf("&updated_since=%s", since.Format("2006-01-02"))
		}

		<-throttle

		if err := utils.RequestJSON(u, &response); err != nil {
			return nil, err
		}

		for _, inatOccurrence := range response.Results {

			// Covered in query, but just to be safe ...
			// https://www.inaturalist.org/pages/help#quality
			if inatOccurrence.QualityGrade != "research" {
				fmt.Println("not research grade", inatOccurrence.QualityGrade)
				continue
			}

			// Covered in query, but just to be safe ...
			if inatOccurrence.Captive {
				fmt.Println("is captive")
				continue
			}

			// Covered in query, but just to be safe ...
			if inatOccurrence.Taxon.ID != taxonID {
				// Ignore descendents
				//fmt.Println("mismatched taxon id", inatOccurrence.Taxon.ID, taxonID)
				continue
			}

			if inatOccurrence.Geojson == nil {
				return nil, errors.Newf("Invalid Geojson [%d]. Should have been solved by changing privacy setting", inatOccurrence.ID)
			}

			output = append(output, inatOccurrence)
		}

		if (response.Page * response.PerPage) < response.TotalResults {
			page += 1
			continue
		}

		break
	}

	return output, nil
}
