package nameusage

import (
	"bitbucket.org/heindl/process/utils"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var usageJSON = []byte(`
{
  "ID": "",
  "Name": {
    "Rank": "species",
    "ScientificName": "cantharellus jebbi"
  },
  "ScientificNames": null,
  "Sources": {
    "14": {
      "133044": {
        "Name": {
          "Rank": "species",
          "ScientificName": "cantharellus jebbi"
        },
      }
    },
    "26": {
      "133044": {
        "Name": {
          "Rank": "species",
          "ScientificName": "cantharellus jebbi"
        },
      }
    },
    "27": {
      "5184832": {
        "TaxonomicReference": true,
        "Name": {
          "Rank": "species",
          "ScientificName": "cantharellus jebbi"
        },
      }
    },
    "INAT": {
      "96709": {
        "TaxonomicReference": true,
        "Name": {
          "Rank": "species",
          "ScientificName": "cantharellus jebbi"
        },
        "ModifiedAt": "2018-02-10T14:21:10.599344-05:00",
        "CreatedAt": "2018-02-10T14:21:10.599342-05:00"
      }
    }
  },
  "CreatedAt": "2018-02-10T14:21:10.599349-05:00",
  "ModifiedAt": "2018-02-10T14:21:13.080988-05:00"
}
`)

func TestNameUsage(t *testing.T) {

	t.Parallel()

	Convey("should parse nameusage", t, func() {
		id, err := newNameUsageID()
		So(err, ShouldBeNil)

		u, err := FromJSON(id, usageJSON)
		So(err, ShouldBeNil)
		fmt.Println(utils.JsonOrSpew(u))

	})
}
