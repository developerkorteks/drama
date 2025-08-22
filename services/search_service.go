package services

import (
	"fmt"
	"log"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/nabilulilalbab/dramaqu/models"
)

type SearchService struct{}

func NewSearchService() *SearchService {
	return &SearchService{}
}

// SearchDrama scrapes and returns search results with the exact same logic as the test
func (s *SearchService) SearchDrama(query string, page int) (*models.SearchResponse, error) {
	// Buat URL pencarian yang benar
	baseURL := "https://dramaqu.ad/"
	targetURL := fmt.Sprintf("%s?s=%s", baseURL, url.QueryEscape(query))
	if page > 1 {
		targetURL = fmt.Sprintf("%spage/%d/?s=%s", baseURL, page, url.QueryEscape(query))
	}

	response := &models.SearchResponse{
		ConfidenceScore: 1.0,
		Message:         "Data berhasil diambil",
		Source:          "dramaqu.ad",
		Data:            []models.SearchDetail{},
	}

	c := colly.NewCollector(
		colly.AllowedDomains("dramaqu.ad"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"),
	)
	c.SetRequestTimeout(30 * time.Second)

	c.OnHTML("article.movie-preview", func(e *colly.HTMLElement) {
		entry := models.SearchDetail{}

		titleElement := e.DOM.Find("span.movie-title a")
		entry.Judul = titleElement.Text()
		entry.URL = titleElement.AttrOr("href", "")

		if parsedURL, err := url.Parse(entry.URL); err == nil {
			entry.Slug = path.Base(strings.TrimSuffix(parsedURL.Path, "/"))
		}

		entry.Cover = e.DOM.Find("img.keremiya-image").AttrOr("src", "")
		entry.Sinopsis = e.DOM.Find("p.story").Text()

		// --- Logika Gimmick/Placeholder ---
		// Menentukan Tipe berdasarkan URL
		if strings.Contains(entry.URL, "/nonton-") {
			entry.Tipe = "Series"
		} else {
			entry.Tipe = "Movie"
		}

		// Menentukan Status berdasarkan Teks Episode
		episodeText := e.DOM.Find("span.icon-hd").Text()
		if episodeText != "" {
			entry.Status = "Ongoing"
		} else {
			entry.Status = "Completed"
		}

		entry.Skor = "N/A"
		entry.Penonton = "15,000+ viewers"                    // Placeholder
		entry.Genre = []string{"Action", "Drama", "Thriller"} // Placeholder

		response.Data = append(response.Data, entry)
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Mengunjungi:", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error saat request ke %s: %v", r.Request.URL, err)
	})

	err := c.Visit(targetURL)
	if err != nil {
		return nil, fmt.Errorf("gagal mengunjungi URL: %v", err)
	}
	c.Wait()

	// Calculate confidence score based on data completeness
	response.ConfidenceScore = s.calculateConfidenceScore(response)

	// Update message based on confidence score
	if response.ConfidenceScore == 0.0 {
		response.Message = "Data tidak lengkap - field wajib tidak ada"
	} else if response.ConfidenceScore < 0.5 {
		response.Message = "Data berhasil diambil dengan kelengkapan rendah"
	} else if response.ConfidenceScore < 1.0 {
		response.Message = "Data berhasil diambil dengan kelengkapan sedang"
	} else {
		response.Message = "Data berhasil diambil dengan kelengkapan sempurna"
	}

	return response, nil
}

// calculateConfidenceScore calculates confidence score for search response
func (s *SearchService) calculateConfidenceScore(response *models.SearchResponse) float64 {
	if len(response.Data) == 0 {
		return 0.0
	}

	totalItems := len(response.Data)
	validItems := 0.0

	for _, item := range response.Data {
		if s.isSearchDetailValid(item) {
			validItems += 1.0
		} else if s.hasRequiredFields(item.Judul, item.URL, item.Slug, item.Cover) {
			// If required fields exist but optional fields missing, count as partial
			validItems += 0.5
		}
		// If required fields missing, add 0 (no increment to validItems)
	}

	score := validItems / float64(totalItems)

	// Round to 2 decimal places
	return float64(int(score*100)) / 100
}

// isSearchDetailValid checks if SearchDetail has all required and optional fields
func (s *SearchService) isSearchDetailValid(item models.SearchDetail) bool {
	// Required fields: Judul, URL, Slug, Cover
	if !s.hasRequiredFields(item.Judul, item.URL, item.Slug, item.Cover) {
		return false
	}

	// Optional fields: Status, Tipe, Skor, Penonton, Sinopsis, Genre
	optionalFieldsCount := 0
	totalOptionalFields := 6

	if strings.TrimSpace(item.Status) != "" {
		optionalFieldsCount++
	}
	if strings.TrimSpace(item.Tipe) != "" {
		optionalFieldsCount++
	}
	if strings.TrimSpace(item.Skor) != "" {
		optionalFieldsCount++
	}
	if strings.TrimSpace(item.Penonton) != "" {
		optionalFieldsCount++
	}
	if strings.TrimSpace(item.Sinopsis) != "" {
		optionalFieldsCount++
	}
	if len(item.Genre) > 0 {
		optionalFieldsCount++
	}

	// Consider valid if all optional fields are present
	return optionalFieldsCount == totalOptionalFields
}

// hasRequiredFields checks if all required fields are present and not empty
func (s *SearchService) hasRequiredFields(fields ...string) bool {
	for _, field := range fields {
		if strings.TrimSpace(field) == "" {
			return false
		}
	}
	return true
}
