package cli

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ugent-library/people-service/gismo"
	"github.com/ugent-library/people-service/models"
)

var orgsCmd = &cobra.Command{
	Use:   "orgs",
	Short: "orgs",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo, err := newRepository()
		if err != nil {
			return err
		}

		err = repo.EachOrganization(context.TODO(), func(org *models.Organization) bool {
			data, err := json.Marshal(org)
			if err != nil {
				return false
			}
			os.Stdout.Write(data)
			os.Stdout.Write([]byte("\n"))

			return true
		})
		return err
	},
}

var procOrgsCmd = &cobra.Command{
	Use:   "procorgs",
	Short: "procorgs",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo, err := newRepository()
		if err != nil {
			return err
		}

		gismo := gismo.NewOrganizationProcessor(repo)

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			file := scanner.Text()
			if strings.HasSuffix(file, ".xml.gz") {
				cmd := exec.Command("gunzip", "-c", file)
				data, err := cmd.Output()
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s : %s\n", file, err)
					continue
				}
				msg, err := gismo.Process(data)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s : %s\n", file, err)
					continue
				}
				msgData, _ := json.Marshal(msg)
				os.Stdout.Write(msgData)
				os.Stdout.Write([]byte("\n"))
			} else {
				data, err := os.ReadFile(file)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s : %s\n", file, err)
					continue
				}
				msg, err := gismo.Process(data)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s : %s\n", file, err)
					continue
				}
				msgData, _ := json.Marshal(msg)
				os.Stdout.Write(msgData)
				os.Stdout.Write([]byte("\n"))
			}
		}
		if err := scanner.Err(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(procOrgsCmd)
	rootCmd.AddCommand(orgsCmd)
}
