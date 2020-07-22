package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetURLShouldReturn308(t *testing.T) {
	key := "mockKey"
	value := "mockURL"

	// Override for test
	shortURLDbPath = "testDb"

	dao := NewDao(shortURLDbPath)

	dao.Save(key, value)
	defer dao.RemoveByKey(key)

	r, _ := http.NewRequest("GET", "/api/"+key, nil)
	w := httptest.NewRecorder()

	vars := map[string]string{
		"urlKey": key,
	}

	r = mux.SetURLVars(r, vars)

	GetURL(w, r)

	if w.Code != http.StatusPermanentRedirect {
		t.Errorf("got: %v, want: %v", w.Code, http.StatusPermanentRedirect)
	}
}
func TestGetURLShouldReturn404(t *testing.T) {
	key := "notFound"
	r, _ := http.NewRequest("GET", "/api/"+key, nil)
	w := httptest.NewRecorder()

	vars := map[string]string{
		"urlKey": key,
	}

	r = mux.SetURLVars(r, vars)

	GetURL(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("got: %v, want: %v", w.Code, http.StatusNotFound)
	}
}
