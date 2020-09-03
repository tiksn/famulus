package scraper

import (
	"fmt"
	"log"
	"net/http"

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
		s.Each(func(i1 int, mainDiv *goquery.Selection) {
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
}
