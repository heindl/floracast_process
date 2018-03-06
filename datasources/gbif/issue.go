package gbif

type occurrenceIssue string

type occurrenceIssues []occurrenceIssue

func (Ω occurrenceIssues) intersects(query occurrenceIssues) bool {
	for i := range Ω {
		if query.hasIssue(Ω[i]) {
			return true
		}
	}
	return false
}

func (Ω occurrenceIssues) hasIssue(query occurrenceIssue) bool {
	for i := range Ω {
		if Ω[i] == query {
			return true
		}
	}
	return false
}

func (Ω occurrenceIssues) isUnacceptable() bool {

	if Ω.hasIssue(occurrenceIssueGeodeticDatumInvalid) &&
		!Ω.hasIssue(occurrenceIssueGeodeticDatumAssumedWGS84) {
		return true
	}
	return Ω.intersects(occurrenceIssues{
		occurrenceIssueBasisOfRecordInvalid,
		occurrenceIssueCoordinateInvalid,
		occurrenceIssueCoordinateOutOfRange,
		occurrenceIssueCoordinateReprojectionFailed,
		occurrenceIssueZeroCoordinate,
		occurrenceIssueRecordedDateInvalid,
		occurrenceIssueRecordedDateUnlikely,
		//occurrenceIssueCoordinatePrecisionInvalid,
		//occurrenceIssueCoordinateUncertaintyMetersInvalid,
		occurrenceIssueCoordinateReprojectionSuspicious,
	})
}

func (Ω occurrenceIssues) isUncertain() bool {
	return Ω.intersects(occurrenceIssues{
		occurrenceIssuePresumedNegatedLatitude,
		occurrenceIssuePresumedNegatedLongitude,
		occurrenceIssueInterpretationError,
	})
}

const (
	occurrenceIssueBasisOfRecordInvalid = occurrenceIssue("BASIS_OF_RECORD_INVALID")

	//The given basis of record is impossible to interpret or seriously different from the recommended vocabulary.
	//occurrenceIssueContinentCountryMismatch = occurrenceIssue("CONTINENT_COUNTRY_MISMATCH")

	//The interpreted continent and country do not match up.
	//occurrenceIssueContinentDerivedFromCoordinates = occurrenceIssue("DERIVED_FROM_COORDINATES")

	//The interpreted continent is based on the coordinates, not the verbatim string information.
	//occurrenceIssueContinentInvalid = occurrenceIssue("CONTINENT_INVALID")

	//Uninterpretable continent values found.
	occurrenceIssueCoordinateInvalid = occurrenceIssue("COORDINATE_INVALID")

	//Coordinate value given in some form but GBIF is unable to interpret it.
	occurrenceIssueCoordinateOutOfRange = occurrenceIssue("COORDINATE_OUT_OF_RANGE")

	//Coordinate has invalid lat/lon values out of their decimal max range.
	//occurrenceIssueCoordinatePrecisionInvalid = occurrenceIssue("COORDINATE_PRECISION_INVALID")

	//Indicates an invalid or very unlikely coordinatePrecision
	//occurrenceIssueCoordinateReprojected = occurrenceIssue("COORDINATE_REPROJECTED")

	//The original coordinate was successfully reprojected from a different geodetic datum to WGS84.
	occurrenceIssueCoordinateReprojectionFailed = occurrenceIssue("COORDINATE_REPROJECTION_FAILED")

	//The given decimal latitude and longitude could not be reprojected to WGS84 based on the provided datum.
	occurrenceIssueCoordinateReprojectionSuspicious = occurrenceIssue("COORDINATE_REPROJECTION_SUSPICIOUS")

	//Indicates successful coordinate reprojection according to provided datum, but which results in a datum shift larger than 0.1 decimal degrees.
	//occurrenceIssueCoordinateRounded = occurrenceIssue("COORDINATE_ROUNDED")

	//Original coordinate modified by rounding to 5 decimals.
	//occurrenceIssueCoordinateUncertaintyMetersInvalid = occurrenceIssue("COORDINATE_UNCERTAINTY_METERS_INVALID")

	//Indicates an invalid or very unlikely dwc:uncertaintyInMeters.
	//occurrenceIssueCountryCoordinateMismatch = occurrenceIssue("COUNTRY_COORDINATE_MISMATCH")

	//The interpreted occurrence coordinates fall outside of the indicated country.
	//occurrenceIssueCountryDerivedFromCoordinates = occurrenceIssue("COUNTRY_DERIVED_FROM_COORDINATES")

	//The interpreted country is based on the coordinates, not the verbatim string information.
	//occurrenceIssueCountryInvalid = occurrenceIssue("COUNTRY_INVALID")

	//Uninterpretable country values found.
	//occurrenceIssueCountryMismatch = occurrenceIssue("COUNTRY_MISMATCH")

	//Interpreted country for dwc:country and dwc:countryCode contradict each other.
	//occurrenceIssueDepthMinMaxSwapped = occurrenceIssue("DEPTH_MIN_MAX_SWAPPED")

	//Set if supplied min>max
	//occurrenceIssueDepthNonNumeric = occurrenceIssue("DEPTH_NON_NUMERIC")

	//Set if depth is a non numeric value
	//occurrenceIssueDepthNotMetric = occurrenceIssue("DEPTH_NOT_METRIC")

	//Set if supplied depth is not given in the metric system, for example using feet instead of meters
	//occurrenceIssueDepthUnlikely = occurrenceIssue("DEPTH_UNLIKELY")

	//Set if depth is larger than 11.000m or negative.
	//occurrenceIssueElevationMinMaxSwapped = occurrenceIssue("ELEVATION_MIN_MAX_SWAPPED")

	//Set if supplied min > max elevation
	//occurrenceIssueElevationNonNumeric = occurrenceIssue("ELEVATION_NON_NUMERIC")

	//Set if elevation is a non numeric value
	//occurrenceIssueElevationNotMetric = occurrenceIssue("ELEVATION_NOT_METRIC")

	//Set if supplied elevation is not given in the metric system, for example using feet instead of meters
	//occurrenceIssueElevationUnlikely = occurrenceIssue("ELEVATION_UNLIKELY")

	//Set if elevation is above the troposphere (17km) or below 11km (Mariana Trench).
	occurrenceIssueGeodeticDatumAssumedWGS84 = occurrenceIssue("GEODETIC_DATUM_ASSUMED_WGS84")

	//Indicating that the interpreted coordinates assume they are based on WGS84 datum as the datum was either not indicated or interpretable.
	occurrenceIssueGeodeticDatumInvalid = occurrenceIssue("GEODETIC_DATUM_INVALID")

	//The geodetic datum given could not be interpreted.
	//occurrenceIssueIdentifiedDateInvalid = occurrenceIssue("IDENTIFIED_DATE_INVALID")

	//The date given for dwc:dateIdentified is invalid and cant be interpreted at all.
	//occurrenceIssueIdentifiedDateUnlikely = occurrenceIssue("IDENTIFIED_DATE_UNLIKELY")

	//The date given for dwc:dateIdentified is in the future or before Linnean times (1700).
	//occurrenceIssueIndividualCountInvalid = occurrenceIssue("INDIVIDUAL_COUNT_INVALID")

	//Individual count value not parsable into an integer.
	occurrenceIssueInterpretationError = occurrenceIssue("INTERPRETATION_ERROR")

	//An error occurred during interpretation, leaving the record interpretation incomplete.
	//occurrenceIssueModifiedDateInvalid = occurrenceIssue("MODIFIED_DATE_INVALID")

	//A (partial) invalid date is given for dc:modified, such as a non existing date, invalid zero month, etc.
	//occurrenceIssueModifiedDateUnlikely = occurrenceIssue("MODIFIED_DATE_UNLIKELY")

	//The date given for dc:modified is in the future or predates unix time (1970).
	//occurrenceIssueMultimediaDateInvalid = occurrenceIssue("MULTIMEDIA_DATE_INVALID")

	//An invalid date is given for dc:created of a multimedia object.
	//occurrenceIssueMultimediaUriInvalid = occurrenceIssue("MULTIMEDIA_URI_INVALID")

	//An invalid uri is given for a multimedia object.
	occurrenceIssuePresumedNegatedLatitude = occurrenceIssue("PRESUMED_NEGATED_LATITUDE")

	//Latitude appears to be negated, e.g.
	occurrenceIssuePresumedNegatedLongitude = occurrenceIssue("PRESUMED_NEGATED_LONGITUDE")

	//Longitude appears to be negated, e.g.
	occurrenceIssuePresumedSwappedCoordinate = occurrenceIssue("PRESUMED_SWAPPED_COORDINATE")

	//Latitude and longitude appear to be swapped.
	occurrenceIssueRecordedDateInvalid = occurrenceIssue("RECORDED_DATE_INVALID")

	//A (partial) invalid date is given, such as a non existing date, invalid zero month, etc.
	//occurrenceIssueRecordedDateMismatch = occurrenceIssue("RECORDED_DATE_MISMATCH")

	//The recording date specified as the eventDate string and the individual year, month, day are contradicting.
	occurrenceIssueRecordedDateUnlikely = occurrenceIssue("RECORDED_DATE_UNLIKELY")

	//The recording date is highly unlikely, falling either into the future or represents a very old date before 1600 that predates modern taxonomy.
	//occurrenceIssueReferencesUriInvalid = occurrenceIssue("REFERENCES_URI_INVALID")

	//An invalid uri is given for dc:references.
	//occurrenceIssueTaxonMatchFuzzy = occurrenceIssue("TAXON_MATCH_FUZZY")

	//Matching to the taxonomic backbone can only be done using a fuzzy, non exact match.
	//occurrenceIssueTaxonMatchHigherrank = occurrenceIssue("TAXON_MATCH_HIGHERRANK")

	//Matching to the taxonomic backbone can only be done on a higher rank and not the scientific name.
	//occurrenceIssueTaxonMatchNone = occurrenceIssue("TAXON_MATCH_NONE")

	//Matching to the taxonomic backbone cannot be done cause there was no match at all or several matches with too little information to keep them apart (homonyms).
	//occurrenceIssueTypeStatusInvalid = occurrenceIssue("TYPE_STATUS_INVALID")

	//The given type status is impossible to interpret or seriously different from the recommended vocabulary.
	occurrenceIssueZeroCoordinate = occurrenceIssue("ZERO_COORDINATE")

	//Coordinate is the exact 0/0 coordinate, often indicating a bad null coordinate.

)
