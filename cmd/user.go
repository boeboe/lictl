package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Interact with LinkedIn user functionalities",
	Long:  `The user command provides functionalities related to LinkedIn users.`,
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
	rootCmd.AddCommand(userCmd)
}
