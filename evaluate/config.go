package evaluate

import "regexp"

type (
	Config struct {
		DeploymentsDir string
		SpecialCases   SpecialCases
	}

	SpecialCases struct {
		RetestTriggers    []*regexp.Regexp
		FullScaleTriggers []*regexp.Regexp
	}
)

func NewConfig(deploymentsDir string, specialRetestCases, specialRedeployCases []string) Config {
	retest := make([]*regexp.Regexp, len(specialRetestCases))
	redeploy := make([]*regexp.Regexp, len(specialRedeployCases))

	for i, retestCase := range specialRetestCases {
		exp, errCompile := regexp.Compile(retestCase)
		if errCompile != nil {
			panic("could not compile retest special case: " + errCompile.Error())
		}

		retest[i] = exp
	}

	for i, redeployCase := range specialRedeployCases {
		exp, errCompile := regexp.Compile(redeployCase)
		if errCompile != nil {
			panic("could not compile full scale special case: " + errCompile.Error())
		}

		redeploy[i] = exp
	}

	return Config{
		DeploymentsDir: deploymentsDir,
		SpecialCases: SpecialCases{
			RetestTriggers:    retest,
			FullScaleTriggers: redeploy,
		},
	}
}
