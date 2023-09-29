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
		fmt.Println("Use a sub-command with 'job'. For help, use 'lictl job --help'")
	},
}

func init() {
	rootCmd.AddCommand(jobCmd)
}
