package util_test

import (
	"github.com/go-lean/bevaluate/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMapKeys(t *testing.T) {
	testCases := []struct {
		name     string
		m        map[string]struct{}
		expected []string
	}{
		{
			name:     "nil should be empty",
			m:        nil,
			expected: []string{},
		},
		{
			name:     "not nil should not be empty",
			m:        map[string]struct{}{"baba": {}, "is": {}, "you": {}},
			expected: []string{"baba", "is", "you"},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			keys := util.MapKeys(c.m)

			require.ElementsMatch(t, c.expected, keys)
		})
	}
}
