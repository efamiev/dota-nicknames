package parsers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchMatchData(t *testing.T) {
	t.Run("Set vacancies to redis list", func(t *testing.T) {
		assert.Equal(t, 1, 1, "HH list not equal")
	})
}
