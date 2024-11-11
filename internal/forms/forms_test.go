package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid instead of valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields are missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	r, _ = http.NewRequest("POST", "/whatever", nil)

	r.PostForm = postedData
}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	// Positive Scenario
	has := form.Has("non-existent field")
	if has {
		t.Error("form shows that it has a field when it does not")
	}

	// Negative Scenario
	postedData := url.Values{}
	postedData.Add("a", "a")

	form = New(postedData)

	has = form.Has("a")
	if !has {
		t.Error("form shows that it does not have field when it should")
	}

}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.MinLength("x", 10)
	if form.Valid() {
		t.Error("form shows min length for non-existent field")
	}

	isError := form.Errors.Get("x")
	if isError == "" {
		t.Error("should have an error, but did not get one")
	}

	postedData := url.Values{}
	postedData.Add("correctLength", "12345")

	form = New(postedData)

	form.MinLength("correctLength", 5)
	if !form.Valid() {
		t.Error("form shows not min length when it has exact length")
	}

	isError = form.Errors.Get("correctLength")
	if isError != "" {
		t.Error("should not have an error, but got one")
	}

	postedData.Add("tooShort", "123")
	form = New(postedData)

	form.MinLength("tooShort", 5)
	if form.Valid() {
		t.Error("form shows min length when input has shorter length")
	}

	postedData.Add("empty", "")
	form = New(postedData)

	form.MinLength("empty", 5)
	if form.Valid() {
		t.Error("form shows min length when input is empty")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.IsEmail("x")
	if form.Valid() {
		t.Error("form shows is email for non-existent field")
	}

	postedData = url.Values{}

	postedData.Add("validEmail", "myEmail@mail.com")
	form = New(postedData)

	form.IsEmail("validEmail")
	if !form.Valid() {
		t.Error("form shows invalid when email is valid")
	}

	postedData.Add("validEmail_alias", "myEmail+123@mail.com")
	form = New(postedData)

	form.IsEmail("validEmail_alias")
	if !form.Valid() {
		t.Error("form shows invalid when email alias is valid")
	}

	postedData.Add("invalidEmail_no_at", "123.com")
	form = New(postedData)

	form.IsEmail("invalidEmail_no_at")
	if form.Valid() {
		t.Error("form shows valid when email has no 'at' sign")
	}

	postedData.Add("invalidEmail_trailing_period_before_at", "123.@mail.com")
	form = New(postedData)

	form.IsEmail("invalidEmail_trailing_period_before_at")
	if form.Valid() {
		t.Error("form shows valid when email has '.' right before 'at' sign")
	}

	postedData.Add("invalidEmail_no_postfix", "123@mail")
	form = New(postedData)

	form.IsEmail("invalidEmail_no_postfix")
	if form.Valid() {
		t.Error("form shows valid when email has no postfix (.com, .ca, etc.)")
	}
}
