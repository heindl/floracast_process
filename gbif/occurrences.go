package gbif

import (
	"bitbucket.org/heindl/taxa/gbif/api"
	"bitbucket.org/heindl/taxa/store"
	"time"
	"bitbucket.org/heindl/taxa/occurrences"
	"github.com/dropbox/godropbox/errors"
	"fmt"
	"bitbucket.org/heindl/taxa/utils"
	"context"
	"bitbucket.org/heindl/taxa/geofeatures"
)

func FetchOccurrences(cxt context.Context, targetID store.DataSourceTargetID, since *time.Time) (*occurrences.Occurrences, error) {

	taxonID := TaxonIDFromTargetID(targetID)

	if !taxonID.Valid() {
		return nil, errors.New("Invalid TaxonID")
	}

	lastInterpreted := ""
	if since != nil && !since.IsZero() {
		lastInterpreted = since.Format("20060102")
	}

	apiList, err := api.Occurrences(api.OccurrenceSearchQuery{
		TaxonKey: int(taxonID),
		LastInterpreted: lastInterpreted,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not fetch occurrences [%d] from the gbif", taxonID)
	}

	res := occurrences.Occurrences{}

	fmt.Println("TOTAL OCCURRENCE LENGTH", len(apiList))

	for _, gbifO := range apiList {

		// Check geospatial issue.
		if coordinateIssueIsUnacceptable(gbifO.Issues) {
			continue
		}

		if coordinateIssueIsUncertain(gbifO.Issues) {
			fmt.Println("coordinateIssueIsUncertain", utils.JsonOrSpew(gbifO))
			continue
		}

		o, err := occurrences.NewOccurrence(store.DataSourceTypeGBIF, targetID, gbifO.GbifID)
		if err != nil {
			return nil, err
		}

		// Rounded to 5 decimal place. Not what I expected.
		// isEstimated := gbifO.Issues.HasIssue(ogbif.OCCURRENCE_ISSUE_COORDINATE_ROUNDED)

		err = o.SetGeospatial(gbifO.DecimalLatitude, gbifO.DecimalLongitude, gbifO.EventDate.Time.Format("20060102"), false)
		if err != nil && utils.ContainsError(err, geofeatures.ErrInvalidCoordinate) {
			continue
		}
		if err != nil && utils.ContainsError(err, occurrences.ErrInvalidDate) {
			continue
		}
		if err != nil {
			return nil, err
		}

		if err := res.Add(o); err != nil {
			return nil, err
		}
	}

	return &res, nil

}

func coordinateIssueIsUnacceptable(issues api.OccurrenceIssues) bool {

	if issues.HasIssue(api.OCCURRENCE_ISSUE_GEODETIC_DATUM_INVALID) &&
		!issues.HasIssue(api.OCCURRENCE_ISSUE_GEODETIC_DATUM_ASSUMED_WGS84) {
		return true
	}
	return issues.Intersects(api.OccurrenceIssues{
		api.OCCURRENCE_ISSUE_BASIS_OF_RECORD_INVALID,
		api.OCCURRENCE_ISSUE_COORDINATE_INVALID,
		api.OCCURRENCE_ISSUE_COORDINATE_OUT_OF_RANGE,
		api.OCCURRENCE_ISSUE_COORDINATE_REPROJECTION_FAILED,
		api.OCCURRENCE_ISSUE_ZERO_COORDINATE,
		api.OCCURRENCE_ISSUE_RECORDED_DATE_INVALID,
		api.OCCURRENCE_ISSUE_RECORDED_DATE_UNLIKELY,
		//ogbif.OCCURRENCE_ISSUE_COORDINATE_PRECISION_INVALID,
		//ogbif.OCCURRENCE_ISSUE_COORDINATE_UNCERTAINTY_METERS_INVALID,
		api.OCCURRENCE_ISSUE_COORDINATE_REPROJECTION_SUSPICIOUS,
	})
}

func coordinateIssueIsUncertain(issues api.OccurrenceIssues) bool {
	return issues.Intersects(api.OccurrenceIssues{
		api.OCCURRENCE_ISSUE_PRESUMED_SWAPPED_COORDINATE,
		api.OCCURRENCE_ISSUE_PRESUMED_NEGATED_LATITUDE,
		api.OCCURRENCE_ISSUE_PRESUMED_NEGATED_LONGITUDE,
		api.OCCURRENCE_ISSUE_INTERPRETATION_ERROR,
	})

}