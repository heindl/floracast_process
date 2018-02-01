package gbif

import (
	ogbif "github.com/heindl/gbif"
	"bitbucket.org/heindl/taxa/store"
	"time"
	"bitbucket.org/heindl/taxa/occurrences"
	"github.com/dropbox/godropbox/errors"
	"fmt"
	"bitbucket.org/heindl/taxa/utils"
	"context"
	"bitbucket.org/heindl/taxa/geofeatures"
)

func FetchOccurrences(cxt context.Context, targetID store.DataSourceTargetID, since *time.Time) (occurrences.Occurrences, error) {

	taxonID := TaxonIDFromTargetID(targetID)

	if !taxonID.Valid() {
		return nil, errors.New("Invalid TaxonID")
	}

	lastInterpreted := ""
	if since != nil && !since.IsZero() {
		lastInterpreted = since.Format("20060102")
	}

	apiList, err := ogbif.Occurrences(ogbif.OccurrenceSearchQuery{
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
		if coordinateIssueIsUnacceptable(gbifO.Issues...) {
			continue
		}

		if coordinateIssueIsUncertain(gbifO.Issues...) {
			fmt.Println("coordinateIssueIsUncertain", utils.JsonOrSpew(gbifO))
			continue
		}

		o, err := occurrences.NewOccurrence(store.DataSourceTypeGBIF, targetID, gbifO.GbifID)
		if err != nil {
			return nil, err
		}
		err = o.SetGeospatial(gbifO.DecimalLatitude, gbifO.DecimalLongitude, gbifO.EventDate.Time)
		if err != nil && utils.ContainsError(err, geofeatures.ErrInvalidCoordinate) {
			continue
		}
		if err != nil && utils.ContainsError(err, occurrences.ErrInvalidDate) {
			continue
		}
		if err != nil {
			return nil, err
		}
		res = append(res, o)
	}

	return res, nil

}

func coordinateIssueIsUnacceptable(issues ...ogbif.OccurrenceIssue) bool {

	hasInvalidGeodeticDatum := false
	geodeticAssumed := false
	for _, issue := range issues {
		if issue == ogbif.OCCURRENCE_ISSUE_GEODETIC_DATUM_INVALID {
			hasInvalidGeodeticDatum = true
		}
		if issue == ogbif.OCCURRENCE_ISSUE_GEODETIC_DATUM_ASSUMED_WGS84 {
			geodeticAssumed = true
		}
	}

	if hasInvalidGeodeticDatum && !geodeticAssumed {
		return true
	}

	for _, issue := range issues {
		for _, invalid := range []ogbif.OccurrenceIssue {
			ogbif.OCCURRENCE_ISSUE_BASIS_OF_RECORD_INVALID,
			ogbif.OCCURRENCE_ISSUE_COORDINATE_INVALID,
			ogbif.OCCURRENCE_ISSUE_COORDINATE_OUT_OF_RANGE,
			ogbif.OCCURRENCE_ISSUE_COORDINATE_REPROJECTION_FAILED,
			ogbif.OCCURRENCE_ISSUE_ZERO_COORDINATE,
			ogbif.OCCURRENCE_ISSUE_RECORDED_DATE_INVALID,
			ogbif.OCCURRENCE_ISSUE_RECORDED_DATE_UNLIKELY,
			//ogbif.OCCURRENCE_ISSUE_COORDINATE_PRECISION_INVALID,
			//ogbif.OCCURRENCE_ISSUE_COORDINATE_UNCERTAINTY_METERS_INVALID,
			ogbif.OCCURRENCE_ISSUE_COORDINATE_REPROJECTION_SUSPICIOUS,
		} {
			if invalid == issue {
				return true
			}
		}
	}
	return false

}

func coordinateIssueIsUncertain(issues ...ogbif.OccurrenceIssue) bool {
	for _, issue := range issues {
		for _, invalid := range []ogbif.OccurrenceIssue{
			ogbif.OCCURRENCE_ISSUE_PRESUMED_SWAPPED_COORDINATE,
		} {
			if invalid == issue {
				return true
			}
		}
	}
	return false

}