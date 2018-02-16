package ecoregions

import "strconv"

// Order

type Biome int
type EcoCode string

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

type Realm int

func (Ω Realm) Valid() bool {
	return Ω > 0 && Ω < 9
}

const (
	RealmAustralasia Realm = iota + 1
	RealmAntarctic
	RealmAfrotropics
	RealmIndoMalay
	RealmNearctic
	RealmNeotropics
	RealmOceania
	RealmPalearctic
)

type RealmCode string //  Biogeographical realm
var RealmCodeDefinitions = map[RealmCode]Realm{
	RealmCode("AA"): RealmAustralasia,
	RealmCode("AN"): RealmAntarctic,
	RealmCode("AT"): RealmAfrotropics,
	RealmCode("IM"): RealmIndoMalay,
	RealmCode("NA"): RealmNearctic,
	RealmCode("NT"): RealmNeotropics,
	RealmCode("OC"): RealmOceania,
	RealmCode("PA"): RealmPalearctic,
}

func (Ω RealmCode) Valid() bool {
	_, ok := RealmCodeDefinitions[Ω]
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
	_, isNeartic := NearticEcoIDDefinitions[Ω]
	_, isNeotropic := NeotropicEcoIDDefinitions[Ω]
	return isNeartic || isNeotropic
}

func (Ω EcoID) Name() string {
	if v, ok := NearticEcoIDDefinitions[Ω]; ok {
		return v
	}
	if v, ok := NeotropicEcoIDDefinitions[Ω]; ok {
		return v
	}
	return ""
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

func (Ω EcoID) Realm() Realm {
	a := Ω.split()
	if len(a) != 3 {
		return Realm(0)
	}
	i, _ := strconv.Atoi(a[0])
	return Realm(i)
}

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

var NearticEcoIDDefinitions = map[EcoID]string{
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

var NeotropicEcoIDDefinitions = map[EcoID]string{
	EcoID(61404): "Northern Mesoamerican Pacific mangroves",
	EcoID(60301): "Bahamian pine mosaic",
	EcoID(60124): "Guianan Highlands moist forests",
	EcoID(61402): "Bahamian-Antillean mangroves",
	EcoID(60129): "Isthmian-Atlantic moist forests",
	EcoID(60216): "Islas Revillagigedo dry forests",
	EcoID(60213): "Cuban dry forests",
	EcoID(61314): "San Lucan xeric scrub",
	EcoID(60169): "Pantepui",
	EcoID(61403): "Mesoamerican Gulf-Caribbean mangroves",
	EcoID(60127): "Hispaniolan moist forests",
	EcoID(60164): "South Florida rocklands",
	EcoID(61407): "Southern Mesoamerican Pacific mangroves",
	EcoID(60709): "Llanos",
	EcoID(60904): "Everglades",
	EcoID(61305): "Caribbean shrublands",
	EcoID(60402): "Magellanic subpolar forests",
	EcoID(61309): "La Costa xeric shrublands",
	EcoID(61401): "Amazon-Orinoco-Southern Caribbean mangroves",
	EcoID(61006): "Northern Andean páramo",
	EcoID(60805): "Patagonian steppe",
	EcoID(60148): "Pantanos de Centla",
	EcoID(60112): "Central American montane forests",
	EcoID(60303): "Central American pine-oak forests",
	EcoID(61405): "South American Pacific mangroves",
	EcoID(60309): "Sierra Madre del Sur pine-oak forests",
	EcoID(60158): "Rio Negro campinarana",
	EcoID(60201): "Apure-Villavicencio dry forests",
	EcoID(60167): "Talamancan montane forests",
	EcoID(60117): "Cordillera La Costa montane forests",
	EcoID(60310): "Trans-Mexican Volcanic Belt pine-oak forests",
	EcoID(61306): "Cuban cactus scrub",
	EcoID(60302): "Belizian pine forests",
	EcoID(61308): "Guajira-Barranquilla xeric scrub",
	EcoID(60906): "Orinoco wetlands",
	EcoID(60120): "Cuban moist forests",
	EcoID(60134): "Leeward Islands moist forests",
	EcoID(60902): "Cuban wetlands",
	EcoID(60205): "Balsas dry forests",
	EcoID(60128): "Iquitos varzeá",
	EcoID(60182): "Guianan piedmont and lowland moist forests",
	EcoID(60179): "Windward Islands moist forests",
	EcoID(60107): "Caqueta moist forests",
	EcoID(60306): "Miskito pine forests",
	EcoID(60209): "Central American dry forests",
	EcoID(60903): "Enriquillo wetlands",
	EcoID(60146): "Oaxacan montane forests",
	EcoID(60133): "Juruá-Purus moist forests",
	EcoID(60404): "Valdivian temperate forests",
	EcoID(60307): "Sierra de la Laguna pine-oak forests",
	EcoID(60220): "Lesser Antillean dry forests",
	EcoID(60218): "Jamaican dry forests",
	EcoID(60181): "Yucatán moist forests",
	EcoID(60131): "Jamaican moist forests",
	EcoID(60707): "Guianan savanna",
	EcoID(60171): "Trinidad and Tobago moist forests",
	EcoID(60155): "Puerto Rican moist forests",
	EcoID(60113): "Chiapas montane forests",
	EcoID(60154): "Petén-Veracruz moist forests",
	EcoID(60308): "Sierra Madre de Oaxaca pine-oak forests",
	EcoID(61004): "Cordillera Central páramo",
	EcoID(60217): "Jalisco dry forests",
	EcoID(60229): "Sinú Valley dry forests",
	EcoID(60176): "Veracruz moist forests",
	EcoID(60226): "Puerto Rican dry forests",
	EcoID(60111): "Central American Atlantic moist forests",
	EcoID(60110): "Cayos Miskitos-San Andrés and Providencia moist forests",
	EcoID(60149): "Guianan freshwater swamp forests",
	EcoID(61301): "Araya and Paria xeric scrub",
	EcoID(60130): "Isthmian-Pacific moist forests",
	EcoID(60125): "Guianan moist forests",
	EcoID(61316): "Tehuacán Valley matorral",
	EcoID(60115): "Chocó-Darién moist forests",
	EcoID(60147): "Orinoco Delta swamp forests",
	EcoID(60305): "Hispaniolan pine forests",
	EcoID(60138): "Marajó varzeá",
	EcoID(60114): "Chimalapas montane forests",
	EcoID(60102): "Atlantic Coast restingas",
	EcoID(60227): "Sierra de la Laguna dry forests",
	EcoID(60304): "Cuban pine forests",
	EcoID(60156): "Purus varzeá",
	EcoID(60215): "Hispaniolan dry forests",
	EcoID(60123): "Fernando de Noronha-Atol das Rocas moist forests",
	EcoID(61406): "Southern Atlantic mangroves",
	EcoID(60703): "Campos Rupestres montane savanna",
	EcoID(60907): "Pantanal",
	EcoID(61307): "Galápagos Islands scrubland mosaic",
	EcoID(60165): "Southern Andean Yungas",
	EcoID(60104): "Bahia interior forests",
	EcoID(60214): "Ecuadorian dry forests",
	EcoID(60704): "Cerrado",
	EcoID(60150): "Alto Paraná Atlantic forests",
	EcoID(60119): "Costa Rican seasonal moist forests",
	EcoID(60180): "Xingu-Tocantins-Araguaia moist forests",
	EcoID(60159): "Santa Marta montane forests",
	EcoID(61005): "Cordillera de Merida páramo",
	EcoID(60206): "Bolivian montane dry forests",
	EcoID(60122): "Eastern Panamanian montane forests",
	EcoID(60108): "Catatumbo moist forests",
	EcoID(61303): "Atacama desert",
	EcoID(60177): "Veracruz montane forests",
	EcoID(60235): "Yucatán dry forests",
	EcoID(60224): "Panamanian dry forests",
	EcoID(60173): "Uatuma-Trombetas moist forests",
	EcoID(60116): "Cocos Island moist forests",
	EcoID(60219): "Lara-Falcón dry forests",
	EcoID(61007): "Santa Marta páramo",
	EcoID(60228): "Sinaloan dry forests",
	EcoID(61304): "Caatinga",
	EcoID(60233): "Veracruz dry forests",
	EcoID(60702): "Beni savanna",
	EcoID(60166): "Southwest Amazon moist forests",
	EcoID(61315): "Sechura desert",
	EcoID(61003): "Central Andean wet puna",
	EcoID(60106): "Caatinga Enclaves moist forests",
	EcoID(60141): "Monte Alegre varzeá",
	EcoID(60140): "Mato Grosso seasonal forests",
	EcoID(60152): "Pernambuco interior forests",
	EcoID(60126): "Gurupa varzeá",
	EcoID(60144): "Northeastern Brazil restingas",
	EcoID(60161): "Sierra de los Tuxtlas",
	EcoID(60109): "Cauca Valley montane forests",
	EcoID(61312): "Motagua Valley thornscrub",
	EcoID(60157): "Purus-Madeira moist forests",
	EcoID(60204): "Bajío dry forests",
	EcoID(60137): "Magdalena-Urabá moist forests",
	EcoID(61311): "Malpelo Island xeric scrub",
	EcoID(60174): "Ucayali moist forests",
	EcoID(60211): "Chiapas Depression dry forests",
	EcoID(60139): "Maranhão Babaçu forests",
	EcoID(60143): "Negro-Branco moist forests",
	EcoID(61313): "Paraguana xeric scrub",
	EcoID(60222): "Maracaibo dry forests",
	EcoID(60170): "Tocantins/Pindare moist forests",
	EcoID(60162): "Sierra Madre de Chiapas moist forests",
	EcoID(60905): "Guayaquil flooded grasslands",
	EcoID(60207): "Cauca Valley dry forests",
	EcoID(60168): "Tapajós-Xingu moist forests",
	EcoID(60210): "Dry Chaco",
	EcoID(60135): "Madeira-Tapajós moist forests",
	EcoID(60151): "Pernambuco coastal forests",
	EcoID(60202): "Atlantic dry forests",
	EcoID(60221): "Magdalena Valley dry forests",
	EcoID(60232): "Tumbes-Piura dry forests",
	EcoID(61010): "High Monte",
	EcoID(60153): "Peruvian Yungas",
	EcoID(60223): "Marañón dry forests",
	EcoID(60801): "Espinal",
	EcoID(60103): "Bahia coastal forests",
	EcoID(60160): "Serra do Mar coastal forests",
	EcoID(60708): "Humid Chaco",
	EcoID(60101): "Araucaria moist forests",
	EcoID(60909): "Southern Cone Mesopotamian savanna",
	EcoID(60710): "Uruguayan savanna",
	EcoID(61201): "Chilean matorral",
	EcoID(61008): "Southern Andean steppe",
	EcoID(60802): "Low Monte",
	EcoID(60803): "Humid Pampas",
	EcoID(61001): "Central Andean dry puna",
	EcoID(60401): "Juan Fernández Islands temperate forests",
	EcoID(60908): "Paraná flooded savanna",
	EcoID(61002): "Central Andean puna",
	EcoID(60403): "San Félix-San Ambrosio Islands temperate forests",
	EcoID(60132): "Japurá-Solimoes-Negro moist forests",
	EcoID(60163): "Solimões-Japurá moist forests",
	EcoID(60212): "Chiquitano dry forests",
	EcoID(60225): "Patía Valley dry forests",
	EcoID(60105): "Bolivian Yungas",
	EcoID(60230): "Southern Pacific dry forests",
	EcoID(60175): "Venezuelan Andes montane forests",
	EcoID(60136): "Magdalena Valley montane forests",
	EcoID(60178): "Western Ecuador moist forests",
	EcoID(60118): "Cordillera Oriental montane forests",
	EcoID(60145): "Northwestern Andean montane forests",
	EcoID(61318): "St. Peter and St. Paul rocks",
	EcoID(60705): "Clipperton Island shrub and grasslands",
	EcoID(60121): "Eastern Cordillera real montane forests",
	EcoID(60142): "Napo moist forests",
	EcoID(60172): "Trindade-Martin Vaz Islands tropical forests",
}