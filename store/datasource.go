package store

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/saleswise/errors/errors"
	"strings"
	"time"
	"google.golang.org/api/iterator"
)

type DataSourceID string

const (
	DataSourceIDGBIF        = DataSourceID("27")
	DataSourceIDINaturalist = DataSourceID("INAT")
	DataSourceIDMushroomObserver = DataSourceID("MUOB")
	DataSourceNatureServe = DataSourceID("11")
)

func (Ω DataSourceID) Valid() bool {
	_, ok := SchemeSourceIDMap[Ω]
	return ok
}

var SchemeSourceIDMap = map[DataSourceID]string{
	// Floracast
	DataSourceIDINaturalist: "iNaturalist",
	// INaturalist
	DataSourceID("1"):  "IUCN Red List of Threatened Species. Version 2012.1",
	DataSourceID("3"):  "Amphibian Species of the World 5.6",
	DataSourceID("2"):  "Amphibiaweb. 2012",
	DataSourceID("5"):  "Amphibian Species of the World 5.5",
	DataSourceID("17"): "New England Wild Flower Society's Flora Novae Angliae",
	DataSourceNatureServe: "NatureServe Explorer: An online encyclopedia of life. Version 7.1",
	DataSourceID("12"): "Calflora",
	DataSourceID("13"): "Odonata Central",
	DataSourceID("14"): "IUCN Red List of Threatened Species. Version 2012.2",
	DataSourceID("10"): "eBird/Clements Checklist 6.7",
	DataSourceID("15"): "CONABIO",
	DataSourceID("6"):  "The Reptile Database",
	DataSourceID("16"): "Afribats",
	DataSourceID("18"): "Norma 059, 2010",
	DataSourceID("4"):  "Draft IUCN/SSC, 2013.1",
	DataSourceID("19"): "Draft IUCN/SSC Amphibian Specialist Group, 2011",
	DataSourceID("20"): "eBird/Clements Checklist 6.8",
	DataSourceID("21"): "IUCN Red List of Threatened Species. Version 2013.2",
	DataSourceID("22"): "eBird/Clements Checklist 6.9",
	DataSourceID("23"): "NatureWatch NZ",
	DataSourceID("24"): "The world spider catalog, version 15.5",
	DataSourceID("25"): "Carabidae of the World",
	DataSourceID("26"): "IUCN Red List of Threatened Species. Version 2014.3",
	DataSourceIDGBIF:   "GBIF",
	DataSourceID("28"): "NPSpecies",
	DataSourceID("29"): "Esslinger&#39;s North American Lichens",
	DataSourceID("30"): "Amphibian Species of the World 6.0",
	DataSourceID("31"): "Esslinger&#39;s North American Lichens, Version 21",
	DataSourceIDMushroomObserver: "MushroomObserver.org",
}

type DataSourceKind string

const (
	DataSourceKindOccurrence  = DataSourceKind("occurrence")
	DataSourceKindPhoto       = DataSourceKind("photo")
	DataSourceKindDescription = DataSourceKind("description")
)

func (Ω DataSourceKind) Valid() bool {
	return Ω == DataSourceKindOccurrence || Ω == DataSourceKindPhoto || Ω == DataSourceKindDescription
}

type DataSourceTargetID string

type DataSourceTargetIDs []DataSourceTargetID

func (Ω DataSourceTargetIDs) Strings() (res []string) {
	for _, id := range Ω {
		res = append(res, string(id))
	}
	return
}

func (Ω DataSourceTargetIDs) AddToSet(ids ...DataSourceTargetID) DataSourceTargetIDs {
	for _, id := range ids {
		if Ω.Contains(id) {
			continue
		}
		Ω = append(Ω, id)
	}
	return Ω
}

func (Ω DataSourceTargetIDs) Contains(id DataSourceTargetID) bool {
	for i := range Ω {
		if Ω[i] == id {
			return true
		}
	}
	return false
}

type OccurrenceSource struct {
	SourceID DataSourceID   `firestore:",omitempty"`
	TargetID      DataSourceTargetID `firestore:",omitempty"`
	CanonicalName       string `firestore:",omitempty"`
	CreatedAt     *time.Time         `firestore:",omitempty"`
	ModifiedAt    *time.Time         `firestore:",omitempty"`
	LastFetchedAt *time.Time         `firestore:",omitempty"`
}

type DataSource struct {
	Kind     DataSourceKind `firestore:",omitempty"`
	SourceID DataSourceID   `firestore:",omitempty"`

	TargetID      DataSourceTargetID `firestore:",omitempty"`
	TaxonID       INaturalistTaxonID `firestore:",omitempty"`
	CreatedAt     *time.Time         `firestore:",omitempty"`
	ModifiedAt    *time.Time         `firestore:",omitempty"`
	LastFetchedAt *time.Time         `firestore:",omitempty"`
}

var DataSourceFieldsToMerge = []firestore.FieldPath{
	firestore.FieldPath{"Kind"},
	firestore.FieldPath{"SourceID"},
	firestore.FieldPath{"TargetID"},
	firestore.FieldPath{"INaturalistTaxonID"},
	firestore.FieldPath{"ModifiedAt"},
	firestore.FieldPath{"CreatedAt"},
}

func (Ω *DataSource) Validate() error {

	if !Ω.TaxonID.Valid() {
		return errors.New("invalid taxon id")
	}

	if Ω.TargetID == "" {
		return errors.New("invalid target id")
	}

	if !Ω.SourceID.Valid() {
		return errors.Newf("invalid source id %s", Ω.SourceID)
	}

	return nil
}
func (Ω *store) GetSourceLastCreated(cxt context.Context, kind DataSourceKind, srcID DataSourceID) (*time.Time, error) {
	iter := Ω.FirestoreClient.Collection(CollectionTypeDataSources).
		Where("Kind", "==", kind).
			Where("SourceID", "==", srcID).
				OrderBy("CreatedAt", firestore.Desc).
					Limit(1).
						Documents(cxt)

	for {
		ref, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, errors.Wrap(err, "could no fetch latest source")
		}

		src := DataSource{}
		if err := ref.DataTo(&src); err != nil {
			return nil, errors.Wrap(err, "could not cast source")
		}

		return src.ModifiedAt, nil
	}

	return nil, nil
}


func (Ω *store) GetOccurrenceDataSources(context context.Context, taxonID INaturalistTaxonID) (res DataSources, err error) {

	q := Ω.FirestoreClient.Collection(CollectionTypeDataSources).
		Where("Kind", "==", DataSourceKindOccurrence)

	if taxonID.Valid() {
		q = q.Where("INaturalistTaxonID", "==", taxonID)
	}

	snaps, err := q.Documents(context).GetAll()
	if err != nil {
		return nil, errors.Wrapf(err, "could not fetch occurrence data sources for taxon id[%s]", taxonID)
	}

	for _, s := range snaps {
		var src DataSource
		if err := s.DataTo(&src); err != nil {
			return nil, errors.Wrap(err, "could not typecast data source")
		}
		res = append(res, src)
	}

	return
}

func (Ω *store) NewDataSourceDocumentRef(taxonID INaturalistTaxonID, dataSourceID DataSourceID, targetID DataSourceTargetID, kind DataSourceKind) (*firestore.DocumentRef, error) {

	if !taxonID.Valid() {
		return nil, errors.New("invalid data source document reference id")
	}
	if !dataSourceID.Valid() {
		return nil, errors.New("invalid data source id")
	}
	if !kind.Valid() {
		return nil, errors.New("invalid kind")
	}
	if targetID == "" {
		return nil, errors.New("invalid target id")
	}

	return Ω.FirestoreClient.Collection(CollectionTypeDataSources).
		Doc(fmt.Sprintf("%s|%s|%s|%s", string(taxonID), dataSourceID, targetID, kind)), nil

}

func (Ω *store) UpsertDataSource(cxt context.Context, src DataSource) error {

	ref, err := Ω.NewDataSourceDocumentRef(src.TaxonID, src.SourceID, src.TargetID, src.Kind)
	if err != nil {
		return err
	}

	if err := Ω.FirestoreClient.RunTransaction(cxt, func(cxt context.Context, tx *firestore.Transaction) error {
		if _, err := tx.Get(ref); err != nil {
			if strings.Contains(err.Error(), "not found") {
				return tx.Set(ref, src)
			} else {
				return err
			}
		}
		return tx.Set(ref, src, firestore.Merge(DataSourceFieldsToMerge...))
	}); err != nil {
		return errors.Wrap(err, "could not update data source")
	}
	return nil
}

func (Ω *store) UpdateDataSourceLastFetched(cxt context.Context, src DataSource) error {

	ref, err := Ω.NewDataSourceDocumentRef(src.TaxonID, src.SourceID, src.TargetID, src.Kind)
	if err != nil {
		return err
	}

	if _, err := ref.Set(cxt, map[string]interface{}{
		"LastFetchedAt": time.Now(),
	}, firestore.Merge(firestore.FieldPath{"LastFetchedAt"})); err != nil {
		return errors.Wrap(err, "could not set data source last fetched")
	}

	return nil
}

type DataSources []DataSource
