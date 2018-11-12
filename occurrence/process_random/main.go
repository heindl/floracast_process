package main

import (
	"github.com/heindl/floracast_process/occurrence"
	"github.com/heindl/floracast_process/store"
	"github.com/heindl/floracast_process/utils"
	"context"
	"flag"
	"github.com/golang/glog"
	"gopkg.in/tomb.v2"
)

// Level 4: 216
// Level 3: 64

func main() {
	batches := flag.Int("batches", 0, "Random point batches")
	cellLevel := flag.Int("level", 0, "S2 cell level")
	flag.Parse()

	ctx := context.Background()

	floraStore, err := store.NewFloraStore(ctx)
	if err != nil {
		panic(err)
	}

	//if err = occurrence.ClearRandomPoints(ctx, floraStore); err != nil {
	//	panic(err)
	//}

	cxt := context.Background()

	limiter := utils.NewLimiter(5)
	tmb := tomb.Tomb{}
	tmb.Go(func() error {
		for i := 1; i <= *batches; i++ {
			batch_number := i
			release := limiter.Go()
			tmb.Go(func() error {
				defer release()
				aggr, err := occurrence.GenerateRandomOccurrences(*cellLevel, batch_number)
				if err != nil {
					return err
				}
				glog.Infof("Uploading Random Batch [%d] with %d Occurrences", batch_number, aggr.Count())
				return aggr.UnsafeUpload(cxt, floraStore)
			})
		}
		return nil
	})
	if err := tmb.Wait(); err != nil {
		panic(err)
	}

}
