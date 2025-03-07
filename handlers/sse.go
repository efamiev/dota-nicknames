package handlers

import (
	"bufio"
	"bytes"
	"context"
	"dota-nicknames/components"
	"dota-nicknames/services/llm"
	"dota-nicknames/services/parsers"
	"dota-nicknames/types"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"
	"github.com/valyala/fasthttp"
)

var ctx = context.Background()
var ch = cache.New(40*time.Minute, 80*time.Minute)

func Sse(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	id, _ := strconv.Atoi(c.Params("id"))

	matches, err := getMatchData(ch, id, parsers.FetchMatchData)
	if err != nil {
	}

	apiUrl := "https://openrouter.ai/api/v1/chat/completions"
	reqBody, err := json.Marshal(types.OpenAIRequest{
		Model:          "deepseek/deepseek-chat:free",
		TextMessage:    types.Message[string]{Role: "system", Content: types.LLMContent},
		MatchesMessage: types.Message[[]types.MatchData]{Role: "user", Content: matches},
	})
	if err != nil {
	}

	nicks, err := llm.GenerateNicknames(apiUrl, reqBody)
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

		log.Println("SSE завершено для ID:", id, nicks)
	}))

	return nil
}

func getMatchData(c *cache.Cache, id int, fetcher parsers.Fetcher) ([]types.MatchData, error) {
	cm, found := c.Get(strconv.Itoa(id))
	if found {
		return cm.([]types.MatchData), nil
	}

	matches, err := fetcher(id)
	if err != nil {
		return nil, fmt.Errorf("FetchMatchData error %s", err)
	}

	return matches, nil
}

func renderList(id int, nicknames []types.Nickname) *bytes.Buffer {
	var buf bytes.Buffer

	component := components.List(id, nicknames)

	if err := component.Render(ctx, &buf); err != nil {
	}

	return &buf
}
