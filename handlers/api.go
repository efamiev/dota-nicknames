package handlers

import (
	"dota-nicknames/helpers"
	// "dota-nicknames/services/parsers"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Request struct {
	URL string `json:"url" validate:"required"`
}

// func GetMatches(c *fiber.Ctx) error {
// 	var req Request
//
// 	if err := c.BodyParser(&req); err != nil {
// 		return err
// 	}
//
// 	if req.URL == "" {
// 		return c.Status(400).JSON(fiber.Map{
// 			"status":  "error",
// 			"message": "URL field is required",
// 		})
// 	}
//
// 	if !helpers.ValidateUrl(req.URL) {
// 		return c.Status(400).JSON(fiber.Map{
// 			"status":  "error",
// 			"message": "Invalid URL format. Expected format: https://www.dotabuff.com/players/{number}/matches",
// 		})
// 	}
//
// 	matches, _ := parsers.FetchMatchData(req.URL)
//
// 	return c.JSON(matches)
// }

func AddTask(c *fiber.Ctx) error {
	var req Request

	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if req.URL == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "URL field is required",
		})
	}

	if !helpers.ValidateUrl(req.URL) {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid URL format. Expected format: https://www.dotabuff.com/players/{number}/matches",
		})
	}

	id, _ := helpers.ExtractPlayerID(req.URL)

	c.Set("HX-Location", "/"+fmt.Sprint(id))

	return c.JSON(map[string]int{"id": id})
}
