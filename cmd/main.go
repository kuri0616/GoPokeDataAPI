package main

import (
	"github.com/gorilla/mux"
	"github.com/rikuya98/go-poke-data-api/handlers"
	"log"
	"net/http"
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/pokemon/{id:[0-9]+}", handlers.GetPokeDataHandler).Methods(http.MethodGet)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
