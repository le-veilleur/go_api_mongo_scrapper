package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maxime-louis14/api-golang/configs" //add this
	"github.com/maxime-louis14/api-golang/routes"
)

func main() {
	app := fiber.New()

	//run database
	configs.ConnectDB()

	//routes
	routes.UserRoute(app) //add this

	app.Listen(":3000")
}
