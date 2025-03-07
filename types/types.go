package types

import (
	"encoding/json"
	"fmt"
)

type Nickname struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Предположим, что у вас есть тип MatchData:
type MatchData struct {
	Hero     string   `json:"hero"`
	Result   string   `json:"result"`
	KDA      string   `json:"kda"`
	Duration string   `json:"duration "`
	Role     string   `json:"role"`
	Lane     string   `json:"lane"`
	Items    []string `json:"items"`
}

// Message — обобщённый тип, где T может быть любым (например, string или []MatchData).
type Message[T any] struct {
	Role    string `json:"role"`
	Content T      `json:"content"`
}

// Вспомогательная функция для обобщённого Unmarshal сообщения.
func unmarshalMessage[T any](data json.RawMessage) (Message[T], error) {
	var m Message[T]
	err := json.Unmarshal(data, &m)
	return m, err
}

// OpenAIRequest содержит два сообщения: один с текстом, второй с []MatchData.
type OpenAIRequest struct {
	Model          string               `json:"model"`
	TextMessage    Message[string]      `json:"-"`
	MatchesMessage Message[[]MatchData] `json:"-"`
	Temp           float64              `json:"temperature,omitempty"`
}

func (r OpenAIRequest) MarshalJSON() ([]byte, error) {
	// Сначала сериализуем каждое сообщение отдельно
	textMsgBytes, err := json.Marshal(r.TextMessage)
	if err != nil {
		return nil, fmt.Errorf("error marshaling TextMessage: %w", err)
	}
	matchesMsgBytes, err := json.Marshal(r.MatchesMessage)
	if err != nil {
		return nil, fmt.Errorf("error marshaling MatchesMessage: %w", err)
	}

	// Создаем временную структуру с нужной структурой JSON
	tmp := struct {
		Model    string            `json:"model"`
		Messages []json.RawMessage `json:"messages"`
		Temp     float64           `json:"temperature,omitempty"`
	}{
		Model:    r.Model,
		Messages: []json.RawMessage{textMsgBytes, matchesMsgBytes},
		Temp:     r.Temp,
	}

	return json.Marshal(tmp)
}

// Реализуем кастомный UnmarshalJSON для OpenAIRequest, используя generics.
func (r *OpenAIRequest) UnmarshalJSON(data []byte) error {
	// Временная структура для извлечения общих полей.
	var tmp struct {
		Model    string            `json:"model"`
		Messages []json.RawMessage `json:"messages"`
		Temp     float64           `json:"temperature,omitempty"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	if len(tmp.Messages) != 2 {
		return fmt.Errorf("expected exactly 2 messages, got %d", len(tmp.Messages))
	}

	// Распаковываем первое сообщение как Message[string]
	textMsg, err := unmarshalMessage[string](tmp.Messages[0])
	if err != nil {
		return fmt.Errorf("error unmarshaling first message: %w", err)
	}
	// Распаковываем второе сообщение как Message[[]MatchData]
	matchesMsg, err := unmarshalMessage[[]MatchData](tmp.Messages[1])
	if err != nil {
		return fmt.Errorf("error unmarshaling second message: %w", err)
	}

	r.Model = tmp.Model
	r.TextMessage = textMsg
	r.MatchesMessage = matchesMsg
	r.Temp = tmp.Temp

	return nil
}

const LLMContent = `Ты работаешь в сервисе по генерации никнеймов для игроков дота 2. Твоя задача заключается в том, чтобы анализировать матчи и основываясь на том, каким героем он играл больше всего, а так же процент побед и уровень KDA, генерировать никнеймы. Нашим сервисом пользуются обычные игроки, профессионалы и стримеры. Всех их ты заставляешь смеяться от формулировок, точности и абсурда. Придерживайся инструкции:
1. Проанализируй стиль игры (KDA, предметы, роль)
2. Ищи смешные ассоциации (огонь, таверна, бутылка, сетки и т.д.)
3. Используй каламбуры, мемы, сочетания слов
4. Поиграй с уменьшительными, забавными окончаниями
Пример: Если ты много фидишь на Pudge → можно назвать себя "Крюк в таверну".
Если ты делаешь много ассистов на Lion → "АссистенТоп".
Главное — юмор, узнаваемость и связь с стилем игрока!
Никнеймы могут быть обидными и указывать на плохую игру.
Никнеймы могут содержать нецензурную лексику на русском или английском языке.
Также добавь объяснение ника так, как будто ты общаешься с близким другом.
Сгенерируй 10 никнеймов и оформи их в строку вида [{name: предложенный ник, description: объяснение ника}, следующие элементы...]. Присылай только элементы, без лишнего текста. Вот данные о последних 20 матчах:`
