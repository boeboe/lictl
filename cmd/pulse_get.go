package cmd

import (
	"fmt"
	"net/http"

	"github.com/boeboe/lictl/pkg/linkedin"
	"github.com/spf13/cobra"
)

// pulseGetCmd represents the pulse get command
var pulseGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get LinkedIn pulse details",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return ValidateFlags(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {

		// Fetching pulse details
		pulse, err := linkedin.GetPulseFromUrl(urlString, debug)
		if err != nil {
			if httpErr, ok := err.(*linkedin.HTTPError); ok && httpErr.StatusCode == http.StatusTooManyRequests {
				fmt.Println("Warning: You've hit the rate limit (HTTP 429 Too Many Requests). Please avoid making further requests for some time.")
			} else {
				fmt.Println("Error:", err)
			}
			return
		}

		// Writing pulse details to output file
		var outErr error
		var filePath string
		format, _ := linkedin.SetFormat(formatString)

		switch format {
		case linkedin.JSON:
			filePath, outErr = writeOutput(pulse.Json(), outputDir, "pulse", "json")
		case linkedin.CSV:
			filePath, outErr = writeOutput(pulse.CsvHeader()+"\n"+pulse.CsvContent(), outputDir, "pulse", "json")
		}

		if outErr != nil {
			fmt.Println("Error writing pulse details:", outErr)
			fmt.Println("Falling back to printing pulse details:")
			fmt.Printf("Pulse details: %+v\n", pulse)
			return
		}

		fmt.Printf("Pulse details written to file %s\n", filePath)
	},
}

func init() {
	pulseCmd.AddCommand(pulseGetCmd)
	addRequiredUrlFlag(pulseGetCmd)
}
