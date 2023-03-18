package main

import (
	"log"
	"mushturd/pkg/config"
	"mushturd/pkg/handlers"
	"mushturd/pkg/render"
	"net/http"
	"os"
)

func main() {

	var app config.AppConfig

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cac")
	}
	app.TemplateCache = tc

	render.NewTempates(&app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
		log.Printf("Defaulting to port %s", port)
	}

	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/picks", handlers.PicksHandler)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
