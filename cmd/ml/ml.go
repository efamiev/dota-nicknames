package main

import (
	"dota-nicknames/internal"
	"encoding/json"
	"log"
)

func nicks(url string) []internal.MatchData {
	matches := internal.FetchMatchData(url)
	matchesJson, _ := json.MarshalIndent(&matches, "", " ")

	log.Println("NICKS", string(matchesJson))
	return matches
}

func main() {
	url := "https://www.dotabuff.com/players/321580662/matches"
	
	nicks(url)
}

