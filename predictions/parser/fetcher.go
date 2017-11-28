package parser

import (
	"bitbucket.org/heindl/taxa/store"
	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"context"
	"github.com/saleswise/errors/errors"
	"strings"
	"sort"
	"bufio"
	"strconv"
	"path"
)

const GCSPredictionsPath = "predictions"

type PredictionResult struct{
	Latitude, Longitude float64
	Date string
	Target, Random float64
}

type GCSFetcher interface{
	FetchLatestPredictionFileNames(cxt context.Context, id store.TaxonID, date string) ([]string, error)
	FetchPredictions(cxt context.Context, gcsFilePath string) ([]PredictionResult, error)
}

func NewGCSFetcher(cxt context.Context, bucketName string) (GCSFetcher, error) {
	client, err := storage.NewClient(cxt)
	if err != nil {
		return nil ,errors.Wrap(err,"Could not create Google Cloud Storage client.")
	}
	return &gcsFetcher{Bucket: client.Bucket(bucketName)}, nil
}

type gcsFetcher struct{
	Bucket *storage.BucketHandle
}

func (Ω *gcsFetcher) FetchLatestPredictionFileNames(cxt context.Context, id store.TaxonID, date string) ([]string, error) {

	if !id.Valid() {
		return nil, errors.New("Invalid TaxonID")
	}

	if date != "" && len(date) != 8 {
		return nil, errors.New("Date must be in format YYYYMMDD")
	}

	prefix := path.Join(GCSPredictionsPath, string(id))
	if date != "" {
		prefix = path.Join(prefix, date)
	}

	q := &storage.Query{
		Prefix:    prefix,
	}

	iter := Ω.Bucket.Objects(cxt, q)
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
		for _, n := range objectNames{
			if strings.Contains(n, path.Join(d, l[len(l)-1])) {
				response = append(response, n)
				break
			}
		}
	}

	return response, nil

}

func (Ω *gcsFetcher) FetchPredictions(cxt context.Context, gcsFilePath string) ([]PredictionResult, error) {

	r, err := Ω.Bucket.Object(gcsFilePath).NewReader(cxt)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get prediction object: %s", gcsFilePath)
	}
	defer r.Close()
	//buf := new(bytes.Buffer)
	//buf.ReadFrom(r)
	//fmt.Println(buf.String())
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	res := []PredictionResult{}
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ",")
		r := PredictionResult{
			Date: s[2],
		}
		r.Latitude, err = strconv.ParseFloat(s[0], 64)
		if err != nil {
			return nil, errors.Wrap(err, "could not parse latitude")
		}
		r.Longitude, err = strconv.ParseFloat(s[1], 64)
		if err != nil {
			return nil, errors.Wrap(err, "could not parse longitude")
		}
		r.Target, err = strconv.ParseFloat(s[3], 64)
		if err != nil {
			return nil, errors.Wrap(err, "could not parse target")
		}
		r.Random, err = strconv.ParseFloat(s[4], 64)
		if err != nil {
			return nil, errors.Wrap(err, "could not parse random")
		}
		res = append(res, r)
	}
	return res, nil
}
