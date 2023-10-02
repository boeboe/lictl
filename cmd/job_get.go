package cmd

import (
	"fmt"

	"github.com/boeboe/lictl/pkg/linkedin"
	"github.com/spf13/cobra"
)

// jobGetCmd represents the job get command
var jobGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get LinkedIn job details",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return ValidateFlags(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {

		// Fetching job details
		// job, err := linkedin.GetJobFromUrl(urlString, debug)
		job := linkedin.Job{}
		// if err != nil {
		// 	if httpErr, ok := err.(*linkedin.HTTPError); ok && httpErr.StatusCode == http.StatusTooManyRequests {
		// 		fmt.Println("Warning: You've hit the rate limit (HTTP 429 Too Many Requests). Please avoid making further requests for some time.")
		// 	} else {
		// 		fmt.Println("Error:", err)
		// 	}
		// 	return
		// }

		// Writing job details to output file
		var outErr error
		var filePath string
		format, _ := linkedin.SetFormat(formatString)

		switch format {
		case linkedin.JSON:
			filePath, outErr = writeOutput(job.Json(), outputDir, "job", "json")
		case linkedin.CSV:
			filePath, outErr = writeOutput(job.CsvHeader()+"\n"+job.CsvContent(), outputDir, "job", "json")
		}

		if outErr != nil {
			fmt.Println("Error writing job details:", outErr)
			fmt.Println("Falling back to printing job details:")
			fmt.Printf("Job details: %+v\n", job)
			return
		}

		fmt.Printf("Job details written to file %s\n", filePath)
	},
}

func init() {
	jobCmd.AddCommand(jobGetCmd)
	addRequiredUrlFlag(jobGetCmd)
}
