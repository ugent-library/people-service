package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/go-ldap/ldap/v3"
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

		importer := ldapsync.NewSynchronizer(repo, ugentLdapClient)
		return importer.Sync(func(msg string) {
			logger.Infof(msg)
		})
	},
}

var ldapTestCmd = &cobra.Command{
	Use:   "ldaptest",
	Short: "ldaptest",
	RunE: func(cmd *cobra.Command, args []string) error {
		i := 0
		return newUgentLdapClient().SearchPeople(ldapsync.PersonQuery, func(e *ldap.Entry) error {
			i++
			fmt.Fprintf(os.Stderr, "LDAP ENTRY: %d\n", i)
			for _, attr := range e.Attributes {
				for _, val := range attr.Values {
					fmt.Fprintf(os.Stderr, "  %s : %s\n", attr.Name, val)
				}
			}
			if i%100 == 0 {
				time.Sleep(1 * time.Second)
			}
			return nil
		})
	},
}

func init() {
	rootCmd.AddCommand(ldapSyncCmd)
	rootCmd.AddCommand(ldapTestCmd)
}
