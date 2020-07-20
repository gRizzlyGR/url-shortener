package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var host = "localhost"
var port = "8080"
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

func getURL(w http.ResponseWriter, r *http.Request) {
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

	go increaseRedirections(urlKey)
	turnOffCache(w)
	http.Redirect(w, r, url, http.StatusPermanentRedirect)
}

func postURL(w http.ResponseWriter, r *http.Request) {
	log.Println("POST URL")
	// Generate key for storage
	urlKey := GenerateKey(false)

	processURLInfo(w, r, urlKey, http.StatusCreated)
	
	// body, err := ioutil.ReadAll(r.Body)

	// if err != nil {
	// 	log.Printf("Error reading body: %v\n", err)
	// 	sendJSON(w, http.StatusBadRequest, genericResponse{"Unable to read body"})
	// 	return
	// }

	// url := string(body)
	// valid := IsURLValid(url)

	// if !valid {
	// 	log.Printf("Invalid URL: %v\n", url)
	// 	sendJSON(w, http.StatusBadRequest, genericResponse{"Invalid URL: " + url + ". Be sure it's in the form http://www.placeholder.com"})
	// 	return
	// }

	// // Save the urlKey and the url for direct access
	// status := NewDao(shortURLDbPath).Save(urlKey, url)

	// if status == Error {
	// 	log.Println("Internal Server Error")
	// 	sendJSON(w, http.StatusInternalServerError, genericResponse{"Internal Server Error"})
	// 	return
	// }

	// // Init count redirections to 0
	// go NewDao(shortURLCountDbPath).Save(urlKey, "0")

	// fullPath := fmt.Sprintf("http://%s:%s/api/%s", host, port, urlKey)
	// log.Printf("URL %s is located at: %s\n", url, fullPath)

	// // Set Location header, that can be accessed with a get
	// w.Header().Set("Location", fullPath)

	// // Send location as JSON as well
	// sendJSON(w, http.StatusCreated, locationResponse{urlKey, fullPath, "OK"})
}

func putURL(w http.ResponseWriter, r *http.Request) {
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

	processURLInfo(w, r, urlKey, http.StatusOK)

	// body, err := ioutil.ReadAll(r.Body)

	// url := string(body)
	// valid := IsURLValid(url)

	// if !valid {
	// 	log.Printf("Invalid URL: %v\n", url)
	// 	sendJSON(w, http.StatusBadRequest, genericResponse{"Invalid URL: " + url + ". Be sure it's in the form http://www.placeholder.com"})
	// 	return
	// }

	// if err != nil {
	// 	log.Printf("Error reading body: %v\n", err)
	// 	sendJSON(w, http.StatusBadRequest, genericResponse{"Unable to read body"})
	// 	return
	// }

	// status = dao.Save(urlKey, string(body))

	// if status == Error {
	// 	log.Println("Internal Server Error")
	// 	sendJSON(w, http.StatusInternalServerError, genericResponse{"Internal Server Error"})
	// 	return
	// }

	// // Restore count redirections to 0 for new url key
	// go NewDao(shortURLCountDbPath).Save(urlKey, "0")

	// fullPath := fmt.Sprintf("http://%s:%s/api/%s", host, port, urlKey)
	// log.Printf("New URL %s is located at: %s\n", url, fullPath)

	// // Send location as JSON as well
	// sendJSON(w, http.StatusCreated, locationResponse{urlKey, fullPath, "OK"})
}

func deleteURL(w http.ResponseWriter, r *http.Request) {
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

func notFound(w http.ResponseWriter, r *http.Request) {
	sendJSON(w, http.StatusNotFound, genericResponse{"Not Found"})
}

func countRedirections(w http.ResponseWriter, r *http.Request) {
	log.Println("GET URL")
	urlKey := mux.Vars(r)[urlKeyPattern]
	count, status := NewDao(shortURLCountDbPath).FindByKey(urlKey)

	if status == Error {
		log.Println("Internal Server Error")
		sendJSON(w, http.StatusInternalServerError, genericResponse{"Internal Server Error"})
		return
	}

	if status == NotFound {
		sendJSON(w, http.StatusNotFound, genericResponse{"URL key not found"})
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

func processURLInfo(w http.ResponseWriter, r *http.Request, urlKey string, httpStatus int) {
	body, err := ioutil.ReadAll(r.Body)

	url := string(body)
	valid := IsURLValid(url)

	if !valid {
		log.Printf("Invalid URL: %v\n", url)
		sendJSON(w, http.StatusBadRequest, genericResponse{"Invalid URL: " + url + ". Be sure it's in the form http://www.placeholder.com"})
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

	fullPath := fmt.Sprintf("http://%s:%s/api/%s", host, port, urlKey)
	log.Printf("URL %s is located at: %s\n", url, fullPath)

	// Send location as JSON as well
	sendJSON(w, httpStatus, locationResponse{urlKey, fullPath, "OK"})
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api", postURL).Methods(http.MethodPost)
	r.HandleFunc("/api/count/{urlKey}", countRedirections).Methods(http.MethodGet)
	r.HandleFunc("/api/{urlKey}", getURL).Methods(http.MethodGet)
	r.HandleFunc("/api/{urlKey}", deleteURL).Methods(http.MethodDelete)
	r.HandleFunc("/api/{urlKey}", putURL).Methods(http.MethodPut)
	r.HandleFunc("/api", notFound)

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public/")))

	log.Printf("Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
