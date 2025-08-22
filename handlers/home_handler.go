package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nabilulilalbab/dramaqu/services"
)

// HomeHandler handles home-related HTTP requests
type HomeHandler struct {
	homeService *services.HomeService
}

// NewHomeHandler creates a new instance of HomeHandler
func NewHomeHandler(homeService *services.HomeService) *HomeHandler {
	return &HomeHandler{
		homeService: homeService,
	}
}

// GetHome godoc
// @Summary Get homepage data
// @Description Mengambil data homepage termasuk top 10 anime, episode terbaru, film terbaru, dan jadwal rilis
// @Tags Home
// @Accept json
// @Produce json
// @Success 200 {object} models.FinalResponse
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/home [get]
func (h *HomeHandler) GetHome(c *gin.Context) {
	data, err := h.homeService.GetHomeData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch home data",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, data)
}
