package operations

import (
	"fmt"
	"github.com/go-lean/bevaluate/config"
	"github.com/go-lean/bevaluate/evaluate"
	"github.com/go-lean/bevaluate/info"
	"github.com/go-lean/bevaluate/storage"
	"path/filepath"
	"strings"
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
		o.cfg.Evaluations.SpecialCases.RetestTriggers,
		o.cfg.Evaluations.SpecialCases.FullScaleTriggers)

	evaluator := evaluate.NewBuildEvaluator(evalCfg)
	result, errEvaluate := evaluator.Evaluate(packages, changes)
	if errEvaluate != nil {
		return fmt.Errorf("could not evaluate build: %w", errEvaluate)
	}

	if errWrite := o.writeResult(result); errWrite != nil {
		return fmt.Errorf("could not write result: %w", errWrite)
	}

	return nil
}

func (o EvaluateBuildOperation) writeResult(result evaluate.Evaluation) error {
	retestContent := strings.Join(result.Retest, storage.NewLine)
	if errWrite := storage.CreateFileWithText(o.cfg.Evaluations.RetestOut, retestContent, o.store.FileOpener); errWrite != nil {
		return fmt.Errorf("could not write retest result: %w", errWrite)
	}

	redeployContent := strings.Join(result.Redeploy, storage.NewLine)
	if errWrite := storage.CreateFileWithText(o.cfg.Evaluations.RedeployOut, redeployContent, o.store.FileOpener); errWrite != nil {
		return fmt.Errorf("could not write redeploy result: %w", errWrite)
	}

	return nil
}
