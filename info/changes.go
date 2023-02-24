package info

import (
	"errors"
	"strings"
	"unicode"
)

var (
	errInvalidChangesFormat = errors.New("invalid changes format")
)

func ParseGitChanges(changesContent string) ([]ChangeInfo, error) {
	lines := strings.Split(changesContent, "\n")
	result := make([]ChangeInfo, 0, len(lines))
	for _, line := range lines {
		if line == "" {
			break
		}

		if isValidChange(line) == false {
			return nil, errInvalidChangesFormat
		}

		path := strings.Trim(line[2:], " ")
		if path == "" {
			return nil, errInvalidChangesFormat
		}
		result = append(result, ChangeInfo{
			Path:      path,
			IsDeleted: line[0] == 'D',
		})
	}

	return result, nil
}

func isValidChange(line string) bool {
	return len(line) >= 3 && unicode.IsDigit(rune(line[0])) == false && line[1] == '\t'
}
