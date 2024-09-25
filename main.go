package main

import (
	"Fiber/database"
	"Fiber/handlers"
	"Fiber/middleware"
	"Fiber/repository"
	"Fiber/routes"
	"Fiber/services"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"log"
	"os"
)

var logger *zap.Logger

func initLogger() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		log.Fatalf("Can't initialize zap logger: %v", err)
	}
}

func main() {
	// Load environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize logger
	initLogger()
	defer logger.Sync()

	// Get the database path from the environment variable
	databasePath := os.Getenv("LOCAL_DATABASE_PATH")
	if databasePath == "" {
		log.Fatalf("LOCAL_DATABASE_PATH is not set in the environment")
	}

	// Initialize the database
	database.InitDB(databasePath)
	defer database.CloseDB()

	// Create the Fiber app
	app := fiber.New()

	// Use custom logging middleware
	app.Use(middleware.ZapLoggerMiddleware(logger))

	// Initialize the repository
	userRepo := repository.NewBuntDBUserRepository(database.DB, logger)

	// Initialize the service with the repository and logger
	userService := services.NewUserService(userRepo, logger)

	// Initialize the handler with the service and logger
	userHandler := handlers.NewUserHandler(userService, logger)

	// Set up routes
	routes.SetupRoutes(app, userHandler)

	// Start the server using the port defined in the environment
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000" // default to 3000 if APP_PORT is not set
	}

	log.Fatal(app.Listen(":" + port))
}
