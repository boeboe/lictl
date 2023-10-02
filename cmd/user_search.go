package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/boeboe/lictl/pkg/linkedin"
	"github.com/boeboe/lictl/pkg/utils"
	"github.com/spf13/cobra"
)

// userSearchCmd represents the search command
var userSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for LinkedIn users based on keywords",
	Run: func(cmd *cobra.Command, args []string) {
		if err := checkSharedFlags(); err != nil {
			if err := cmd.Help(); err != nil {
				fmt.Printf("Failed to display help: %v\n", err)
			}
			return
		}

		users, err := linkedin.SearchUsersOnline(keywords, interval, debug)
		if err != nil {
			if httpErr, ok := err.(*utils.HTTPError); ok && httpErr.StatusCode == http.StatusTooManyRequests {
				fmt.Println("Warning: You've hit the rate limit (HTTP 429 Too Many Requests). Please avoid making further requests for some time.")
			} else {
				fmt.Println("Error:", err)
			}
			return
		}

		// Dump the users to the specified format
		var dumpErr error
		var filePath string

		switch format {
		case utils.JSON:
			filePath, dumpErr = utils.DumpToJSON(users, output)
		case utils.CSV:
			filePath, dumpErr = utils.DumpToCSV(users, output)
		}

		if dumpErr != nil {
			fmt.Println("Error dumping data:", dumpErr)
			fmt.Println("Falling back to printing users:")

			for _, user := range users {
				userJSON, err := json.MarshalIndent(user, "", "  ")
				if err != nil {
					fmt.Println("Error:", err)
					return
				}
				fmt.Println(string(userJSON))
			}
			return
		}

		fmt.Printf("Data dumped to: %s\n", filePath)
	},
}

func init() {
	userCmd.AddCommand(userSearchCmd)

	// Add flags
	addSharedFlags(userSearchCmd)
}
