package gohem

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/spf13/cobra"
	"github.com/zachath/gohem/pkg/gohem"
)

var bostaderCmd = &cobra.Command{
	Use:   "bostader",
	Short: "Scrape properties of search result",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		arg := args[0]
		if !gohem.BostaderRegex.MatchString(arg) {
			fmt.Fprintln(os.Stderr, "Provided link does not match a hemnet /bostader link.")
			return
		}

		properties, err := gohem.ScrapeSearch(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
			return
		}

		if toFile {
			jason, err := json.Marshal(properties)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Encountered an error: %s", err.Error())
				return
			}

			url, err := url.Parse(args[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Encountered an error: %s", err.Error())
				return
			}

			err = os.WriteFile(fmt.Sprintf("%s.json", url.RawQuery), jason, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Encountered an error: %s", err.Error())
				return
			}
		} else {
			s, err := json.MarshalIndent(properties, "", "\t")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Encountered an error: %s", err.Error())
				return
			}

			fmt.Print(string(s))
		}
	},
}

func init() {
	bostaderCmd.Flags().BoolVarP(&toFile, "file", "f", false, "Write to file")
	rootCmd.AddCommand(bostaderCmd)
}
