package main

import (
	"bytes"
	"dota-nicknames/internal"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Temp     float64   `json:"temperature"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func nicks(url string) string {
	matches := internal.FetchMatchData(url)
	matchesJson, _ := json.MarshalIndent(&matches, "", " ")

	log.Println("NICKS", string(matchesJson))
	return string(matchesJson)
}

func generateNick() ([]string, error) {
	// matches := nicks("https://www.dotabuff.com/players/321580662/matches")

	apiKey := os.Getenv("API_KEY")

	reqBody := OpenAIRequest{
		Model: "deepseek/deepseek-chat:free",
		Messages: []Message{
			{
				Role:    "system",
				Content: "Предложим, ты работаешь в сервисе по генерации никнеймов для игроков дота 2. Твоя задача заключается в том, чтобы анализировать страничку профиля игрока на ресурсе Dotabuff, и, основываяюсь на том, каким героем он играл больше всего за полседние 20 игр, а так же процент побед и уровень KDA, предлагать пользователю на выбор 5 никнеймов в юмористическом стиле (предпочитая абсурдизм и постмодернизм). При генерации никнеймов, нужно так же учитывать локальные мемы русскоязычного сообщества Dota 2",
			},
			{
				Role: "user",
				Content: `[
					{
						"Hero": "Troll Warlord",
						"Result": "Lost Match",
						"KDA": "7/6/13",
						"Duration": "40:25",
						"Role": "Core Role",
						"Lane": "Safe Lane",
						"Items": [
							"Phase Boots",
							"Battle Fury",
							"Manta Style",
							"Black King Bar",
							"Silver Edge",
							"Butterfly"
						]
					}
				]`,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return []string{}, err
	}
	log.Println(string(jsonData))

	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return []string{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: time.Second * 30}
	resp, err := client.Do(req)
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []string{}, err
	}

	// jres, _ := json.Unmarshal(body, "")

	log.Println(string(body))

	return []string{}, nil
}

func main() {
	// url := "https://www.dotabuff.com/players/321580662/matches"

	generateNick()
}
