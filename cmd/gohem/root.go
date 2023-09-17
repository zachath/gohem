package gohem

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "gohem",
	Version: "0.1.0",
	Short:   "gohem - tool for scraping https://www.hemnet.se",
	Long:    `gohem is a tool for scraping either single or multiple properties listed on hemnet.se All output is in JSON.`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Encounterd and error during execution: '%s'", err.Error())
	}
}
