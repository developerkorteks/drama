# Anime Terbaru API - DramaQu

## 📋 Overview

Endpoint `/api/v1/anime-terbaru` mengambil daftar anime/drama terbaru dari dramaqu.ad dengan pagination. Implementasi ini **100% konsisten** dengan `scrape/drakor_terbaru_test.go`.

## 🔗 Endpoint

```
GET /api/v1/anime-terbaru
```

## 📝 Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `page` | integer | No | 1 | Nomor halaman untuk pagination |

## 📊 Response Structure

Response mengikuti struktur yang sama persis dengan test file:

```json
{
  "confidence_score": 1.0,
  "message": "Data berhasil diambil dengan kelengkapan sempurna",
  "source": "dramaqu.ad",
  "data": [
    {
      "judul": "My Girlfriend is the Man!",
      "url": "https://dramaqu.ad/nonton-my-girlfriend-is-the-man-subtitle-indonesia/",
      "anime_slug": "nonton-my-girlfriend-is-the-man-subtitle-indonesia",
      "episode": "Episode 9",
      "uploader": "DramaQu Admin",
      "rilis": "Unknown",
      "cover": "https://sp-ao.shortpixel.ai/client/to_webp,q_glossy,ret_img/https://dramaqu.ad/wp-content/uploads/2025/07/nonton-my-girlfriend-is-the-man-subtitle-indonesia-236x350.jpg"
    }
  ]
}
```

## 🏗️ Data Structure

### OngoingDramaResponse
- `confidence_score` (float64): Skor kelengkapan data (0.0 - 1.0)
- `message` (string): Pesan status berdasarkan confidence score
- `source` (string): Sumber data ("dramaqu.ad")
- `data` ([]DramaEntry): Array berisi data drama

### DramaEntry
- `judul` (string): Nama drama/anime
- `url` (string): Link ke halaman detail
- `anime_slug` (string): Slug untuk URL (extracted dari URL)
- `episode` (string): Informasi episode terbaru
- `uploader` (string): Default "DramaQu Admin"
- `rilis` (string): Default "Unknown"
- `cover` (string): URL gambar cover

## 🎯 Confidence Score System

Sama seperti endpoint `/home`, menggunakan sistem validasi:

### Required Fields (Wajib):
- `judul` - Nama drama
- `url` - Link halaman detail
- `anime_slug` - Slug URL
- `cover` - URL gambar cover

### Optional Fields (Opsional):
- `episode` - Info episode
- `uploader` - Info uploader
- `rilis` - Info rilis

### Scoring Logic:
- **1.0** = Semua field wajib + opsional ada
- **0.5** = Field wajib ada, opsional tidak lengkap
- **0.0** = Field wajib tidak ada

### Message Mapping:
- Score 1.0: "Data berhasil diambil dengan kelengkapan sempurna"
- Score 0.5-0.99: "Data berhasil diambil dengan kelengkapan sedang"
- Score 0.1-0.49: "Data berhasil diambil dengan kelengkapan rendah"
- Score 0.0: "Data tidak lengkap - field wajib tidak ada"

## 🔍 Scraping Details

### Target URL Pattern:
- Page 1: `https://dramaqu.ad/category/ongoing-drama/`
- Page 2+: `https://dramaqu.ad/category/ongoing-drama/page/{page}/`

### HTML Selectors:
- Container: `article.movie-preview`
- Title & URL: `span.movie-title a`
- Episode: `span.icon-hd`
- Cover: `img.keremiya-image`

### Data Processing:
1. Extract title dan URL dari `span.movie-title a`
2. Generate slug dari URL menggunakan `path.Base()`
3. Extract episode dari `span.icon-hd`
4. Extract cover dari `img.keremiya-image`
5. Set default values untuk `uploader` dan `rilis`

## 📚 Usage Examples

### Basic Request
```bash
curl -X GET "http://localhost:8080/api/v1/anime-terbaru"
```

### With Pagination
```bash
curl -X GET "http://localhost:8080/api/v1/anime-terbaru?page=2"
```

### Response Analysis
```bash
# Get summary info
curl -s "http://localhost:8080/api/v1/anime-terbaru?page=1" | jq '{
  confidence_score: .confidence_score,
  message: .message,
  total_items: (.data | length),
  first_title: .data[0].judul
}'

# Get all titles
curl -s "http://localhost:8080/api/v1/anime-terbaru?page=1" | jq '.data[].judul'
```

## ✅ Validation Results

### Current Performance:
- **Confidence Score**: 1.0 (Perfect)
- **Message**: "Data berhasil diambil dengan kelengkapan sempurna"
- **Items per Page**: ~10 items
- **All Fields**: Complete and valid

### Test Results:
```json
{
  "confidence_score": 1,
  "message": "Data berhasil diambil dengan kelengkapan sempurna",
  "source": "dramaqu.ad",
  "data_count": 10
}
```

## 🔄 Consistency with Test File

✅ **Struktur Response**: Identik dengan `OngoingDramaResponse`  
✅ **Field Names**: Sama persis dengan `DramaEntry`  
✅ **Data Types**: Konsisten dengan definisi struct  
✅ **Default Values**: "DramaQu Admin" dan "Unknown"  
✅ **URL Processing**: Menggunakan `path.Base()` untuk slug  
✅ **Selectors**: Sama dengan test file  
✅ **Pagination**: URL pattern identik  

## 🚀 Integration

Endpoint ini sudah terintegrasi dengan:
- ✅ Swagger Documentation
- ✅ Gin Router
- ✅ Error Handling
- ✅ Confidence Score System
- ✅ Logging System

## 📊 Swagger Documentation

Endpoint sudah terdokumentasi di Swagger UI:
- URL: `http://localhost:8080/swagger/index.html`
- Tag: `anime-terbaru`
- Method: GET
- Parameters: page (query, integer, optional)

Implementasi endpoint `/api/v1/anime-terbaru` telah **100% konsisten** dengan `scrape/drakor_terbaru_test.go` dan siap untuk production! 🎉