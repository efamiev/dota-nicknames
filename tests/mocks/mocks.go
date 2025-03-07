package mocks

import (
	"io"
	"log"
	"os"
)

func ReadHTML(name string) string {
	file, err := os.Open("./" + name)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("failed reading data from file: %s", err)
	}

	content := string(data)
	return content
}
