package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"dota-nicknames/components"
	"dota-nicknames/internal"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

var globalCtx = context.Background()

type Request struct {
	URL string `json:"url" validate:"required"`
}

func Render(c *fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	return component.Render(c.Context(), c.Response().BodyWriter())
}

var dotaBuffRegex = regexp.MustCompile(`^https://www\.dotabuff\.com/players/\d+/matches$`)

func extractPlayerID(url string) (int, error) {
	re := regexp.MustCompile(`https://www\.dotabuff\.com/players/(\d+)/matches`)

	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		return 0, fmt.Errorf("ID не найден")
	}

	id, _ := strconv.Atoi(matches[1])

	return id, nil
}

func main() {
	app := fiber.New()

	app.Static("/static/", "./static")

	app.Get("/", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")

		return Render(c, components.Index(components.Form(), components.List(0, []string{})))
	})

	app.Get("/nicknames/:id", func(c *fiber.Ctx) error {
		id, _ := strconv.Atoi(c.Params("id"))
		cache := map[int][]string{
			1: {"Бобр добр", "Саня Лох"},
			5: {"sdsd", "asdsdsd"},
		}

		return Render(c, components.Index(components.Form(), components.List(id, cache[id])))
	})

	app.Get("/sse/:id", func(c *fiber.Ctx) error {
		log.Println("SSE")
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")

		id := c.Params("id")

		c.Status(fiber.StatusOK).Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
			var i int
			for {
				i++
				msg := fmt.Sprintf("%d - the time is %v", i, time.Now())
				fmt.Fprintf(w, "event: nicknames-%s\ndata: <li>%s</li>\n\n", id, "Текст"+id)
				fmt.Println(msg)

				err := w.Flush()
				if err != nil {
					// Refreshing page in web browser will establish a new
					// SSE connection, but only (the last) one is alive, so
					// dead connections must be closed here.
					fmt.Printf("Error while flushing: %v. Closing http connection.\n", err)

					break
				}
				time.Sleep(2 * time.Second)
			}
		}))

		return nil
	})

	app.Post("/generate", func(c *fiber.Ctx) error {
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

		if !dotaBuffRegex.MatchString(req.URL) {
			return c.Status(400).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid URL format. Expected format: https://www.dotabuff.com/players/{number}/matches",
			})
		}

		id, _ := extractPlayerID(req.URL)

		c.Set("HX-Location", "/nicknames/"+fmt.Sprint(id))
		return c.SendStatus(200)
	})

	app.Post("/api/matches", func(c *fiber.Ctx) error {
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

		if !dotaBuffRegex.MatchString(req.URL) {
			return c.Status(400).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid URL format. Expected format: https://www.dotabuff.com/players/{number}/matches",
			})
		}

		return c.JSON(internal.FetchMatchData(req.URL))
	})

	log.Fatal(app.Listen(":3000"))
}
