package helpers

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

func Render(c *fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")

	return component.Render(c.Context(), c.Response().BodyWriter())
}

func ValidateUrl(url string) bool {
	return regexp.MustCompile(`^https://www\.dotabuff\.com/players/\d+/matches$`).MatchString(url)
}

func ExtractPlayerID(url string) (int, error) {
	re := regexp.MustCompile(`https://www\.dotabuff\.com/players/(\d+)/matches`)

	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		return 0, fmt.Errorf("ID не найден")
	}

	id, _ := strconv.Atoi(matches[1])

	return id, nil
}
