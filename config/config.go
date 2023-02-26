package config

type (
	Config struct {
		Packages    Packages    `yaml:"packages"`
		Evaluations Evaluations `yaml:"evaluations"`
	}

	Packages struct {
		IgnoredDirs []string `yaml:"ignoredDirs,flow"`
	}

	Evaluations struct {
		DeploymentsDir string       `yaml:"deploymentsDir"`
		SpecialCases   SpecialCases `yaml:"specialCases"`
	}

	SpecialCases struct {
		Retest    []string `yaml:"retest,flow"`
		FullScale []string `yaml:"fullScale,flow"`
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
		Evaluations: Evaluations{
			DeploymentsDir: "cmd/",
			SpecialCases: SpecialCases{
				FullScale: []string{"go.mod"},
			},
		},
	}
}
