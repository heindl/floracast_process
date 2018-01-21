package inaturalist

//import "bitbucket.org/heindl/taxa/store"

type Taxon struct {
	CommonName       string             `firestore:",omitempty"`
	ScientificName    string      `firestore:",omitempty"`
	ParentID         INaturalistTaxonID `firestore:",omitempty"`
	PhotoURL         string             `firestore:",omitempty"`
	Rank             TaxonRank          `firestore:",omitempty"`
	RankLevel        RankLevel          `firestore:",omitempty"`
	ModifiedAt       time.Time          `firestore:",omitempty"`
	CreatedAt        time.Time          `firestore:",omitempty"`
	States           []State            `firestore:",omitempty"`
	WikipediaSummary string             `firestore:",omitempty"`
	EcoRegions       map[string]int     `firestore:",omitempty"`
}


func ParseINaturalistTaxon(taxa []*INaturalistTaxon) store.Taxa {

	for _, txn := range taxa {
		txn.CurrentSynonymousTaxonIds
		txn := store.Taxon{

		}
	}

}