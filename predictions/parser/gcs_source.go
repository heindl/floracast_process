package parser

import (
	"bitbucket.org/heindl/process/nameusage/nameusage"
	"bitbucket.org/heindl/process/store"
	"bitbucket.org/heindl/process/utils"
	"cloud.google.com/go/storage"
	"context"
	"github.com/dropbox/godropbox/errors"
	"google.golang.org/api/iterator"
	"path"
	"sort"
	"strings"
)

func NewGCSPredictionSource(cxt context.Context, florastore store.FloraStore) (PredictionSource, error) {
	gcsHandle, err := florastore.CloudStorageBucket()
	if err != nil {
		return nil, err
	}
	return &gcsSource{gcsBucketHandle: gcsHandle}, nil
}

type gcsSource struct {
	gcsBucketHandle *storage.BucketHandle
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

	q := &storage.Query{
		Prefix: prefix,
	}

	iter := Ω.gcsBucketHandle.Objects(cxt, q)
	dateFiles := map[string][]string{}
	objectNames := []string{}
	for {
		o, err := iter.Next()
		if err != nil && err == iterator.Done {
			break
		}
		if err != nil {
			return nil, errors.Wrap(err, "error iterating over object names")
		}

		objectNames = append(objectNames, o.Name)

		s := strings.Split(o.Name, "/")
		// 0: predictions
		// 1: taxon
		// 2: date
		// 3. filename
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

func (Ω *gcsSource) FetchPredictions(cxt context.Context, gcsFilePath string) (res []*PredictionResult, err error) {

	r, err := Ω.gcsBucketHandle.Object(gcsFilePath).NewReader(cxt)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get prediction object: %s", gcsFilePath)
	}
	defer utils.SafeClose(r, &err)

	nameUsageID, err := parseNameUsageIDFromFilePath(gcsFilePath)
	if !nameUsageID.Valid() {
		return nil, errors.Newf("Invalid ID [%s] from FilePath", nameUsageID)
	}

	return parsePredictionReader(nameUsageID, r)
}
