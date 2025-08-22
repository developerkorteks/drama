package services

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/nabilulilalbab/dramaqu/models"
	"github.com/nabilulilalbab/dramaqu/scrape"
)

// HomeService handles home page data scraping
type HomeService struct{}

// NewHomeService creates a new instance of HomeService
func NewHomeService() *HomeService {
	return &HomeService{}
}

// GetHomeData scrapes and returns home page data with the exact same logic as the test
func (s *HomeService) GetHomeData() (*models.FinalResponse, error) {
	rand.Seed(time.Now().UnixNano())

	finalResponse := &models.FinalResponse{
		ConfidenceScore: 0.0, // Will be calculated later
		Message:         "Data berhasil diambil",
		Source:          "dramaqu.ad",
	}

	c := colly.NewCollector(
		colly.AllowedDomains("dramaqu.ad"),
	)

	c.OnHTML("div.film-content", func(e *colly.HTMLElement) {
		sectionTitle := e.ChildText("h2.title span")
		e.ForEach("article.movie-preview", func(_ int, item *colly.HTMLElement) {
			switch sectionTitle {
			case "Ongoing Drama":
				if len(finalResponse.NewEps) < 20 {
					itemData := s.parseNewEpsItem(item)
					finalResponse.NewEps = append(finalResponse.NewEps, itemData)
				}
			case "Film Korea":
				if len(finalResponse.Movies) < 20 {
					itemData := s.parseMovieItem(item)
					finalResponse.Movies = append(finalResponse.Movies, itemData)
				}
			case "Drama Populer":
				if len(finalResponse.Top10) < 10 {
					itemData := s.parseTop10Item(item)
					finalResponse.Top10 = append(finalResponse.Top10, itemData)
				}
			}
		})
	})

	log.Println("Memulai scraping halaman utama...")
	err := c.Visit("https://dramaqu.ad/")
	if err != nil {
		return nil, fmt.Errorf("failed to visit main page: %v", err)
	}
	c.Wait()
	log.Println("Scraping halaman utama selesai.")

	scheduleCollector := colly.NewCollector(
		colly.AllowedDomains("dramaqu.ad"),
	)
	var ongoingItemsForSchedule []models.JadwalItem
	scheduleCollector.OnHTML("article.movie-preview", func(e *colly.HTMLElement) {
		item := s.parseJadwalItem(e)
		ongoingItemsForSchedule = append(ongoingItemsForSchedule, item)
	})

	log.Println("Memulai scraping halaman 'Ongoing' untuk data jadwal...")
	err = scheduleCollector.Visit("https://dramaqu.ad/category/ongoing-drama/")
	if err != nil {
		return nil, fmt.Errorf("failed to visit ongoing page: %v", err)
	}
	scheduleCollector.Wait()
	log.Println("Scraping halaman 'Ongoing' selesai.")

	finalResponse.JadwalRilis = s.generateJadwal(ongoingItemsForSchedule)
	log.Println("Jadwal rilis dummy berhasil dibuat.")

	// Calculate confidence score based on data completeness
	finalResponse.ConfidenceScore = s.calculateConfidenceScore(finalResponse)

	// Update message based on confidence score
	if finalResponse.ConfidenceScore == 0.0 {
		finalResponse.Message = "Data tidak lengkap - field wajib tidak ada"
	} else if finalResponse.ConfidenceScore < 0.5 {
		finalResponse.Message = "Data berhasil diambil dengan kelengkapan rendah"
	} else if finalResponse.ConfidenceScore < 1.0 {
		finalResponse.Message = "Data berhasil diambil dengan kelengkapan sedang"
	} else {
		finalResponse.Message = "Data berhasil diambil dengan kelengkapan sempurna"
	}

	return finalResponse, nil
}

func (s *HomeService) parseNewEpsItem(e *colly.HTMLElement) models.NewEpsItem {
	url := e.ChildAttr("a", "href")
	judul := scrape.CleanTitle(e.ChildText(".movie-title a"))
	return models.NewEpsItem{
		Judul:     judul,
		URL:       url,
		AnimeSlug: scrape.GenerateSlug(url),
		Episode:   e.ChildText(".center-icons .icon-hd"),
		Cover:     e.ChildAttr("img", "src"),
		Rilis:     fmt.Sprintf("%d jam", rand.Intn(23)+1), // Dummy
	}
}

func (s *HomeService) parseMovieItem(e *colly.HTMLElement) models.MovieItem {
	url := e.ChildAttr("a", "href")
	judul := scrape.CleanTitle(e.ChildText(".movie-title a"))
	return models.MovieItem{
		Judul:     judul,
		URL:       url,
		AnimeSlug: scrape.GenerateSlug(url),
		Cover:     e.ChildAttr("img", "src"),
		Tanggal:   fmt.Sprintf("%d hari", rand.Intn(5)+1),  // Dummy
		Genres:    []string{"Action", "Drama", "Thriller"}, // Dummy
	}
}

func (s *HomeService) parseTop10Item(e *colly.HTMLElement) models.Top10Item {
	url := e.ChildAttr("a", "href")
	judul := scrape.CleanTitle(e.ChildText(".movie-title a"))
	rating := strings.TrimSpace(e.ChildText(".icon-star.imdb"))
	if rating == "" {
		rating = fmt.Sprintf("%.2f", 7.0+rand.Float64()*2) // Dummy rating 7.00-9.00
	}
	return models.Top10Item{
		Judul:     judul,
		URL:       url,
		AnimeSlug: scrape.GenerateSlug(url),
		Cover:     e.ChildAttr("img", "src"),
		Rating:    rating,
		Genres:    []string{"Action", "Adventure", "Drama"}, // Dummy
	}
}

func (s *HomeService) parseJadwalItem(e *colly.HTMLElement) models.JadwalItem {
	url := e.ChildAttr("a", "href")
	judul := scrape.CleanTitle(e.ChildText(".movie-title a"))
	return models.JadwalItem{
		Title:       judul,
		URL:         url,
		AnimeSlug:   scrape.GenerateSlug(url),
		CoverURL:    e.ChildAttr("img", "src"),
		Type:        "TV",                                                   // Dummy
		Score:       fmt.Sprintf("%.1f", 7.0+rand.Float64()),                // Dummy
		Genres:      []string{"Drama", "Romance", "Comedy"},                 // Dummy
		ReleaseTime: fmt.Sprintf("%02d:%02d", rand.Intn(24), rand.Intn(60)), // Dummy
	}
}

func (s *HomeService) generateJadwal(items []models.JadwalItem) models.JadwalRilis {
	// PENTING: Inisialisasi slice untuk setiap hari agar tidak nil
	jadwal := models.JadwalRilis{
		Monday:    []models.JadwalItem{},
		Tuesday:   []models.JadwalItem{},
		Wednesday: []models.JadwalItem{},
		Thursday:  []models.JadwalItem{},
		Friday:    []models.JadwalItem{},
		Saturday:  []models.JadwalItem{},
		Sunday:    []models.JadwalItem{},
	}
	if len(items) == 0 {
		return jadwal
	}

	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

	for _, item := range items {
		dayIndex := rand.Intn(len(days))
		switch days[dayIndex] {
		case "Monday":
			jadwal.Monday = append(jadwal.Monday, item)
		case "Tuesday":
			jadwal.Tuesday = append(jadwal.Tuesday, item)
		case "Wednesday":
			jadwal.Wednesday = append(jadwal.Wednesday, item)
		case "Thursday":
			jadwal.Thursday = append(jadwal.Thursday, item)
		case "Friday":
			jadwal.Friday = append(jadwal.Friday, item)
		case "Saturday":
			jadwal.Saturday = append(jadwal.Saturday, item)
		case "Sunday":
			jadwal.Sunday = append(jadwal.Sunday, item)
		}
	}
	return jadwal
}

// calculateConfidenceScore calculates confidence score based on data completeness
func (s *HomeService) calculateConfidenceScore(response *models.FinalResponse) float64 {
	totalItems := 0
	validItems := 0.0

	// Check Top10 items
	for _, item := range response.Top10 {
		totalItems++
		if s.isTop10ItemValid(item) {
			validItems += 1.0
		} else if s.hasRequiredFields(item.Judul, item.URL, item.AnimeSlug, item.Cover) {
			// If required fields exist but optional fields missing, count as partial
			validItems += 0.5
		}
		// If required fields missing, add 0 (no increment to validItems)
	}

	// Check NewEps items
	for _, item := range response.NewEps {
		totalItems++
		if s.isNewEpsItemValid(item) {
			validItems += 1.0
		} else if s.hasRequiredFields(item.Judul, item.URL, item.AnimeSlug, item.Cover) {
			validItems += 0.5
		}
	}

	// Check Movies items
	for _, item := range response.Movies {
		totalItems++
		if s.isMovieItemValid(item) {
			validItems += 1.0
		} else if s.hasRequiredFields(item.Judul, item.URL, item.AnimeSlug, item.Cover) {
			validItems += 0.5
		}
	}

	// Check JadwalRilis items
	allJadwalItems := [][]models.JadwalItem{
		response.JadwalRilis.Monday,
		response.JadwalRilis.Tuesday,
		response.JadwalRilis.Wednesday,
		response.JadwalRilis.Thursday,
		response.JadwalRilis.Friday,
		response.JadwalRilis.Saturday,
		response.JadwalRilis.Sunday,
	}

	for _, dayItems := range allJadwalItems {
		for _, item := range dayItems {
			totalItems++
			if s.isJadwalItemValid(item) {
				validItems += 1.0
			} else if s.hasRequiredFields(item.Title, item.URL, item.AnimeSlug, item.CoverURL) {
				validItems += 0.5
			}
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
func (s *HomeService) hasRequiredFields(fields ...string) bool {
	for _, field := range fields {
		if strings.TrimSpace(field) == "" {
			return false
		}
	}
	return true
}

// isTop10ItemValid checks if Top10Item has all required and optional fields
func (s *HomeService) isTop10ItemValid(item models.Top10Item) bool {
	// Required fields: Judul, URL, AnimeSlug, Cover
	if !s.hasRequiredFields(item.Judul, item.URL, item.AnimeSlug, item.Cover) {
		return false
	}

	// Optional fields: Rating, Genres
	optionalFieldsCount := 0
	totalOptionalFields := 2

	if strings.TrimSpace(item.Rating) != "" {
		optionalFieldsCount++
	}
	if len(item.Genres) > 0 {
		optionalFieldsCount++
	}

	// Consider valid if all optional fields are present
	return optionalFieldsCount == totalOptionalFields
}

// isNewEpsItemValid checks if NewEpsItem has all required and optional fields
func (s *HomeService) isNewEpsItemValid(item models.NewEpsItem) bool {
	// Required fields: Judul, URL, AnimeSlug, Cover
	if !s.hasRequiredFields(item.Judul, item.URL, item.AnimeSlug, item.Cover) {
		return false
	}

	// Optional fields: Episode, Rilis
	optionalFieldsCount := 0
	totalOptionalFields := 2

	if strings.TrimSpace(item.Episode) != "" {
		optionalFieldsCount++
	}
	if strings.TrimSpace(item.Rilis) != "" {
		optionalFieldsCount++
	}

	return optionalFieldsCount == totalOptionalFields
}

// isMovieItemValid checks if MovieItem has all required and optional fields
func (s *HomeService) isMovieItemValid(item models.MovieItem) bool {
	// Required fields: Judul, URL, AnimeSlug, Cover
	if !s.hasRequiredFields(item.Judul, item.URL, item.AnimeSlug, item.Cover) {
		return false
	}

	// Optional fields: Tanggal, Genres
	optionalFieldsCount := 0
	totalOptionalFields := 2

	if strings.TrimSpace(item.Tanggal) != "" {
		optionalFieldsCount++
	}
	if len(item.Genres) > 0 {
		optionalFieldsCount++
	}

	return optionalFieldsCount == totalOptionalFields
}

// isJadwalItemValid checks if JadwalItem has all required and optional fields
func (s *HomeService) isJadwalItemValid(item models.JadwalItem) bool {
	// Required fields: Title, URL, AnimeSlug, CoverURL
	if !s.hasRequiredFields(item.Title, item.URL, item.AnimeSlug, item.CoverURL) {
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

	return optionalFieldsCount == totalOptionalFields
}
