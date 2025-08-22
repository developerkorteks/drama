package models

// ReleaseScheduleResponse represents the response structure for release schedule
type ReleaseScheduleResponse struct {
	ConfidenceScore float64                   `json:"confidence_score"`
	Message         string                    `json:"message"`
	Source          string                    `json:"source"`
	Data            map[string][]ReleaseEntry `json:"data"`
}

// ReleaseEntry represents each release item in the schedule
type ReleaseEntry struct {
	Title       string   `json:"title"`
	URL         string   `json:"url"`
	Slug        string   `json:"anime_slug"`
	CoverURL    string   `json:"cover_url"`
	Type        string   `json:"type"`
	Score       string   `json:"score"`
	Genres      []string `json:"genres"`
	ReleaseTime string   `json:"release_time"`
}

// ScheduleByDayResponse represents the response structure for schedule by specific day
type ScheduleByDayResponse struct {
	ConfidenceScore float64         `json:"confidence_score"`
	Message         string          `json:"message"`
	Source          string          `json:"source"`
	Data            []ScheduleEntry `json:"data"`
}

// ScheduleEntry represents each schedule item for specific day
type ScheduleEntry struct {
	Title       string   `json:"title"`
	URL         string   `json:"url"`
	Slug        string   `json:"anime_slug"`
	CoverURL    string   `json:"cover_url"`
	Type        string   `json:"type"`
	Score       string   `json:"score"`
	Genres      []string `json:"genres"`
	ReleaseTime string   `json:"release_time"`
}
