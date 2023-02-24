package info

import (
	"fmt"
	"os/exec"
	"strings"
)

func CollectChanges(target string) ([]ChangeInfo, error) {
	cmd := exec.Command("git", "diff", target, "--name-status")
	data, errCmd := cmd.Output()

	if errCmd != nil {
		return nil, fmt.Errorf("could not execute git diff command")
	}

	lines := strings.Split(string(data), "\n")
	result := make([]ChangeInfo, 0, len(lines))
	for _, line := range lines {
		if line == "" {
			break
		}

		result = append(result, ChangeInfo{
			Path:      line[2:],
			IsDeleted: line[0] == 'D',
		})
	}

	return result, nil
}
