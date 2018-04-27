package cmd

import (
	"bitbucket.org/heindl/process/store"
	"cloud.google.com/go/firestore"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
)

func init() {
	predictionCmd.AddCommand(countPredictionCmd)
	rootCmd.AddCommand(predictionCmd)
}

var predictionCmd = &cobra.Command{
	Use: "predictions",
}

var countPredictionCmd = &cobra.Command{
	Use:   "count",
	Short: "Count all documents in the given collection",
	RunE:  CountPredictions,
}

func CountPredictions(cmd *cobra.Command, args []string) error {

	ctx := context.Background()

	floraStore, err := store.NewFloraStore(ctx)
	if err != nil {
		return err
	}
	col, err := floraStore.FirestoreCollection(store.CollectionPredictionIndex)
	if err != nil {
		return err
	}

	iter := &firestore.DocumentIterator{}
	if len(args) > 0 {
		iter = col.Where("NameUsageID", "==", args[0]).Documents(ctx)
	} else {
		iter = col.Documents(ctx)
	}

	i := 0
	for {
		_, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		i += 1
	}
	fmt.Println(i)
	return nil

}
