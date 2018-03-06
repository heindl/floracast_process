package gbif

type rank string

const (
	//rankABERRATION = rank("ABERRATION")
	//// Zoological legacy rank
	//rankBIOVAR = rank("BIOVAR")
	//// Microbial rank based on biochemical or physiological properties.
	//rankCHEMOFORM = rank("CHEMOFORM")
	//// Microbial infrasubspecific rank based on chemical constitution.
	//rankCHEMOVAR = rank("CHEMOVAR")
	//// Microbial rank based on production or amount of production of a particular chemical.
	//rankCLASS  = rank("CLASS")
	//rankCOHORT = rank("COHORT")
	//// Sometimes used in zoology, e.g.
	//rankCONVARIETY = rank("CONVARIETY")
	//// A group of cultivars.
	//rankCULTIVAR       = rank("CULTIVAR")
	//rankCULTIVAR_GROUP = rank("CULTIVAR_GROUP")
	//// rank in use from the code for cultivated plants.
	//rankDOMAIN          = rank("DOMAIN")
	//rankFAMILY          = rank("FAMILY")
	//rankFORM            = rank("FORM")
	//rankFORMA_SPECIALIS = rank("FORMA_SPECIALIS")
	//// Microbial infrasubspecific rank.
	rankGenus = rank("GENUS")
	//rankGRANDORDER = rank("GRANDORDER")
	//rankGREX       = rank("GREX")
	//// The term grex has been coined to expand botanical nomenclature to describe hybrids of orchids.
	//rankINFRACLASS        = rank("INFRACLASS")
	//rankINFRACOHORT       = rank("INFRACOHORT")
	//rankINFRAFAMILY       = rank("INFRAFAMILY")
	//rankINFRAGENERIC_NAME = rank("INFRAGENERIC_NAME")
	//// used for any other unspecific rank below genera and above species.
	//rankINFRAGENUS         = rank("INFRAGENUS")
	//rankINFRAKINGDOM       = rank("INFRAKINGDOM")
	//rankINFRALEGION        = rank("INFRALEGION")
	//rankINFRAORDER         = rank("INFRAORDER")
	//rankINFRAPHYLUM        = rank("INFRAPHYLUM")
	//rankINFRASPECIFIC_NAME = rank("INFRASPECIFIC_NAME")
	//// used for any other unspecific rank below species.
	//rankINFRASUBSPECIFIC_NAME = rank("INFRASUBSPECIFIC_NAME")
	//// used also for any other unspecific rank below subspecies.
	//rankINFRATRIBE = rank("INFRATRIBE")
	//rankKINGDOM    = rank("KINGDOM")
	//rankLEGION     = rank("LEGION")
	//// Sometimes used in zoology, e.g.
	//rankMAGNORDER = rank("MAGNORDER")
	//rankMORPH     = rank("MORPH")
	//// Zoological legacy rank
	//rankMORPHOVAR = rank("MORPHOVAR")
	//// Microbial rank based on morphological characterislics.
	//rankNATIO = rank("NATIO")
	//// Zoological legacy rank
	//rankORDER = rank("ORDER")
	//rankOTHER = rank("OTHER")
	//// Any other rank we cannot map to this enumeration
	//rankPARVCLASS = rank("PARVCLASS")
	//rankPARVORDER = rank("PARVORDER")
	//rankPATHOVAR  = rank("PATHOVAR")
	//// Microbial rank based on pathogenic reactions in one or more hosts.
	//rankPHAGOVAR = rank("PHAGOVAR")
	//// Microbial infrasubspecific rank based on reactions to bacteriophage.
	//rankPHYLUM = rank("PHYLUM")
	//rankPROLES = rank("PROLES")
	//// Botanical legacy rank
	//rankRACE = rank("RACE")
	//// Botanical legacy rank
	//rankSECTION = rank("SECTION")
	//rankSERIES  = rank("SERIES")
	//rankSEROVAR = rank("SEROVAR")
	//// Microbial infrasubspecific rank based on antigenic characteristics.
	//rankSPECIES           = rank("SPECIES")
	//rankSPECIES_AGGREGATE = rank("SPECIES_AGGREGATE")
	//// A loosely defined group of species.
	//rankSTRAIN = rank("STRAIN")
	//// A microbial strain.
	//rankSUBCLASS          = rank("SUBCLASS")
	//rankSUBCOHORT         = rank("SUBCOHORT")
	//rankSUBFAMILY         = rank("SUBFAMILY")
	//rankSUBFORM           = rank("SUBFORM")
	//rankSUBGENUS          = rank("SUBGENUS")
	//rankSUBKINGDOM        = rank("SUBKINGDOM")
	//rankSUBLEGION         = rank("SUBLEGION")
	//rankSUBORDER          = rank("SUBORDER")
	//rankSUBPHYLUM         = rank("SUBPHYLUM")
	//rankSUBSECTION        = rank("SUBSECTION")
	//rankSUBSERIES         = rank("SUBSERIES")
	//rankSUBSPECIES        = rank("SUBSPECIES")
	//rankSUBTRIBE          = rank("SUBTRIBE")
	//rankSUBVARIETY        = rank("SUBVARIETY")
	//rankSUPERCLASS        = rank("SUPERCLASS")
	//rankSUPERCOHORT       = rank("SUPERCOHORT")
	//rankSUPERFAMILY       = rank("SUPERFAMILY")
	//rankSUPERKINGDOM      = rank("SUPERKINGDOM")
	//rankSUPERLEGION       = rank("SUPERLEGION")
	//rankSUPERORDER        = rank("SUPERORDER")
	//rankSUPERPHYLUM       = rank("SUPERPHYLUM")
	//rankSUPERTRIBE        = rank("SUPERTRIBE")
	//rankSUPRAGENERIC_NAME = rank("SUPRAGENERIC_NAME")
	//// Used for any other unspecific rank above genera.
	//rankTRIBE    = rank("TRIBE")
	rankUnranked = rank("UNRANKED")
	//rankVARIETY  = rank("VARIETY")
)
