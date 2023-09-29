package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/boeboe/lictl/pkg/linkedin"
	"github.com/boeboe/lictl/pkg/utils"
	"github.com/spf13/cobra"
)

var regions []string
var keywords []string
var count int

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for LinkedIn jobs based on regions and keywords",
	Run: func(cmd *cobra.Command, args []string) {
		if len(regions) == 0 || len(keywords) == 0 {
			fmt.Println("Error: Both 'regions' and 'keywords' flags are mandatory.")
			cmd.Help() // Display the help message
			return
		}

		jobs, err := linkedin.SearchJobsOnline(regions, keywords)
		if err != nil {
			if httpErr, ok := err.(*utils.HTTPError); ok && httpErr.StatusCode == http.StatusTooManyRequests {
				fmt.Println("Warning: You've hit the rate limit (HTTP 429 Too Many Requests). Please avoid making further requests for some time.")
			} else {
				fmt.Println("Error:", err)
			}
			return
		}

		for _, job := range jobs {
			jobJSON, err := json.MarshalIndent(job, "", "  ")
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println(string(jobJSON))
		}
	},
}

func init() {
	jobCmd.AddCommand(searchCmd)

	// Add flags for regions and keywords.
	searchCmd.Flags().StringSliceVarP(&regions, "regions", "r", nil, "Specify one or more regions")
	searchCmd.Flags().StringSliceVarP(&keywords, "keywords", "k", nil, "Specify one or more keywords")
	searchCmd.Flags().IntVarP(&count, "count", "c", 0, "Specify the number of results to fetch (optional)")
}
