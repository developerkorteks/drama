// File: ongoing_drama_test.go
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

// ========= STRUCT UNTUK ONGOING DRAMA (Sesuai Permintaan) =========

// Stuktur utama untuk respons API
type OngoingDramaResponse struct {
	ConfidenceScore float64      `json:"confidence_score"`
	Message         string       `json:"message"`
	Source          string       `json:"source"`
	Data            []DramaEntry `json:"data"`
}

// Struktur untuk setiap item drama dalam daftar
type DramaEntry struct {
	Judul    string `json:"judul"`
	URL      string `json:"url"`
	Slug     string `json:"anime_slug"`
	Episode  string `json:"episode"`
	Uploader string `json:"uploader"` // Tidak tersedia, akan diisi nilai default
	Rilis    string `json:"rilis"`    // Tidak tersedia, akan diisi nilai default
	Cover    string `json:"cover"`
}

// ========= FUNGSI TES UNTUK MENGAMBIL ONGOING DRAMA =========

func TestGetOngoingDrama(t *testing.T) {
	// --- INPUT: Ganti nomor halaman di sini ---
	page := 2
	// ----------------------------------------

	// Buat URL target berdasarkan nomor halaman
	baseURL := "https://dramaqu.ad/category/ongoing-drama/"
	targetURL := baseURL
	if page > 1 {
		targetURL = fmt.Sprintf("%spage/%d/", baseURL, page)
	}

	// Inisialisasi struct respons utama
	response := OngoingDramaResponse{
		ConfidenceScore: 1.0,
		Message:         "Data berhasil diambil",
		Source:          "dramaqu.ad",
		Data:            []DramaEntry{}, // Inisialisasi slice kosong
	}

	c := colly.NewCollector(
		colly.AllowedDomains("dramaqu.ad"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"),
	)
	c.SetRequestTimeout(30 * time.Second)

	// Callback untuk setiap elemen <article> yang berisi detail drama
	c.OnHTML("article.movie-preview", func(e *colly.HTMLElement) {
		entry := DramaEntry{}

		// Ekstrak Judul dan URL
		titleElement := e.DOM.Find("span.movie-title a")
		entry.Judul = titleElement.Text()
		entry.URL = titleElement.AttrOr("href", "")

		// Buat Slug dari URL
		if parsedURL, err := url.Parse(entry.URL); err == nil {
			// path.Base akan mengambil bagian terakhir dari path, misal: "my-girlfriend-is-the-man"
			entry.Slug = path.Base(strings.TrimSuffix(parsedURL.Path, "/"))
		}

		// Ekstrak Episode
		entry.Episode = e.DOM.Find("span.icon-hd").Text()

		// Ekstrak Gambar Cover
		entry.Cover = e.DOM.Find("img.keremiya-image").AttrOr("src", "")

		// Isi data yang tidak tersedia dengan nilai default agar struktur konsisten
		entry.Uploader = "DramaQu Admin" // Nilai default
		entry.Rilis = "Unknown"          // Nilai default

		response.Data = append(response.Data, entry)
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Mengunjungi:", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		t.Fatalf("Gagal saat request ke %s: %v", r.Request.URL, err)
	})

	// Kunjungi halaman target
	err := c.Visit(targetURL)
	if err != nil {
		t.Fatalf("Gagal mengunjungi URL: %v", err)
	}
	c.Wait()

	if len(response.Data) == 0 {
		t.Fatal("Tidak ada data drama yang berhasil di-scrape.")
	}

	// Simpan hasil ke file JSON
	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		t.Fatalf("Gagal marshal JSON: %v", err)
	}

	fileName := fmt.Sprintf("ongoing_drama_page_%d_%s.json", page, time.Now().Format("20060102150405"))
	err = ioutil.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		t.Fatalf("Gagal menyimpan file: %v", err)
	}

	log.Printf("SUCCESS: File '%s' berhasil dibuat dengan %d item.", fileName, len(response.Data))
}
