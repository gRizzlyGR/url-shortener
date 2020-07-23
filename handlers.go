package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"fmt"
	"io/ioutil"

	"github.com/gorilla/mux"
)

var urlKeyPattern = "urlKey"

// For convenience, DBs are stored in the source code folder
var shortURLDbPath = "shortURLDb"
var shortURLCountDbPath = "shortURLCountDb"

type genericResponse struct {
	Message string `json:"message"`
}

type locationResponse struct {
	URLKey   string `json:"urlKey"`
	Location string `json:"location"`
	Message  string `json:"message"`
}

type redirectionsCountResponse struct {
	RedirectionsCount int `json:"redirectionsCount"`
}

// GetURL redirects to the original URL using the provided URL key
func GetURL(w http.ResponseWriter, r *http.Request) {
	log.Println("GET URL")
	urlKey := mux.Vars(r)[urlKeyPattern]
	url, status := NewDao(shortURLDbPath).FindByKey(urlKey)

	if status == NotFound {
		sendJSON(w, http.StatusNotFound, genericResponse{"No URL found for key " + urlKey})
		return
	}

	if status == Error {
		log.Println("Internal Server Error")
		sendJSON(w, http.StatusInternalServerError, genericResponse{"Internal Server Error"})
		return
	}

	increaseRedirections(urlKey)
	turnOffCache(w)
	http.Redirect(w, r, url, http.StatusPermanentRedirect)
}

// PostURL creates a short URL and returns its location
func PostURL(w http.ResponseWriter, r *http.Request) {
	log.Println("POST URL")

	// Generate key for storage
	urlKey := GenerateKey(false)

	writeURLInfo(w, r, urlKey, http.StatusCreated)
}

// PutURL replaces a URL with a new one, using the provided URL key
func PutURL(w http.ResponseWriter, r *http.Request) {
	log.Println("PUT URL")
	urlKey := mux.Vars(r)[urlKeyPattern]
	dao := NewDao(shortURLDbPath)
	exists, status := dao.DoesExist(urlKey)

	if status == Error {
		log.Println("Internal Server Error")
		sendJSON(w, http.StatusInternalServerError, genericResponse{"Internal Server Error"})
		return
	}

	if !exists {
		sendJSON(w, http.StatusNotFound, genericResponse{"No URL found for key " + urlKey})
		return
	}

	writeURLInfo(w, r, urlKey, http.StatusOK)
}

// DeleteURL removes a URL using the provided URL key
func DeleteURL(w http.ResponseWriter, r *http.Request) {
	log.Println("DELETE URL")

	urlKey := mux.Vars(r)[urlKeyPattern]

	shortURLDao := NewDao(shortURLDbPath)

	status := shortURLDao.RemoveByKey(urlKey)
	if status == NotFound {
		sendJSON(w, http.StatusNotFound, genericResponse{"No URL found for key " + urlKey})
		return
	}

	if status == Error {
		sendJSON(w, http.StatusInternalServerError, genericResponse{"Internal Server Error"})
		return
	}

	// Once the url key has been removed, remove its count as well
	NewDao(shortURLCountDbPath).RemoveByKey(urlKey)

	sendJSON(w, http.StatusOK, genericResponse{"URL successfully deleted for key " + urlKey})
}

// SendNotFound throws a 404 Not Found error for missing handling
func SendNotFound(w http.ResponseWriter, r *http.Request) {
	sendJSON(w, http.StatusNotFound, genericResponse{"Not Found"})
}

// CountRedirections returns how many times a URL redirection occurred
func CountRedirections(w http.ResponseWriter, r *http.Request) {
	log.Println("GET COUNT")
	urlKey := mux.Vars(r)[urlKeyPattern]
	count, status := NewDao(shortURLCountDbPath).FindByKey(urlKey)

	if status == NotFound {
		sendJSON(w, http.StatusNotFound, genericResponse{"URL key not found"})
		return
	}

	if status == Error {
		log.Println("Internal Server Error")
		sendJSON(w, http.StatusInternalServerError, genericResponse{"Internal Server Error"})
		return
	}

	intCount := 0

	if count != "" {
		intCount, _ = strconv.Atoi(count)
	}

	message := fmt.Sprintf("Redirections count for %s: %s", urlKey, count)
	log.Println(message)

	sendJSON(w, http.StatusOK, redirectionsCountResponse{intCount})
}

func increaseRedirections(urlKey string) {
	dao := NewDao(shortURLCountDbPath)
	count, _ := dao.FindByKey(urlKey)
	intCount, _ := strconv.Atoi(count)
	intCount++
	count = strconv.Itoa(intCount)
	log.Printf("Increased redirections count for %s: %s", urlKey, count)
	dao.Save(urlKey, count)
}

func writeURLInfo(w http.ResponseWriter, r *http.Request, urlKey string, httpStatus int) {
	body, err := ioutil.ReadAll(r.Body)

	url := string(body)
	valid := IsURLValid(url)

	if !valid {
		log.Printf("Invalid URL: '%v'\n", url)
		sendJSON(w, http.StatusBadRequest, genericResponse{"Invalid URL: '" + url + "'. Be sure it's in the form http://www.placeholder.com"})
		return
	}

	if err != nil {
		log.Printf("Error reading body: %v\n", err)
		sendJSON(w, http.StatusBadRequest, genericResponse{"Unable to read body"})
		return
	}

	status := NewDao(shortURLDbPath).Save(urlKey, string(body))

	if status == Error {
		log.Println("Internal Server Error")
		sendJSON(w, http.StatusInternalServerError, genericResponse{"Internal Server Error"})
		return
	}

	// Set count redirections to 0
	go NewDao(shortURLCountDbPath).Save(urlKey, "0")

	fullPath := fmt.Sprintf("http://%s:%s/api/%s", host, defaultPort, urlKey)
	log.Printf("URL %s is located at: %s\n", url, fullPath)

	// Send location as JSON as well
	sendJSON(w, httpStatus, locationResponse{urlKey, fullPath, "OK"})
}

func sendJSON(w http.ResponseWriter, code int, response interface{}) {
	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)
	w.Write(jsonResponse)
}

// turnOffCache disables caching in almost all browser using different headers.
// It may not work for the redirection of unsecured http requests (i.e. no
// https)
func turnOffCache(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "private, no-cache, no-store, must-revalidate, max-age=0, proxy-revalidate, s-maxage=0")
	w.Header().Set("Expires", "0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Vary", "*")
}
