package tests

import (
	"fmt"
	"src"
	"testing"
)

func TestAPSP(t *testing.T) {

	startID := 9
	endID := 12

	src.AllPairsShortestPath(g)
	path := src.GetShortestPath(g, startID, endID)

	// Print the shortest path
	if len(path) == 0 {
		fmt.Printf("No path exists from %d to %d\n", startID, endID)
	} else {
		fmt.Printf("Shortest path from %d to %d: %v\n", startID, endID, path)
	}
}

func TestAPSPNE(t *testing.T) {

	startID := 9
	endID := 12
	
	nodes := g.GetNodes()
	edges := g.GetEdges()
	parition := src.NewPartition()

	src.AllPairsShortestPathNE(nodes, edges, parition)

	asps_dist := parition.Apsp_Dist

	// Print the shortest path
	fmt.Printf("Shortest path from %d to %d: %v\n", startID, endID, asps_dist[startID][endID])
}
