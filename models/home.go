package models

// FinalResponse represents the complete API response structure
type FinalResponse struct {
	ConfidenceScore float64      `json:"confidence_score"`
	Message         string       `json:"message"`
	Source          string       `json:"source"`
	Top10           []Top10Item  `json:"top10"`
	NewEps          []NewEpsItem `json:"new_eps"`
	Movies          []MovieItem  `json:"movies"`
	JadwalRilis     JadwalRilis  `json:"jadwal_rilis"`
}

// Top10Item represents a top 10 drama item
type Top10Item struct {
	Judul     string   `json:"judul"`
	URL       string   `json:"url"`
	AnimeSlug string   `json:"anime_slug"`
	Rating    string   `json:"rating"`
	Cover     string   `json:"cover"`
	Genres    []string `json:"genres"`
}

// NewEpsItem represents a new episode item
type NewEpsItem struct {
	Judul     string `json:"judul"`
	URL       string `json:"url"`
	AnimeSlug string `json:"anime_slug"`
	Episode   string `json:"episode"`
	Rilis     string `json:"rilis"`
	Cover     string `json:"cover"`
}

// MovieItem represents a movie item
type MovieItem struct {
	Judul     string   `json:"judul"`
	URL       string   `json:"url"`
	AnimeSlug string   `json:"anime_slug"`
	Tanggal   string   `json:"tanggal"`
	Cover     string   `json:"cover"`
	Genres    []string `json:"genres"`
}

// JadwalRilis represents the release schedule for all days
type JadwalRilis struct {
	Monday    []JadwalItem `json:"Monday"`
	Tuesday   []JadwalItem `json:"Tuesday"`
	Wednesday []JadwalItem `json:"Wednesday"`
	Thursday  []JadwalItem `json:"Thursday"`
	Friday    []JadwalItem `json:"Friday"`
	Saturday  []JadwalItem `json:"Saturday"`
	Sunday    []JadwalItem `json:"Sunday"`
}

// JadwalItem represents a schedule item
type JadwalItem struct {
	Title       string   `json:"title"`
	URL         string   `json:"url"`
	AnimeSlug   string   `json:"anime_slug"`
	CoverURL    string   `json:"cover_url"`
	Type        string   `json:"type"`
	Score       string   `json:"score"`
	Genres      []string `json:"genres"`
	ReleaseTime string   `json:"release_time"`
}
