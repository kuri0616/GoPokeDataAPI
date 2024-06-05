package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, "Hello, World!")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	http.HandleFunc("/test", testHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
