package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maxime-louis14/api-golang/controllers"
)

func UserRoute(app *fiber.App) {
	app.Post("/user", controllers.CreateUser)
	app.Get("/users", controllers.GetAllUsers)
	app.Get("/user/:userId", controllers.GetAUser)
	app.Put("/user/:userId", controllers.EditAUser)
	app.Delete("/user/:userId", controllers.DeleteAUser)
	
}
