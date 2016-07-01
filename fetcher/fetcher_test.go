package main

import (
	"github.com/bitly/go-nsq"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"bitbucket.org/heindl/species/store"
	"bitbucket.org/heindl/nsqeco"
)

func TestTaxonFetcher(t *testing.T) {

	t.Parallel()

	Convey("given an environment", t, func() {
		producer := &nsqeco.MockProducer{}
		store := store.NewMockStore(t)

		Convey("should handle taxon fetch message without error", func() {

			id := nsq.MessageID{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', 'a', 's', 'd', 'f', 'g', 'h'}

			fetcher := &SpeciesFetchHandler{
				NSQProducer: producer,
				SpeciesStore:    store,
			}
			So(fetcher.HandleMessage(nsq.NewMessage(id, []byte("Limenitis arthemis"))), ShouldBeNil)

			So(producer.Count(nsqeco.NSQOccurrenceFetch), ShouldEqual, 25)
			So(producer.Count(nsqeco.NSQSpeciesMetaFetch), ShouldEqual, 5)
			specs, err := store.Read()
			So(err, ShouldBeNil)
			So(len(specs), ShouldEqual, 5)
			for _, n := range []string{
				"Limenitis arthemis arizonensis",
				"Limenitis arthemis rubrofasciata",
				"Limenitis arthemis albofasciata",
				"Limenitis arthemis arthemis",
				"Limenitis arthemis astyanax",
			} {
				var exists bool
				for _, s := range specs {
					if n == string(s.CanonicalName) {
						exists = true
					}
				}
				So(exists, ShouldBeTrue)
			}

			//So(specs[0].CanonicalName, ShouldEqual, "Limenitis arthemis arizonensis")
			//So(len(specs[0].Sources), ShouldEqual, 1)
			//So(specs[0].Sources[0].IndexKey, ShouldEqual, 6225972)

		})

		Reset(func() {
			store.Close()
		})

	})
}
