package mushroom_observer

import (
	"fmt"
	"strings"
	"bitbucket.org/heindl/taxa/occurrences"
	"context"
	"bitbucket.org/heindl/taxa/store"
	"time"
)

func FetchOccurrences(cxt context.Context, targetID store.DataSourceTargetID, since *time.Time) (occurrences.Occurrences, error) {

	// Get last fetched time.
	startDateStr := "20180101"
	endDateStr := "20180119"

	parameters := []string{
		"date=%s-%s",
		"format=json",
		"has_images=true",
		"has_location=true",
		"is_collection_location=true",
		"east=-49.0",
		"north=83.3",
		"west=-178.2",
		"south=6.6",
	}
	url := "http://mushroomobserver.org/api/observations?" + fmt.Sprintf(strings.Join(parameters, "&"), startDateStr, endDateStr)


	fmt.Println(url)

	return nil, nil
}