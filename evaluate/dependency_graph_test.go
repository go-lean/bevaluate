package evaluate_test

import (
	"github.com/go-lean/bevaluate/evaluate"
	"github.com/go-lean/bevaluate/info"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDependencyGraph_Build_Nil_NoError(t *testing.T) {
	g := evaluate.NewDependencyGraph(nil)
	err := g.Build()

	require.NoError(t, err)
	require.Empty(t, g.Nodes)
	require.Empty(t, g.NodesMap)
}

func TestDependencyGraph_Build_MissingPackage_Error(t *testing.T) {
	packages := []info.PackageInfo{
		{
			Path:         "baba",
			Dependencies: []string{"verymissing"},
		},
	}

	g := evaluate.NewDependencyGraph(packages)
	err := g.Build()

	require.Error(t, err)
	require.Contains(t, err.Error(), "verymissing")
}

func TestDependencyGraph_Build_NoDependencies_NoDependants(t *testing.T) {
	packages := []info.PackageInfo{
		{
			Path: "baba",
		},
		{
			Path: "other",
		},
	}

	g := evaluate.NewDependencyGraph(packages)
	err := g.Build()

	require.NoError(t, err)
	require.Len(t, g.Nodes, 2)
	require.Len(t, g.NodesMap, 2)

	require.Empty(t, g.NodesMap["baba"].Dependants)
	require.Empty(t, g.NodesMap["other"].Dependants)
}

func TestDependencyGraph_Build_OnePackageDependantOnOther_OneDependant(t *testing.T) {
	packages := []info.PackageInfo{
		{
			Path:         "baba",
			Dependencies: []string{"other"},
		},
		{
			Path: "other",
		},
	}

	g := evaluate.NewDependencyGraph(packages)
	err := g.Build()

	require.NoError(t, err)
	require.Len(t, g.Nodes, 2)
	require.Len(t, g.NodesMap, 2)

	require.Empty(t, g.NodesMap["baba"].Dependants)
	require.Len(t, g.NodesMap["other"].Dependants, 1)
	require.Equal(t, g.NodesMap["baba"], g.NodesMap["other"].Dependants[0])
}

func TestDependencyGraph_Build_ThreeDependantsInARow(t *testing.T) {
	packages := []info.PackageInfo{
		{
			Path:         "baba",
			Dependencies: []string{"other"},
		},
		{
			Path:         "other",
			Dependencies: []string{"common"},
		},
		{
			Path:         "common",
			Dependencies: []string{"dal"},
		},
		{
			Path: "dal",
		},
	}

	g := evaluate.NewDependencyGraph(packages)
	err := g.Build()

	require.NoError(t, err)
	require.Len(t, g.Nodes, 4)
	require.Len(t, g.NodesMap, 4)

	require.Empty(t, g.NodesMap["baba"].Dependants)

	require.Len(t, g.NodesMap["other"].Dependants, 1)
	require.Equal(t, g.NodesMap["baba"], g.NodesMap["other"].Dependants[0])

	require.Len(t, g.NodesMap["common"].Dependants, 1)
	require.Equal(t, g.NodesMap["other"], g.NodesMap["common"].Dependants[0])

	require.Len(t, g.NodesMap["dal"].Dependants, 1)
	require.Equal(t, g.NodesMap["common"], g.NodesMap["dal"].Dependants[0])
}

func TestDependencyGraph_Build_AllDependOnOne(t *testing.T) {
	packages := []info.PackageInfo{
		{
			Path:         "baba",
			Dependencies: []string{"dal"},
		},
		{
			Path:         "other",
			Dependencies: []string{"dal"},
		},
		{
			Path:         "common",
			Dependencies: []string{"dal"},
		},
		{
			Path: "dal",
		},
	}

	g := evaluate.NewDependencyGraph(packages)
	err := g.Build()

	require.NoError(t, err)
	require.Len(t, g.Nodes, 4)
	require.Len(t, g.NodesMap, 4)

	require.Empty(t, g.NodesMap["baba"].Dependants)
	require.Empty(t, g.NodesMap["other"].Dependants)
	require.Empty(t, g.NodesMap["common"].Dependants)

	require.Len(t, g.NodesMap["dal"].Dependants, 3)
	require.Equal(t, g.NodesMap["baba"], g.NodesMap["dal"].Dependants[0])
	require.Equal(t, g.NodesMap["other"], g.NodesMap["dal"].Dependants[1])
	require.Equal(t, g.NodesMap["common"], g.NodesMap["dal"].Dependants[2])
}
