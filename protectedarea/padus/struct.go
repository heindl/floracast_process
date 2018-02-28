package padus

// https://gapanalysis.usgs.gov/padus/data/metadata/
type ProtectedArea struct {
	Access  PublicAccess `json:"Access"`
	DAccess string       `json:"d_Access"`
	// Documents the Source of Access domain assignment (e.g. State Data Steward or ‘GAP Default’ categorical assignment).
	AccessSource string `json:"Access_Src"`
	// "Aggregator Source" identifies the organization credited with data aggregation, PAD-US publication and name of aggregated data set or file.
	// Attributed in the format ‘organization name_PADUSversion_filename.filetype’ (e.g. TNC_PADUS1_4_SecuredAreas2008.shp).
	// Use acronym or replace spaces with underscore.
	// A data aggregator submits data in the PAD-US schema according to standards and/or aggregates regional or national datasets with required fields for PAD-US translation
	AggregatorSource string   `json:"Agg_Src"`
	Category         Category `json:"Category"`
	DCategory        string   `json:"d_Category"`
	// Comments from either the original data source or aggregator.
	Comments string `json:"Comments"`
	// The Year (yyyy) the protected area was designated, decreed or otherwise established. Date is assigned to each unit by name, without event status(e.g. Yellowstone National Park: 1872, Frank Church-River of No Return Wilderness Area: 1980)
	DateEst     string      `json:"Date_Est"`
	Designation Designation `json:"Des_Tp"`
	DDesTp      string      `json:"d_Des_Tp"`
	// The most current Year (yyyy) the GAP Status Code was assigned to the polygon.
	GAPCdDt string `json:"GAPCdDt"`
	// An acronym to describe the organization(s) that applied Gap Status Code to the polygon. This field also describes the methods used for assigning GAP Status as follows: ‘GAP – Default’ is assigned when GAP’s categorical assignment of status has been applied, without more detailed review or inquiry. ‘GAP’ is assigned when standard methods (management plan reviewed and/or land manager interviewed to assign GAP Status to a protected area) apply as provided above. ‘GAP – other organization’ (e.g. GAP – NPS) applies when the measure is assigned in partnership with GAP, including review. When another organization applied GAP Status according to their methods then ‘other organization’ (e.g. TNC) is assigned. See the PAD-US Standards Manual for more information.
	GapStatusCodeSource string    `json:"GAPCdSrc"`
	GAPStatusCode       GAPStatus `json:"GAP_Sts"`
	DGAPSts             string    `json:"d_GAP_Sts"`
	// Acres calculated for each polygon converted from the Shape_Area Field
	GISAcres int64 `json:"GIS_Acres"`
	// The source of spatial data the aggregator obtained (e.g. WYGF_whmas08.shp) for each record.
	// Files names match original source data to increase update efficiency.
	GISSource string `json:"GIS_Src"`
	// The most current Year (yyyy) the IUCN Category was assigned to the polygon.
	IUCNCategoryDate string `json:"IUCNCtDt"`
	// An acronym to describe the organization(s) that applied IUCN Category to the polygon.
	// This field also describes the methods used for assigning IUCN Category as follows:
	// ‘GAP – Default’ is assigned when GAP’s categorical assignment of status has been applied, without additional review.
	// ‘GAP – other organization’ (e.g. GAP – NPS) applies when the measure is assigned in partnership with GAP, including review.
	// When another organization applies IUCN Category according to their methods then ‘other organization’ (e.g. TNC) is assigned.
	// See the PAD-US Standards Manual for more information.
	IUCNCategorySource string       `json:"IUCNCtSrc"`
	IUCNCategory       IUCNCategory `json:"IUCN_Cat"`
	DIUCNCat           string       `json:"d_IUCN_Cat"`
	// The unit’s land management description or designation as provided by data source. "Local Designation" is not standardized and complements "Designation Type", null values indicate data gaps in source files.
	LocalDesignation string `json:"Loc_Ds"`
	// The name of the land manager as provided by the data source or "Local Manager", to complement the standardized "Manager Name" field (e.g. ‘State Fish and Wildlife’ is a standard ‘Manager Name’, Washington Department of Fish and Wildlife is a ‘Local Manager’). Null values indicate data gaps in source files.
	LocalManager string `json:"Loc_Mang"`
	// The name of the protected area as provided by the data source; "Local Name" is not standardized. This field may include designations, different formats, spelling errors, unit or area identifiers unique to parcels; however, it links directly to source data files
	LocalName string `json:"Loc_Nm"`
	// The name of the land owner as provided by the data source or "Local Owner", to complement the standardized "Owner Name" field (e.g. ‘State Fish and Wildlife’ is a standard ‘Owner Name’, Washington Department of Fish and Wildlife is a ‘Local Owner’). Null values indicate data gaps in source files.
	LocalOwner  string      `json:"Loc_Own"`
	ManagerName ManagerName `json:"Mang_Name"`
	DMangNam    string      `json:"d_Mang_Nam"`
	ManagerType ManagerType `json:"Mang_Type"`
	DMangTyp    string      `json:"d_Mang_Typ"`
	OwnerName   OwnerName   `json:"Own_Name"`
	DOwnName    string      `json:"d_Own_Name"`
	OwnerType   OwnerType   `json:"Own_Type"`
	DOwnType    string      `json:"d_Own_Type"`
	// The date (yyyy/mm/dd) GIS data was published or obtained (in the case of infrequently updated files) by the data aggregator. If month or day is unknown date is yyyy/00/00.
	SourceDate string `json:"Src_Date"`
	// Name of state or territory spelled out in Proper Case. See domain descriptions in PAD-US Standards Manual or geodatabase look up table for details.
	StateNm  string `json:"State_Nm"`
	DStateNm string `json:"d_State_Nm"`
	UnitNm   string `json:"Unit_Nm"`
	// Source PA Identifier
	SourcePAI string `json:"Source_PAI"`
	// Site Id assigned by UNEP-World Conservation Monitoring Center (WCMC) to all multi-part polygons of same protected
	//	area (GAP Status Code 1 and 2 only) in the World Database for Protected Areas (WDPA). Site Id's are maintained in
	//	PAD-US to facilitate WDPA updates and provide users with the ability to select all the multiple parts of one protected
	//	area. Currently, WDPA Site Id's are only assigned to GAP Status 1 and 2 protected areas as only these areas are defined
	//	as protected for the preservation of biodiversity as defined by the International Union for the Conservation of Nature
	//(IUCN) for the WDPA. WDPA Site ID Codes from PAD-US version 1 have been assigned in PAD-USv1.1, except for
	//	updates. Work is underway to reassign WDPA Codes to the Northwest, California and the Northeast that were completely
	//	replaced. GAP and WCMC are working on automated methods to maintain WDPA Site ID codes and extend them to all
	//	the protected areas in PAD-US.
	WDPACd int64 `json:"WDPA_Cd"`
}

// Note: Missing from records but found in documentation:
// Attribute_Label: Des_Nm
// Attribute_Definition:
// The name of protected area following the PAD-US Standard (i.e. full name including designation type to Proper Case without acronyms, special characters, space or return errors). As null values are not permitted in this standardized field, categorical assignments are made from "Manager Name" when data gaps occur in source files.
