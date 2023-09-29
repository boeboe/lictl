package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// jobCmd represents the job command
var jobCmd = &cobra.Command{
	Use:   "job",
	Short: "Interact with LinkedIn job functionalities",
	Long:  `The job command provides functionalities related to LinkedIn jobs.`,
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
	rootCmd.AddCommand(jobCmd)
}
