package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// postCmd represents the user post
var postCmd = &cobra.Command{
	Use:   "user",
	Short: "Interact with LinkedIn user functionalities",
	Long:  `The post command provides functionalities related to LinkedIn posts.`,
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
	rootCmd.AddCommand(postCmd)
}
