package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/zachath/gohem/cmd/gohem"
)

func init() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

// Usage: gohem bostad/bostader <link>
func main() {
	/*url := "https://www.hemnet.se/bostad/lagenhet-1rum-kallhall-jarfalla-kommun-rondellen-14,-1-tr-20347297" //såld
	//url := "https://www.hemnet.se/bostad/villa-8rum-kallhalls-villastad-jarfalla-kommun-nackrosvagen-16-20378175"
	//url := "https://www.hemnet.se/bostad/lagenhet-1rum-kallhall-jarfalla-kommun-fabelvagen-7-19883846"

	scraper := PropertyScraper{url: url}
	property, err := scraper.ScrapeProperty()
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to scrape property")
		os.Exit(1)
	}

	log.Info().Interface("Property", property).Msg("")

	os.Exit(0)*/
	gohem.Execute()
}
