package cli

import (
	"context"
	"fmt"
	"os"

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

		ctx, cancel := context.WithCancel(context.TODO())
		defer cancel()

		importer := ldapsync.NewSynchronizer(repo, ugentLdapClient)
		return importer.Sync(ctx, func(msg string) {
			logger.Infof(msg)
		})
	},
}

var ldapTestCmd = &cobra.Command{
	Use:   "ldaptest",
	Short: "ldaptest",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.TODO())
		defer cancel()

		i := 0
		err := newUgentLdapClient().SearchPeople(ctx, ldapsync.PersonQuery, func(e *ldap.Entry) error {
			i++
			fmt.Fprintf(os.Stderr, "LDAP ENTRY: %d\n", i)
			for _, attr := range e.Attributes {
				for _, val := range attr.Values {
					fmt.Fprintf(os.Stderr, "  %s : %s\n", attr.Name, val)
				}
			}
			return nil
		})
		fmt.Fprintf(os.Stderr, "error: %+v\n", err)
		return err
	},
}

func init() {
	rootCmd.AddCommand(ldapSyncCmd)
	rootCmd.AddCommand(ldapTestCmd)
}
