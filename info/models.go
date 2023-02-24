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
)
