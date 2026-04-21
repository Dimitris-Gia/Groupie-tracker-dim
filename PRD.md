# Product Requirements Document — Groupie Tracker

## 1. Overview

Groupie Tracker is a web application that consumes a public REST API to display information about music artists and bands. Users can browse, search, and filter artists, and navigate to a dedicated detail page for each artist showing their concert dates, locations, and dates-locations mapping.

---

## 2. Goals

- Display a browsable, paginated grid of music artists fetched from an external API.
- Allow users to filter artists by creation date range, first album year range, and number of members.
- Provide a live search with autocomplete suggestions across artist names, members, locations, first album dates, and creation dates.
- Provide a dedicated artist detail page showing all related data (dates, locations, relations).
- Allow users to view only a specific section of the detail page (dates, locations, or relations) via a tab parameter.

---

## 3. Users

Single user type: anonymous visitor. No authentication or user accounts are required.

---

## 4. Functional Requirements

### 4.1 Home Page (`/`)

| ID | Requirement |
|----|-------------|
| F1 | Fetch and display all artists from `https://groupietrackers.herokuapp.com/api/artists`. |
| F2 | Display artists in a responsive 3-column card grid. Each card shows: name (linked), image (linked), first album date, creation date, members list, and links to concert dates, locations, and dates+locations. |
| F3 | Paginate the grid at 6 artists per page. Show Previous / Next controls and numbered page links. Pagination links must preserve all active filter values. |
| F4 | Filter by creation date range (`yearFrom` / `yearTo`) using a dual range slider with synced number inputs. |
| F5 | Filter by first album year range (`albumFrom` / `albumTo`) using a dual range slider with synced number inputs. |
| F6 | Filter by number of members using checkboxes (1–7). Multiple values may be selected simultaneously. |
| F7 | All filters are applied server-side on form submission via GET parameters. |

### 4.2 Live Search (`/search`)

| ID | Requirement |
|----|-------------|
| F8 | A search input in the header triggers a fetch to `/search?q=...` on every keystroke. |
| F9 | Results replace the card grid without a page reload. Pagination is hidden during search. |
| F10 | Clearing the search input restores the original server-rendered grid and pagination. |
| F11 | An autocomplete dropdown shows deduplicated suggestions with their match type (artist/band, member, location, first album date, creation date). |
| F12 | Clicking a suggestion fills the input and re-runs the search. Clicking outside the search box closes the dropdown. |
| F13 | Search matches against: artist name, member names, concert locations, first album date string, and creation date year. |

### 4.3 Artist Detail Page (`/artist/{id}`)

| ID | Requirement |
|----|-------------|
| F14 | Fetch and display the artist's full profile: image, name, creation date, first album, members list. |
| F15 | Display concert dates, concert locations, and dates-locations mapping. |
| F16 | An optional `?tab=dates`, `?tab=locations`, or `?tab=relations` query parameter shows only the corresponding section; all three sections are shown when no tab is specified. |
| F17 | The left panel contains navigation links to each tab and an "All Info" link. |
| F18 | Clicking the artist name or image on any card navigates to the detail page with no tab (all info visible). |

### 4.4 Error Handling

| ID | Requirement |
|----|-------------|
| F19 | Any unknown path returns a 404 page rendered with `error.html`. |
| F20 | Invalid artist IDs (non-numeric or < 1) return a 400 response. |
| F21 | API or template failures return appropriate 500 responses. |

---

## 5. Non-Functional Requirements

| ID | Requirement |
|----|-------------|
| NF1 | The server must start and be ready to serve requests within 2 seconds on a standard development machine. |
| NF2 | All pages must render correctly in modern browsers (Chrome, Firefox, Safari, Edge). |
| NF3 | Static assets (CSS, JS) are served from the `/static/` directory with path prefix stripping. |
| NF4 | No external Go dependencies beyond the standard library. |
| NF5 | The application must handle API unavailability gracefully by returning an error page rather than panicking. |

---

## 6. External API

Base URL: `https://groupietrackers.herokuapp.com/api`

| Endpoint | Used for |
|----------|----------|
| `GET /artists` | Full artist list |
| `GET /artists/{id}` | Single artist |
| `GET /locations` | All concert locations |
| `GET /locations/{id}` | Locations for one artist |
| `GET /dates/{id}` | Concert dates for one artist |
| `GET /relation/{id}` | Dates-locations map for one artist |

---

## 7. Project Structure

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

## 8. Out of Scope

- User authentication or accounts
- Favouriting or saving artists
- Server-side caching of API responses
- Mobile-specific layouts
