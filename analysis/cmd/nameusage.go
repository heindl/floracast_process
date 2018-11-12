package cmd

import (
	"github.com/heindl/floracast_process/store"
	"github.com/heindl/floracast_process/utils"
	"context"
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"github.com/spf13/cobra"
)

func init() {
	nameUsageCmd.AddCommand(nameUsageJSONCmd)
	rootCmd.AddCommand(nameUsageCmd)
}

var nameUsageCmd = &cobra.Command{
	Use:   "nameusage",
	Short: "Handle operations on the NameUsage collection",
}

var nameUsageJSONCmd = &cobra.Command{
	Use:   "json",
	Short: "Read NameUsage into JSON",
	RunE:  NameUsageJSON,
}

func NameUsageJSON(cmd *cobra.Command, args []string) error {

	if len(args) == 0 || args[0] == "" {
		return errors.New("NameUsageID expected")
	}

	ctx := context.Background()

	floraStore, err := store.NewFloraStore(ctx)
	if err != nil {
		return err
	}
	col, err := floraStore.FirestoreCollection(store.CollectionNameUsages)
	if err != nil {
		return err
	}
	snap, err := col.Doc(args[0]).Get(ctx)
	if err != nil {
		return errors.Wrapf(err, "Could not get NameUsageID [%s]", args[0])
	}

	fmt.Println(utils.JsonOrSpew(snap.Data()))
	return nil

}
