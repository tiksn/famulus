package scraper

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func ContactScrape(url string) {
	fmt.Printf("Scraping: %s", url)
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
			element := mainDiv.Find("table")
			if element != nil {
				element := element.Find("tbody")
				if element != nil {
					element := element.Find("tr")
					if element != nil {
						element := element.Find("td")
						if element != nil {
							element := element.First().Find("span a")
							link, exists := element.Attr("href")
							if exists {
								fmt.Println(link)
							}
						}
					}
				}
			}
		})
	})
}
