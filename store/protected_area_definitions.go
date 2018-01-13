package store

type AreaGAPStatus string

var AreaGAPStatusDefinitions = map[AreaGAPStatus]string{
	AreaGAPStatus("1"): "1 - managed for biodiversity - disturbance events proceed or are mimicked",
	AreaGAPStatus("2"): "2 - managed for biodiversity - disturbance events suppressed",
}

type AreaCategory string

var AreaCategoryDefinitions = map[AreaCategory]string{
	AreaCategory("Designation"): "Designation",
	AreaCategory("Easement"):    "Easement",
	AreaCategory("Fee"):         "Fee",
	AreaCategory("Other"):       "Other",
	AreaCategory("Unknown"):     "Unknown",
}

type AreaDesignation string

var AreaDesignationDefinitions = map[AreaDesignation]string{
	AreaDesignation("ACC"):   "Access Area",
	AreaDesignation("ACEC"):  "Area of Critical Environmental Concern",
	AreaDesignation("CONE"):  "Conservation Easement",
	AreaDesignation("FOTH"):  "Federal Other or Unknown",
	AreaDesignation("HCA"):   "Historic or Cultural Area",
	AreaDesignation("LCA"):   "Local Conservation Area",
	AreaDesignation("LHCA"):  "Local Historic or Cultural Area",
	AreaDesignation("LOTH"):  "Local Other or Unknown",
	AreaDesignation("LP"):    "Local Park",
	AreaDesignation("LREC"):  "Local Recreation Area",
	AreaDesignation("LRMA"):  "Local Resource Management Area",
	AreaDesignation("MIL"):   "Military Land",
	AreaDesignation("MIT"):   "Mitigation Land or Bank",
	AreaDesignation("MPA"):   "Marine Protected Area",
	AreaDesignation("NCA"):   "Conservation Area",
	AreaDesignation("ND"):    "Not Designated",
	AreaDesignation("NF"):    "National Forest",
	AreaDesignation("NM"):    "National Monument or Landmark",
	AreaDesignation("NP"):    "National Park",
	AreaDesignation("NRA"):   "National Recreation Area",
	AreaDesignation("NSBV"):  "National Scenic, Botanical or Volcanic Area",
	AreaDesignation("NT"):    "National Scenic or Historic Trail",
	AreaDesignation("NWR"):   "National Wildlife Refuge",
	AreaDesignation("OTHE"):  "Other Easement",
	AreaDesignation("PAGR"):  "Private Agricultural",
	AreaDesignation("PCON"):  "Private Conservation",
	AreaDesignation("PFOR"):  "Private Forest Stewardship",
	AreaDesignation("PHCA"):  "Private Historic or Cultural",
	AreaDesignation("POTH"):  "Private Other or Unknown",
	AreaDesignation("PREC"):  "Private Recreation or Education",
	AreaDesignation("PROC"):  "Approved or Proclamation Boundary",
	AreaDesignation("REA"):   "Research or Educational Area",
	AreaDesignation("REC"):   "Recreation Management Area",
	AreaDesignation("RECE"):  "Recreation or Education Easement",
	AreaDesignation("RMA"):   "Resource Management Area",
	AreaDesignation("RNA"):   "Research Natural Area",
	AreaDesignation("SCA"):   "State Conservation Area",
	AreaDesignation("SDA"):   "Special Designation Area",
	AreaDesignation("SHCA"):  "State Historic or Cultural Area",
	AreaDesignation("SOTH"):  "State Other or Unknown",
	AreaDesignation("SP"):    "State Park",
	AreaDesignation("SREC"):  "State Recreation Area",
	AreaDesignation("SRMA"):  "State Resource Management Area",
	AreaDesignation("SW"):    "State Wilderness",
	AreaDesignation("TRIBL"): "Native American Land",
	AreaDesignation("UNK"):   "Unknown",
	AreaDesignation("WA"):    "Wilderness Area",
	AreaDesignation("WPA"):   "Watershed Protection Area",
	AreaDesignation("WSA"):   "Wilderness Study Area",
	AreaDesignation("WSR"):   "Wild and Scenic River",
}

type AreaManagerName string

var AreaManagerNameDefinitions = map[AreaManagerName]string{
	AreaManagerName("ARS"):   "Agricultural Research Service",
	AreaManagerName("BIA"):   "Bureau of Indian Affairs",
	AreaManagerName("BLM"):   "Bureau of Land Management",
	AreaManagerName("CITY"):  "City Land",
	AreaManagerName("CNTY"):  "County Land",
	AreaManagerName("DOD"):   "Department of Defense",
	AreaManagerName("FWS"):   "U.S. Fish & Wildlife Service",
	AreaManagerName("JNT"):   "Joint",
	AreaManagerName("NGO"):   "Non-Governmental Organization",
	AreaManagerName("NOAA"):  "National Oceanic and Atmospheric Administration",
	AreaManagerName("NPS"):   "National Park Service",
	AreaManagerName("NRCS"):  "Natural Resources Conservation Service",
	AreaManagerName("OTHF"):  "Other or Unknown Federal Land",
	AreaManagerName("OTHR"):  "Other",
	AreaManagerName("OTHS"):  "Other or Unknown State Land",
	AreaManagerName("PVT"):   "Private",
	AreaManagerName("REG"):   "Regional Agency Land",
	AreaManagerName("RWD"):   "Regional Water Districts",
	AreaManagerName("SDC"):   "State Department of Conservation",
	AreaManagerName("SDNR"):  "State Department of Natural Resources",
	AreaManagerName("SDOL"):  "State Department of Land",
	AreaManagerName("SFW"):   "State Fish and Wildlife",
	AreaManagerName("SLB"):   "State Land Board",
	AreaManagerName("SPR"):   "State Park and Recreation",
	AreaManagerName("TVA"):   "Tennessee Valley Authority",
	AreaManagerName("UNK"):   "Unknown",
	AreaManagerName("UNKL"):  "Other or Unknown Local Government",
	AreaManagerName("USACE"): "Army Corps of Engineers",
	AreaManagerName("USBR"):  "Bureau of Reclamation",
	AreaManagerName("USFS"):  "Forest Service",
}

type AreaManagerType string

var AreaManagerTypeDefinitions = map[AreaManagerType]string{
	AreaManagerType("DIST"): "Regional Agency Special District",
	AreaManagerType("FED"):  "Federal",
	AreaManagerType("JNT"):  "Joint",
	AreaManagerType("LOC"):  "Local Government",
	AreaManagerType("NGO"):  "Non-Governmental Organization",
	AreaManagerType("PVT"):  "Private",
	AreaManagerType("STAT"): "State",
	AreaManagerType("UNK"):  "Unknown",
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
