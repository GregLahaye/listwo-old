package main

import (
	"net/http"
	"net/url"
	"testing"
)

func TestHandleGetColumns(t *testing.T) {
	w := testHandler("GET", "/columns?list=housework-uuid", nil, true)

	assert(t, w.Code, http.StatusOK)
}

func TestHandleGetColumnsUnauthorized(t *testing.T) {
	w := testHandler("GET", "/columns?list=housework-uuid", nil, false)

	assert(t, w.Code, http.StatusUnauthorized)
}

func TestHandleGetColumnsForbidden(t *testing.T) {
	w := testHandler("GET", "/columns?list=gardening-uuid", nil, true)

	assert(t, w.Code, http.StatusForbidden)
}

func TestHandleGetColumnsNotFound(t *testing.T) {
	w := testHandler("GET", "/columns?list=ghost-uuid", nil, true)

	assert(t, w.Code, http.StatusForbidden)
}

func TestHandleGetColumnsBadRequest(t *testing.T) {
	w := testHandler("GET", "/columns", nil, true)

	assert(t, w.Code, http.StatusBadRequest)
}

func TestHandleCreateColumn(t *testing.T) {
	form := url.Values{}
	form.Add("list", "housework-uuid")
	form.Add("title", "Vacuum Lounge")

	w := testHandler("POST", "/columns", form, true)

	assert(t, w.Code, http.StatusOK)
}

func TestHandleCreateColumnUnauthorized(t *testing.T) {
	form := url.Values{}
	form.Add("list", "housework-uuid")
	form.Add("title", "Lounge")

	w := testHandler("POST", "/columns", form, false)

	assert(t, w.Code, http.StatusUnauthorized)
}

func TestHandleCreateColumnForbidden(t *testing.T) {
	form := url.Values{}
	form.Add("list", "gardening-uuid")
	form.Add("title", "Roses")

	w := testHandler("POST", "/columns", form, true)

	assert(t, w.Code, http.StatusForbidden)
}

func TestHandleCreateColumnNotFound(t *testing.T) {
	form := url.Values{}
	form.Add("list", "ghost-uuid")
	form.Add("title", "Attic")

	w := testHandler("POST", "/columns", form, true)

	assert(t, w.Code, http.StatusForbidden)
}

func TestHandleCreateColumnBadRequest(t *testing.T) {
	form := url.Values{}
	form.Add("list", "ghost-uuid")

	w := testHandler("POST", "/columns", form, true)

	assert(t, w.Code, http.StatusBadRequest)
}

func TestHandleDeleteColumn(t *testing.T) {
	form := url.Values{}
	form.Add("id", "bathroom-uuid")

	w := testHandler("DELETE", "/columns", form, true)

	assert(t, w.Code, http.StatusOK)
}

func TestHandleDeleteColumnUnauthorized(t *testing.T) {
	form := url.Values{}
	form.Add("id", "laundry-uuid")

	w := testHandler("DELETE", "/columns", form, false)

	assert(t, w.Code, http.StatusUnauthorized)
}

func TestHandleDeleteColumnForbidden(t *testing.T) {
	form := url.Values{}
	form.Add("id", "gardening-uuid")

	w := testHandler("DELETE", "/columns", form, true)

	assert(t, w.Code, http.StatusForbidden)
}

func TestHandleDeleteColumnNotFound(t *testing.T) {
	form := url.Values{}
	form.Add("id", "ghost-uuid")

	w := testHandler("DELETE", "/columns", form, true)

	assert(t, w.Code, http.StatusForbidden)
}

func TestHandleDeleteColumnBadRequest(t *testing.T) {
	w := testHandler("DELETE", "/columns", nil, true)

	assert(t, w.Code, http.StatusBadRequest)
}

func TestHandleUpdateColumn(t *testing.T) {
	form := url.Values{}
	form.Add("id", "laundry-uuid")
	form.Add("title", "My list name")

	w := testHandler("PATCH", "/columns", form, true)

	assert(t, w.Code, http.StatusOK)
}

func TestHandleUpdateColumnUnauthorized(t *testing.T) {
	form := url.Values{}
	form.Add("id", "laundry-uuid")
	form.Add("title", "My new list")

	w := testHandler("PATCH", "/columns", form, false)

	assert(t, w.Code, http.StatusUnauthorized)
}

func TestHandleUpdateColumnForbidden(t *testing.T) {
	form := url.Values{}
	form.Add("id", "gardening-uuid")
	form.Add("title", "Outdoor List")

	w := testHandler("PATCH", "/columns", form, true)

	assert(t, w.Code, http.StatusForbidden)
}

func TestHandleUpdateColumnNotFound(t *testing.T) {
	form := url.Values{}
	form.Add("id", "ghost-uuid")
	form.Add("title", "Spooky List")

	w := testHandler("PATCH", "/columns", form, true)

	assert(t, w.Code, http.StatusForbidden)
}

func TestHandleUpdateColumnBadRequest(t *testing.T) {
	form := url.Values{}
	form.Add("id", "laundry-uuid")

	w := testHandler("PATCH", "/columns", nil, true)

	assert(t, w.Code, http.StatusBadRequest)
}
