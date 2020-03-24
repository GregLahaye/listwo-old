package main

import (
	"net/http"
	"net/url"
	"testing"
)

func TestHandleGetList(t *testing.T) {
	w := testHandler("GET", "/list?id=schoolwork-uuid", nil, true)

	assert(t, w.Code, http.StatusOK)
}

func TestHandleGetListUnauthorized(t *testing.T) {
	w := testHandler("GET", "/list?id=schoolwork-uuid", nil, false)

	assert(t, w.Code, http.StatusUnauthorized)
}

func TestHandleGetListForbidden(t *testing.T) {
	w := testHandler("GET", "/list?id=gardening-uuid", nil, true)

	assert(t, w.Code, http.StatusForbidden)
}

func TestHandleGetListNotFound(t *testing.T) {
	w := testHandler("GET", "/list?id=ghost-uuid", nil, true)

	assert(t, w.Code, http.StatusForbidden)
}

func TestHandleGetLists(t *testing.T) {
	w := testHandler("GET", "/lists", nil, true)

	assert(t, w.Code, http.StatusOK)
}

func TestHandleGetListsUnauthorized(t *testing.T) {
	w := testHandler("GET", "/lists", nil, false)

	assert(t, w.Code, http.StatusUnauthorized)
}

func TestHandleCreateList(t *testing.T) {
	form := url.Values{}
	form.Add("title", "Bills")

	w := testHandler("POST", "/lists", form, true)

	assert(t, w.Code, http.StatusOK)
}

func TestHandleCreateListUnauthorized(t *testing.T) {
	form := url.Values{}
	form.Add("title", "Bills")

	w := testHandler("POST", "/lists", form, false)

	assert(t, w.Code, http.StatusUnauthorized)
}

func TestHandleCreateListBadRequest(t *testing.T) {
	form := url.Values{}
	form.Add("title", "")

	w := testHandler("POST", "/lists", form, true)

	assert(t, w.Code, http.StatusBadRequest)
}

func TestHandleDeleteList(t *testing.T) {
	form := url.Values{}
	form.Add("id", "schoolwork-uuid")

	w := testHandler("DELETE", "/lists", form, true)

	assert(t, w.Code, http.StatusOK)
}

func TestHandleDeleteListUnauthorized(t *testing.T) {
	form := url.Values{}
	form.Add("id", "housework-uuid")

	w := testHandler("DELETE", "/lists", form, false)

	assert(t, w.Code, http.StatusUnauthorized)
}

func TestHandleDeleteListForbidden(t *testing.T) {
	form := url.Values{}
	form.Add("id", "gardening-uuid")

	w := testHandler("DELETE", "/lists", form, true)

	assert(t, w.Code, http.StatusForbidden)
}

func TestHandleDeleteListBadRequest(t *testing.T) {
	form := url.Values{}
	form.Add("id", "")

	w := testHandler("DELETE", "/lists", form, true)

	assert(t, w.Code, http.StatusBadRequest)
}

func TestHandleUpdateList(t *testing.T) {
	form := url.Values{}
	form.Add("id", "housework-uuid")
	form.Add("title", "My new list")

	w := testHandler("PATCH", "/lists", form, true)

	assert(t, w.Code, http.StatusOK)
}

func TestHandleUpdateListUnauthorized(t *testing.T) {
	form := url.Values{}
	form.Add("id", "housework-uuid")
	form.Add("title", "Updated list name")

	w := testHandler("PATCH", "/lists", form, false)

	assert(t, w.Code, http.StatusUnauthorized)
}

func TestHandleUpdateListForbidden(t *testing.T) {
	form := url.Values{}
	form.Add("id", "gardening-uuid")
	form.Add("title", "Another list name")

	w := testHandler("PATCH", "/lists", form, true)

	assert(t, w.Code, http.StatusForbidden)
}

func TestHandleUpdateListBadRequest(t *testing.T) {
	form := url.Values{}
	form.Add("id", "housework-uuid")
	form.Add("title", "")

	w := testHandler("PATCH", "/lists", form, true)

	assert(t, w.Code, http.StatusBadRequest)
}
