package config

type (
	Config struct {
		Packages Packages `yaml:"packages"`
	}

	Packages struct {
		IgnoredDirs []string `yaml:"ignoredDirs,flow"`
	}
)

func Default() Config {
	return Config{
		Packages: Packages{
			IgnoredDirs: []string{
				"build$",
				"vendor$",
				".*/mocks$",
			},
		},
	}
}
