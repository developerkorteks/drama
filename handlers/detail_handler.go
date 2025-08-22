package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nabilulilalbab/dramaqu/services"
)

type DetailHandler struct {
	service *services.DetailService
}

func NewDetailHandler(service *services.DetailService) *DetailHandler {
	return &DetailHandler{service: service}
}

// GetAnimeDetail handles GET /api/v1/anime-detail
// @Summary Get anime/movie/series detail
// @Description Mengambil detail anime, film, atau series termasuk episode, sinopsis, dan rekomendasi. Slug dapat berupa 'nama-anime', 'film/nama-film', atau 'series/nama-series'
// @Tags anime-detail
// @Accept json
// @Produce json
// @Param anime_slug query string true "Anime/Movie/Series slug (contoh: 'kobane-2022', 'film/kobane-2022', 'series/legend-of-the-female-general')"
// @Success 200 {object} models.DetailResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/anime-detail [get]
func (h *DetailHandler) GetAnimeDetail(c *gin.Context) {
	// Get anime_slug parameter
	animeSlug := c.Query("anime_slug")
	if strings.TrimSpace(animeSlug) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Anime slug parameter is required",
			"message": "Please provide an anime_slug parameter",
		})
		return
	}

	// Get detail data from service
	data, err := h.service.GetDetailDrama(animeSlug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch anime detail",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, data)
}
