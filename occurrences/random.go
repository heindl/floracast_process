package occurrences

import (
	"math"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/terra/ecoregions"
	"bitbucket.org/heindl/process/terra/ecoregions/cache"
	"bitbucket.org/heindl/process/terra/grid"
	"gopkg.in/tomb.v2"
)

type Season struct {
	Min, Max int64
}

var Spring = Season{
	Min: time.Date(2011, time.March, 1, 0, 0, 0, 0, time.UTC).Unix(),
	Max: time.Date(2011, time.May, 31, 0, 0, 0, 0, time.UTC).Unix(),
}

var Summer = Season{
	Min: time.Date(2011, time.June, 1, 0, 0, 0, 0, time.UTC).Unix(),
	Max: time.Date(2011, time.August, 31, 0, 0, 0, 0, time.UTC).Unix(),
}

var Autumn = Season{
	Min: time.Date(2011, time.September, 1, 0, 0, 0, 0, time.UTC).Unix(),
	Max: time.Date(2011, time.November, 30, 0, 0, 0, 0, time.UTC).Unix(),
}

var Winter = Season{
	Min: time.Date(2010, time.December, 1, 0, 0, 0, 0, time.UTC).Unix(),
	Max: time.Date(2011, time.February, 28, 0, 0, 0, 0, time.UTC).Unix(),
}

type randomOccurrenceGenerator struct {
	gridGenerator        grid.Generator
	occurrenceAggregator *OccurrenceAggregation
	counter              int
	sync.Mutex
}

func newRandomOccurrenceGenerator() (*randomOccurrenceGenerator, error) {

	rand.Seed(time.Now().Unix())

	gridGenerator, err := grid.NewGridGenerator()
	if err != nil {
		return nil, err
	}

	return &randomOccurrenceGenerator{
		gridGenerator:        gridGenerator,
		occurrenceAggregator: NewOccurrenceAggregation(),
	}, nil
}

func GenerateRandomOccurrences(number float64) (*OccurrenceAggregation, error) {
	gen, err := newRandomOccurrenceGenerator()
	if err != nil {
		return nil, err
	}

	bounds, err := gen.gridGenerator.SubDivide(grid.NorthAmerica, 4)
	if err != nil {
		panic(err)
	}

	numberOfBatches := math.Ceil(number / float64(len(bounds)*4))

	counter := 1

	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _batch := 1; _batch <= int(numberOfBatches); _batch++ {
			batch := _batch
			for _, _bound := range bounds {
				bound := *_bound
				for _, _season := range []Season{Spring, Summer, Autumn, Winter} {
					season := _season
					tmb.Go(func() error {
						gen.Lock()
						recordNumber := counter
						counter += 1
						gen.Unlock()
						return gen.generateRandomOccurrence(batch, recordNumber, bound, season)
					})
				}
			}
		}
		return nil
	})
	if err := tmb.Wait(); err != nil {
		return nil, err
	}

	return gen.occurrenceAggregator, nil
}

func (Ω *randomOccurrenceGenerator) generateRandomOccurrence(batch, recordNumber int, bounds grid.Bound, season Season) error {

	timeDelta := season.Max - season.Min
	xDelta := int(math.Ceil(bounds.East) - math.Ceil(bounds.West))
	yDelta := int(math.Ceil(bounds.North) - math.Ceil(bounds.South))

	for {
		//rand.Seed(time.Now().Unix())

		seconds := rand.Int63n(timeDelta) + int64(season.Min)
		dateStr := time.Unix(seconds, 0).Format("20060102")

		lat := (float64(rand.Intn(yDelta)) * rand.Float64()) + bounds.South
		lng := (float64(rand.Intn(xDelta)) * rand.Float64()) + bounds.West

		_, err := cache.FetchEcologicalRegion(lat, lng)
		if err == ecoregions.ErrNotFound {
			continue
		}
		if err != nil {
			return err
		}

		o, err := NewOccurrence(datasources.TypeRandom, datasources.TargetID(strconv.Itoa(batch)), strconv.Itoa(recordNumber))
		if err != nil {
			return err
		}

		if err := o.SetGeospatial(lat, lng, dateStr, false); err != nil {
			return err
		}

		err = Ω.occurrenceAggregator.AddOccurrence(o)
		if err == ErrCollision {
			continue
		}
		if err != nil {
			return err
		}

		break

	}

	return nil
}
