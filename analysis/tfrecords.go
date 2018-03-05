package main

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

func countTFRecordsInCloudStoragePath(ctx context.Context, floraStore store.FloraStore, gcsFilePath string) (int, error) {

	gcsObjects, err := listGCSObjects(ctx, floraStore, gcsFilePath)
	if err != nil {
		return 0, err
	}

	sumTotal := 0
	for _, gcsObject := range gcsObjects {
		sum, err := countRecordObject(ctx, gcsObject)
		if err != nil {
			return 0, err
		}
		sumTotal += sum
	}
	return sumTotal, nil
}

func countRecordObject(ctx context.Context, gcsObject *storage.ObjectHandle) (int, error) {

	// This slows things down, but to be sure of data consistency ...
	attrs, err := gcsObject.Attrs(ctx)
	if err != nil {
		return 0, err
	}

	if !strings.HasSuffix(attrs.Name, ".gz") {
		return 0, errors.Newf("Expected a gzip encoded file [%s]", attrs.Name)
	}

	if attrs.ContentEncoding != "gzip" {
		if _, err := gcsObject.Update(ctx, storage.ObjectAttrsToUpdate{ContentEncoding: "gzip"}); err != nil {
			return 0, errors.Wrapf(err, "Could not update ContentEncoding of GCS Object [%s]", attrs.Name)
		}
	}

	gcsObjectReader, err := gcsObject.ReadCompressed(false).NewReader(ctx)
	if err != nil {
		return 0, err
	}

	reader, err := tfutils.NewReader(bufio.NewReader(gcsObjectReader), &tfutils.RecordReaderOptions{
		CompressionType: tfutils.CompressionTypeNone,
	})
	if err != nil {
		return 0, errors.Wrap(err, "Could not get TfRecordReader")
	}
	sum := 0
	for {
		_, err := reader.ReadRecord()
		if err == io.EOF {
			return sum, nil
		}
		if err != nil {
			return 0, errors.Wrap(err, "Could not read TfRecord")
		}
		sum++
	}
	return sum, nil

}

func listGCSObjects(cxt context.Context, floraStore store.FloraStore, f string) ([]*storage.ObjectHandle, error) {

	gcsHandle, err := floraStore.CloudStorageBucket()
	if err != nil {
		return nil, err
	}

	f = strings.TrimSpace(f)
	if strings.HasPrefix(f, "gs://") {
		f = strings.TrimPrefix(f, "gs://")
		l := strings.Split(f, "/")
		f = strings.Join(l[1:], "/")
	}

	if strings.HasSuffix(f, ".tfrecord.gz") {
		return []*storage.ObjectHandle{gcsHandle.Object(f)}, nil
	}

	iter := gcsHandle.Objects(cxt, &storage.Query{
		Prefix: f,
	})

	gcsObjects := []*storage.ObjectHandle{}

	for {
		o, err := iter.Next()
		if err != nil && err == iterator.Done {
			break
		}
		if err != nil {
			return nil, errors.Wrap(err, "Could not iterate over object names")
		}
		if strings.HasSuffix(o.Name, ".tfrecord.gz") {
			gcsObjects = append(gcsObjects, gcsHandle.Object(o.Name))
		}
	}

	return gcsObjects, nil
}
