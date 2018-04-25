package occurrence

import (
	"math"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"bitbucket.org/heindl/process/datasources"
	"bitbucket.org/heindl/process/terra/ecoregions"
	"bitbucket.org/heindl/process/terra/ecoregions/cache"
	"bitbucket.org/heindl/process/terra/geo"
	"bitbucket.org/heindl/process/terra/grid"
	"bitbucket.org/heindl/process/utils"
	"fmt"
	"github.com/golang/glog"
	"gopkg.in/tomb.v2"
)

// maxGenerationAttemps is the max number of times to try to generate a random occurrence.
// A safeguard against cases in which a grid does not contain an EcoRegion.
const maxGenerationAttempts = 20

// Season is an object to handle Season based date calculations.
type Season struct {
	Min, Max int64
}

// Spring is March through May
var Spring = Season{
	Min: time.Date(2016, time.March, 1, 0, 0, 0, 0, time.UTC).Unix(),
	Max: time.Date(2016, time.May, 31, 0, 0, 0, 0, time.UTC).Unix(),
}

// Summer is June through August
var Summer = Season{
	Min: time.Date(2016, time.June, 1, 0, 0, 0, 0, time.UTC).Unix(),
	Max: time.Date(2016, time.August, 31, 0, 0, 0, 0, time.UTC).Unix(),
}

// Autumn is September through November
var Autumn = Season{
	Min: time.Date(2016, time.September, 1, 0, 0, 0, 0, time.UTC).Unix(),
	Max: time.Date(2016, time.November, 30, 0, 0, 0, 0, time.UTC).Unix(),
}

// Winter is December through February
var Winter = Season{
	Min: time.Date(2015, time.December, 1, 0, 0, 0, 0, time.UTC).Unix(),
	Max: time.Date(2016, time.February, 28, 0, 0, 0, 0, time.UTC).Unix(),
}

type randomOccurrenceGenerator struct {
	gridGenerator        grid.Generator
	occurrenceAggregator *Aggregation
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
func GenerateRandomOccurrences(gridLevel, batch int) (*Aggregation, error) {
	gen, err := newRandomOccurrenceGenerator()
	if err != nil {
		return nil, err
	}

	bounds, err := gen.gridGenerator.SubDivide(grid.NorthAmerica, gridLevel)
	if err != nil {
		panic(err)
	}
	batchSize := len(bounds) * 4

	glog.Infof("Generating %d Random Occurrences in Batch %d for %d Bounds", batchSize, batch, len(bounds))

	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for _boundCount, _bound := range bounds {
			// 0, 1, 2 ...
			boundCount := _boundCount
			bound := *_bound

			tmb.Go(func() error {
				for _seasonCount, _season := range []Season{Spring, Summer, Autumn, Winter} {
					// 0, 1, 2, 3
					season := _season
					seasonCount := _seasonCount
					tmb.Go(func() error {
						err := gen.generateRandomOccurrence(
							batch,
							(boundCount*4)+seasonCount+1,
							bound,
							season)
						if utils.ContainsError(err, errMaxRandomGenerationAttempts) {
							glog.Warningf("Max Attempts [%d] Reached while Generating Random Occurrences Point [%v]", maxGenerationAttempts, bound)
							return nil
						}
						return err
					})
				}
				return nil
			})
		}
		return nil
	})
	if err := tmb.Wait(); err != nil {
		return nil, err
	}

	glog.Infof("Uploading %d Random Occurrences Points", gen.occurrenceAggregator.Count())

	return gen.occurrenceAggregator, nil
}

var errMaxRandomGenerationAttempts = fmt.Errorf("reached maximum attempts at generating random occurrence point")

func (Ω *randomOccurrenceGenerator) generateRandomOccurrence(batch, recordNumber int, bounds geo.Bound, season Season) error {

	timeDelta := season.Max - season.Min
	xDelta := int(math.Ceil(bounds.East) - math.Ceil(bounds.West))
	yDelta := int(math.Ceil(bounds.North) - math.Ceil(bounds.South))

	attempts := 0

	for {
		attempts++
		if attempts > maxGenerationAttempts {
			return errMaxRandomGenerationAttempts
		}
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

		// Including all of this information so that the transform function can infer the size of the batch,
		// to correctly balance against taxon occurrence count.
		occurrenceID := fmt.Sprintf("%d", recordNumber)

		o, err := NewOccurrence(nil, datasources.TypeRandom, datasources.TargetID(strconv.Itoa(batch)), occurrenceID)
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
