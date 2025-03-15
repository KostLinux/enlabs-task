package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"enlabs-task/pkg/config"
	"enlabs-task/pkg/controller"
	"enlabs-task/pkg/database"
	"enlabs-task/pkg/middleware"
	"enlabs-task/pkg/repository"
	"enlabs-task/pkg/service"
)

func main() {
	// Load environment variables from .env file if it exists
	_ = godotenv.Load()

	// Initialize configuration
	cfg := config.New()
	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database connection
	db, err := database.NewPostgresConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize repository, service, and controller layers
	repos := repository.NewRepositoryManager(db.DB)
	services := service.NewServices(repos, db.DB)
	controllers := controller.NewController(services)

	// Initialize router
	router := gin.New()

	// Apply global middlewares
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS(cfg))

	// Register routes
	controllers.RegisterRoutes(router)

	// Start the server (simplified approach)
	log.Printf("Starting server on port %s in %s mode", cfg.App.Port, cfg.App.Environment)
	if err := router.Run(":" + cfg.App.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
