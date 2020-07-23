package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

var host = "localhost"
var defaultPort = "8080"

func main() {
	if len(os.Args) > 1 {
		defaultPort = initPort(os.Args[1])
	}

	r := mux.NewRouter()

	r.HandleFunc("/api", PostURL).Methods(http.MethodPost)
	r.HandleFunc("/api/count/{urlKey}", CountRedirections).Methods(http.MethodGet)
	r.HandleFunc("/api/{urlKey}", GetURL).Methods(http.MethodGet)
	r.HandleFunc("/api/{urlKey}", DeleteURL).Methods(http.MethodDelete)
	r.HandleFunc("/api/{urlKey}", PutURL).Methods(http.MethodPut)
	r.HandleFunc("/api", SendNotFound)

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public/")))

	log.Printf("Listening on port %s\n", defaultPort)
	log.Fatal(http.ListenAndServe(":"+defaultPort, r))
}

func initPort(port string) string {
	_, err := strconv.Atoi(port)
	if err != nil {
		log.Printf("Invalid port %s, using default one: %s", port, defaultPort)
		return defaultPort
	}

	return port
}
