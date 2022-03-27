package routes

import (
	"github.com/LeoRosaGarcia/go-api-rest/controllers"
	"github.com/gofiber/fiber/v2"
)

func TodoContentsRoute(route fiber.Router) {
	route.Get("/", controllers.GetAllTodoContents)
	route.Get("/:id", controllers.GetTodoContent)
	route.Post("/", controllers.AddTodoContent)
	route.Put("/:id", controllers.UpdateTodoContent)
	route.Delete("/:id", controllers.DeleteTodoContent)
}
