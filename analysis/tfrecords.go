package main

import (
	"bitbucket.org/heindl/process/store"
	"bufio"
	"cloud.google.com/go/storage"
	"context"
	"github.com/dropbox/godropbox/errors"
	"github.com/heindl/tfutils"
	"io"
	"strings"
)

func countTFRecordsInCloudStoragePath(ctx context.Context, floraStore store.FloraStore, gcsFilePath string) (int, error) {

	gcsObjects, err := floraStore.CloudStorageObjects(ctx, gcsFilePath, ".tfrecord.gz")
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
