package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vivekv96/auth/handlers"
)

// Setup registers all the routes with their respective fiber.HandlerFuncs
func Setup(app *fiber.App) {
	app.Post("/api/register", handlers.Register)
	app.Post("/api/login", handlers.Login)
	app.Get("/api/user", handlers.User)
	app.Post("/api/logout", handlers.Logout)
	app.Post("/api/forgot-password", handlers.Forgot)
	app.Post("/api/reset", handlers.Reset)
}
