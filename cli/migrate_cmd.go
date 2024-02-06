package cli

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/ugent-library/people-service/db"
)

func init() {
	rootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:       "migrate [up|down]",
	Short:     "Run database migrations",
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"up", "down"},
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "up":
			return db.MigrateUp(context.Background(), config.Repo.Conn)
		case "down":
			return db.MigrateDown(context.Background(), config.Repo.Conn)
		}
		return nil
	},
}
