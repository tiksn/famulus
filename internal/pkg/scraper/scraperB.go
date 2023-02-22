package scraper

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type scraperB struct {
	sourceAddress      string
	sourcePhoneAddress string
	interval           time.Duration
}

func (c *scraperB) ListScrape() ([]Contact, error) {
	fmt.Printf("Scraping list: %s", c.sourceAddress)
	fmt.Println()

	time.Sleep(c.interval)
	res, err := http.Get(c.sourceAddress)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	result := []Contact{}

	doc.Find("div.main").ChildrenFiltered("div.items").ChildrenFiltered("div.item").Each(func(i int, s *goquery.Selection) {
		s.Each(func(i int, itemDiv *goquery.Selection) {
			attrClass, exists := itemDiv.Attr("class")
			if exists && attrClass == "item" {
				element := itemDiv.Find("a")
				if element != nil {
					link, exists := element.Attr("href")
					if exists {
						sourceUrl, err := url.Parse(c.sourceAddress)
						if err != nil {
							fmt.Println(err)
						} else {
							contactUrl := sourceUrl
							contactUrl.Path = link
							contact, err := c.contactScrape(contactUrl.String(), c.sourcePhoneAddress, c.interval)
							if err != nil {
								fmt.Println(err)
							} else {
								result = append(result, contact)
								fmt.Println(contact.GetNumbers())
							}
						}
					}
				}
			}
		})
	})

	return result, nil
}

func (c *scraperB) contactScrape(url string, phoneUrl string, interval time.Duration) (Contact, error) {
	fmt.Printf("Scraping contact: %s", url)
	fmt.Println()

	time.Sleep(interval)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	element := doc.Find("div.right")
	if element != nil {
		t := element.Find("h1").Text()
		nums := strings.TrimSpace(element.Find("div.phone0mt div.phone02 a.phoneMT").Text())
		d := strings.TrimSpace(element.Find("div.text").Text())

		return &contact{
			numbers:     nums,
			website:     url,
			title:       t,
			description: d,
		}, nil
	}

	return nil, errors.New("Unable to find contact information")
}
