package main

import (
	"Fiber/database"
	"Fiber/routes"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	// DB start.
	database.InitDB()
	defer database.CloseDB()

	// Fiber framework start.
	app := fiber.New()

	// Route starting.
	routes.SetupRoutes(app)

	// Server run.
	log.Fatal(app.Listen(":3000"))
}
