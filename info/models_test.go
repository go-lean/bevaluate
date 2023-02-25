package info_test

import (
	"github.com/go-lean/bevaluate/info"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewConfig_InvalidExpression_Panic(t *testing.T) {
	require.Panics(t, func() {
		_ = info.NewConfig("[")
	})
}
