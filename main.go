package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello world")
	})
	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
