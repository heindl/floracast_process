package pad_us

// https://gapanalysis.usgs.gov/padus/data/metadata/
type ProtectedArea struct {
	// Level of public access permitted. Open requires no special requirements for public access to the property (may include regular hours available); Restricted requires a special permit from the owner for access, a registration permit on public land or has highly variable times when open to use; Closed occurs where no public access allowed (land bank property, special ecological study areas, military bases, etc. Unknown is assigned where information is not currently available. Access is assigned categorically by Designation Type or provided by PAD-US State Data Stewards, federal or NGO partners. Contact the PAD-US Coordinator with available public access information
	Access  PublicAccess `json:"Access"`
	DAccess string       `json:"d_Access"`
	// Documents the Source of Access domain assignment (e.g. State Data Steward or ‘GAP Default’ categorical assignment).
	AccessSource string `json:"Access_Src"`
	// "Aggregator Source" identifies the organization credited with data aggregation, PAD-US publication and name of aggregated data set or file. Attributed in the format ‘organization name_PADUSversion_filename.filetype’ (e.g. TNC_PADUS1_4_SecuredAreas2008.shp). Use acronym or replace spaces with underscore. A data aggregator submits data in the PAD-US schema according to standards and/or aggregates regional or national datasets with required fields for PAD-US translation
	AggregatorSource    string `json:"Agg_Src"`
	// General category for the protection mechanism associated with the protected area. ‘Fee’ is the most common way real estate is owned. A conservation ‘easement’ creates a legally enforceable land preservation agreement between a landowner and government agency or qualified land protection organization (i.e. land trust). ‘Other’ types of protection include leases, agreements or those over marine waters. ‘Designation’ is applied to designations in the federal theme not tied to title documents (e.g. National Monument, Wild and Scenic River). These may be removed to reduce overlaps for area based analyses.
	Category  Category `json:"Category"`
	DCategory string `json:"d_Category"`
	// Comments from either the original data source or aggregator.
	Comments  string `json:"Comments"`
	// The Year (yyyy) the protected area was designated, decreed or otherwise established. Date is assigned to each unit by name, without event status(e.g. Yellowstone National Park: 1872, Frank Church-River of No Return Wilderness Area: 1980)
	DateEst   string `json:"Date_Est"`
	// The unit’s land management description or designation, standardized for nation (e.g. Area of Critical Environmental Concern, Wilderness Area, State Park, Local Rec Area, Conservation Easement). See the PAD-US Data Standard for a crosswalk of "Designation Type" from source data files or the geodatabase look up table for "Designation Type" for domain descriptions. "Designation Type" supports PAD-US queries and categorical conservation measures or public access assignments in the absence of other information.
	Designation Designation `json:"Des_Tp"`
	DDesTp      string      `json:"d_Des_Tp"`
	// The most current Year (yyyy) the GAP Status Code was assigned to the polygon.
	GAPCdDt   string `json:"GAPCdDt"`
	// An acronym to describe the organization(s) that applied Gap Status Code to the polygon. This field also describes the methods used for assigning GAP Status as follows: ‘GAP – Default’ is assigned when GAP’s categorical assignment of status has been applied, without more detailed review or inquiry. ‘GAP’ is assigned when standard methods (management plan reviewed and/or land manager interviewed to assign GAP Status to a protected area) apply as provided above. ‘GAP – other organization’ (e.g. GAP – NPS) applies when the measure is assigned in partnership with GAP, including review. When another organization applied GAP Status according to their methods then ‘other organization’ (e.g. TNC) is assigned. See the PAD-US Standards Manual for more information.
	GapStatusCodeSource  string `json:"GAPCdSrc"`
	// The GAP Status Code is a measure of management intent to conserve biodiversity defined as: Status 1: An area having permanent protection from conversion of natural land cover and a mandated management plan in operation to maintain a natural state within which disturbance events (of natural type, frequency, intensity, and legacy) are allowed to proceed without interference or are mimicked through management. Status 2: An area having permanent protection from conversion of natural land cover and a mandated management plan in operation to maintain a primarily natural state, but which may receive uses or management practices that degrade the quality of existing natural communities, including suppression of natural disturbance. Status 3: An area having permanent protection from conversion of natural land cover for the majority of the area, but subject to extractive uses of either a broad, low-intensity type (e.g., logging, OHV recreation) or localized intense type (e.g., mining). It also confers protection to federally listed endangered and threatened species throughout the area. Status 4: There are no known public or private institutional mandates or legally recognized easements or deed restrictions held by the managing entity to prevent conversion of natural habitat types to anthropogenic habitat types. The area generally allows conversion to unnatural land cover throughout or management intent is unknown. See the PAD-US Standards Manual for a summary of methods or the geodatabase look up table for short descriptions.
	GAPStatusCode    GAPStatus `json:"GAP_Sts"`
	DGAPSts   string `json:"d_GAP_Sts"`
	// Acres calculated for each polygon converted from the Shape_Area Field
	GISAcres  int64  `json:"GIS_Acres"`
	// The source of spatial data the aggregator obtained (e.g. WYGF_whmas08.shp) for each record. Files names match original source data to increase update efficiency.
	GISSource    string `json:"GIS_Src"`
	// The most current Year (yyyy) the IUCN Category was assigned to the polygon.
	IUCNCategoryDate  string `json:"IUCNCtDt"`
	// An acronym to describe the organization(s) that applied IUCN Category to the polygon. This field also describes the methods used for assigning IUCN Category as follows: ‘GAP – Default’ is assigned when GAP’s categorical assignment of status has been applied, without additional review. ‘GAP – other organization’ (e.g. GAP – NPS) applies when the measure is assigned in partnership with GAP, including review. When another organization applies IUCN Category according to their methods then ‘other organization’ (e.g. TNC) is assigned. See the PAD-US Standards Manual for more information.
	IUCNCategorySource string `json:"IUCNCtSrc"`
	// International Union for the Conservation of Nature (IUCN) management categories assigned to protected areas for inclusion in the UNEP- World Conservation Monitoring Center’s (WCMC) World Database for Protected Areas (WDPA) and the Commission for Environmental Cooperation’s (CEC) North American Terrestrial Protected Areas Database. IUCN defines a protected area as, "A clearly defined geographical space, recognized, dedicated and managed, through legal or other effective means, to achieve the long-term conservation of nature with associated ecosystem services and cultural values" (includes GAP Status Code 1 and 2 only). Categorization follows as: Category Ia: Strict Nature Reserves are strictly protected areas set aside to protect biodiversity and also possibly geological / geomorphological features, where human visitation, use and impacts are strictly controlled and limited to ensure preservation of the conservation values. Such protected areas can serve as indispensable reference areas for scientific research and monitoring. Category Ib: Wilderness Areas are protected areas are usually large unmodified or slightly modified areas, retaining their natural character and influence, without permanent or significant human habitation, which are protected and managed so as to preserve their natural condition. Category II: National Park protected areas are large natural or near natural areas set aside to protect large-scale ecological processes, along with the complement of species and ecosystems characteristic of the area, which also provide a foundation for environmentally and culturally compatible spiritual, scientific, educational, recreational and visitor opportunities. Category III: Natural Monument or Feature protected areas are set aside to protect a specific natural monument, which can be a land form, sea mount, submarine caverns, geological feature such as caves or even a living feature such as an ancient grove. They are generally quite small protected areas and often have high visitor value. Category IV: Habitat/species management protected areas aim to protect particular species or habitats and management reflects this priority. Many category IV protected areas will need regular, active interventions to address the requirements of particular species or to maintain habitats, but this is not a requirement of this category. Category V: Protected landscape/seascape protected areas occur where the interaction of people and nature over time has produced an area of distinct character with significant ecological, biological, cultural and scenic value. Category VI: Protected area with sustainable use (community based, non industrial) of natural resources are generally large, with much of the area in a more-or-less natural condition and whereas a proportion is under sustainable natural resource management and where such exploitation is seen as one of the main aims of the area. "Other Conservation Areas" are not recognized by IUCN at this time; however, they are included in the CEC’s database referenced above. These areas (includes GAP Status Code 3 only)are included in the IUCN Category Domain along with "Unassigned" areas (includes GAP Status Code 4) and "Not Reported" areas that meet the definition of IUCN protection (i.e. GAP Status Code 1 or 2) but where IUCN Category has not yet been assigned and categorical assignment is not appropriate. See the PAD-US Standards Manual for a summary of methods.
	IUCNCategory   IUCNCategory `json:"IUCN_Cat"`
	DIUCNCat  string `json:"d_IUCN_Cat"`
	// The unit’s land management description or designation as provided by data source. "Local Designation" is not standardized and complements "Designation Type", null values indicate data gaps in source files.
	LocalDesignation     string `json:"Loc_Ds"`
	// The name of the land manager as provided by the data source or "Local Manager", to complement the standardized "Manager Name" field (e.g. ‘State Fish and Wildlife’ is a standard ‘Manager Name’, Washington Department of Fish and Wildlife is a ‘Local Manager’). Null values indicate data gaps in source files.
	LocalManager   string `json:"Loc_Mang"`
	// The name of the protected area as provided by the data source; "Local Name" is not standardized. This field may include designations, different formats, spelling errors, unit or area identifiers unique to parcels; however, it links directly to source data files
	LocalName     string `json:"Loc_Nm"`
	// The name of the land owner as provided by the data source or "Local Owner", to complement the standardized "Owner Name" field (e.g. ‘State Fish and Wildlife’ is a standard ‘Owner Name’, Washington Department of Fish and Wildlife is a ‘Local Owner’). Null values indicate data gaps in source files.
	LocalOwner    string `json:"Loc_Own"`
	// Land manager or administrative agency (e.g. USFS, State Fish and Game, City Land, TNC) standardized for the US. See PAD-US Data Standard or geodatabase look up table for "Agency Name" for full domain descriptions. Use "Manager Name" for the best depiction of federal lands as many overlapping designations (i.e. "Designation") and ownership related data gaps (i.e. "Unknown") occur in the federal theme. GAP attributes the applicable "Agency Name" to all records provided by the agency data steward in the "Manager Name" field. "Owner Name" contains ‘unknown’ values where parcel level ownership data are not yet available from authoritative data sources.
	ManagerName  ManagerName `json:"Mang_Name"`
	DMangNam  string `json:"d_Mang_Nam"`
	// General land manager description (e.g. Federal, Tribal, State, Private) standardized for the US. See PAD-US Data Standard for "Agency Name" to "Agency Type" crosswalk or geodatabase look up table for full domain descriptions. Use "Manager Type" for the most complete depiction of federal lands as overlapping designations (i.e. "Designation") and ownership related data gaps (i.e. "Unknown") occur in the federal theme.
	ManagerType  ManagerType `json:"Mang_Type"`
	DMangTyp  string `json:"d_Mang_Typ"`
	// Land owner or holding agency (e.g. USFS, State Fish and Game, City Land, TNC) standardized for the US. See PAD-US Data Standard or geodatabase look up table for "Agency Name" for full domain descriptions . Use "Manager Name" for the best depiction of federal lands as many overlapping designations (i.e. "Designation") and ownership related data gaps (i.e. "Unknown") occur in the federal theme.
	OwnerName   OwnerName `json:"Own_Name"`
	DOwnName  string `json:"d_Own_Name"`
	// General land owner description (e.g. Federal, Tribal, State, Private) standardized for the US. See PAD-US Data Standard for "Agency Name" to "Agency Type" crosswalk or geodatabase look up table for full domain descriptions. Use "Manager Type" for the best depiction of federal lands as many overlapping designations (i.e. "Designation") and ownership related data gaps (i.e. "Unknown") occur in the federal theme.
	OwnerType   OwnerType `json:"Own_Type"`
	DOwnType  string `json:"d_Own_Type"`
	// The date (yyyy/mm/dd) GIS data was published or obtained (in the case of infrequently updated files) by the data aggregator. If month or day is unknown date is yyyy/00/00.
	SourceDate   string `json:"Src_Date"`
	// Name of state or territory spelled out in Proper Case. See domain descriptions in PAD-US Standards Manual or geodatabase look up table for details.
	StateNm   string `json:"State_Nm"`
	DStateNm  string `json:"d_State_Nm"`
	UnitNm    string `json:"Unit_Nm"`
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
	WDPACd    int64  `json:"WDPA_Cd"`
}

// Attribute_Label: Des_Nm
// Attribute_Definition:
// The name of protected area following the PAD-US Standard (i.e. full name including designation type to Proper Case without acronyms, special characters, space or return errors). As null values are not permitted in this standardized field, categorical assignments are made from "Manager Name" when data gaps occur in source files.
// Attribute_Definition_Source: