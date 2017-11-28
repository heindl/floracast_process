package filecache

import (
	"strings"
	"fmt"
	"bitbucket.org/heindl/taxa/store"
	"os"
	"path"
	"github.com/saleswise/errors/errors"
	"github.com/elgs/gostrgen"
)

type FileCache struct{
	FilePath string
}

func NewFileCache() (FileCache, error) {


	random_string, err := gostrgen.RandGen(10, gostrgen.Lower, "", "")
	if err != nil {
		return "", errors.Wrap(err, "could not create temp file random string name")
	}
	tmp := path.Join(os.TempDir(), fmt.Sprintf("predictions-%s", random_string))
	if err := os.Mkdir(tmp, os.ModePerm); err != nil {
		return "", errors.Wrap(err, "could not create tmp path")
	}
	return tmp, nil

	fmt.Println("TEMP_FILE_DIRECTORY", tmp)
	return &FileCache{
		FilePath: tmp,
	}, nil
}

func (Ω *FileCache) WritePredictionLine(p store.Prediction) error {
	if err := Ω.writeTaxaFile(p); err != nil {
		return err
	}
	return Ω.writeTaxonFile(p)
}

func (*FileCache) Close() error {
	return nil
}

func (Ω *FileCache) writeTaxonFile(p store.Prediction) error {
	// First write predictions.
	taxonFile, err := os.OpenFile(path.Join(Ω.FilePath, string(p.TaxonID)+".taxon"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "could not open taxon file")
	}
	defer taxonFile.Close()

	line := fmt.Sprintf("%s,%s,%.6f",
		p.FormattedDate,
		p.WildernessAreaID,
		p.PredictionValue)

	line = strings.Replace(line, ".", "|", -1)

	line = fmt.Sprintf("{\"%s\": [%.6f, %.6f]}\n", line, p.Location.Latitude, p.Location.Longitude)

	if _, err := taxonFile.WriteString(line); err != nil {
		return errors.Wrap(err, "could not write taxon file")
	}
	return nil
}
func (Ω *FileCache) writeTaxaFile(p store.Prediction) error {
	taxaFile, err := os.OpenFile(path.Join(Ω.FilePath, "/species.txt"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer taxaFile.Close()
	// TODO: There is the potential tha two percentiles could be exactly the same in which case one would override the other, though this seems unlikely
	line := fmt.Sprintf("%s,%s,%.6f,%.6f,%.6f",
		p.TaxonID,
		p.FormattedDate,
		p.PredictionValue,
		p.PercentileOverAllTaxonPredictions,
		p.PercentileOverAllTaxaPredictionsForDay)
	line = strings.Replace(line, ".", "|", -1)
	line = fmt.Sprintf("{\"%s\": [%.6f, %.6f]}\n", line, p.Location.Latitude, p.Location.Longitude)

	if _, err := taxaFile.WriteString(line); err != nil {
		return errors.Wrap(err, "could not write taxa file")
	}
	return nil
}
