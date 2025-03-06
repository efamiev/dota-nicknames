package main

import (
	"context"
	"dota-nicknames/components"
	"dota-nicknames/handlers"
	"dota-nicknames/helpers"
	"dota-nicknames/types"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/gofiber/fiber/v2"
)

var globalCtx = context.Background()

type Request struct {
	URL string `json:"url" validate:"required"`
}

func main() {
	app := fiber.New()
	
	app.Use(pprof.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${ip} - ${method} ${path} -> ${latency}\n",
	}))

	app.Static("/static/", "./static")

	app.Get("/", func(c *fiber.Ctx) error {
		return helpers.Render(c, components.Index())
	})

	app.Get("/:id", func(c *fiber.Ctx) error {
		id, _ := strconv.Atoi(c.Params("id"))

		nicknames := []types.Nickname{
			{Name: "InvokerGod", Description: "Если твои комбо быстрее, чем анимация заклинания, этот ник идеально подойдёт для тебя. 10/10 каток заканчиваются твоими хайлайтами."},
			{Name: "HookMaster", Description: "Ты попадёшь даже в тени деревьев. Враги знают, что если ты на Пудже, лучше не выходить из таверны."},
			{Name: "RampageHunter", Description: "Твой стиль игры — идти только на рампагу. Союзники могут не понять, но в конце игры всё равно дадут лайк."},
			{Name: "CarryOrFeed", Description: "Ты знаешь только два состояния: 20/2 или 2/20. Ты либо легенда, либо мем, но в любом случае тебя запомнят."},
			{Name: "SilentSupport", Description: "Ты никогда не жалуешься, но всегда спасаешь команду. Идеальный саппорт, который понимает игру лучше, чем керри."},
		}

		cache := map[int][]types.Nickname{
			321580662: nicknames,
			1:         {},
		}

		return helpers.Render(c, components.Index(components.List(id, cache[id])))
	})

	app.Get("/sse/:id", handlers.Sse)

	api := app.Group("/api")
	api.Post("/matches", handlers.GetMatches)
	api.Post("/add-task", handlers.AddTask)
	
	log.Fatal(app.Listen(":3000"))
}
