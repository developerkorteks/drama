package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/nabilulilalbab/dramaqu/models"
)

type EpisodeDetailService struct{}

func NewEpisodeDetailService() *EpisodeDetailService {
	return &EpisodeDetailService{}
}

// GetEpisodeDetail scrapes and returns episode detail with the exact same logic as the test
func (s *EpisodeDetailService) GetEpisodeDetail(episodeURL string) (*models.EpisodeDetailResponse, error) {
	episodeResponse := &models.EpisodeDetailResponse{
		ConfidenceScore: 1.0,
		Message:         "Success",
		Source:          "dramaqu.ad",
		ReleaseInfo:     fmt.Sprintf("Released on %s %d", time.Now().Month().String(), time.Now().Year()),
		DownloadLinks: models.DownloadLinks{
			MKV:  make(map[string][]models.DownloadProvider),
			MP4:  make(map[string][]models.DownloadProvider),
			X265: make(map[string][]models.DownloadProvider),
		},
		AnimeInfo:        models.AnimeInfo{},
		StreamingServers: []models.StreamingServer{}, // Initialize slice to avoid null
		OtherEpisodes:    []models.OtherEpisode{},
	}

	c := colly.NewCollector(
		colly.AllowedDomains("dramaqu.ad"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"),
	)
	c.SetRequestTimeout(30 * time.Second)

	c.OnHTML("body", func(e *colly.HTMLElement) {
		log.Println("Mem-parsing HTML dari halaman utama...")
		doc := e.DOM

		// Parsing data statis
		mainContent := doc.Find("div.single-content.movie")
		episodeResponse.Title = mainContent.Find("div.title span").Text() + " " + mainContent.Find("div.release").Text()
		thumbnail := mainContent.Find("div.poster img").AttrOr("src", "")
		episodeResponse.ThumbnailURL = thumbnail
		episodeResponse.AnimeInfo.ThumbnailURL = thumbnail
		episodeResponse.AnimeInfo.Synopsis = strings.TrimSpace(mainContent.Find("div.excerpt").Text())
		mainContent.Find("div.categories a").Each(func(i int, s *goquery.Selection) {
			episodeResponse.AnimeInfo.Genres = append(episodeResponse.AnimeInfo.Genres, s.Text())
		})

		baseDramaURL, currentEpNum := "", "1"
		reEp := regexp.MustCompile(`(https?://[^/]+/[^/]+)/(\d+)/?$`)
		reBase := regexp.MustCompile(`(https?://[^/]+/[^/]+)/?$`)
		if matches := reEp.FindStringSubmatch(episodeURL); len(matches) > 2 {
			baseDramaURL, currentEpNum = matches[1], matches[2]
		} else if matches := reBase.FindStringSubmatch(episodeURL); len(matches) > 1 {
			baseDramaURL = strings.TrimSuffix(matches[1], "/")
		}

		num, _ := strconv.Atoi(currentEpNum)
		if baseDramaURL != "" {
			episodeResponse.Navigation.AllEpisodesURL = baseDramaURL + "/"
			if num > 1 {
				if num == 2 {
					episodeResponse.Navigation.PreviousEpisodeURL = baseDramaURL + "/"
				} else {
					episodeResponse.Navigation.PreviousEpisodeURL = fmt.Sprintf("%s/%d/", baseDramaURL, num-1)
				}
			}
			if doc.Find(fmt.Sprintf("a.post-page-numbers[href*='/%d/']", num+1)).Length() > 0 {
				episodeResponse.Navigation.NextEpisodeURL = fmt.Sprintf("%s/%d/", baseDramaURL, num+1)
			}
		}

		doc.Find("div#action-parts a.post-page-numbers").Each(func(i int, s *goquery.Selection) {
			epNum := s.Find("span").Text()
			epURL := s.AttrOr("href", "")
			episodeResponse.OtherEpisodes = append(episodeResponse.OtherEpisodes, models.OtherEpisode{
				Title: "Episode " + epNum, URL: epURL, ThumbnailURL: thumbnail, ReleaseDate: "Unknown",
			})
		})

		reTitle := regexp.MustCompile(`(.*)\s+\(Episode\s+\d+\)`)
		matches := reTitle.FindStringSubmatch(episodeResponse.Title)
		if len(matches) > 1 {
			episodeResponse.AnimeInfo.Title = strings.TrimSpace(matches[1])
		} else {
			reClean := regexp.MustCompile(`\s+Episode\s+\d+.*`)
			episodeResponse.AnimeInfo.Title = reClean.ReplaceAllString(s.cleanTitle(episodeResponse.Title), "")
		}

		// Mendapatkan parameter AJAX
		log.Println("Mencari parameter AJAX...")
		playerID, exists := doc.Find("div.apicodes-container").Attr("id")
		if !exists {
			log.Println("PERINGATAN: Tidak dapat menemukan Player ID.")
			return
		}

		var nonce string
		dataUri, exists := doc.Find("script#dramagu-player-js-extra").Attr("src")
		if !exists {
			log.Println("PERINGATAN: Tidak dapat menemukan script tag '#dramagu-player-js-extra'.")
			return
		}

		parts := strings.SplitN(dataUri, ",", 2)
		if len(parts) != 2 {
			log.Println("PERINGATAN: Format data URI pada script tidak valid.")
			return
		}

		decodedScript, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			log.Printf("PERINGATAN: Gagal men-decode base64: %v", err)
			return
		}

		reNonce := regexp.MustCompile(`"nonce":"(\w+)"`)
		nonceMatches := reNonce.FindStringSubmatch(string(decodedScript))
		if len(nonceMatches) > 1 {
			nonce = nonceMatches[1]
		} else {
			log.Println("PERINGATAN: Gagal mengekstrak nonce dari script yang di-decode.")
			return
		}

		// Kirim permintaan AJAX
		log.Printf("Parameter ditemukan: PlayerID=%s, Nonce=%s", playerID, nonce)
		ajaxURL := e.Request.AbsoluteURL("/wp-admin/admin-ajax.php")
		formData := map[string]string{
			"action":    "get_player_url",
			"player_id": playerID,
			"nonce":     nonce,
		}
		log.Printf("Mengirim permintaan POST ke: %s", ajaxURL)
		if err := c.Post(ajaxURL, formData); err != nil {
			log.Printf("Gagal mengirimkan permintaan AJAX: %v", err)
		}
	})

	c.OnResponse(func(r *colly.Response) {
		if strings.Contains(r.Request.URL.String(), "admin-ajax.php") {
			log.Println("Menerima respons dari AJAX call.")
			var ajaxResp models.AjaxPlayerResponse
			if err := json.Unmarshal(r.Body, &ajaxResp); err != nil {
				log.Printf("Gagal mem-parsing JSON dari respons AJAX: %v", err)
				return
			}

			if ajaxResp.Success && ajaxResp.Data.IframeURL != "" {
				iframeSrc := ajaxResp.Data.IframeURL
				log.Println("BERHASIL! Link streaming dari AJAX ditemukan:", iframeSrc)

				// Dapatkan nama host dari URL untuk dijadikan nama server
				parsedURL, err := url.Parse(iframeSrc)
				serverName := "Default Server" // Fallback
				if err == nil {
					serverName = strings.ReplaceAll(parsedURL.Hostname(), "www.", "")
				}

				// Tambahkan hanya satu server streaming yang valid
				episodeResponse.StreamingServers = append(episodeResponse.StreamingServers, models.StreamingServer{
					ServerName:   serverName,
					StreamingURL: iframeSrc,
				})

				// Tambahkan juga satu link download sebagai contoh (opsional)
				provider := models.DownloadProvider{Provider: serverName, URL: iframeSrc}
				episodeResponse.DownloadLinks.MKV["720p"] = []models.DownloadProvider{provider}

			} else {
				log.Println("Respons AJAX tidak berhasil atau URL iframe kosong.")
			}
		}
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Mengunjungi:", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Gagal saat request ke %s: %v", r.Request.URL, err)
	})

	err := c.Visit(episodeURL)
	if err != nil {
		return nil, fmt.Errorf("gagal mengunjungi URL: %v", err)
	}

	c.Wait()

	// Calculate confidence score based on data completeness
	episodeResponse.ConfidenceScore = s.calculateConfidenceScore(episodeResponse)

	// Update message based on confidence score
	if episodeResponse.ConfidenceScore == 0.0 {
		episodeResponse.Message = "Data tidak lengkap - field wajib tidak ada"
	} else if episodeResponse.ConfidenceScore < 0.5 {
		episodeResponse.Message = "Data berhasil diambil dengan kelengkapan rendah"
	} else if episodeResponse.ConfidenceScore < 1.0 {
		episodeResponse.Message = "Data berhasil diambil dengan kelengkapan sedang"
	} else {
		episodeResponse.Message = "Data berhasil diambil dengan kelengkapan sempurna"
	}

	return episodeResponse, nil
}

// calculateConfidenceScore calculates confidence score for episode detail response
func (s *EpisodeDetailService) calculateConfidenceScore(response *models.EpisodeDetailResponse) float64 {
	// Required fields: Title, ThumbnailURL, StreamingServers (at least 1), Navigation.AllEpisodesURL
	if !s.hasRequiredFields(response.Title, response.ThumbnailURL, response.Navigation.AllEpisodesURL) ||
		len(response.StreamingServers) == 0 {
		return 0.0
	}

	// Count optional fields
	optionalFieldsCount := 0
	totalOptionalFields := 10 // Total optional fields to check

	// Check optional fields
	if strings.TrimSpace(response.ReleaseInfo) != "" {
		optionalFieldsCount++
	}
	if len(response.DownloadLinks.MKV) > 0 || len(response.DownloadLinks.MP4) > 0 || len(response.DownloadLinks.X265) > 0 {
		optionalFieldsCount++
	}
	if strings.TrimSpace(response.Navigation.PreviousEpisodeURL) != "" {
		optionalFieldsCount++
	}
	if strings.TrimSpace(response.Navigation.NextEpisodeURL) != "" {
		optionalFieldsCount++
	}
	if strings.TrimSpace(response.AnimeInfo.Title) != "" {
		optionalFieldsCount++
	}
	if strings.TrimSpace(response.AnimeInfo.Synopsis) != "" {
		optionalFieldsCount++
	}
	if len(response.AnimeInfo.Genres) > 0 {
		optionalFieldsCount++
	}
	if len(response.OtherEpisodes) > 0 {
		optionalFieldsCount++
	}
	// Source is always present, so count it
	optionalFieldsCount++
	// StreamingServers is required but count as optional for scoring
	optionalFieldsCount++

	// Calculate score based on optional fields completeness
	score := float64(optionalFieldsCount) / float64(totalOptionalFields)

	// Round to 2 decimal places
	return float64(int(score*100)) / 100
}

// hasRequiredFields checks if all required fields are present and not empty
func (s *EpisodeDetailService) hasRequiredFields(fields ...string) bool {
	for _, field := range fields {
		if strings.TrimSpace(field) == "" {
			return false
		}
	}
	return true
}

// cleanTitle cleans and formats the title
func (s *EpisodeDetailService) cleanTitle(title string) string {
	return strings.TrimSpace(title)
}

// generateSlug generates slug from URL
func (s *EpisodeDetailService) generateSlug(url string) string {
	if url == "" {
		return ""
	}
	return path.Base(strings.TrimSuffix(url, "/"))
}
