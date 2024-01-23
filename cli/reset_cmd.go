package cli

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/ugent-library/people-service/models"
	"github.com/ugent-library/people-service/repository"
)

func init() {
	rootCmd.AddCommand(resetCmd)
	resetCmd.Flags().Bool("force", false, "force destructive reset of all data")
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Destructive reset",
	RunE: func(cmd *cobra.Command, args []string) error {
		if force, _ := cmd.Flags().GetBool("force"); !force {
			cmd.Println("The --force flag is required to perform a destructive reset.")
			return nil
		}

		repo, err := repository.NewRepository(&repository.Config{
			DbUrl:  config.Db.Url,
			AesKey: config.Db.AesKey,
		})
		if err != nil {
			return err
		}

		ctx := context.TODO()

		err = repo.EachPerson(ctx, func(p *models.Person) bool {
			repo.DeletePerson(ctx, p.ID)
			return true
		})

		if err != nil {
			return err
		}

		err = repo.EachOrganization(ctx, func(o *models.Organization) bool {
			repo.DeleteOrganization(ctx, o.ID)
			return true
		})

		if err != nil {
			return err
		}

		cmd.Println("Finished destructive reset.")

		return nil
	},
}
