package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

// SETUP

// Silence log output during unit tests. Importantly you need to call Run() once
// you've done what you need
func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

// GET
func TestGetURLShouldReturn308(t *testing.T) {
	// Override for test
	shortURLDbPath = "testDb"

	key := "mockKey"
	value := "mockURL"

	dao := NewDao(shortURLDbPath)

	dao.Save(key, value)

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

	os.RemoveAll(shortURLDbPath)
}
func TestGetURLShouldReturn404(t *testing.T) {
	// Override for test
	shortURLDbPath = "testDb"

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

	os.RemoveAll(shortURLDbPath)
}

// POST
func TestPostShouldReturn201(t *testing.T) {
	// Override for test
	shortURLDbPath = "testDb"

	r, _ := http.NewRequest("POST", "/api", strings.NewReader("https://www.verylongtestwebsitethatneedstobeshortened.com"))
	w := httptest.NewRecorder()

	PostURL(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("got: %v, want: %v", w.Code, http.StatusCreated)
	}

	var response locationResponse

	json.Unmarshal(w.Body.Bytes(), &response)

	if response.Message != "OK" {
		t.Errorf("got: %v, want: %v", response.Message, "OK")
	}

	os.RemoveAll(shortURLDbPath)
}

func TestPostUrlShouldReturn400(t *testing.T) {
	// Override for test
	shortURLDbPath = "testDb"

	r, _ := http.NewRequest("POST", "/api", strings.NewReader("invalid url"))
	w := httptest.NewRecorder()

	PostURL(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got: %v, want: %v", w.Code, http.StatusBadRequest)
	}

	os.RemoveAll(shortURLDbPath)
}

// PUT
func TestPutShouldReturn200(t *testing.T) {
	// Override for test
	shortURLDbPath = "testDb"

	key := "mockKey"
	value := "mockURL"

	dao := NewDao(shortURLDbPath)

	dao.Save(key, value)

	r, _ := http.NewRequest("PUT", "/api/"+key, strings.NewReader("https://www.verylongtestwebsitethatneedstobeshortened.com"))
	w := httptest.NewRecorder()

	vars := map[string]string{
		"urlKey": key,
	}

	r = mux.SetURLVars(r, vars)

	PutURL(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("got: %v, want: %v", w.Code, http.StatusOK)
	}

	if w.Body == nil {
		t.Errorf("Null body")
	}

	os.RemoveAll(shortURLDbPath)
}
func TestPutShouldReturn400(t *testing.T) {
	// Override for test
	shortURLDbPath = "testDb"

	key := "mockKey"
	value := "mockURL"

	dao := NewDao(shortURLDbPath)

	dao.Save(key, value)

	r, _ := http.NewRequest("PUT", "/api/"+key, strings.NewReader("invalid url"))
	w := httptest.NewRecorder()

	vars := map[string]string{
		"urlKey": key,
	}

	r = mux.SetURLVars(r, vars)

	PutURL(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got: %v, want: %v", w.Code, http.StatusBadRequest)
	}

	os.RemoveAll(shortURLDbPath)
}
func TestPutShouldReturn404(t *testing.T) {
	// Override for test
	shortURLDbPath = "testDb"

	key := "notFound"

	r, _ := http.NewRequest("PUT", "/api/"+key, strings.NewReader("https://www.verylongtestwebsitethatneedstobeshortened.com"))
	w := httptest.NewRecorder()

	PutURL(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("got: %v, want: %v", w.Code, http.StatusNotFound)
	}

	os.RemoveAll(shortURLDbPath)
}

//DELETE
func TestDeleteShouldReturn200(t *testing.T) {
	// Override for test
	shortURLDbPath = "testDb"

	key := "mockKey"
	value := "mockURL"

	dao := NewDao(shortURLDbPath)

	dao.Save(key, value)

	r, _ := http.NewRequest("DELETE", "/api/"+key, nil)
	w := httptest.NewRecorder()

	vars := map[string]string{
		"urlKey": key,
	}

	r = mux.SetURLVars(r, vars)

	DeleteURL(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("got: %v, want: %v", w.Code, http.StatusOK)
	}

	os.RemoveAll(shortURLDbPath)
}
func TestDeleteShouldReturn404(t *testing.T) {
	// Override for test
	shortURLDbPath = "testDb"

	key := "mockKey"

	r, _ := http.NewRequest("DELETE", "/api/"+key, nil)
	w := httptest.NewRecorder()

	DeleteURL(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("got: %v, want: %v", w.Code, http.StatusNotFound)
	}

	os.RemoveAll(shortURLDbPath)
}

// Redirections count
func TestCountShouldReturn200AndNumbersOfRedirections(t *testing.T) {
	// Override for test
	shortURLDbPath = "testDb"

	key := "mockKey"
	value := "mockURL"

	dao := NewDao(shortURLDbPath)

	dao.Save(key, value)

	r, _ := http.NewRequest("GET", "/api/"+key, nil)
	w := httptest.NewRecorder()

	vars := map[string]string{
		"urlKey": key,
	}

	r = mux.SetURLVars(r, vars)

	// Get it three times to increase by 3
	GetURL(w, r)
	GetURL(w, r)
	GetURL(w, r)

	// Get redirections count
	r, _ = http.NewRequest("GET", "/api/count"+key, nil)
	w = httptest.NewRecorder()

	vars = map[string]string{
		"urlKey": key,
	}

	r = mux.SetURLVars(r, vars)
	CountRedirections(w, r)

	var response redirectionsCountResponse

	json.Unmarshal(w.Body.Bytes(), &response)

	if response.RedirectionsCount != 3 {
		t.Errorf("got: %v, want: %v", response.RedirectionsCount, 3)
	}

	os.RemoveAll(shortURLDbPath)
}
func TestCountShouldReturn404(t *testing.T) {
	// Override for test
	shortURLDbPath = "testDb"

	key := "notFound"

	r, _ := http.NewRequest("GET", "/api/count"+key, nil)
	w := httptest.NewRecorder()

	CountRedirections(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("got: %v, want: %v", w.Code, http.StatusNotFound)
	}

	os.RemoveAll(shortURLDbPath)
}
