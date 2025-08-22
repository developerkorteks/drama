package models

// EpisodeDetailResponse represents the response structure for episode detail
type EpisodeDetailResponse struct {
	ConfidenceScore  float64           `json:"confidence_score"`
	Message          string            `json:"message"`
	Source           string            `json:"source"`
	Title            string            `json:"title"`
	ThumbnailURL     string            `json:"thumbnail_url"`
	StreamingServers []StreamingServer `json:"streaming_servers"`
	ReleaseInfo      string            `json:"release_info"`
	DownloadLinks    DownloadLinks     `json:"download_links"`
	Navigation       Navigation        `json:"navigation"`
	AnimeInfo        AnimeInfo         `json:"anime_info"`
	OtherEpisodes    []OtherEpisode    `json:"other_episodes"`
}

// StreamingServer represents each streaming server
type StreamingServer struct {
	ServerName   string `json:"server_name"`
	StreamingURL string `json:"streaming_url"`
}

// DownloadLinks represents download links organized by format and quality
type DownloadLinks struct {
	MKV  map[string][]DownloadProvider `json:"MKV"`
	MP4  map[string][]DownloadProvider `json:"MP4"`
	X265 map[string][]DownloadProvider `json:"x265 [Mode Irit Kuota tapi Kualitas Sama Beningnya]"`
}

// DownloadProvider represents each download provider
type DownloadProvider struct {
	Provider string `json:"provider"`
	URL      string `json:"url"`
}

// Navigation represents episode navigation links
type Navigation struct {
	PreviousEpisodeURL string `json:"previous_episode_url,omitempty"`
	AllEpisodesURL     string `json:"all_episodes_url"`
	NextEpisodeURL     string `json:"next_episode_url,omitempty"`
}

// AnimeInfo represents anime information
type AnimeInfo struct {
	Title        string   `json:"title"`
	ThumbnailURL string   `json:"thumbnail_url"`
	Synopsis     string   `json:"synopsis"`
	Genres       []string `json:"genres"`
}

// OtherEpisode represents other episodes
type OtherEpisode struct {
	Title        string `json:"title"`
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnail_url"`
	ReleaseDate  string `json:"release_date"`
}

// AjaxPlayerResponse represents AJAX response structure
type AjaxPlayerResponse struct {
	Success bool `json:"success"`
	Data    struct {
		IframeURL string `json:"iframe_url"`
	} `json:"data"`
}
