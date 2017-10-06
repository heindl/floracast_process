package store

import (
	"github.com/saleswise/errors/errors"
	"context"
	"cloud.google.com/go/firestore"
	"time"
	"fmt"
	"github.com/fatih/structs"
)



type DataSourceID string
const (
	DataSourceIDGBIF        = DataSourceID("27")
	DataSourceIDINaturalist = DataSourceID("F1")
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
	DataSourceID("11"):           "NatureServe Explorer: An online encyclopedia of life. Version 7.1",
	DataSourceID("12"):           "Calflora",
	DataSourceID("13"):           "Odonata Central",
	DataSourceID("14"):           "IUCN Red List of Threatened Species. Version 2012.2",
	DataSourceID("10"):           "eBird/Clements Checklist 6.7",
	DataSourceID("15"):           "CONABIO",
	DataSourceID("6"):            "The Reptile Database",
	DataSourceID("16"):           "Afribats",
	DataSourceID("18"):           "Norma 059, 2010",
	DataSourceID("4"):            "Draft IUCN/SSC, 2013.1",
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
}

type DataSourceKind string
const (
	DataSourceKindOccurrence        = DataSourceKind("occurrence")
	DataSourceKindPhoto = DataSourceKind("photo")
	DataSourceKindDescription = DataSourceKind("description")
)

func (Ω DataSourceKind) Valid() bool {
	return Ω == DataSourceKindOccurrence || Ω == DataSourceKindPhoto || Ω == DataSourceKindDescription
}
//
//const EntityKindMetaScheme = "MetaScheme"
//func NewMetaScheme(origin DataSourceDocRef, id DataSourceTargetID, txn *datastore.Key) *DataSource {
//	return newScheme(EntityKindMetaScheme, origin, id, txn)
//}
//const EntityKindOccurrenceScheme = "OccurrenceScheme"
//func NewOccurrenceScheme(origin DataSourceDocRef, id DataSourceTargetID, parent *datastore.Key) *DataSource {
//	return newScheme(EntityKindOccurrenceScheme, origin, id, parent)
//}

// TargetID is the query object for the data source. So if the INaturalist taxon id is 12345, and the GBIF taxon id (targetid) is 678910.
type DataSourceTargetID string

type DataSource struct{
	Kind DataSourceKind `firestore:",omitempty"`
	SourceID      DataSourceID `firestore:",omitempty"`

	TargetID      DataSourceTargetID `firestore:",omitempty"`
	TaxonID TaxonID `firestore:",omitempty"`
	CreatedAt     *time.Time          `firestore:",omitempty"`
	ModifiedAt    *time.Time          `firestore:",omitempty"`
	LastFetchedAt *time.Time          `firestore:",omitempty"`
}

//type SchemeKeyName string
//const schemeKeyNameSeperator  = "|||"

//func (Ω SchemeKeyName) Parse() (DataSourceDocRef, DataSourceTargetID) {
//	s := strings.Split(string(Ω), schemeKeyNameSeperator)
//	return DataSourceDocRef(s[0]), DataSourceTargetID(s[1])
//}
//
//func NewSchemeKeyName(id DataSourceDocRef, targetID DataSourceTargetID) SchemeKeyName {
//	return SchemeKeyName(fmt.Sprintf("%s%s%s", id, schemeKeyNameSeperator, targetID))
//}
//
//func (Ω SchemeKeyName) Validate() error {
//	if !strings.Contains(string(Ω), schemeKeyNameSeperator) {
//		return errors.New("invalid scheme key name").SetState(M{utils.LogkeyScheme: Ω})
//	}
//	return nil
//}

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

	//if err := SchemeKeyName(Ω.Key.Name).Validate(); err != nil {
	//	return err
	//}
	//
	//if Ω.Key.Kind != EntityKindMetaScheme && Ω.Key.Kind != EntityKindOccurrenceScheme {
	//	return errors.New("wrong key kind for scheme")
	//}
	//
	//if !ValidTaxonKey(Ω.Key.Parent) {
	//	return errors.New("invalid scheme parent taxon")
	//}

	return nil
}

//func (Ω DataSource) Combine(s *DataSource) *DataSource {
//	if s.CreatedAt.Before(Ω.CreatedAt) {
//		Ω.CreatedAt = s.CreatedAt
//	}
//	if Ω.LastFetchedAt.Before(s.LastFetchedAt) {
//		Ω.LastFetchedAt = s.LastFetchedAt
//	}
//	return &Ω
//}

//func newScheme(entityKind string, origin DataSourceDocRef, target DataSourceTargetID, txn *datastore.Key) *DataSource {
//	return &DataSource{
//			Key: datastore.NameKey(entityKind, string(NewSchemeKeyName(origin, target)), txn),
//			CreatedAt: time.Now(),
//			ModifiedAt: time.Now(),
//	}
//}

//func (Ω *store) getScheme(kind string, origin DataSourceDocRef, id DataSourceTargetID, txn *datastore.Key) (*DataSource, error) {
//
//	q := datastore.NewQuery(kind).
//		Filter("SourceID =", string(origin)).
//		Filter("TargetID =", string(id)).
//		Ancestor(txn)
//
//	res := []*DataSource{}
//	keys, err := Ω.FirebaseClient.GetAll(context.Background(), q, &res)
//	if err != nil {
//		return nil, errors.Wrap(err, "could not get datastore scheme")
//	}
//
//	if len(keys) > 1 {
//		return nil, errors.New("unexpectedly found more than one meta scheme").SetState(M{
//			utils.LogkeyIdentifier: id,
//			utils.LogkeyStringValue: origin,
//			utils.LogkeyDatastoreKey: txn,
//		})
//	}
//
//	if len(keys) == 0 {
//		return nil, nil
//	}
//
//	return res[0], nil
//
//}


func (Ω *store) GetOccurrenceDataSources(context context.Context, taxonID TaxonID) (res DataSources, err error) {

	q := Ω.FirestoreClient.Collection("DataSources").
		Where("Kind", "==", DataSourceKindOccurrence);

	if taxonID.Valid() {
		q = q.Where("TaxonID", "==", taxonID)
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

func (Ω *store) NewDataSourceDocumentRef(taxonID TaxonID, dataSourceID DataSourceID, targetID DataSourceTargetID, kind DataSourceKind) (*firestore.DocumentRef, error) {

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
		Doc(fmt.Sprintf("%s|%s|%s|%s", taxonID, dataSourceID, targetID, kind)), nil

}

func (Ω *store) UpsertDataSource(cxt context.Context, src DataSource) error {

	ref, err := Ω.NewDataSourceDocumentRef(src.TaxonID, src.SourceID, src.TargetID, src.Kind)
	if err != nil {
		return err
	}

	if err := Ω.FirestoreClient.RunTransaction(cxt, func(cxt context.Context, tx *firestore.Transaction) error {
		return tx.UpdateMap(ref, structs.Map(src))
	}); err != nil {
		return errors.Wrap(err, "could not update data source")
	}
}

//func (Ω *store) CreateScheme(s *DataSource) (*datastore.Key, error) {
//
//	if !ValidTaxonKey(s.Key.Parent) {
//		return nil, errors.New("invalid scheme parent taxon key").SetState(M{utils.LogkeyDatastoreKey: s.Key.Parent})
//	}
//	//if s.SourceID == "" {
//	//	return nil, errors.New("invalid scheme id origin")
//	//}
//	//if s.TargetID == "" {
//	//	return nil, errors.New("invalid scheme id")
//	//}
//
//	// No opt if we find the scheme. Can have an additional update function below if
//	// we need to update additional fields.
//	if sch, err := Ω.getScheme(s.Key.Kind, s.SourceID, s.TargetID, s.Key.Parent); err != nil {
//		return nil, err
//	} else if sch != nil {
//		return sch.Key, nil
//	}
//
//	k, err := Ω.FirebaseClient.Put(context.Background(), s.Key, s)
//	if err != nil {
//		return nil, errors.Wrap(err, "could not set meta scheme")
//	}
//
//	return k, nil
//}

func (Ω *store) UpdateSchemaLastFetched(cxt context.Context, src DataSource) error {

	ref, err := Ω.NewDataSourceDocumentRef(src.TaxonID, src.SourceID, src.TargetID, src.Kind)
	if err != nil {
		return err
	}

	if err := Ω.FirestoreClient.RunTransaction(cxt, func(cxt context.Context, tx *firestore.Transaction) error {
		return tx.UpdateMap(ref, map[string]interface{}{
			"LastFetchedAt": time.Now(),
		})
	}); err != nil {
		return errors.Wrap(err, "could not update data source last fetched")
	}

	return nil
}

type DataSources []DataSource

//func (Ω DataSources) RemoveDuplicates() (response DataSources) {
//	for _, s := range Ω {
//		if response.Find(s) == nil {
//			response = append(response, s)
//		}
//	}
//	return
//}

//func (Ω DataSources) Find(needle DataSource) *DataSource {
//	for _, s := range Ω {
//		if s.TaxonID == needle.TaxonID && s.Kind == needle.Kind {
//			return &s
//		}
//	}
//	return nil
//}

//func (Ω DataSources) AddToSet(s *DataSource) (DataSources, error) {
//
//	if err := s.Validate(); err != nil {
//		return nil, err
//	}
//
//	for i := range Ω {
//		if Ω[i].Key.Kind != s.Key.Kind {
//			continue
//		}
//		if Ω[i].Key.Name != s.Key.Name {
//			continue
//		}
//		if Ω[i].Key.Parent.ID != s.Key.Parent.ID {
//			continue
//		}
//		Ω[i] = Ω[i].Combine(s)
//		return Ω, nil
//	}
//
//	return append(Ω, s), nil
//
//}
