package info

import (
	"fmt"
	"github.com/go-lean/bevaluate/models"
	"github.com/go-lean/bevaluate/util"
	"github.com/zyedidia/generic/stack"
	"go/parser"
	"go/token"
	"io"
	"path/filepath"
	"strings"
)

type (
	PackageReader struct {
		fileOpener FileOpener
		dirReader  DirReader
		config     Config
	}

	FileOpener interface {
		OpenRead(path string) (io.ReadCloser, error)
	}

	DirReader interface {
		Read(path string) ([]models.DirEntry, error)
	}
)

func NewPackageReader(dirReader DirReader, fileOpener FileOpener, cfg Config) PackageReader {
	return PackageReader{
		dirReader:  dirReader,
		fileOpener: fileOpener,
		config:     cfg,
	}
}

func (r PackageReader) ReadRecursively(root, moduleName string) ([]PackageInfo, error) {
	entries, errRead := r.dirReader.Read(root)
	if errRead != nil {
		return nil, fmt.Errorf("could not read root directory: %w", errRead)
	}

	dirs := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() == false || r.config.IgnoredDirs.Contains(entry.Name()) {
			continue
		}

		dirs = append(dirs, entry.Name())
	}

	result, errRead := r.readSubDirsRecursively(root, moduleName, dirs)
	if errRead != nil {
		return nil, fmt.Errorf("could not read root sub dirs: %w", errRead)
	}

	return result, nil
}

func (r PackageReader) readSubDirsRecursively(root, moduleName string, dirs []string) ([]PackageInfo, error) {
	errChan := make(chan error, 1)
	pkgChan := make(chan []PackageInfo)

	for _, dir := range dirs {
		go func(dirPath string) {
			packages, errRead := r.readSubDirRecursively(root, moduleName, dirPath)
			if errRead != nil {
				errChan <- errRead
				return
			}

			pkgChan <- packages
		}(dir)
	}

	result := make([]PackageInfo, 0, len(dirs))
	for i := 0; i < len(dirs); i++ {
		select {
		case errRead := <-errChan:
			return nil, errRead
		case packages := <-pkgChan:
			result = append(result, packages...)
		}
	}

	return result, nil
}

func (r PackageReader) readSubDirRecursively(root, moduleName, subDir string) ([]PackageInfo, error) {
	dirsStack := stack.New[string]()
	dirsStack.Push(subDir)

	result := make([]PackageInfo, 0)

	for dirsStack.Size() > 0 {
		dir := dirsStack.Pop()

		entries, errRead := r.dirReader.Read(filepath.Join(root, dir))
		if errRead != nil {
			return nil, fmt.Errorf("could not read dir: %w", errRead)
		}

		sourceFiles := r.processEntries(dir, entries, dirsStack)
		if len(sourceFiles) == 0 {
			continue
		}

		pkg, errRead := r.readPackage(root, dir, moduleName, sourceFiles)
		if errRead != nil {
			return nil, fmt.Errorf("could not read package: %w", errRead)
		}

		result = append(result, pkg)
	}

	return result, nil
}

func (r PackageReader) readPackage(root, dir, moduleName string, sourceFiles []string) (PackageInfo, error) {
	dependencies := make(map[string]struct{}, 0)
	containsTests := false

	for _, filePath := range sourceFiles {
		file, errRead := r.fileOpener.OpenRead(filepath.Join(root, filePath))
		if errRead != nil {
			return PackageInfo{}, fmt.Errorf("could not read source file: %w", errRead)
		}

		parsedFile, errParse := parser.ParseFile(&token.FileSet{}, filePath, file, parser.ImportsOnly)
		if errParse != nil {
			return PackageInfo{}, fmt.Errorf("could not parse source file: %w", errParse)
		}

		for _, imp := range parsedFile.Imports {
			impPath := strings.Trim(imp.Path.Value, "\"")

			if impPath == "testing" && strings.HasSuffix(filePath, "_test.go") {
				containsTests = true
				continue
			}

			if strings.HasPrefix(impPath, moduleName) == false {
				continue // non internal dependency
			}

			dependency, _ := filepath.Rel(moduleName, impPath)
			if dependency == dir {
				continue
			}

			dependencies[dependency] = struct{}{}
		}
	}

	return PackageInfo{
		Path:          dir,
		Dependencies:  util.MapKeys(dependencies),
		ContainsTests: containsTests,
	}, nil
}

func (r PackageReader) processEntries(dirPath string, entries []models.DirEntry, dirsStack *stack.Stack[string]) []string {
	sourceFiles := make([]string, 0, len(entries))

	for _, entry := range entries {
		entryPath := filepath.Join(dirPath, entry.Name())
		if entry.IsDir() {
			if r.config.IgnoredDirs.Contains(entryPath) == false {
				dirsStack.Push(entryPath)
			}
			continue
		}

		if strings.HasSuffix(entryPath, ".go") == false {
			continue
		}

		sourceFiles = append(sourceFiles, entryPath)
	}

	return sourceFiles
}
