// File: drama_list_test.go
package scrape

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"path"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/gocolly/colly/v2"
)

// ========= STRUCT UNTUK DRAMA LIST (Sesuai Permintaan) =========

type DramaListResponse struct {
	ConfidenceScore float64       `json:"confidence_score"`
	Message         string        `json:"message"`
	Source          string        `json:"source"`
	Data            []DramaDetail `json:"data"`
}

type DramaDetail struct {
	Judul    string   `json:"judul"`
	URL      string   `json:"url"`
	Slug     string   `json:"anime_slug"`
	Status   string   `json:"status"`
	Skor     string   `json:"skor"`
	Sinopsis string   `json:"sinopsis"`
	Views    string   `json:"views"`
	Cover    string   `json:"cover"`
	Genres   []string `json:"genres"`
	Tanggal  string   `json:"tanggal"`
}

// ========= FUNGSI TES UNTUK MENGAMBIL DRAMA LIST =========

func TestGetDramaList(t *testing.T) {
	// --- INPUT: Ganti nomor halaman di sini ---
	page := 1
	// ----------------------------------------

	baseURL := "https://dramaqu.ad/drama-list/"
	targetURL := baseURL
	if page > 1 {
		targetURL = fmt.Sprintf("%spage/%d/", baseURL, page)
	}

	response := DramaListResponse{
		ConfidenceScore: 1.0,
		Message:         "Data berhasil diambil",
		Source:          "dramaqu.ad",
		Data:            []DramaDetail{},
	}

	c := colly.NewCollector(
		colly.AllowedDomains("dramaqu.ad"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"),
	)
	c.SetRequestTimeout(30 * time.Second)

	// Regex untuk membersihkan angka dari string views
	reViews := regexp.MustCompile(`[0-9,]+`)

	c.OnHTML("article.movie-preview", func(e *colly.HTMLElement) {
		entry := DramaDetail{}

		titleElement := e.DOM.Find("span.movie-title a")
		entry.Judul = titleElement.Text()
		entry.URL = titleElement.AttrOr("href", "")

		// Buat Slug dari URL (mempertahankan 'nonton-')
		if parsedURL, err := url.Parse(entry.URL); err == nil {
			entry.Slug = path.Base(strings.TrimSuffix(parsedURL.Path, "/"))
		}

		entry.Cover = e.DOM.Find("img.keremiya-image").AttrOr("src", "")
		entry.Sinopsis = e.DOM.Find("p.story").Text()
		entry.Tanggal = e.DOM.Find("span.movie-release").Text()

		viewsText := e.DOM.Find("span.views").Text()
		entry.Views = reViews.FindString(viewsText)

		// Isi data yang tidak tersedia dengan nilai default/tetap
		entry.Status = "Completed"
		entry.Skor = "N/A"
		entry.Genres = []string{"Action", "Drama", "Fantasy"} // Data gimmick sesuai permintaan

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
		t.Fatal("Tidak ada data drama yang berhasil di-scrape.")
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		t.Fatalf("Gagal marshal JSON: %v", err)
	}

	fileName := fmt.Sprintf("drama_list_page_%d_%s.json", page, time.Now().Format("20060102150405"))
	err = ioutil.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		t.Fatalf("Gagal menyimpan file: %v", err)
	}

	log.Printf("SUCCESS: File '%s' berhasil dibuat dengan %d item.", fileName, len(response.Data))
}
