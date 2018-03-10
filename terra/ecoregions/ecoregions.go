package ecoregions

import (
	"errors"
	dropboxError "github.com/dropbox/godropbox/errors"
	"strconv"
)

// Region contains all necessary information about a World Wildlife Fund ecological region.
type Region struct {
	realm Realm
	name  string
	//ecoCode      ecoCode
	ecoNum EcoNum
	ecoID  ecoID
	biome  Biome
}

// Name of the region EcoID, or a combination of Realm, Biome and EcoNum.
func (Ω *Region) Name() string {
	return Ω.name
}

// Realm returns the id of one of seven major habitats.
func (Ω *Region) Realm() Realm {
	return Ω.realm
}

// Biome returns the id of one of 14 habitats
func (Ω *Region) Biome() Biome {
	return Ω.biome
}

// EcoNum returns a unique number for each ecoregion within each biome nested within each realm.
func (Ω *Region) EcoNum() EcoNum {
	return Ω.ecoNum
}

// ErrNotFound flags when a coordinate is not within a cached EcoRegion
var ErrNotFound = errors.New("EcoRegion Not Found")

// NewRegion creates and validates a new Region structure.
// Returns error if all fields are not Valid.
func NewRegion(ecoIDInt int) (*Region, error) {

	region := Region{
		ecoID: ecoID(ecoIDInt),
	}
	if !region.ecoID.Valid() {
		return nil, dropboxError.Newf("Invalid ecoID [%ecoIDStr]", region.ecoID)
	}
	var err error
	region.name, err = region.ecoID.name()
	if err != nil {
		return nil, err
	}

	ecoIDStr := strconv.Itoa(ecoIDInt)
	// Pad to ensure the realm isn't throwing us off, though realm should always be five.
	if len(ecoIDStr) != 6 {
		ecoIDStr = "0" + ecoIDStr
	}

	realmInt, err := strconv.Atoi(ecoIDStr[0:2])
	if err != nil {
		return nil, dropboxError.Wrap(err, "Invalid Realm")
	}
	region.realm = Realm(realmInt)
	if !region.realm.Valid() {
		return nil, dropboxError.Newf("Invalid Realm [%d]", region.realm)
	}

	biomeInt, err := strconv.Atoi(ecoIDStr[2:4])
	if err != nil {
		return nil, dropboxError.Wrap(err, "Invalid Biome")
	}
	region.biome = Biome(biomeInt)
	if !region.biome.Valid() {
		return nil, dropboxError.Newf("Invalid Biome [%d]", region.biome)
	}

	ecoNumInt, err := strconv.Atoi(ecoIDStr[4:6])
	if err != nil {
		return nil, dropboxError.Wrap(err, "Invalid EcoNum")
	}
	region.ecoNum = EcoNum(ecoNumInt)

	return &region, nil

}

// EcoCode is an alphanumeric code that is similar to eco_ID but a little easier to interpret.
// The first 2 characters (letters) are the realm the ecoregion is in.
// The 2nd 2 characters are the biome and the last 2 characters are the ecoregion number.
//type ecoCode string

// Biome is a large naturally occurring community of flora and fauna occupying a major habitat.
type Biome int

var biomeDefinitions = map[Biome]string{
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

// Valid checks the biome against a known list.
func (Ω Biome) Valid() bool {
	_, ok := biomeDefinitions[Ω]
	return ok
}

// Realm is one of seven major habitats: Afrotropical, Australasia, Indo-Malayan, Nearctic, Neotropical, Oceania, Palearctic
type Realm int

// Valid checks the realm against as a known value.
func (Ω Realm) Valid() bool {
	return Ω > 0 && Ω < 9
}

const (
	realmAustralasia Realm = iota + 1
	realmAntarctic
	realmAfrotropics
	realmIndoMalay
	realmNearctic
	realmNeotropics
	realmOceania
	realmPalearctic
)

type realmCode string //  Biogeographical realm
var realmCodeDefinitions = map[realmCode]Realm{
	realmCode("AA"): realmAustralasia,
	realmCode("AN"): realmAntarctic,
	realmCode("AT"): realmAfrotropics,
	realmCode("IM"): realmIndoMalay,
	realmCode("NA"): realmNearctic,
	realmCode("NT"): realmNeotropics,
	realmCode("OC"): realmOceania,
	realmCode("PA"): realmPalearctic,
}

func (Ω realmCode) Valid() bool {
	_, ok := realmCodeDefinitions[Ω]
	return ok
}

//type globalStatus int
//// A 30-year prediction of future conservation status given current conservation status and trajectories.
//var globalStatusDefinitions = map[globalStatus]string{
//	globalStatus(1): "CRITICAL OR ENDANGERED",
//	globalStatus(2): "VULNERABLE",
//	globalStatus(3): "RELATIVELY STABLE OR INTACT",
//}

// EcoNum is a unique number for each ecoregion within each biome nested within each realm.
type EcoNum int

// Valid checks that the EcoNum is not empty
func (Ω EcoNum) Valid() bool {
	return Ω != 0
}

// ecoID is created by combining realm, BIOME, and ECO_NUM, thus creating a unique numeric ID for each ecoregion.
type ecoID int

func (Ω ecoID) Valid() bool {
	_, isNeartic := nearticEcoIDDefinitions[Ω]
	_, isNeotropic := neotropicEcoIDDefinitions[Ω]
	return isNeartic || isNeotropic
}

func (Ω ecoID) name() (string, error) {
	if v, ok := nearticEcoIDDefinitions[Ω]; ok {
		return v, nil
	}
	if v, ok := neotropicEcoIDDefinitions[Ω]; ok {
		return v, nil
	}
	return "", dropboxError.Newf("Name not found [%s]", Ω)
}

var nearticEcoIDDefinitions = map[ecoID]string{
	ecoID(50513): "Florida sand pine scrub",
	ecoID(51306): "Gulf of California xeric scrub",
	ecoID(50812): "Northern tall grasslands",
	ecoID(50401): "Allegheny Highlands forests",
	ecoID(50813): "Palouse grasslands",
	ecoID(50417): "Willamette Valley forests",
	ecoID(51308): "Mojave desert",
	ecoID(50523): "Piney Woods forests",
	ecoID(51114): "Low Arctic tundra",
	ecoID(51101): "Alaska-St. Elias Range tundra",
	ecoID(50414): "Southern Great Lakes forests",
	ecoID(50406): "Eastern forest-boreal transition",
	ecoID(51106): "Beringia lowland tundra",
	ecoID(50516): "Klamath-Siskiyou forests",
	ecoID(50613): "Northern Cordillera forests",
	ecoID(50512): "Eastern Cascades forests",
	ecoID(51203): "California montane chaparral and woodlands",
	ecoID(51301): "Baja California desert",
	ecoID(50519): "Northern California coastal forests",
	ecoID(51304): "Colorado Plateau shrublands",
	ecoID(50606): "Eastern Canadian Shield taiga",
	ecoID(50804): "Central forest-grasslands transition",
	ecoID(50413): "Southeastern mixed forests",
	ecoID(51103): "Arctic coastal tundra",
	ecoID(50416): "Western Great Lakes forests",
	ecoID(50411): "Northeastern coastal forests",
	ecoID(50808): "Montana Valley and Foothill grasslands",
	ecoID(50617): "Yukon Interior dry forests",
	ecoID(50527): "Sierra Nevada forests",
	ecoID(50409): "Mississippi lowland forests",
	ecoID(50602): "Central Canadian Shield forests",
	ecoID(50302): "Sierra Madre Occidental pine-oak forests",
	ecoID(50505): "Blue Mountains forests",
	ecoID(51309): "Snake-Columbia shrub steppe",
	ecoID(50612): "Northern Canadian Shield taiga",
	ecoID(50810): "Northern mixed grasslands",
	ecoID(51310): "Sonoran desert",
	ecoID(50528): "South Central Rockies forests",
	ecoID(51302): "Central Mexican matorral",
	ecoID(50601): "Alaska Peninsula montane taiga",
	ecoID(51108): "Brooks-British Range tundra",
	ecoID(50501): "Alberta Mountain forests",
	ecoID(50811): "Northern short grasslands",
	ecoID(50524): "Puget lowland forests",
	ecoID(51311): "Tamaulipan matorral",
	ecoID(50605): "Eastern Canadian forests",
	ecoID(50303): "Sierra Madre Oriental pine-oak forests",
	ecoID(50610): "Muskwa-Slave Lake forests",
	ecoID(50609): "Midwestern Canadian Shield forests",
	ecoID(50410): "New England-Acadian forests",
	ecoID(50815): "Western short grasslands",
	ecoID(50511): "Colorado Rockies forests",
	ecoID(51105): "Baffin coastal tundra",
	ecoID(51107): "Beringia upland tundra",
	ecoID(50607): "Interior Alaska-Yukon lowland taiga",
	ecoID(50616): "Southern Hudson Bay taiga",
	ecoID(50506): "British Columbia mainland coastal forests",
	ecoID(51307): "Meseta Central matorral",
	ecoID(50517): "Middle Atlantic coastal forests",
	ecoID(50614): "Northwest Territories taiga",
	ecoID(50518): "North Central Rockies forests",
	ecoID(50301): "Bermuda subtropical conifer forests",
	ecoID(50604): "Copper Plateau taiga",
	ecoID(50809): "Nebraska Sand Hills mixed grasslands",
	ecoID(50801): "California Central Valley grasslands",
	ecoID(51110): "High Arctic tundra",
	ecoID(50608): "Mid-Continental Canadian forests",
	ecoID(50515): "Great Basin montane forests",
	ecoID(50507): "Cascade Mountains leeward forests",
	ecoID(50807): "Flint Hills tall grasslands",
	ecoID(51113): "Kalaallit Nunaat low arctic tundra",
	ecoID(51111): "Interior Yukon-Alaska alpine tundra",
	ecoID(50802): "Canadian Aspen forests and parklands",
	ecoID(51116): "Ogilvie-MacKenzie alpine tundra",
	ecoID(50611): "Newfoundland Highland forests",
	ecoID(51118): "Torngat Mountain tundra",
	ecoID(51104): "Arctic foothills tundra",
	ecoID(51313): "Wyoming Basin shrub steppe",
	ecoID(50510): "Central Pacific coastal forests",
	ecoID(50201): "Sonoran-Sinaloan transition subtropical dry forest",
	ecoID(50522): "Okanagan dry forests",
	ecoID(50806): "Edwards Plateau savanna",
	ecoID(50805): "Central tall grasslands",
	ecoID(50412): "Ozark Mountain forests",
	ecoID(50526): "Sierra Juarez and San Pedro Martir pine-oak forests",
	ecoID(50403): "Appalachian-Blue Ridge forests",
	ecoID(50503): "Arizona Mountains forests",
	ecoID(50525): "Queen Charlotte Islands",
	ecoID(50508): "Central and Southern Cascades forests",
	ecoID(50504): "Atlantic coastal pine barrens",
	ecoID(50407): "Eastern Great Lakes lowland forests",
	ecoID(50502): "Alberta-British Columbia foothills forests",
	ecoID(51303): "Chihuahuan desert",
	ecoID(50405): "East Central Texas forests",
	ecoID(50404): "Central U.S. hardwood forests",
	ecoID(51202): "California interior chaparral and woodlands",
	ecoID(50408): "Gulf of St. Lawrence lowland forests",
	ecoID(50814): "Texas blackland prairies",
	ecoID(51312): "Tamaulipan mezquital",
	ecoID(51109): "Davis Highlands tundra",
	ecoID(51112): "Kalaallit Nunaat high arctic tundra",
	ecoID(50520): "Northern Pacific coastal forests",
	ecoID(51117): "Pacific Coastal Mountain icefields and tundra",
	ecoID(51102): "Aleutian Islands tundra",
	ecoID(50603): "Cook Inlet taiga",
	ecoID(50701): "Western Gulf coastal grasslands",
	ecoID(50803): "Central and Southern mixed grasslands",
	ecoID(51201): "California coastal sage and chaparral",
	ecoID(50509): "Central British Columbia Mountain forests",
	ecoID(51115): "Middle Arctic tundra",
	ecoID(50529): "Southeastern conifer forests",
	ecoID(50530): "Wasatch and Uinta montane forests",
	ecoID(50521): "Northern transitional alpine forests",
	ecoID(50402): "Appalachian mixed mesophytic forests",
	ecoID(50514): "Fraser Plateau and Basin complex",
	ecoID(51305): "Great Basin shrub steppe",
	ecoID(50415): "Upper Midwest forest-savanna transition",
	ecoID(50615): "South Avalon-Burin oceanic barrens",
}

var neotropicEcoIDDefinitions = map[ecoID]string{
	ecoID(61404): "Northern Mesoamerican Pacific mangroves",
	ecoID(60301): "Bahamian pine mosaic",
	ecoID(60124): "Guianan Highlands moist forests",
	ecoID(61402): "Bahamian-Antillean mangroves",
	ecoID(60129): "Isthmian-Atlantic moist forests",
	ecoID(60216): "Islas Revillagigedo dry forests",
	ecoID(60213): "Cuban dry forests",
	ecoID(61314): "San Lucan xeric scrub",
	ecoID(60169): "Pantepui",
	ecoID(61403): "Mesoamerican Gulf-Caribbean mangroves",
	ecoID(60127): "Hispaniolan moist forests",
	ecoID(60164): "South Florida rocklands",
	ecoID(61407): "Southern Mesoamerican Pacific mangroves",
	ecoID(60709): "Llanos",
	ecoID(60904): "Everglades",
	ecoID(61305): "Caribbean shrublands",
	ecoID(60402): "Magellanic subpolar forests",
	ecoID(61309): "La Costa xeric shrublands",
	ecoID(61401): "Amazon-Orinoco-Southern Caribbean mangroves",
	ecoID(61006): "Northern Andean páramo",
	ecoID(60805): "Patagonian steppe",
	ecoID(60148): "Pantanos de Centla",
	ecoID(60112): "Central American montane forests",
	ecoID(60303): "Central American pine-oak forests",
	ecoID(61405): "South American Pacific mangroves",
	ecoID(60309): "Sierra Madre del Sur pine-oak forests",
	ecoID(60158): "Rio Negro campinarana",
	ecoID(60201): "Apure-Villavicencio dry forests",
	ecoID(60167): "Talamancan montane forests",
	ecoID(60117): "Cordillera La Costa montane forests",
	ecoID(60310): "Trans-Mexican Volcanic Belt pine-oak forests",
	ecoID(61306): "Cuban cactus scrub",
	ecoID(60302): "Belizian pine forests",
	ecoID(61308): "Guajira-Barranquilla xeric scrub",
	ecoID(60906): "Orinoco wetlands",
	ecoID(60120): "Cuban moist forests",
	ecoID(60134): "Leeward Islands moist forests",
	ecoID(60902): "Cuban wetlands",
	ecoID(60205): "Balsas dry forests",
	ecoID(60128): "Iquitos varzeá",
	ecoID(60182): "Guianan piedmont and lowland moist forests",
	ecoID(60179): "Windward Islands moist forests",
	ecoID(60107): "Caqueta moist forests",
	ecoID(60306): "Miskito pine forests",
	ecoID(60209): "Central American dry forests",
	ecoID(60903): "Enriquillo wetlands",
	ecoID(60146): "Oaxacan montane forests",
	ecoID(60133): "Juruá-Purus moist forests",
	ecoID(60404): "Valdivian temperate forests",
	ecoID(60307): "Sierra de la Laguna pine-oak forests",
	ecoID(60220): "Lesser Antillean dry forests",
	ecoID(60218): "Jamaican dry forests",
	ecoID(60181): "Yucatán moist forests",
	ecoID(60131): "Jamaican moist forests",
	ecoID(60707): "Guianan savanna",
	ecoID(60171): "Trinidad and Tobago moist forests",
	ecoID(60155): "Puerto Rican moist forests",
	ecoID(60113): "Chiapas montane forests",
	ecoID(60154): "Petén-Veracruz moist forests",
	ecoID(60308): "Sierra Madre de Oaxaca pine-oak forests",
	ecoID(61004): "Cordillera Central páramo",
	ecoID(60217): "Jalisco dry forests",
	ecoID(60229): "Sinú Valley dry forests",
	ecoID(60176): "Veracruz moist forests",
	ecoID(60226): "Puerto Rican dry forests",
	ecoID(60111): "Central American Atlantic moist forests",
	ecoID(60110): "Cayos Miskitos-San Andrés and Providencia moist forests",
	ecoID(60149): "Guianan freshwater swamp forests",
	ecoID(61301): "Araya and Paria xeric scrub",
	ecoID(60130): "Isthmian-Pacific moist forests",
	ecoID(60125): "Guianan moist forests",
	ecoID(61316): "Tehuacán Valley matorral",
	ecoID(60115): "Chocó-Darién moist forests",
	ecoID(60147): "Orinoco Delta swamp forests",
	ecoID(60305): "Hispaniolan pine forests",
	ecoID(60138): "Marajó varzeá",
	ecoID(60114): "Chimalapas montane forests",
	ecoID(60102): "Atlantic Coast restingas",
	ecoID(60227): "Sierra de la Laguna dry forests",
	ecoID(60304): "Cuban pine forests",
	ecoID(60156): "Purus varzeá",
	ecoID(60215): "Hispaniolan dry forests",
	ecoID(60123): "Fernando de Noronha-Atol das Rocas moist forests",
	ecoID(61406): "Southern Atlantic mangroves",
	ecoID(60703): "Campos Rupestres montane savanna",
	ecoID(60907): "Pantanal",
	ecoID(61307): "Galápagos Islands scrubland mosaic",
	ecoID(60165): "Southern Andean Yungas",
	ecoID(60104): "Bahia interior forests",
	ecoID(60214): "Ecuadorian dry forests",
	ecoID(60704): "Cerrado",
	ecoID(60150): "Alto Paraná Atlantic forests",
	ecoID(60119): "Costa Rican seasonal moist forests",
	ecoID(60180): "Xingu-Tocantins-Araguaia moist forests",
	ecoID(60159): "Santa Marta montane forests",
	ecoID(61005): "Cordillera de Merida páramo",
	ecoID(60206): "Bolivian montane dry forests",
	ecoID(60122): "Eastern Panamanian montane forests",
	ecoID(60108): "Catatumbo moist forests",
	ecoID(61303): "Atacama desert",
	ecoID(60177): "Veracruz montane forests",
	ecoID(60235): "Yucatán dry forests",
	ecoID(60224): "Panamanian dry forests",
	ecoID(60173): "Uatuma-Trombetas moist forests",
	ecoID(60116): "Cocos Island moist forests",
	ecoID(60219): "Lara-Falcón dry forests",
	ecoID(61007): "Santa Marta páramo",
	ecoID(60228): "Sinaloan dry forests",
	ecoID(61304): "Caatinga",
	ecoID(60233): "Veracruz dry forests",
	ecoID(60702): "Beni savanna",
	ecoID(60166): "Southwest Amazon moist forests",
	ecoID(61315): "Sechura desert",
	ecoID(61003): "Central Andean wet puna",
	ecoID(60106): "Caatinga Enclaves moist forests",
	ecoID(60141): "Monte Alegre varzeá",
	ecoID(60140): "Mato Grosso seasonal forests",
	ecoID(60152): "Pernambuco interior forests",
	ecoID(60126): "Gurupa varzeá",
	ecoID(60144): "Northeastern Brazil restingas",
	ecoID(60161): "Sierra de los Tuxtlas",
	ecoID(60109): "Cauca Valley montane forests",
	ecoID(61312): "Motagua Valley thornscrub",
	ecoID(60157): "Purus-Madeira moist forests",
	ecoID(60204): "Bajío dry forests",
	ecoID(60137): "Magdalena-Urabá moist forests",
	ecoID(61311): "Malpelo Island xeric scrub",
	ecoID(60174): "Ucayali moist forests",
	ecoID(60211): "Chiapas Depression dry forests",
	ecoID(60139): "Maranhão Babaçu forests",
	ecoID(60143): "Negro-Branco moist forests",
	ecoID(61313): "Paraguana xeric scrub",
	ecoID(60222): "Maracaibo dry forests",
	ecoID(60170): "Tocantins/Pindare moist forests",
	ecoID(60162): "Sierra Madre de Chiapas moist forests",
	ecoID(60905): "Guayaquil flooded grasslands",
	ecoID(60207): "Cauca Valley dry forests",
	ecoID(60168): "Tapajós-Xingu moist forests",
	ecoID(60210): "Dry Chaco",
	ecoID(60135): "Madeira-Tapajós moist forests",
	ecoID(60151): "Pernambuco coastal forests",
	ecoID(60202): "Atlantic dry forests",
	ecoID(60221): "Magdalena Valley dry forests",
	ecoID(60232): "Tumbes-Piura dry forests",
	ecoID(61010): "High Monte",
	ecoID(60153): "Peruvian Yungas",
	ecoID(60223): "Marañón dry forests",
	ecoID(60801): "Espinal",
	ecoID(60103): "Bahia coastal forests",
	ecoID(60160): "Serra do Mar coastal forests",
	ecoID(60708): "Humid Chaco",
	ecoID(60101): "Araucaria moist forests",
	ecoID(60909): "Southern Cone Mesopotamian savanna",
	ecoID(60710): "Uruguayan savanna",
	ecoID(61201): "Chilean matorral",
	ecoID(61008): "Southern Andean steppe",
	ecoID(60802): "Low Monte",
	ecoID(60803): "Humid Pampas",
	ecoID(61001): "Central Andean dry puna",
	ecoID(60401): "Juan Fernández Islands temperate forests",
	ecoID(60908): "Paraná flooded savanna",
	ecoID(61002): "Central Andean puna",
	ecoID(60403): "San Félix-San Ambrosio Islands temperate forests",
	ecoID(60132): "Japurá-Solimoes-Negro moist forests",
	ecoID(60163): "Solimões-Japurá moist forests",
	ecoID(60212): "Chiquitano dry forests",
	ecoID(60225): "Patía Valley dry forests",
	ecoID(60105): "Bolivian Yungas",
	ecoID(60230): "Southern Pacific dry forests",
	ecoID(60175): "Venezuelan Andes montane forests",
	ecoID(60136): "Magdalena Valley montane forests",
	ecoID(60178): "Western Ecuador moist forests",
	ecoID(60118): "Cordillera Oriental montane forests",
	ecoID(60145): "Northwestern Andean montane forests",
	ecoID(61318): "St. Peter and St. Paul rocks",
	ecoID(60705): "Clipperton Island shrub and grasslands",
	ecoID(60121): "Eastern Cordillera real montane forests",
	ecoID(60142): "Napo moist forests",
	ecoID(60172): "Trindade-Martin Vaz Islands tropical forests",
}
