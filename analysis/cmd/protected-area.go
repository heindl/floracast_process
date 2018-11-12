package cmd

import (
	"github.com/heindl/floracast_process/store"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
)

func init() {
	rootCmd.AddCommand(protectedAreaCmd)
}

var protectedAreaCmd = &cobra.Command{
	Use:   "list-protected-area-kilometers",
	Short: "List the square kilomters for all protected areas",
	RunE:  ListProtectedAreaSquareKilometers,
}

func ListProtectedAreaSquareKilometers(cmd *cobra.Command, args []string) error {

	ctx := context.Background()

	floraStore, err := store.NewFloraStore(ctx)
	if err != nil {
		return err
	}
	col, err := floraStore.FirestoreCollection(store.CollectionProtectedAreas)
	if err != nil {
		return err
	}
	iter := col.Documents(ctx)
	for {
		snap, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		data, err := snap.DataAt("SquareKilometers")
		if err != nil {
			return err
		}
		fmt.Println(fmt.Sprintf(`"%s": %f,`, snap.Ref.ID, data))
	}
	return nil

}
