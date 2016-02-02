package metafetcher

import (
	"bitbucket.org/heindl/cxt"
	"bitbucket.org/heindl/species"
	"encoding/json"
	"github.com/bitly/go-nsq"
	"github.com/omidnikta/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

func TestMetaFetcher(t *testing.T) {

	t.Parallel()

	Convey("given an environment", t, func() {

		producer := &cxt.TestProducer{}
		server, session := cxt.TestMongo(t)

		Convey("should handle eol fetch message without error", func() {

			c := &cxt.Context{
				Producer: producer,
				Mongo:    session,
				Log:      logrus.StandardLogger(),
			}
			c.Log.ShowCaller(true)

			id := nsq.MessageID{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', 'a', 's', 'd', 'f', 'g', 'h'}

			b, err := json.Marshal("Limenitis arthemis astyanax")
			So(err, ShouldBeNil)
			fetcher := &MetaFetcher{
				Context: c,
			}
			So(fetcher.HandleMessage(nsq.NewMessage(id, b)), ShouldBeNil)

			Convey("should return expected eol page", func() {
				var list []species.Species
				So(c.Mongo.Coll(cxt.SpeciesColl).Find(bson.M{}).All(&list), ShouldBeNil)
				So(len(list), ShouldEqual, 1)
				So(list[0].CanonicalName, ShouldEqual, species.CanonicalName("Limenitis arthemis astyanax"))
				So(len(list[0].Sources), ShouldEqual, 1)
				So(list[0].Image, ShouldNotBeNil)
				So(list[0].Image.Value, ShouldEqual, "http://media.eol.org/content/2014/05/12/09/43805_orig.jpg")
				So(list[0].Description, ShouldNotBeNil)
				So(list[0].Description.Value, ShouldEqual, "<i>Limenitis arthemis astyanax</i> is resident to eastern Canada and the northeastern US., migratory to the west in some parts of its range. This is the southern sub-species of this species and mimics the pipevine swallowtail, with which it co-occurs.  The northern subspecies, <i>L. a. arthemis</i>, does not mimic the pipevine and so shows a different color pattern  (Scott 1986).  Habitats are deciduous wooded areas.  Host plants are trees species from families including Saliaceae, Betulaceae, Rosaceae, Fagaceae, Ulmaceae and Tiliaceae.  Eggs are laid on the host plant singly.  Individuals overwinter in a hibernaculum as third instar larvae.  There is a variable number flights each year depending on latitude, with one flight approximately July15-Aug15 in the northern part of the range and multiple flights in the south between Mar 1-Nov. 30 (Scott 1986).")
			})

		})

		Reset(func() {
			session.Close()
			server.Stop()
		})

	})
}
