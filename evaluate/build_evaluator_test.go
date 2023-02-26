package evaluate_test

import (
	"github.com/go-lean/bevaluate/evaluate"
	"github.com/go-lean/bevaluate/info"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

func testCfg() evaluate.Config {
	return evaluate.NewConfig("cmd/", nil, nil)
}

func TestBuildEvaluator_Evaluate_Nil_Empty(t *testing.T) {
	eval := evaluate.NewBuildEvaluator(testCfg())
	result, errEval := eval.Evaluate(nil, nil)

	require.NoError(t, errEval)
	require.Empty(t, result.Retest)
	require.Empty(t, result.Redeploy)
}

func TestBuildEvaluator_Evaluate_BadPackages_Error(t *testing.T) {
	eval := evaluate.NewBuildEvaluator(testCfg())
	packages := []info.PackageInfo{{Path: "baba", Dependencies: []string{"missing"}}}

	_, errEval := eval.Evaluate(packages, nil)

	require.Error(t, errEval)
	require.Contains(t, errEval.Error(), "missing")
}

func TestBuildEvaluator_Evaluate_ChangedSourceFileMissingPackage_Error(t *testing.T) {
	packages := []info.PackageInfo{
		{
			Path:          "other",
			ContainsTests: true,
		},
	}
	changes := []info.ChangeInfo{
		{
			Path: "other/server.go",
		},
		{
			Path: "baba/server.go",
		},
	}

	eval := evaluate.NewBuildEvaluator(testCfg())
	result, errEval := eval.Evaluate(packages, changes)

	require.Error(t, errEval)
	require.Contains(t, errEval.Error(), "missing")

	require.Empty(t, result.Retest)
	require.Empty(t, result.Redeploy)
}

func TestBuildEvaluator_Evaluate_DeletedSourceFileMissingPackage_NoError(t *testing.T) {
	packages := []info.PackageInfo{
		{
			Path:         "cmd/other",
			Dependencies: []string{"other"},
		},
		{
			Path:          "other",
			ContainsTests: true,
		},
	}
	changes := []info.ChangeInfo{
		{
			Path: "other/server.go",
		},
		{
			Path:      "baba/server.go",
			IsDeleted: true,
		},
	}

	eval := evaluate.NewBuildEvaluator(testCfg())
	result, errEval := eval.Evaluate(packages, changes)

	require.NoError(t, errEval)

	require.Len(t, result.Retest, 1)
	require.Equal(t, "other", result.Retest[0])

	require.Len(t, result.Redeploy, 1)
	require.Equal(t, "cmd/other", result.Redeploy[0])
}

func TestBuildEvaluator_Evaluate_ChangedSourceFileInsideRoot_IsIgnored(t *testing.T) {
	packages := []info.PackageInfo{
		{
			Path:         "cmd/other",
			Dependencies: []string{"other"},
		},
		{
			Path:          "other",
			ContainsTests: true,
		},
	}
	changes := []info.ChangeInfo{
		{
			Path: "other/server.go",
		},
		{
			Path: "common.go",
		},
	}

	eval := evaluate.NewBuildEvaluator(testCfg())
	result, errEval := eval.Evaluate(packages, changes)

	require.NoError(t, errEval)

	require.Len(t, result.Retest, 1)
	require.Equal(t, "other", result.Retest[0])

	require.Len(t, result.Redeploy, 1)
	require.Equal(t, "cmd/other", result.Redeploy[0])
}

func TestBuildEvaluator_Evaluate_ChangedTestSourceFile_OnlyAffectsRetesting(t *testing.T) {
	packages := []info.PackageInfo{
		{
			Path:         "cmd/other",
			Dependencies: []string{"other"},
		},
		{
			Path:          "other",
			ContainsTests: true,
		},
	}
	changes := []info.ChangeInfo{
		{
			Path: "other/server_test.go",
		},
	}

	eval := evaluate.NewBuildEvaluator(testCfg())
	result, errEval := eval.Evaluate(packages, changes)

	require.NoError(t, errEval)

	require.Len(t, result.Retest, 1)
	require.Equal(t, "other", result.Retest[0])

	require.Empty(t, result.Redeploy)
}

func TestBuildEvaluator_Evaluate_ChangedFileWithoutPackage_FindsParentPackage(t *testing.T) {
	packages := []info.PackageInfo{
		{
			Path:         "cmd/baba",
			Dependencies: []string{"baba"},
		},
		{
			Path:         "cmd/other",
			Dependencies: []string{"other"},
		},
		{
			Path:          "baba",
			ContainsTests: true,
		},
		{
			Path:          "other",
			ContainsTests: true,
		},
	}
	changes := []info.ChangeInfo{
		{
			Path: "other/templates/emails.en.json",
		},
	}

	eval := evaluate.NewBuildEvaluator(testCfg())
	result, errEval := eval.Evaluate(packages, changes)

	require.NoError(t, errEval)

	require.Len(t, result.Retest, 1)
	require.Equal(t, "other", result.Retest[0])

	require.Len(t, result.Redeploy, 1)
	require.Equal(t, "cmd/other", result.Redeploy[0])
}

func TestBuildEvaluator_Evaluate_ChangedFileWithoutPackage_NoParentPackage_FullScale(t *testing.T) {
	packages := []info.PackageInfo{
		{
			Path:         "cmd/baba",
			Dependencies: []string{"baba"},
		},
		{
			Path:         "cmd/other",
			Dependencies: []string{"other"},
		},
		{
			Path:          "baba",
			ContainsTests: true,
		},
		{
			Path:          "other",
			ContainsTests: true,
		},
	}
	changes := []info.ChangeInfo{
		{
			Path: "templates/emails/emails.en.json",
		},
	}

	eval := evaluate.NewBuildEvaluator(testCfg())
	result, errEval := eval.Evaluate(packages, changes)

	require.NoError(t, errEval)

	require.Len(t, result.Retest, 2)
	expectedRetest := []string{"other", "baba"}
	require.ElementsMatch(t, expectedRetest, result.Retest)

	require.Len(t, result.Redeploy, 2)
	expectedRedeploy := []string{"cmd/other", "cmd/baba"}
	require.ElementsMatch(t, expectedRedeploy, result.Redeploy)
}

func TestBuildEvaluator_Evaluate_ChangedFileMatchingSpecialRetestCase_FullRetest(t *testing.T) {
	packages := []info.PackageInfo{
		{
			Path:         "cmd/baba",
			Dependencies: []string{"baba"},
		},
		{
			Path:         "cmd/other",
			Dependencies: []string{"other"},
		},
		{
			Path:          "baba",
			ContainsTests: true,
		},
		{
			Path:          "other",
			ContainsTests: true,
		},
	}
	changes := []info.ChangeInfo{
		{
			Path: "helm/any_file",
		},
	}

	cfg := evaluate.NewConfig("cmd/", []string{"helm/.*"}, nil)
	eval := evaluate.NewBuildEvaluator(cfg)

	result, errEval := eval.Evaluate(packages, changes)

	require.NoError(t, errEval)

	require.Len(t, result.Retest, 2)
	expectedRetest := []string{"other", "baba"}
	require.ElementsMatch(t, expectedRetest, result.Retest)

	require.Empty(t, result.Redeploy)
}

func TestBuildEvaluator_Evaluate_ChangedFileNotMatchingSpecialRetestCase_NothingHappens(t *testing.T) {
	packages := []info.PackageInfo{
		{
			Path:         "cmd/baba",
			Dependencies: []string{"baba"},
		},
		{
			Path:          "baba",
			ContainsTests: true,
		},
		{
			Path: "common",
		},
	}
	changes := []info.ChangeInfo{
		{
			Path: "common/errors.go",
		},
	}

	cfg := testCfg()
	exp := regexp.MustCompile("Jenkinsfile.*")
	cfg.SpecialCases.RetestAll = append(cfg.SpecialCases.RetestAll, exp)
	eval := evaluate.NewBuildEvaluator(cfg)

	result, errEval := eval.Evaluate(packages, changes)

	require.NoError(t, errEval)

	require.Empty(t, result.Retest)
	require.Empty(t, result.Redeploy)
}

func TestBuildEvaluator_Evaluate_ChangedFileMatchingSpecialRedeployCase_FullRedeploy(t *testing.T) {
	packages := []info.PackageInfo{
		{
			Path:         "cmd/baba",
			Dependencies: []string{"baba"},
		},
		{
			Path:         "cmd/other",
			Dependencies: []string{"other"},
		},
		{
			Path:          "baba",
			ContainsTests: true,
		},
		{
			Path:          "other",
			ContainsTests: true,
		},
	}
	changes := []info.ChangeInfo{
		{
			Path: "helm/any_file",
		},
	}

	cfg := evaluate.NewConfig("cmd/", nil, []string{"helm/.*"})
	eval := evaluate.NewBuildEvaluator(cfg)

	result, errEval := eval.Evaluate(packages, changes)

	require.NoError(t, errEval)

	require.Len(t, result.Retest, 2)
	expectedRetest := []string{"other", "baba"}
	require.ElementsMatch(t, expectedRetest, result.Retest)

	require.Len(t, result.Redeploy, 2)
	expectedRedeploy := []string{"cmd/other", "cmd/baba"}
	require.ElementsMatch(t, expectedRedeploy, result.Redeploy)
}

func TestBuildEvaluator_Evaluate_ChangedFileNotMatchingSpecialRedeployCase_NothingHappens(t *testing.T) {
	packages := []info.PackageInfo{
		{
			Path:         "cmd/baba",
			Dependencies: []string{"baba"},
		},
		{
			Path:          "baba",
			ContainsTests: true,
		},
		{
			Path: "common",
		},
	}
	changes := []info.ChangeInfo{
		{
			Path: "common/errors.go",
		},
	}

	cfg := testCfg()
	exp := regexp.MustCompile("Jenkinsfile.*")
	cfg.SpecialCases.RedeployAll = append(cfg.SpecialCases.RedeployAll, exp)
	eval := evaluate.NewBuildEvaluator(cfg)

	result, errEval := eval.Evaluate(packages, changes)

	require.NoError(t, errEval)

	require.Empty(t, result.Retest)
	require.Empty(t, result.Redeploy)
}

func TestBuildEvaluator_Evaluate_SharedDependencies_ShouldPropagateCorrectly(t *testing.T) {
	packages := []info.PackageInfo{
		{
			Path:         "cmd/baba",
			Dependencies: []string{"baba"},
		},
		{
			Path:         "cmd/other",
			Dependencies: []string{"other"},
		},
		{
			Path:         "cmd/otherrpc",
			Dependencies: []string{"otherrpc"},
		},
		{
			Path:          "baba",
			ContainsTests: true,
			Dependencies:  []string{"dbtype"},
		},
		{
			Path:          "otherrpc",
			ContainsTests: true,
		},
		{
			Path:          "other",
			ContainsTests: true,
			Dependencies:  []string{"common", "dbtype"},
		},
		{
			Path:          "common",
			ContainsTests: true,
			Dependencies:  []string{"dbtype"},
		},
		{
			Path: "dbtype",
		},
	}
	changes := []info.ChangeInfo{
		{
			Path: "dbtype/id.go",
		},
	}

	cfg := testCfg()
	exp := regexp.MustCompile("Jenkinsfile.*")
	cfg.SpecialCases.RedeployAll = append(cfg.SpecialCases.RedeployAll, exp)
	eval := evaluate.NewBuildEvaluator(cfg)

	result, errEval := eval.Evaluate(packages, changes)

	require.NoError(t, errEval)

	require.Len(t, result.Retest, 3)
	expectedRetest := []string{"baba", "other", "common"}
	require.ElementsMatch(t, expectedRetest, result.Retest)

	require.Len(t, result.Redeploy, 2)
	expectedRedeploy := []string{"cmd/baba", "cmd/other"}
	require.ElementsMatch(t, expectedRedeploy, result.Redeploy)
}
