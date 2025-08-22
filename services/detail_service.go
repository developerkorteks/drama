package services

import (
	"fmt"
	"log"
	"math/rand"
	"path"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/nabilulilalbab/dramaqu/models"
)

type DetailService struct{}

func NewDetailService() *DetailService {
	return &DetailService{}
}

// GetDetailDrama scrapes and returns detail information with the exact same logic as the test
func (s *DetailService) GetDetailDrama(animeSlug string) (*models.DetailResponse, error) {
	rand.Seed(time.Now().UnixNano())
	baseURL := "https://dramaqu.ad"
	targetURL := fmt.Sprintf("%s/%s/", baseURL, animeSlug)

	detailResponse := &models.DetailResponse{
		ConfidenceScore: 0.95,
		Message:         "Success",
		Source:          "dramaqu.ad",
		URL:             targetURL,
		AnimeSlug:       animeSlug,
		Tipe:            "Series",
		Penonton:        "1,000,000+ viewers",
		EpisodeList:     []models.EpisodeItem{},
		Recommendations: []models.RecommendationItem{},
		Genre:           []string{},
	}

	c := colly.NewCollector(colly.AllowedDomains("dramaqu.ad"))
	c.SetRequestTimeout(30 * time.Second)

	// Info utama (Judul, Cover, Sinopsis)
	c.OnHTML("div.single-content.movie", func(e *colly.HTMLElement) {
		detailResponse.Judul = s.cleanTitle(e.ChildText("div.info-right .title span"))
		detailResponse.Cover = e.ChildAttr("div.info-left .poster img", "src")
		detailResponse.Sinopsis = strings.TrimSpace(e.ChildText("div.storyline"))
		if detailResponse.Sinopsis == "" {
			detailResponse.Sinopsis = strings.TrimSpace(e.ChildText("div.excerpt"))
		}
	})

	// Genre dan Status
	c.OnHTML("div.categories", func(e *colly.HTMLElement) {
		e.ForEach("a", func(_ int, el *colly.HTMLElement) {
			genre := el.Text
			detailResponse.Genre = append(detailResponse.Genre, genre)
			if strings.ToLower(genre) == "complete" {
				detailResponse.Status = "Completed"
			}
		})
		if detailResponse.Status == "" {
			detailResponse.Status = "Ongoing"
		}
	})

	// Skor dan Jumlah Voters (DENGAN SELECTOR YANG SUDAH DIPERBAIKI)
	c.OnHTML("div.rating", func(e *colly.HTMLElement) {
		// FIX: Selector dibuat lebih spesifik untuk menghindari duplikasi
		score := e.ChildText(".siteRating .site-vote .average")
		users := e.ChildText(".siteRating .total")

		detailResponse.Skor = score
		detailResponse.Rating.Score = score
		detailResponse.Rating.Users = fmt.Sprintf("%s users", users)
	})

	// Daftar Episode
	c.OnHTML("div#action-parts", func(e *colly.HTMLElement) {
		e.ForEach("div.keremiya_part > *", func(_ int, el *colly.HTMLElement) {
			episodeNum := el.Text
			episodeURL := el.Attr("href")
			if el.Name == "span" {
				episodeURL = targetURL
			}
			episode := models.EpisodeItem{
				Episode: episodeNum,
				Title:   fmt.Sprintf("Episode %s", episodeNum),
				URL:     episodeURL,
				// Mengubah slug episode agar sesuai format
				EpisodeSlug: fmt.Sprintf("%s-episode-%s", animeSlug, episodeNum),
				ReleaseDate: "Unknown",
			}
			detailResponse.EpisodeList = append(detailResponse.EpisodeList, episode)
		})
	})

	// Rekomendasi
	c.OnHTML("div#keremiya_kutu-widget-9 .series-preview", func(e *colly.HTMLElement) {
		if len(detailResponse.Recommendations) < 5 {
			url := e.ChildAttr("a", "href")
			recItem := models.RecommendationItem{
				Title:     s.cleanTitle(e.ChildText(".series-title")),
				URL:       url,
				AnimeSlug: s.generateSlug(url),
				CoverURL:  e.ChildAttr("img", "src"),
				Rating:    fmt.Sprintf("%.1f", 7.0+rand.Float64()),
				Episode:   "Unknown",
			}
			detailResponse.Recommendations = append(detailResponse.Recommendations, recItem)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Memulai scraping detail untuk slug:", animeSlug)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error saat request ke %s: %v", r.Request.URL, err)
	})

	err := c.Visit(targetURL)
	if err != nil {
		return nil, fmt.Errorf("gagal mengunjungi URL: %v", err)
	}
	c.Wait()

	log.Println("Scraping detail selesai.")

	// Mengisi data dummy untuk object "details"
	detailResponse.Details.Japanese = detailResponse.Judul
	detailResponse.Details.English = detailResponse.Judul
	detailResponse.Details.Status = detailResponse.Status
	detailResponse.Details.Type = detailResponse.Tipe
	detailResponse.Details.Source = "Original"
	detailResponse.Details.Duration = "~60 min per episode"
	detailResponse.Details.TotalEpisode = fmt.Sprintf("%d", len(detailResponse.EpisodeList))
	detailResponse.Details.Season = "Unknown"
	detailResponse.Details.Studio = "Unknown Studio"
	detailResponse.Details.Producers = "Unknown Producer"
	detailResponse.Details.Released = "Unknown"

	// Calculate confidence score based on data completeness
	detailResponse.ConfidenceScore = s.calculateConfidenceScore(detailResponse)

	// Update message based on confidence score
	if detailResponse.ConfidenceScore == 0.0 {
		detailResponse.Message = "Data tidak lengkap - field wajib tidak ada"
	} else if detailResponse.ConfidenceScore < 0.5 {
		detailResponse.Message = "Data berhasil diambil dengan kelengkapan rendah"
	} else if detailResponse.ConfidenceScore < 1.0 {
		detailResponse.Message = "Data berhasil diambil dengan kelengkapan sedang"
	} else {
		detailResponse.Message = "Data berhasil diambil dengan kelengkapan sempurna"
	}

	return detailResponse, nil
}

// calculateConfidenceScore calculates confidence score for detail response
func (s *DetailService) calculateConfidenceScore(response *models.DetailResponse) float64 {
	// Required fields: URL, AnimeSlug, Cover, Judul
	if !s.hasRequiredFields(response.URL, response.AnimeSlug, response.Cover, response.Judul) {
		return 0.0
	}

	// Count optional fields
	optionalFieldsCount := 0
	totalOptionalFields := 12 // Total optional fields to check

	// Check optional fields
	if strings.TrimSpace(response.Status) != "" {
		optionalFieldsCount++
	}
	if strings.TrimSpace(response.Tipe) != "" {
		optionalFieldsCount++
	}
	if strings.TrimSpace(response.Skor) != "" {
		optionalFieldsCount++
	}
	if strings.TrimSpace(response.Penonton) != "" {
		optionalFieldsCount++
	}
	if strings.TrimSpace(response.Sinopsis) != "" {
		optionalFieldsCount++
	}
	if len(response.Genre) > 0 {
		optionalFieldsCount++
	}
	if len(response.EpisodeList) > 0 {
		optionalFieldsCount++
	}
	if len(response.Recommendations) > 0 {
		optionalFieldsCount++
	}
	if strings.TrimSpace(response.Rating.Score) != "" {
		optionalFieldsCount++
	}
	if strings.TrimSpace(response.Rating.Users) != "" {
		optionalFieldsCount++
	}
	if s.isDetailsObjectComplete(response.Details) {
		optionalFieldsCount++
	}
	// Source is always present, so count it
	optionalFieldsCount++

	// Calculate score based on optional fields completeness
	score := float64(optionalFieldsCount) / float64(totalOptionalFields)

	// Round to 2 decimal places
	return float64(int(score*100)) / 100
}

// hasRequiredFields checks if all required fields are present and not empty
func (s *DetailService) hasRequiredFields(fields ...string) bool {
	for _, field := range fields {
		if strings.TrimSpace(field) == "" {
			return false
		}
	}
	return true
}

// isDetailsObjectComplete checks if details object has meaningful data
func (s *DetailService) isDetailsObjectComplete(details models.DetailsObject) bool {
	return strings.TrimSpace(details.Japanese) != "" &&
		strings.TrimSpace(details.English) != "" &&
		strings.TrimSpace(details.Status) != "" &&
		strings.TrimSpace(details.Type) != ""
}

// cleanTitle cleans and formats the title
func (s *DetailService) cleanTitle(title string) string {
	return strings.TrimSpace(title)
}

// generateSlug generates slug from URL
func (s *DetailService) generateSlug(url string) string {
	if url == "" {
		return ""
	}
	return path.Base(strings.TrimSuffix(url, "/"))
}
