package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sync"

	"bitbucket.org/heindl/process/nameusage/nameusage"
	"bitbucket.org/heindl/process/predictions"
	"github.com/dropbox/godropbox/errors"
	"github.com/elgs/gostrgen"
	"gopkg.in/tomb.v2"
)

type localFileCache struct {
	filePath string
	sync.Mutex
	filePointers map[string]*os.File
}

func NewLocalFileCache() (PredictionCache, func() error, error) {
	random_string, err := gostrgen.RandGen(10, gostrgen.Lower, "", "")
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not create temp file random string name")
	}
	tmp := path.Join("/tmp/", fmt.Sprintf("predictions-%s", random_string))
	if err := os.Mkdir(tmp, os.ModePerm); err != nil {
		return nil, nil, errors.Wrap(err, "could not create tmp path")
	}
	fmt.Println("TEMP_FILE_DIRECTORY", tmp)

	c := localFileCache{
		filePath:     tmp,
		filePointers: map[string]*os.File{},
	}

	return &c, c.Close, nil
}

func (Ω *localFileCache) ReadPredictions(lat, lng, radius float64, qDate string, usageID *nameusage.NameUsageID) ([]string, error) {
	return nil, nil
}

func (Ω *localFileCache) getFilePointer(p predictions.Prediction) (*os.File, error) {

	formattedDate, err := p.Date()
	if err != nil {
		return nil, err
	}

	usageID, err := p.UsageID()
	if err != nil {
		return nil, err
	}

	filename := fmt.Sprintf("%s-%s.jsonl", formattedDate, usageID)

	Ω.Lock()
	defer Ω.Unlock()

	if _, ok := Ω.filePointers[filename]; ok {
		return Ω.filePointers[filename], nil
	}

	taxonFile, err := os.OpenFile(path.Join(Ω.filePath, filename), os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not open TaxonFile [%s]", filename)
	}

	Ω.filePointers[filename] = taxonFile

	return taxonFile, nil
}

func (Ω *localFileCache) WritePredictions(prediction_list predictions.Predictions) error {

	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _, _prediction := range prediction_list {
			prediction := _prediction
			tmb.Go(func() error {
				b, err := json.Marshal(prediction)
				if err != nil {
					return errors.Wrap(err, "could not marshal prediction")
				}
				f, err := Ω.getFilePointer(prediction)
				if err != nil {
					return err
				}
				if _, err := fmt.Fprintln(f, string(b)); err != nil {
					return errors.Wrap(err, "could not write line")
				}
				return nil
			})
		}
		return nil
	})
	return tmb.Wait()
}

func (Ω *localFileCache) Close() error {

	for _, filePointer := range Ω.filePointers {
		if err := filePointer.Close(); err != nil {
			return errors.Wrap(err, "Could not close local prediction file")
		}
	}

	return nil
}
