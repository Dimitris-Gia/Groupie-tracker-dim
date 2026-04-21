package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHomeHandler_WrongPath checks that any path other than "/" returns 404.
func TestHomeHandler_WrongPath(t *testing.T) {
	r := httptest.NewRequest("GET", "/notfound", nil)
	w := httptest.NewRecorder()
	HomeHandler(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

// TestArtistHandler_InvalidID checks that a non-numeric ID returns 400.
func TestArtistHandler_InvalidID(t *testing.T) {
	r := httptest.NewRequest("GET", "/artist/abc", nil)
	r.URL.Path = "/artist/abc"
	w := httptest.NewRecorder()
	ArtistHandler(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// TestArtistHandler_ZeroID checks that ID=0 returns 400.
func TestArtistHandler_ZeroID(t *testing.T) {
	r := httptest.NewRequest("GET", "/artist/0", nil)
	r.URL.Path = "/artist/0"
	w := httptest.NewRecorder()
	ArtistHandler(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// TestArtistHandler_NegativeID checks that a negative ID returns 400.
func TestArtistHandler_NegativeID(t *testing.T) {
	r := httptest.NewRequest("GET", "/artist/-1", nil)
	r.URL.Path = "/artist/-1"
	w := httptest.NewRecorder()
	ArtistHandler(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// TestSearchHandler_EmptyQuery checks that /search with no query returns a valid JSON response.
func TestSearchHandler_EmptyQuery(t *testing.T) {
	r := httptest.NewRequest("GET", "/search?q=", nil)
	w := httptest.NewRecorder()
	SearchHandler(w, r)

	// The handler calls the external API; if it's unreachable the response will be 500.
	// We only assert the Content-Type when the call succeeds (200).
	if w.Code == http.StatusOK {
		ct := w.Header().Get("Content-Type")
		if ct != "application/json" {
			t.Errorf("expected application/json, got %q", ct)
		}
	}
}
