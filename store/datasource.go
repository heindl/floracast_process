package store

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/saleswise/errors/errors"
	"strings"
	"time"
	"google.golang.org/api/iterator"
	"strconv"
)

type DataSourceType string

const (
	DataSourceTypeGBIF             = DataSourceType("27")
	DataSourceTypeINaturalist      = DataSourceType("INAT")
	DataSourceTypeMushroomObserver = DataSourceType("MUOB")
	DataSourceTypeNatureServe      = DataSourceType("11")
)

func (Ω DataSourceType) Valid() bool {
	_, ok := SchemeSourceIDMap[Ω]
	return ok
}

var SchemeSourceIDMap = map[DataSourceType]string{
	// Floracast
	DataSourceTypeINaturalist: "iNaturalist",
	// INaturalist
	DataSourceType("1"):       "IUCN Red List of Threatened Species. Version 2012.1",
	DataSourceType("3"):       "Amphibian Species of the World 5.6",
	DataSourceType("2"):       "Amphibiaweb. 2012",
	DataSourceType("5"):       "Amphibian Species of the World 5.5",
	DataSourceType("17"):      "New England Wild Flower Society's Flora Novae Angliae",
	DataSourceTypeNatureServe: "NatureServe Explorer: An online encyclopedia of life. Version 7.1",
	DataSourceType("12"):      "Calflora",
	DataSourceType("13"):      "Odonata Central",
	DataSourceType("14"):      "IUCN Red List of Threatened Species. Version 2012.2",
	DataSourceType("10"):      "eBird/Clements Checklist 6.7",
	DataSourceType("15"):      "CONABIO",
	DataSourceType("6"):       "The Reptile Database",
	DataSourceType("16"):      "Afribats",
	DataSourceType("18"):      "Norma 059, 2010",
	DataSourceType("4"):          "Draft IUCN/SSC, 2013.1",
	DataSourceType("19"):         "Draft IUCN/SSC Amphibian Specialist Group, 2011",
	DataSourceType("20"):         "eBird/Clements Checklist 6.8",
	DataSourceType("21"):         "IUCN Red List of Threatened Species. Version 2013.2",
	DataSourceType("22"):         "eBird/Clements Checklist 6.9",
	DataSourceType("23"):         "NatureWatch NZ",
	DataSourceType("24"):           "The world spider catalog, version 15.5",
	DataSourceType("25"):           "Carabidae of the World",
	DataSourceType("26"):           "IUCN Red List of Threatened Species. Version 2014.3",
	DataSourceTypeGBIF:             "GBIF",
	DataSourceType("28"):           "NPSpecies",
	DataSourceType("29"):           "Esslinger&#39;s North American Lichens",
	DataSourceType("30"):           "Amphibian Species of the World 6.0",
	DataSourceType("31"):           "Esslinger&#39;s North American Lichens, Version 21",
	DataSourceTypeMushroomObserver: "MushroomObserver.org",
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

func (Ω DataSourceTargetID) Valid() bool {
	return string(Ω) != ""
}

func (Ω DataSourceTargetID) ToInt() (int, error) {
	i, err := strconv.Atoi(string(Ω))
	if err != nil {
		return 0, errors.Wrap(err, "Could not cast TargetID as int")
	}
	return i, nil
}

func NewDataSourceTargetIDFromInt(i int) (*DataSourceTargetID, error) {
	if i == 0 {
		return nil, errors.New("Invalid DataSourceTargetID: Received zero.")
	}
	id := DataSourceTargetID(strconv.Itoa(i))
	return &id, nil
}

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
	SourceID      DataSourceType     `firestore:",omitempty"`
	TargetID      DataSourceTargetID `firestore:",omitempty"`
	CanonicalName string             `firestore:",omitempty"`
	CreatedAt     *time.Time         `firestore:",omitempty"`
	ModifiedAt    *time.Time         `firestore:",omitempty"`
	LastFetchedAt *time.Time         `firestore:",omitempty"`
}

type DataSource struct {
	Kind     DataSourceKind `firestore:",omitempty"`
	SourceID DataSourceType `firestore:",omitempty"`

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
func (Ω *store) GetSourceLastCreated(cxt context.Context, kind DataSourceKind, srcID DataSourceType) (*time.Time, error) {
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

func (Ω *store) NewDataSourceDocumentRef(taxonID INaturalistTaxonID, dataSourceID DataSourceType, targetID DataSourceTargetID, kind DataSourceKind) (*firestore.DocumentRef, error) {

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
