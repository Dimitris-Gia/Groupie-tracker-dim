package handlers

import (
	"groupie-tracker-dim/api"
	"net/http"
	"strconv"
)

func Pagination(ItemsPerPage int, artists []api.Artist, r *http.Request) struct {
	Artists    []api.Artist
	Page       int
	TotalPages int
	HasNext    bool
	HasPrev    bool
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
	}{
		Artists:    paginatedArtists,
		Page:       page,
		TotalPages: totalPages,
		HasNext:    end < len(artists),
		HasPrev:    start > 0,
	}
	return data

}
