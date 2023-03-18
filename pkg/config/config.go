package config

import "html/template"

// holds application config
type AppConfig struct {
	TemplateCache map[string]*template.Template
}
