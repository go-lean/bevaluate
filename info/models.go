package info

type (
	ChangeInfo struct {
		Path      string
		IsDeleted bool
	}

	PackageInfo struct {
		Path          string
		Dependencies  []string
		ContainsTests bool
	}

	Config struct {
		IgnoredDirs Ignored
	}

	Ignored map[string]struct{}
)

func (i Ignored) Contains(path string) bool {
	_, ok := i[path]
	return ok
}
