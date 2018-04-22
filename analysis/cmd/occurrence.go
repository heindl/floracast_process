package cmd

import (
	"bitbucket.org/heindl/process/store"
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
)

func init() {
	occurrenceCmd.AddCommand(countOccurrenceCmd)
	rootCmd.AddCommand(occurrenceCmd)
}

var occurrenceCmd = &cobra.Command{
	Use: "occurrences",
}

var countOccurrenceCmd = &cobra.Command{
	Use:   "count",
	Short: "Count all documents in the given collection",
	RunE:  CountOccurrences,
}

func CountOccurrences(cmd *cobra.Command, args []string) error {

	ctx := context.Background()

	floraStore, err := store.NewFloraStore(ctx)
	if err != nil {
		return err
	}
	col, err := floraStore.FirestoreCollection(store.CollectionOccurrences)
	if err != nil {
		return err
	}
	iter := col.Documents(ctx)
	m := map[string]int{}
	for {
		snap, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		usageID, err := snap.DataAt("NameUsageID")
		if err != nil {
			return errors.Wrap(err, "Could not get NameUsageID")
		}
		if _, ok := m[usageID.(string)]; !ok {
			m[usageID.(string)] = 0
		}
		m[usageID.(string)] += 1
	}
	for usageID, count := range m {
		fmt.Println(fmt.Sprintf("%s: %d", usageID, count))
	}
	return nil

}
