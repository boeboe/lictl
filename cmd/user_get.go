package cmd

import (
	"fmt"
	"net/http"

	"github.com/boeboe/lictl/pkg/linkedin"
	"github.com/spf13/cobra"
)

// userGetCmd represents the user get command
var userGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get LinkedIn user details",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return ValidateFlags(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {

		// Fetching user details
		user, err := linkedin.GetUserFromUrl(urlString, debug)
		if err != nil {
			if httpErr, ok := err.(*linkedin.HTTPError); ok && httpErr.StatusCode == http.StatusTooManyRequests {
				fmt.Println("Warning: You've hit the rate limit (HTTP 429 Too Many Requests). Please avoid making further requests for some time.")
			} else {
				fmt.Println("Error:", err)
			}
			return
		}

		// Writing user details to output file
		var outErr error
		var filePath string
		format, _ := linkedin.SetFormat(formatString)

		switch format {
		case linkedin.JSON:
			filePath, outErr = writeOutput(user.Json(), outputDir, "user", "json")
		case linkedin.CSV:
			filePath, outErr = writeOutput(user.CsvHeader()+"\n"+user.CsvContent(), outputDir, "user", "json")
		}

		if outErr != nil {
			fmt.Println("Error writing user details:", outErr)
			fmt.Println("Falling back to printing user details:")
			fmt.Printf("User details: %+v\n", user)
			return
		}

		fmt.Printf("User details written to file %s\n", filePath)
	},
}

func init() {
	userCmd.AddCommand(userGetCmd)
	addRequiredUrlFlag(userGetCmd)
}
