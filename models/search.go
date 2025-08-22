package models

// SearchResponse represents the response structure for search results
type SearchResponse struct {
	ConfidenceScore float64        `json:"confidence_score"`
	Message         string         `json:"message"`
	Source          string         `json:"source"`
	Data            []SearchDetail `json:"data"`
}

// SearchDetail represents each search result item
type SearchDetail struct {
	Judul    string   `json:"judul"`
	URL      string   `json:"url"`
	Slug     string   `json:"anime_slug"`
	Status   string   `json:"status"`
	Tipe     string   `json:"tipe"`
	Skor     string   `json:"skor"`
	Penonton string   `json:"penonton"`
	Sinopsis string   `json:"sinopsis"`
	Genre    []string `json:"genre"`
	Cover    string   `json:"cover"`
}
