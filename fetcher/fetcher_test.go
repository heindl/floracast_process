package fetcher

import (
	"bitbucket.org/heindl/cxt"
	"encoding/json"
	"github.com/bitly/go-nsq"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2/bson"
	"testing"
	"bitbucket.org/heindl/taxon"
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

			b, err := json.Marshal("1937885,5133635,1892729")
			So(err, ShouldBeNil)

			processor := &TaxonFetcher{
				Context:  c,
			}
			So(processor.HandleMessage(nsq.NewMessage(id, b)), ShouldBeNil)

			Convey("should send expected number of nsq message and contain three records", func() {
				So(producer.Count(cxt.NSQOccurrencesFetch), ShouldEqual, 3)
				var txns []taxon.Taxon
				So(c.Mongo.Coll(cxt.TaxonColl).Find(bson.M{}).All(&txns), ShouldBeNil)
				So(len(txns), ShouldEqual, 3)
				So(txns[0].ID, ShouldEqual, taxon.Key(1937885))
				So(txns[0].LastModified.IsZero(), ShouldEqual, false)
			})

			Convey("message reprocess should maintain same number of database results", func() {
				So(processor.HandleMessage(nsq.NewMessage(id, b)), ShouldBeNil)
				So(producer.Count(cxt.NSQOccurrencesFetch), ShouldEqual, 6)
				var txns []taxon.Taxon
				So(c.Mongo.Coll(cxt.TaxonColl).Find(bson.M{}).All(&txns), ShouldBeNil)
				So(len(txns), ShouldEqual, 3)
				So(txns[0].ID, ShouldEqual, taxon.Key(1937885))
				So(txns[0].LastModified.IsZero(), ShouldEqual, false)
			})

		})

		Reset(func() {
			session.Close()
			server.Stop()
		})

	})
}
