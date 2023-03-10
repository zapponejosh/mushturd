package main

import (
	"mushturd/pkg/handlers"
	"net/http"
)

const PORT = ":3000"

func main() {

	// fmt.Println("Running go in main repos dir")
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/picks", handlers.PicksHandler)
	http.ListenAndServe(PORT, nil)

}
