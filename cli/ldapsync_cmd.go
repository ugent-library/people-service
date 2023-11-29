package cli

import (
	"github.com/spf13/cobra"
	"github.com/ugent-library/people-service/ldapsync"
	"github.com/ugent-library/people-service/models"
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

		importer := ldapsync.NewSynchronizer(repo, ugentLdapClient)
		return importer.Sync(func(person *models.Person) {
			logger.Infof("saved person %s to database", person.ID)
		})
	},
}

func init() {
	rootCmd.AddCommand(ldapSyncCmd)
}
