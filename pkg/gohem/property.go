package gohem

import (
	"strconv"
	"strings"

	"github.com/pingcap/errors"
)

type Property struct {
	Address            string      `json:"address"`
	Price              int         `json:"price"`
	Type               string      `json:"type"`
	Tenure             string      `json:"tenure"`
	Rooms              int         `json:"rooms"`
	Area               string      `json:"area"`
	Year               string      `json:"constructionYear"`
	HousingCooperative string      `json:"housingCooperative"`
	Fee                CostPerUnit `json:"fee"`
	OperatingCost      CostPerUnit `json:"operatingCost"`
	PriceSquareMeter   CostPerUnit `json:"priceSquareMeter"`
	Visits             int         `json:"visists"`
	Realtor            string      `json:"realtor"`
	Agency             string      `json:"agency"`
	Removed            bool        `json:"removed"`
	RemovedDate        string      `json:"removedDate,omitempty"`
}

type Properties struct {
	Properties []Property `json:"properties"`
}

func (p *Property) AddField(keyValuePairs ...string) ([]string, error) {
	if len(keyValuePairs) < 2 {
		return keyValuePairs, nil
	}

	key := keyValuePairs[0]
	value := keyValuePairs[1]

	err := p.addField(key, value)
	if err != nil {
		return keyValuePairs, errors.Wrap(err, "failed to add field to property")
	}

	return p.AddField(keyValuePairs[2:]...)
}

func (p *Property) addField(key, value string) error {
	switch key {
	case "Address":
		p.Address = value
	case "Bostadstyp":
		p.Type = value
	case "Upplåtelseform":
		p.Tenure = value
	case "Antal rum":
		rooms, err := parseRoomsField(value)
		if err != nil {
			return errors.Annotate(err, "failed to parse rooms field")
		}
		p.Rooms = rooms
	case "Boarea":
		p.Area = value
	case "Byggår":
		p.Year = value
	case "Förening":
		p.HousingCooperative = strings.Split(value, "\n")[0]
	case "Avgift":
		fee, err := parseCostPerUnitField(value)
		if err != nil {
			return errors.Annotate(err, "failed to parse fee field")
		}
		p.Fee = fee
	case "Driftkostnad":
		operatingCost, err := parseCostPerUnitField(value)
		if err != nil {
			return errors.Annotate(err, "failed to parse operating cost field")
		}
		p.OperatingCost = operatingCost
	case "Pris/m²":
		priceSquareMeter, err := parseCostPerUnitField(value)
		if err != nil {
			return errors.Annotate(err, "failed to parse price square meter field")
		}
		p.PriceSquareMeter = priceSquareMeter
	}

	return nil
}

func parseRoomsField(s string) (int, error) {
	ss := strings.Split(s, "")
	i, err := strconv.ParseInt(ss[0], 10, 0)
	if err != nil {
		return 1, errors.Wrap(err, "failed to parse int")
	}
	return int(i), nil
}

func parseCostPerUnitField(s string) (CostPerUnit, error) {
	ss := strings.Fields(s)

	currencyByUnit := strings.Split(ss[len(ss)-1], "/")
	currency := currencyByUnit[0]
	unit := currencyByUnit[1]

	ss = ss[:len(ss)-1]

	c := ""
	for _, d := range ss {
		c += d
	}

	cost, err := strconv.ParseInt(c, 10, 0)
	if err != nil {
		return CostPerUnit{}, errors.Wrap(err, "failed to parse cost")
	}

	return NewCostPerUnit(int(cost), currency, unit)
}

type CostPerUnit struct {
	Cost     int    `json:"cost"`
	Currency string `json:"currency"`
	Unit     string `json:"unit"`
}

func NewCostPerUnit(cost int, currency, unit string) (CostPerUnit, error) {
	return CostPerUnit{
		Cost:     cost,
		Currency: currency,
		Unit:     unit}, nil
}
