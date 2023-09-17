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
		arg := args[0]
		if !gohem.BostadRegex.MatchString(arg) {
			fmt.Fprintln(os.Stderr, "Provided link does not match a hemnet /bostad link.")
			return
		}

		property, err := gohem.ScrapeProperty(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		}

		if toFile {
			jason, err := json.Marshal(property)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Encountered an error: %s", err.Error())
				return
			}
			err = os.WriteFile(fmt.Sprintf("%s.json", property.Address), jason, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Encountered an error: %s", err.Error())
				return
			}
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
