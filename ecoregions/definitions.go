package ecoregions

import "strconv"

// Order

type Biome int

var BiomeDefinitions = map[Biome]string{
	Biome(1):  "Tropical & Subtropical Moist Broadleaf Forests",
	Biome(2):  "Tropical & Subtropical Dry Broadleaf Forests",
	Biome(3):  "Tropical & Subtropical Coniferous Forests",
	Biome(4):  "Temperate Broadleaf & Mixed Forests",
	Biome(5):  "Temperate Conifer Forests",
	Biome(6):  "Boreal Forests/Taiga",
	Biome(7):  "Tropical & Subtropical Grasslands, Savannas & Shrublands",
	Biome(8):  "Temperate Grasslands, Savannas & Shrublands",
	Biome(9):  "Flooded Grasslands & Savannas",
	Biome(10): "Montane Grasslands & Shrublands",
	Biome(11): "Tundra",
	Biome(12): "Mediterranean Forests, Woodlands & Scrub",
	Biome(13): "Deserts & Xeric Shrublands",
	Biome(14): "Mangroves",
}

func (Ω Biome) Valid() bool {
	_, ok := BiomeDefinitions[Ω]
	return ok
}

type Realm string //  Biogeographical Realm
var RealmDefinitions = map[Realm]string{
	Realm("AA"): "Australasia",
	Realm("AN"): "Antarctic",
	Realm("AT"): "Afrotropics",
	Realm("IM"): "IndoMalay",
	Realm("NA"): "Nearctic",
	Realm("NT"): "Neotropics",
	Realm("OC"): "Oceania",
	Realm("PA"): "Palearctic",
}

func (Ω Realm) Valid() bool {
	_, ok := RealmDefinitions[Ω]
	return ok
}

type GlobalStatus int

// A 30-year prediction of future conservation status given current conservation status and trajectories.
var GlobalStatusDefinitions = map[GlobalStatus]string{
	GlobalStatus(1): "CRITICAL OR ENDANGERED",
	GlobalStatus(2): "VULNERABLE",
	GlobalStatus(3): "RELATIVELY STABLE OR INTACT",
}

type EcoNum int // A unique number for each ecoregion within each biome nested within each realm.

func (Ω EcoNum) Valid() bool {
	return Ω != 0
}

type EcoID int

func (Ω EcoID) Valid() bool {
	_, ok := EcoIDDefinitions[Ω]
	return ok
}

func (Ω EcoID) Name() string {
	v, _ := EcoIDDefinitions[Ω]
	return v
}

func (Ω EcoID) split() []string {
	if !Ω.Valid() {
		return nil
	}
	s := strconv.Itoa(int(Ω))
	// Pad to ensure the realm isn't throwing us off, though realm should always be five.
	if len(s) != 6 {
		s = "0" + s
	}
	return []string{s[0:2], s[2:4], s[4:6]}
}

//func (Ω EcoID) Realm() Realm {
//	a := Ω.split()
//	if len(a) != 3 {
//		return Realm(0)
//	}
//	i, _ := strconv.Atoi(a[0])
//	return Realm(i)
//}

func (Ω EcoID) Biome() Biome {
	a := Ω.split()
	if len(a) != 3 {
		return Biome(0)
	}
	i, _ := strconv.Atoi(a[1])
	return Biome(i)
}

func (Ω EcoID) EcoNum() EcoNum {
	a := Ω.split()
	if len(a) != 3 {
		return EcoNum(0)
	}
	i, _ := strconv.Atoi(a[2])
	return EcoNum(i)
}

var EcoIDDefinitions = map[EcoID]string{
	EcoID(50513): "Florida sand pine scrub",
	EcoID(51306): "Gulf of California xeric scrub",
	EcoID(50812): "Northern tall grasslands",
	EcoID(50401): "Allegheny Highlands forests",
	EcoID(50813): "Palouse grasslands",
	EcoID(50417): "Willamette Valley forests",
	EcoID(51308): "Mojave desert",
	EcoID(50523): "Piney Woods forests",
	EcoID(51114): "Low Arctic tundra",
	EcoID(51101): "Alaska-St. Elias Range tundra",
	EcoID(50414): "Southern Great Lakes forests",
	EcoID(50406): "Eastern forest-boreal transition",
	EcoID(51106): "Beringia lowland tundra",
	EcoID(50516): "Klamath-Siskiyou forests",
	EcoID(50613): "Northern Cordillera forests",
	EcoID(50512): "Eastern Cascades forests",
	EcoID(51203): "California montane chaparral and woodlands",
	EcoID(51301): "Baja California desert",
	EcoID(50519): "Northern California coastal forests",
	EcoID(51304): "Colorado Plateau shrublands",
	EcoID(50606): "Eastern Canadian Shield taiga",
	EcoID(50804): "Central forest-grasslands transition",
	EcoID(50413): "Southeastern mixed forests",
	EcoID(51103): "Arctic coastal tundra",
	EcoID(50416): "Western Great Lakes forests",
	EcoID(50411): "Northeastern coastal forests",
	EcoID(50808): "Montana Valley and Foothill grasslands",
	EcoID(50617): "Yukon Interior dry forests",
	EcoID(50527): "Sierra Nevada forests",
	EcoID(50409): "Mississippi lowland forests",
	EcoID(50602): "Central Canadian Shield forests",
	EcoID(50302): "Sierra Madre Occidental pine-oak forests",
	EcoID(50505): "Blue Mountains forests",
	EcoID(51309): "Snake-Columbia shrub steppe",
	EcoID(50612): "Northern Canadian Shield taiga",
	EcoID(50810): "Northern mixed grasslands",
	EcoID(51310): "Sonoran desert",
	EcoID(50528): "South Central Rockies forests",
	EcoID(51302): "Central Mexican matorral",
	EcoID(50601): "Alaska Peninsula montane taiga",
	EcoID(51108): "Brooks-British Range tundra",
	EcoID(50501): "Alberta Mountain forests",
	EcoID(50811): "Northern short grasslands",
	EcoID(50524): "Puget lowland forests",
	EcoID(51311): "Tamaulipan matorral",
	EcoID(50605): "Eastern Canadian forests",
	EcoID(50303): "Sierra Madre Oriental pine-oak forests",
	EcoID(50610): "Muskwa-Slave Lake forests",
	EcoID(50609): "Midwestern Canadian Shield forests",
	EcoID(50410): "New England-Acadian forests",
	EcoID(50815): "Western short grasslands",
	EcoID(50511): "Colorado Rockies forests",
	EcoID(51105): "Baffin coastal tundra",
	EcoID(51107): "Beringia upland tundra",
	EcoID(50607): "Interior Alaska-Yukon lowland taiga",
	EcoID(50616): "Southern Hudson Bay taiga",
	EcoID(50506): "British Columbia mainland coastal forests",
	EcoID(51307): "Meseta Central matorral",
	EcoID(50517): "Middle Atlantic coastal forests",
	EcoID(50614): "Northwest Territories taiga",
	EcoID(50518): "North Central Rockies forests",
	EcoID(50301): "Bermuda subtropical conifer forests",
	EcoID(50604): "Copper Plateau taiga",
	EcoID(50809): "Nebraska Sand Hills mixed grasslands",
	EcoID(50801): "California Central Valley grasslands",
	EcoID(51110): "High Arctic tundra",
	EcoID(50608): "Mid-Continental Canadian forests",
	EcoID(50515): "Great Basin montane forests",
	EcoID(50507): "Cascade Mountains leeward forests",
	EcoID(50807): "Flint Hills tall grasslands",
	EcoID(51113): "Kalaallit Nunaat low arctic tundra",
	EcoID(51111): "Interior Yukon-Alaska alpine tundra",
	EcoID(50802): "Canadian Aspen forests and parklands",
	EcoID(51116): "Ogilvie-MacKenzie alpine tundra",
	EcoID(50611): "Newfoundland Highland forests",
	EcoID(51118): "Torngat Mountain tundra",
	EcoID(51104): "Arctic foothills tundra",
	EcoID(51313): "Wyoming Basin shrub steppe",
	EcoID(50510): "Central Pacific coastal forests",
	EcoID(50201): "Sonoran-Sinaloan transition subtropical dry forest",
	EcoID(50522): "Okanagan dry forests",
	EcoID(50806): "Edwards Plateau savanna",
	EcoID(50805): "Central tall grasslands",
	EcoID(50412): "Ozark Mountain forests",
	EcoID(50526): "Sierra Juarez and San Pedro Martir pine-oak forests",
	EcoID(50403): "Appalachian-Blue Ridge forests",
	EcoID(50503): "Arizona Mountains forests",
	EcoID(50525): "Queen Charlotte Islands",
	EcoID(50508): "Central and Southern Cascades forests",
	EcoID(50504): "Atlantic coastal pine barrens",
	EcoID(50407): "Eastern Great Lakes lowland forests",
	EcoID(50502): "Alberta-British Columbia foothills forests",
	EcoID(51303): "Chihuahuan desert",
	EcoID(50405): "East Central Texas forests",
	EcoID(50404): "Central U.S. hardwood forests",
	EcoID(51202): "California interior chaparral and woodlands",
	EcoID(50408): "Gulf of St. Lawrence lowland forests",
	EcoID(50814): "Texas blackland prairies",
	EcoID(51312): "Tamaulipan mezquital",
	EcoID(51109): "Davis Highlands tundra",
	EcoID(51112): "Kalaallit Nunaat high arctic tundra",
	EcoID(50520): "Northern Pacific coastal forests",
	EcoID(51117): "Pacific Coastal Mountain icefields and tundra",
	EcoID(51102): "Aleutian Islands tundra",
	EcoID(50603): "Cook Inlet taiga",
	EcoID(50701): "Western Gulf coastal grasslands",
	EcoID(50803): "Central and Southern mixed grasslands",
	EcoID(51201): "California coastal sage and chaparral",
	EcoID(50509): "Central British Columbia Mountain forests",
	EcoID(51115): "Middle Arctic tundra",
	EcoID(50529): "Southeastern conifer forests",
	EcoID(50530): "Wasatch and Uinta montane forests",
	EcoID(50521): "Northern transitional alpine forests",
	EcoID(50402): "Appalachian mixed mesophytic forests",
	EcoID(50514): "Fraser Plateau and Basin complex",
	EcoID(51305): "Great Basin shrub steppe",
	EcoID(50415): "Upper Midwest forest-savanna transition",
	EcoID(50615): "South Avalon-Burin oceanic barrens",
}

type EcoCode string

var EcoCodeDefinitions = map[EcoCode]string{
	EcoCode("NA0515"): "Great Basin montane forests",
	EcoCode("NA0405"): "East Central Texas forests",
	EcoCode("NA0403"): "Appalachian-Blue Ridge forests",
	EcoCode("NA0813"): "Palouse grasslands",
	EcoCode("NA0514"): "Fraser Plateau and Basin complex",
	EcoCode("NA0507"): "Cascade Mountains leeward forests",
	EcoCode("NA0415"): "Upper Midwest forest-savanna transition",
	EcoCode("NA1309"): "Snake-Columbia shrub steppe",
	EcoCode("NA0613"): "Northern Cordillera forests",
	EcoCode("NA1313"): "Wyoming Basin shrub steppe",
	EcoCode("NA0414"): "Southern Great Lakes forests",
	EcoCode("NA0701"): "Western Gulf coastal grasslands",
	EcoCode("NA0801"): "California Central Valley grasslands",
	EcoCode("NA0408"): "Gulf of St. Lawrence lowland forests",
	EcoCode("NA0417"): "Willamette Valley forests",
	EcoCode("NA1202"): "California interior chaparral and woodlands",
	EcoCode("NA0517"): "Middle Atlantic coastal forests",
	EcoCode("NA0815"): "Western short grasslands",
	EcoCode("NA1110"): "High Arctic tundra",
	EcoCode("NA0513"): "Florida sand pine scrub",
	EcoCode("NA0505"): "Blue Mountains forests",
	EcoCode("NA0201"): "Sonoran-Sinaloan transition subtropical dry forest",
	EcoCode("NA0808"): "Montana Valley and Foothill grasslands",
	EcoCode("NA1104"): "Arctic foothills tundra",
	EcoCode("NA1301"): "Baja California desert",
	EcoCode("NA0406"): "Eastern forest-boreal transition",
	EcoCode("NA0614"): "Northwest Territories taiga",
	EcoCode("NA1115"): "Middle Arctic tundra",
	EcoCode("NA1203"): "California montane chaparral and woodlands",
	EcoCode("NA1305"): "Great Basin shrub steppe",
	EcoCode("NA0525"): "Queen Charlotte Islands",
	EcoCode("NA1105"): "Baffin coastal tundra",
	EcoCode("NA0604"): "Copper Plateau taiga",
	EcoCode("NA0301"): "Bermuda subtropical conifer forests",
	EcoCode("NA0526"): "Sierra Juarez and San Pedro Martir pine-oak forests",
	EcoCode("NA1102"): "Aleutian Islands tundra",
	EcoCode("NA1112"): "Kalaallit Nunaat high arctic tundra",
	EcoCode("NA0302"): "Sierra Madre Occidental pine-oak forests",
	EcoCode("NA0402"): "Appalachian mixed mesophytic forests",
	EcoCode("NA0810"): "Northern mixed grasslands",
	EcoCode("NA0527"): "Sierra Nevada forests",
	EcoCode("NA1304"): "Colorado Plateau shrublands",
	EcoCode("NA1118"): "Torngat Mountain tundra",
	EcoCode("NA0303"): "Sierra Madre Oriental pine-oak forests",
	EcoCode("NA0509"): "Central British Columbia Mountain forests",
	EcoCode("NA0413"): "Southeastern mixed forests",
	EcoCode("NA1114"): "Low Arctic tundra",
	EcoCode("NA0411"): "Northeastern coastal forests",
	EcoCode("NA1108"): "Brooks-British Range tundra",
	EcoCode("NA0610"): "Muskwa-Slave Lake forests",
	EcoCode("NA0616"): "Southern Hudson Bay taiga",
	EcoCode("NA0528"): "South Central Rockies forests",
	EcoCode("NA0523"): "Piney Woods forests",
	EcoCode("NA0804"): "Central forest-grasslands transition",
	EcoCode("NA0404"): "Central U.S. hardwood forests",
	EcoCode("NA0502"): "Alberta-British Columbia foothills forests",
	EcoCode("NA0501"): "Alberta Mountain forests",
	EcoCode("NA0522"): "Okanagan dry forests",
	EcoCode("NA0611"): "Newfoundland Highland forests",
	EcoCode("NA0608"): "Mid-Continental Canadian forests",
	EcoCode("NA1310"): "Sonoran desert",
	EcoCode("NA0603"): "Cook Inlet taiga",
	EcoCode("NA1306"): "Gulf of California xeric scrub",
	EcoCode("NA1312"): "Tamaulipan mezquital",
	EcoCode("NA1109"): "Davis Highlands tundra",
	EcoCode("NA0607"): "Interior Alaska-Yukon lowland taiga",
	EcoCode("NA0401"): "Allegheny Highlands forests",
	EcoCode("NA1308"): "Mojave desert",
	EcoCode("NA0805"): "Central tall grasslands",
	EcoCode("NA0512"): "Eastern Cascades forests",
	EcoCode("NA1307"): "Meseta Central matorral",
	EcoCode("NA0602"): "Central Canadian Shield forests",
	EcoCode("NA0516"): "Klamath-Siskiyou forests",
	EcoCode("NA0809"): "Nebraska Sand Hills mixed grasslands",
	EcoCode("NA0409"): "Mississippi lowland forests",
	EcoCode("NA0524"): "Puget lowland forests",
	EcoCode("NA0519"): "Northern California coastal forests",
	EcoCode("NA1303"): "Chihuahuan desert",
	EcoCode("NA0511"): "Colorado Rockies forests",
	EcoCode("NA1106"): "Beringia lowland tundra",
	EcoCode("NA0504"): "Atlantic coastal pine barrens",
	EcoCode("NA0503"): "Arizona Mountains forests",
	EcoCode("NA1302"): "Central Mexican matorral",
	EcoCode("NA0605"): "Eastern Canadian forests",
	EcoCode("NA0510"): "Central Pacific coastal forests",
	EcoCode("NA0410"): "New England-Acadian forests",
	EcoCode("NA0617"): "Yukon Interior dry forests",
	EcoCode("NA0530"): "Wasatch and Uinta montane forests",
	EcoCode("NA0812"): "Northern tall grasslands",
	EcoCode("NA0615"): "South Avalon-Burin oceanic barrens",
	EcoCode("NA0802"): "Canadian Aspen forests and parklands",
	EcoCode("NA0806"): "Edwards Plateau savanna",
	EcoCode("NA1201"): "California coastal sage and chaparral",
	EcoCode("NA1311"): "Tamaulipan matorral",
	EcoCode("NA0412"): "Ozark Mountain forests",
	EcoCode("NA0612"): "Northern Canadian Shield taiga",
	EcoCode("NA1111"): "Interior Yukon-Alaska alpine tundra",
	EcoCode("NA1101"): "Alaska-St. Elias Range tundra",
	EcoCode("NA0807"): "Flint Hills tall grasslands",
	EcoCode("NA0520"): "Northern Pacific coastal forests",
	EcoCode("NA0803"): "Central and Southern mixed grasslands",
	EcoCode("NA0606"): "Eastern Canadian Shield taiga",
	EcoCode("NA0506"): "British Columbia mainland coastal forests",
	EcoCode("NA0601"): "Alaska Peninsula montane taiga",
	EcoCode("NA0518"): "North Central Rockies forests",
	EcoCode("NA0814"): "Texas blackland prairies",
	EcoCode("NA1113"): "Kalaallit Nunaat low arctic tundra",
	EcoCode("NA0416"): "Western Great Lakes forests",
	EcoCode("NA1103"): "Arctic coastal tundra",
	EcoCode("NA0407"): "Eastern Great Lakes lowland forests",
	EcoCode("NA1116"): "Ogilvie-MacKenzie alpine tundra",
	EcoCode("NA0609"): "Midwestern Canadian Shield forests",
	EcoCode("NA1107"): "Beringia upland tundra",
	EcoCode("NA1117"): "Pacific Coastal Mountain icefields and tundra",
	EcoCode("NA0508"): "Central and Southern Cascades forests",
	EcoCode("NA0529"): "Southeastern conifer forests",
	EcoCode("NA0521"): "Northern transitional alpine forests",
	EcoCode("NA0811"): "Northern short grasslands",
}
