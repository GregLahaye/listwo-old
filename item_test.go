package main

import (
	"net/http"
	"net/url"
	"testing"
)

func TestHandleGetItems(t *testing.T) {
	w := testHandler("GET", "/items?column=laundry-uuid", nil, true)

	assert(t, w.Code, http.StatusOK)
}

func TestHandleGetItemsUnauthorized(t *testing.T) {
	w := testHandler("GET", "/items?column=laundry-uuid", nil, false)

	assert(t, w.Code, http.StatusUnauthorized)
}

func TestHandleGetItemsForbidden(t *testing.T) {
	w := testHandler("GET", "/items?column=roses-uuid", nil, true)

	assert(t, w.Code, http.StatusForbidden)
}

func TestHandleGetItemsNotFound(t *testing.T) {
	w := testHandler("GET", "/items?column=ghost-uuid", nil, true)

	assert(t, w.Code, http.StatusForbidden)
}

func TestHandleGetItemsBadRequest(t *testing.T) {
	w := testHandler("GET", "/items", nil, true)

	assert(t, w.Code, http.StatusBadRequest)
}

func TestHandleCreateItem(t *testing.T) {
	form := url.Values{}
	form.Add("column", "kitchen-uuid")
	form.Add("title", "Wipe Surfaces")

	w := testHandler("POST", "/items", form, true)

	assert(t, w.Code, http.StatusOK)
}

func TestHandleCreateItemUnauthorized(t *testing.T) {
	form := url.Values{}
	form.Add("column", "kitchen-uuid")
	form.Add("title", "Wipe Surfaces")

	w := testHandler("POST", "/items", form, false)

	assert(t, w.Code, http.StatusUnauthorized)
}

func TestHandleCreateItemForbidden(t *testing.T) {
	form := url.Values{}
	form.Add("column", "gardening-uuid")
	form.Add("title", "Water Roses")

	w := testHandler("POST", "/items", form, true)

	assert(t, w.Code, http.StatusForbidden)
}

func TestHandleCreateItemNotFound(t *testing.T) {
	form := url.Values{}
	form.Add("column", "ghost-uuid")
	form.Add("title", "Clear Attic")

	w := testHandler("POST", "/items", form, true)

	assert(t, w.Code, http.StatusForbidden)
}

func TestHandleCreateItemBadRequest(t *testing.T) {
	form := url.Values{}
	form.Add("column", "kitchen-uuid")

	w := testHandler("POST", "/items", form, true)

	assert(t, w.Code, http.StatusBadRequest)
}

func TestHandleDeleteItem(t *testing.T) {
	form := url.Values{}
	form.Add("id", "iron-uuid")

	w := testHandler("DELETE", "/items", form, true)

	assert(t, w.Code, http.StatusOK)
}

func TestHandleDeleteItemUnauthorized(t *testing.T) {
	form := url.Values{}
	form.Add("id", "iron-uuid")

	w := testHandler("DELETE", "/items", form, false)

	assert(t, w.Code, http.StatusUnauthorized)
}

func TestHandleDeleteItemForbidden(t *testing.T) {
	form := url.Values{}
	form.Add("id", "sing-uuid")

	w := testHandler("DELETE", "/items", form, true)

	assert(t, w.Code, http.StatusForbidden)
}

func TestHandleDeleteItemNotFound(t *testing.T) {
	form := url.Values{}
	form.Add("id", "ghost-uuid")

	w := testHandler("DELETE", "/items", form, true)

	assert(t, w.Code, http.StatusForbidden)
}

func TestHandleDeleteItemBadRequest(t *testing.T) {
	w := testHandler("DELETE", "/items", nil, true)

	assert(t, w.Code, http.StatusBadRequest)
}

func TestHandleUpdateItem(t *testing.T) {
	form := url.Values{}
	form.Add("id", "hang-uuid")
	form.Add("title", "My item name")
	form.Add("dst", "0")

	w := testHandler("PATCH", "/items", form, true)

	assert(t, w.Code, http.StatusOK)
}

func TestHandleUpdateItemUnauthorized(t *testing.T) {
	form := url.Values{}
	form.Add("id", "hang-uuid")
	form.Add("title", "My new item")
	form.Add("dst", "0")

	w := testHandler("PATCH", "/items", form, false)

	assert(t, w.Code, http.StatusUnauthorized)
}

func TestHandleUpdateItemForbidden(t *testing.T) {
	form := url.Values{}
	form.Add("id", "sing-uuid")
	form.Add("title", "Chorus item")
	form.Add("dst", "1")

	w := testHandler("PATCH", "/items", form, true)

	assert(t, w.Code, http.StatusForbidden)
}

func TestHandleUpdateItemNotFound(t *testing.T) {
	form := url.Values{}
	form.Add("id", "ghost-uuid")
	form.Add("title", "Spooky item")
	form.Add("dst", "1")

	w := testHandler("PATCH", "/items", form, true)

	assert(t, w.Code, http.StatusForbidden)
}

func TestHandleUpdateItemBadRequest(t *testing.T) {
	form := url.Values{}
	form.Add("id", "hang-uuid")

	w := testHandler("PATCH", "/items", nil, true)

	assert(t, w.Code, http.StatusBadRequest)
}
