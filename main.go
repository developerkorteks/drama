package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nabilulilalbab/dramaqu/handlers"
	"github.com/nabilulilalbab/dramaqu/routes"
	"github.com/nabilulilalbab/dramaqu/services"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/nabilulilalbab/dramaqu/docs" // Import generated docs
)

// @title DramaQu API
// @version 1.0
// @description API untuk scraping data drama Korea dari dramaqu.ad
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {
	// Set Gin to release mode in production
	// gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	// Initialize services
	homeService := services.NewHomeService()
	animeTerbaruService := services.NewAnimeTerbaruService()
	movieService := services.NewMovieService()
	scheduleService := services.NewScheduleService()
	searchService := services.NewSearchService()
	detailService := services.NewDetailService()
	episodeDetailService := services.NewEpisodeDetailService()

	// Initialize handlers
	homeHandler := handlers.NewHomeHandler(homeService)
	animeTerbaruHandler := handlers.NewAnimeTerbaruHandler(animeTerbaruService)
	movieHandler := handlers.NewMovieHandler(movieService)
	scheduleHandler := handlers.NewScheduleHandler(scheduleService)
	searchHandler := handlers.NewSearchHandler(searchService)
	detailHandler := handlers.NewDetailHandler(detailService)
	episodeDetailHandler := handlers.NewEpisodeDetailHandler(episodeDetailService)

	// Setup routes
	routes.SetupRoutes(r, homeHandler, animeTerbaruHandler, movieHandler, scheduleHandler, searchHandler, detailHandler, episodeDetailHandler)

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Alternative swagger endpoint
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Server starting on :8080")
	log.Println("Swagger documentation available at: http://localhost:8080/swagger/index.html")

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
