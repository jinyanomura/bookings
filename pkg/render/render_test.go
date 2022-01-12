package render

import (
	"net/http"
	"testing"

	"github.com/jinyanomura/bookings/pkg/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "123")

	result := AddDefaultData(r, &td)

	if result.Flash != "123" {
		t.Error("flash value of 123 not found in session")
	}
}

func TestTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc

	var w myWriter
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	err = Template(w, r, "home.page.html", &models.TemplateData{})
	if err != nil {
		t.Error("error writing template to blowser")
	}

	err = Template(w, r, "dummy.page.html", &models.TemplateData{})
	if err == nil {
		t.Error("rendered template which does not exist")
	}
}

func TestNewTemplates(t *testing.T) {
	SetNewTemplates(app)
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil
}