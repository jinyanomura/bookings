package forms

import (
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	var v url.Values
	form := New(v)

	if !form.Valid() {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	postData := url.Values{}
	form := New(postData)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("Required() shows valid when required fields are missing")
	}

	postData.Add("a", "a")
	postData.Add("b", "a")
	postData.Add("c", "a")

	form = New(postData)

	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("Required() shows invalid when should have been valid")
	}
}

func TestFormMinLength(t *testing.T) {
	postData := url.Values{}
	form := New(postData)

	form.MinLength("a", 3)
	if form.Valid() {
		t.Error("MinLength() shows valid when given field does not satisfy the minimum length")
	}

	err := form.Errors.Get("a")
	if err == "" {
		t.Error("Get() returned no error when there should be one")
	}

	postData.Add("a", "This should now work!")

	form = New(postData)

	form.MinLength("a", 3)
	if !form.Valid() {
		t.Error("MinLength() shows invalid when should have been valid")
	}

	err = form.Errors.Get("a")
	if err != "" {
		t.Error("Get() returned error when there should not be one")
	}
}

func TestFormIsEmail(t *testing.T) {
	postData := url.Values{}
	form := New(postData)

	form.IsEmail("a")
	if form.Valid() {
		t.Error("IsEmail() shows valid when given string is not a valid email address")
	}

	postData.Add("a", "test@something.com")

	form = New(postData)

	form.IsEmail("a")
	if !form.Valid() {
		t.Error("IsEmail() shows invalid when should have been valid")
	}
}

func TestFormHas(t *testing.T) {
	postData := url.Values{}
	form := New(postData)

	if form.Has("a") {
		t.Error("Has() shows valid when given field does not exist")
	}

	postData.Add("a", "This should now work!")

	form = New(postData)
	if !form.Has("a") {
		t.Error("Has() shows invalid when should have been valid")
	}
}