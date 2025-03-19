package main

import (
	"auth_project/database"
	"auth_project/router"

	"github.com/gofiber/fiber/v2"
)

func main() {

	// Create the Fiber App Instance
	app := fiber.New(fiber.Config{
		AppName: "jwo_auth",
	})
	client := database.DBConnect()
	router.SetupRputes(app, client)
	app.Listen(":3000")
}
