package info_test

import (
	"github.com/go-lean/bevaluate/info"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseGitChanges_BadContent_Error(t *testing.T) {
	testCases := []struct {
		name    string
		content string
	}{
		{
			name:    "bad content should be invalid",
			content: "baba is you",
		},
		{
			name:    "empty path content should be invalid",
			content: "D\t ",
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			changes, errParse := info.ParseGitChanges(c.content)

			require.Error(t, errParse)
			require.Contains(t, errParse.Error(), "invalid")
			require.Empty(t, changes)
		})
	}
}

func TestParseGitChanges_OK(t *testing.T) {
	changes, errParse := info.ParseGitChanges("D\tbaba.go\nM\tflag.go")
	require.NoError(t, errParse)
	require.Len(t, changes, 2)
	expected := []info.ChangeInfo{{Path: "baba.go", IsDeleted: true}, {Path: "flag.go"}}
	require.ElementsMatch(t, expected, changes)
}

func TestParseGitChanges_EndingWithNewLine_OK(t *testing.T) {
	changes, errParse := info.ParseGitChanges("D\tbaba.go\nM\tflag.go\n")
	require.NoError(t, errParse)
	require.Len(t, changes, 2)
	expected := []info.ChangeInfo{{Path: "baba.go", IsDeleted: true}, {Path: "flag.go"}}
	require.ElementsMatch(t, expected, changes)
}

func TestParseGitChanges_EmptyContent_OK(t *testing.T) {
	changes, errParse := info.ParseGitChanges("")
	require.NoError(t, errParse)
	require.Empty(t, changes)
}
