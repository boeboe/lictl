// cmd/flags.go

package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/boeboe/lictl/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	debug        bool
	format       utils.FormatType
	formatString string
	interval     time.Duration
	keywords     []string
	output       string
)

func addSharedFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable or disable debug mode")
	cmd.Flags().StringVarP(&formatString, "format", "f", "json", "Specify the format")
	cmd.Flags().DurationVarP(&interval, "interval", "i", 100*time.Millisecond, "Specify the interval between web calls")
	cmd.Flags().StringSliceVarP(&keywords, "keywords", "k", nil, "Specify one or more keywords")
	cmd.Flags().StringVarP(&output, "output", "o", "", "Specify the output (default is current working directory)")
}

func checkSharedFlags() error {
	var err error
	if len(keywords) == 0 {
		return fmt.Errorf("error: 'regions' flag is mandatory")
	}
	format, err = utils.SetFormat(formatString)
	if err != nil {
		return fmt.Errorf("error parsing format: %w", err)
	}
	if output == "" {
		output, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("error getting current working directory: %w", err)
		}
	}
	if _, err := os.Stat(output); os.IsNotExist(err) {
		err = os.MkdirAll(output, 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return fmt.Errorf("error creating directory: %w", err)
		}
	}
	return nil
}
