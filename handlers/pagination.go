package handlers

import (
	"net/http"
	"strconv"

	"groupie-tracker-dim/api"
)

func Pagination(ItemsPerPage int, artists []api.Artist, r *http.Request) struct {
	Artists    []api.Artist
	Page       int
	TotalPages int
	HasNext    bool
	HasPrev    bool
	Filter     string
	Year       string
	Members    []string
	Album      string
} {
	// 👉 Get page from URL
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	perPage := ItemsPerPage
	start := (page - 1) * perPage
	end := start + perPage

	if start > len(artists) {
		start = len(artists)
	}
	if end > len(artists) {
		end = len(artists)
	}

	paginatedArtists := artists[start:end]

	// Calculate total pages
	totalPages := (len(artists) + perPage - 1) / perPage

	// 👉 Send extra data to template
	data := struct {
		Artists    []api.Artist
		Page       int
		TotalPages int
		HasNext    bool
		HasPrev    bool
		Filter     string
		Year       string
		Members    []string
		Album      string
	}{
		Artists:    paginatedArtists,
		Page:       page,
		TotalPages: totalPages,
		HasNext:    end < len(artists),
		HasPrev:    start > 0,
		Filter:     r.URL.Query().Get("filter"),
		Year:       r.URL.Query().Get("year"),
		Members:    r.URL.Query()["members"],
		Album:      r.URL.Query().Get("album"),
	}
	return data
}
