package dependency

import "fmt"

type Node interface {
	DependsOn(d Node)
	UnsetDependency(d Node)
	Name() string
	Dependencies() []Node
	Depended() []Node
	SetDependencies([]Node)
	SetDepended([]Node)
}

type BaseNode struct {
	name       string
	dependsOn  []Node
	dependedBy []Node
}

type DependencyRelation struct {
	Name string
	Deps []string
}

func (n *BaseNode) Name() string {
	return n.name
}

func (n *BaseNode) DependsOn(d Node) {
	d.SetDepended(append(d.Depended(), n))
	n.SetDependencies(append(n.Dependencies(), d))
}

func (n *BaseNode) UnsetDependency(d Node) {
	d.SetDepended(removeNodeFromList(d.Depended(), n))
	n.SetDependencies(removeNodeFromList(n.Dependencies(), d))
}

func (n *BaseNode) Dependencies() []Node {
	return n.dependsOn
}

func (n *BaseNode) Depended() []Node {
	return n.dependedBy
}

func (n *BaseNode) SetDependencies(nodeList []Node) {
	n.dependsOn = nodeList
}

func (n *BaseNode) SetDepended(nodeList []Node) {
	n.dependedBy = nodeList
}

func ManageDependency(items []DependencyRelation) (map[string]Node, error) {
	nodeMap := map[string]Node{}
	for _, item := range items {
		nodeMap[item.Name] = &BaseNode{name: item.Name}
	}

	for _, item := range items {
		node := nodeMap[item.Name]
		for _, dependency := range item.Deps {
			depNode, ok := nodeMap[dependency]
			if !ok {
				continue
			}
			node.DependsOn(depNode)
		}
	}

	return nodeMap, nil
}

func IsRoot(node Node) bool {
	return len(node.Dependencies()) == 0
}

func IsBottom(node Node) bool {
	return len(node.Depended()) == 0
}

func removeNodeFromList(nodeList []Node, removedNode Node) []Node {
	newList := []Node{}
	for _, n := range nodeList {
		if n.Name() != removedNode.Name() {
			newList = append(newList, n)
		}
	}
	return newList
}

func CheckCyclicDependency(n Node, name string) error {
	for _, d := range n.Dependencies() {
		if d.Name() == name {
			return fmt.Errorf("cyclic dependency! node %s has dependency on node %s", n.Name(), name)
		}
		if err := CheckCyclicDependency(d, name); err != nil {
			return err
		}
	}
	return nil
}
