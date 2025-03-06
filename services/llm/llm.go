package llm

import (
	"bytes"
	"dota-nicknames/helpers"
	"dota-nicknames/services/parsers"
	"dota-nicknames/types"
	"fmt"

	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Temp     float64   `json:"temperature,omitempty"`
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

type APIError struct {
	Error struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	} `json:"error"`
}

var apiKey = os.Getenv("API_KEY")

func GenerateNicknames(id int) ([]types.Nickname, error) {
	matches, err := fetchMatches(id)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения матчей: %w", err)
	}
	log.Printf("Для аккаунта %d получены матчи %s", id, matches)

	reqBody, err := json.Marshal(OpenAIRequest{
		Model: "deepseek/deepseek-chat:free",
		Messages: []Message{
			{Role: "system",	Content: types.LLMContent},
			{Role: "user", Content: string(matches)},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка сериализации запроса: %w", err)
	}

	body, err := callLLM(reqBody)
	if err != nil {
		return nil, err
	}

	var data APIResponse

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	var nicknames []types.Nickname

	// Добавить проверку на длину data.Choices
	json.Unmarshal([]byte(data.Choices[0].Message.Content), &nicknames)

	return nicknames, nil
}

func fetchMatches(id int) ([]byte, error) {
	url := fmt.Sprintf("https://www.dotabuff.com/players/%d/matches", id)
	
	matches, err := parsers.FetchMatchData(url)
	if matches == nil || err != nil {
		return nil, fmt.Errorf("не удалось получить матчи для ID %d %w", id, err)
	}
	
	matchesJson, err := json.Marshal(matches)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить матчи для ID %d %w", id, err)
	}

	return matchesJson, nil
}

func callLLM(jsonData []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания HTTP-запроса: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{
		Transport: &helpers.LoggingTransport{Transport: http.DefaultTransport},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка отправки запроса: %w", err)
	}
	defer resp.Body.Close()

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа LLM API: %w", err)
	}

	log.Printf("Ответ от LLM %s", body)

	// Обрабатываем статус код
	if resp.StatusCode != http.StatusOK {
		var apiErr APIError
		if err := json.Unmarshal(body, &apiErr); err != nil {
			return nil, fmt.Errorf("ошибка LLM API (%d): %s", resp.StatusCode, string(body))
		}
		return nil, fmt.Errorf("ошибка LLM API (%d): %s", apiErr.Error.Code, apiErr.Error.Message)
	}

	return body, nil
}
