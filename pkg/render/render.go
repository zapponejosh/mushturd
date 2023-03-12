package render

import (
	"html/template"
	"log"
	"net/http"
)

func RenderTemplate(w http.ResponseWriter, tmpl string, data any) {
	parsedTemplate, _ := template.ParseFiles("./templates/"+tmpl, "./templates/layout.gohtml")

	err := parsedTemplate.Execute(w, data)
	if err != nil {
		log.Println("error parsing template: ", err)
		return
	}
}
