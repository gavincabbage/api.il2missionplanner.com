package models

type Plan struct {
	// Since multiple front lines can be present in a mission area, a quaqrupally nested list is necessary.
	// The field is an array of front lines. A front line is an array of two lines (gray and red).
	// A line is an array of points. Finally, a point is an array of two floats representing its coordinates.
	Frontline [][][][]float64 `json:"frontline"`
	MapHash   string          `json:"mapHash"`
	Routes    []Route         `json:"routes"`
	Points    []Point         `json:"points"`
}

type Route struct {
	LatLngs []LatLng `json:"latLngs"`
	Name    string   `json:"name"`
	Owner   string   `json:"owner"`
	Color   string   `json:"color"`
	Speed   int      `json:"speed"`
	Speeds  []int    `json:"speeds"`
}

type Point struct {
	LatLng LatLng `json:"latLng"`
	Name   string `json:"name"`
	Owner  string `json:"owner"`
	Color  string `json:"color"`
	Type   string `json:"type"`
	Notes  string `json:"notes"`
}

type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
