package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/jinyanomura/bookings/pkg/config"
	"github.com/jinyanomura/bookings/pkg/handlers"
	"github.com/jinyanomura/bookings/pkg/models"
	"github.com/jinyanomura/bookings/pkg/render"
)

var app config.AppConfig
var session *scs.SessionManager

func main() {
	// what am I gonna store in session info
	gob.Register(models.Reservation{})

	//change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	c, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = c
	app.UseCache = false

	handlers.NewHandler(&app)
	render.SetNewTemplates(&app)

	fmt.Println("Starting server on port 8080...")

	srv := &http.Server {
		Addr: ":8080",
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}