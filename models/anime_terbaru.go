package models

// OngoingDramaResponse represents the response structure for ongoing drama list
type OngoingDramaResponse struct {
	ConfidenceScore float64      `json:"confidence_score"`
	Message         string       `json:"message"`
	Source          string       `json:"source"`
	Data            []DramaEntry `json:"data"`
}

// DramaEntry represents each drama item in the list
type DramaEntry struct {
	Judul    string `json:"judul"`
	URL      string `json:"url"`
	Slug     string `json:"anime_slug"`
	Episode  string `json:"episode"`
	Uploader string `json:"uploader"`
	Rilis    string `json:"rilis"`
	Cover    string `json:"cover"`
}
