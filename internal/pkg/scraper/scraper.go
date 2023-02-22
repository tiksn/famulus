package scraper

import (
	"errors"
	"time"
)

type Scraper interface {
	ListScrape() ([]Contact, error)
}

func GetScraper(kind string, url string, phoneUrl string, interval time.Duration) (Scraper, error) {
	switch kind {
	case "A":
		return &scraperA{
			sourceAddress:      url,
			sourcePhoneAddress: phoneUrl,
			interval:           interval,
		}, nil
	case "B":
		return &scraperB{
			sourceAddress:      url,
			sourcePhoneAddress: phoneUrl,
			interval:           interval,
		}, nil
	default:
		return nil, errors.New("unknown source kind")
	}
}
