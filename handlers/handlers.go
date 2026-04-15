package handlers

import (
	"net/http"
	"strconv"
	"text/template"

	"groupie-tracker-dim/api"
)

func renderError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	tmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		w.Write([]byte(message))
		return
	}
	data := struct {
		Code    int
		Message string
	}{Code: status, Message: message}

	// If template execution fails, we've already called WriteHeader
	// So just write the error message
	if err := tmpl.Execute(w, data); err != nil {
		w.Write([]byte(message))
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		renderError(w, http.StatusNotFound, "Page not found")
		return
	}
	tmpl, err := template.New("index.html").Funcs(template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"gt":  func(a, b int) bool { return a > b },
		"lt":  func(a, b int) bool { return a < b },
		"seq": func(start, end int) []int {
			// Generate a slice of integers from start to end
			nums := make([]int, end-start+1)
			for i := range nums {
				nums[i] = start + i
			}
			return nums
		},
	}).ParseFiles("templates/index.html")
	if err != nil {
		renderError(w, http.StatusNotFound, "Template not found")
		return
	}
	url := "https://groupietrackers.herokuapp.com/api/artists"
	Artists, err := api.Api(url)
	if err != nil {
		renderError(w, http.StatusNotFound, "Api not found")
		return
	}

	// filter := r.URL.Query().Get("filter")
	year := r.URL.Query().Get("year")
	members := r.URL.Query()["members"]
	// album := r.URL.Query().Get("album")
	// location := r.URL.Query().Get("location")

	filtered := []api.Artist{}

	for _, artist := range Artists {
		if year != "" {
			y, _ := strconv.Atoi(year)
			if artist.CreationDate < y {
				continue
			}
		}
		if len(members) > 0 {
			match := false

			for _, mStr := range members {
				m, err := strconv.Atoi(mStr)
				if err == nil && len(artist.Members) == m {
					match = true
					break
				}
			}

			if !match {
				continue
			}
		}
		filtered = append(filtered, artist)
	}

	data := Pagination(6, filtered, r)

	err = tmpl.Execute(w, data)
	if err != nil {
		// Return 500 if template execution fails
		renderError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
}
