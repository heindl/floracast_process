package pad_us

type GAPStatus string

var GAPStatusDefinitions = map[GAPStatus]string{
	GAPStatus("1"): "1 - managed for biodiversity - disturbance events proceed or are mimicked",
	GAPStatus("2"): "2 - managed for biodiversity - disturbance events suppressed",
}

type Category string

var CategoryDefinitions = map[Category]string{
	Category("Designation"): "Designation",
	Category("Easement"):    "Easement",
	Category("Fee"):         "Fee",
	Category("Other"):       "Other",
	Category("Unknown"):     "Unknown",
}

type Designation string

var DesignationDefinitions = map[Designation]string{
	Designation("ACC"):   "Access Area",
	Designation("ACEC"):  "Area of Critical Environmental Concern",
	Designation("CONE"):  "Conservation Easement",
	Designation("FOTH"):  "Federal Other or Unknown",
	Designation("HCA"):   "Historic or Cultural Area",
	Designation("LCA"):   "Local Conservation Area",
	Designation("LHCA"):  "Local Historic or Cultural Area",
	Designation("LOTH"):  "Local Other or Unknown",
	Designation("LP"):    "Local Park",
	Designation("LREC"):  "Local Recreation Area",
	Designation("LRMA"):  "Local Resource Management Area",
	Designation("MIL"):   "Military Land",
	Designation("MIT"):   "Mitigation Land or Bank",
	Designation("MPA"):   "Marine Protected Area",
	Designation("NCA"):   "Conservation Area",
	Designation("ND"):    "Not Designated",
	Designation("NF"):    "National Forest",
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
	Designation("TRIBL"): "Native American Land",
	Designation("UNK"):   "Unknown",
	Designation("WA"):    "Wilderness Area",
	Designation("WPA"):   "Watershed Protection Area",
	Designation("WSA"):   "Wilderness Study Area",
	Designation("WSR"):   "Wild and Scenic River",
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

type AreaOwnerName string

var AreaOwnerNameDefinitions = map[AreaOwnerName]string{
	AreaOwnerName("CITY"): "City Land",
	AreaOwnerName("CNTY"): "County Land",
	AreaOwnerName("DESG"): "Designation",
	AreaOwnerName("DOD"):  "Department of Defense",
	AreaOwnerName("FWS"):  "U.S. Fish & Wildlife Service",
	AreaOwnerName("JNT"):  "Joint",
	AreaOwnerName("NGO"):  "Non-Governmental Organization",
	AreaOwnerName("NOAA"): "National Oceanic and Atmospheric Administration",
	AreaOwnerName("NPS"):  "National Park Service",
	AreaOwnerName("OTHF"): "Other or Unknown Federal Land",
	AreaOwnerName("OTHS"): "Other or Unknown State Land",
	AreaOwnerName("PVT"):  "Private",
	AreaOwnerName("REG"):  "Regional Agency Land",
	AreaOwnerName("RWD"):  "Regional Water Districts",
	AreaOwnerName("SDC"):  "State Department of Conservation",
	AreaOwnerName("SDNR"): "State Department of Natural Resources",
	AreaOwnerName("SDOL"): "State Department of Land",
	AreaOwnerName("SFW"):  "State Fish and Wildlife",
	AreaOwnerName("SLB"):  "State Land Board",
	AreaOwnerName("SPR"):  "State Park and Recreation",
	AreaOwnerName("TRIB"): "American Indian Lands",
	AreaOwnerName("UNK"):  "Unknown",
	AreaOwnerName("UNKL"): "Other or Unknown Local Government",
	AreaOwnerName("USBR"): "Bureau of Reclamation",
	AreaOwnerName("USFS"): "Forest Service",
}

type AreaPublicAccess string

var AreaPublicAccessDefinitions = map[AreaPublicAccess]string{
	AreaPublicAccess("OA"): "Open Access",
	AreaPublicAccess("RA"): "Restricted Access",
	AreaPublicAccess("UK"): "Unknown",
	AreaPublicAccess("XA"): "Closed",
}

type AreaIUCNCategory string

var AreaIUCNCategoryDefinitions = map[AreaIUCNCategory]string{
	AreaIUCNCategory("II"):  "II: National park",
	AreaIUCNCategory("III"): "III: Natural monument or feature",
	AreaIUCNCategory("IV"):  "IV: Habitat / species management",
	AreaIUCNCategory("Ia"):  "Ia: Strict nature reserves",
	AreaIUCNCategory("Ib"):  "Ib: Wilderness areas",
	AreaIUCNCategory("N/R"): "Not Reported",
	AreaIUCNCategory("V"):   "V: Protected landscape / seascape",
	AreaIUCNCategory("VI"):  "VI: Protected area with sustainable use of natural resources",
}

type AreaOwnerType string

var AreaOwnerTypeDefinitions = map[AreaOwnerType]string{
	AreaOwnerType("DESG"): "Designation",
	AreaOwnerType("DIST"): "Regional Agency Special District",
	AreaOwnerType("FED"):  "Federal",
	AreaOwnerType("JNT"):  "Joint",
	AreaOwnerType("LOC"):  "Local Government",
	AreaOwnerType("NGO"):  "Non-Governmental Organization",
	AreaOwnerType("PVT"):  "Private",
	AreaOwnerType("STAT"): "State",
	AreaOwnerType("TRIB"): "American Indian Lands",
	AreaOwnerType("UNK"):  "Unknown",
}
