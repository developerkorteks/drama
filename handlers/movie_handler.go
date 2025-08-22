package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nabilulilalbab/dramaqu/services"
)

// MovieHandler handles movie related requests
type MovieHandler struct {
	service *services.MovieService
}

// NewMovieHandler creates a new MovieHandler
func NewMovieHandler(service *services.MovieService) *MovieHandler {
	return &MovieHandler{
		service: service,
	}
}

// GetMovies handles GET /api/v1/movie
// @Summary Get movies
// @Description Mengambil daftar film dengan pagination
// @Tags movie
// @Accept json
// @Produce json
// @Param page query int false "Nomor halaman" default(1)
// @Success 200 {object} models.DramaListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/movie [get]
func (h *MovieHandler) GetMovies(c *gin.Context) {
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

	// Get movie data from service
	data, err := h.service.GetMovies(page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch movie data",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, data)
}
