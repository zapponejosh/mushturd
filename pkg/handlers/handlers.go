package handlers

import (
	"mushturd/pkg/models"
	"mushturd/pkg/redis"
	"mushturd/pkg/scraper"
	"net/http"
	"sort"

	"mushturd/pkg/render"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Define the Redis cache key for the data
	cacheKey := "musher_data"

	mushers, err := redis.GetMushersFromCacheOrAPI(cacheKey, scraper.Scraper)
	if err != nil {
		http.Error(w, "Error with cache or fetching data. Look at logs.", http.StatusInternalServerError)
		return
	}
	render.RenderTemplate(w, "home.page.gohtml", mushers)
}

func PicksHandler(w http.ResponseWriter, r *http.Request) {
	type Pick struct {
		Name       string
		Bib        int
		Position   int
		PointValue int
		Rookie     bool
		Status     string
	}

	type User struct {
		Picks      []*Pick
		Name       string
		PointTotal int
	}
	// get musher data
	cacheKey := "musher_data"
	musherData, err := redis.GetMushersFromCacheOrAPI(cacheKey, scraper.Scraper)
	if err != nil {
		http.Error(w, "Error with cache or fetching data. Look at logs.", http.StatusInternalServerError)
		return
	}
	topRookie := models.Musher{
		CurrentPos: 33,
	}
	mushersAsPicks := make(map[int]Pick)
	for _, m := range musherData {
		var pk Pick
		pk.Bib = m.Bib
		pk.Name = m.Name
		pk.Position = m.CurrentPos
		pk.Status = m.Status
		// inverse of position 1st place = 33 points plus 5 for rookie?
		pk.PointValue = 33 - m.CurrentPos + 1 // calc rookie after
		pk.Rookie = m.Rookie
		if m.Rookie && m.Status != "Scratched" {
			if m.CurrentPos < topRookie.CurrentPos {
				topRookie = m
			}
		}
		if pk.Status == "Scratched" {
			pk.PointValue = 0
			pk.Position = 100
		}
		mushersAsPicks[m.Bib] = pk
	}
	rookieOfYear := mushersAsPicks[topRookie.Bib]
	rookieOfYear.PointValue += 5
	mushersAsPicks[rookieOfYear.Bib] = rookieOfYear
	// find pick bib in data and attach position and point value

	PickData := []*User{
		{
			Name: "Josh",
			Picks: []*Pick{
				{Bib: 14},
				{Bib: 19, Rookie: true},
				{Bib: 5},
				{Bib: 20},
			},
		},
		{
			Name: "Chelsea",
			Picks: []*Pick{
				{Bib: 26},
				{Bib: 8},
				{Bib: 28},
				{Bib: 18, Rookie: true},
			},
		},
		{
			Name: "JoDee",
			Picks: []*Pick{
				{Bib: 15},
				{Bib: 2},
				{Bib: 16},
				{Bib: 29, Rookie: true},
			},
		},
		{
			Name: "Marissa",
			Picks: []*Pick{
				{Bib: 4},
				{Bib: 6},
				{Bib: 25},
				{Bib: 13, Rookie: true},
			},
		},
		{
			Name: "Bill",
			Picks: []*Pick{
				{Bib: 33},
				{Bib: 23},
				{Bib: 9},
				{Bib: 22, Rookie: true},
			},
		},
		{
			Name: "Gigi",
			Picks: []*Pick{
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
				p.Status = d.Status
			}
		}
		u.PointTotal = points

		sort.Slice(u.Picks, func(i, j int) bool {

			if u.Picks[i].Rookie && !u.Picks[j].Rookie {
				return false
			} else {
				if u.Picks[i].Position < u.Picks[j].Position {
					return true
				} else {
					return false
				}
			}
		})
	}

	sort.Slice(PickData, func(i, j int) bool {
		return PickData[i].PointTotal > PickData[j].PointTotal
	})
	render.RenderTemplate(w, "picks.page.gohtml", PickData)
}

func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
