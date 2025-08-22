# Confidence Score System - DramaQu API

## 📊 Overview

Sistem confidence score menghitung kelengkapan data berdasarkan validasi field wajib dan opsional untuk setiap item yang di-scrape dari dramaqu.ad.

## 🎯 Scoring Logic

### Score Calculation:
- **1.0** = Item memiliki semua field wajib + semua field opsional
- **0.5** = Item memiliki semua field wajib tetapi field opsional tidak lengkap
- **0.0** = Item tidak memiliki field wajib (judul, url, anime_slug, cover)

### Final Confidence Score:
```
confidence_score = total_valid_items / total_items
```

## 📋 Field Validation Rules

### 🔴 Required Fields (Wajib)
Jika salah satu field ini kosong/tidak ada, item mendapat score **0.0**:

#### For All Items:
- `judul` / `title` - Nama drama/film
- `url` - Link ke halaman detail
- `anime_slug` - Slug untuk URL
- `cover` / `cover_url` - URL gambar cover

### 🟡 Optional Fields (Opsional)
Field ini mempengaruhi score antara 0.5 dan 1.0:

#### Top10Item:
- `rating` - Rating drama
- `genres` - Array genre

#### NewEpsItem:
- `episode` - Nomor episode
- `rilis` - Waktu rilis

#### MovieItem:
- `tanggal` - Tanggal rilis
- `genres` - Array genre

#### JadwalItem:
- `type` - Tipe konten (TV/Movie)
- `score` - Score rating
- `genres` - Array genre
- `release_time` - Waktu rilis

## 📈 Message System

Berdasarkan confidence score, API akan mengembalikan message yang sesuai:

### Score 0.0:
```json
{
  "confidence_score": 0.0,
  "message": "Data tidak lengkap - field wajib tidak ada"
}
```

### Score < 0.5:
```json
{
  "confidence_score": 0.25,
  "message": "Data berhasil diambil dengan kelengkapan rendah"
}
```

### Score < 1.0:
```json
{
  "confidence_score": 0.75,
  "message": "Data berhasil diambil dengan kelengkapan sedang"
}
```

### Score = 1.0:
```json
{
  "confidence_score": 1.0,
  "message": "Data berhasil diambil dengan kelengkapan sempurna"
}
```

## 🧮 Calculation Examples

### Example 1: Perfect Data
```
Total Items: 50
- 40 items with all required + optional fields = 40 × 1.0 = 40.0
- 10 items with required fields only = 10 × 0.5 = 5.0
- 0 items missing required fields = 0 × 0.0 = 0.0

Final Score: (40.0 + 5.0 + 0.0) / 50 = 0.90
Message: "Data berhasil diambil dengan kelengkapan sedang"
```

### Example 2: Missing Required Fields
```
Total Items: 30
- 20 items with all fields = 20 × 1.0 = 20.0
- 5 items with required only = 5 × 0.5 = 2.5
- 5 items missing required fields = 5 × 0.0 = 0.0

Final Score: (20.0 + 2.5 + 0.0) / 30 = 0.75
Message: "Data berhasil diambil dengan kelengkapan sedang"
```

### Example 3: All Perfect
```
Total Items: 25
- 25 items with all fields = 25 × 1.0 = 25.0

Final Score: 25.0 / 25 = 1.00
Message: "Data berhasil diambil dengan kelengkapan sempurna"
```

## 🔍 Validation Functions

### Core Validation Functions:

1. **`hasRequiredFields(fields ...string)`**
   - Checks if all required fields are present and not empty
   - Returns `false` if any field is empty or whitespace-only

2. **`isTop10ItemValid(item)`**
   - Validates Top10Item with all required + optional fields
   - Required: judul, url, anime_slug, cover
   - Optional: rating, genres

3. **`isNewEpsItemValid(item)`**
   - Validates NewEpsItem with all required + optional fields
   - Required: judul, url, anime_slug, cover
   - Optional: episode, rilis

4. **`isMovieItemValid(item)`**
   - Validates MovieItem with all required + optional fields
   - Required: judul, url, anime_slug, cover
   - Optional: tanggal, genres

5. **`isJadwalItemValid(item)`**
   - Validates JadwalItem with all required + optional fields
   - Required: title, url, anime_slug, cover_url
   - Optional: type, score, genres, release_time

## 🎯 Benefits

1. **Quality Assurance**: Memastikan data yang dikembalikan memiliki kualitas yang baik
2. **Monitoring**: Dapat memantau kualitas scraping dari website sumber
3. **Debugging**: Membantu identifikasi masalah dalam proses scraping
4. **User Experience**: User dapat mengetahui tingkat kelengkapan data yang diterima

## 📊 Real-time Testing

```bash
# Test confidence score
curl -s http://localhost:8080/api/v1/home | jq '.confidence_score, .message'

# Test data completeness
curl -s http://localhost:8080/api/v1/home | jq '{
  confidence_score: .confidence_score,
  message: .message,
  total_items: ((.top10 | length) + (.new_eps | length) + (.movies | length) + (.jadwal_rilis | to_entries | map(.value | length) | add))
}'
```

## ✅ Current Performance

Berdasarkan testing terbaru:
- **Confidence Score**: 1.0 (Perfect)
- **Message**: "Data berhasil diambil dengan kelengkapan sempurna"
- **Data Quality**: Semua field wajib dan opsional terisi dengan baik

Sistem confidence score telah berhasil diimplementasikan dan berfungsi dengan baik! 🎉