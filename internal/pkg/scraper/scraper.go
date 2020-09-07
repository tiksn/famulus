package scraper

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

func ListScrape(url string) {
	fmt.Printf("Scraping list: %s", url)
	fmt.Println()

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("div.main-content").ChildrenFiltered("div").Each(func(i int, s *goquery.Selection) {
		s.Each(func(i int, mainDiv *goquery.Selection) {
			element := mainDiv.Find("table tbody tr td")
			if element != nil {
				element := element.First().Find("span a")
				link, exists := element.Attr("href")
				if exists {
					ContactScrape(link)
				}
			}
		})
	})
}

func ContactScrape(url string) {
	fmt.Printf("Scraping contact: %s", url)
	fmt.Println()

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	element := doc.Find("body table:nth-child(8) tbody tr td:nth-child(1) table tbody tr")
	if element != nil {
		fmt.Println(element.Eq(0).Find("td h1").Text())
		fmt.Println(element.Eq(2).Find("td div div span").Text())
		onclick, exists := element.Eq(2).Find("td div div input").Attr("onclick")
		if exists {
			re := regexp.MustCompile("^\\s*showPhonesWithDigits\\(\\'(?P<id>\\d*)\\',\\s*\\'(?P<key>[A-Fa-f0-9]{40})\\'\\)\\;\\s*return\\s*false\\s*\\;\\s*$")
			matches := re.FindStringSubmatch(onclick)
			fmt.Println(matches[re.SubexpIndex("id")])
			fmt.Println(matches[re.SubexpIndex("key")])
		}
	}
}
