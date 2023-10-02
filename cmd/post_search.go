package cmd

import (
	"fmt"
	"net/http"

	"github.com/boeboe/lictl/pkg/linkedin"
	"github.com/spf13/cobra"
)

// postSearchCmd represents the search command
var postSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for LinkedIn posts based on keywords",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return ValidateFlags(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {

		// Fetching posts
		posts, err := linkedin.SearchPostsOnline(keywords, interval, debug)
		if err != nil {
			if httpErr, ok := err.(*linkedin.HTTPError); ok && httpErr.StatusCode == http.StatusTooManyRequests {
				fmt.Println("Warning: You've hit the rate limit (HTTP 429 Too Many Requests). Please avoid making further requests for some time.")
			} else {
				fmt.Println("Error:", err)
			}
			return
		}

		// Writing posts to output file
		var outErr error
		var filePath string
		format, _ := linkedin.SetFormat(formatString)

		switch format {
		case linkedin.JSON:
			filePath, outErr = writeOutput(linkedin.ConvertToJSON(posts), outputDir, "posts", "json")
		case linkedin.CSV:
			filePath, outErr = writeOutput(linkedin.ConvertToCSV(posts), outputDir, "posts", "json")
		}

		if outErr != nil {
			fmt.Println("Error writing posts:", outErr)
			fmt.Println("Falling back to printing posts:")
			fmt.Printf("Posts: %+v\n", posts)
			return
		}

		fmt.Printf("Posts written to file %s\n", filePath)
	},
}

func init() {
	postCmd.AddCommand(postSearchCmd)
	addRequiredKeywordsFlag(postSearchCmd)
	addIntervalFlag(postSearchCmd)
}
