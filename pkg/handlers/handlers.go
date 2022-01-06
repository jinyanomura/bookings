package handlers

import (
	"net/http"

	"github.com/jinyanomura/bookings/pkg/config"
	"github.com/jinyanomura/bookings/pkg/models"
	"github.com/jinyanomura/bookings/pkg/render"
)

type Repository struct {
	App *config.AppConfig
}

var Repo *Repository

func NewHandler(a *config.AppConfig) {
	Repo = &Repository{
		App: a,
	}
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")

	stringMap := make(map[string]string)
	stringMap["test"] = "Hello there!"
	stringMap["remote_ip"] = remoteIp
	render.RenderTemplate(w, "home.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, "about.page.html", &models.TemplateData{})
}

