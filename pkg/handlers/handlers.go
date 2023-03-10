package handlers

import (
	"fmt"
	"mushturd/pkg/scraper"
	"net/http"

	"mushturd/pkg/render"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	musherData := scraper.Scraper()
	render.RenderTemplate(w, "home.gohtml", musherData)
}
func PicksHandler(w http.ResponseWriter, r *http.Request) {
	type Pick struct {
		Name       string
		Bib        int
		Position   int
		PointValue int
		Rookie     bool
	}

	type User struct {
		Picks      []Pick
		Name       string
		PointTotal int
	}
	// get musher data
	musherData := scraper.Scraper()
	topRookie := scraper.Musher{
		CurrentPos: 33,
	}
	mushersAsPicks := make(map[int]Pick)
	for _, m := range musherData {
		var pk Pick
		pk.Bib = m.Bib
		pk.Name = m.Name
		pk.Position = m.CurrentPos
		// inverse of position 1st place = 33 points plus 5 for rookie?
		pk.PointValue = 33 - m.CurrentPos + 1 // calc rookie after
		pk.Rookie = m.Rookie
		if m.Rookie {
			if m.CurrentPos < topRookie.CurrentPos {
				topRookie = m
			}
		}
		mushersAsPicks[m.Bib] = pk
	}
	rookieOfYear := mushersAsPicks[topRookie.Bib]
	rookieOfYear.PointValue += 5
	fmt.Println("Adding rookie of year points")
	mushersAsPicks[rookieOfYear.Bib] = rookieOfYear
	// find pick bib in data and attach position and point value

	// TODO need to adjust this to be pointer to user so I can modify in range loop
	PickData := []User{
		{
			Name: "Josh",
			Picks: []Pick{
				{Bib: 14},
				{Bib: 5},
				{Bib: 20},
				{Bib: 19, Rookie: true},
			},
		},
		{
			Name: "Chelsea",
			Picks: []Pick{
				{Bib: 26},
				{Bib: 8},
				{Bib: 28},
				{Bib: 18, Rookie: true},
			},
		},
		{
			Name: "JoDee",
			Picks: []Pick{
				{Bib: 15},
				{Bib: 2},
				{Bib: 16},
				{Bib: 29, Rookie: true},
			},
		},
		{
			Name: "Marissa",
			Picks: []Pick{
				{Bib: 4},
				{Bib: 6},
				{Bib: 25},
				{Bib: 13, Rookie: true},
			},
		},
		{
			Name: "Bill",
			Picks: []Pick{
				{Bib: 33},
				{Bib: 23},
				{Bib: 9},
				{Bib: 22, Rookie: true},
			},
		},
		{
			Name: "Gigi",
			Picks: []Pick{
				{Bib: 14},
				{Bib: 31},
				{Bib: 26},
				{Bib: 19, Rookie: true},
			},
		},
	}

	for _, u := range PickData {
		var points int
		for _, p := range u.Picks {
			d, ok := mushersAsPicks[p.Bib]
			if ok {
				p.Name = d.Name
				p.PointValue = d.PointValue
				p.Position = d.Position
				points += d.PointValue
			}
		}
		u.PointTotal = points
		fmt.Printf("%s: %d", u.Name, u.PointTotal)
	}

	render.RenderTemplate(w, "picks.gohtml", nil)
}
