package cli

import "github.com/spf13/cobra"

var importCmd = &cobra.Command{
	Use: "import",
}

func init() {
	rootCmd.AddCommand(importCmd)
}
