package main

import (
	"Fiber/database"
	"Fiber/handlers"
	"Fiber/middleware"
	"Fiber/repository"
	"Fiber/routes"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"log"
)

var logger *zap.Logger

// initLogger initializes the global logger
func initLogger() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		log.Fatalf("Can't initialize zap logger: %v", err)
	}
}

func main() {
	// Initialize the logger
	initLogger()
	defer logger.Sync()

	// Initialize the database
	database.InitDB()
	defer database.CloseDB()

	// Create a new Fiber app
	app := fiber.New()

	// Set up middleware for logging HTTP requests
	app.Use(middleware.ZapLoggerMiddleware(logger))

	// Initialize the user repository with the database and logger
	userRepo := repository.NewBuntDBUserRepository(database.DB, logger)

	// Initialize the user handler with the repository and logger
	userHandler := handlers.NewUserHandler(userRepo, logger)

	// Set up routes with the user handler
	routes.SetupRoutes(app, userHandler)

	// Start the server and listen on port 3000
	log.Fatal(app.Listen(":3000"))
}
