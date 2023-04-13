package scraper

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type scraperC struct {
	sourceAddress      string
	sourcePhoneAddress string
	interval           time.Duration
}

type scraperCPhoneResponse struct {
	Success     bool   `json:"success"`
	PhoneNumber string `json:"phone"`
}

func (c *scraperC) ListScrape() ([]Contact, error) {
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

	doc.Find("div.wrapper").ChildrenFiltered("div.row:nth-child(2)").ChildrenFiltered("div.item-col").ChildrenFiltered("a").Each(func(i int, s *goquery.Selection) {
		s.Each(func(i int, itemDiv *goquery.Selection) {
			_, exists := itemDiv.Attr("title")
			if exists {
				link, exists := itemDiv.Attr("href")
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

func (c *scraperC) contactScrape(url string, phoneUrl string, interval time.Duration) (Contact, error) {
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
	element := doc.Find("div.wrapper div.row:nth-child(3)")
	if element != nil {
		t := strings.TrimSpace(element.Find("div.content-info-wrapper div.content-info-col div.content-name h1").Text())

		nums, err := c.numbersScrape(element, phoneUrl, interval)
		if err != nil {
			return nil, err
		}
		descriptions := make(map[string]string)
		element.Find("div.content-info-wrapper div.content-info-col div.content-desc div.tab-content div.tab-pane").Each(func(i int, sDescription *goquery.Selection) {
			lang, exists := sDescription.Attr("id")
			if exists {
				desc := sDescription.Text()
				desc = strings.TrimSpace(desc)
				descriptions[lang] = desc
				fmt.Printf("[%s]: %s", lang, desc)
			}
		})

		return &contact{
			numbers:     nums,
			website:     url,
			title:       t,
			description: descriptions["EN"],
		}, nil
	}

	return nil, errors.New("Unable to find contact information")
}

func (c *scraperC) numbersScrape(element *goquery.Selection, phoneUrl string, interval time.Duration) (string, error) {
	var phoneId string
	var phoneIdExists bool
	element.Find("div.content-info-wrapper div.content-info-col div.content-name div.content-contact div.contact-elem").Each(func(i int, sContact *goquery.Selection) {
		contactElement := sContact.Find("a")
		if !phoneIdExists {
			phoneId, phoneIdExists = contactElement.Attr("data-phone-id")
		}

		href, exists := contactElement.Attr("href")
		if exists && href != "#" {
			fmt.Printf("Contact Reference is %s", href)
		}
	})

	if phoneIdExists {
		payload := strings.NewReader(fmt.Sprintf("id=%s", phoneId))

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

		decoder := json.NewDecoder(res.Body)
		var response scraperCPhoneResponse
		err = decoder.Decode(&response)
		if err != nil {
			return "", err
		}

		if response.Success {
			return strings.ReplaceAll(strings.TrimSpace(response.PhoneNumber), " ", ""), nil
		} else {
			return "", errors.New("phone number resolution request failed")
		}
	}

	return "", errors.New("number is not found")
}
