package cli

import (
	"context"

	"github.com/spf13/cobra"
)

var expirePersonCmd = &cobra.Command{
	Use:   "expire-person",
	Short: "auto expire person records",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo, err := newRepository()
		if err != nil {
			return err
		}
		nAffected, err := repo.AutoExpirePeople(context.TODO())
		if err != nil {
			return err
		}
		logger.Infof("%d person records expired", nAffected)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(expirePersonCmd)
}
