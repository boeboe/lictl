package cmd

import (
	"fmt"
	"net/http"

	"github.com/boeboe/lictl/pkg/linkedin"
	"github.com/spf13/cobra"
)

// postGetCmd represents the post get command
var postGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get LinkedIn post details",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return ValidateFlags(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {

		// Fetching post details
		post, err := linkedin.GetPostFromUrl(urlString, debug)
		if err != nil {
			if httpErr, ok := err.(*linkedin.HTTPError); ok && httpErr.StatusCode == http.StatusTooManyRequests {
				fmt.Println("Warning: You've hit the rate limit (HTTP 429 Too Many Requests). Please avoid making further requests for some time.")
			} else {
				fmt.Println("Error:", err)
			}
			return
		}

		// Writing post details to output file
		var outErr error
		var filePath string
		format, _ := linkedin.SetFormat(formatString)

		switch format {
		case linkedin.JSON:
			filePath, outErr = writeOutput(post.Json(), outputDir, "post", "json")
		case linkedin.CSV:
			filePath, outErr = writeOutput(post.CsvHeader()+"\n"+post.CsvContent(), outputDir, "post", "json")
		}

		if outErr != nil {
			fmt.Println("Error writing post details:", outErr)
			fmt.Println("Falling back to printing post details:")
			fmt.Printf("Post details: %+v\n", post)
			return
		}

		fmt.Printf("Post details written to file %s\n", filePath)
	},
}

func init() {
	postCmd.AddCommand(postGetCmd)
	addRequiredUrlFlag(postGetCmd)
}
