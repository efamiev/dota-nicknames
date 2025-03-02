package internal

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	// "time"
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

type Choice struct {
	Message Message `json:"message"`
}

type APIResponse struct {
	Choices []Choice `json:"choices"`
}

var apiKey = os.Getenv("API_KEY")

func nicks(url string) string {
	matches := FetchMatchData(url)
	matchesJson, _ := json.MarshalIndent(matches[:20], "", " ")

	return string(matchesJson)
}

func GenerateNicknames() ([]string, error) {
	matches := nicks("https://www.dotabuff.com/players/321580662/matches")

	reqBody := OpenAIRequest{
		Model: "deepseek/deepseek-chat:free",
		Messages: []Message{
			{
				Role:    "system",
				Content: "Предложим, ты работаешь в сервисе по генерации никнеймов для игроков дота 2. Твоя задача заключается в том, чтобы анализировать страничку профиля игрока на ресурсе Dotabuff, и, основываяюсь на том, каким героем он играл больше всего за полседние 20 игр, а так же процент побед и уровень KDA, предлагать пользователю на выбор 5 никнеймов в юмористическом стиле (предпочитая абсурдизм и постмодернизм). При генерации никнеймов, нужно так же учитывать локальные мемы русскоязычного сообщества Dota 2",
			},
			{
				Role:    "user",
				Content: matches,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return []string{}, err
	}

	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return []string{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// client := &http.Client{Timeout: time.Second * 30}
	client := &http.Client{}
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

	var data APIResponse

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println("Ошибка парсинга JSON:", err)
		return []string{}, nil
	}

	res := data.Choices[0].Message.Content

	// Регулярное выражение для поиска никнеймов и их описаний
	re := regexp.MustCompile(`\d+\.\s\*\*(.*?)\*\*\s*\((.*?)\)`)

	// Извлекаем все совпадения
	rawNicks := re.FindAllStringSubmatch(res, -1)

	// Заполняем слайс никнеймами и описаниями
	var nicknames []string
	for _, match := range rawNicks {
		if len(match) > 2 {
			nicknames = append(nicknames, match[1]+" - "+match[2])
		}
	}

	log.Println("Извлеченные никнеймы:", nicknames)

	return nicknames, nil
}
