package main

import (
	"fmt"

	"github.com/emilio-riv20/proyecto1/Analyzer"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())

	app.Static("/", "./api")

	app.Get("/api", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"data": "Comandos",
		})
	})

	app.Post("/api", func(c *fiber.Ctx) error {
		body := string(c.Body())
		fmt.Println(body)
		// Llama a la función Analyzer y maneja la respuesta
		resultado, err := Analyzer.Analyzer(body)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"resultado": resultado,
		})
	})

	app.Listen(":3000")
	fmt.Println("Server on port 3000")
}
