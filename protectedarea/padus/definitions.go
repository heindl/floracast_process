package padus

import (
	"github.com/dropbox/godropbox/errors"
	"strconv"
)

// ErrInvalidPADUSProperty is the flag for function that returns an invalid property warning.
var ErrInvalidPADUSProperty = errors.New("Invalid PAD-US Property")

// GAPStatus is the GAP Status Code is a measure of management intent to conserve biodiversity defined as:
// Status 1: An area having permanent protection from conversion of natural land cover and a mandated
// management plan in operation to maintain a natural state within which disturbance events (of natural type,
// frequency, intensity, and legacy) are allowed to proceed without interference or are mimicked through management.
// Status 2: An area having permanent protection from conversion of natural land cover and a mandated
// management plan in operation to maintain a primarily natural state, but which may receive uses or
// management practices that degrade the quality of existing natural communities, including suppression of natural disturbance.
// Status 3: An area having permanent protection from conversion of natural land cover for the majority of
// the area, but subject to extractive uses of either a broad, low-intensity type
// (e.g., logging, OHV recreation) or localized intense type (e.g., mining). It also confers protection to
// federally listed endangered and threatened species throughout the area.
// Status 4: There are no known public or private institutional mandates or legally recognized easements
// or deed restrictions held by the managing entity to prevent conversion of natural habitat types to
// anthropogenic habitat types. The area generally allows conversion to unnatural land cover throughout
// or management intent is unknown. See the PAD-US Standards Manual for a summary of methods or the
// geodatabase look up table for short descriptions.
type GAPStatus string

// Valid returns true is the property is accepted for use.
func (Ω GAPStatus) Valid() bool {
	_, ok := gapStatusDefinitions[Ω]
	return ok
}

// ProtectionLevel converts GAPStatus to integer between one and four
func (Ω GAPStatus) ProtectionLevel() (int, error) {
	i, err := strconv.Atoi(string(Ω))
	if err != nil {
		return 0, errors.Wrapf(err, "Invalid GAPStatus [%s]", Ω)
	}

	if i < 1 || i > 4 {
		return 0, errors.Wrapf(ErrInvalidPADUSProperty, "GAPStatus [%s]", Ω)
	}

	return i, nil
}

var gapStatusDefinitions = map[GAPStatus]string{
	GAPStatus("1"): "1 - managed for biodiversity - disturbance events proceed or are mimicked",
	GAPStatus("2"): "2 - managed for biodiversity - disturbance events suppressed",
	GAPStatus("3"): "3 - managed for multiple uses - subject to extractive (e.g. mining or logging) or OHV use",
	GAPStatus("4"): "4 - no known mandate for protection", // Ends up covering many Wildlife Refuges
}

// Category is the general category for the protection mechanism associated with the protected area.
// ‘Fee’ is the most common way real estate is owned.
// A conservation ‘easement’ creates a legally enforceable land preservation agreement
// between a landowner and government agency or qualified land protection organization (i.e. land trust).
// ‘Other’ types of protection include leases, agreements or those over marine waters.
// ‘Designation’ is applied to designations in the federal theme not tied to title documents (e.g. National Monument, Wild and Scenic River).
// These may be removed to reduce overlaps for area based analyses.
type Category string

// Valid returns true is the property is accepted for use.
func (Ω Category) Valid() bool {
	_, ok := categoryDefinitions[Ω]
	return ok
}

var categoryDefinitions = map[Category]string{
	Category("Designation"): "Designation",
	Category("Easement"):    "Easement",
	Category("Fee"):         "Fee",
	Category("Other"):       "Other",   // Mostly Wildlife Management Areas in Alabama
	Category("Unknown"):     "Unknown", // Mostly Preserves in Florida
}

// Designation is the unit’s land management description or designation, standardized for nation
// (e.g. Area of Critical Environmental Concern, Wilderness Area, State Park, Local Rec Area,
// Conservation Easement). See the PAD-US Data Standard for a crosswalk of "Designation Type"
// from source data files or the geodatabase look up table for "Designation Type" for domain descriptions.
// "Designation Type" supports PAD-US queries and categorical conservation measures or public access
// assignments in the absence of other information.
type Designation string

// Valid returns true is the property is accepted for use.
func (Ω Designation) Valid() bool {
	_, ok := designationDefinitions[Ω]
	return ok
}

var designationDefinitions = map[Designation]string{
	Designation("ACC"):  "Access Area",
	Designation("ACEC"): "Area of Critical Environmental Concern",
	//Designation("AGRE"): "Agricultural Easement",
	Designation("CONE"): "Conservation Easement",
	Designation("FORE"): "Forest Stewardship Easement",
	Designation("FOTH"): "Federal Other or Unknown",
	Designation("HCA"):  "Historic or Cultural Area",
	Designation("IRA"):  "Inventoried Roadless Area",
	Designation("LCA"):  "Local Conservation Area",
	Designation("LHCA"): "Local Historic or Cultural Area",
	Designation("LOTH"): "Local Other or Unknown",
	// Designation("LP"):    "Local Park", // These tend to be baseball fields or golf courses, so skip.
	Designation("LREC"): "Local Recreation Area",
	Designation("LRMA"): "Local Resource Management Area",
	//Designation("MIL"):   "Military Land",
	Designation("MIT"): "Mitigation Land or Bank",
	//Designation("MPA"):   "Marine Protected Area",
	Designation("NCA"):  "Conservation Area",
	Designation("ND"):   "Not Designated",
	Designation("NF"):   "National Forest",
	Designation("NG"):   "National Grassland",
	Designation("NLS"):  "National Lakeshore or Seashore",
	Designation("NM"):   "National Monument or Landmark",
	Designation("NP"):   "National Park",
	Designation("NRA"):  "National Recreation Area",
	Designation("NSBV"): "National Scenic, Botanical or Volcanic Area",
	Designation("NT"):   "National Scenic or Historic Trail",
	Designation("NWR"):  "National Wildlife Refuge",
	Designation("OTHE"): "Other Easement",
	Designation("PAGR"): "Private Agricultural",
	Designation("PCON"): "Private Conservation",
	Designation("PFOR"): "Private Forest Stewardship",
	Designation("PHCA"): "Private Historic or Cultural",
	Designation("POTH"): "Private Other or Unknown",
	Designation("PREC"): "Private Recreation or Education",
	Designation("PROC"): "Approved or Proclamation Boundary",
	Designation("PUB"):  "National Public Lands",
	Designation("RANE"): "Ranch Easement",
	Designation("REA"):  "Research or Educational Area",
	Designation("REC"):  "Recreation Management Area",
	Designation("RECE"): "Recreation or Education Easement",
	Designation("RMA"):  "Resource Management Area",
	Designation("RNA"):  "Research Natural Area",
	Designation("SCA"):  "State Conservation Area",
	Designation("SDA"):  "Special Designation Area",
	Designation("SHCA"): "State Historic or Cultural Area",
	Designation("SOTH"): "State Other or Unknown",
	Designation("SP"):   "State Park",
	Designation("SREC"): "State Recreation Area",
	Designation("SRMA"): "State Resource Management Area",
	Designation("SW"):   "State Wilderness",
	//Designation("TRIBL"): "Native American Land",
	Designation("WA"):   "Wilderness Area",
	Designation("WPA"):  "Watershed Protection Area",
	Designation("WSA"):  "Wilderness Study Area",
	Designation("WSR"):  "Wild and Scenic River",
	Designation("UNK"):  "Unknown",
	Designation("UNKE"): "Unknown Easement",
}

// ManagerName is the land manager or administrative agency (e.g. USFS, State Fish and Game, City Land, TNC)
// standardized for the US. See PAD-US Data Standard or geodatabase look up table for "Agency Name" for
// full domain descriptions. Use "Manager Name" for the best depiction of federal lands as many overlapping
// designations (i.e. "Designation") and ownership related data gaps (i.e. "Unknown") occur in the federal theme.
// GAP attributes the applicable "Agency Name" to all records provided by the agency data steward in the "Manager Name" field.
// "Owner Name" contains ‘unknown’ values where parcel level ownership data are not yet available from authoritative data sources.
type ManagerName string

// Valid returns true is the property is accepted for use.
func (Ω ManagerName) Valid() bool {
	_, ok := managerNameDefinitions[Ω]
	return ok
}

var managerNameDefinitions = map[ManagerName]string{
	ManagerName("ARS"):   "Agricultural Research Service",
	ManagerName("BIA"):   "Bureau of Indian Affairs",
	ManagerName("BLM"):   "Bureau of Land Management",
	ManagerName("CITY"):  "City Land",
	ManagerName("CNTY"):  "County Land",
	ManagerName("DOD"):   "Department of Defense",
	ManagerName("FWS"):   "U.S. Fish & Wildlife Service",
	ManagerName("JNT"):   "Joint",
	ManagerName("NGO"):   "Non-Governmental Organization",
	ManagerName("NOAA"):  "National Oceanic and Atmospheric Administration",
	ManagerName("NPS"):   "National Park Service",
	ManagerName("NRCS"):  "Natural Resources Conservation Service",
	ManagerName("OTHF"):  "Other or Unknown Federal Land",
	ManagerName("OTHR"):  "Other",
	ManagerName("OTHS"):  "Other or Unknown State Land",
	ManagerName("PVT"):   "Private",
	ManagerName("REG"):   "Regional Agency Land",
	ManagerName("RWD"):   "Regional Water Districts",
	ManagerName("SDC"):   "State Department of Conservation",
	ManagerName("SDNR"):  "State Department of Natural Resources",
	ManagerName("SDOL"):  "State Department of Land",
	ManagerName("SFW"):   "State Fish and Wildlife",
	ManagerName("SLB"):   "State Land Board",
	ManagerName("SPR"):   "State Park and Recreation",
	ManagerName("TVA"):   "Tennessee Valley Authority",
	ManagerName("UNK"):   "Unknown",
	ManagerName("UNKL"):  "Other or Unknown Local Government",
	ManagerName("USACE"): "Army Corps of Engineers",
	ManagerName("USBR"):  "Bureau of Reclamation",
	ManagerName("USFS"):  "Forest Service",
}

// ManagerType is the general land manager description (e.g. Federal, Tribal, State, Private)
// standardized for the US. See PAD-US Data Standard for "Agency Name" to "Agency Type" crosswalk
// or geodatabase look up table for full domain descriptions. Use "Manager Type" for the most
// complete depiction of federal lands as overlapping designations (i.e. "Designation") and
// ownership related data gaps (i.e. "Unknown") occur in the federal theme.
type ManagerType string

// Valid returns true is the property is accepted for use.
func (Ω ManagerType) Valid() bool {
	_, ok := managerTypeDefinitions[Ω]
	return ok
}

var managerTypeDefinitions = map[ManagerType]string{
	ManagerType("DIST"): "Regional Agency Special District",
	ManagerType("FED"):  "Federal",
	ManagerType("JNT"):  "Joint",
	ManagerType("LOC"):  "Local Government",
	ManagerType("NGO"):  "Non-Governmental Organization",
	ManagerType("PVT"):  "Private",
	ManagerType("STAT"): "State",
	ManagerType("UNK"):  "Unknown",
}

// OwnerName is the land owner or holding agency (e.g. USFS, State Fish and Game, City Land, TNC)
// standardized for the US. See PAD-US Data Standard or geodatabase look up table for "Agency Name"
// for full domain descriptions . Use "Manager Name" for the best depiction of federal lands
// as many overlapping designations (i.e. "Designation") and ownership related data gaps (i.e. "Unknown")
// occur in the federal theme.
type OwnerName string

// Valid returns true is the property is accepted for use.
func (Ω OwnerName) Valid() bool {
	_, ok := ownerNameDefinitions[Ω]
	return ok
}

var ownerNameDefinitions = map[OwnerName]string{
	OwnerName("BLM"):  "Bureau of Land Management",
	OwnerName("CITY"): "City Land",
	OwnerName("CNTY"): "County Land",
	OwnerName("DESG"): "Designation",
	OwnerName("DOD"):  "Department of Defense",
	OwnerName("FWS"):  "U.S. Fish & Wildlife Service",
	OwnerName("JNT"):  "Joint",
	OwnerName("NGO"):  "Non-Governmental Organization",
	OwnerName("NOAA"): "National Oceanic and Atmospheric Administration",
	OwnerName("NPS"):  "National Park Service",
	OwnerName("OTHF"): "Other or Unknown Federal Land",
	OwnerName("OTHS"): "Other or Unknown State Land",
	OwnerName("PVT"):  "Private",
	OwnerName("REG"):  "Regional Agency Land",
	OwnerName("RWD"):  "Regional Water Districts",
	OwnerName("SDC"):  "State Department of Conservation",
	OwnerName("SDNR"): "State Department of Natural Resources",
	OwnerName("SDOL"): "State Department of Land",
	OwnerName("SFW"):  "State Fish and Wildlife",
	OwnerName("SLB"):  "State Land Board",
	OwnerName("SPR"):  "State Park and Recreation",
	//OwnerName("TRIB"): "American Indian Lands",
	OwnerName("UNK"):  "Unknown",
	OwnerName("UNKL"): "Other or Unknown Local Government",
	OwnerName("USBR"): "Bureau of Reclamation",
	OwnerName("USFS"): "Forest Service",
}

// PublicAccess is the level of public access permitted.
// Open requires no special requirements for public access to the property (may include regular hours available);
// Restricted requires a special permit from the owner for access, a registration permit on public land or has highly variable times when open to use;
// Closed occurs where no public access allowed (land bank property, special ecological study areas, military bases, etc.
// Unknown is assigned where information is not currently available.
// Access is assigned categorically by Designation Type or provided by PAD-US State Data Stewards, federal or NGO partners.
// Contact the PAD-US Coordinator with available public access information
type PublicAccess string

// Valid returns true is the property is accepted for use.
func (Ω PublicAccess) Valid() bool {
	_, ok := publicAccessDefinitions[Ω]
	return ok
}

// AccessLevel converts PublicAccess to integer between one and four
func (Ω PublicAccess) AccessLevel() (int, error) {
	switch Ω {
	case PublicAccessOpen:
		return 1, nil
	case PublicAccessRestricted:
		return 2, nil
	case PublicAccessUnknown:
		return 3, nil
	case PublicAccessClosed:
		return 4, nil
	default:
		return 0, errors.Wrapf(ErrInvalidPADUSProperty, "PublicAccess [%s]", Ω)
	}
}

// PublicAccessOpen [1, "OA"]
const PublicAccessOpen = PublicAccess("OA")

// PublicAccessRestricted [2, "RA"]
const PublicAccessRestricted = PublicAccess("RA")

// PublicAccessClosed [4, "XA"]
const PublicAccessClosed = PublicAccess("XA")

// PublicAccessUnknown [3, "UK"]
const PublicAccessUnknown = PublicAccess("UK")

var publicAccessDefinitions = map[PublicAccess]string{
	PublicAccessOpen:       "Open Access",
	PublicAccessRestricted: "Restricted Access",
	PublicAccessUnknown:    "Unknown",
	PublicAccessClosed:     "Closed",
}

// IUCNCategory is the International Union for the Conservation of Nature (IUCN) management categories
// assigned to protected areas for inclusion in the UNEP- World Conservation Monitoring Center’s (WCMC)
// World Database for Protected Areas (WDPA) and the Commission for Environmental Cooperation’s (CEC)
// North American Terrestrial Protected Areas Database.
type IUCNCategory string

// IUCN defines a protected area as,
// "A clearly defined geographical space, recognized, dedicated and managed, through legal or
// other effective means, to achieve the long-term conservation of nature with associated ecosystem
// services and cultural values" (includes GAP Status Code 1 and 2 only).
//
// Categorization follows as:
//
// Category Ia: Strict Nature Reserves are strictly protected areas set aside to protect biodiversity
// and also possibly geological / geomorphological features, where human visitation, use and impacts
// are strictly controlled and limited to ensure preservation of the conservation values.
// Such protected areas can serve as indispensable reference areas for scientific research and monitoring.

// Category Ib: Wilderness Areas are protected areas are usually large unmodified or slightly modified areas,
// retaining their natural character and influence, without permanent or significant human habitation,
// which are protected and managed so as to preserve their natural condition.
//
// Category II: National Park protected areas are large natural or near natural areas set aside to protect
// large-scale ecological processes, along with the complement of species and ecosystems characteristic of
// the area, which also provide a foundation for environmentally and culturally compatible spiritual,
// scientific, educational, recreational and visitor opportunities.
//
// Category III: Natural Monument or Feature protected areas are set aside to protect a specific
// natural monument, which can be a land form, sea mount, submarine caverns, geological feature
// such as caves or even a living feature such as an ancient grove. They are generally quite small
// protected areas and often have high visitor value.
//
// Category IV: Habitat/species management protected areas aim to protect particular species or habitats and management reflects this priority.
// Many category IV protected areas will need regular, active interventions to address the requirements
// of particular species or to maintain habitats, but this is not a requirement of this category.
//
// Category V: Protected landscape/seascape protected areas occur where the interaction of people and
// nature over time has produced an area of distinct character with significant ecological, biological,
// cultural and scenic value. Category VI: Protected area with sustainable use
// (community based, non industrial) of natural resources are generally large, with much of the
// area in a more-or-less natural condition and whereas a proportion is under sustainable natural
// resource management and where such exploitation is seen as one of the main aims of the area.
//
// "Other Conservation Areas" are not recognized by IUCN at this time; however, they are included
// in the CEC’s database referenced above. These areas (includes GAP Status Code 3 only)are included in the
// IUCN Category Domain along with "Unassigned" areas (includes GAP Status Code 4) and
// "Not Reported" areas that meet the definition of IUCN protection (i.e. GAP Status Code 1 or 2)
// but where IUCN Category has not yet been assigned and categorical assignment is not appropriate.
// See the PAD-US Standards Manual for a summary of methods.

// Valid returns true is the property is accepted for use.
func (Ω IUCNCategory) Valid() bool {
	_, ok := iucnCategoryDefinitions[Ω]
	return ok
}

var iucnCategoryDefinitions = map[IUCNCategory]string{
	IUCNCategory("II"):                      "II: National park",
	IUCNCategory("III"):                     "III: Natural monument or feature",
	IUCNCategory("IV"):                      "IV: Habitat / species management",
	IUCNCategory("Ia"):                      "Ia: Strict nature reserves",
	IUCNCategory("Ib"):                      "Ib: Wilderness areas",
	IUCNCategory("N/R"):                     "Not Reported",
	IUCNCategory("V"):                       "V: Protected landscape / seascape",
	IUCNCategory("VI"):                      "VI: Protected area with sustainable use of natural resources",
	IUCNCategory("Other Conservation Area"): "Other Conservation Area",
	IUCNCategory("Unassigned"):              "Unassigned", // Looks like state parks in smaller states are generally unassigned, so let it pass.

}

// OwnerType is the general land owner description (e.g. Federal, Tribal, State, Private) standardized for the US.
// See PAD-US Data Standard for "Agency Name" to "Agency Type" crosswalk or geodatabase look up table
// for full domain descriptions. Use "Manager Type" for the best depiction of federal lands as many
// overlapping designations (i.e. "Designation") and ownership related data gaps (i.e. "Unknown")
// occur in the federal theme.
type OwnerType string

// Valid returns true is the property is accepted for use.
func (Ω OwnerType) Valid() bool {
	_, ok := ownerTypeDefinitions[Ω]
	return ok
}

var ownerTypeDefinitions = map[OwnerType]string{
	OwnerType("DESG"): "Designation",
	OwnerType("DIST"): "Regional Agency Special District",
	OwnerType("FED"):  "Federal",
	OwnerType("JNT"):  "Joint",
	OwnerType("LOC"):  "Local Government",
	OwnerType("NGO"):  "Non-Governmental Organization",
	OwnerType("PVT"):  "Private",
	OwnerType("STAT"): "State",
	//OwnerType("TRIB"): "American Indian Lands",
	OwnerType("UNK"): "Unknown",
}
