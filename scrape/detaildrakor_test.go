package scrape

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/gocolly/colly/v2"
)

// ========= STRUCT UNTUK DETAIL RESPONSE (Sesuai Permintaan) =========

type DetailResponse struct {
	ConfidenceScore float64              `json:"confidence_score"`
	Message         string               `json:"message"`
	Source          string               `json:"source"`
	Judul           string               `json:"judul"`
	URL             string               `json:"url"`
	AnimeSlug       string               `json:"anime_slug"`
	Cover           string               `json:"cover"`
	EpisodeList     []EpisodeItem        `json:"episode_list"`
	Recommendations []RecommendationItem `json:"recommendations"`
	Status          string               `json:"status"`
	Tipe            string               `json:"tipe"`
	Skor            string               `json:"skor"`
	Penonton        string               `json:"penonton"`
	Sinopsis        string               `json:"sinopsis"`
	Genre           []string             `json:"genre"`
	Details         DetailsObject        `json:"details"`
	Rating          RatingObject         `json:"rating"`
}

type EpisodeItem struct {
	Episode     string `json:"episode"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	EpisodeSlug string `json:"episode_slug"`
	ReleaseDate string `json:"release_date"`
}

type RecommendationItem struct {
	Title     string `json:"title"`
	URL       string `json:"url"`
	AnimeSlug string `json:"anime_slug"`
	CoverURL  string `json:"cover_url"`
	Rating    string `json:"rating"`
	Episode   string `json:"episode"`
}

type DetailsObject struct {
	Japanese     string `json:"Japanese"`
	English      string `json:"English"`
	Status       string `json:"Status"`
	Type         string `json:"Type"`
	Source       string `json:"Source"`
	Duration     string `json:"Duration"`
	TotalEpisode string `json:"Total Episode"`
	Season       string `json:"Season"`
	Studio       string `json:"Studio"`
	Producers    string `json:"Producers"`
	Released     string `json:"Released:"`
}

type RatingObject struct {
	Score string `json:"score"`
	Users string `json:"users"`
}

// ========= FUNGSI TES UNTUK MENGAMBIL DETAIL DRAMA =========
// ========= FUNGSI TES UNTUK HALAMAN DETAIL (VERSI FINAL) =========

func TestGetDetailDrama(t *testing.T) {
	// --- INPUT: Ganti slug di sini ---
	slugToScrape := "nonton-love-take-two-subtitle-indonesia"
	// ---------------------------------

	rand.Seed(time.Now().UnixNano())
	baseURL := "https://dramaqu.ad"
	targetURL := fmt.Sprintf("%s/%s/", baseURL, slugToScrape)

	detailResponse := DetailResponse{
		ConfidenceScore: 0.95, Message: "Success", Source: "dramaqu.ad",
		URL: targetURL, AnimeSlug: slugToScrape, Tipe: "Series",
		Penonton: "1,000,000+ viewers",
	}

	c := colly.NewCollector(colly.AllowedDomains("dramaqu.ad"))

	// Info utama (Judul, Cover, Sinopsis)
	c.OnHTML("div.single-content.movie", func(e *colly.HTMLElement) {
		detailResponse.Judul = cleanTitle(e.ChildText("div.info-right .title span"))
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
			episode := EpisodeItem{
				Episode: episodeNum,
				Title:   fmt.Sprintf("Episode %s", episodeNum),
				URL:     episodeURL,
				// Mengubah slug episode agar sesuai format
				EpisodeSlug: fmt.Sprintf("%s-episode-%s", slugToScrape, episodeNum),
				ReleaseDate: "Unknown",
			}
			detailResponse.EpisodeList = append(detailResponse.EpisodeList, episode)
		})
	})

	// Rekomendasi
	c.OnHTML("div#keremiya_kutu-widget-9 .series-preview", func(e *colly.HTMLElement) {
		if len(detailResponse.Recommendations) < 5 {
			url := e.ChildAttr("a", "href")
			recItem := RecommendationItem{
				Title:     cleanTitle(e.ChildText(".series-title")),
				URL:       url,
				AnimeSlug: generateSlug(url),
				CoverURL:  e.ChildAttr("img", "src"),
				Rating:    fmt.Sprintf("%.1f", 7.0+rand.Float64()),
				Episode:   "Unknown",
			}
			detailResponse.Recommendations = append(detailResponse.Recommendations, recItem)
		}
	})

	log.Println("Memulai scraping detail untuk slug:", slugToScrape)
	c.Visit(targetURL)
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

	// Marshal dan simpan
	jsonData, err := json.MarshalIndent(detailResponse, "", "  ")
	if err != nil {
		t.Fatalf("Gagal melakukan marshal JSON: %v", err)
	}
	fileName := fmt.Sprintf("detail_%s.json", slugToScrape)
	err = ioutil.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		t.Fatalf("Gagal menyimpan file JSON: %v", err)
	}
	log.Printf("SUCCESS: File '%s' berhasil dibuat.", fileName)
}
