package tfrecords

import (
	"bitbucket.org/heindl/process/store"
	"bufio"
	"cloud.google.com/go/storage"
	"context"
	"github.com/dropbox/godropbox/errors"
	"github.com/heindl/tfutils"
	"google.golang.org/api/iterator"
	"io"
	"strings"
)

type Iterator interface {
	Next(cxt context.Context) (Record, error)
}

func NewIterator(ctx context.Context, floraStore store.FloraStore, gcsFilePath string) (Iterator, error) {
	gcsObjects, err := floraStore.CloudStorageObjects(ctx, gcsFilePath, ".tfrecords")
	if err != nil {
		return nil, err
	}
	i := &tfIterator{
		gcsObjects: gcsObjects,
	}
	if err := i.callNextObject(ctx); err != nil {
		return nil, err
	}

	return i, nil
}

type tfIterator struct {
	gcsObjects    []*storage.ObjectHandle
	currentReader *tfutils.RecordReader
}

type FeatureType int

const (
	FeatureTypeBytes FeatureType = iota + 1
	FeatureTypeInt
	FeatureTypeFloat
)

func (Ω *tfIterator) Next(ctx context.Context) (Record, error) {
	for {
		if Ω.currentReader == nil {
			return nil, iterator.Done
		}
		r, err := Ω.currentReader.ReadRecord()

		if err == io.EOF {
			Ω.currentReader = nil
			if err := Ω.callNextObject(ctx); err != nil {
				return nil, err
			}
			continue
		}
		if err != nil {
			return nil, errors.Wrap(err, "Could not read TfRecord")
		}
		return NewRecord(r)
	}

}

func (Ω *tfIterator) callNextObject(ctx context.Context) error {
	// This slows things down, but to be sure of data consistency ...
	if len(Ω.gcsObjects) == 0 {
		return nil
	}

	attrs, err := Ω.gcsObjects[0].Attrs(ctx)
	if err != nil {
		return err
	}

	if !strings.HasSuffix(attrs.Name, ".tfrecords") {
		return errors.Newf("Expected a tfrecords encoded file [%s]", attrs.Name)
	}
	//
	//if attrs.ContentEncoding != "gzip" {
	//	if _, err := gcsObject.Update(ctx, storage.ObjectAttrsToUpdate{ContentEncoding: "gzip"}); err != nil {
	//		return 0, errors.Wrapf(err, "Could not update ContentEncoding of GCS Object [%s]", attrs.Name)
	//	}
	//}

	gcsObjectReader, err := Ω.gcsObjects[0].ReadCompressed(false).NewReader(ctx)
	if err != nil {
		return err
	}

	reader, err := tfutils.NewReaderFromBufio(bufio.NewReader(gcsObjectReader), &tfutils.RecordReaderOptions{
		CompressionType: tfutils.CompressionTypeNone,
	})
	if err != nil {
		return errors.Wrap(err, "Could not get TfRecordReader")
	}

	Ω.currentReader = reader

	Ω.gcsObjects = Ω.gcsObjects[1:]
	return nil
}
