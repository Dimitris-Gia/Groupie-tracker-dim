package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"groupie-tracker/api"
)

// Filters applies the active query-string filters to the full artist list and returns
// only the artists that match all provided criteria.
// Supported filters:
//   - yearFrom / yearTo   : creation date range
//   - albumFrom / albumTo : first album year range
//   - members             : one or more required member counts (multi-value param)
//   - location            : substring match against the artist's concert locations
//
// locMap maps artist ID → list of raw location strings from the API.
func Filters(Artists []api.Artist, locMap map[int][]string, r *http.Request) (filtered []api.Artist) {
	yearFrom := r.URL.Query().Get("yearFrom")
	yearTo := r.URL.Query().Get("yearTo")
	members := r.URL.Query()["members"]
	albumFrom := r.URL.Query().Get("albumFrom")
	albumTo := r.URL.Query().Get("albumTo")
	// location is the value selected in the dropdown, already in display form (spaces, not underscores)
	location := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("location")))

	for _, artist := range Artists {
		// --- Creation date range filter ---
		if yearFrom != "" && yearTo != "" {
			fromCreYear, err1 := strconv.Atoi(yearFrom)
			toCreYear, err2 := strconv.Atoi(yearTo)
			if err1 != nil || err2 != nil {
				continue
			}
			// Swap if the user entered the range backwards
			if fromCreYear > toCreYear {
				fromCreYear, toCreYear = toCreYear, fromCreYear
			}
			if artist.CreationDate < fromCreYear || artist.CreationDate > toCreYear {
				continue
			}
		}

		// --- Member count filter ---
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

		// --- Album year range filter ---
		if albumFrom != "" && albumTo != "" {
			fromYear, err1 := strconv.Atoi(albumFrom)
			toYear, err2 := strconv.Atoi(albumTo)
			if err1 != nil || err2 != nil {
				continue
			}
			if fromYear > toYear {
				fromYear, toYear = toYear, fromYear
			}
			if len(artist.FirstAlbum) < 4 {
				continue
			}
			yearStr := artist.FirstAlbum[len(artist.FirstAlbum)-4:]
			a, err := strconv.Atoi(yearStr)
			if err != nil || a < fromYear || a > toYear {
				continue
			}
		}

		// --- Location filter ---
		// Uses substring matching so that e.g. "washington" matches "seattle-washington-usa".
		// The hint states Seattle, Washington, USA is part of Washington, USA — substring
		// matching on the full location string handles this naturally.
		if location != "" {
			locs, ok := locMap[artist.Id]
			if !ok {
				continue
			}
			matched := false
			for _, l := range locs {
				// Normalise: replace underscores with spaces before comparing
				if strings.Contains(strings.ToLower(strings.ReplaceAll(l, "_", " ")), location) {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}

		filtered = append(filtered, artist)
	}

	return filtered
}
