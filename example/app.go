package example

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New(fiber.Config{})

	app.Listen(":8080")
}
