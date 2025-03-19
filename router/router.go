package router

import (
	"auth_project/cmd/jwt_auth/ent"
	"auth_project/handler"
	"auth_project/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRputes(app *fiber.App, client *ent.Client) {
	app.Post("/register", handler.Register(client))
	app.Post("/login", handler.Login(client))
	app.Get("/home", middleware.Protected(), handler.Home)
}
