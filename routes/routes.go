package routes

import (
	"Fiber/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, userHandler *handlers.UserHandler) {

	app.Post("/users", userHandler.CreateUser)
	app.Get("/users/search", userHandler.SearchUsers)
	app.Get("/users/:id", userHandler.GetUserByID)
	app.Get("/users", userHandler.GetAllUsers)
	app.Put("/users/:id", userHandler.UpdateUser)
	app.Delete("/users/:id", userHandler.DeleteUser)

}
