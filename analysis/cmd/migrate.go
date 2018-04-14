package cmd

import (
	"bitbucket.org/heindl/process/store"
	"cloud.google.com/go/firestore"
	"context"
	"github.com/spf13/cobra"
	"google.golang.org/api/iterator"
)

func init() {
	rootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "A migration skip that can change often.",
	RunE:  migrate,
}

func migrate(cmd *cobra.Command, args []string) error {

	ctx := context.Background()

	floraStore, err := store.NewFloraStore(ctx)
	if err != nil {
		return err
	}

	ref, err := floraStore.FirestoreCollection(store.CollectionOccurrences)
	if err != nil {
		return err
	}

	iter := ref.Documents(ctx)

	for {
		snap, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		date, err := snap.DataAt("FormattedDate")
		if err != nil {
			return err
		}

		//geopoint, err := snap.DataAt("GeoFeatureSet.GeoPoint")
		//if err != nil {
		//	return err
		//}
		//
		//var lat, lng float64
		//
		//if p, ok := geopoint.(map[string]interface{}); ok {
		//	lat = p["latitude"].(float64)
		//	lng = p["longitude"].(float64)
		//}
		//
		//if p, ok := geopoint.(*latlng.LatLng); ok {
		//	lat = p.Latitude
		//	lng = p.Longitude
		//}
		//
		//if lat == 0 || lng == 0 {
		//	return errors.New("Geopoint [%s] matches neither expected type")
		//}

		//point, err := geo.NewPoint(lat, lng)
		//if err != nil {
		//	return err
		//}
		//
		//cells := point.S2TokenMap()

		if _, err := snap.Ref.Update(ctx, []firestore.Update{{Path: "FormattedMonth", Value: date.(string)[4:6]}}); err != nil {
			return err
		}

	}

	return nil

	//snaps, err := ref.Documents(ctx).GetAll()
	//if err != nil {
	//	return err
	//}
	//for _, snap := range snaps {
	//
	//	//data := map
	//	//if _, err := snap.Ref.Update(ctx, []firestore.Update{{Path: "NameUsageID", Value: "9sYKdRe6OUgzTwabsjjuFiwVU"}}); err != nil {
	//	//	return err
	//	//}
	//	//snap.
	//}
}
