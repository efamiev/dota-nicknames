package main

import (
	"dota-nicknames/services/llm"
)

func main() {
	print(llm.GenerateNicknames(321580662))
}
