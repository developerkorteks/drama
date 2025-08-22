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

// AnimeTerbaruService handles anime terbaru data scraping
type AnimeTerbaruService struct{}

// NewAnimeTerbaruService creates a new instance of AnimeTerbaruService
func NewAnimeTerbaruService() *AnimeTerbaruService {
	return &AnimeTerbaruService{}
}

// GetAnimeTerbaru scrapes and returns anime terbaru data with the exact same logic as the test
func (s *AnimeTerbaruService) GetAnimeTerbaru(page int) (*models.OngoingDramaResponse, error) {
	// Build target URL based on page number
	baseURL := "https://dramaqu.ad/category/ongoing-drama/"
	targetURL := baseURL
	if page > 1 {
		targetURL = fmt.Sprintf("%spage/%d/", baseURL, page)
	}

	// Initialize main response struct
	response := &models.OngoingDramaResponse{
		ConfidenceScore: 0.0, // Will be calculated later
		Message:         "Data berhasil diambil",
		Source:          "dramaqu.ad",
		Data:            []models.DramaEntry{}, // Initialize empty slice
	}

	c := colly.NewCollector(
		colly.AllowedDomains("dramaqu.ad"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"),
	)
	c.SetRequestTimeout(30 * time.Second)

	// Callback for each <article> element containing drama details
	c.OnHTML("article.movie-preview", func(e *colly.HTMLElement) {
		entry := models.DramaEntry{}

		// Extract Title and URL
		titleElement := e.DOM.Find("span.movie-title a")
		entry.Judul = titleElement.Text()
		entry.URL = titleElement.AttrOr("href", "")

		// Create Slug from URL
		if parsedURL, err := url.Parse(entry.URL); err == nil {
			// path.Base will take the last part of the path, e.g.: "my-girlfriend-is-the-man"
			entry.Slug = path.Base(strings.TrimSuffix(parsedURL.Path, "/"))
		}

		// Extract Episode
		entry.Episode = e.DOM.Find("span.icon-hd").Text()

		// Extract Cover Image
		entry.Cover = e.DOM.Find("img.keremiya-image").AttrOr("src", "")

		// Fill unavailable data with default values for consistent structure
		entry.Uploader = "DramaQu Admin" // Default value
		entry.Rilis = "Unknown"          // Default value

		response.Data = append(response.Data, entry)
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Mengunjungi:", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error saat request ke %s: %v", r.Request.URL, err)
	})

	// Visit target page
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

// calculateConfidenceScore calculates confidence score based on data completeness
func (s *AnimeTerbaruService) calculateConfidenceScore(response *models.OngoingDramaResponse) float64 {
	if len(response.Data) == 0 {
		return 0.0
	}

	totalItems := len(response.Data)
	validItems := 0.0

	for _, item := range response.Data {
		if s.isDramaEntryValid(item) {
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

// hasRequiredFields checks if all required fields are present and not empty
func (s *AnimeTerbaruService) hasRequiredFields(fields ...string) bool {
	for _, field := range fields {
		if strings.TrimSpace(field) == "" {
			return false
		}
	}
	return true
}

// isDramaEntryValid checks if DramaEntry has all required and optional fields
func (s *AnimeTerbaruService) isDramaEntryValid(item models.DramaEntry) bool {
	// Required fields: Judul, URL, Slug, Cover
	if !s.hasRequiredFields(item.Judul, item.URL, item.Slug, item.Cover) {
		return false
	}

	// Optional fields: Episode, Uploader, Rilis
	optionalFieldsCount := 0
	totalOptionalFields := 3

	if strings.TrimSpace(item.Episode) != "" {
		optionalFieldsCount++
	}
	if strings.TrimSpace(item.Uploader) != "" {
		optionalFieldsCount++
	}
	if strings.TrimSpace(item.Rilis) != "" {
		optionalFieldsCount++
	}

	// Consider valid if all optional fields are present
	return optionalFieldsCount == totalOptionalFields
}
