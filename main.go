package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nabilulilalbab/dramaqu/config"
	"github.com/nabilulilalbab/dramaqu/handlers"
	"github.com/nabilulilalbab/dramaqu/middleware"
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

// @host DYNAMIC_HOST
// @BasePath /
// @schemes http https
func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Set Gin mode based on environment
	if cfg.Environment == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Add middleware for dynamic host detection
	r.Use(middleware.DynamicSwaggerHost())

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

	// Dynamic swagger config endpoint
	r.GET("/swagger-config", middleware.SwaggerConfigHandler())

	// Dynamic swagger docs endpoint
	r.GET("/api/swagger.json", middleware.DynamicSwaggerDocsHandler())

	// Swagger endpoint with dynamic host support
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/api/swagger.json")))

	// Alternative swagger endpoint
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/api/swagger.json")))

	log.Printf("Server starting on :%s", cfg.Port)
	log.Printf("Environment: %s", cfg.Environment)
	if cfg.IsDynamic {
		log.Printf("Swagger Host: Dynamic (will detect from request)")
		log.Printf("Swagger documentation available at: http://[your-domain]/swagger/index.html")
	} else {
		log.Printf("Swagger Host: %s", cfg.SwaggerHost)
		log.Printf("Swagger documentation available at: http://%s/swagger/index.html", cfg.SwaggerHost)
	}

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
