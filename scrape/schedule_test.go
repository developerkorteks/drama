// File: release_schedule_test.go
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

// ========= STRUCT UNTUK JADWAL RILIS (Sesuai Permintaan) =========

type ReleaseScheduleResponse struct {
	ConfidenceScore float64                   `json:"confidence_score"`
	Message         string                    `json:"message"`
	Source          string                    `json:"source"`
	Data            map[string][]ReleaseEntry `json:"data"`
}

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

// ========= FUNGSI TES UNTUK MENGAMBIL JADWAL RILIS =========

func TestGetReleaseSchedule(t *testing.T) {
	targetURL := "https://dramaqu.ad/category/ongoing-drama/"

	// Map untuk menampung data yang dikelompokkan berdasarkan hari
	scheduleData := make(map[string][]ReleaseEntry)
	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	for _, day := range days {
		scheduleData[day] = []ReleaseEntry{} // Inisialisasi setiap hari dengan slice kosong
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

		entry := ReleaseEntry{
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
		t.Fatalf("Gagal saat request ke %s: %v", r.Request.URL, err)
	})

	err := c.Visit(targetURL)
	if err != nil {
		t.Fatalf("Gagal mengunjungi URL: %v", err)
	}
	c.Wait()

	if itemCounter == 0 {
		t.Fatal("Tidak ada data drama yang berhasil di-scrape.")
	}

	// Buat respons akhir
	finalResponse := ReleaseScheduleResponse{
		ConfidenceScore: 1.0,
		Message:         "Data berhasil diambil",
		Source:          "dramaqu.ad",
		Data:            scheduleData,
	}

	jsonData, err := json.MarshalIndent(finalResponse, "", "  ")
	if err != nil {
		t.Fatalf("Gagal marshal JSON: %v", err)
	}

	fileName := fmt.Sprintf("release_schedule_%s.json", time.Now().Format("20060102150405"))
	err = ioutil.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		t.Fatalf("Gagal menyimpan file: %v", err)
	}

	log.Printf("SUCCESS: File '%s' berhasil dibuat dengan %d total item.", fileName, itemCounter)
}
