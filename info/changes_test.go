package info_test

import (
	"fmt"
	"github.com/go-lean/bevaluate/info"
	"github.com/stretchr/testify/require"
	"os/exec"
	"testing"
)

func TestCollectChanges_BadTarget_Error(t *testing.T) {
	changes, errCollect := info.CollectChanges("--baba")

	require.Error(t, errCollect)
	require.Empty(t, changes)
}

func TestCollectChanges_OK(t *testing.T) {
	cmd := exec.Command("git", "status")
	data, errCmd := cmd.Output()

	if errCmd != nil {
		fmt.Println("git not available")
		return
	}

	require.NotEmpty(t, data)
	_, errCollect := info.CollectChanges("master")

	require.NoError(t, errCollect)
}
