package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/boeboe/lictl/pkg/linkedin"
	"github.com/boeboe/lictl/pkg/utils"
	"github.com/spf13/cobra"
)

// pulseSearchCmd represents the search command
var pulseSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for LinkedIn pulses based on keywords",
	Run: func(cmd *cobra.Command, args []string) {
		if err := checkSharedFlags(); err != nil {
			if err := cmd.Help(); err != nil {
				fmt.Printf("Failed to display help: %v\n", err)
			}
			return
		}

		pulses, err := linkedin.SearchPulsesOnline(keywords, interval, debug)
		if err != nil {
			if httpErr, ok := err.(*utils.HTTPError); ok && httpErr.StatusCode == http.StatusTooManyRequests {
				fmt.Println("Warning: You've hit the rate limit (HTTP 429 Too Many Requests). Please avoid making further requests for some time.")
			} else {
				fmt.Println("Error:", err)
			}
			return
		}

		// Dump the pulses to the specified format
		var dumpErr error
		var filePath string

		switch format {
		case utils.JSON:
			filePath, dumpErr = utils.DumpToJSON(pulses, output)
		case utils.CSV:
			filePath, dumpErr = utils.DumpToCSV(pulses, output)
		}

		if dumpErr != nil {
			fmt.Println("Error dumping data:", dumpErr)
			fmt.Println("Falling back to printing pulses:")

			for _, pulse := range pulses {
				pulseJSON, err := json.MarshalIndent(pulse, "", "  ")
				if err != nil {
					fmt.Println("Error:", err)
					return
				}
				fmt.Println(string(pulseJSON))
			}
			return
		}

		fmt.Printf("Data dumped to: %s\n", filePath)
	},
}

func init() {
	pulseCmd.AddCommand(pulseSearchCmd)

	// Add flags
	addSharedFlags(pulseSearchCmd)
}
