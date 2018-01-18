package pad_us

type GAPStatus string

var GAPStatusDefinitions = map[GAPStatus]string{
	GAPStatus("1"): "1 - managed for biodiversity - disturbance events proceed or are mimicked",
	GAPStatus("2"): "2 - managed for biodiversity - disturbance events suppressed",
	GAPStatus("3"): "3 - managed for multiple uses - subject to extractive (e.g. mining or logging) or OHV use",
	GAPStatus("4"): "4 - no known mandate for protection", // Ends up covering many Wildlife Refuges
}

type Category string

var CategoryDefinitions = map[Category]string{
	Category("Designation"): "Designation",
	Category("Easement"):    "Easement",
	Category("Fee"):         "Fee",
	Category("Other"):       "Other", // Mostly Wildlife Management Areas in Alabama
	Category("Unknown"):     "Unknown",  // Mostly Preserves in Florida
}

type Designation string

var DesignationDefinitions = map[Designation]string{
	Designation("ACC"):   "Access Area",
	Designation("ACEC"):  "Area of Critical Environmental Concern",
	//Designation("AGRE"): "Agricultural Easement",
	Designation("CONE"):  "Conservation Easement",
	Designation("FORE"): "Forest Stewardship Easement",
	Designation("FOTH"):  "Federal Other or Unknown",
	Designation("HCA"):   "Historic or Cultural Area",
	Designation("IRA"): "Inventoried Roadless Area",
	Designation("LCA"):   "Local Conservation Area",
	Designation("LHCA"):  "Local Historic or Cultural Area",
	Designation("LOTH"):  "Local Other or Unknown",
	// Designation("LP"):    "Local Park", // These tend to be baseball fields or golf courses, so skip.
	Designation("LREC"):  "Local Recreation Area",
	Designation("LRMA"):  "Local Resource Management Area",
	//Designation("MIL"):   "Military Land",
	Designation("MIT"):   "Mitigation Land or Bank",
	//Designation("MPA"):   "Marine Protected Area",
	Designation("NCA"):   "Conservation Area",
	Designation("ND"):    "Not Designated",
	Designation("NF"):    "National Forest",
	Designation("NG"): "National Grassland",
	Designation("NLS"): "National Lakeshore or Seashore",
	Designation("NM"):    "National Monument or Landmark",
	Designation("NP"):    "National Park",
	Designation("NRA"):   "National Recreation Area",
	Designation("NSBV"):  "National Scenic, Botanical or Volcanic Area",
	Designation("NT"):    "National Scenic or Historic Trail",
	Designation("NWR"):   "National Wildlife Refuge",
	Designation("OTHE"):  "Other Easement",
	Designation("PAGR"):  "Private Agricultural",
	Designation("PCON"):  "Private Conservation",
	Designation("PFOR"):  "Private Forest Stewardship",
	Designation("PHCA"):  "Private Historic or Cultural",
	Designation("POTH"):  "Private Other or Unknown",
	Designation("PREC"):  "Private Recreation or Education",
	Designation("PROC"):  "Approved or Proclamation Boundary",
	Designation("PUB"): "National Public Lands",
	Designation("RANE"): "Ranch Easement",
	Designation("REA"):   "Research or Educational Area",
	Designation("REC"):   "Recreation Management Area",
	Designation("RECE"):  "Recreation or Education Easement",
	Designation("RMA"):   "Resource Management Area",
	Designation("RNA"):   "Research Natural Area",
	Designation("SCA"):   "State Conservation Area",
	Designation("SDA"):   "Special Designation Area",
	Designation("SHCA"):  "State Historic or Cultural Area",
	Designation("SOTH"):  "State Other or Unknown",
	Designation("SP"):    "State Park",
	Designation("SREC"):  "State Recreation Area",
	Designation("SRMA"):  "State Resource Management Area",
	Designation("SW"):    "State Wilderness",
	//Designation("TRIBL"): "Native American Land",
	Designation("WA"):    "Wilderness Area",
	Designation("WPA"):   "Watershed Protection Area",
	Designation("WSA"):   "Wilderness Study Area",
	Designation("WSR"):   "Wild and Scenic River",
	Designation("UNK"): "Unknown",
	Designation("UNKE"): "Unknown Easement",
}

type ManagerName string

var ManagerNameDefinitions = map[ManagerName]string{
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

type ManagerType string

var ManagerTypeDefinitions = map[ManagerType]string{
	ManagerType("DIST"): "Regional Agency Special District",
	ManagerType("FED"):  "Federal",
	ManagerType("JNT"):  "Joint",
	ManagerType("LOC"):  "Local Government",
	ManagerType("NGO"):  "Non-Governmental Organization",
	ManagerType("PVT"):  "Private",
	ManagerType("STAT"): "State",
	ManagerType("UNK"):  "Unknown",
}

type OwnerName string

var OwnerNameDefinitions = map[OwnerName]string{
	OwnerName("BLM"):   "Bureau of Land Management",
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

type PublicAccess string

const PublicAccessClosed = PublicAccess("XA");
const PublicAccessRestricted = PublicAccess("RA")
const PublicAccessUnknown = PublicAccess("UK")

var PublicAccessDefinitions = map[PublicAccess]string{
	PublicAccess("OA"): "Open Access",
	PublicAccess("RA"): "Restricted Access",
	PublicAccess("UK"): "Unknown",
	PublicAccess("XA"): "Closed",
}

type IUCNCategory string

var AreaIUCNCategoryDefinitions = map[IUCNCategory]string{
	IUCNCategory("II"):  "II: National park",
	IUCNCategory("III"): "III: Natural monument or feature",
	IUCNCategory("IV"):  "IV: Habitat / species management",
	IUCNCategory("Ia"):  "Ia: Strict nature reserves",
	IUCNCategory("Ib"):  "Ib: Wilderness areas",
	IUCNCategory("N/R"): "Not Reported",
	IUCNCategory("V"):   "V: Protected landscape / seascape",
	IUCNCategory("VI"):  "VI: Protected area with sustainable use of natural resources",
	IUCNCategory("Other Conservation Area"): "Other Conservation Area",
	IUCNCategory("Unassigned"): "Unassigned", // Looks like state parks in smaller states are generally unassigned, so let it pass.

}

type OwnerType string

var OwnerTypeDefinitions = map[OwnerType]string{
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
