package main

import (
	"net/http"
	"net/url"
	"testing"
)

func TestHandleSignUp(t *testing.T) {
	form := url.Values{}
	form.Add("email", "barry@example.com")
	form.Add("password", "Password1")

	w := testHandler("POST", "/signup", form)

	assert(t, w.Code, http.StatusOK)
}

func TestHandleSignUpAlreadyExists(t *testing.T) {
	form := url.Values{}
	form.Add("email", "john@example.com")
	form.Add("password", "Password1")

	w := testHandler("POST", "/signup", form)

	assert(t, w.Code, http.StatusConflict)
}

func TestHandleSignUpNoEmail(t *testing.T) {
	form := url.Values{}
	form.Add("email", "john@example.com")
	form.Add("password", "Password1")

	w := testHandler("POST", "/signup", form)

	assert(t, w.Code, http.StatusConflict)
}

func TestHandleSignUpShortPassword(t *testing.T) {
	form := url.Values{}
	form.Add("email", "simon@example.com")
	form.Add("password", "short")

	w := testHandler("POST", "/signup", form)

	assert(t, w.Code, http.StatusBadRequest)
}

func TestHandleSignIn(t *testing.T) {
	form := url.Values{}
	form.Add("email", "john@example.com")
	form.Add("password", "Password1")

	w := testHandler("POST", "/signin", form)

	assert(t, w.Code, http.StatusOK)
}

func TestHandleSignInInvalidPassword(t *testing.T) {
	form := url.Values{}
	form.Add("email", "john@example.com")
	form.Add("password", "NotThePassword")

	w := testHandler("POST", "/signin", form)

	assert(t, w.Code, http.StatusUnauthorized)
}

func TestHandleSignInInvalidEmail(t *testing.T) {
	form := url.Values{}
	form.Add("email", "ghost@example.com")
	form.Add("password", "Password1")

	w := testHandler("POST", "/signin", form)

	assert(t, w.Code, http.StatusUnauthorized)
}
