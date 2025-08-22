# DramaQu API

API untuk scraping data drama Korea dari dramaqu.ad menggunakan Go dan Gin framework.

## Features

- **Home Endpoint**: Mengambil data homepage termasuk top 10 drama, episode terbaru, film terbaru, dan jadwal rilis
- **Swagger Documentation**: Dokumentasi API yang lengkap dan interaktif
- **Structured Response**: Response JSON yang terstruktur dan konsisten

## Installation

1. Clone repository ini
2. Install dependencies:
```bash
go mod tidy
```

3. Generate Swagger documentation:
```bash
swag init
```

4. Run the server:
```bash
go run main.go
```

Server akan berjalan di `http://localhost:8080`

## API Endpoints

### GET /api/v1/home

Mengambil data homepage termasuk:
- Top 10 drama populer
- Episode terbaru (ongoing drama)
- Film Korea terbaru
- Jadwal rilis mingguan

**Response Structure:**
```json
{
  "confidence_score": 1.0,
  "message": "Data berhasil diambil",
  "source": "dramaqu.ad",
  "top10": [
    {
      "judul": "Drama Title",
      "url": "https://dramaqu.ad/drama-url",
      "anime_slug": "drama-slug",
      "rating": "8.5",
      "cover": "https://image-url.jpg",
      "genres": ["Action", "Adventure", "Drama"]
    }
  ],
  "new_eps": [
    {
      "judul": "Drama Title",
      "url": "https://dramaqu.ad/drama-url",
      "anime_slug": "drama-slug",
      "episode": "Episode 12",
      "rilis": "2 jam",
      "cover": "https://image-url.jpg"
    }
  ],
  "movies": [
    {
      "judul": "Movie Title",
      "url": "https://dramaqu.ad/movie-url",
      "anime_slug": "movie-slug",
      "tanggal": "3 hari",
      "cover": "https://image-url.jpg",
      "genres": ["Action", "Drama", "Thriller"]
    }
  ],
  "jadwal_rilis": {
    "Monday": [...],
    "Tuesday": [...],
    "Wednesday": [...],
    "Thursday": [...],
    "Friday": [...],
    "Saturday": [...],
    "Sunday": [...]
  }
}
```

### GET /health

Health check endpoint untuk memastikan API berjalan dengan baik.

## Documentation

Swagger documentation tersedia di: `http://localhost:8080/swagger/index.html`

## Project Structure

```
├── docs/                 # Generated Swagger documentation
├── handlers/            # HTTP request handlers
├── models/              # Data structures/models
├── routes/              # Route definitions
├── scrape/              # Original scraping logic and utilities
├── services/            # Business logic layer
├── main.go              # Application entry point
└── README.md           # This file
```

## Development

### Adding New Endpoints

1. Create model in `models/` directory
2. Create service in `services/` directory
3. Create handler in `handlers/` directory
4. Add route in `routes/routes.go`
5. Add Swagger comments to handler
6. Regenerate docs: `swag init`

### Testing

Run the API and test endpoints:
```bash
# Start server
go run main.go

# Test health endpoint
curl http://localhost:8080/health

# Test home endpoint
curl http://localhost:8080/api/v1/home
```

## Notes

- API menggunakan scraping real-time dari dramaqu.ad
- Response time tergantung pada kecepatan website target
- Beberapa data menggunakan dummy values untuk konsistensi struktur JSON
- Pastikan koneksi internet stabil untuk scraping yang optimal