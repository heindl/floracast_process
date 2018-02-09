package gbif

import (
	"bitbucket.org/heindl/taxa/datasources/gbif/api"
	"time"
	"github.com/dropbox/godropbox/errors"
	"fmt"
	"context"
	"bitbucket.org/heindl/taxa/datasources"
)

func FetchOccurrences(cxt context.Context, targetID datasources.TargetID, since *time.Time) ([]*api.Occurrence, error) {

	taxonID, err := TaxonIDFromTargetID(targetID)
	if err != nil {
		return nil, err
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

	res := []*api.Occurrence{}

	for _, gbifO := range apiList {

		// Check geospatial issue.
		if coordinateIssueIsUnacceptable(gbifO.Issues) {
			continue
		}

		if gbifO.Issues.HasIssue(api.OCCURRENCE_ISSUE_PRESUMED_SWAPPED_COORDINATE) {
			fmt.Println("WARNING:", api.OCCURRENCE_ISSUE_PRESUMED_SWAPPED_COORDINATE, fmt.Sprintf("Latitude/Longitude [%f, %f]", gbifO.DecimalLatitude, gbifO.DecimalLongitude))
			continue
			}

		if coordinateIssueIsUncertain(gbifO.Issues) {
			continue
		}

		res = append(res, gbifO)

	}

	return res, nil

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
		api.OCCURRENCE_ISSUE_PRESUMED_NEGATED_LATITUDE,
		api.OCCURRENCE_ISSUE_PRESUMED_NEGATED_LONGITUDE,
		api.OCCURRENCE_ISSUE_INTERPRETATION_ERROR,
	})

}