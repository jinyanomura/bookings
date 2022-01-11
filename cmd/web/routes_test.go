package main

import (
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jinyanomura/bookings/pkg/config"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch mux.(type) {
	case *chi.Mux:
	default: t.Error("type is not *chi.Mux")
	}
}