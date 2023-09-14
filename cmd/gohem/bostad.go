package gohem

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zachath/gohem/pkg/gohem"
)

var toFile bool
var bostadCmd = &cobra.Command{
	Use:   "bostad",
	Short: "Scrape a property",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		property, err := gohem.ScrapeProperty(args[0]) //TODO func ensureing this is an actual hemnet link.
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		}

		if toFile {
			jason, err := json.Marshal(property)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Encountered an error: %s", err.Error())
				return
			}
			os.WriteFile(fmt.Sprintf("%s.json", property.Address), jason, 0644)
		} else {
			s, err := json.MarshalIndent(property, "", "\t")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Encountered an error: %s", err.Error())
				return
			}

			fmt.Print(string(s))
		}
	},
}

func init() {
	bostadCmd.Flags().BoolVarP(&toFile, "file", "f", false, "Write to file")
	rootCmd.AddCommand(bostadCmd)
}
