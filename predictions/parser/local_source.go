package parser

import (
	"path"
	"context"
	"github.com/dropbox/godropbox/errors"
	"bitbucket.org/heindl/process/nameusage/nameusage"
	"io/ioutil"
	"sort"
	"os"
)

func NewLocalPredictionSource(cxt context.Context, localPath string) (PredictionSource, error) {
	if localPath == "" {
		return nil, errors.New("Invalid Local Path")
	}
	return &localSource{
		path: path.Join(localPath, "predictions"),
	}, nil
}

type localSource struct{
	path string
}

func (Ω *localSource) FetchLatestPredictionFileNames(cxt context.Context, id nameusage.NameUsageID, date string) ([]string, error) {

	if !id.Valid() {
		return nil, errors.New("Invalid NameUsageID")
	}

	if len(date) != 8 && date != "*" {
		return nil, errors.New("Date must be in format YYYYMMDD")
	}

	dates := []string{}
	if date == "*" {
		files, err := ioutil.ReadDir(path.Join(Ω.path, string(id)))
		if err != nil {
			return nil, errors.Wrapf(err, "could not read taxon prediction file [%s]", path.Join(Ω.path, string(id)))
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
		p := path.Join(Ω.path, string(id), date)
		files, err := ioutil.ReadDir(p)
		if err != nil {
			return nil, errors.Wrapf(err, "could not read taxon prediction file [%s]", path.Join(Ω.path, string(id)))
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

func (Ω *localSource) FetchPredictions(cxt context.Context, filePath string) (res []*PredictionResult, err error) {

	f, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "could not open prediction file [%s]", filePath)
	}
	defer func(){
		if closeErr := f.Close(); closeErr != nil && err != nil {
			err = closeErr
			res = nil
			return
		}
	}()

	nameUsageID, err := parseNameUsageIDFromFilePath(filePath)
	if err != nil {
		return nil, err
	}

	if !nameUsageID.Valid() {
		return nil, errors.Newf("invalid taxon [%s] from file path", nameUsageID)
	}

	return parsePredictionReader(nameUsageID, f)
}

