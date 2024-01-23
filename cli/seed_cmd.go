package cli

import (
	"context"
	"encoding/json"
	"io"
	"os"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/spf13/cobra"
	"github.com/ugent-library/people-service/models"
	"github.com/ugent-library/people-service/repository"
)

func init() {
	rootCmd.AddCommand(seedCmd)
	seedCmd.Flags().Bool("force", false, "force seeding the database")
	seedCmd.Flags().String("people-file", "", "json formatted file containing people to import")
}

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the application with dummy data",
	RunE: func(cmd *cobra.Command, args []string) error {
		// setup services
		repo, err := repository.NewRepository(&repository.Config{
			DbUrl:  config.Db.Url,
			AesKey: config.Db.AesKey,
		})
		if err != nil {
			return err
		}

		ctx := context.TODO()

		if force, _ := cmd.Flags().GetBool("force"); !force {
			people, _, err := repo.GetPeople(ctx)
			if err != nil {
				return err
			}

			if len(people) > 0 {
				cmd.Println("Not seeding dummy data because the database is not empty")
				return nil
			}
		}

		// Generate organisations
		orgs := []string{"CA", "DS", "DI", "EB", "FW", "GE", "LA", "LW", "PS", "PP", "RE", "TW", "WE", "GUK"}
		for _, val := range orgs {
			org := models.NewOrganization()
			org.NameEng = val
			org.AddIdentifier(models.NewURN("biblio_id", val))

			if _, err = repo.SaveOrganization(ctx, org); err != nil {
				return err
			}

		}

		// Read users from an optional JSON file
		if file, _ := cmd.Flags().GetString("people-file"); file != "" {
			fh, err := os.Open(file)
			if err != nil {
				return err
			}

			defer fh.Close()

			var people []models.Person
			raw, _ := io.ReadAll(fh)
			err = json.Unmarshal([]byte(raw), &people)
			if err != nil {
				return err
			}

			for _, person := range people {
				if _, err = repo.SavePerson(ctx, &person); err != nil {
					return err
				}
			}
		}

		// Generate 100 people
		for i := 0; i < 100; i++ {
			var person models.Person
			gofakeit.Struct(&person)

			// Hook this person to a random organization
			org := gofakeit.RandomString(orgs)
			urn := models.NewURN("biblio_id", org)

			orgs, _ := repo.GetOrganizationsByIdentifier(ctx, urn)

			if len(orgs) > 0 {
				org := orgs[0]
				newOrgMember := models.NewOrganizationMember(org.ID)
				person.AddOrganizationMember(newOrgMember)
			}

			if _, err = repo.SavePerson(ctx, &person); err != nil {
				return err
			}
		}

		cmd.Println("Finished seeding dummy data.")

		return nil
	},
}
