package cmd

import (
	"fmt"
	"net/http"

	"github.com/boeboe/lictl/pkg/linkedin"
	"github.com/spf13/cobra"
)

// companySearchCmd represents the company search command
var companySearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for LinkedIn companies based on keywords",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return ValidateFlags(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {

		// Fetching companies
		companies, err := linkedin.SearchCompaniesOnline(keywords, interval, debug)
		if err != nil {
			if httpErr, ok := err.(*linkedin.HTTPError); ok && httpErr.StatusCode == http.StatusTooManyRequests {
				fmt.Println("Warning: You've hit the rate limit (HTTP 429 Too Many Requests). Please avoid making further requests for some time.")
			} else {
				fmt.Println("Error:", err)
			}
			return
		}

		// Writing companies to output file
		var outErr error
		var filePath string
		format, _ := linkedin.SetFormat(formatString)

		switch format {
		case linkedin.JSON:
			filePath, outErr = writeOutput(linkedin.ConvertToJSON(companies), outputDir, "companies", "json")
		case linkedin.CSV:
			filePath, outErr = writeOutput(linkedin.ConvertToCSV(companies), outputDir, "companies", "json")
		}

		if outErr != nil {
			fmt.Println("Error writing companies:", outErr)
			fmt.Println("Falling back to printing companies:")
			fmt.Printf("Companies: %+v\n", companies)
			return
		}

		fmt.Printf("Companies written to file %s\n", filePath)
	},
}

func init() {
	companyCmd.AddCommand(companySearchCmd)
	addRequiredKeywordsFlag(companySearchCmd)
	addIntervalFlag(companySearchCmd)
}
