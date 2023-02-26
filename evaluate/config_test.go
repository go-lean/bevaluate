package evaluate_test

import (
	"github.com/go-lean/bevaluate/evaluate"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewConfig_BadRetestExpression_Panic(t *testing.T) {
	require.Panics(t, func() {
		_ = evaluate.NewConfig("cmd/", []string{"["}, nil)
	})
}

func TestNewConfig_BadRedeployExpression_Panic(t *testing.T) {
	require.Panics(t, func() {
		_ = evaluate.NewConfig("cmd/", nil, []string{"["})
	})
}
