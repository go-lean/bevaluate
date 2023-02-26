package evaluate

import (
	"fmt"
	"github.com/go-lean/bevaluate/info"
)

type (
	DependencyGraph struct {
		Nodes    []*DependencyNode
		NodesMap map[string]*DependencyNode
	}

	DependencyNode struct {
		info.PackageInfo
		Dependants []*DependencyNode
		retest     bool
		redeploy   bool
	}
)

func NewDependencyGraph(packages []info.PackageInfo) DependencyGraph {
	l := len(packages)
	nodes := make([]*DependencyNode, l)
	nodesMap := make(map[string]*DependencyNode, l)

	for i := 0; i < l; i++ {
		node := &DependencyNode{PackageInfo: packages[i]}
		nodes[i] = node
		nodesMap[node.Path] = node
	}

	return DependencyGraph{
		Nodes:    nodes,
		NodesMap: nodesMap,
	}
}

func (g DependencyGraph) Build() error {
	for _, node := range g.Nodes {
		for _, dependency := range node.Dependencies {
			dependencyNode, ok := g.NodesMap[dependency]
			if ok == false {
				return fmt.Errorf("could not find dependency node with path: %q", dependency)
			}

			dependencyNode.Dependants = append(dependencyNode.Dependants, node)
		}
	}

	return nil
}
