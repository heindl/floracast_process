package parser

import (
	"bitbucket.org/heindl/taxa/store"
	"time"
	"google.golang.org/genproto/googleapis/type/latlng"
	"bitbucket.org/heindl/taxa/utils"
	"context"
	"gopkg.in/tomb.v2"
	"github.com/saleswise/errors/errors"
	"fmt"
)

type Writer interface{
	WritePredictionLine(p store.Prediction) error
	Close() error
}

type PredictionParser interface {
	FetchWritePredictions(cxt context.Context, taxon store.TaxonID, date string) error
}

func NewPredictionParser(cxt context.Context, gcsBucketName string, writer Writer) (PredictionParser, error) {

	taxastore, err := store.NewTaxaStore()
	if err != nil {
		return nil, err
	}

	gcsFetcher, err := NewGCSFetcher(cxt, gcsBucketName)
	if err != nil {
		return nil, err
	}

	return &predictionParser{
		WildernessAreaFetcher: NewWildernessAreaFetcher(taxastore),
		GCSFetcher: gcsFetcher,
		Writer: writer,
	}, nil
}

type predictionParser struct {
	WildernessAreaFetcher WildernessAreaFetcher
	GCSFetcher            GCSFetcher
	Writer                Writer
}

func (Ω *predictionParser) FetchWritePredictions(cxt context.Context, taxon store.TaxonID, date string) error {

	gcsFilePaths, err := Ω.GCSFetcher.FetchLatestPredictionFileNames(cxt, taxon, date)
	if err != nil {
		return err
	}

	fmt.Println("gcsFilePaths", len(gcsFilePaths))

	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _, _fpath := range gcsFilePaths {
			fpath := _fpath
			tmb.Go(func() error {
				predictions, err := Ω.GCSFetcher.FetchPredictions(cxt, fpath)
				if err != nil {
					return err
				}
				for _, p := range predictions {
					if p.Target <= p.Random {
						continue
					}
					wa, err := Ω.WildernessAreaFetcher.GetWildernessArea(cxt, p.Latitude, p.Longitude)
					if err != nil {
						return err
					}
					d, err := time.ParseInLocation("20060102", p.Date, time.UTC)
					if err != nil {
						return errors.Wrap(err, "could not parse date")
					}
					return Ω.Writer.WritePredictionLine(store.Prediction{
						CreatedAt: utils.TimePtr(time.Now()),
						Location: latlng.LatLng{p.Latitude, p.Longitude},
						PredictionValue: p.Target,
						TaxonID: taxon,
						Date: utils.TimePtr(d),
						FormattedDate: p.Date,
						Month: d.Month(),
						WildernessAreaID: wa.ID,
						WildernessAreaName: wa.Name,
					})
				}
				return nil
			})
		}
		return nil
	})
	return tmb.Wait()
}