package parser

import (
	"github.com/heindl/floracast_process/nameusage/nameusage"
	"github.com/heindl/floracast_process/store"
	"github.com/heindl/floracast_process/utils"
	"cloud.google.com/go/storage"
	"context"
	"github.com/dropbox/godropbox/errors"
	"gopkg.in/tomb.v2"
	"path"
	"sort"
	"strings"
	"sync"
)

// NewGCSPredictionSource returns a source for fetching or loading predictions.
func NewGCSPredictionSource(cxt context.Context, floraStore store.FloraStore) (PredictionSource, error) {

	return &gcsSource{floraStore: floraStore}, nil
}

type gcsSource struct {
	floraStore store.FloraStore
	//gcsBucketHandle *storage.BucketHandle
}

func (Ω *gcsSource) FetchLatestPredictionFileNames(cxt context.Context, id nameusage.ID, date string) ([]string, error) {

	if !id.Valid() {
		return nil, errors.Newf("Invalid ID [%s]", id)
	}

	if len(date) != 8 && date != "*" {
		return nil, errors.New("Date must be in format YYYYMMDD")
	}

	prefix := path.Join(GCSPredictionsPath, string(id))
	if date != "" {
		prefix = path.Join(prefix, date)
	}

	gcsObjects, err := Ω.floraStore.CloudStorageObjects(cxt, prefix, ".csv")
	if err != nil {
		return nil, err
	}

	return sortedFileNames(cxt, gcsObjects)
}

func sortedFileNames(cxt context.Context, gcsObjects []*storage.ObjectHandle) ([]string, error) {
	dateFiles := map[string][]string{}
	objectNames := []string{}
	for _, gcsObject := range gcsObjects {
		attrs, err := gcsObject.Attrs(cxt)
		if err != nil {
			return nil, errors.Wrap(err, "Could not get GCS Object Attributes")
		}
		s := strings.Split(attrs.Name, "/")
		// [0] predictions, [1] taxon, [2] date, [3] filename
		if _, ok := dateFiles[s[2]]; !ok {
			dateFiles[s[2]] = []string{}
		}
		dateFiles[s[2]] = append(dateFiles[s[2]], s[3])
	}
	response := []string{}
	for d, l := range dateFiles {
		sort.Strings(l)
		for _, n := range objectNames {
			if strings.Contains(n, path.Join(d, l[len(l)-1])) {
				response = append(response, n)
				break
			}
		}
	}

	return response, nil
}

func (Ω *gcsSource) FetchPredictions(cxt context.Context, gcsFilePath string) ([]*PredictionResult, error) {

	gcsObjects, err := Ω.floraStore.CloudStorageObjects(cxt, gcsFilePath, "")
	if err != nil {
		return nil, err
	}

	res := []*PredictionResult{}
	lock := sync.Mutex{}

	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _, _gcsObject := range gcsObjects {
			gcsObject := _gcsObject
			tmb.Go(func() (err error) {
				r, err := gcsObject.NewReader(cxt)
				if err != nil {
					return err
				}
				defer utils.SafeClose(r, &err)
				nameUsageID, err := parseNameUsageIDFromFilePath(gcsFilePath)
				if !nameUsageID.Valid() {
					return errors.Newf("Invalid ID [%s] from FilePath", nameUsageID)
				}
				newPredictions, err := parsePredictionReader(nameUsageID, r)
				if err != nil {
					return err
				}
				lock.Lock()
				defer lock.Unlock()
				res = append(res, newPredictions...)
				return nil
			})
		}
		return nil
	})
	if err := tmb.Wait(); err != nil {
		return nil, err
	}
	return res, nil
}
