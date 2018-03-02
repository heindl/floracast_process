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
	"bitbucket.org/heindl/process/utils"
	"gopkg.in/tomb.v2"
)

// Season is an object to handle Season based date calculations.
type Season struct {
	Min, Max int64
}

// Spring is March through May
var Spring = Season{
	Min: time.Date(2011, time.March, 1, 0, 0, 0, 0, time.UTC).Unix(),
	Max: time.Date(2011, time.May, 31, 0, 0, 0, 0, time.UTC).Unix(),
}

// Summer is June through August
var Summer = Season{
	Min: time.Date(2011, time.June, 1, 0, 0, 0, 0, time.UTC).Unix(),
	Max: time.Date(2011, time.August, 31, 0, 0, 0, 0, time.UTC).Unix(),
}

// Autumn is September through November
var Autumn = Season{
	Min: time.Date(2011, time.September, 1, 0, 0, 0, 0, time.UTC).Unix(),
	Max: time.Date(2011, time.November, 30, 0, 0, 0, 0, time.UTC).Unix(),
}

// Winter is December through February
var Winter = Season{
	Min: time.Date(2010, time.December, 1, 0, 0, 0, 0, time.UTC).Unix(),
	Max: time.Date(2011, time.February, 28, 0, 0, 0, 0, time.UTC).Unix(),
}

type randomOccurrenceGenerator struct {
	gridGenerator        grid.Generator
	occurrenceAggregator *OccurrenceAggregation
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

// GenerateRandomOccurrences generates an evenly distributed list of Random Occurrences
// across southern Canada, northern Mexico, and the United States.
// There will be one for each season.
func GenerateRandomOccurrences(gridLevel, numberOfBatches int) (*OccurrenceAggregation, error) {
	gen, err := newRandomOccurrenceGenerator()
	if err != nil {
		return nil, err
	}

	bounds, err := gen.gridGenerator.SubDivide(grid.NorthAmerica, gridLevel)
	if err != nil {
		panic(err)
	}

	recordNumberCounter := 1

	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for 𝝨 := 1; 𝝨 <= numberOfBatches; 𝝨++ {
			batch := 𝝨
			for _, _bound := range bounds {
				bound := *_bound
				for _, _season := range []Season{Spring, Summer, Autumn, Winter} {
					season := _season
					tmb.Go(func() error {
						gen.Lock()
						recordNumber := recordNumberCounter
						recordNumberCounter++
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

		seconds := rand.Int63n(timeDelta) + season.Min
		dateStr := time.Unix(seconds, 0).Format("20060102")

		lat := (float64(rand.Intn(yDelta)) * rand.Float64()) + bounds.South
		lng := (float64(rand.Intn(xDelta)) * rand.Float64()) + bounds.West

		_, err := cache.FetchEcologicalRegion(lat, lng)
		if utils.ContainsError(err, ecoregions.ErrNotFound) {
			continue
		}
		if err != nil {
			return err
		}

		o, err := NewOccurrence(datasources.TypeRandom, datasources.TargetID(strconv.Itoa(batch)), strconv.Itoa(recordNumber))
		if err != nil {
			return err
		}

		if err := o.SetGeoSpatial(lat, lng, dateStr, false); err != nil {
			return err
		}

		err = Ω.occurrenceAggregator.AddOccurrence(o)
		if utils.ContainsError(err, ErrCollision) {
			continue
		}
		if err != nil {
			return err
		}

		break

	}

	return nil
}
