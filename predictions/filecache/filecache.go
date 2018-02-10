package filecache

import (
	"bitbucket.org/heindl/processors/store"
	"encoding/json"
	"fmt"
	"github.com/elgs/gostrgen"
	"github.com/dropbox/godropbox/errors"
	"os"
	"path"
)

type FileCache struct {
	FilePath string
}

func NewFileCache() (*FileCache, error) {
	random_string, err := gostrgen.RandGen(10, gostrgen.Lower, "", "")
	if err != nil {
		return nil, errors.Wrap(err, "could not create temp file random string name")
	}
	tmp := path.Join("/tmp/", fmt.Sprintf("predictions-%s", random_string))
	if err := os.Mkdir(tmp, os.ModePerm); err != nil {
		return nil, errors.Wrap(err, "could not create tmp path")
	}
	fmt.Println("TEMP_FILE_DIRECTORY", tmp)
	return &FileCache{
		FilePath: tmp,
	}, nil
}

func (Ω *FileCache) ReadTaxa(lat, lng, radius float64, qDate string, taxon string) ([]string, error) {
	return nil, nil
}

func (Ω *FileCache) WritePredictionLine(p *store.Prediction) error {

	b, err := json.Marshal(p)
	if err != nil {
		return errors.Wrap(err, "could not marshal prediction")
	}
	filename := fmt.Sprintf("%s-%s.jsonl", p.FormattedDate, p.TaxonID)
	taxonFile, err := os.OpenFile(path.Join(Ω.FilePath, filename), os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "could not open taxon file")
	}
	defer taxonFile.Close()

	if _, err := fmt.Fprintln(taxonFile, string(b)); err != nil {
		return errors.Wrap(err, "could not write line")
	}

	return nil
}

func (*FileCache) Close() error {
	return nil
}
