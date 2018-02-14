package utils

import "fmt"

func GetMorchellaAggregateTestData() []byte {
  return []byte(fmt.Sprintf(`{"eM3R8X2YQyLJWiLMVIGzZaU1I": %s}`, morchella))
}

func GetMorchellaUsageTestData() []byte {
  return []byte(morchella)
}

var morchella = `{
    "CanonicalName": {
      "ScientificName": "morchella esculenta",
      "Rank": "species"
    },
    "Sources": {
      "11": {
        "ELEMENT_GLOBAL.2.121953": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta",
            "Rank": "species"
          },
          "Synonyms": [
            {
              "ScientificName": "morchella conica",
              "Rank": "species"
            }
          ],
          "CommonNames": [
            "common morel",
            "sponge morel",
            "yellow morel"
          ],
          "ModifiedAt": "2018-02-11T18:58:09.261401-05:00",
          "CreatedAt": "2018-02-11T18:58:09.261399-05:00"
        }
      },
      "27": {
        "2594602": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta",
            "Rank": "species"
          },
          "Occurrences": 204,
          "LastFetchedAt": "2018-02-13T16:18:21.965554-05:00",
          "ModifiedAt": "2018-02-11T18:57:47.281841-05:00",
          "CreatedAt": "2018-02-11T18:57:47.281839-05:00"
        },
        "2594603": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta umbrina",
            "Rank": "variety"
          },
          "Occurrences": 1,
          "LastFetchedAt": "2018-02-13T16:17:48.412803-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577832-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577831-05:00"
        },
        "2594604": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella umbrina",
            "Rank": "species"
          },
          "Occurrences": 1,
          "LastFetchedAt": "2018-02-13T16:18:02.839718-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577947-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577946-05:00"
        },
        "2594605": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morellus esculentus",
            "Rank": "species"
          },
          "LastFetchedAt": "2018-02-13T16:18:01.504723-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.578-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577999-05:00"
        },
        "2594606": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella rotunda",
            "Rank": "species"
          },
          "LastFetchedAt": "2018-02-13T16:17:44.746656-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577878-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577877-05:00"
        },
        "2594607": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella rotunda esculenta",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:17:50.804324-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.57791-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577909-05:00"
        },
        "2594609": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "phallus esculentus",
            "Rank": "species"
          },
          "LastFetchedAt": "2018-02-13T16:18:26.82774-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.57802-05:00",
          "CreatedAt": "2018-02-11T18:57:48.578018-05:00"
        },
        "2594610": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta rotunda",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:18:01.637398-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577811-05:00",
          "CreatedAt": "2018-02-11T18:57:48.57781-05:00"
        },
        "2594611": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta esculenta",
            "Rank": "subspecies"
          },
          "LastFetchedAt": "2018-02-13T16:17:47.641199-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.57772-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577719-05:00"
        },
        "2594617": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella conica",
            "Rank": "species"
          },
          "Occurrences": 40,
          "LastFetchedAt": "2018-02-13T16:18:38.643026-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577516-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577515-05:00"
        },
        "3314945": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "phallus tremelloides",
            "Rank": "species"
          },
          "LastFetchedAt": "2018-02-13T16:17:52.58985-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.578043-05:00",
          "CreatedAt": "2018-02-11T18:57:48.578042-05:00"
        },
        "3495583": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella distans longissima",
            "Rank": "form"
          },
          "LastFetchedAt": "2018-02-13T16:17:59.787348-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577652-05:00",
          "CreatedAt": "2018-02-11T18:57:48.57765-05:00"
        },
        "3495586": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella rotunda alba",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:17:46.988087-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577894-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577893-05:00"
        },
        "3495595": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta alba",
            "Rank": "form"
          },
          "LastFetchedAt": "2018-02-13T16:17:59.937664-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577698-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577697-05:00"
        },
        "3495599": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella rigida",
            "Rank": "species"
          },
          "LastFetchedAt": "2018-02-13T16:17:52.186688-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577874-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577872-05:00"
        },
        "3495604": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella rotunda crassipes",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:17:48.541629-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577906-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577905-05:00"
        },
        "3495608": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella conica pusilla",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:17:55.459473-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577584-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577582-05:00"
        },
        "3495618": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta stipitata",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:18:30.866582-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577824-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577823-05:00"
        },
        "3495623": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta atrotomentosa",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:18:00.175077-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577744-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577742-05:00"
        },
        "3495625": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella rotunda cinerea",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:17:46.578518-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577898-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577897-05:00"
        },
        "3495629": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella rotunda rigida",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:18:22.361901-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577931-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577929-05:00"
        },
        "3495630": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta grisea",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:18:38.803506-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577775-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577774-05:00"
        },
        "3495638": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta pubescens",
            "Rank": "subspecies"
          },
          "LastFetchedAt": "2018-02-13T16:17:43.775737-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577724-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577723-05:00"
        },
        "3495643": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella conica angusticeps",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:17:56.619727-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577537-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577536-05:00"
        },
        "3495653": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella conica serotina",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:18:30.426972-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577602-05:00",
          "CreatedAt": "2018-02-11T18:57:48.5776-05:00"
        },
        "3495654": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta viridis",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:17:48.697426-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577845-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577844-05:00"
        },
        "3495655": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta lutescens",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:17:55.625537-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577779-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577778-05:00"
        },
        "3495659": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta abietina",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:17:45.29127-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577728-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577727-05:00"
        },
        "3495662": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella umbrina macroalveola",
            "Rank": "form"
          },
          "LastFetchedAt": "2018-02-13T16:17:52.892613-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577952-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577951-05:00"
        },
        "3495664": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella tremelloides",
            "Rank": "species"
          },
          "LastFetchedAt": "2018-02-13T16:17:50.94858-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577938-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577937-05:00"
        },
        "3495665": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella distans distans",
            "Rank": "form"
          },
          "LastFetchedAt": "2018-02-13T16:17:52.318174-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577645-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577644-05:00"
        },
        "3495667": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella distans spathulata",
            "Rank": "form"
          },
          "LastFetchedAt": "2018-02-13T16:17:44.887984-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577673-05:00",
          "CreatedAt": "2018-02-11T18:57:48.57767-05:00"
        },
        "3495671": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta corrugata",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:17:56.994997-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577759-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577758-05:00"
        },
        "3495673": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella dunensis sterilis",
            "Rank": "form"
          },
          "LastFetchedAt": "2018-02-13T16:18:27.532906-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577694-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577693-05:00"
        },
        "3495675": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella dunensis sterile",
            "Rank": "form"
          },
          "LastFetchedAt": "2018-02-13T16:17:51.084598-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.57769-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577689-05:00"
        },
        "3495676": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella umbrina umbrina",
            "Rank": "form"
          },
          "LastFetchedAt": "2018-02-13T16:18:38.949888-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577957-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577956-05:00"
        },
        "3495679": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta fulva",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:17:52.454011-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577771-05:00",
          "CreatedAt": "2018-02-11T18:57:48.57777-05:00"
        },
        "3495683": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella lutescens",
            "Rank": "species"
          },
          "LastFetchedAt": "2018-02-13T16:17:45.436557-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577854-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577853-05:00"
        },
        "3495684": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta rigida",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:17:45.625043-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.5778-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577798-05:00"
        },
        "3495691": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella viridis",
            "Rank": "species"
          },
          "LastFetchedAt": "2018-02-13T16:18:29.870241-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577961-05:00",
          "CreatedAt": "2018-02-11T18:57:48.57796-05:00"
        },
        "3495699": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella rotunda pallida",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:18:01.771477-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577922-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577921-05:00"
        },
        "3495700": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella abietina",
            "Rank": "species"
          },
          "LastFetchedAt": "2018-02-13T16:17:46.711539-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577512-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577511-05:00"
        },
        "3495706": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella pubescens",
            "Rank": "species"
          },
          "LastFetchedAt": "2018-02-13T16:17:51.403869-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.57787-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577868-05:00"
        },
        "3495709": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella conica cylindrica",
            "Rank": "form"
          },
          "LastFetchedAt": "2018-02-13T16:18:02.972675-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577527-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577526-05:00"
        },
        "3495712": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella prunarii",
            "Rank": "species"
          },
          "LastFetchedAt": "2018-02-13T16:17:57.418788-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577862-05:00",
          "CreatedAt": "2018-02-11T18:57:48.57786-05:00"
        },
        "3495713": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella vulgaris vulgaris",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:17:56.019451-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577996-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577995-05:00"
        },
        "3495715": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella cylindrica",
            "Rank": "species"
          },
          "LastFetchedAt": "2018-02-13T16:18:27.924425-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577613-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577612-05:00"
        },
        "3495717": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella conica metheformis",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:18:27.672878-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577573-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577572-05:00"
        },
        "3495719": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta dunensis",
            "Rank": "form"
          },
          "LastFetchedAt": "2018-02-13T16:18:22.495495-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577702-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577701-05:00"
        },
        "3495721": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella conica rigida",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:17:56.170593-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577596-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577595-05:00"
        },
        "3495725": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella rotunda minutula",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:17:50.516484-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577918-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577917-05:00"
        },
        "3495726": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella vulgaris parvula",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:17:45.022128-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577988-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577987-05:00"
        },
        "3495727": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta albida",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:18:00.760382-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.57774-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577738-05:00"
        },
        "3495730": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta prunarii",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:18:01.088112-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577791-05:00",
          "CreatedAt": "2018-02-11T18:57:48.57779-05:00"
        },
        "3495734": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella vulgaris albida",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:18:31.005944-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.57798-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577979-05:00"
        },
        "3495740": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta conica",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:17:51.548035-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577755-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577754-05:00"
        },
        "6284792": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta vulgaris",
            "Rank": "variety"
          },
          "Occurrences": 2,
          "LastFetchedAt": "2018-02-13T16:18:00.3932-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.57785-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577849-05:00"
        },
        "6357642": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta pubescens",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:18:03.112598-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577796-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577795-05:00"
        },
        "6357649": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella conica conica",
            "Rank": "form"
          },
          "LastFetchedAt": "2018-02-13T16:18:22.699853-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.57752-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577519-05:00"
        },
        "7081015": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "phallus esculentus esculentus",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:17:43.915805-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.578031-05:00",
          "CreatedAt": "2018-02-11T18:57:48.57803-05:00"
        },
        "7258325": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella rotunda rotunda",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:17:45.770674-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577934-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577933-05:00"
        },
        "7258328": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella conica conica",
            "Rank": "subspecies"
          },
          "LastFetchedAt": "2018-02-13T16:18:31.669886-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577533-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577532-05:00"
        },
        "7258329": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella distans",
            "Rank": "species"
          },
          "LastFetchedAt": "2018-02-13T16:17:44.54227-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577635-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577633-05:00"
        },
        "7258343": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella dunensis",
            "Rank": "species"
          },
          "LastFetchedAt": "2018-02-13T16:17:46.296643-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577679-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577677-05:00"
        },
        "7258344": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella dunensis dunensis",
            "Rank": "form"
          },
          "LastFetchedAt": "2018-02-13T16:17:47.121769-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577685-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577683-05:00"
        },
        "7258347": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella rotunda rotunda",
            "Rank": "form"
          },
          "LastFetchedAt": "2018-02-13T16:17:57.739206-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577886-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577885-05:00"
        },
        "7349747": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta albida",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:17:50.654226-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577736-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577735-05:00"
        },
        "7402568": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta violacea",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:18:30.008543-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577842-05:00",
          "CreatedAt": "2018-02-11T18:57:48.57784-05:00"
        },
        "7434880": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morilla esculenta",
            "Rank": "species"
          },
          "LastFetchedAt": "2018-02-13T16:17:56.3762-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.578007-05:00",
          "CreatedAt": "2018-02-11T18:57:48.578006-05:00"
        },
        "7436343": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "phalloboletus esculentus",
            "Rank": "species"
          },
          "LastFetchedAt": "2018-02-13T16:18:01.223095-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.578016-05:00",
          "CreatedAt": "2018-02-11T18:57:48.578015-05:00"
        },
        "7448333": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta rubroris",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:17:45.156002-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.57782-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577819-05:00"
        },
        "7474391": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella ovalis pallida",
            "Rank": "form"
          },
          "LastFetchedAt": "2018-02-13T16:18:30.575236-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577858-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577857-05:00"
        },
        "7482360": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta mahoniae",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:17:51.711448-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577783-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577782-05:00"
        },
        "7517103": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella rotunda pubescens",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:17:53.220385-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577926-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577925-05:00"
        },
        "7534204": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella conica cilicicae",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:18:01.903633-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577545-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577543-05:00"
        },
        "7552760": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "phallus esculentus rotundus",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:18:24.382789-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.578039-05:00",
          "CreatedAt": "2018-02-11T18:57:48.578038-05:00"
        },
        "7563370": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella rotunda alba",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:17:53.405993-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.57789-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577889-05:00"
        },
        "7564823": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella conica ceracea",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:18:30.152507-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577541-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577539-05:00"
        },
        "7635914": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta dunensis",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:17:48.968403-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577764-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577763-05:00"
        },
        "7654589": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella conica pygmaea",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:18:31.148126-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.57759-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577588-05:00"
        },
        "7671808": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella conica distans",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:18:23.00542-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577556-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577555-05:00"
        },
        "7696772": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "helvella esculenta",
            "Rank": "species"
          },
          "LastFetchedAt": "2018-02-13T16:18:01.364365-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577507-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577506-05:00"
        },
        "7742106": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella vulgaris tremelloides",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:18:28.532115-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577992-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577991-05:00"
        },
        "7743604": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella rotunda cinerea",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:17:46.443288-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577902-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577901-05:00"
        },
        "7754457": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella conica meandriformis",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:17:51.851759-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577569-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577568-05:00"
        },
        "7755971": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella vulgaris cinerascens",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:18:02.322502-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577984-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577983-05:00"
        },
        "7801825": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella conica flexuosa",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:18:28.67357-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577565-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577564-05:00"
        },
        "7826978": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella conica crassa",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:18:23.951498-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577553-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577552-05:00"
        },
        "7844394": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella conica nigra",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:17:49.107575-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577577-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577576-05:00"
        },
        "7863363": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta umbrinoides",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:17:46.15481-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577838-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577836-05:00"
        },
        "7903027": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella conica conica",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:18:24.095257-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577549-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577548-05:00"
        },
        "7922204": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella rotunda fulva",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:17:52.725953-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577914-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577913-05:00"
        },
        "7936814": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella conica violeipes",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:17:47.778763-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577608-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577606-05:00"
        },
        "7965505": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta esculenta",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:17:53.548868-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577768-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577766-05:00"
        },
        "7968619": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella vulgaris alba",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:17:51.22701-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577972-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577971-05:00"
        },
        "7986342": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta aurantiaca",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:17:47.914614-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577748-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577746-05:00"
        },
        "8022920": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta sterilis",
            "Rank": "form"
          },
          "LastFetchedAt": "2018-02-13T16:18:30.727376-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577716-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577715-05:00"
        },
        "8040650": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morilla conica",
            "Rank": "species"
          },
          "LastFetchedAt": "2018-02-13T16:17:53.063628-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.578004-05:00",
          "CreatedAt": "2018-02-11T18:57:48.578002-05:00"
        },
        "8071344": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella vulgaris",
            "Rank": "species"
          },
          "LastFetchedAt": "2018-02-13T16:17:50.378508-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577968-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577967-05:00"
        },
        "8072850": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella conica elata",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:18:02.487688-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.57756-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577559-05:00"
        },
        "8106927": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "phallus esculentus albus",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:17:48.831156-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.578024-05:00",
          "CreatedAt": "2018-02-11T18:57:48.578022-05:00"
        },
        "8148999": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella pubescens",
            "Rank": "species"
          },
          "LastFetchedAt": "2018-02-13T16:18:30.287117-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577865-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577864-05:00"
        },
        "8210888": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella vulgaris alba",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:18:02.640343-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577976-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577975-05:00"
        },
        "8259908": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morilla tremelloides",
            "Rank": "species"
          },
          "LastFetchedAt": "2018-02-13T16:18:28.257111-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.578012-05:00",
          "CreatedAt": "2018-02-11T18:57:48.578011-05:00"
        },
        "8271210": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta ovalis",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:18:23.667598-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577787-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577786-05:00"
        },
        "8272775": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta roseostraminea",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:18:00.615719-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577808-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577806-05:00"
        },
        "8281396": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta alba",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:18:28.809517-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577732-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577731-05:00"
        },
        "8298968": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta rigida",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:17:53.734753-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577803-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577802-05:00"
        },
        "8307722": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella tremelloides",
            "Rank": "species"
          },
          "LastFetchedAt": "2018-02-13T16:17:47.278225-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577943-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577942-05:00"
        },
        "8309385": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta brunnea",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:18:24.237766-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577751-05:00",
          "CreatedAt": "2018-02-11T18:57:48.57775-05:00"
        },
        "8320714": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "phallus esculentus cinereus",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:17:46.849416-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.578027-05:00",
          "CreatedAt": "2018-02-11T18:57:48.578026-05:00"
        },
        "8347054": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta rotunda",
            "Rank": "form"
          },
          "LastFetchedAt": "2018-02-13T16:17:51.989151-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577712-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577711-05:00"
        },
        "8348709": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta esculenta",
            "Rank": "form"
          },
          "LastFetchedAt": "2018-02-13T16:17:57.877996-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577708-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577707-05:00"
        },
        "8365335": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta theobromichroa",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:18:29.731406-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577828-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577826-05:00"
        },
        "8398362": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "phallus esculentus fuscus",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:18:23.800078-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.578035-05:00",
          "CreatedAt": "2018-02-11T18:57:48.578034-05:00"
        },
        "8574619": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta",
            "Rank": "species"
          },
          "LastFetchedAt": "2018-02-13T16:17:44.049915-05:00",
          "ModifiedAt": "2018-02-11T18:57:47.398261-05:00",
          "CreatedAt": "2018-02-11T18:57:47.398259-05:00"
        },
        "8586978": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella rotunda",
            "Rank": "species"
          },
          "LastFetchedAt": "2018-02-13T16:18:28.394001-05:00",
          "ModifiedAt": "2018-02-11T18:57:47.455536-05:00",
          "CreatedAt": "2018-02-11T18:57:47.455533-05:00"
        },
        "8693950": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella rotunda",
            "Rank": "species"
          },
          "LastFetchedAt": "2018-02-13T16:17:59.629427-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577882-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577881-05:00"
        },
        "8976256": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella vulgaris",
            "Rank": "species"
          },
          "LastFetchedAt": "2018-02-13T16:18:26.671627-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577965-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577963-05:00"
        },
        "9225951": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta rotunda",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:18:28.120427-05:00",
          "ModifiedAt": "2018-02-11T18:57:48.577816-05:00",
          "CreatedAt": "2018-02-11T18:57:48.577815-05:00"
        }
      },
      "INAT": {
        "206086": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta umbrina",
            "Rank": "variety"
          },
          "LastFetchedAt": "2018-02-13T16:18:39.437376-05:00",
          "ModifiedAt": "2018-02-11T18:57:46.971267-05:00",
          "CreatedAt": "2018-02-11T18:57:46.971265-05:00"
        },
        "58682": {
          "TaxonomicReference": true,
          "CanonicalName": {
            "ScientificName": "morchella esculenta",
            "Rank": "species"
          },
          "CommonNames": [
            "morel"
          ],
          "Occurrences": 120,
          "LastFetchedAt": "2018-02-13T16:18:43.45229-05:00",
          "ModifiedAt": "2018-02-11T18:57:46.977884-05:00",
          "CreatedAt": "2018-02-11T18:57:46.977882-05:00"
        }
      }
    },
    "CreatedAt": "2018-02-11T18:57:46.971291-05:00",
    "ModifiedAt": "2018-02-11T18:58:09.279379-05:00"
  }`