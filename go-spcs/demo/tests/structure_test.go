package tests

import (
	"fmt"
	"math"
	"sort"
	"src"
	"testing"
)

func TestGraphCreation(t *testing.T) {

	// Get nodes and store their IDs
	nodeIDs := []int{}
	for id := range g.GetNodes() {
		nodeIDs = append(nodeIDs, id)
	}

	// Sort and compare node IDs with vertices
	sort.Ints(nodeIDs)
	sort.Ints(vertices)
	actualEdges := [][]int{}
	for _, id := range nodeIDs {
		node := g.GetNodes()[id]
		for _, edge := range node.GetEdges() {
			actualEdges = append(actualEdges, []int{node.GetId(), (edge.GetOtherNode(node)).GetId(), edge.GetWeight()})
		}
	}

	// Check if all edges are present
	for _, e := range edges {
		found := false
		for _, a := range actualEdges {
			if e[0] == a[0] && e[1] == a[1] && e[2] == a[2] {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Edge %v not found", e)
		}
	}

	t.Log("Graph creation is Ok!")
}

func TestMLPCreation(t *testing.T) {

	// Create a new MLP
	m := src.NewMLP(levels)
	m.MLPConstruction(g)

	// These would need to be adjusted based on your expectations and the specifics of how MLPConstruction works.
	if m.GetLevelNum() != levels {
		t.Errorf("Expected 3 levels, got %d", m.GetLevelNum())
	}

	// Check the number of partitions at the first level
	expectedPartitionsAtLevel0 := 1 // Assuming initial construction puts all nodes into a single partition at level 0
	if count := m.GetPartitionCountAtLevel(0); count != expectedPartitionsAtLevel0 {
		t.Errorf("Expected %d partitions at level 0, got %d", expectedPartitionsAtLevel0, count)
	}

	// Check the number of partitions at the last level
	expectedPartitionsAtLevelNum := math.Pow(2, float64(levels-1))
	if count := m.GetPartitionCountAtLevel(levels - 1); count != int(expectedPartitionsAtLevelNum) {
		t.Errorf("Expected %d partitions at level %d, got %d", int(expectedPartitionsAtLevelNum), levels-1, count)
	}

	for i := 0; i < m.GetLevelNum(); i++ {

		// Print the partitions at the level
		fmt.Println("Partitions at level", i, ":")

		for _, p := range m.GetPartitions(i) {

			// Get all the nodes in the partition
			nodes := p.GetNodes()
			nodeIDs := []int{}
			for id := range nodes {
				nodeIDs = append(nodeIDs, id)
			}

			// Sort and print the node IDs
			sort.Ints(nodeIDs)
			fmt.Println("Nodes in partition", p.GetId(), ":", nodeIDs)

			dist := p.Apsp_Dist
			for node := range p.Shortcut_Nodes {
				for otherNode := range p.Shortcut_Nodes {
					// Check if there is an entry in the distance matrix
					if _, ok := dist[node][otherNode]; !ok {
						continue	
					}
					fmt.Println("Distance from", node, "to", otherNode, ":", dist[node][otherNode])
				}
			}
		}
	}

	filePath := "graphs/mlp.dot"
	m.GenerateMLPDOT(filePath)

	// Save MLP 
	err := src.SaveMLPToJsonFile(m, "mlp.json")
	if err != nil {
		t.Errorf("Error saving MLP to file: %v", err)
	}


	t.Log("MLP creation is Ok!")
}

func TestSaveandLoad(t *testing.T) {
	fmt.Print("Test Save and Load")
	// Create a new MLP
	m := src.NewMLP(levels)
	m.MLPConstruction(g)

	// Print the MLP
	fmt.Println("MLP:")
	for i := 0; i < m.GetLevelNum(); i++ {
		fmt.Println("Level", i, ":")
		for _, p := range m.GetPartitions(i) {
			fmt.Println("Partition", p.GetId(), ":")
			fmt.Println("Shortcut Nodes:", p.Shortcut_Nodes)
			fmt.Println("Shortcut Edges:", p.Shortcut_Edges)
			fmt.Println("Border Nodes:", p.Border_Nodes)
			fmt.Println("Border Edges:", p.Border_Edges)
		}
	}


	// Save MLP
	err := src.SaveMLPToJsonFile(m, "mlp.json")
	if err != nil {
		t.Errorf("Error saving MLP to file: %v", err)
	}

	// Load MLP
	loadedMLP, err := src.LoadMLPFromJSONFile("mlp.json")
	if err != nil {
		t.Errorf("Error loading MLP from file: %v", err)
	}

	// Save the loaded MLP	
	err1 := src.SaveMLPToJsonFile(loadedMLP, "mlp_loaded.json")
	if err1 != nil {
		t.Errorf("Error saving loaded MLP to file: %v", err)
	}
	
}
