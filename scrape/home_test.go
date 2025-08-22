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

// ========= STRUCT UNTUK OUTPUT JSON FINAL (Sesuai Permintaan) =========

type FinalResponse struct {
	ConfidenceScore float64      `json:"confidence_score"`
	Message         string       `json:"message"`
	Source          string       `json:"source"`
	Top10           []Top10Item  `json:"top10"`
	NewEps          []NewEpsItem `json:"new_eps"`
	Movies          []MovieItem  `json:"movies"`
	JadwalRilis     JadwalRilis  `json:"jadwal_rilis"`
}

type Top10Item struct {
	Judul     string   `json:"judul"`
	URL       string   `json:"url"`
	AnimeSlug string   `json:"anime_slug"`
	Rating    string   `json:"rating"`
	Cover     string   `json:"cover"`
	Genres    []string `json:"genres"`
}

type NewEpsItem struct {
	Judul     string `json:"judul"`
	URL       string `json:"url"`
	AnimeSlug string `json:"anime_slug"`
	Episode   string `json:"episode"`
	Rilis     string `json:"rilis"`
	Cover     string `json:"cover"`
}

type MovieItem struct {
	Judul     string   `json:"judul"`
	URL       string   `json:"url"`
	AnimeSlug string   `json:"anime_slug"`
	Tanggal   string   `json:"tanggal"`
	Cover     string   `json:"cover"`
	Genres    []string `json:"genres"`
}

type JadwalRilis struct {
	Monday    []JadwalItem `json:"Monday"`
	Tuesday   []JadwalItem `json:"Tuesday"`
	Wednesday []JadwalItem `json:"Wednesday"`
	Thursday  []JadwalItem `json:"Thursday"`
	Friday    []JadwalItem `json:"Friday"`
	Saturday  []JadwalItem `json:"Saturday"`
	Sunday    []JadwalItem `json:"Sunday"`
}

type JadwalItem struct {
	Title       string   `json:"title"`
	URL         string   `json:"url"`
	AnimeSlug   string   `json:"anime_slug"`
	CoverURL    string   `json:"cover_url"`
	Type        string   `json:"type"`
	Score       string   `json:"score"`
	Genres      []string `json:"genres"`
	ReleaseTime string   `json:"release_time"`
}

// ========= FUNGSI UTAMA UNTUK SCRAPING DAN GENERATE JSON =========

func TestGenerateFullJSON(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	finalResponse := FinalResponse{
		ConfidenceScore: 1.0,
		Message:         "Data berhasil diambil",
		Source:          "dramaqu.ad", // Diubah sesuai sumber
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
					itemData := parseNewEpsItem(item)
					finalResponse.NewEps = append(finalResponse.NewEps, itemData)
				}
			case "Film Korea":
				if len(finalResponse.Movies) < 20 {
					itemData := parseMovieItem(item)
					finalResponse.Movies = append(finalResponse.Movies, itemData)
				}
			case "Drama Populer":
				if len(finalResponse.Top10) < 10 {
					itemData := parseTop10Item(item)
					finalResponse.Top10 = append(finalResponse.Top10, itemData)
				}
			}
		})
	})

	log.Println("Memulai scraping halaman utama...")
	c.Visit("https://dramaqu.ad/")
	c.Wait()
	log.Println("Scraping halaman utama selesai.")

	scheduleCollector := colly.NewCollector(
		colly.AllowedDomains("dramaqu.ad"),
	)
	var ongoingItemsForSchedule []JadwalItem
	scheduleCollector.OnHTML("article.movie-preview", func(e *colly.HTMLElement) {
		item := parseJadwalItem(e)
		ongoingItemsForSchedule = append(ongoingItemsForSchedule, item)
	})

	log.Println("Memulai scraping halaman 'Ongoing' untuk data jadwal...")
	scheduleCollector.Visit("https://dramaqu.ad/category/ongoing-drama/")
	scheduleCollector.Wait()
	log.Println("Scraping halaman 'Ongoing' selesai.")

	finalResponse.JadwalRilis = generateJadwal(ongoingItemsForSchedule)
	log.Println("Jadwal rilis dummy berhasil dibuat.")

	jsonData, err := json.MarshalIndent(finalResponse, "", "  ")
	if err != nil {
		t.Fatalf("Gagal melakukan marshal JSON: %v", err)
	}

	err = ioutil.WriteFile("dramaqu_response.json", jsonData, 0644)
	if err != nil {
		t.Fatalf("Gagal menyimpan file JSON: %v", err)
	}

	log.Println("SUCCESS: File 'dramaqu_response.json' berhasil dibuat.")
}

// ========= FUNGSI HELPER UNTUK PARSING DAN GENERATE DATA =========

// cleanTitle membersihkan judul dari teks yang tidak diinginkan
// func cleanTitle(rawTitle string) string {
// 	// Hapus "Nonton" dan "Drama Korea Subtitle Indonesia"
// 	re := regexp.MustCompile(`(?i)Nonton\s*|\s*Drama\s*Korea\s*Subtitle\s*Indonesia`)
// 	title := re.ReplaceAllString(rawTitle, "")

// 	// Hapus tahun dalam kurung, contoh: (2024)
// 	re = regexp.MustCompile(`\s*\(\d{4}\)`)
// 	title = re.ReplaceAllString(title, "")

// 	return strings.TrimSpace(title)
// }

// KODE BARU YANG BENAR
// func generateSlug(url string) string {
// 	parts := strings.Split(strings.Trim(url, "/"), "/")
// 	// Cukup ambil bagian terakhir dari URL tanpa modifikasi
// 	return parts[len(parts)-1]
// }

func parseNewEpsItem(e *colly.HTMLElement) NewEpsItem {
	url := e.ChildAttr("a", "href")
	judul := cleanTitle(e.ChildText(".movie-title a"))
	return NewEpsItem{
		Judul:     judul,
		URL:       url,
		AnimeSlug: generateSlug(url),
		Episode:   e.ChildText(".center-icons .icon-hd"),
		Cover:     e.ChildAttr("img", "src"),
		Rilis:     fmt.Sprintf("%d jam", rand.Intn(23)+1), // Dummy
	}
}

func parseMovieItem(e *colly.HTMLElement) MovieItem {
	url := e.ChildAttr("a", "href")
	judul := cleanTitle(e.ChildText(".movie-title a"))
	return MovieItem{
		Judul:     judul,
		URL:       url,
		AnimeSlug: generateSlug(url),
		Cover:     e.ChildAttr("img", "src"),
		Tanggal:   fmt.Sprintf("%d hari", rand.Intn(5)+1),  // Dummy
		Genres:    []string{"Action", "Drama", "Thriller"}, // Dummy
	}
}

func parseTop10Item(e *colly.HTMLElement) Top10Item {
	url := e.ChildAttr("a", "href")
	judul := cleanTitle(e.ChildText(".movie-title a"))
	rating := strings.TrimSpace(e.ChildText(".icon-star.imdb"))
	if rating == "" {
		rating = fmt.Sprintf("%.2f", 7.0+rand.Float64()*2) // Dummy rating 7.00-9.00
	}
	return Top10Item{
		Judul:     judul,
		URL:       url,
		AnimeSlug: generateSlug(url),
		Cover:     e.ChildAttr("img", "src"),
		Rating:    rating,
		Genres:    []string{"Action", "Adventure", "Drama"}, // Dummy
	}
}

func parseJadwalItem(e *colly.HTMLElement) JadwalItem {
	url := e.ChildAttr("a", "href")
	judul := cleanTitle(e.ChildText(".movie-title a"))
	return JadwalItem{
		Title:       judul,
		URL:         url,
		AnimeSlug:   generateSlug(url),
		CoverURL:    e.ChildAttr("img", "src"),
		Type:        "TV",                                                   // Dummy
		Score:       fmt.Sprintf("%.1f", 7.0+rand.Float64()),                // Dummy
		Genres:      []string{"Drama", "Romance", "Comedy"},                 // Dummy
		ReleaseTime: fmt.Sprintf("%02d:%02d", rand.Intn(24), rand.Intn(60)), // Dummy
	}
}

func generateJadwal(items []JadwalItem) JadwalRilis {
	// PENTING: Inisialisasi slice untuk setiap hari agar tidak nil
	jadwal := JadwalRilis{
		Monday:    []JadwalItem{},
		Tuesday:   []JadwalItem{},
		Wednesday: []JadwalItem{},
		Thursday:  []JadwalItem{},
		Friday:    []JadwalItem{},
		Saturday:  []JadwalItem{},
		Sunday:    []JadwalItem{},
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
