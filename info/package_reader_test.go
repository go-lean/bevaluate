package info_test

import (
	"github.com/go-lean/bevaluate/info"
	"github.com/go-lean/bevaluate/models"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPackageReader_ReadRecursively_DirReaderError(t *testing.T) {
	dirReader := NewDirReader().CanRead(false)

	opener := NewFileOpener()
	r := info.NewPackageReader(dirReader, opener, emptyConfig)

	packages, errRead := r.ReadRecursively("baba", testModuleName)

	require.Error(t, errRead)
	require.Contains(t, errRead.Error(), "kaboom")
	require.NotEqual(t, "kaboom", errRead.Error())
	require.Empty(t, packages)
}

func TestPackageReader_ReadRecursively_EmptyRootDir(t *testing.T) {
	dirReader := NewDirReader()
	dirReader.MockAt("baba", []models.DirEntry{})

	opener := NewFileOpener()
	r := info.NewPackageReader(dirReader, opener, emptyConfig)

	packages, errRead := r.ReadRecursively("baba", testModuleName)

	require.NoError(t, errRead)
	require.Empty(t, packages)
}

func TestPackageReader_ReadRecursively_ExplosiveSubDir_Error(t *testing.T) {
	dirReader := NewDirReader()
	dirReader.MockAt("baba", []models.DirEntry{
		DirEntry{
			name:  "kaboom",
			isDir: true,
		},
	})

	opener := NewFileOpener()
	r := info.NewPackageReader(dirReader, opener, emptyConfig)

	packages, errRead := r.ReadRecursively("baba", testModuleName)

	require.Error(t, errRead)
	require.Contains(t, errRead.Error(), "kaboom")
	require.NotEqual(t, "kaboom", errRead.Error())
	require.Empty(t, packages)
}

func TestPackageReader_ReadRecursively_ExplosiveInnerSubDir_Error(t *testing.T) {
	dirReader := NewDirReader()
	dirReader.MockAt("baba", []models.DirEntry{
		DirEntry{
			name:  "serviceone",
			isDir: true,
		},
	})
	dirReader.MockAt("baba/serviceone", []models.DirEntry{
		DirEntry{
			name:  "kaboom",
			isDir: true,
		},
	})

	opener := NewFileOpener()
	r := info.NewPackageReader(dirReader, opener, emptyConfig)

	packages, errRead := r.ReadRecursively("baba", testModuleName)

	require.Error(t, errRead)
	require.Contains(t, errRead.Error(), "kaboom")
	require.NotEqual(t, "kaboom", errRead.Error())
	require.Empty(t, packages)
}

func TestPackageReader_ReadRecursively_EmptySubDir_ShouldBeEmpty(t *testing.T) {
	dirReader := NewDirReader()
	dirReader.MockAt("baba", []models.DirEntry{
		DirEntry{
			name:  "serviceone",
			isDir: true,
		},
	})
	dirReader.MockAt("baba/serviceone", []models.DirEntry{})

	opener := NewFileOpener()
	r := info.NewPackageReader(dirReader, opener, emptyConfig)

	packages, errRead := r.ReadRecursively("baba", testModuleName)

	require.NoError(t, errRead)
	require.Empty(t, packages)
}

func TestPackageReader_ReadRecursively_NonSourceFile_ShouldBeEmpty(t *testing.T) {
	dirReader := NewDirReader()
	dirReader.MockAt("baba", []models.DirEntry{
		DirEntry{
			name:  "serviceone",
			isDir: true,
		},
	})
	dirReader.MockAt("baba/serviceone", []models.DirEntry{
		DirEntry{
			name: "baba.txt",
		},
	})

	opener := NewFileOpener()

	r := info.NewPackageReader(dirReader, opener, emptyConfig)

	packages, errRead := r.ReadRecursively("baba", testModuleName)

	require.NoError(t, errRead)
	require.Empty(t, packages)
}

func TestPackageReader_ReadRecursively_NonEmptySubDir_ErrOpen(t *testing.T) {
	dirReader := NewDirReader()
	dirReader.MockAt("baba", []models.DirEntry{
		DirEntry{
			name:  "serviceone",
			isDir: true,
		},
	})
	dirReader.MockAt("baba/serviceone", []models.DirEntry{
		DirEntry{
			name: "baba.go",
		},
	})

	opener := NewFileOpener().CanOpen(false)

	r := info.NewPackageReader(dirReader, opener, emptyConfig)

	packages, errRead := r.ReadRecursively("baba", testModuleName)

	require.Error(t, errRead)
	require.Contains(t, errRead.Error(), "kaboom")
	require.NotEqual(t, "kaboom", errRead.Error())
	require.Empty(t, packages)
}

func TestPackageReader_ReadRecursively_NonEmptyIgnoredSubDir_ShouldBeEmpty(t *testing.T) {
	dirReader := NewDirReader()
	dirReader.MockAt("baba", []models.DirEntry{
		DirEntry{
			name:  "serviceone",
			isDir: true,
		},
	})
	dirReader.MockAt("baba/serviceone", []models.DirEntry{
		DirEntry{
			name: "baba.go",
		},
	})

	fakeFile := NewFakeFile("package serviceone").CanClose(false)

	opener := NewFileOpener()
	opener.MockAt("baba/serviceone/baba.go", fakeFile)
	cfg := info.NewConfig(".*serviceone$")

	r := info.NewPackageReader(dirReader, opener, cfg)

	packages, errRead := r.ReadRecursively("baba", testModuleName)

	require.NoError(t, errRead)
	require.Empty(t, packages)
}

func TestPackageReader_ReadRecursively_NonEmptySubDir_ShouldNotBeEmpty(t *testing.T) {
	dirReader := NewDirReader()
	dirReader.MockAt("baba", []models.DirEntry{
		DirEntry{
			name:  "serviceone",
			isDir: true,
		},
	})
	dirReader.MockAt("baba/serviceone", []models.DirEntry{
		DirEntry{
			name: "baba.go",
		},
	})

	fakeFile := NewFakeFile("package serviceone").CanClose(false)

	opener := NewFileOpener()
	opener.MockAt("baba/serviceone/baba.go", fakeFile)

	r := info.NewPackageReader(dirReader, opener, emptyConfig)

	packages, errRead := r.ReadRecursively("baba", testModuleName)

	require.NoError(t, errRead)
	require.Len(t, packages, 1)

	require.Equal(t, "serviceone", packages[0].Path)
	require.False(t, packages[0].ContainsTests)
	require.Empty(t, packages[0].Dependencies)
}

func TestPackageReader_ReadRecursively_BadGoCodeBeforeImports_Error(t *testing.T) {
	dirReader := NewDirReader()
	dirReader.MockAt("baba", []models.DirEntry{
		DirEntry{
			name:  "serviceone",
			isDir: true,
		},
	})
	dirReader.MockAt("baba/serviceone", []models.DirEntry{
		DirEntry{
			name: "baba.go",
		},
	})

	opener := NewFileOpener()
	opener.MockAt("baba/serviceone/baba.go", NewFakeFile("bad go code"))

	r := info.NewPackageReader(dirReader, opener, emptyConfig)

	packages, errRead := r.ReadRecursively("baba", testModuleName)

	require.Error(t, errRead)
	require.Contains(t, errRead.Error(), "parse")
	require.Empty(t, packages)
}

func TestPackageReader_ReadRecursively_NonInternalDependencies_ShouldHaveNoDependencies(t *testing.T) {
	dirReader := NewDirReader()
	dirReader.MockAt("baba", []models.DirEntry{
		DirEntry{
			name:  "serviceone",
			isDir: true,
		},
	})
	dirReader.MockAt("baba/serviceone", []models.DirEntry{
		DirEntry{
			name: "baba.go",
		},
	})

	opener := NewFileOpener()
	opener.MockAt("baba/serviceone/baba.go", NewFakeFile(
		`
package serviceone

import (
	"io"
	"github.com/some/dependency"
)
`))

	r := info.NewPackageReader(dirReader, opener, emptyConfig)

	packages, errRead := r.ReadRecursively("baba", testModuleName)

	require.NoError(t, errRead)
	require.Len(t, packages, 1)

	require.Empty(t, packages[0].Dependencies)
}

func TestPackageReader_ReadRecursively_InternalDependency_ShouldHaveOneDependency(t *testing.T) {
	dirReader := NewDirReader()
	dirReader.MockAt("baba", []models.DirEntry{
		DirEntry{
			name:  "serviceone",
			isDir: true,
		},
	})
	dirReader.MockAt("baba/serviceone", []models.DirEntry{
		DirEntry{
			name: "baba.go",
		},
	})

	opener := NewFileOpener()
	opener.MockAt("baba/serviceone/baba.go", NewFakeFile(
		`
package serviceone

import (
	"io"
	"github.com/some/dependency"
	"github.com/baba/is/you/common" // internal dependency
)
`))

	r := info.NewPackageReader(dirReader, opener, emptyConfig)

	packages, errRead := r.ReadRecursively("baba", testModuleName)

	require.NoError(t, errRead)
	require.Len(t, packages, 1)

	require.Len(t, packages[0].Dependencies, 1)
	require.Equal(t, "common", packages[0].Dependencies[0])
}

func TestPackageReader_ReadRecursively_InternalDependencyInTestFile_ShouldHaveOneDependency(t *testing.T) {
	dirReader := NewDirReader()
	dirReader.MockAt("baba", []models.DirEntry{
		DirEntry{
			name:  "serviceone",
			isDir: true,
		},
	})
	dirReader.MockAt("baba/serviceone", []models.DirEntry{
		DirEntry{
			name: "baba.go",
		},
		DirEntry{
			name: "baba_test.go",
		},
	})

	opener := NewFileOpener()
	opener.MockAt("baba/serviceone/baba.go", NewFakeFile(
		`
package serviceone

import (
	"io"
	"github.com/some/dependency"
)
`))
	opener.MockAt("baba/serviceone/baba_test.go", NewFakeFile(
		`
package serviceone_test

import (
	"testing"
	"github.com/baba/is/you/serviceone" // internal dependency
	"github.com/baba/is/you/common" // internal dependency
)
`))

	r := info.NewPackageReader(dirReader, opener, emptyConfig)

	packages, errRead := r.ReadRecursively("baba", testModuleName)

	require.NoError(t, errRead)
	require.Len(t, packages, 1)

	require.Len(t, packages[0].Dependencies, 1)
	require.Equal(t, "common", packages[0].Dependencies[0])
	require.True(t, packages[0].ContainsTests)
}

func TestPackageReader_ReadRecursively_InnerIgnoredSubDir_ShouldNotBeEmpty(t *testing.T) {
	dirReader := NewDirReader()
	dirReader.MockAt("baba", []models.DirEntry{
		DirEntry{
			name:  "serviceone",
			isDir: true,
		},
	})
	dirReader.MockAt("baba/serviceone", []models.DirEntry{
		DirEntry{
			name:  "inner",
			isDir: true,
		},
	})
	dirReader.MockAt("baba/serviceone/inner", []models.DirEntry{
		DirEntry{
			name: "baba.go",
		},
		DirEntry{
			name: "baba_test.go",
		},
	})

	opener := NewFileOpener()
	opener.MockAt("baba/serviceone/inner/baba.go", NewFakeFile("package inner"))
	opener.MockAt("baba/serviceone/inner/baba_test.go", NewFakeFile(`package inner_test)
`))

	cfg := info.NewConfig("serviceone/inner$")
	r := info.NewPackageReader(dirReader, opener, cfg)

	packages, errRead := r.ReadRecursively("baba", testModuleName)

	require.NoError(t, errRead)
	require.Empty(t, packages)
}

func TestPackageReader_ReadRecursively_InnerSubDir_ShouldNotBeEmpty(t *testing.T) {
	dirReader := NewDirReader()
	dirReader.MockAt("baba", []models.DirEntry{
		DirEntry{
			name:  "serviceone",
			isDir: true,
		},
	})
	dirReader.MockAt("baba/serviceone", []models.DirEntry{
		DirEntry{
			name:  "inner",
			isDir: true,
		},
	})
	dirReader.MockAt("baba/serviceone/inner", []models.DirEntry{
		DirEntry{
			name: "baba.go",
		},
		DirEntry{
			name: "baba_test.go",
		},
	})

	opener := NewFileOpener()
	opener.MockAt("baba/serviceone/inner/baba.go", NewFakeFile(`
package inner

import (
	"io"
	"github.com/some/third-party/dependency"
	"github.com/baba/is/you/common"
	"github.com/baba/is/you/serviceone"
)
`))
	opener.MockAt("baba/serviceone/inner/baba_test.go", NewFakeFile(`
package inner

import (
	"io"
	"github.com/some/third-party/dependency"
	"github.com/baba/is/you/serviceone/inner"
	"github.com/baba/is/you/serviceone"
)
`))

	r := info.NewPackageReader(dirReader, opener, emptyConfig)

	packages, errRead := r.ReadRecursively("baba", testModuleName)

	require.NoError(t, errRead)
	require.Len(t, packages, 1)
	require.Equal(t, "serviceone/inner", packages[0].Path)

	require.Len(t, packages[0].Dependencies, 2)
	expected := []string{"serviceone", "common"}
	require.ElementsMatch(t, expected, packages[0].Dependencies)
}

func TestPackageReader_ReadRecursively_AllRootFilesAreIgnored(t *testing.T) {
	dirReader := NewDirReader()
	dirReader.MockAt("baba", []models.DirEntry{
		DirEntry{
			name: "go.mod",
		},
		DirEntry{
			name: "baba.go",
		},
		DirEntry{
			name: "baba_test.go",
		},
	})

	opener := NewFileOpener()

	r := info.NewPackageReader(dirReader, opener, emptyConfig)

	packages, errRead := r.ReadRecursively("baba", testModuleName)

	require.NoError(t, errRead)
	require.Empty(t, packages)
}

func TestPackageReader_ReadRecursively_NamedImportsAreAlsoRelative(t *testing.T) {
	dirReader := NewDirReader()
	dirReader.MockAt("baba", []models.DirEntry{
		DirEntry{
			name:  "service",
			isDir: true,
		},
	})
	dirReader.MockAt("baba/service", []models.DirEntry{
		DirEntry{
			name: "baba.go",
		},
	})

	opener := NewFileOpener()
	opener.MockAt("baba/service/baba.go", NewFakeFile(`
package service

import apps "github.com/baba/is/you/plugins"
import named "github.com/baba/is/you/service/inner"
`))

	r := info.NewPackageReader(dirReader, opener, emptyConfig)

	packages, errRead := r.ReadRecursively("baba", testModuleName)

	require.NoError(t, errRead)
	require.Len(t, packages, 1)

	require.Len(t, packages[0].Dependencies, 2)
	expected := []string{"plugins", "service/inner"}
	require.ElementsMatch(t, expected, packages[0].Dependencies)
}

func TestPackageReader_ReadRecursively_IgnoringInnerSubDirs(t *testing.T) {
	dirReader := NewDirReader()
	dirReader.MockAt("baba", []models.DirEntry{
		DirEntry{
			name:  "service",
			isDir: true,
		},
	})
	dirReader.MockAt("baba/service", []models.DirEntry{
		DirEntry{
			name:  "mocks",
			isDir: true,
		},
		DirEntry{
			name: "server.go",
		},
	})
	dirReader.MockAt("baba/service/mocks", []models.DirEntry{
		DirEntry{
			name: "baba.go",
		},
	})

	opener := NewFileOpener()
	opener.MockAt("baba/service/server.go", NewFakeFile(`
package service

import "github.com/baba/is/you/service/mocks"
`))
	opener.MockAt("baba/service/mocks/baba.go", NewFakeFile(`
package mocks

import "github.com/baba/is/you/other"
`))

	r := info.NewPackageReader(dirReader, opener, info.NewConfig(".*/mocks$"))

	packages, errRead := r.ReadRecursively("baba", testModuleName)

	require.NoError(t, errRead)
	require.Len(t, packages, 1)
	require.Equal(t, "service", packages[0].Path)
	require.Empty(t, packages[0].Dependencies)
}
