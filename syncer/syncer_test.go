package fetcher

import (
	"bitbucket.org/heindl/cxt"
	"bitbucket.org/heindl/species"
	"bitbucket.org/heindl/utils"
	"encoding/json"
	"github.com/bitly/go-nsq"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

func TestTaxonFetcher(t *testing.T) {

	t.Parallel()

	Convey("given an environment", t, func() {

		producer := &cxt.TestProducer{}
		server, session := cxt.TestMongo(t)

		Convey("should handle taxon fetch message without error", func() {

			c := &cxt.Context{
				Producer: producer,
				Mongo:    session,
			}

			id := nsq.MessageID{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', 'a', 's', 'd', 'f', 'g', 'h'}

			b, err := json.Marshal("Limenitis arthemis")
			So(err, ShouldBeNil)

			fetcher := &SpeciesFetcher{
				Context: c,
			}
			So(fetcher.HandleMessage(nsq.NewMessage(id, b)), ShouldBeNil)

			Convey("should send expected number of nsq message and contain three records", func() {
				So(producer.Count(cxt.NSQOccurrencesFetch), ShouldEqual, 26)
				So(producer.Count(cxt.NSQSpeciesDataFetch), ShouldEqual, 6)
				var specs []species.Species
				So(c.Mongo.Coll(cxt.SpeciesColl).Find(bson.M{}).All(&specs), ShouldBeNil)
				So(len(specs), ShouldEqual, 6)
				So(specs[0].CanonicalName, ShouldEqual, "Limenitis arthemis virithemis")
				So(len(specs[0].Sources), ShouldEqual, 1)
				So(specs[0].Sources[0].IndexKey, ShouldEqual, 6225972)
			})

		})

		Reset(func() {
			session.Close()
			server.Stop()
		})

	})
}
