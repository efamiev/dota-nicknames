package llm

import (
	"dota-nicknames/types"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var matches = []types.MatchData{
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
	// t.Run("Генерируем никнеймы", func(t *testing.T) {
	// 	url := "https://openrouter.ai/api/v1/chat/completions"
	// 	reqBody, _ := json.Marshal(types.OpenAIRequest{
	// 		Model: "deepseek/deepseek-chat:free",
	// 		TextMessage: types.Message[string]{Role: "system", Content: types.LLMContent},
	// 		MatchesMessage: types.Message[[]types.MatchData]{Role: "user", Content: matches},
	// 	})
	//
	// 	GenerateNicknames(url, reqBody)
	// 	assert.Equal(t, 1,2)
	// })

	t.Run("Обрабатываем ошибку от АПИ", func(t *testing.T) {
		apiError := "{ \"error\": { \"message\": \"Provider returned error\", \"code\": 400, \"metadata\": { \"raw\": \"status 400: err targon only supports basic chat requests with `role:string` and `content:string`\", \"provider_name\": \"Targon\", \"isDownstreamPipeClean\": true, \"isErrorUpstreamFault\": false } }, \"user_id\": \"user_2tiPJAs3bq0sSICjYZJ2ntCwaYq\"}"

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, apiError)
		}))
		defer server.Close()

		url := server.URL + "/api/v1/chat/completions"
		reqBody, _ := json.Marshal(types.OpenAIRequest{
			Model:          "deepseek/deepseek-chat:free",
			TextMessage:    types.Message[string]{Role: "system", Content: types.LLMContent},
			MatchesMessage: types.Message[[]types.MatchData]{Role: "user", Content: matches},
		})

		res, err := GenerateNicknames(url, reqBody)

		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "ошибка чтения ответа LLM API:")
	})
}
