package main

import (
	"fmt"
	"github.com/cockroachdb/apd"
	"github.com/dropbox/godropbox/errors"
	"math"
	"strconv"
)

func main() {
	uncertainty, err := UncertaintyAssociatedWithCoordinatePrecision("10.27", "-123.6")
	if err != nil {
		panic(err)
	}
	fmt.Println(uncertainty)
}

//const EARTH_EQUATORIAL_CIRCUMFERENCE = 40070000
//const EARTH_MERIDONIAL_CIRCUMFERENCE = 39931000

//const LONGITUDAL_DEGREE = EARTH_EQUATORIAL_CIRCUMFERENCE / 360
//const LATITUDAL_DEGREE = EARTH_MERIDONIAL_CIRCUMFERENCE / 360

//const EQUATORIAL_RADIUS = 6378137
const SEMI_MAJOR_AXIS = 6378137 // Equatorial Radius

//const POLAR_RADIUS = 6356752.3141
//const SEMI_MINOR_AXIS = POLAR_RADIUS

//var FLATTENING = (SEMI_MAJOR_AXIS - SEMI_MINOR_AXIS) / SEMI_MAJOR_AXIS

var FLATTENING = 1 / 298.25722356

var FIRST_ECCENTRICITY = (2 * FLATTENING) - math.Pow(FLATTENING, 2)

//( math.Pow(FLATTENING, 2) - b2 ) / a2
//
//var FIRST_ECCENTRICITY = 0.00669438002290

func CoordinatePrecisionAsFractionOfOneDegree(lat, lng string) (float64, float64, error) {
	latDecimal, _, err := apd.NewFromString(lat)
	if err != nil {
		return 0, 0, errors.Wrap(err, "Could not read lat from string")
	}
	//lngDecimal, _, err := apd.NewFromString(lng)
	//if err != nil {
	//	return 0, 0, errors.Wrap(err, "Could not read lng from string")
	//}
	latPrecision := math.Pow10(int(latDecimal.Exponent))
	//lngPrecision := math.Pow10(int(lngDecimal.Exponent))
	return latPrecision, latPrecision, nil
}


func Radius(B float64) float64 {

// https://rechneronline.de/earth-radius/
	[ (r1² * math.Cos(B))² + (r2² * math.Sin(B))² ] / [ (r1 * cos(B))² + (r2 * sin(B))² ]f
}


func RadiusOfCurvatureOfTheMeridian(latitude float64) float64 {
	// R = a(1-e2)/(1-e2sin2(latitude))3/2

	sinSquared := math.Pow(math.Sin(latitude), 2)

	numerator := (1 - FIRST_ECCENTRICITY)
	denominator := math.Sqrt(math.Pow(1-(FIRST_ECCENTRICITY*sinSquared), 3))
	return SEMI_MAJOR_AXIS * numerator / denominator
}

func DistanceToThePolarAxis(latitude float64) float64 {
	// X = abs(Ncos(latitude))
	var N = RadiusOfCurvatureInThePrimeVertical(latitude)
	return math.Abs(N * math.Cos(latitude))
}

func RadiusOfCurvatureInThePrimeVertical(latitude float64) float64 {
	// N = a/sqrt(1-e2sin2(latitude))
	sinSquared := math.Pow(math.Sin(latitude), 2)
	numerator := float64(SEMI_MAJOR_AXIS)
	denominator := math.Sqrt(1 - FIRST_ECCENTRICITY*sinSquared)
	return numerator / denominator
}

// UncertaintyAssociatedWithCoordinatePrecision generates an uncertainty in meters for a coordinate.
// https://www.tandfonline.com/doi/pdf/10.1080/13658810412331280211
func UncertaintyAssociatedWithCoordinatePrecision(latitude, longitude string) (float64, error) {

	latPrecision, lngPrecision, err := CoordinatePrecisionAsFractionOfOneDegree(latitude, longitude)
	if err != nil {
		return 0, err
	}

	fmt.Println(latPrecision, lngPrecision)

	latFloat, err := strconv.ParseFloat(latitude, 64)
	if err != nil {
		return 0, errors.Wrap(err, "Could not parse float")
	}

	lngFloat, err := strconv.ParseFloat(longitude, 64)
	if err != nil {
		return 0, errors.Wrap(err, "Could not parse float")
	}

	fmt.Println("RadiusOfCurvatureOfTheMeridian", latFloat, RadiusOfCurvatureOfTheMeridian(latFloat))

	latError := (math.Pi * RadiusOfCurvatureOfTheMeridian(latFloat) * latPrecision) / 180

	fmt.Println("latError", fmt.Sprintf("%.4f", latError))

	lngError := (math.Pi * DistanceToThePolarAxis(lngFloat) * latPrecision) / 180

	fmt.Println("DistanceToThePolarAxisOrthogonalToThePolarAxis", latFloat, DistanceToThePolarAxis(latFloat))

	fmt.Println("lngError", fmt.Sprintf("%.4f", lngError))

	return math.Sqrt(math.Pow(latError, 2) + math.Pow(lngError, 2)), nil

}

//
//
//// http://mycoportal.org/portal/collections/listtabledisplay.php?eventdate1=2001-01-01&eventdate2=2002-06-17&llbound=59.5747563;24.5465169;-145.1767463;-49&db=all&occindex=3&sortfield1=Number&sortfield2=&sortorder=asc
//
//func url(page int, request *providers.OccurrenceFetchRequest) {
//	url := "http://mycoportal.org/portal/collections/listtabledisplay.php?"
//	if request.StartTime != nil {
//		url += fmt.Sprintf("eventdate1=", request.StartTime.Format("2006-01-2"))
//	}
//	if request.EndTime != nil {
//		url += fmt.Sprintf("eventdate2=", request.EndTime.Format("2006-01-2"))
//	}
//
//	eventdate1=2001-01-01&eventdate2=2002-06-17&llbound=59.5747563;24.5465169;-145.1767463;-49&db=all&occindex=3&sortfield1=Number&sortfield2=&sortorder=asc'
//}
//
//// http://mycoportal.org/portal/collections/listtabledisplay.php?eventdate1=2001-01-01&eventdate2=2002-06-17&llbound=83.3;6.6;-178.2;-49
//func FetchOccurrences() {
//
//
//	// Create another collector to scrape course details
//	detailCollector := c.Clone()
//
//	courses := make([]Course, 0, 200)
//	// Next 1000 records
//	// On every a element which has href attribute call callback
//	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
//		// If attribute class is this long string return from callback
//		// As this a is irrelevant
//		if e.Attr("class") == "Button_1qxkboh-o_O-primary_cv02ee-o_O-md_28awn8-o_O-primaryLink_109aggg" {
//			return
//		}
//		link := e.Attr("href")
//		// If link start with browse or includes either signup or login return from callback
//		if !strings.HasPrefix(link, "/browse") || strings.Index(link, "=signup") > -1 || strings.Index(link, "=login") > -1 {
//			return
//		}
//		// start scaping the page under the link found
//		e.Request.Visit(link)
//	})
//
//	// Before making a request print "Visiting ..."
//	c.OnRequest(func(r *colly.Request) {
//		log.Println("visiting", r.URL.String())
//	})
//
//	// On every a HTML element which has name attribute call callback
//	c.OnHTML(`a[name]`, func(e *colly.HTMLElement) {
//		// Activate detailCollector if the link contains "coursera.org/learn"
//		courseURL := e.Request.AbsoluteURL(e.Attr("href"))
//		if strings.Index(courseURL, "coursera.org/learn") != -1 {
//			detailCollector.Visit(courseURL)
//		}
//	})
//
//	// Extract details of the course
//	detailCollector.OnHTML(`div[id=rendered-content]`, func(e *colly.HTMLElement) {
//		log.Println("Course found", e.Request.URL)
//		title := e.ChildText(".course-title")
//		if title == "" {
//			log.Println("No title found", e.Request.URL)
//		}
//		course := Course{
//			Title:       title,
//			URL:         e.Request.URL.String(),
//			Description: e.ChildText("div.content"),
//			Creator:     e.ChildText("div.creator-names > span"),
//		}
//		// Iterate over rows of the table which contains different information
//		// about the course
//		e.ForEach("table.basic-info-table tr", func(_ int, el *colly.HTMLElement) {
//			switch el.ChildText("td:first-child") {
//			case "Language":
//				course.Language = el.ChildText("td:nth-child(2)")
//			case "Level":
//				course.Level = el.ChildText("td:nth-child(2)")
//			case "Commitment":
//				course.Commitment = el.ChildText("td:nth-child(2)")
//			case "How To Pass":
//				course.HowToPass = el.ChildText("td:nth-child(2)")
//			case "User Ratings":
//				course.Rating = el.ChildText("td:nth-child(2) div:nth-of-type(2)")
//			}
//		})
//		courses = append(courses, course)
//	})
//
//	// Start scraping on http://coursera.com/browse
//	c.Visit("https://coursera.org/browse")
//
//	enc := json.NewEncoder(os.Stdout)
//	enc.SetIndent("", "  ")
//
//	// Dump json to the standard output
//	enc.Encode(courses)
//}
