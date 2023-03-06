package scraper

import (
	"fmt"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

// bib - int (primary key)

// current place - int

// current checkpoint (in or out, dogs) - jsonb(name, status)

// mandatory rests taken - string[]
// speed

// rookie - bool

// final place? - int

type Checkpoint struct {
	in      bool
	dogsIn  int
	out     bool
	dogsOut int
}

type Musher struct {
	currentPos   int    // 0
	name         string //
	rookie       bool
	bib          int
	checkpoint   Checkpoint
	speed        float32
	inCheckpoint bool
}

func Scraper() {
	fmt.Println("Running scraper")
	c := colly.NewCollector(
		// Visit only domains: https://iditarod.com/
		colly.AllowedDomains("iditarod.com"),
	)

	// On every a element which has href attribute call callback
	c.OnHTML("#current-standings-data", func(e *colly.HTMLElement) {
		var musherSlice []Musher
		e.DOM.ChildrenFiltered("tr").EachWithBreak(func(i int, s *goquery.Selection) bool {
			var m Musher
			// fmt.Println(i + 1)
			s.Children().Each(func(idx int, col *goquery.Selection) {
				// fmt.Printf("%d - %s  |  ", i, col.Text())
				switch idx {
				case 0:
					m.currentPos = i + 1
				case 1:
					m.name = col.Text()
				case 3:
					b, err := strconv.Atoi(col.Text())
					if err != nil {
						return
					}
					m.bib = b
				}
			})
			musherSlice = append(musherSlice, m)
			return true
		})
		for _, m := range musherSlice {
			fmt.Printf("%+v\n", m)
		}

	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Start scraping on https://iditarod.com/race/2023/standings/
	c.Visit("https://iditarod.com/race/2023/standings/")
}
