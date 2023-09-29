package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/boeboe/lictl/pkg/linkedin"
	"github.com/boeboe/lictl/pkg/utils"
	"github.com/spf13/cobra"
)

var debug bool
var formatString string
var interval time.Duration
var keywords []string
var output string
var regions []string

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for LinkedIn jobs based on regions and keywords",
	Run: func(cmd *cobra.Command, args []string) {
		if len(regions) == 0 || len(keywords) == 0 {
			fmt.Println("Error: Both 'regions' and 'keywords' flags are mandatory.")
			if err := cmd.Help(); err != nil {
				fmt.Printf("Failed to display help: %v\n", err)
			}
			return
		}

		format, err := utils.SetFormat(formatString)
		if err != nil {
			fmt.Println("Error parsing format:", err)
			return
		}

		if output == "" {
			output, err = os.Getwd()
			if err != nil {
				fmt.Println("Error getting current working directory:", err)
				return
			}
		}

		if _, err := os.Stat(output); os.IsNotExist(err) {
			err = os.MkdirAll(output, 0755)
			if err != nil {
				fmt.Println("Error creating directory:", err)
				return
			}
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

			for _, job := range jobs {
				jobJSON, err := json.MarshalIndent(job, "", "  ")
				if err != nil {
					fmt.Println("Error:", err)
					return
				}
				fmt.Println(string(jobJSON))
			}
			return
		}

		fmt.Printf("Data dumped to: %s\n", filePath)
	},
}

func init() {
	jobCmd.AddCommand(searchCmd)

	// Add flags
	searchCmd.Flags().StringSliceVarP(&regions, "regions", "r", nil, "Specify one or more regions")
	searchCmd.Flags().StringSliceVarP(&keywords, "keywords", "k", nil, "Specify one or more keywords")
	searchCmd.Flags().StringVarP(&output, "output", "o", "", "Specify the output (default is current working directory)")
	searchCmd.Flags().StringVarP(&formatString, "format", "f", "json", "Specify the format")
	searchCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable or disable debug mode")
	searchCmd.Flags().DurationVarP(&interval, "interval", "i", 100*time.Millisecond, "Specify the interval between web calls")
}
