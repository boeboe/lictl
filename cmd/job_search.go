package cmd

import (
	"fmt"
	"net/http"

	"github.com/boeboe/lictl/pkg/linkedin"
	"github.com/spf13/cobra"
)

var regions []string

// jobSearchCmd represents the job search command
var jobSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for LinkedIn jobs based on regions and keywords",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return ValidateFlags(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {

		// Fetching jobs
		jobs, err := linkedin.SearchJobsOnline(regions, keywords, interval, debug)
		if err != nil {
			if httpErr, ok := err.(*linkedin.HTTPError); ok && httpErr.StatusCode == http.StatusTooManyRequests {
				fmt.Println("Warning: You've hit the rate limit (HTTP 429 Too Many Requests). Please avoid making further requests for some time.")
			} else {
				fmt.Println("Error:", err)
			}
			return
		}

		// Writing jobs to output file
		var outErr error
		var filePath string
		format, _ := linkedin.SetFormat(formatString)

		switch format {
		case linkedin.JSON:
			filePath, outErr = writeOutput(linkedin.ConvertToJSON(jobs), outputDir, "jobs", "json")
		case linkedin.CSV:
			filePath, outErr = writeOutput(linkedin.ConvertToCSV(jobs), outputDir, "jobs", "json")
		}

		if outErr != nil {
			fmt.Println("Error writing jobs:", outErr)
			fmt.Println("Falling back to printing jobs:")
			fmt.Printf("Jobs: %+v\n", jobs)
			return
		}

		fmt.Printf("Jobs written to file %s\n", filePath)
	},
}

func init() {
	jobCmd.AddCommand(jobSearchCmd)
	addRequiredKeywordsFlag(jobSearchCmd)
	addRequiredRegionsFlag(jobSearchCmd)
	addIntervalFlag(jobSearchCmd)
}
