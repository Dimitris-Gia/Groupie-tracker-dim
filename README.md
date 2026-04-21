# Groupie Tracker

A web application that displays music artists and bands fetched from the [Groupie Trackers API](https://groupietrackers.herokuapp.com/api). Browse, search, and filter artists, and explore their concert dates and locations.

---

## Features

- Paginated artist grid (6 per page)
- Filter by creation date range, first album year range, and number of members
- Live search with autocomplete suggestions (artist, member, location, date)
- Artist detail page with concert dates, locations, and dates-locations mapping
- Tab navigation on the detail page to isolate a specific section

---

## Requirements

- [Go](https://golang.org/) 1.18 or later
- Internet access (the app fetches data from an external API at runtime)

---

## Getting Started

**Clone the repository**

```bash
git clone <repository-url>
cd groupie-tracker
```

**Run the server**

```bash
go run main.go
```

**Open in your browser**

```
http://localhost:8080
```

---

## Project Structure

```
groupie-tracker/
├── main.go                  # Entry point, route registration
├── api/
│   └── api.go               # External API types and fetch functions
├── handlers/
│   ├── handlers.go          # HTTP handlers (Home, Artist, Search)
│   ├── filters.go           # Server-side filter logic
│   └── pagination.go        # Pagination logic
├── static/
│   ├── styles.css           # Application styles
│   └── main.js              # Client-side search and slider logic
└── templates/
    ├── index.html           # Home page template
    ├── artist.html          # Artist detail page template
    └── error.html           # Error page template
```

---

## Routes

| Method | Path | Description |
|--------|------|-------------|
| GET | `/` | Home page — artist grid with filters and pagination |
| GET | `/artist/{id}` | Artist detail page (all sections) |
| GET | `/artist/{id}?tab=dates` | Artist detail — concert dates only |
| GET | `/artist/{id}?tab=locations` | Artist detail — concert locations only |
| GET | `/artist/{id}?tab=relations` | Artist detail — dates & locations only |
| GET | `/search?q={query}` | Live search — returns JSON |
| GET | `/static/` | Static file server (CSS, JS) |

---

## Filters (Home Page)

All filters are applied server-side via GET query parameters.

| Parameter | Type | Description |
|-----------|------|-------------|
| `yearFrom` | int | Creation date range — lower bound |
| `yearTo` | int | Creation date range — upper bound |
| `albumFrom` | int | First album year range — lower bound |
| `albumTo` | int | First album year range — upper bound |
| `members` | int (repeatable) | Required member count(s), e.g. `members=2&members=4` |
| `page` | int | Page number (default: 1) |

---

## Running Tests

```bash
go test ./...
```

---

## External API

All data is sourced from `https://groupietrackers.herokuapp.com/api`. The application makes no writes to the API.
