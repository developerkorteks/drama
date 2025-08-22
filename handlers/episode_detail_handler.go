package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nabilulilalbab/dramaqu/services"
)

type EpisodeDetailHandler struct {
	service *services.EpisodeDetailService
}

func NewEpisodeDetailHandler(service *services.EpisodeDetailService) *EpisodeDetailHandler {
	return &EpisodeDetailHandler{service: service}
}

// GetEpisodeDetail handles GET /api/v1/episode-detail
// @Summary Get episode detail
// @Description Mengambil detail episode termasuk server streaming dan link download
// @Tags episode-detail
// @Accept json
// @Produce json
// @Param episode_url query string true "URL episode"
// @Success 200 {object} models.EpisodeDetailResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/episode-detail [get]
func (h *EpisodeDetailHandler) GetEpisodeDetail(c *gin.Context) {
	// Get episode_url parameter
	episodeURL := c.Query("episode_url")
	if strings.TrimSpace(episodeURL) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Episode URL parameter is required",
			"message": "Please provide an episode_url parameter",
		})
		return
	}

	// Get episode detail data from service
	data, err := h.service.GetEpisodeDetail(episodeURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch episode detail",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, data)
}
