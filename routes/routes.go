package routes

import (
	"Fiber/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Post("/users", handlers.CreateUser)
	app.Get("/users/search", handlers.SearchUsers)
	app.Get("/users/:id", handlers.ReadUser)
	app.Put("/users/:id", handlers.UpdateUser)
	app.Delete("/users/:id", handlers.DeleteUser)
	app.Get("/users", handlers.GetAllUsers)
}
