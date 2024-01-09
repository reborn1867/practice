package main

import (
	"fmt"

	"github.wdf.sap.corp/practice-learning/dependency/dependency"
)

func main() {
	testList := []dependency.DependencyRelation{
		{
			Name: "a",
			Deps: []string{"b"},
		},
		{
			Name: "b",
		},
		{
			Name: "c",
			Deps: []string{"b", "e"},
		},
		{
			Name: "d",
			Deps: []string{"e"},
		},
		{
			Name: "e",
		},
	}

	nodeMap, err := dependency.ManageDependency(testList)
	if err != nil {
		fmt.Printf("failed to manage dependency: %s\n", err)
		return
	}

	for _, node := range nodeMap {
		var l1, l2 []string
		for _, d := range node.Dependencies() {
			l1 = append(l1, d.Name())
		}
		for _, d := range node.Depended() {
			l2 = append(l2, d.Name())
		}
		fmt.Printf("name: %s, depend on: %+v, depended by %+v\n", node.Name(), l1, l2)
	}

	fmt.Println("\ninstallation test")
	for _, node := range nodeMap {
		if dependency.IsRoot(node) {
			fmt.Printf("run install task in node %s \n", node.Name())
		} else {
			fmt.Printf("skip task in node %s due to dependencies\n", node.Name())
		}
	}

	fmt.Println("\nuninstallation test")
	for _, node := range nodeMap {
		if dependency.IsBottom(node) {
			fmt.Printf("run install task in node %s \n", node.Name())
		} else {
			fmt.Printf("skip task in node %s due to dependencies\n", node.Name())
		}
	}
}
