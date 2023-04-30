package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maxime-louis14/api-golang/controllers"
)

func RecetteRoute(app *fiber.App) {
	app.Get("/recettes", controllers.GetRecettes)
}
