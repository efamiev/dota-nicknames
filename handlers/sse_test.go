package handlers

import (
	"dota-nicknames/types"
	"strconv"

	"testing"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
)

var matches = []types.MatchData{
	{
		Hero: "Queen of Pain", Result: "Lost Match",
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

func fetcher(id int) ([]types.MatchData, error) {
	return matches, nil
}

func TestGetMatchData(t *testing.T) {
	cacheExpiration := 5
	c := cache.New(time.Duration(cacheExpiration)*time.Minute, time.Duration(cacheExpiration)*2*time.Minute)

	t.Run("Получаем матчи из кеша", func(t *testing.T) {
		cacheId := 176586336
		c.Set(strconv.Itoa(cacheId), matches, cache.DefaultExpiration)

		res, _ := getMatchData(c, cacheId, fetcher)

		assert.Equal(t, matches, res)
	})

	t.Run("Делаем запрос, если данных в кеше нет", func(t *testing.T) {
		cachedMatches := []types.MatchData{
			{Hero: "Queen of Pain", Result: "Lost Match", KDA: "8/12/20", Duration: "46:57", Role: "Support Role", Lane: "Off Lane", Items: nil},
			{Hero: "Nature's Prophet", Result: "Won Match", KDA: "2/10/19", Duration: "37:47", Role: "Support Role", Lane: "Off Lane", Items: nil},
		}
		cacheId := 176586336
		c.Set(strconv.Itoa(cacheId), cachedMatches, cache.DefaultExpiration)

		res, _ := getMatchData(c, 176586337, fetcher)

		assert.Equal(t, matches, res)
	})
}
