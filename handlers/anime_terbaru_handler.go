package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nabilulilalbab/dramaqu/services"
)

// AnimeTerbaruHandler handles anime terbaru related requests
type AnimeTerbaruHandler struct {
	service *services.AnimeTerbaruService
}

// NewAnimeTerbaruHandler creates a new AnimeTerbaruHandler
func NewAnimeTerbaruHandler(service *services.AnimeTerbaruService) *AnimeTerbaruHandler {
	return &AnimeTerbaruHandler{
		service: service,
	}
}

// GetAnimeTerbaru handles GET /api/v1/anime-terbaru
// @Summary Get anime terbaru
// @Description Mengambil daftar anime terbaru dengan pagination
// @Tags anime-terbaru
// @Accept json
// @Produce json
// @Param page query int false "Nomor halaman" default(1)
// @Success 200 {object} models.OngoingDramaResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/anime-terbaru [get]
func (h *AnimeTerbaruHandler) GetAnimeTerbaru(c *gin.Context) {
	// Get page parameter from query, default to 1
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid page parameter",
			"message": "Page must be a positive integer",
		})
		return
	}

	// Get anime terbaru data from service
	data, err := h.service.GetAnimeTerbaru(page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch anime terbaru data",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, data)
}
