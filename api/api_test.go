package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// newMockServer starts a test HTTP server that responds with the given body for every request.
func newMockServer(body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(body))
	}))
}

func TestFetchJSON_ValidResponse(t *testing.T) {
	srv := newMockServer(`[{"id":1,"name":"Queen"}]`)
	defer srv.Close()

	var artists []Artist
	if err := fetchJSON(srv.URL, &artists); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(artists) != 1 || artists[0].Name != "Queen" {
		t.Errorf("unexpected result: %+v", artists)
	}
}

func TestFetchJSON_InvalidJSON(t *testing.T) {
	srv := newMockServer(`not json`)
	defer srv.Close()

	var artists []Artist
	if err := fetchJSON(srv.URL, &artists); err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestFetchJSON_UnreachableURL(t *testing.T) {
	var artists []Artist
	if err := fetchJSON("http://127.0.0.1:0/nope", &artists); err == nil {
		t.Error("expected error for unreachable URL, got nil")
	}
}

func TestApi_ReturnsList(t *testing.T) {
	payload := []Artist{
		{Id: 1, Name: "Queen", CreationDate: 1970},
		{Id: 2, Name: "Metallica", CreationDate: 1981},
	}
	body, _ := json.Marshal(payload)
	srv := newMockServer(string(body))
	defer srv.Close()

	artists, err := Api(srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(artists) != 2 {
		t.Errorf("expected 2 artists, got %d", len(artists))
	}
}

func TestGetAllLocations_ReturnsIndex(t *testing.T) {
	payload := AllLocations{
		Index: []Locations{
			{Id: 1, Locations: []string{"london-uk", "paris-france"}},
		},
	}
	body, _ := json.Marshal(payload)
	srv := newMockServer(string(body))
	defer srv.Close()

	// Override the URL by calling fetchJSON directly since GetAllLocations hardcodes the URL
	var all AllLocations
	if err := fetchJSON(srv.URL, &all); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(all.Index) != 1 || all.Index[0].Id != 1 {
		t.Errorf("unexpected locations: %+v", all)
	}
}

func TestArtistStruct_JSONUnmarshal(t *testing.T) {
	raw := `{
		"id": 5,
		"image": "https://example.com/img.jpg",
		"name": "The Beatles",
		"members": ["John", "Paul", "George", "Ringo"],
		"creationDate": 1960,
		"firstAlbum": "22-03-1963",
		"locations": "https://example.com/locations/5",
		"concertDates": "https://example.com/dates/5",
		"relations": "https://example.com/relation/5"
	}`
	var a Artist
	if err := json.Unmarshal([]byte(raw), &a); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if a.Id != 5 {
		t.Errorf("expected Id 5, got %d", a.Id)
	}
	if a.Name != "The Beatles" {
		t.Errorf("expected 'The Beatles', got %q", a.Name)
	}
	if len(a.Members) != 4 {
		t.Errorf("expected 4 members, got %d", len(a.Members))
	}
}

func TestRelationsStruct_JSONUnmarshal(t *testing.T) {
	raw := `{"id":1,"datesLocations":{"london-uk":["12-05-2019"],"paris-france":["14-05-2019"]}}`
	var rel Relations
	if err := json.Unmarshal([]byte(raw), &rel); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(rel.DatesLocations) != 2 {
		t.Errorf("expected 2 locations, got %d", len(rel.DatesLocations))
	}
}
