package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var host = "localhost"
var port = "8080"
var urlKeyPattern = "urlKey"

// For convenience, DBs are stored in the source code folder
var shortURLDbPath = "shortURLDb"
var fullURLDbPath = "fullURLDb"

type genericResponse struct {
	Message string `json:"message"`
}

type locationResponse struct {
	URLKey   string `json:"urlKey"`
	Location string `json:"location"`
	Message string `json:"message"`
}

func sendJSON(w http.ResponseWriter, code int, res interface{}) {
	jsonResponse, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)
	w.Write(jsonResponse)
}

func get(w http.ResponseWriter, r *http.Request) {
	log.Println("GET")
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

	http.Redirect(w, r, url, http.StatusPermanentRedirect)
}

func post(w http.ResponseWriter, r *http.Request) {
	log.Println("POST")
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Printf("Error reading body: %v\n", err)
		sendJSON(w, http.StatusBadRequest, genericResponse{"Unable to read body"})
		return
	}

	url := string(body)
	valid := IsURLValid(url)

	if !valid {
		log.Printf("Invalid URL: %v\n", url)
		sendJSON(w, http.StatusBadRequest, genericResponse{"Invalid URL: " + url})
		return
	}

	fullURLDao := NewDao(fullURLDbPath)
	urlKey, status := fullURLDao.FindByKey(url)

	if status == Error {
		log.Println("Internal Server Error")
		sendJSON(w, http.StatusInternalServerError, genericResponse{"Internal Server Error"})
		return
	}

	if urlKey != "" {
		log.Printf("URL already exists: %v\n", url)
		fullPath := fmt.Sprintf("http://%s:%s/api/%s", host, port, urlKey)
		
		// Set Location header, that can be accessed with a get
		w.Header().Set("Location", fullPath)
		sendJSON(w, http.StatusForbidden, locationResponse{urlKey, fullPath, "URL already exists in the DB"})
		return		
	}

	// Geneate key for storage
	urlKey = GenerateKey(false)

	// Save the urlKey and the url for direct access
	status = NewDao(shortURLDbPath).Save(urlKey, url)

	if status == Error {
		log.Println("Internal Server Error")
		sendJSON(w, http.StatusInternalServerError, genericResponse{"Internal Server Error"})
		return
	}

	// To check if the URL is already present in the DB, we need to save the URL
	// itself in the DB as key, so to prevent in the future that the same URL is
	// shortened again. This approach makes the DB a Bimap de facto, with a
	// short URL that is key of a full URL, and a full URL that is key of a
	// short URL
	status = fullURLDao.Save(url, urlKey)

	fullPath := fmt.Sprintf("http://%s:%s/api/%s", host, port, urlKey)
	log.Printf("URL %s is located at: %s\n", url, fullPath)

	// Set Location header, that can be accessed with a get
	w.Header().Set("Location", fullPath)

	// Send location as JSON as well
	sendJSON(w, http.StatusCreated, locationResponse{urlKey, fullPath, "OK"})
}

func delete(w http.ResponseWriter, r *http.Request) {
	log.Println("DELETE")

	urlKey := mux.Vars(r)[urlKeyPattern]

	dao := NewDao(shortURLDbPath)

	// To delete the key, we need to delete the URL treated as key as well.
	// Ignore if no url has been found
	url, _ := dao.FindByKey(urlKey)
	
	status := dao.RemoveByKey(urlKey)
	if status == NotFound {
		sendJSON(w, http.StatusNotFound, genericResponse{"No URL found for key " + urlKey})
		return
	}
	
	if status == Error {
		sendJSON(w, http.StatusInternalServerError, genericResponse{"Internal Server Error"})
		return
	}

	// Once the url key has been removed, remove the URL as well.
	NewDao(fullURLDbPath).RemoveByKey(url)
	
	sendJSON(w, http.StatusOK, genericResponse{"URL successfully deleted for key " + urlKey})
}

func notFound(w http.ResponseWriter, r *http.Request) {
	sendJSON(w, http.StatusNotFound, genericResponse{"Not Found"})
}

func main() {
	r := mux.NewRouter()

	//r.PathPrefix("/").Handler(http.FileServer(http.Dir("public/")))
	r.HandleFunc("/api", post).Methods(http.MethodPost)
	r.HandleFunc("/api/{urlKey}", get).Methods(http.MethodGet)
	r.HandleFunc("/api/{urlKey}", delete).Methods(http.MethodDelete)
	r.HandleFunc("/api", notFound)

	log.Printf("Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
