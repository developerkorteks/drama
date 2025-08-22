// File: utils.go
package scrape

import (
	"regexp"
	"strings"
)

// CleanTitle membersihkan judul dari teks yang tidak diinginkan
func CleanTitle(rawTitle string) string {
	re := regexp.MustCompile(`(?i)Nonton\s*|\s*Subtitle\s*Indonesia.*|\s*Drama\s*Korea\s*Subtitle\s*Indonesia`)
	title := re.ReplaceAllString(rawTitle, "")

	// Hapus tahun dalam kurung, contoh: (2024)
	re = regexp.MustCompile(`\s*\(\d{4}\)`)
	title = re.ReplaceAllString(title, "")

	return strings.TrimSpace(title)
}

// GenerateSlug mengambil bagian terakhir dari URL path.
func GenerateSlug(url string) string {
	parts := strings.Split(strings.Trim(url, "/"), "/")
	return parts[len(parts)-1]
}
