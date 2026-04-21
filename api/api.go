// Package api handles all communication with the Groupie Trackers external REST API.
// Base URL: https://groupietrackers.herokuapp.com/api
package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Artist represents a music artist or band returned by the /api/artists endpoint.
type Artist struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"` // URL to the artist's locations endpoint
	ConcertDates string   `json:"concertDates"` // URL to the artist's dates endpoint
	Relations    string   `json:"relations"`    // URL to the artist's relations endpoint
}

// Locations holds the concert locations for a single artist.
type Locations struct {
	Id        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

// AllLocations wraps the index array returned by /api/locations.
type AllLocations struct {
	Index []Locations `json:"index"`
}

// Dates holds the concert dates for a single artist.
type Dates struct {
	Id    int      `json:"id"`
	Dates []string `json:"dates"`
}

// Relations maps each concert location to its list of dates for a single artist.
type Relations struct {
	Id             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

// Api fetches the full list of artists from the given URL.
func Api(url string) ([]Artist, error) {
	var artists []Artist
	err := fetchJSON(url, &artists)
	return artists, err
}

// fetchJSON is a shared helper that performs a GET request and unmarshals the JSON body into target.
func fetchJSON(url string, target interface{}) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, target)
}

// GetArtist fetches a single artist by ID from /api/artists/{id}.
func GetArtist(id int) (*Artist, error) {
	var artist Artist
	err := fetchJSON(fmt.Sprintf("https://groupietrackers.herokuapp.com/api/artists/%d", id), &artist)
	return &artist, err
}

// GetLocations fetches the concert locations for a single artist by ID.
func GetLocations(id int) (*Locations, error) {
	var loc Locations
	err := fetchJSON(fmt.Sprintf("https://groupietrackers.herokuapp.com/api/locations/%d", id), &loc)
	return &loc, err
}

// GetDates fetches the concert dates for a single artist by ID.
func GetDates(id int) (*Dates, error) {
	var dates Dates
	err := fetchJSON(fmt.Sprintf("https://groupietrackers.herokuapp.com/api/dates/%d", id), &dates)
	return &dates, err
}

// GetRelations fetches the dates-locations mapping for a single artist by ID.
func GetRelations(id int) (*Relations, error) {
	var rel Relations
	err := fetchJSON(fmt.Sprintf("https://groupietrackers.herokuapp.com/api/relation/%d", id), &rel)
	return &rel, err
}

// GetAllLocations fetches the full locations list for all artists from /api/locations.
func GetAllLocations() ([]Locations, error) {
	var all AllLocations
	err := fetchJSON("https://groupietrackers.herokuapp.com/api/locations", &all)
	return all.Index, err
}
