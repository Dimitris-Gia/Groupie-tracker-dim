package main

import (
	"groupie-tracker-dim/handlers"
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("static"))
	// When someone visits /static/, strip that part and serve files from static folder
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/search", handlers.SearchHandler)

	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
