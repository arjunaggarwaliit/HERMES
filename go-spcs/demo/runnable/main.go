package main

import (
	"src"
	"simulator"
)

func ExampleNewGraph(vertices []int, edges [][]int) *src.Graph {
	g := src.NewGraph()

	// Add nodes
	for _, v := range vertices {
		g.AddNode(v)
	}

	// Add edges
	for _, e := range edges {
		g.AddEdge(e[0], e[1], e[2])
	}

	return g
}

func main() {

	simulator.Simulate()
}
