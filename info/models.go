package info

import "regexp"

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
		Ignored
	}

	Ignored struct {
		expressions []*regexp.Regexp
	}
)

func NewConfig(ignoredExpressions ...string) Config {
	expressions := make([]*regexp.Regexp, 0, len(ignoredExpressions))
	for _, expression := range ignoredExpressions {
		exp, err := regexp.Compile(expression)
		if err != nil {
			panic("could not parse ignored dir regex: " + err.Error())
		}

		expressions = append(expressions, exp)
	}

	return Config{Ignored{
		expressions: expressions,
	}}
}

func (i Ignored) IsIgnored(path string) bool {
	for _, exp := range i.expressions {
		if exp.MatchString(path) == false {
			continue
		}

		return true
	}

	return false
}
