package models

// DramaListResponse represents the response structure for drama list
type DramaListResponse struct {
	ConfidenceScore float64       `json:"confidence_score"`
	Message         string        `json:"message"`
	Source          string        `json:"source"`
	Data            []DramaDetail `json:"data"`
}

// DramaDetail represents each drama item in the list
type DramaDetail struct {
	Judul    string   `json:"judul"`
	URL      string   `json:"url"`
	Slug     string   `json:"anime_slug"`
	Status   string   `json:"status"`
	Skor     string   `json:"skor"`
	Sinopsis string   `json:"sinopsis"`
	Views    string   `json:"views"`
	Cover    string   `json:"cover"`
	Genres   []string `json:"genres"`
	Tanggal  string   `json:"tanggal"`
}
