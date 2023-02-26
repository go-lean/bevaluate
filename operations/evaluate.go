package operations

import (
	"fmt"
	"github.com/go-lean/bevaluate/config"
	"github.com/go-lean/bevaluate/evaluate"
	"github.com/go-lean/bevaluate/info"
	"github.com/go-lean/bevaluate/storage"
	"path/filepath"
)

type (
	EvaluateBuildOperation struct {
		cfg   config.Config
		store storage.Store
	}
)

func NewEvaluateOperation(store storage.Store, cfg config.Config) EvaluateBuildOperation {
	return EvaluateBuildOperation{
		cfg:   cfg,
		store: store,
	}
}

func (o EvaluateBuildOperation) Run(root, changesContent string) error {
	changes, errParse := info.ParseGitChanges(changesContent)
	if errParse != nil {
		return fmt.Errorf("could not parse changes: %w", errParse)
	}

	if len(changes) == 0 {
		return nil
	}

	moduleName, errName := storage.ReadModuleName(filepath.Join(root, "go.mod"), o.store.FileOpener)
	if errName != nil {
		return fmt.Errorf("could not read go module name: %w", errName)
	}

	infoCfg := info.NewConfig(o.cfg.Packages.IgnoredDirs...)
	packageReader := info.NewPackageReader(o.store.DirReader, o.store.FileOpener, infoCfg)

	packages, errRead := packageReader.ReadRecursively(root, moduleName)
	if errRead != nil {
		return fmt.Errorf("could not read packages: %w", errRead)
	}

	if len(packages) == 0 {
		return nil
	}

	evalCfg := evaluate.NewConfig(
		o.cfg.Evaluations.DeploymentsDir,
		o.cfg.Evaluations.SpecialCases.Retest,
		o.cfg.Evaluations.SpecialCases.FullScale)

	evaluator := evaluate.NewBuildEvaluator(evalCfg)
	result, errEvaluate := evaluator.Evaluate(packages, changes)
	if errEvaluate != nil {
		return fmt.Errorf("could not evaluate build: %w", errEvaluate)
	}

	fmt.Println(result)
	return nil
}
