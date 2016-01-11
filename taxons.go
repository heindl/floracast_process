package taxon

import (
	"bitbucket.org/heindl/cxt"
	"gopkg.in/mgo.v2/bson"
	"github.com/dropbox/godropbox/errors"
	"encoding/binary"
	"strconv"
	"time"
)

type Key int

func init() {
	cxt.RegisterCollection(cxt.CollectionInitiator{
		Name: cxt.TaxonColl,
	})
}

func NewKeyFromBytes(b []byte) Key {
	k, _ := strconv.ParseInt(string(b), 10, 64)
	return Key(k)
}

func NewKeyFromString(b string) (Key, error) {
	s, err := strconv.Atoi(b)
	if err != nil {
		return Key(0), errors.Wrapf(err, "could not parse string: %s", b)
	}
	return Key(s), nil
}

func (k Key) Bytes() []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(k))
	return b
}

func (k Key) String() string {
	return strconv.Itoa(int(k))
}

func AllKeys(c *cxt.Context) ([]Key, error) {
	var taxons []Taxon
	if err := c.Mongo.Coll(cxt.TaxonColl).Find(bson.M{}).All(&taxons); err != nil {
		return nil, errors.Wrap(err, "could not find taxon list")
	}
	var response []Key
	for _, taxon := range taxons {
		response = append(response, taxon.ID)
	}
	return response, nil
}

type Taxon struct {
	ID           Key    `bson:"_id" json:"_id"`
	LastModified time.Time `bson:"lastModified" json:"lastModified"`
}
