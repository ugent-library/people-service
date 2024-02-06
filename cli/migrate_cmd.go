package cli

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/ugent-library/people-service/db"
)

var (
	migrateUpFlag   bool
	migrateDownFlag bool
)

func init() {
	migrateCmd.Flags().BoolVar(&migrateUpFlag, "up", false, "migrate up")
	migrateCmd.Flags().BoolVar(&migrateDownFlag, "down", false, "migrate down")
	migrateCmd.MarkFlagsOneRequired("up", "down")
	migrateCmd.MarkFlagsMutuallyExclusive("up", "down")
	rootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if migrateUpFlag {
			return db.MigrateUp(context.Background(), config.Repo.Conn)
		}
		return db.MigrateDown(context.Background(), config.Repo.Conn)
	},
}
