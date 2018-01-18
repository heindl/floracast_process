package terra

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	//"github.com/golang/geo/s1"
	//"github.com/golang/geo/s2"
	//"fmt"
	"github.com/golang/geo/s2"
	"github.com/golang/geo/s1"
	"fmt"
)

func TestEcoRegionFetch(t *testing.T) {

	t.Parallel()

	SkipConvey("should explain routes", t, func() {
		concentricLoopsPolygon := func(center s2.Point, numLoops, verticesPerLoop int) *s2.Polygon {
			var loops []*s2.Loop
			for li := 0; li < numLoops; li++ {
				radius := s1.Angle(0.005 * float64(li+1) / float64(numLoops))
				l := s2.RegularLoop(center, radius, verticesPerLoop+li)
				loops = append(loops, l)
			}
			for i, l := range loops {
				fmt.Println(i, len(l.Vertices()), isPositivelyOriented(l))
			}
			polygon := s2.PolygonFromLoops(loops)
			for i, p := range polygon.Loops() {
				fmt.Println(i, p.IsHole(), len(p.Vertices()), isPositivelyOriented(p))
			}
			return polygon
		}
		So(concentricLoopsPolygon(s2.PointFromLatLng(s2.LatLngFromDegrees(33.745252,-118.0801775)), 10, 10), ShouldBeNil)
	})

	Convey("should parse geojson", t, func() {

		seal_beach := []byte(`{  
   "type":"Feature",
   "properties":{  
      "Category":"Designation",
      "d_Category":"Designation",
      "Own_Type":"FED",
      "d_Own_Type":"Federal",
      "Own_Name":"FWS",
      "d_Own_Name":"U.S. Fish & Wildlife Service",
      "Loc_Own":"U.S. Fish and Wildlife Service",
      "Mang_Type":"FED",
      "d_Mang_Typ":"Federal",
      "Mang_Name":"FWS",
      "d_Mang_Nam":"U.S. Fish & Wildlife Service",
      "Loc_Mang":"U.S. Fish and Wildlife Service",
      "Des_Tp":"MPA",
      "d_Des_Tp":"Marine Protected Area",
      "Loc_Ds":"National Wildlife Refuge",
      "Unit_Nm":"Seal Beach National Wildlife Refuge",
      "Loc_Nm":"Seal Beach National Wildlife Refuge",
      "State_Nm":"CA",
      "d_State_Nm":"California",
      "Agg_Src":"NOAA_PADUS1_4MPA_MPAIMember_Eligble2016",
      "GIS_Src":"MPAI_2013_padus_revised.gdb\/mpa_inventory_2013_NSMember_Eligible_revised",
      "Src_Date":"2014\/02\/25",
      "GIS_Acres":967,
      "Source_PAI":"NWR103",
      "WDPA_Cd":75128,
      "Access":"UK",
      "d_Access":"Unknown",
      "Access_Src":"GAP - Default",
      "GAP_Sts":"1",
      "d_GAP_Sts":"1 - managed for biodiversity - disturbance events proceed or are mimicked",
      "GAPCdSrc":"GAP - NOAA",
      "GAPCdDt":"2016",
      "IUCN_Cat":"Ia",
      "d_IUCN_Cat":"Ia: Strict nature reserves",
      "IUCNCtSrc":"GAP - NOAA",
      "IUCNCtDt":"2016",
      "Date_Est":"1974",
      "Comments":null
   },
   "geometry":{  
      "type":"MultiPolygon",
      "coordinates":[  
         [
            [  
               [  
                  -118.069896698325621,
                  33.749107649377173
               ],
               [  
                  -118.06991890813751,
                  33.748788732670299
               ],
               [  
                  -118.069948075154912,
                  33.748420347641229
               ],
               [  
                  -118.069947086748073,
                  33.746869246239491
               ],
               [  
                  -118.069902051561627,
                  33.745070423240357
               ],
               [  
                  -118.070775898664323,
                  33.745063532472329
               ],
               [  
                  -118.069009652105166,
                  33.744092658141874
               ],
               [  
                  -118.069030201941445,
                  33.738960926266003
               ],
               [  
                  -118.068629524564329,
                  33.738953534840057
               ],
               [  
                  -118.068162294572247,
                  33.739055837795604
               ],
               [  
                  -118.067709367010806,
                  33.739009696111061
               ],
               [  
                  -118.066769802976623,
                  33.739027253671622
               ],
               [  
                  -118.065908983033736,
                  33.739056178945631
               ],
               [  
                  -118.065245540283726,
                  33.739047533150412
               ],
               [  
                  -118.064016889496727,
                  33.739069199233221
               ],
               [  
                  -118.060292213276384,
                  33.739045934000913
               ],
               [  
                  -118.060260038202671,
                  33.737307666949711
               ],
               [  
                  -118.061264511322577,
                  33.736558920103448
               ],
               [  
                  -118.061350044220504,
                  33.73074544214122
               ],
               [  
                  -118.067328239479735,
                  33.730669345229259
               ],
               [  
                  -118.0677311683379,
                  33.730385230388407
               ],
               [  
                  -118.075779865950295,
                  33.730396917680089
               ],
               [  
                  -118.076523584442029,
                  33.730397971766216
               ],
               [  
                  -118.08130864042414,
                  33.730404642594422
               ],
               [  
                  -118.081301176590486,
                  33.730906338482434
               ],
               [  
                  -118.081292378543409,
                  33.731497706109849
               ],
               [  
                  -118.08128139269256,
                  33.732236139098092
               ],
               [  
                  -118.084421298327129,
                  33.732239593766963
               ],
               [  
                  -118.084472539985896,
                  33.732410340359635
               ],
               [  
                  -118.084602480164889,
                  33.732597949232023
               ],
               [  
                  -118.084659617827938,
                  33.732856729314996
               ],
               [  
                  -118.08478884169368,
                  33.733137840112136
               ],
               [  
                  -118.084839326443657,
                  33.733407589716542
               ],
               [  
                  -118.084936340382328,
                  33.733606047766344
               ],
               [  
                  -118.085078873034306,
                  33.733865217692461
               ],
               [  
                  -118.085148728805436,
                  33.734179060271657
               ],
               [  
                  -118.085193193761071,
                  33.734377276309246
               ],
               [  
                  -118.085296189478953,
                  33.734652767304986
               ],
               [  
                  -118.085385795334091,
                  33.734961196551858
               ],
               [  
                  -118.085553889614758,
                  33.735313991158151
               ],
               [  
                  -118.085651161922385,
                  33.735479446723062
               ],
               [  
                  -118.085800436416875,
                  33.735716644739135
               ],
               [  
                  -118.085943144060394,
                  33.735953814627109
               ],
               [  
                  -118.086066438855809,
                  33.736152391862127
               ],
               [  
                  -118.086202453353621,
                  33.736406032262373
               ],
               [  
                  -118.086358551635001,
                  33.73661025880363
               ],
               [  
                  -118.086317581301145,
                  33.736813584165077
               ],
               [  
                  -118.086369630104983,
                  33.736879827014072
               ],
               [  
                  -118.086525939437024,
                  33.737056552776274
               ],
               [  
                  -118.086643300083765,
                  33.737172597866248
               ],
               [  
                  -118.086812667240039,
                  33.73736038523338
               ],
               [  
                  -118.087061071920516,
                  33.737521029707743
               ],
               [  
                  -118.087374371005069,
                  33.737786479136766
               ],
               [  
                  -118.087629349231591,
                  33.737947153146237
               ],
               [  
                  -118.087897547665222,
                  33.738096886759699
               ],
               [  
                  -118.088224620019474,
                  33.738279889820703
               ],
               [  
                  -118.088610731493148,
                  33.738474164254384
               ],
               [  
                  -118.089023203624663,
                  33.738657553921541
               ],
               [  
                  -118.08931785297365,
                  33.738785403681305
               ],
               [  
                  -118.081449610792561,
                  33.749629022961528
               ],
               [  
                  -118.077221439558897,
                  33.747541313076965
               ],
               [  
                  -118.075409043054194,
                  33.74998564008451
               ],
               [  
                  -118.073499182230918,
                  33.752738823020366
               ],
               [  
                  -118.071286208628734,
                  33.75150194226083
               ],
               [  
                  -118.070876662261028,
                  33.751274106004388
               ],
               [  
                  -118.070609870683398,
                  33.750863839894429
               ],
               [  
                  -118.070430196176659,
                  33.750628817717171
               ],
               [  
                  -118.070105673375835,
                  33.750012360709299
               ],
               [  
                  -118.069990937566118,
                  33.749658126154841
               ],
               [  
                  -118.069907037737352,
                  33.749465219258965
               ],
               [  
                  -118.069896698325621,
                  33.749107649377173
               ]
            ],
			[
				[
					-118.0780506134033,
						33.73140667964194
					],
				[
					-118.06260108947754,
						33.73140667964194
					],
				[
					-118.06260108947754,
						33.73754523293169
					],
				[
					-118.0780506134033,
						33.73754523293169
					],
				[
					-118.0780506134033,
						33.73140667964194
					]
			],
[
				[
				  -118.08002471923827,
				  33.733333947157824
				],
				[
				  -118.0788230895996,
				  33.73190634574705
				],
				[
				  -118.07599067687988,
				  33.73240600894261
				],
				[
				  -118.07564735412598,
				  33.73126391736321
				],
				[
				  -118.06852340698242,
				  33.732548769321205
				],
				[
				  -118.06440353393556,
				  33.73347670599252
				],
				[
				  -118.07101249694826,
				  33.73833036504399
				],
				[
				  -118.0780506134033,
				  33.738401740334254
				],
				[
				  -118.08002471923827,
				  33.733333947157824
				]
			]
         ]
      ]
   }
}`)

		multipolygon := MultiPolygon{}
		So(ParseGeoJSONFeature(seal_beach, func(encoded_properties []byte, polygon MultiPolygon) error {
			multipolygon = polygon
			return nil
		}), ShouldBeNil)

		_, err := multipolygon.ToGeoJSONFeature(nil)
		So(err, ShouldBeNil)
		//fmt.Println(string(f))
		So(multipolygon.Contains(33.735760815044635,-118.06564807891844), ShouldBeFalse) // Within Hole
		So(multipolygon.Contains(33.745252,-118.0801775), ShouldBeTrue)
		So(multipolygon.Contains(33.7387304,-118.0735578), ShouldBeTrue)
		So(multipolygon.Contains(33.8041577,-84.4721115), ShouldBeFalse) // Atlanta
		So(multipolygon.Contains(33.7215987,-118.1046553), ShouldBeFalse) // In the Ocean.
	})
}

