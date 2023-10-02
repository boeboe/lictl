package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/boeboe/lictl/pkg/linkedin"
	"github.com/boeboe/lictl/pkg/utils"
	"github.com/spf13/cobra"
)

// postSearchCmd represents the search command
var postSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for LinkedIn posts based on keywords",
	Run: func(cmd *cobra.Command, args []string) {
		if err := checkSharedFlags(); err != nil {
			if err := cmd.Help(); err != nil {
				fmt.Printf("Failed to display help: %v\n", err)
			}
			return
		}

		posts, err := linkedin.SearchPostsOnline(keywords, interval, debug)
		if err != nil {
			if httpErr, ok := err.(*utils.HTTPError); ok && httpErr.StatusCode == http.StatusTooManyRequests {
				fmt.Println("Warning: You've hit the rate limit (HTTP 429 Too Many Requests). Please avoid making further requests for some time.")
			} else {
				fmt.Println("Error:", err)
			}
			return
		}

		// Dump the posts to the specified format
		var dumpErr error
		var filePath string

		switch format {
		case utils.JSON:
			filePath, dumpErr = utils.DumpToJSON(posts, output)
		case utils.CSV:
			filePath, dumpErr = utils.DumpToCSV(posts, output)
		}

		if dumpErr != nil {
			fmt.Println("Error dumping data:", dumpErr)
			fmt.Println("Falling back to printing posts:")

			for _, post := range posts {
				postJSON, err := json.MarshalIndent(post, "", "  ")
				if err != nil {
					fmt.Println("Error:", err)
					return
				}
				fmt.Println(string(postJSON))
			}
			return
		}

		fmt.Printf("Data dumped to: %s\n", filePath)
	},
}

func init() {
	postCmd.AddCommand(postSearchCmd)

	// Add flags
	addSharedFlags(postSearchCmd)
}
