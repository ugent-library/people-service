package cli

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/ugent-library/people-service/db"
)

func init() {
	migrateCmd.Flags().Int32("version", -1, "version number to migrate up or down to")
	rootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		versionFlag, err := cmd.Flags().GetInt32("version")
		if err != nil {
			return err
		}
		if versionFlag >= 0 {
			return db.MigrateTo(ctx, config.Repo.Conn, versionFlag)
		}
		return db.Migrate(ctx, config.Repo.Conn)
	},
}
