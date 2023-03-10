package evaluate

import (
	"errors"
	"fmt"
	"github.com/go-lean/bevaluate/info"
	"github.com/zyedidia/generic/stack"
	"path/filepath"
	"strings"
)

type (
	BuildEvaluator struct {
		config Config
	}

	Evaluation struct {
		Retest   []string
		Redeploy []string
	}
)

const (
	defaultDependencyLevels = 5
)

var (
	ErrUnsupportedScenario  = errors.New("could not evaluate scenario")
	ErrSpecialRetestCase    = errors.New("could not continue change evaluation because of a special retest case")
	ErrSpecialFullScaleCase = errors.New("could not continue change evaluation because of a special full scale case")
)

func NewBuildEvaluator(config Config) BuildEvaluator {
	return BuildEvaluator{config: config}
}

func (e BuildEvaluator) Evaluate(packages []info.PackageInfo, changes []info.ChangeInfo) (Evaluation, error) {
	graph := NewDependencyGraph(packages)
	if errBuild := graph.Build(); errBuild != nil {
		return Evaluation{}, fmt.Errorf("could not build dependency graph: %w", errBuild)
	}

	issuedFullRetest := false
	for _, change := range changes {
		errEvaluate := e.evaluateChange(change, graph)
		if errEvaluate == nil {
			continue
		}

		if errors.Is(errEvaluate, ErrSpecialRetestCase) {
			if issuedFullRetest == false {
				issueFullScaleRetest(graph)
				issuedFullRetest = true
			}

			continue
		}

		if errors.Is(errEvaluate, ErrUnsupportedScenario) ||
			errors.Is(errEvaluate, ErrSpecialFullScaleCase) {
			issueFullScaleRetest(graph)
			e.issueFullScaleRedeploy(graph)

			break
		}

		return Evaluation{}, fmt.Errorf("could not evaluate change: %w", errEvaluate)
	}

	result := prepareEvaluation(graph)
	return result, nil
}

func (e BuildEvaluator) evaluateChange(change info.ChangeInfo, graph DependencyGraph) error {
	if errSpecialCase := e.evaluateSpecialCase(change); errSpecialCase != nil {
		return errSpecialCase
	}

	pkgPath := filepath.Dir(change.Path)
	if pkgPath == "." {
		return nil // unhandled special case, should be added in config
	}

	pkg, ok := graph.NodesMap[pkgPath]
	if ok == false {
		return e.handleMissingPackage(pkgPath, change, graph)
	}

	if strings.HasSuffix(change.Path, "_test.go") {
		if pkg.ContainsTests {
			pkg.retest = true
		}
		return nil
	}

	e.markPackageDirtyRecursively(pkg)
	return nil
}

func (e BuildEvaluator) handleMissingPackage(pkgPath string, change info.ChangeInfo, graph DependencyGraph) error {
	if strings.HasSuffix(change.Path, ".go") {
		if change.IsDeleted {
			return nil
		}

		return fmt.Errorf("missing package at: %q", pkgPath)
	}

	parent, ok := findParentRecursively(pkgPath, graph)
	if ok == false {
		return ErrUnsupportedScenario
	}

	e.markPackageDirtyRecursively(parent)
	return nil
}

func (e BuildEvaluator) markPackageDirtyRecursively(pkg *DependencyNode) {
	pkgStack := stack.New[*DependencyNode]()
	pkgStack.Push(pkg)
	visited := make(map[string]struct{}, defaultDependencyLevels)

	for pkgStack.Size() > 0 {
		p := pkgStack.Pop()

		if _, ok := visited[p.Path]; ok {
			continue
		}

		visited[p.Path] = struct{}{}

		if p.ContainsTests {
			p.retest = true
		}

		if e.canBeDeployed(p) {
			p.redeploy = true
		}

		for _, dependant := range p.Dependants {
			pkgStack.Push(dependant)
		}
	}
}

func findParentRecursively(pkgPath string, graph DependencyGraph) (*DependencyNode, bool) {
	path := filepath.Dir(pkgPath)

	for path != "." {
		pkg, ok := graph.NodesMap[path]
		if ok == false {
			path = filepath.Dir(path)
			continue
		}

		return pkg, true
	}

	return nil, false
}

func prepareEvaluation(graph DependencyGraph) Evaluation {
	l := len(graph.Nodes)
	retest := make([]string, 0, l)
	redeploy := make([]string, 0, l)

	for _, node := range graph.Nodes {
		if node.retest {
			retest = append(retest, node.Path)
		}
		if node.redeploy {
			redeploy = append(redeploy, node.Path)
		}
	}

	return Evaluation{
		Retest:   retest,
		Redeploy: redeploy,
	}
}

func issueFullScaleRetest(g DependencyGraph) {
	for _, node := range g.Nodes {
		if node.ContainsTests == false {
			continue
		}

		node.retest = true
	}
}

func (e BuildEvaluator) issueFullScaleRedeploy(g DependencyGraph) {
	for _, node := range g.Nodes {
		if e.canBeDeployed(node) == false {
			continue
		}

		node.redeploy = true
	}
}

func (e BuildEvaluator) canBeDeployed(node *DependencyNode) bool {
	return strings.HasPrefix(node.Path, e.config.DeploymentsDir)
}

func (e BuildEvaluator) evaluateSpecialCase(change info.ChangeInfo) error {
	for _, trigger := range e.config.SpecialCases.RetestTriggers {
		if trigger.MatchString(change.Path) == false {
			continue
		}

		return ErrSpecialRetestCase
	}

	for _, trigger := range e.config.SpecialCases.FullScaleTriggers {
		if trigger.MatchString(change.Path) == false {
			continue
		}

		return ErrSpecialFullScaleCase
	}

	return nil
}
