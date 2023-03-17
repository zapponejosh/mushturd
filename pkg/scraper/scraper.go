package scraper

import (
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type Checkpoint struct {
	Name    string
	In      bool
	DogsIn  int
	Out     bool
	DogsOut int
}

type Musher struct {
	CurrentPos       int    // 0
	Name             string //
	Rookie           bool
	Bib              int
	LatestCheckpoint Checkpoint
	Speed            float32
	InCheckpoint     bool
	Status           string // scratched, finished
}

func Scraper() []Musher {
	var (
		finished, racing, scratched, musherData []Musher
	)

	// TODO process all rows using a function and each passes into the proper array
	processRows := func(i int, s *goquery.Selection) bool {
		var m Musher
		// log.Println(i + 1)

		// TODO - just add values into a map instead of iterating over each col
		s.Children().Each(func(idx int, col *goquery.Selection) {
			// log.Printf("%d - %s  |  ", idx, col.Text())
			switch idx {
			// position
			case 0:
				m.CurrentPos = i + 1
			// name
			case 1:
				m.Name = col.Text()
			// bib number
			case 2:
				b, err := strconv.Atoi(col.Text())
				if err != nil {
					return
				}
				m.Bib = b
			// latest checkpoint
			case 3:
				cp := col.Text()
				if cp == "Nome" {
					m.Status = "Finished"
				}
				m.LatestCheckpoint.Name = col.Text()
			// time in latest checkpoint
			case 4:
				if len(col.Text()) >= 1 {
					m.LatestCheckpoint.In = true
				}
			// dogs in
			case 5:
				num, err := strconv.Atoi(col.Text())
				if err != nil {
					return
				}
				m.LatestCheckpoint.DogsIn = num
			// time out of latest checkpoint
			case 6:
				if len(col.Text()) >= 1 {
					m.LatestCheckpoint.Out = true
				}
			// Dogs out
			case 7:
				num, err := strconv.Atoi(col.Text())
				if err != nil {
					return
				}
				m.LatestCheckpoint.DogsOut = num

			// rest in latest checkpoint
			case 8:
				//nothing

			// time enroute to latest checkpoint
			case 9:
				//nothing

			// previous checkpoint
			case 10:
				//nothing

			// time out of previous checkpoint
			case 11:
				// nothing
			// speed between checkpoints
			case 12:
				num, err := strconv.ParseFloat(col.Text(), 32)
				if err != nil {
					return
				}
				m.Speed = float32(num)
			// 8 hour rest complete
			case 13:
				// if len(col.Children().Nodes) >= 1 {
				// 	log.Println("8 hour complete")
				// }
			// 24 hours rest complete
			case 14:
				// if len(col.Children().Nodes) >= 1 {
				// 	log.Println("24 hour complete")
				// }
			// status
			case 15:
				// if len(col.Text()) >= 1 {
				// 	m.Status = col.Text()
				// }
				m.Status = "Scratched"
			// do nothing
			default:
				//do nothing
			}
		})

		if strings.Contains(m.Name, "(r)") {
			m.Rookie = true
			m.Name = m.Name[0:(len(m.Name) - 4)]
		}
		if !m.LatestCheckpoint.Out {
			m.InCheckpoint = true
		}
		scratched = append(scratched, m)
		return true
	}
	// var musherSlice []Musher
	log.Println("Running scraper")
	c := colly.NewCollector(
		// Visit only domains: https://iditarod.com/
		colly.AllowedDomains("iditarod.com"),
	)

	// On every a element which has href attribute call callback
	c.OnHTML(".current-standings tbody", func(e *colly.HTMLElement) {

		parentAttr, _ := e.DOM.Parent().Attr("class")
		if strings.Contains(parentAttr, "stats-table-racing") { // racing table
			e.DOM.ChildrenFiltered("tr").EachWithBreak(func(i int, s *goquery.Selection) bool {
				var m Musher
				// log.Println(i + 1)

				// TODO - just add values into a map instead of iterating over each col
				s.Children().Each(func(idx int, col *goquery.Selection) {
					// log.Printf("%d - %s  |  ", idx, col.Text())
					switch idx {
					// position
					case 0:
						str := col.Text()
						str = strings.TrimSuffix(str, "*")
						p, err := strconv.Atoi(str)
						if err != nil {
							return
						}
						m.CurrentPos = p
					// name
					case 1:
						m.Name = col.Text()
					// bib number
					case 2:
						b, err := strconv.Atoi(col.Text())
						if err != nil {
							return
						}
						m.Bib = b
					// latest checkpoint
					case 3:
						cp := col.Text()
						if cp == "Nome" {
							m.Status = "Finished"
						}
						m.LatestCheckpoint.Name = col.Text()
					// time in latest checkpoint
					case 4:
						if len(col.Text()) >= 1 {
							m.LatestCheckpoint.In = true
						}
					// dogs in
					case 5:
						num, err := strconv.Atoi(col.Text())
						if err != nil {
							return
						}
						m.LatestCheckpoint.DogsIn = num
					// time out of latest checkpoint
					case 6:
						if len(col.Text()) >= 1 {
							m.LatestCheckpoint.Out = true
						}
					// Dogs out
					case 7:
						num, err := strconv.Atoi(col.Text())
						if err != nil {
							return
						}
						m.LatestCheckpoint.DogsOut = num

					// rest in latest checkpoint
					case 8:
						//nothing

					// time enroute to latest checkpoint
					case 9:
						//nothing

					// previous checkpoint
					case 10:
						//nothing

					// time out of previous checkpoint
					case 11:
						// nothing
					// speed between checkpoints
					case 12:
						num, err := strconv.ParseFloat(col.Text(), 32)
						if err != nil {
							return
						}
						m.Speed = float32(num)
					// 8 hour rest complete
					case 13:
						// if len(col.Children().Nodes) >= 1 {
						// 	log.Println("8 hour complete")
						// }
					// 24 hours rest complete
					case 14:
						// if len(col.Children().Nodes) >= 1 {
						// 	log.Println("24 hour complete")
						// }
					// status
					case 15:
						if len(col.Text()) >= 1 {
							m.Status = col.Text()
						}
					// do nothing
					default:
						//do nothing
					}
				})

				if strings.Contains(m.Name, "(r)") {
					m.Rookie = true
					m.Name = m.Name[0:(len(m.Name) - 4)]
				}
				if !m.LatestCheckpoint.Out {
					m.InCheckpoint = true
				}
				racing = append(racing, m)
				return true
			})
			musherData = append(musherData, racing...)

		} else if strings.Contains(parentAttr, "stats-table-out_of_race") { // scratched table
			e.DOM.ChildrenFiltered("tr").EachWithBreak(processRows)
			musherData = append(musherData, scratched...)

		} else { // finished table
			e.DOM.ChildrenFiltered("tr").EachWithBreak(func(i int, s *goquery.Selection) bool {
				var m Musher
				// log.Println(i + 1)

				// TODO - just add values into a map instead of iterating over each col
				s.Children().Each(func(idx int, col *goquery.Selection) {
					// log.Printf("%d - %s  |  ", idx, col.Text())
					switch idx {
					// position
					case 0:
						str := col.Text()
						str = strings.TrimSuffix(str, "*")
						p, err := strconv.Atoi(str)
						if err != nil {
							return
						}
						m.CurrentPos = p
					// name
					case 1:
						m.Name = col.Text()
					// bib number
					case 2:
						b, err := strconv.Atoi(col.Text())
						if err != nil {
							return
						}
						m.Bib = b
					// latest checkpoint
					case 3:
						cp := col.Text()
						if cp == "Nome" {
							m.Status = "Finished"
						}
						m.LatestCheckpoint.Name = col.Text()
					// time in latest checkpoint
					case 4:
						if len(col.Text()) >= 1 {
							m.LatestCheckpoint.In = true
						}
					// dogs in
					case 5:
						num, err := strconv.Atoi(col.Text())
						if err != nil {
							return
						}
						m.LatestCheckpoint.DogsIn = num
					// time out of latest checkpoint
					case 6:
						if len(col.Text()) >= 1 {
							m.LatestCheckpoint.Out = true
						}
					// Dogs out
					case 7:
						num, err := strconv.Atoi(col.Text())
						if err != nil {
							return
						}
						m.LatestCheckpoint.DogsOut = num

					// rest in latest checkpoint
					case 8:
						//nothing

					// time enroute to latest checkpoint
					case 9:
						//nothing

					// previous checkpoint
					case 10:
						//nothing

					// time out of previous checkpoint
					case 11:
						// nothing
					// speed between checkpoints
					case 12:
						num, err := strconv.ParseFloat(col.Text(), 32)
						if err != nil {
							return
						}
						m.Speed = float32(num)
					// 8 hour rest complete
					case 13:
						// if len(col.Children().Nodes) >= 1 {
						// 	log.Println("8 hour complete")
						// }
					// 24 hours rest complete
					case 14:
						// if len(col.Children().Nodes) >= 1 {
						// 	log.Println("24 hour complete")
						// }
					// status
					case 15:
						if len(col.Text()) >= 1 {
							m.Status = col.Text()
						}
					// do nothing
					default:
						//do nothing
					}
				})

				if strings.Contains(m.Name, "(r)") {
					m.Rookie = true
					m.Name = m.Name[0:(len(m.Name) - 4)]
				}
				if !m.LatestCheckpoint.Out {
					m.InCheckpoint = true
				}
				finished = append(finished, m)
				return true
			})
			musherData = append(musherData, finished...)
		}

	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Start scraping on https://iditarod.com/race/2023/standings/
	c.Visit("https://iditarod.com/race/2023/standings/")

	return musherData
}
