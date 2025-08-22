package models

// DetailResponse represents the response structure for anime/drama detail
type DetailResponse struct {
	ConfidenceScore float64              `json:"confidence_score"`
	Message         string               `json:"message"`
	Source          string               `json:"source"`
	Judul           string               `json:"judul"`
	URL             string               `json:"url"`
	AnimeSlug       string               `json:"anime_slug"`
	Cover           string               `json:"cover"`
	EpisodeList     []EpisodeItem        `json:"episode_list"`
	Recommendations []RecommendationItem `json:"recommendations"`
	Status          string               `json:"status"`
	Tipe            string               `json:"tipe"`
	Skor            string               `json:"skor"`
	Penonton        string               `json:"penonton"`
	Sinopsis        string               `json:"sinopsis"`
	Genre           []string             `json:"genre"`
	Details         DetailsObject        `json:"details"`
	Rating          RatingObject         `json:"rating"`
}

// EpisodeItem represents each episode in the episode list
type EpisodeItem struct {
	Episode     string `json:"episode"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	EpisodeSlug string `json:"episode_slug"`
	ReleaseDate string `json:"release_date"`
}

// RecommendationItem represents each recommendation item
type RecommendationItem struct {
	Title     string `json:"title"`
	URL       string `json:"url"`
	AnimeSlug string `json:"anime_slug"`
	CoverURL  string `json:"cover_url"`
	Rating    string `json:"rating"`
	Episode   string `json:"episode"`
}

// DetailsObject represents detailed information about the anime/drama
type DetailsObject struct {
	Japanese     string `json:"Japanese"`
	English      string `json:"English"`
	Status       string `json:"Status"`
	Type         string `json:"Type"`
	Source       string `json:"Source"`
	Duration     string `json:"Duration"`
	TotalEpisode string `json:"Total Episode"`
	Season       string `json:"Season"`
	Studio       string `json:"Studio"`
	Producers    string `json:"Producers"`
	Released     string `json:"Released:"`
}

// RatingObject represents rating information
type RatingObject struct {
	Score string `json:"score"`
	Users string `json:"users"`
}
