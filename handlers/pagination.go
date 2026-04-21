package handlers

import (
	"net/http"
	"strconv"

	"groupie-tracker/api"
)

// Pagination slices the artist list into a single page and returns a struct
// containing the page's artists plus all metadata needed by the template
// to render pagination controls and preserve active filter state in links.
//
// Parameters:
//   - ItemsPerPage : number of artists to show per page
//   - artists      : the already-filtered list of artists
//   - allLocations : sorted list of unique location strings for the dropdown
//   - r            : the incoming request (used to read the "page" query param and filter values)
func Pagination(ItemsPerPage int, artists []api.Artist, allLocations []string, r *http.Request) struct {
	Artists      []api.Artist
	Page         int
	TotalPages   int
	HasNext      bool
	HasPrev      bool
	Filter       string
	YearFrom     string
	YearTo       string
	Members      []string
	AlbumFrom    string
	AlbumTo      string
	AllLocations []string
	Location     string
} {
	// Default to page 1 if the param is missing or invalid
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	perPage := ItemsPerPage
	start := (page - 1) * perPage
	end := start + perPage

	// Clamp slice bounds to the actual list length
	if start > len(artists) {
		start = len(artists)
	}
	if end > len(artists) {
		end = len(artists)
	}

	paginatedArtists := artists[start:end]

	// Integer ceiling division to get total page count
	totalPages := (len(artists) + perPage - 1) / perPage

	data := struct {
		Artists      []api.Artist
		Page         int
		TotalPages   int
		HasNext      bool
		HasPrev      bool
		Filter       string
		YearFrom     string
		YearTo       string
		Members      []string
		AlbumFrom    string
		AlbumTo      string
		AllLocations []string
		Location     string
	}{
		Artists:    paginatedArtists,
		Page:       page,
		TotalPages: totalPages,
		HasNext:    end < len(artists),
		HasPrev:    start > 0,
		// Preserve all active filter values so pagination links keep the current filter state
		Filter:       r.URL.Query().Get("filter"),
		YearFrom:     r.URL.Query().Get("yearFrom"),
		YearTo:       r.URL.Query().Get("yearTo"),
		Members:      r.URL.Query()["members"],
		AlbumFrom:    r.URL.Query().Get("albumFrom"),
		AlbumTo:      r.URL.Query().Get("albumTo"),
		AllLocations: allLocations,
		Location:     r.URL.Query().Get("location"),
	}
	return data
}
