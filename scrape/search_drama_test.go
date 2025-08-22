// File: search_drama_test.go
package scrape

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/gocolly/colly/v2"
)

// ========= STRUCT UNTUK HASIL PENCARIAN (Sesuai Permintaan) =========

type SearchResponse struct {
	ConfidenceScore float64        `json:"confidence_score"`
	Message         string         `json:"message"`
	Source          string         `json:"source"`
	Data            []SearchDetail `json:"data"`
}

type SearchDetail struct {
	Judul    string   `json:"judul"`
	URL      string   `json:"url"`
	Slug     string   `json:"anime_slug"`
	Status   string   `json:"status"`
	Tipe     string   `json:"tipe"`
	Skor     string   `json:"skor"`
	Penonton string   `json:"penonton"`
	Sinopsis string   `json:"sinopsis"`
	Genre    []string `json:"genre"`
	Cover    string   `json:"cover"`
}

// ========= FUNGSI TES UNTUK PENCARIAN DRAMA =========

func TestSearchDrama(t *testing.T) {
	// --- INPUT: Ganti query pencarian dan nomor halaman di sini ---
	query := "a"
	page := 2
	// -----------------------------------------------------------

	// Buat URL pencarian yang benar
	baseURL := "https://dramaqu.ad/"
	targetURL := fmt.Sprintf("%s?s=%s", baseURL, url.QueryEscape(query))
	if page > 1 {
		targetURL = fmt.Sprintf("%spage/%d/?s=%s", baseURL, page, url.QueryEscape(query))
	}

	response := SearchResponse{
		ConfidenceScore: 1.0,
		Message:         "Data berhasil diambil",
		Source:          "dramaqu.ad",
		Data:            []SearchDetail{},
	}

	c := colly.NewCollector(
		colly.AllowedDomains("dramaqu.ad"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"),
	)
	c.SetRequestTimeout(30 * time.Second)

	c.OnHTML("article.movie-preview", func(e *colly.HTMLElement) {
		entry := SearchDetail{}

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
		t.Fatalf("Gagal saat request ke %s: %v", r.Request.URL, err)
	})

	err := c.Visit(targetURL)
	if err != nil {
		t.Fatalf("Gagal mengunjungi URL: %v", err)
	}
	c.Wait()

	if len(response.Data) == 0 {
		t.Fatalf("Tidak ada hasil pencarian untuk query '%s' di halaman %d.", query, page)
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		t.Fatalf("Gagal marshal JSON: %v", err)
	}

	fileName := fmt.Sprintf("search_results_%s_page_%d_%s.json", query, page, time.Now().Format("20060102150405"))
	err = ioutil.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		t.Fatalf("Gagal menyimpan file: %v", err)
	}

	log.Printf("SUCCESS: File '%s' berhasil dibuat dengan %d item.", fileName, len(response.Data))
}
