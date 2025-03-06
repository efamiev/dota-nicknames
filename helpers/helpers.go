package helpers

import (
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
	"time"

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

type LoggingTransport struct {
	Transport http.RoundTripper
}

func (t *LoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	resp, err := t.Transport.RoundTrip(req)
	if err != nil {
		slog.Error("Ошибка HTTP-запроса", "error", err)
		return nil, err
	}

	duration := time.Since(start)
	slog.Info("HTTP Request",
		"method", req.Method,
		"url", req.URL.String(),
		"status", resp.Status,
		"duration", duration.String(),
	)

	return resp, nil
}
