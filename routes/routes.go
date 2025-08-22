package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nabilulilalbab/dramaqu/handlers"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(r *gin.Engine, homeHandler *handlers.HomeHandler, animeTerbaruHandler *handlers.AnimeTerbaruHandler, movieHandler *handlers.MovieHandler, scheduleHandler *handlers.ScheduleHandler, searchHandler *handlers.SearchHandler, detailHandler *handlers.DetailHandler, episodeDetailHandler *handlers.EpisodeDetailHandler) {
	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Home endpoint
		v1.GET("/home", homeHandler.GetHome)

		// Anime terbaru endpoint
		v1.GET("/anime-terbaru", animeTerbaruHandler.GetAnimeTerbaru)

		// Movie endpoint
		v1.GET("/movie", movieHandler.GetMovies)

		// Jadwal rilis endpoint
		v1.GET("/jadwal-rilis", scheduleHandler.GetReleaseSchedule)
		v1.GET("/jadwal-rilis/:day", scheduleHandler.GetScheduleByDay)

		// Search endpoint
		v1.GET("/search", searchHandler.SearchDrama)

		// Detail endpoint
		v1.GET("/anime-detail", detailHandler.GetAnimeDetail)

		// Episode Detail endpoint
		v1.GET("/episode-detail", episodeDetailHandler.GetEpisodeDetail)
	}

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "DramaQu API is running",
		})
	})
}
