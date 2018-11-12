package cmd

import (
	"github.com/heindl/floracast_process/store"
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
)

func init() {
	rootCmd.AddCommand(countCmd)
}

var countCmd = &cobra.Command{
	Use:   "count",
	Short: "Count all documents in the given collection",
	RunE:  Count,
}

func Count(cmd *cobra.Command, args []string) error {

	if len(args) == 0 || args[0] == "" {
		return errors.New("Collection name expected as the first argument")
	}

	ctx := context.Background()

	floraStore, err := store.NewFloraStore(ctx)
	if err != nil {
		return err
	}
	col, err := floraStore.FirestoreCollection(store.FireStoreCollection(args[0]))
	if err != nil {
		return err
	}
	iter := col.Documents(ctx)
	total := 0
	for {
		_, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		total += 1
	}
	fmt.Println(total)
	return nil

}
