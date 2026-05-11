package views

import (
	"errors"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

type Renderer struct {
	templates *template.Template
}

func New(files fs.FS) (*Renderer, error) {
	tmpl, err := template.ParseFS(files, "templates/*.html")
	if err != nil {
		return nil, err
	}
	return &Renderer{templates: tmpl}, nil
}

func (r *Renderer) HTML(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := r.templates.ExecuteTemplate(w, name, data); err != nil {
		if errors.Is(err, http.ErrHandlerTimeout) {
			return
		}
		log.Printf("template render failed: %v", err)
	}
}
