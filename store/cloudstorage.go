package store

import (
	"github.com/heindl/floracast_process/utils"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"google.golang.org/api/iterator"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

func (Ω *store) SyncGCSPathWithLocal(ctx context.Context, cloudPath string, localPath string) error {

	files, err := ioutil.ReadDir(localPath)
	if err != nil {
		return errors.Wrapf(err, "Could not read files in local directory [%s]", localPath)
	}
	if len(files) > 0 {
		return errors.Newf("Expected local directory [%s] to be empty, instead have %d files", len(files))
	}

	objects, err := Ω.CloudStorageObjects(ctx, cloudPath)
	if err != nil {
		return err
	}

	for _, obj := range objects {

		attrs, err := obj.Attrs(ctx)
		if err != nil {
			return errors.Wrap(err, "Could not get object attributes")
		}

		isDir := strings.HasSuffix(attrs.Name, "/")

		fName := path.Join(localPath, strings.TrimPrefix(attrs.Name, cloudPath))

		fmt.Println("filename", fName)

		if isDir {
			fmt.Println("isdir")
			if err := os.MkdirAll(fName, os.ModePerm); err != nil {
				return errors.Wrapf(err, "Could not make directories [%s]", fName)
			}
			continue
		} else {
			dir := filepath.Dir(fName)
			fmt.Println("base", dir)
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				return errors.Wrapf(err, "Could not make base directory [%s]", dir)
			}
		}
		reader, err := obj.NewReader(ctx)
		if err != nil {
			return errors.Wrapf(err, "Could not read GCS objects [%s]", attrs.Name)
		}

		b, err := ioutil.ReadAll(reader)
		if err != nil {
			return errors.Wrapf(err, "Could not read GCS object [%s]", attrs.Name)
		}
		if err := ioutil.WriteFile(fName, b, os.ModePerm); err != nil {
			return errors.Wrapf(err, "Could not write file [%s]", fName)
		}
	}

	return nil

}

// ".tfrecord.gz"
func (Ω *store) CloudStorageObjects(ctx context.Context, pathname string, requiredSuffixes ...string) ([]*storage.ObjectHandle, error) {

	pathname = strings.TrimSpace(pathname)
	if strings.HasPrefix(pathname, "gs://") {
		pathname = strings.TrimPrefix(pathname, "gs://")
		l := strings.Split(pathname, "/")
		pathname = strings.Join(l[1:], "/")
	}

	if len(requiredSuffixes) > 0 && utils.HasSuffix(pathname, requiredSuffixes...) {
		return []*storage.ObjectHandle{Ω.gcsBucketHandle.Object(pathname)}, nil
	}

	iter := Ω.gcsBucketHandle.Objects(ctx, &storage.Query{
		Prefix: pathname,
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
		if len(requiredSuffixes) == 0 || utils.HasSuffix(o.Name, requiredSuffixes...) {
			gcsObjects = append(gcsObjects, Ω.gcsBucketHandle.Object(o.Name))
		}
	}

	return gcsObjects, nil
}

func (Ω *store) CloudStorageObjectNames(ctx context.Context, pathname string, requiredSuffixes ...string) ([]string, error) {

	objs, err := Ω.CloudStorageObjects(ctx, pathname, requiredSuffixes...)
	if err != nil {
		return nil, err
	}
	names := []string{}
	for _, obj := range objs {
		attrs, err := obj.Attrs(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "Could not get object attributes")
		}
		name := strings.TrimPrefix(attrs.Name, pathname)
		name = strings.TrimPrefix(name, "/")
		names = append(names, name)
	}

	sort.Strings(names)

	return names, nil
}

//func (Ω *store) CloudStorageBucket() (*storage.BucketHandle, error) {
//	return Ω.gcsBucketHandle, nil
//}
