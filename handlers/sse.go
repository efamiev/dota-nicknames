package handlers

import (
	"bufio"
	"bytes"
	"context"
	"dota-nicknames/components"
	"dota-nicknames/types"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

var ctx = context.Background()

func Sse(c *fiber.Ctx) error {
	log.Println("SSE", c.Params("id"))
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	id, _ := strconv.Atoi(c.Params("id"))
	strId := c.Params("id")

	var buf bytes.Buffer

	nicknames := []types.Nickname{
		{Name: "InvokerGod", Description: "Если твои комбо быстрее, чем анимация заклинания, этот ник идеально подойдёт для тебя. 10/10 каток заканчиваются твоими хайлайтами."},
		{Name: "HookMaster", Description: "Ты попадёшь даже в тени деревьев. Враги знают, что если ты на Пудже, лучше не выходить из таверны."},
		{Name: "RampageHunter", Description: "Твой стиль игры — идти только на рампагу. Союзники могут не понять, но в конце игры всё равно дадут лайк."},
		{Name: "CarryOrFeed", Description: "Ты знаешь только два состояния: 20/2 или 2/20. Ты либо легенда, либо мем, но в любом случае тебя запомнят."},
		{Name: "SilentSupport", Description: "Ты никогда не жалуешься, но всегда спасаешь команду. Идеальный саппорт, который понимает игру лучше, чем керри."},
	}

	// Создаём HTML-компонент и рендерим в буфер с контекстом
	component := components.List(id, nicknames)
	component.Render(ctx, &buf)

	c.Status(fiber.StatusOK).Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		time.Sleep(5 * time.Second)
		fmt.Fprintf(w, "event: nicknames-%s\ndata: %s\n\n", strId, buf.String())

		err := w.Flush()
		if err != nil {
			// Refreshing page in web browser will establish a new
			// SSE connection, but only (the last) one is alive, so
			// dead connections must be closed here.
			fmt.Printf("Error while flushing: %v. Closing http connection.\n", err)
		}

		log.Println("SSE завершено для ID:", strId)
	}))

	return nil
}
