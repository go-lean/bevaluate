package config

type (
	Config struct {
		Packages    Packages    `yaml:"packages"`
		Evaluations Evaluations `yaml:"evaluations"`
	}

	Packages struct {
		IgnoredDirs []string `yaml:"ignored_dirs,flow"`
	}

	Evaluations struct {
		DeploymentsDir string       `yaml:"deployments_dir"`
		RetestOut      string       `yaml:"retest_out"`
		RedeployOut    string       `yaml:"redeploy_out"`
		SpecialCases   SpecialCases `yaml:"special_cases"`
	}

	SpecialCases struct {
		RetestTriggers    []string `yaml:"retest_triggers,flow"`
		FullScaleTriggers []string `yaml:"full_scale_triggers,flow"`
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
			RetestOut:      "bevaluate/retest.out",
			RedeployOut:    "bevaluate/redeploy.out",
			SpecialCases: SpecialCases{
				FullScaleTriggers: []string{"go.mod$"},
			},
		},
	}
}
