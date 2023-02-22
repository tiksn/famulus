package scraper

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/charmap"
)

type scraperA struct {
	sourceAddress      string
	sourcePhoneAddress string
	interval           time.Duration
}

func (c *scraperA) ListScrape() ([]Contact, error) {
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

	doc.Find("div.main-content").ChildrenFiltered("div").Each(func(i int, s *goquery.Selection) {
		s.Each(func(i int, mainDiv *goquery.Selection) {
			element := mainDiv.Find("table tbody tr td")
			if element != nil {
				element := element.First().Find("span a")
				link, exists := element.Attr("href")
				if exists {
					contact, err := c.contactScrape(link, c.sourcePhoneAddress, c.interval)
					if err != nil {
						fmt.Println(err)
					} else {
						result = append(result, contact)
						fmt.Println(contact.GetNumbers())
					}
				}
			}
		})
	})

	return result, nil
}

func (c *scraperA) contactScrape(url string, phoneUrl string, interval time.Duration) (Contact, error) {
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

	r := charmap.Windows1251.NewDecoder().Reader(res.Body)
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	element := doc.Find("body table:nth-child(8) tbody tr td:nth-child(1) table tbody tr")
	if element != nil {
		t := element.Eq(0).Find("td h1").Text()
		d := strings.TrimSpace(element.Eq(2).Find("td div").Nodes[0].FirstChild.Data)
		onclick, exists := element.Eq(2).Find("td div div input").Attr("onclick")
		if exists {
			re := regexp.MustCompile("^\\s*showPhonesWithDigits\\(\\'(?P<id>\\d*)\\',\\s*\\'(?P<key>[A-Fa-f0-9]{40})\\'\\)\\;\\s*return\\s*false\\s*\\;\\s*$")
			matches := re.FindStringSubmatch(onclick)

			nums, err := c.numbersScrape(matches[re.SubexpIndex("id")], matches[re.SubexpIndex("key")], phoneUrl, interval)
			if err != nil {
				return nil, err
			}

			return &contact{
				numbers:     nums,
				website:     url,
				title:       t,
				description: d,
			}, nil
		}
	}

	return nil, errors.New("Unable to find contact information")
}

func (c *scraperA) numbersScrape(id string, key string, phoneUrl string, interval time.Duration) (string, error) {
	payload := strings.NewReader(fmt.Sprintf("i=%s&s=%s", id, key))

	client := &http.Client{}
	time.Sleep(interval)
	req, err := http.NewRequest("POST", phoneUrl, payload)

	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}

	element := doc.Find("div span")
	if element != nil {
		return element.Text(), nil
	} else {
		return "", errors.New("Unable to find phone number")
	}
}
