package handlers

import (
	"net/http"
	"testing"

	"groupie-tracker/api"
)

// makeArtists generates a slice of n dummy artists for pagination tests.
func makeArtists(n int) []api.Artist {
	artists := make([]api.Artist, n)
	for i := range artists {
		artists[i] = api.Artist{Id: i + 1, Name: "Artist"}
	}
	return artists
}

var noLocs []string // empty locations list — not relevant for pagination logic

func TestPagination_FirstPage(t *testing.T) {
	artists := makeArtists(15)
	r, _ := http.NewRequest("GET", "/?page=1", nil)
	data := Pagination(6, artists, noLocs, r)

	if data.Page != 1 {
		t.Errorf("expected page 1, got %d", data.Page)
	}
	if len(data.Artists) != 6 {
		t.Errorf("expected 6 artists, got %d", len(data.Artists))
	}
	if data.HasPrev {
		t.Error("expected HasPrev=false on first page")
	}
	if !data.HasNext {
		t.Error("expected HasNext=true on first page")
	}
	if data.TotalPages != 3 {
		t.Errorf("expected 3 total pages, got %d", data.TotalPages)
	}
}

func TestPagination_LastPage(t *testing.T) {
	artists := makeArtists(15)
	r, _ := http.NewRequest("GET", "/?page=3", nil)
	data := Pagination(6, artists, noLocs, r)

	if data.Page != 3 {
		t.Errorf("expected page 3, got %d", data.Page)
	}
	if len(data.Artists) != 3 {
		t.Errorf("expected 3 artists on last page, got %d", len(data.Artists))
	}
	if !data.HasPrev {
		t.Error("expected HasPrev=true on last page")
	}
	if data.HasNext {
		t.Error("expected HasNext=false on last page")
	}
}

func TestPagination_MiddlePage(t *testing.T) {
	artists := makeArtists(15)
	r, _ := http.NewRequest("GET", "/?page=2", nil)
	data := Pagination(6, artists, noLocs, r)

	if !data.HasPrev || !data.HasNext {
		t.Error("expected both HasPrev and HasNext on middle page")
	}
	if len(data.Artists) != 6 {
		t.Errorf("expected 6 artists, got %d", len(data.Artists))
	}
}

func TestPagination_InvalidPageDefaultsToOne(t *testing.T) {
	artists := makeArtists(10)
	r, _ := http.NewRequest("GET", "/?page=abc", nil)
	data := Pagination(6, artists, noLocs, r)

	if data.Page != 1 {
		t.Errorf("expected page 1 for invalid param, got %d", data.Page)
	}
}

func TestPagination_PageBeyondTotal_ClampsToEmpty(t *testing.T) {
	artists := makeArtists(6)
	r, _ := http.NewRequest("GET", "/?page=99", nil)
	data := Pagination(6, artists, noLocs, r)

	if len(data.Artists) != 0 {
		t.Errorf("expected 0 artists for out-of-range page, got %d", len(data.Artists))
	}
}

func TestPagination_ExactlyOneFullPage(t *testing.T) {
	artists := makeArtists(6)
	r, _ := http.NewRequest("GET", "/", nil)
	data := Pagination(6, artists, noLocs, r)

	if data.TotalPages != 1 {
		t.Errorf("expected 1 total page, got %d", data.TotalPages)
	}
	if data.HasNext || data.HasPrev {
		t.Error("expected no navigation on a single full page")
	}
}

func TestPagination_EmptyList(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	data := Pagination(6, []api.Artist{}, noLocs, r)

	if len(data.Artists) != 0 {
		t.Errorf("expected 0 artists, got %d", len(data.Artists))
	}
	if data.TotalPages != 0 {
		t.Errorf("expected 0 total pages, got %d", data.TotalPages)
	}
}

func TestPagination_PreservesFilterParams(t *testing.T) {
	artists := makeArtists(6)
	r, _ := http.NewRequest("GET", "/?yearFrom=1970&yearTo=2000&albumFrom=1975&albumTo=1995&members=4&location=london-uk", nil)
	data := Pagination(6, artists, noLocs, r)

	if data.YearFrom != "1970" || data.YearTo != "2000" {
		t.Errorf("year range not preserved: %s %s", data.YearFrom, data.YearTo)
	}
	if data.AlbumFrom != "1975" || data.AlbumTo != "1995" {
		t.Errorf("album range not preserved: %s %s", data.AlbumFrom, data.AlbumTo)
	}
	if len(data.Members) != 1 || data.Members[0] != "4" {
		t.Errorf("members not preserved: %v", data.Members)
	}
	if data.Location != "london-uk" {
		t.Errorf("location not preserved: %s", data.Location)
	}
}

func TestPagination_AllLocationsPassedThrough(t *testing.T) {
	artists := makeArtists(3)
	locs := []string{"berlin-germany", "london-uk", "paris-france"}
	r, _ := http.NewRequest("GET", "/", nil)
	data := Pagination(6, artists, locs, r)

	if len(data.AllLocations) != 3 {
		t.Errorf("expected 3 locations, got %d", len(data.AllLocations))
	}
}
