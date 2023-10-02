package cmd

import (
	"fmt"
	"net/http"

	"github.com/boeboe/lictl/pkg/linkedin"
	"github.com/spf13/cobra"
)

// userSearchCmd represents the search command
var userSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for LinkedIn users based on keywords",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return ValidateFlags(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Fetching users
		users, err := linkedin.SearchUsersOnline(keywords, interval, debug)
		if err != nil {
			if httpErr, ok := err.(*linkedin.HTTPError); ok && httpErr.StatusCode == http.StatusTooManyRequests {
				fmt.Println("Warning: You've hit the rate limit (HTTP 429 Too Many Requests). Please avoid making further requests for some time.")
			} else {
				fmt.Println("Error:", err)
			}
			return
		}

		// Writing users to output file
		var outErr error
		var filePath string
		format, _ := linkedin.SetFormat(formatString)

		switch format {
		case linkedin.JSON:
			filePath, outErr = writeOutput(linkedin.ConvertToJSON(users), outputDir, "users", "json")
		case linkedin.CSV:
			filePath, outErr = writeOutput(linkedin.ConvertToCSV(users), outputDir, "users", "json")
		}

		if outErr != nil {
			fmt.Println("Error writing users:", outErr)
			fmt.Println("Falling back to printing users:")
			fmt.Printf("Users: %+v\n", users)
			return
		}

		fmt.Printf("Users written to file %s\n", filePath)
	},
}

func init() {
	userCmd.AddCommand(userSearchCmd)
	addRequiredKeywordsFlag(userSearchCmd)
	addIntervalFlag(userSearchCmd)
}
