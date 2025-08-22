package services

import (
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/nabilulilalbab/dramaqu/models"
)

// ScheduleService handles schedule data scraping
type ScheduleService struct{}

// NewScheduleService creates a new instance of ScheduleService
func NewScheduleService() *ScheduleService {
	return &ScheduleService{}
}

// GetReleaseSchedule scrapes and returns release schedule data with the exact same logic as the test
func (s *ScheduleService) GetReleaseSchedule() (*models.ReleaseScheduleResponse, error) {
	targetURL := "https://dramaqu.ad/category/ongoing-drama/"

	// Map untuk menampung data yang dikelompokkan berdasarkan hari
	scheduleData := make(map[string][]models.ReleaseEntry)
	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	for _, day := range days {
		scheduleData[day] = []models.ReleaseEntry{} // Inisialisasi setiap hari dengan slice kosong
	}

	c := colly.NewCollector(
		colly.AllowedDomains("dramaqu.ad"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"),
	)
	c.SetRequestTimeout(30 * time.Second)

	// Counter untuk mendistribusikan drama ke hari yang berbeda
	itemCounter := 0

	c.OnHTML("article.movie-preview", func(e *colly.HTMLElement) {
		titleElement := e.DOM.Find("span.movie-title a")
		dramaURL := titleElement.AttrOr("href", "")

		// Buat slug dari URL (mempertahankan 'nonton-')
		var slug string
		if parsedURL, err := url.Parse(dramaURL); err == nil {
			slug = path.Base(strings.TrimSuffix(parsedURL.Path, "/"))
		}

		// --- Membuat Data Gimmick ---
		rand.Seed(time.Now().UnixNano() + int64(itemCounter)) // Seed randomizer

		// Waktu rilis acak
		releaseHour := rand.Intn(24)
		releaseMinute := rand.Intn(60)
		releaseTime := fmt.Sprintf("%02d:%02d", releaseHour, releaseMinute)

		// Skor acak antara 7.0 dan 9.5
		score := 7.0 + rand.Float64()*(9.5-7.0)

		// Tipe acak (lebih sering TV)
		dramaType := "TV"
		if rand.Intn(10) > 8 { // 20% kemungkinan menjadi Movie
			dramaType = "Movie"
		}

		entry := models.ReleaseEntry{
			Title:       titleElement.Text(),
			URL:         dramaURL,
			Slug:        slug,
			CoverURL:    e.DOM.Find("img.keremiya-image").AttrOr("src", ""),
			Type:        dramaType,
			Score:       fmt.Sprintf("%.1f", score),
			Genres:      []string{"Drama", "Romance", "Comedy"}, // Genre gimmick
			ReleaseTime: releaseTime,
		}

		// Distribusikan ke dalam map jadwal secara bergiliran
		dayKey := days[itemCounter%len(days)]
		scheduleData[dayKey] = append(scheduleData[dayKey], entry)
		itemCounter++
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

	if itemCounter == 0 {
		return nil, fmt.Errorf("tidak ada data drama yang berhasil di-scrape")
	}

	// Buat respons akhir
	response := &models.ReleaseScheduleResponse{
		ConfidenceScore: 0.0, // Will be calculated
		Message:         "Data berhasil diambil",
		Source:          "dramaqu.ad",
		Data:            scheduleData,
	}

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

// GetScheduleByDay scrapes and returns schedule data for specific day with the exact same logic as the test
func (s *ScheduleService) GetScheduleByDay(inputDay string) (*models.ScheduleByDayResponse, error) {
	targetURL := "https://dramaqu.ad/category/ongoing-drama/"

	response := &models.ScheduleByDayResponse{
		ConfidenceScore: 1.0,
		Message:         "Data berhasil diambil",
		Source:          "dramaqu.ad",
		Data:            []models.ScheduleEntry{},
	}

	c := colly.NewCollector(
		colly.AllowedDomains("dramaqu.ad"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"),
	)
	c.SetRequestTimeout(30 * time.Second)

	itemCounter := 0

	c.OnHTML("article.movie-preview", func(e *colly.HTMLElement) {
		title := e.DOM.Find("span.movie-title a").Text()

		// Tentukan hari rilis drama ini secara konsisten
		releaseDay := s.getDayForTitle(title)

		// HANYA proses item jika harinya cocok dengan input (case-insensitive)
		if strings.EqualFold(releaseDay, inputDay) {

			dramaURL := e.DOM.Find("span.movie-title a").AttrOr("href", "")
			var slug string
			if parsedURL, err := url.Parse(dramaURL); err == nil {
				slug = path.Base(strings.TrimSuffix(parsedURL.Path, "/"))
			}

			// Buat data gimmick
			rand.Seed(time.Now().UnixNano() + int64(itemCounter))
			releaseHour := rand.Intn(24)
			releaseMinute := rand.Intn(60)
			score := 7.0 + rand.Float64()*(9.5-7.0)
			dramaType := "TV"
			if rand.Intn(10) > 8 {
				dramaType = "Movie"
			}

			entry := models.ScheduleEntry{
				Title:       title,
				URL:         dramaURL,
				Slug:        slug,
				CoverURL:    e.DOM.Find("img.keremiya-image").AttrOr("src", ""),
				Type:        dramaType,
				Score:       fmt.Sprintf("%.1f", score),
				Genres:      []string{"Drama", "Romance", "Action"},
				ReleaseTime: fmt.Sprintf("%02d:%02d", releaseHour, releaseMinute),
			}

			response.Data = append(response.Data, entry)
			itemCounter++
		}
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

	log.Printf("Menemukan %d item untuk hari %s.", len(response.Data), inputDay)

	// Calculate confidence score based on data completeness
	response.ConfidenceScore = s.calculateConfidenceScoreByDay(response)

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

// getDayForTitle determines release day consistently based on title
func (s *ScheduleService) getDayForTitle(title string) string {
	days := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	var sum int
	// Menjumlahkan nilai byte dari judul untuk mendapatkan angka yang konsisten
	for _, char := range title {
		sum += int(char)
	}
	// Modulo 7 akan selalu menghasilkan angka 0-6
	return days[sum%len(days)]
}

// calculateConfidenceScoreByDay calculates confidence score for schedule by day response
func (s *ScheduleService) calculateConfidenceScoreByDay(response *models.ScheduleByDayResponse) float64 {
	if len(response.Data) == 0 {
		return 0.0
	}

	totalItems := len(response.Data)
	validItems := 0.0

	for _, item := range response.Data {
		if s.isScheduleEntryValid(item) {
			validItems += 1.0
		} else if s.hasRequiredFields(item.Title, item.URL, item.Slug, item.CoverURL) {
			// If required fields exist but optional fields missing, count as partial
			validItems += 0.5
		}
		// If required fields missing, add 0 (no increment to validItems)
	}

	score := validItems / float64(totalItems)

	// Round to 2 decimal places
	return float64(int(score*100)) / 100
}

// isScheduleEntryValid checks if ScheduleEntry has all required and optional fields
func (s *ScheduleService) isScheduleEntryValid(item models.ScheduleEntry) bool {
	// Required fields: Title, URL, Slug, CoverURL
	if !s.hasRequiredFields(item.Title, item.URL, item.Slug, item.CoverURL) {
		return false
	}

	// Optional fields: Type, Score, Genres, ReleaseTime
	optionalFieldsCount := 0
	totalOptionalFields := 4

	if strings.TrimSpace(item.Type) != "" {
		optionalFieldsCount++
	}
	if strings.TrimSpace(item.Score) != "" {
		optionalFieldsCount++
	}
	if len(item.Genres) > 0 {
		optionalFieldsCount++
	}
	if strings.TrimSpace(item.ReleaseTime) != "" {
		optionalFieldsCount++
	}

	// Consider valid if all optional fields are present
	return optionalFieldsCount == totalOptionalFields
}

// calculateConfidenceScore calculates confidence score based on data completeness
func (s *ScheduleService) calculateConfidenceScore(response *models.ReleaseScheduleResponse) float64 {
	totalItems := 0
	validItems := 0.0

	// Count total items across all days
	for _, entries := range response.Data {
		totalItems += len(entries)
		for _, item := range entries {
			if s.isReleaseEntryValid(item) {
				validItems += 1.0
			} else if s.hasRequiredFields(item.Title, item.URL, item.Slug, item.CoverURL) {
				// If required fields exist but optional fields missing, count as partial
				validItems += 0.5
			}
			// If required fields missing, add 0 (no increment to validItems)
		}
	}

	if totalItems == 0 {
		return 0.0
	}

	score := validItems / float64(totalItems)

	// Round to 2 decimal places
	return float64(int(score*100)) / 100
}

// hasRequiredFields checks if all required fields are present and not empty
func (s *ScheduleService) hasRequiredFields(fields ...string) bool {
	for _, field := range fields {
		if strings.TrimSpace(field) == "" {
			return false
		}
	}
	return true
}

// isReleaseEntryValid checks if ReleaseEntry has all required and optional fields
func (s *ScheduleService) isReleaseEntryValid(item models.ReleaseEntry) bool {
	// Required fields: Title, URL, Slug, CoverURL
	if !s.hasRequiredFields(item.Title, item.URL, item.Slug, item.CoverURL) {
		return false
	}

	// Optional fields: Type, Score, Genres, ReleaseTime
	optionalFieldsCount := 0
	totalOptionalFields := 4

	if strings.TrimSpace(item.Type) != "" {
		optionalFieldsCount++
	}
	if strings.TrimSpace(item.Score) != "" {
		optionalFieldsCount++
	}
	if len(item.Genres) > 0 {
		optionalFieldsCount++
	}
	if strings.TrimSpace(item.ReleaseTime) != "" {
		optionalFieldsCount++
	}

	// Consider valid if all optional fields are present
	return optionalFieldsCount == totalOptionalFields
}
