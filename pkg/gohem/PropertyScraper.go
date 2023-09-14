package gohem

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func ScrapeProperty(url string) (Property, error) {
	isRemoved, err := isPropertyRemoved(url)
	if err != nil {
		return Property{}, err
	}

	property := Property{}
	if isRemoved {
		property, err = scrapeRemoved(url)
		if err != nil {
			return Property{}, err
		}
	} else {
		property, err = scrape(url)
		if err != nil {
			return Property{}, err
		}
	}

	return property, nil
}

func isPropertyRemoved(url string) (bool, error) {
	collector := colly.NewCollector()

	var removed bool
	collector.OnHTML("div", func(h *colly.HTMLElement) {
		if h.Attr("class") == "removed-listing qa-removed-listing" {
			removed = true
			collector.OnHTMLDetach("div")
		}
	})

	err := collector.Visit(url)
	if err != nil {
		return false, errors.Wrapf(err, "faild to visit url %s", url)
	}

	return removed, nil
}

func scrape(url string) (Property, error) {
	property := Property{}
	collector := colly.NewCollector()

	collector.OnHTML("div", func(e *colly.HTMLElement) {
		if e.Attr("class") == "property-attributes-table__row" || e.Attr("class") == "property-attributes-table__row qa-living-area-attribute" {
			_, err := property.AddField(e.ChildTexts("*")...)
			if err != nil {
				log.Error().Err(err).Msg("error while scraping")
				return
			}
		}

		if e.Attr("class") == "property-visits-counter__row-value" {
			i, err := strconv.ParseInt(strings.Join(strings.Fields(e.Text), ""), 10, 0)
			if err != nil {
				log.Error().Err(err).Msg("error while scraping")
				return
			}

			property.Visits = int(i)
		}

		if e.Attr("class") == "property-info__price-container" {
			price, err := convertPrice(e.Text)
			if err != nil {
				log.Error().Err(err).Msg("error while scraping")
				return
			}

			property.Price = price
		}
	})

	collector.OnHTML("p", realtorFunc(&property))

	collector.OnHTML("a", agencyFunc(&property))

	collector.OnHTML("h1", addressFunc(&property))

	err := collector.Visit(url)
	if err != nil {
		return Property{}, errors.Wrapf(err, "error when visiting url: %s", url)
	}

	return property, nil
}

func scrapeRemoved(url string) (Property, error) {
	collector := colly.NewCollector()
	property := Property{Removed: true}

	collector.OnHTML("span", func(h *colly.HTMLElement) {
		if h.Attr("class") == "removed-listing__heading" {
			property.RemovedDate = h.ChildTexts("*")[0]
		}
	})

	collector.OnHTML("div", func(e *colly.HTMLElement) {
		if e.Attr("class") == "property-attributes-table__row" || e.Attr("class") == "property-attributes-table__row qa-living-area-attribute" {
			_, err := property.AddField(e.ChildTexts("*")...)
			if err != nil {
				log.Error().Err(err).Msg("error while scraping")
				return
			}
		}
	})

	collector.OnHTML("span", func(h *colly.HTMLElement) {
		if h.Attr("class") == "hcl-subheading hcl-subheading--size1 qa-property-price" {
			price, err := convertPrice(h.Text)
			if err != nil {
				log.Error().Err(err).Msg("error while scraping")
				return
			}

			property.Price = price
		}
	})

	collector.OnHTML("p", realtorFunc(&property))

	collector.OnHTML("a", agencyFunc(&property))

	collector.OnHTML("h1", addressFunc(&property))

	err := collector.Visit(url)
	if err != nil {
		return Property{}, errors.Wrapf(err, "error when visiting url: %s", url)
	}

	return property, nil
}

func realtorFunc(property *Property) func(*colly.HTMLElement) {
	return func(h *colly.HTMLElement) {
		if h.Attr("class") == "broker-card__text qa-broker-name" {
			re := regexp.MustCompile(`\n`)
			property.Realtor = strings.TrimSpace(re.ReplaceAllString(h.Text, ""))
		}
	}
}

func agencyFunc(property *Property) func(*colly.HTMLElement) {
	return func(h *colly.HTMLElement) {
		if h.Attr("class") == "hcl-link qa-broker-agency-listings-link" {
			property.Agency = h.Attr("data-broker-agency-name")
		}
	}
}

func addressFunc(property *Property) func(*colly.HTMLElement) {
	return func(h *colly.HTMLElement) {
		if h.Attr("class") == "qa-property-heading hcl-heading hcl-heading--size2" || h.Attr("class") == "qa-property-heading hcl-heading hcl-heading--size1" {
			property.Address = h.Text
		}
	}
}

func convertPrice(text string) (int, error) {
	ss := strings.Fields(text)
	aString := strings.Join(ss, "")
	ss = strings.Split(aString, "kr")
	i, err := strconv.ParseInt(ss[0], 10, 0)
	if err != nil {
		return 0, err
	}

	return int(i), nil
}