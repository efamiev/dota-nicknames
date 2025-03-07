package types

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var matches = []MatchData{
	{
		Hero:     "Queen of Pain",
		Result:   "Lost Match",
		KDA:      "8/12/20",
		Duration: "46:57",
		Role:     "Support Role",
		Lane:     "Off Lane",
		Items: []string{
			"Magic Wand",
			"Power Treads",
			"Observer Ward",
			"Orchid Malevolence",
			"Aghanim's Scepter",
			"Linken's Sphere",
		},
	},
	{
		Hero:     "Nature's Prophet",
		Result:   "Won Match",
		KDA:      "2/10/19",
		Duration: "37:47",
		Role:     "Support Role",
		Lane:     "Off Lane",
		Items: []string{
			"Power Treads",
			"Null Talisman",
			"Magic Wand",
			"Orchid Malevolence",
			"Mjollnir",
			"Desolator",
		},
	},
}

func TestFetchMatches(t *testing.T) {
	t.Run("Parsing TextMessage and MatchesMessage to Messages array", func(t *testing.T) {
		jsonContent := `{
			"model": "deepseek/deepseek-chat:free",
			"messages": [
			{
				"role": "system",
				"content": "Ты работаешь в сервисе по генерации никнеймов для игроков дота 2. Твоя задача заключается в том, чтобы..."
			},
			{
				"role": "user",
				"content": [
				{
					"hero": "Queen of Pain",
					"result": "Lost Match",
					"kda": "8/12/20",
					"duration ": "46:57",
					"role": "Support Role",
					"lane": "Off Lane",
					"items": [
					"Magic Wand",
					"Power Treads",
					"Observer Ward",
					"Orchid Malevolence",
					"Aghanim's Scepter",
					"Linken's Sphere"
					]
				},
				{
					"hero": "Nature's Prophet",
					"result": "Won Match",
					"kda": "2/10/19",
					"duration ": "37:47",
					"role": "Support Role",
					"lane": "Off Lane",
					"items": [
					"Power Treads",
					"Null Talisman",
					"Magic Wand",
					"Orchid Malevolence",
					"Mjollnir",
					"Desolator"
					]
				}
				]
			}
			]
		}`

		reqBody := OpenAIRequest{
			Model:          "deepseek/deepseek-chat:free",
			TextMessage:    Message[string]{Role: "system", Content: "Ты работаешь в сервисе по генерации никнеймов для игроков дота 2. Твоя задача заключается в том, чтобы..."},
			MatchesMessage: Message[[]MatchData]{Role: "user", Content: matches},
		}

		reqJson, _ := json.Marshal(reqBody)
		var buf bytes.Buffer
		json.Compact(&buf, []byte(jsonContent))

		assert.Equal(t, buf.String(), string(reqJson))
	})
}
