package main

import (
	"mushturd/pkg/handlers"
	"mushturd/pkg/scraper"
	"net/http"
)

const PORT = ":3000"

func main() {
	scraper.Scraper()
	// fmt.Println("Running go in main repos dir")
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/about", handlers.AboutHandler)
	http.ListenAndServe(PORT, nil)

}
