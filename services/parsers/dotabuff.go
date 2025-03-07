package parsers

import (
	"dota-nicknames/helpers"
	"dota-nicknames/types"
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Fetcher func(id int) ([]types.MatchData, error)

func FetchMatchData(id int) ([]types.MatchData, error) {
	url := fmt.Sprintf("https://www.dotabuff.com/players/%d/matches", id)

	req, err := http.NewRequest("GET", url, nil)

	client := &http.Client{
		Transport: &helpers.LoggingTransport{Transport: http.DefaultTransport},
	}
	res, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("Ошибка при обращении к %s: %w", url, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ошибка при обращении к %s: %w", url, err)
	}

	// Добавить обработку закрытого профиля
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	matches := goquery.Map(doc.Find("section tbody tr").Slice(0, 20), func(i int, s *goquery.Selection) types.MatchData {
		match := types.MatchData{}

		match.Hero = s.Find("td").Eq(1).Find("a").Text()
		match.Role = s.Find("td").Eq(2).Find("i").Eq(0).AttrOr("title", "")
		match.Lane = s.Find("td").Eq(2).Find("i").Eq(1).AttrOr("title", "")
		match.Result = s.Find("td").Eq(3).Find("a").Text()
		match.Duration = s.Find("td").Eq(5).Text()
		match.KDA = s.Find("td").Eq(6).Find(".kda-record").Text()
		match.Items = s.Find("td").Eq(7).Find("img").Map(func(_ int, s *goquery.Selection) string {
			return s.AttrOr("title", "")
		})

		return match
	})

	return matches, nil
}
