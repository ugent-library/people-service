package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rebuildAutocompleteOrganizationsCmd = &cobra.Command{
	Use:   "rebuild-autocomplete-organizations",
	Short: "Rebuild autocomplete organization records",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo, err := newRepository()
		if err != nil {
			fmt.Fprintf(os.Stderr, "got error when building repo: %s\n", err)
			return err
		}

		ctx, cancel := context.WithCancel(context.TODO())
		defer cancel()

		return repo.RebuildAutocompleteOrganizations(ctx)
	},
}

var rebuildAutocompletePeopleCmd = &cobra.Command{
	Use:   "rebuild-autocomplete-people",
	Short: "Rebuild autocomplete person records",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo, err := newRepository()
		if err != nil {
			return err
		}

		ctx, cancel := context.WithCancel(context.TODO())
		defer cancel()

		return repo.RebuildAutocompletePeople(ctx)
	},
}

func init() {
	rootCmd.AddCommand(rebuildAutocompleteOrganizationsCmd)
	rootCmd.AddCommand(rebuildAutocompletePeopleCmd)
}
