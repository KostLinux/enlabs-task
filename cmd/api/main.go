package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"enlabs-task/pkg/config"
	"enlabs-task/pkg/controller"
	"enlabs-task/pkg/database"
	"enlabs-task/pkg/middleware"
	"enlabs-task/pkg/repository"
	"enlabs-task/pkg/service"

	_ "enlabs-task/docs" // Required for swagger docs
)

//	@title			Gambling API
//	@version		1.0
//	@description	API for managing user balances and processing transactions.

//	@contact.name	API Support
//	@contact.email	support@example.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

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
	repos := repository.New(db.DB)
	services := service.New(repos, db.DB)
	controllers := controller.New(services)

	// Initialize router
	router := gin.New()

	// Apply global middlewares
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS(cfg))

	// Add swagger documentation route
	router.Static("/docs", "./docs")

	// Add swagger documentation route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/docs", func(ctx *gin.Context) {
		ctx.File("/docs/index.html")
	})

	// Register routes
	controllers.RegisterRoutes(router)

	// Start the server (simplified approach)
	log.Printf("Starting server on port %s in %s mode", cfg.App.Port, cfg.App.Environment)
	if err := router.Run(":" + cfg.App.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
