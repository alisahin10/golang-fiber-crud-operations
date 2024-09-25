package main

import (
	"Fiber/database"
	"Fiber/handlers"
	"Fiber/middleware" // Ensure this import is correct
	"Fiber/routes"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"log"
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
	initLogger()
	defer logger.Sync()

	database.InitDB()
	defer database.CloseDB()

	app := fiber.New()

	// Pass the logger to handlers package
	handlers.SetLogger(logger)

	// Register middleware
	app.Use(middleware.ZapLoggerMiddleware(logger))

	// Route starting.
	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
