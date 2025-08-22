package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nabilulilalbab/dramaqu/services"
)

// ScheduleHandler handles schedule related requests
type ScheduleHandler struct {
	service *services.ScheduleService
}

// NewScheduleHandler creates a new ScheduleHandler
func NewScheduleHandler(service *services.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{
		service: service,
	}
}

// GetReleaseSchedule handles GET /api/v1/jadwal-rilis
// @Summary Get jadwal rilis
// @Description Mengambil jadwal rilis anime per hari
// @Tags jadwal-rilis
// @Accept json
// @Produce json
// @Success 200 {object} models.ReleaseScheduleResponse
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/jadwal-rilis [get]
func (h *ScheduleHandler) GetReleaseSchedule(c *gin.Context) {
	// Get release schedule data from service
	data, err := h.service.GetReleaseSchedule()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch release schedule data",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, data)
}

// GetScheduleByDay handles GET /api/v1/jadwal-rilis/{day}
// @Summary Get jadwal rilis by day
// @Description Mengambil jadwal rilis anime untuk hari tertentu
// @Tags jadwal-rilis
// @Accept json
// @Produce json
// @Param day path string true "Nama hari (monday, tuesday, wednesday, thursday, friday, saturday, sunday)"
// @Success 200 {object} models.ScheduleByDayResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/jadwal-rilis/{day} [get]
func (h *ScheduleHandler) GetScheduleByDay(c *gin.Context) {
	// Get day parameter from URL path
	day := c.Param("day")
	if day == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Day parameter is required",
			"message": "Please provide a valid day (monday, tuesday, wednesday, thursday, friday, saturday, sunday)",
		})
		return
	}

	// Validate day parameter
	validDays := []string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}
	isValidDay := false
	for _, validDay := range validDays {
		if strings.EqualFold(day, validDay) {
			isValidDay = true
			break
		}
	}

	if !isValidDay {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid day parameter",
			"message": "Day must be one of: monday, tuesday, wednesday, thursday, friday, saturday, sunday",
		})
		return
	}

	// Get schedule data for specific day from service
	data, err := h.service.GetScheduleByDay(day)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch schedule data for the day",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, data)
}
