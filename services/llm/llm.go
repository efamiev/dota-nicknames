package llm

import (
	"bytes"
	"dota-nicknames/helpers"
	"dota-nicknames/types"
	"encoding/json"
	"fmt"
	"io"

	"net/http"
	"os"
)

type Choice struct {
	Message types.Message[string] `json:"message"`
}

type APIResponse struct {
	Choices []Choice `json:"choices"`
}

var apiKey = os.Getenv("API_KEY")

func GenerateNicknames(url string, reqBody []byte) ([]types.Nickname, error) {
	resp, err := sendRequest(url, reqBody)
	defer resp.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("ошибка отправки запроса: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("ошибка чтения ответа LLM API: %w", err)
		}

		return nil, fmt.Errorf("ошибка LLM API (%d): %s", resp.StatusCode, string(body))
	}

	return readResp(resp)
}

func readResp(resp *http.Response) ([]types.Nickname, error) {
	const errMsg = "ошибка чтения ответа LLM API"

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	var data APIResponse
	if err = json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	if len(data.Choices) == 0 {
		return nil, fmt.Errorf("%s: %s", errMsg, body)
	}

	var nicknames []types.Nickname
	if err = json.Unmarshal([]byte(data.Choices[0].Message.Content), &nicknames); err != nil {
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	return nicknames, nil
}

func sendRequest(url string, reqBody []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания HTTP-запроса: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{
		Transport: &helpers.LoggingTransport{Transport: http.DefaultTransport},
	}

	return client.Do(req)
}
