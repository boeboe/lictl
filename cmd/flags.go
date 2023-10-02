// cmd/flags.go

package cmd

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/boeboe/licli/pkg/linkedin"
	"github.com/spf13/cobra"
)

var (
	debug        bool
	formatString string
	interval     time.Duration
	keywords     []string
	outputDir    string
	urlString    string
)

func addPersistentFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable or disable debug mode")
	cmd.Flags().StringVarP(&formatString, "format", "f", "json", "Output format")
	cmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output folder (default is current folder)")
}

func addIntervalFlag(cmd *cobra.Command) {
	cmd.Flags().DurationVarP(&interval, "interval", "i", 100*time.Millisecond, "Interval between web calls")
}

func addRequiredKeywordsFlag(cmd *cobra.Command) {
	cmd.Flags().StringSliceVarP(&keywords, "keywords", "k", nil, "One or more keywords")
	if err := cmd.MarkFlagRequired("keywords"); err != nil {
		log.Fatalf("Error marking keywords flag as required: %v", err)
	}
}

func addRequiredRegionsFlag(cmd *cobra.Command) {
	cmd.Flags().StringSliceVarP(&regions, "regions", "r", nil, "One or more regions")
	if err := cmd.MarkFlagRequired("regions"); err != nil {
		log.Fatalf("Error marking regions flag as required: %v", err)
	}
}

func addRequiredUrlFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&urlString, "url", "u", "", "Url of details page")
	if err := cmd.MarkFlagRequired("url"); err != nil {
		log.Fatalf("Error marking url flag as required: %v", err)
	}
}

func ValidateFormatFlag() error {
	_, err := linkedin.SetFormat(formatString)
	if err != nil {
		return fmt.Errorf("invalid format. Valid formats are: %s, %s", "json", "csv")
	}
	return nil
}

func ValidateUrlFlag() error {
	if urlString == "" {
		return errors.New("url cannot be empty")
	}

	parsedUrl, err := url.ParseRequestURI(urlString)
	if err != nil {
		return errors.New("invalid URL format")
	}

	if !strings.HasPrefix(parsedUrl.Scheme, "http") {
		return errors.New("url must start with http or https")
	}

	return nil
}

func ValidateIntervalFlag() error {
	if interval <= 0 {
		return errors.New("interval should be larger then 0")
	}

	return nil
}

func ValidateFlags(cmd *cobra.Command, args []string) error {
	if err := ValidateFormatFlag(); err != nil {
		return err
	}
	if err := ValidateUrlFlag(); err != nil {
		return err
	}
	if err := ValidateIntervalFlag(); err != nil {
		return err
	}
	return nil
}
