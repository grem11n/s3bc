package convert

import (
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func TestIsExcluded(t *testing.T) {
	t.Parallel()

	exList := []string{
		"foo/.*",
		"bizz/bazz",
	}

	testCases := map[string]bool{
		"foo/zoo":   true,
		"foo/meow":  true,
		"fizz/buzz": false,
		"bizz/zum":  false,
		"bizz/bazz": true,
	}

	for tc, expected := range testCases {
		got := isExcluded(exList, tc)
		assert.Equal(t, expected, got)
	}
}
