package main

import (
	"log"
	"mushturd/pkg/handlers"
	"net/http"
	"os"
)

func main() {

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
