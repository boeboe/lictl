package cmd

import (
	"fmt"
	"net/http"

	"github.com/boeboe/lictl/pkg/linkedin"
	"github.com/spf13/cobra"
)

// pulseSearchCmd represents the search command
var pulseSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for LinkedIn pulses based on keywords",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return ValidateFlags(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {

		// Fetching pulses
		pulses, err := linkedin.SearchPulsesOnline(keywords, interval, debug)
		if err != nil {
			if httpErr, ok := err.(*linkedin.HTTPError); ok && httpErr.StatusCode == http.StatusTooManyRequests {
				fmt.Println("Warning: You've hit the rate limit (HTTP 429 Too Many Requests). Please avoid making further requests for some time.")
			} else {
				fmt.Println("Error:", err)
			}
			return
		}

		// Writing pulses to output file
		var outErr error
		var filePath string
		format, _ := linkedin.SetFormat(formatString)

		switch format {
		case linkedin.JSON:
			filePath, outErr = writeOutput(linkedin.ConvertToJSON(pulses), outputDir, "pulses", "json")
		case linkedin.CSV:
			filePath, outErr = writeOutput(linkedin.ConvertToCSV(pulses), outputDir, "pulses", "json")
		}

		if outErr != nil {
			fmt.Println("Error writing pulses:", outErr)
			fmt.Println("Falling back to printing pulses:")
			fmt.Printf("Pulses: %+v\n", pulses)
			return
		}

		fmt.Printf("Pulses written to file %s\n", filePath)
	},
}

func init() {
	pulseCmd.AddCommand(pulseSearchCmd)
	addRequiredKeywordsFlag(pulseSearchCmd)
	addIntervalFlag(pulseSearchCmd)
}
