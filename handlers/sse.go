package handlers

import (
	"bufio"
	"bytes"
	"context"
	"dota-nicknames/components"
	"dota-nicknames/services/llm"
	"dota-nicknames/types"
	"fmt"
	"log"
	"strconv"

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

	nicks, err := llm.GenerateNicknames(id)
	if err != nil {
		// Добавить отправку и отображение ошибки на фронте
		log.Println("GenerateNicknames error", err)
	}

	listComponent := renderList(id, nicks)

	c.Status(fiber.StatusOK).Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		fmt.Fprintf(w, "event: nicknames-%d\ndata: %s\n\n", id, listComponent.String())

		if err := w.Flush(); err != nil {
			fmt.Printf("Error while flushing: %v. Closing http connection.\n", err)
		}

		log.Println("SSE завершено для ID:", id)
	}))

	return nil
}

func renderList(id int, nicknames []types.Nickname) *bytes.Buffer {
	var buf bytes.Buffer

	component := components.List(id, nicknames)
	// Добавить обработку ошибки
	component.Render(ctx, &buf)

	return &buf
}
