package cmd

import (
	"fmt"
	"net/http"

	"github.com/boeboe/lictl/pkg/linkedin"
	"github.com/boeboe/lictl/pkg/utils"
	"github.com/spf13/cobra"
)

// companySearchCmd represents the search command
var companySearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for LinkedIn companies based on keywords",
	Run: func(cmd *cobra.Command, args []string) {
		if err := checkSharedFlags(); err != nil {
			if err := cmd.Help(); err != nil {
				fmt.Printf("Failed to display help: %v\n", err)
			}
			return
		}

		companies, err := linkedin.SearchCompaniesOnline(keywords, interval, debug)
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
			filePath, dumpErr = utils.DumpToJSON(companies, output)
		case utils.CSV:
			filePath, dumpErr = utils.DumpToCSV(companies, output)
		}

		if dumpErr != nil {
			fmt.Println("Error dumping data:", dumpErr)
			fmt.Println("Falling back to printing companies:")
			utils.DumpFallback(companies)
			return
		}

		fmt.Printf("Data dumped to: %s\n", filePath)
	},
}

func init() {
	companyCmd.AddCommand(companySearchCmd)

	// Add flags
	addSharedFlags(companySearchCmd)
}
