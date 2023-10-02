package cmd

import (
	"fmt"
	"net/http"

	"github.com/boeboe/lictl/pkg/linkedin"
	"github.com/boeboe/lictl/pkg/utils"
	"github.com/spf13/cobra"
)

var regions []string

// jobSearchCmd represents the search command
var jobSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for LinkedIn jobs based on regions and keywords",
	Run: func(cmd *cobra.Command, args []string) {
		if len(regions) == 0 {
			fmt.Println("Error: 'keywords' flag is mandatory")
			if err := cmd.Help(); err != nil {
				fmt.Printf("Failed to display help: %v\n", err)
			}
			return
		}

		if err := checkSharedFlags(); err != nil {
			if err := cmd.Help(); err != nil {
				fmt.Printf("Failed to display help: %v\n", err)
			}
			return
		}

		jobs, err := linkedin.SearchJobsOnline(regions, keywords, interval, debug)
		if err != nil {
			if httpErr, ok := err.(*utils.HTTPError); ok && httpErr.StatusCode == http.StatusTooManyRequests {
				fmt.Println("Warning: You've hit the rate limit (HTTP 429 Too Many Requests). Please avoid making further requests for some time.")
			} else {
				fmt.Println("Error:", err)
			}
			return
		}

		// Dump the jobs to the specified format
		var dumpErr error
		var filePath string

		switch format {
		case utils.JSON:
			filePath, dumpErr = utils.DumpToJSON(jobs, output)
		case utils.CSV:
			filePath, dumpErr = utils.DumpToCSV(jobs, output)
		}

		if dumpErr != nil {
			fmt.Println("Error dumping data:", dumpErr)
			fmt.Println("Falling back to printing jobs:")
			utils.DumpFallback(jobs)

			return
		}

		fmt.Printf("Data dumped to: %s\n", filePath)
	},
}

func init() {
	jobCmd.AddCommand(jobSearchCmd)

	// Add flags
	jobSearchCmd.Flags().StringSliceVarP(&regions, "regions", "r", nil, "Specify one or more regions")
	addSharedFlags(jobSearchCmd)
}
