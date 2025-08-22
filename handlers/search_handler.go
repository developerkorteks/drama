package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nabilulilalbab/dramaqu/services"
)

type SearchHandler struct {
	service *services.SearchService
}

func NewSearchHandler(service *services.SearchService) *SearchHandler {
	return &SearchHandler{service: service}
}

// SearchDrama handles GET /api/v1/search
// @Summary Search anime
// @Description Mencari anime berdasarkan judul
// @Tags search
// @Accept json
// @Produce json
// @Param query query string true "Query pencarian"
// @Param page query int false "Nomor halaman (default: 1)"
// @Success 200 {object} models.SearchResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/search [get]
func (h *SearchHandler) SearchDrama(c *gin.Context) {
	// Get query parameter
	query := c.Query("query")
	if strings.TrimSpace(query) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Query parameter is required",
			"message": "Please provide a search query",
		})
		return
	}

	// Get page parameter (optional, default to 1)
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid page parameter",
			"message": "Page must be a positive integer",
		})
		return
	}

	// Get search results from service
	data, err := h.service.SearchDrama(query, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch search results",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, data)
}
