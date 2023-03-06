package handlers

import (
	"net/http"

	"mushturd/pkg/render"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "home.gohtml")
}
func AboutHandler(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "about.gohtml")
}
