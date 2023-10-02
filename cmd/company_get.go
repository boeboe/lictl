package cmd

import (
	"fmt"
	"net/http"

	"github.com/boeboe/lictl/pkg/linkedin"
	"github.com/spf13/cobra"
)

// companyGetCmd represents the company get command
var companyGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get LinkedIn company details",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return ValidateFlags(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {

		// Fetching company details
		company, err := linkedin.GetCompanyFromUrl(urlString, debug)
		if err != nil {
			if httpErr, ok := err.(*linkedin.HTTPError); ok && httpErr.StatusCode == http.StatusTooManyRequests {
				fmt.Println("Warning: You've hit the rate limit (HTTP 429 Too Many Requests). Please avoid making further requests for some time.")
			} else {
				fmt.Println("Error:", err)
			}
			return
		}

		// Writing company details to output file
		var outErr error
		var filePath string
		format, _ := linkedin.SetFormat(formatString)

		switch format {
		case linkedin.JSON:
			filePath, outErr = writeOutput(company.Json(), outputDir, "company", "json")
		case linkedin.CSV:
			filePath, outErr = writeOutput(company.CsvHeader()+"\n"+company.CsvContent(), outputDir, "company", "json")
		}

		if outErr != nil {
			fmt.Println("Error writing company details:", outErr)
			fmt.Println("Falling back to printing company details:")
			fmt.Printf("Company details: %+v\n", company)
			return
		}

		fmt.Printf("Company details written to file %s\n", filePath)
	},
}

func init() {
	companyCmd.AddCommand(companyGetCmd)
	addRequiredUrlFlag(companyGetCmd)
}
