package main

import (
	"groupie-tracker-dim/handlers"
	"net/http"
)

func main() {
	http.HandleFunc("/", handlers.HomeHandler)
}
