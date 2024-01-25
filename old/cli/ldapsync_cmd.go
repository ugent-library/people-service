package cli

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/ugent-library/people-service/ldapsync"
)

var ldapSyncCmd = &cobra.Command{
	Use:   "ldapsync",
	Short: "synchronize person records with UGent LDAP person records",
	RunE: func(cmd *cobra.Command, args []string) error {
		ugentLdapClient := newUgentLdapClient()
		repo, err := newRepository()
		if err != nil {
			return err
		}

		ctx, cancel := context.WithCancel(context.TODO())
		defer cancel()

		importer := ldapsync.NewSynchronizer(repo, ugentLdapClient, logger)
		return importer.Sync(ctx)
	},
}

func init() {
	rootCmd.AddCommand(ldapSyncCmd)
}
