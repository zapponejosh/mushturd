package render

import (
	"bytes"
	"html/template"
	"log"
	"mushturd/pkg/config"
	"net/http"
	"path/filepath"
)

var app *config.AppConfig

func NewTempates(a *config.AppConfig) {
	app = a
}

func RenderTemplate(w http.ResponseWriter, tmpl string, data any) {
	// get tc from app config

	templateCache := app.TemplateCache

	// get template from cache
	if t, ok := templateCache[tmpl]; !ok {
		log.Fatalf("No template found for %s\n", tmpl)
	} else {
		buf := new(bytes.Buffer)
		err := t.Execute(buf, data)
		if err != nil {
			log.Println(err)
		}
		_, err = buf.WriteTo(w)
		if err != nil {
			log.Println(err)
		}
	}

}

func CreateTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// get all pages
	pages, err := filepath.Glob("./templates/*.page.gohtml")
	if err != nil {
		return cache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return cache, err
		}

		layouts, err := filepath.Glob("./templates/*.layout.gohtml")
		if err != nil {
			return cache, err
		}
		if len(layouts) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.gohtml")
			if err != nil {
				return cache, err
			}
		}
		cache[name] = ts
	}
	return cache, nil

}
