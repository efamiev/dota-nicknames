package main

import (
	"log"
	"regexp"

	"github.com/gofiber/fiber/v3"
)

type Request struct {
	URL string `json:"url" validate:"required"`
}

var dotaBuffRegex = regexp.MustCompile(`^https://www\.dotabuff\.com/players/\d+/matches$`)

func main() {
	app := fiber.New()

	app.Post("/stats", func(c fiber.Ctx) error {
		var req Request

		if err := c.Bind().Body(&req); err != nil {
			return err
		}

		log.Println("Req", req.URL)

		if req.URL == "" {
			return c.Status(400).JSON(fiber.Map{
				"status":  "error",
				"message": "URL field is required",
			})
		}

		if !dotaBuffRegex.MatchString(req.URL) {
			return c.Status(400).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid URL format. Expected format: https://www.dotabuff.com/players/{number}/matches",
			})
		}

		return c.JSON(matchData(req.URL))
	})

	log.Fatal(app.Listen(":3000"))
}
