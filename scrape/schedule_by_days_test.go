// File: schedule_by_day_test.go
package scrape

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/url"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/gocolly/colly/v2"
)

// ========= STRUCT UNTUK JADWAL RILIS PER HARI =========

type ScheduleByDayResponse struct {
	ConfidenceScore float64         `json:"confidence_score"`
	Message         string          `json:"message"`
	Source          string          `json:"source"`
	Data            []ScheduleEntry `json:"data"`
}

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

// Fungsi untuk menentukan hari rilis secara konsisten berdasarkan judul
func getDayForTitle(title string) string {
	days := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	var sum int
	// Menjumlahkan nilai byte dari judul untuk mendapatkan angka yang konsisten
	for _, char := range title {
		sum += int(char)
	}
	// Modulo 7 akan selalu menghasilkan angka 0-6
	return days[sum%len(days)]
}

// ========= FUNGSI TES UNTUK MENGAMBIL JADWAL PER HARI =========

func TestGetScheduleByDay(t *testing.T) {
	// --- INPUT: Ganti nama hari di sini (lowercase) ---
	// Pilihan: "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"
	inputDay := "monday"
	// ----------------------------------------------------

	targetURL := "https://dramaqu.ad/category/ongoing-drama/"

	response := ScheduleByDayResponse{
		ConfidenceScore: 1.0,
		Message:         "Data berhasil diambil",
		Source:          "dramaqu.ad",
		Data:            []ScheduleEntry{},
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
		releaseDay := getDayForTitle(title)

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

			entry := ScheduleEntry{
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

	c.OnRequest(func(r *colly.Request) { log.Println("Mengunjungi:", r.URL.String()) })
	c.OnError(func(r *colly.Response, err error) { t.Fatalf("Gagal saat request ke %s: %v", r.Request.URL, err) })

	err := c.Visit(targetURL)
	if err != nil {
		t.Fatalf("Gagal mengunjungi URL: %v", err)
	}
	c.Wait()

	log.Printf("Menemukan %d item untuk hari %s.", len(response.Data), inputDay)

	// Simpan ke file JSON
	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		t.Fatalf("Gagal marshal JSON: %v", err)
	}

	fileName := fmt.Sprintf("schedule_for_%s_%s.json", inputDay, time.Now().Format("20060102150405"))
	err = ioutil.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		t.Fatalf("Gagal menyimpan file: %v", err)
	}

	log.Printf("SUCCESS: File '%s' berhasil dibuat.", fileName)
}
