package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// pulseCmd represents the pulse command
var pulseCmd = &cobra.Command{
	Use:   "pulse",
	Short: "Interact with LinkedIn pulse functionalities",
	Long:  `The pulse command provides functionalities related to LinkedIn pulses.`,
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
	rootCmd.AddCommand(pulseCmd)
}
