package store
//
//import (
//	"cloud.google.com/go/firestore"
//	"context"
//	"fmt"
//	"github.com/saleswise/errors/errors"
//	"strings"
//	"time"
//	"google.golang.org/api/iterator"
//	"bitbucket.org/heindl/taxa/datasources"
//)
//
//type DataSourceKind string
//
//const (
//	DataSourceKindOccurrence  = DataSourceKind("occurrence")
//	DataSourceKindPhoto       = DataSourceKind("photo")
//	DataSourceKindDescription = DataSourceKind("description")
//)
//
//func (Ω DataSourceKind) Valid() bool {
//	return Ω == DataSourceKindOccurrence || Ω == DataSourceKindPhoto || Ω == DataSourceKindDescription
//}
//
//type OccurrenceSource struct {
//	SourceID      datasources.DataSourceType     `firestore:",omitempty"`
//	TargetID      datasources.DataSourceTargetID `firestore:",omitempty"`
//	CanonicalName string             `firestore:",omitempty"`
//	CreatedAt     *time.Time         `firestore:",omitempty"`
//	ModifiedAt    *time.Time         `firestore:",omitempty"`
//	LastFetchedAt *time.Time         `firestore:",omitempty"`
//}
//
//type DataSource struct {
//	Kind     DataSourceKind `firestore:",omitempty"`
//	SourceID datasources.DataSourceType `firestore:",omitempty"`
//
//	TargetID      datasources.DataSourceTargetID `firestore:",omitempty"`
//	TaxonID       INaturalistTaxonID `firestore:",omitempty"`
//	CreatedAt     *time.Time         `firestore:",omitempty"`
//	ModifiedAt    *time.Time         `firestore:",omitempty"`
//	LastFetchedAt *time.Time         `firestore:",omitempty"`
//}
//
//var DataSourceFieldsToMerge = []firestore.FieldPath{
//	firestore.FieldPath{"Kind"},
//	firestore.FieldPath{"SourceID"},
//	firestore.FieldPath{"TargetID"},
//	firestore.FieldPath{"INaturalistTaxonID"},
//	firestore.FieldPath{"ModifiedAt"},
//	firestore.FieldPath{"CreatedAt"},
//}
//
//func (Ω *DataSource) Validate() error {
//
//	if !Ω.TaxonID.Valid() {
//		return errors.New("invalid taxon id")
//	}
//
//	if Ω.TargetID == "" {
//		return errors.New("invalid target id")
//	}
//
//	if !Ω.SourceID.Valid() {
//		return errors.Newf("invalid source id %s", Ω.SourceID)
//	}
//
//	return nil
//}
//func (Ω *store) GetSourceLastCreated(cxt context.Context, kind DataSourceKind, srcID datasources.DataSourceType) (*time.Time, error) {
//	iter := Ω.FirestoreClient.Collection(CollectionTypeDataSources).
//		Where("Kind", "==", kind).
//			Where("SourceID", "==", srcID).
//				OrderBy("CreatedAt", firestore.Desc).
//					Limit(1).
//						Documents(cxt)
//
//	for {
//		ref, err := iter.Next()
//		if err == iterator.Done {
//			break
//		}
//		if err != nil {
//			return nil, errors.Wrap(err, "could no fetch latest source")
//		}
//
//		src := DataSource{}
//		if err := ref.DataTo(&src); err != nil {
//			return nil, errors.Wrap(err, "could not cast source")
//		}
//
//		return src.ModifiedAt, nil
//	}
//
//	return nil, nil
//}
//
//
//func (Ω *store) GetOccurrenceDataSources(context context.Context, taxonID INaturalistTaxonID) (res DataSources, err error) {
//
//	q := Ω.FirestoreClient.Collection(CollectionTypeDataSources).
//		Where("Kind", "==", DataSourceKindOccurrence)
//
//	if taxonID.Valid() {
//		q = q.Where("INaturalistTaxonID", "==", taxonID)
//	}
//
//	snaps, err := q.Documents(context).GetAll()
//	if err != nil {
//		return nil, errors.Wrapf(err, "could not fetch occurrence data sources for taxon id[%s]", taxonID)
//	}
//
//	for _, s := range snaps {
//		var src DataSource
//		if err := s.DataTo(&src); err != nil {
//			return nil, errors.Wrap(err, "could not typecast data source")
//		}
//		res = append(res, src)
//	}
//
//	return
//}
//
//func (Ω *store) NewDataSourceDocumentRef(taxonID INaturalistTaxonID, dataSourceID datasources.DataSourceType, targetID datasources.DataSourceTargetID, kind DataSourceKind) (*firestore.DocumentRef, error) {
//
//	if !taxonID.Valid() {
//		return nil, errors.New("invalid data source document reference id")
//	}
//	if !dataSourceID.Valid() {
//		return nil, errors.New("invalid data source id")
//	}
//	if !kind.Valid() {
//		return nil, errors.New("invalid kind")
//	}
//	if targetID == "" {
//		return nil, errors.New("invalid target id")
//	}
//
//	return Ω.FirestoreClient.Collection(CollectionTypeDataSources).
//		Doc(fmt.Sprintf("%s|%s|%s|%s", string(taxonID), dataSourceID, targetID, kind)), nil
//
//}
//
//func (Ω *store) UpsertDataSource(cxt context.Context, src DataSource) error {
//
//	ref, err := Ω.NewDataSourceDocumentRef(src.TaxonID, src.SourceID, src.TargetID, src.Kind)
//	if err != nil {
//		return err
//	}
//
//	if err := Ω.FirestoreClient.RunTransaction(cxt, func(cxt context.Context, tx *firestore.Transaction) error {
//		if _, err := tx.Get(ref); err != nil {
//			if strings.Contains(err.Error(), "not found") {
//				return tx.Set(ref, src)
//			} else {
//				return err
//			}
//		}
//		return tx.Set(ref, src, firestore.Merge(DataSourceFieldsToMerge...))
//	}); err != nil {
//		return errors.Wrap(err, "could not update data source")
//	}
//	return nil
//}
//
//func (Ω *store) UpdateDataSourceLastFetched(cxt context.Context, src DataSource) error {
//
//	ref, err := Ω.NewDataSourceDocumentRef(src.TaxonID, src.SourceID, src.TargetID, src.Kind)
//	if err != nil {
//		return err
//	}
//
//	if _, err := ref.Set(cxt, map[string]interface{}{
//		"LastFetchedAt": time.Now(),
//	}, firestore.Merge(firestore.FieldPath{"LastFetchedAt"})); err != nil {
//		return errors.Wrap(err, "could not set data source last fetched")
//	}
//
//	return nil
//}
//
//type DataSources []DataSource
