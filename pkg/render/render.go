package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/jinyanomura/bookings/pkg/config"
	"github.com/jinyanomura/bookings/pkg/models"
	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{}
var app *config.AppConfig
var pathToTemplates = "./templates"

// SetNewTemplates sets the config for the template package
func SetNewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(r *http.Request ,td *models.TemplateData) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	return td
}

// Template renders templates using html/template
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	var c map[string]*template.Template

	if app.UseCache {
		c = app.TemplateCache
	} else {
		c, _ = CreateTemplateCache()
	}

	t, ok := c[tmpl]
	if !ok {
		return errors.New("cannot get template from template cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(r, td)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		return err
	}

	return nil
}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplates))
	if err != nil {
		return cache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return cache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
		if err != nil {
			return cache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
			if err != nil {
				return cache, err
			}
		}

		cache[name] = ts
	}

	return cache, nil
}