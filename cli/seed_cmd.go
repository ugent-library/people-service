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
	seedCmd.Flags().Bool("without-generated-people", false, "do not generate dummy people")
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

		// Generate organisations if they don't exist yet
		orgs := []string{"CA", "DS", "DI", "EB", "FW", "GE", "LA", "LW", "PS", "PP", "RE", "TW", "WE", "GUK"}
		for _, val := range orgs {
			urn := models.NewURN("biblio_id", val)

			orgs, err := repo.GetOrganizationsByIdentifier(ctx, urn)
			if err != nil {
				return err
			}

			if len(orgs) == 0 {
				org := models.NewOrganization()
				org.NameEng = val
				org.AddIdentifier(urn)

				if _, err = repo.SaveOrganization(ctx, org); err != nil {
					return err
				}
			}
		}

		// Generate users from an optional JSON file, if they don't exist yet
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
				urns := person.GetIdentifierByNS("ugent_username")

				// Don't save a person if they don't have a ugent_username identifier
				if len(urns) == 0 {
					continue
				}

				// Force upsert
				people, err := repo.GetPeopleByIdentifier(ctx, urns[0])
				if err != nil {
					return err
				}

				if len(people) > 0 {
					p := people[0]
					person.DateCreated = p.DateCreated
					person.ID = p.ID
				}

				if _, err = repo.SavePerson(ctx, &person); err != nil {
					return err
				}
			}
		}

		// Generate 100 people
		if without, _ := cmd.Flags().GetBool("without-generated-people"); !without {
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
		}

		cmd.Println("Finished seeding dummy data.")

		return nil
	},
}
