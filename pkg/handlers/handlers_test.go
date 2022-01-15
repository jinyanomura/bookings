package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/jinyanomura/bookings/pkg/models"
)

var tests = []struct{
	name string
	url string
	method string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"generals-quarters", "/generals-quarters", "GET", http.StatusOK},
	{"majors-suite", "/majors-suite", "GET", http.StatusOK},
	{"search-availability", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()

	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range tests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID: 1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned unexpected status code: %d", rr.Code)
	}

	// test case where reservation is not in session(reset everything)
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned unexpected status code: %d", rr.Code)
	}

	// test with invalid roomID
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned unexpected status code: %d", rr.Code)
	}
}

func testPostReservation(key string) int {
	postData := url.Values{}
	postData.Add("start_date", "2050-01-01")
	postData.Add("end_date", "2050-01-02")
	postData.Add("first_name", "John")
	postData.Add("last_name", "Smith")
	postData.Add("email", "john@smith.com")
	postData.Add("phone", "555-555-5555")
	postData.Add("room_id", "1")

	var req *http.Request
	if key == "missing_body" {
		req, _ = http.NewRequest("POST", "/make-reservation", nil)
	} else {
		switch key {
		case "start_date": postData.Set("start_date", "invalid")
		case "end_date": postData.Set("end_date", "invalid")
		case "room_id": postData.Set("room_id", "invalid")
		case "form_validation": postData.Set("first_name", "a")
		case "insert_reservation": postData.Set("room_id", "2")
		case "insert_restriction": postData.Set("start_date", "2222-02-02")
		}
		req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postData.Encode()))
	}

	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	
	return rr.Code
}

func TestRepository_PostReservation(t *testing.T) {
	code := testPostReservation("")
	if code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned unexpected status code: %d", code)
	}

	// test for missing post body
	code = testPostReservation("missing_body")
	if code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned unexpected status code: %d for missing post-body", code)
	}

	// test for invalid start_date
	code = testPostReservation("start_date")
	if code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned unexpected status code: %d for invalid start_date", code)
	}

	// test for invalid end_date
	code = testPostReservation("end_date")
	if code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned unexpected status code: %d for invalid end_date", code)
	}

	// test for invalid rooom_id
	code = testPostReservation("room_id")
	if code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned unexpected status code: %d for invalid room_id", code)
	}

	// test for invalid form data
	code = testPostReservation("form_validation")
	if code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned unexpected status code: %d for form validation", code)
	}

	// test for failure for inserting reservation into database
	code = testPostReservation("insert_reservation")
	if code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned unexpected status code: %d for inserting reservation failure", code)
	}

	// test for failure for inserting restriction into database
	code = testPostReservation("insert_restriction")
	if code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned unexpected status code: %d for inserting restriction failure", code)
	}
}

func TestRepository_AvailabilityJSON(t *testing.T) {
	// case1: rooms are not available
	reqBody := "start_date=2050-01-01&end_date=2050-01-02&room_id=1"

	// create request with context
	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// make a handler
	handler := http.HandlerFunc(Repo.AvailabilityJSON)

	// make request to the handler
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	var j jsonResponse
	err := json.Unmarshal([]byte(rr.Body.Bytes()), &j)
	if err != nil {
		t.Error("failed to parse json")
	}
}

func getCtx(r *http.Request) context.Context {
	ctx, err := session.Load(r.Context(), r.Header.Get("X-Session"))
	if err != nil {
		log.Panicln(err)
	}
	return ctx
}