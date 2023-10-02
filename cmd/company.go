package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// companyCmd represents the company command
var companyCmd = &cobra.Command{
	Use:   "company",
	Short: "Interact with LinkedIn company functionalities",
	Long:  `The company command provides functionalities related to LinkedIn companies.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			if err := cmd.Help(); err != nil {
				fmt.Printf("Failed to display help: %v\n", err)
			}
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(companyCmd)
}
