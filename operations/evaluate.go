package operations

import (
	"fmt"
	"github.com/go-lean/bevaluate/info"
	"github.com/go-lean/bevaluate/storage"
	"path/filepath"
)

func EvaluateBuild(root, changesContent string) error {
	changes, errParse := info.ParseGitChanges(changesContent)
	if errParse != nil {
		return fmt.Errorf("could not parse changes: %w", errParse)
	}

	if len(changes) == 0 {
		return nil
	}

	opener := storage.FileOpener{}
	reader := storage.DirReader{}

	moduleName, errName := storage.ReadModuleName(filepath.Join(root, "go.mod"), opener)
	if errName != nil {
		return fmt.Errorf("could not read go module name: %w", errName)
	}

	packageReader := info.NewPackageReader(reader, opener, info.Config{})
	packages, errRead := packageReader.ReadRecursively(root, moduleName)
	if errRead != nil {
		return fmt.Errorf("could not read packages: %w", errRead)
	}

	if len(packages) == 0 {
		return nil
	}

	return nil
}
