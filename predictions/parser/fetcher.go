package parser

import (
	"bitbucket.org/heindl/taxa/store"
	"bufio"
	"cloud.google.com/go/storage"
	"context"
	"github.com/saleswise/errors/errors"
	"google.golang.org/api/iterator"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
)

const GCSPredictionsPath = "predictions"

type PredictionResult struct {
	Latitude, Longitude float64
	Date                string
	Target, Random      float64
	Taxon               store.TaxonID
}

type GCSFetcher interface {
	FetchLatestPredictionFileNames(cxt context.Context, id store.TaxonID, date string) ([]string, error)
	FetchPredictions(cxt context.Context, gcsFilePath string) ([]PredictionResult, error)
}

func NewGCSFetcher(cxt context.Context, bucketName string, localPath string) (GCSFetcher, error) {
	if localPath != "" {
		if strings.Contains(localPath, bucketName) {
			return &gcsFetcher{LocalPath: path.Join(localPath, "predictions")}, nil
		} else {
			return &gcsFetcher{LocalPath: path.Join(localPath, bucketName, "predictions")}, nil
		}
	}
	client, err := storage.NewClient(cxt)
	if err != nil {
		return nil, errors.Wrap(err, "Could not create Google Cloud Storage client.")
	}
	return &gcsFetcher{Bucket: client.Bucket(bucketName)}, nil
}

type gcsFetcher struct {
	LocalPath string
	Bucket    *storage.BucketHandle
}

func (Ω *gcsFetcher) FetchLatestPredictionFileNames(cxt context.Context, id store.TaxonID, date string) ([]string, error) {

	if !id.Valid() {
		return nil, errors.New("Invalid TaxonID")
	}

	if len(date) != 8 && date != "*" {
		return nil, errors.New("Date must be in format YYYYMMDD")
	}

	if Ω.LocalPath != "" {
		return Ω.fetchLocalPredictionFileNames(cxt, id, date)
	} else {
		return Ω.fetchRemoteFileNames(cxt, id, date)
	}

}

func (Ω *gcsFetcher) fetchLocalPredictionFileNames(cxt context.Context, id store.TaxonID, date string) ([]string, error) {

	dates := []string{}
	if date == "*" {
		files, err := ioutil.ReadDir(path.Join(Ω.LocalPath, string(id)))
		if err != nil {
			return nil, errors.Wrapf(err, "could not read taxon prediction file [%s]", path.Join(Ω.LocalPath, string(id)))
		}
		for _, file := range files {
			if file.IsDir() {
				dates = append(dates, file.Name())
			}
		}
	} else {
		dates = append(dates, date)
	}

	res := []string{}
	for _, date := range dates {
		p := path.Join(Ω.LocalPath, string(id), date)
		files, err := ioutil.ReadDir(p)
		if err != nil {
			return nil, errors.Wrapf(err, "could not read taxon prediction file [%s]", path.Join(Ω.LocalPath, string(id)))
		}
		sort.Sort(FileNames(files))
		res = append(res, path.Join(p, files[len(files)-1].Name()))
	}

	return res, nil
}

type FileNames []os.FileInfo

func (s FileNames) Len() int {
	return len(s)
}
func (s FileNames) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s FileNames) Less(i, j int) bool {
	return s[i].Name() < s[i].Name()
}

func (Ω *gcsFetcher) fetchRemoteFileNames(cxt context.Context, id store.TaxonID, date string) ([]string, error) {
	prefix := path.Join(GCSPredictionsPath, string(id))
	if date != "" {
		prefix = path.Join(prefix, date)
	}

	q := &storage.Query{
		Prefix: prefix,
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
		for _, n := range objectNames {
			if strings.Contains(n, path.Join(d, l[len(l)-1])) {
				response = append(response, n)
				break
			}
		}
	}

	return response, nil
}

func (Ω *gcsFetcher) FetchPredictions(cxt context.Context, gcsFilePath string) ([]PredictionResult, error) {
	scanner := &bufio.Scanner{}
	if strings.HasPrefix(gcsFilePath, "/tmp") {
		f, err := os.Open(gcsFilePath)
		if err != nil {
			return nil, errors.Wrapf(err, "could not open prediction file [%s]", gcsFilePath)
		}
		scanner = bufio.NewScanner(f)
		defer f.Close()
	} else {
		r, err := Ω.Bucket.Object(gcsFilePath).NewReader(cxt)
		if err != nil {
			return nil, errors.Wrapf(err, "could not get prediction object: %s", gcsFilePath)
		}
		defer r.Close()
		scanner = bufio.NewScanner(r)
	}

	taxon := taxonFromPredictionFilePath(gcsFilePath)
	if !taxon.Valid() {
		return nil, errors.Newf("invalid taxon [%s] from file path", taxon)
	}
	//buf := new(bytes.Buffer)
	//buf.ReadFrom(r)
	//fmt.Println(buf.String())
	scanner.Split(bufio.ScanLines)
	res := []PredictionResult{}
	for scanner.Scan() {
		var err error
		s := strings.Split(scanner.Text(), ",")
		r := PredictionResult{
			Date:  s[2],
			Taxon: taxon,
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
