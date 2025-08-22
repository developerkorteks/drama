// File: episode_detail_test.go
package scrape

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

// ========= STRUCT UNTUK DETAIL EPISODE =========

type EpisodeDetailResponse struct {
	ConfidenceScore  float64           `json:"confidence_score"`
	Message          string            `json:"message"`
	Source           string            `json:"source"`
	Title            string            `json:"title"`
	ThumbnailURL     string            `json:"thumbnail_url"`
	StreamingServers []StreamingServer `json:"streaming_servers"`
	ReleaseInfo      string            `json:"release_info"`
	DownloadLinks    DownloadLinks     `json:"download_links"`
	Navigation       Navigation        `json:"navigation"`
	AnimeInfo        AnimeInfo         `json:"anime_info"`
	OtherEpisodes    []OtherEpisode    `json:"other_episodes"`
}

type StreamingServer struct {
	ServerName   string `json:"server_name"`
	StreamingURL string `json:"streaming_url"`
}

type DownloadLinks struct {
	MKV  map[string][]DownloadProvider `json:"MKV"`
	MP4  map[string][]DownloadProvider `json:"MP4"`
	X265 map[string][]DownloadProvider `json:"x265 [Mode Irit Kuota tapi Kualitas Sama Beningnya]"`
}

type DownloadProvider struct {
	Provider string `json:"provider"`
	URL      string `json:"url"`
}

type Navigation struct {
	PreviousEpisodeURL string `json:"previous_episode_url,omitempty"`
	AllEpisodesURL     string `json:"all_episodes_url"`
	NextEpisodeURL     string `json:"next_episode_url,omitempty"`
}

type AnimeInfo struct {
	Title        string   `json:"title"`
	ThumbnailURL string   `json:"thumbnail_url"`
	Synopsis     string   `json:"synopsis"`
	Genres       []string `json:"genres"`
}

type OtherEpisode struct {
	Title        string `json:"title"`
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnail_url"`
	ReleaseDate  string `json:"release_date"`
}

// Struct untuk mem-parsing respons AJAX dari server
type AjaxPlayerResponse struct {
	Success bool `json:"success"`
	Data    struct {
		IframeURL string `json:"iframe_url"`
	} `json:"data"`
}

// ========= FUNGSI TES MENGGUNAKAN COLLY (FINAL) =========
func TestGetEpisodeDetailColly(t *testing.T) {
	episodeURL := "https://dramaqu.ad/nonton-my-girlfriend-is-the-man-subtitle-indonesia/"

	episodeResponse := &EpisodeDetailResponse{
		ConfidenceScore: 1.0,
		Message:         "Success",
		Source:          "dramaqu.ad",
		ReleaseInfo:     fmt.Sprintf("Released on %s %d", time.Now().Month().String(), time.Now().Year()),
		DownloadLinks: DownloadLinks{
			MKV:  make(map[string][]DownloadProvider),
			MP4:  make(map[string][]DownloadProvider),
			X265: make(map[string][]DownloadProvider),
		},
		AnimeInfo:        AnimeInfo{},
		StreamingServers: []StreamingServer{}, // Inisialisasi slice agar tidak null
	}

	c := colly.NewCollector(
		colly.AllowedDomains("dramaqu.ad"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"),
	)
	c.SetRequestTimeout(30 * time.Second)

	c.OnHTML("body", func(e *colly.HTMLElement) {
		log.Println("Mem-parsing HTML dari halaman utama...")
		doc := e.DOM

		// Parsing data statis
		mainContent := doc.Find("div.single-content.movie")
		episodeResponse.Title = mainContent.Find("div.title span").Text() + " " + mainContent.Find("div.release").Text()
		thumbnail := mainContent.Find("div.poster img").AttrOr("src", "")
		episodeResponse.ThumbnailURL = thumbnail
		episodeResponse.AnimeInfo.ThumbnailURL = thumbnail
		episodeResponse.AnimeInfo.Synopsis = strings.TrimSpace(mainContent.Find("div.excerpt").Text())
		mainContent.Find("div.categories a").Each(func(i int, s *goquery.Selection) {
			episodeResponse.AnimeInfo.Genres = append(episodeResponse.AnimeInfo.Genres, s.Text())
		})
		baseDramaURL, currentEpNum := "", "1"
		reEp := regexp.MustCompile(`(https?://[^/]+/[^/]+)/(\d+)/?$`)
		reBase := regexp.MustCompile(`(https?://[^/]+/[^/]+)/?$`)
		if matches := reEp.FindStringSubmatch(episodeURL); len(matches) > 2 {
			baseDramaURL, currentEpNum = matches[1], matches[2]
		} else if matches := reBase.FindStringSubmatch(episodeURL); len(matches) > 1 {
			baseDramaURL = strings.TrimSuffix(matches[1], "/")
		}
		num, _ := strconv.Atoi(currentEpNum)
		if baseDramaURL != "" {
			episodeResponse.Navigation.AllEpisodesURL = baseDramaURL + "/"
			if num > 1 {
				if num == 2 {
					episodeResponse.Navigation.PreviousEpisodeURL = baseDramaURL + "/"
				} else {
					episodeResponse.Navigation.PreviousEpisodeURL = fmt.Sprintf("%s/%d/", baseDramaURL, num-1)
				}
			}
			if doc.Find(fmt.Sprintf("a.post-page-numbers[href*='/%d/']", num+1)).Length() > 0 {
				episodeResponse.Navigation.NextEpisodeURL = fmt.Sprintf("%s/%d/", baseDramaURL, num+1)
			}
		}
		doc.Find("div#action-parts a.post-page-numbers").Each(func(i int, s *goquery.Selection) {
			epNum := s.Find("span").Text()
			epURL := s.AttrOr("href", "")
			episodeResponse.OtherEpisodes = append(episodeResponse.OtherEpisodes, OtherEpisode{
				Title: "Episode " + epNum, URL: epURL, ThumbnailURL: thumbnail, ReleaseDate: "Unknown",
			})
		})
		reTitle := regexp.MustCompile(`(.*)\s+\(Episode\s+\d+\)`)
		matches := reTitle.FindStringSubmatch(episodeResponse.Title)
		if len(matches) > 1 {
			episodeResponse.AnimeInfo.Title = strings.TrimSpace(matches[1])
		} else {
			reClean := regexp.MustCompile(`\s+Episode\s+\d+.*`)
			episodeResponse.AnimeInfo.Title = reClean.ReplaceAllString(cleanTitle(episodeResponse.Title), "")
		}

		// Mendapatkan parameter AJAX
		log.Println("Mencari parameter AJAX...")
		playerID, exists := doc.Find("div.apicodes-container").Attr("id")
		if !exists {
			log.Println("PERINGATAN: Tidak dapat menemukan Player ID.")
			return
		}
		var nonce string
		dataUri, exists := doc.Find("script#dramagu-player-js-extra").Attr("src")
		if !exists {
			log.Println("PERINGATAN: Tidak dapat menemukan script tag '#dramagu-player-js-extra'.")
			return
		}
		parts := strings.SplitN(dataUri, ",", 2)
		if len(parts) != 2 {
			log.Println("PERINGATAN: Format data URI pada script tidak valid.")
			return
		}
		decodedScript, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			log.Printf("PERINGATAN: Gagal men-decode base64: %v", err)
			return
		}
		reNonce := regexp.MustCompile(`"nonce":"(\w+)"`)
		nonceMatches := reNonce.FindStringSubmatch(string(decodedScript))
		if len(nonceMatches) > 1 {
			nonce = nonceMatches[1]
		} else {
			log.Println("PERINGATAN: Gagal mengekstrak nonce dari script yang di-decode.")
			return
		}

		// Kirim permintaan AJAX
		log.Printf("Parameter ditemukan: PlayerID=%s, Nonce=%s", playerID, nonce)
		ajaxURL := e.Request.AbsoluteURL("/wp-admin/admin-ajax.php")
		formData := map[string]string{
			"action":    "get_player_url",
			"player_id": playerID,
			"nonce":     nonce,
		}
		log.Printf("Mengirim permintaan POST ke: %s", ajaxURL)
		if err := c.Post(ajaxURL, formData); err != nil {
			t.Errorf("Gagal mengirimkan permintaan AJAX: %v", err)
		}
	})

	c.OnResponse(func(r *colly.Response) {
		if strings.Contains(r.Request.URL.String(), "admin-ajax.php") {
			log.Println("Menerima respons dari AJAX call.")
			var ajaxResp AjaxPlayerResponse
			if err := json.Unmarshal(r.Body, &ajaxResp); err != nil {
				log.Printf("Gagal mem-parsing JSON dari respons AJAX: %v", err)
				return
			}

			if ajaxResp.Success && ajaxResp.Data.IframeURL != "" {
				iframeSrc := ajaxResp.Data.IframeURL
				log.Println("BERHASIL! Link streaming dari AJAX ditemukan:", iframeSrc)

				// ==========================================================
				// PERUBAHAN DI SINI: HAPUS GIMMICK, GUNAKAN DATA ASLI
				// ==========================================================

				// Dapatkan nama host dari URL untuk dijadikan nama server
				parsedURL, err := url.Parse(iframeSrc)
				serverName := "Default Server" // Fallback
				if err == nil {
					serverName = strings.ReplaceAll(parsedURL.Hostname(), "www.", "")
				}

				// Tambahkan hanya satu server streaming yang valid
				episodeResponse.StreamingServers = append(episodeResponse.StreamingServers, StreamingServer{
					ServerName:   serverName,
					StreamingURL: iframeSrc,
				})

				// Tambahkan juga satu link download sebagai contoh (opsional)
				// Anda bisa mengembangkan logika ini jika menemukan link download asli
				provider := DownloadProvider{Provider: serverName, URL: iframeSrc}
				episodeResponse.DownloadLinks.MKV["720p"] = []DownloadProvider{provider}

			} else {
				log.Println("Respons AJAX tidak berhasil atau URL iframe kosong.")
			}
		}
	})

	c.OnRequest(func(r *colly.Request) { log.Println("Mengunjungi:", r.URL.String()) })
	c.OnError(func(r *colly.Response, err error) { t.Fatalf("Gagal saat request ke %s: %v", r.Request.URL, err) })

	err := c.Visit(episodeURL)
	if err != nil {
		t.Fatalf("Gagal mengunjungi URL: %v", err)
	}

	c.Wait()

	if len(episodeResponse.StreamingServers) == 0 {
		t.Logf("%+v", episodeResponse)
		t.Fatal("Gagal mengambil data penting (Server Streaming). Proses dihentikan.")
	}

	jsonData, err := json.MarshalIndent(episodeResponse, "", "  ")
	if err != nil {
		t.Fatalf("Gagal marshal JSON: %v", err)
	}

	fileName := fmt.Sprintf("episode_colly_%s.json", generateSlug(episodeURL))
	err = ioutil.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		t.Fatalf("Gagal menyimpan file: %v", err)
	}

	log.Printf("SUCCESS: File '%s' berhasil dibuat.", fileName)
}
