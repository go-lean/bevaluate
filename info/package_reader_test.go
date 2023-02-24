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
	r := info.NewPackageReader(dirReader, opener)

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
	r := info.NewPackageReader(dirReader, opener)

	packages, errRead := r.ReadRecursively("baba", testModuleName)

	require.NoError(t, errRead)
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
	r := info.NewPackageReader(dirReader, opener)

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

	r := info.NewPackageReader(dirReader, opener)

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

	r := info.NewPackageReader(dirReader, opener)

	packages, errRead := r.ReadRecursively("baba", testModuleName)

	require.Error(t, errRead)
	require.Contains(t, errRead.Error(), "kaboom")
	require.NotEqual(t, "kaboom", errRead.Error())
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

	r := info.NewPackageReader(dirReader, opener)

	packages, errRead := r.ReadRecursively("baba", testModuleName)

	require.NoError(t, errRead)
	require.Len(t, packages, 1)

	require.Equal(t, "serviceone", packages[0].Path)
	require.False(t, packages[0].ContainsTests)
	require.Empty(t, packages[0].Dependencies)
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

	r := info.NewPackageReader(dirReader, opener)

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

	r := info.NewPackageReader(dirReader, opener)

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

	r := info.NewPackageReader(dirReader, opener)

	packages, errRead := r.ReadRecursively("baba", testModuleName)

	require.NoError(t, errRead)
	require.Len(t, packages, 1)

	require.Len(t, packages[0].Dependencies, 1)
	require.Equal(t, "common", packages[0].Dependencies[0])
	require.True(t, packages[0].ContainsTests)
}
