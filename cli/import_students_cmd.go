package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ugent-library/people-service/models"
	"github.com/ugent-library/people-service/student"
)

var importStudentsCmd = &cobra.Command{
	Use:   "students",
	Short: "import student records from UGent LDAP",
	RunE: func(cmd *cobra.Command, args []string) error {
		ugentLdapClient := newUgentLdapClient()
		repo, err := newRepository()
		if err != nil {
			return err
		}

		importer := student.NewImporter(repo, ugentLdapClient)
		return importer.Each(func(person *models.Person) error {
			isNewRecord := !person.IsStored()
			person, err := repo.SavePerson(context.TODO(), person)
			if err != nil {
				return fmt.Errorf("unable to save person %s: %w", person.ID, err)
			}
			if isNewRecord {
				logger.Infof("successfully created person %s", person.ID)
			} else {
				logger.Infof("successfully updated person %s", person.ID)
			}
			return nil
		})
	},
}

func init() {
	importCmd.AddCommand(importStudentsCmd)
}
