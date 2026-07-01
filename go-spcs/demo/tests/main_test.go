package tests

import (
	"bufio"
	"fmt"
	"os"

	//"runtime"
	"slices"
	"src"
	"strconv"
	"strings"
	"sync"
	"testing"
)

func sample_graphFromText(vertices *[]int, edges *[][]int, fname string) error {
	// Write sequential code to read the graph from the text file
	file, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		err := processLine(line, vertices, edges)
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func processLine(line string, vertices *[]int, edges *[][]int) error {
	fields := strings.Split(line, ";")
	if len(fields) < 4 {
		return fmt.Errorf("invalid line format: %s", line)
	}

	source, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return err
	}
	target, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return err
	}

	oneway := fields[2]
	weight, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return err
	}

	// Use a mutex to synchronize access to shared data
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()

	// Check if the source and target vertices are in the vertices array
	if !slices.Contains(*vertices, int(source)) {
		*vertices = append(*vertices, int(source))
	}

	if !slices.Contains(*vertices, int(target)) {
		*vertices = append(*vertices, int(target))
	}

	// Add the edge to the edges array
	*edges = append(*edges, []int{int(source), int(target), int(weight)})

	if oneway == "B" {
		// Add the reverse edge if it's bidirectional
		*edges = append(*edges, []int{int(target), int(source), int(weight)})
	}

	return nil
}

var (
	g   *src.Graph // Global graph variable
	err error      // Error variable
	// Nodes
	// vertices = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}

	// Graph 1:
	// // Edges: Each edge is represented as a slice of three integers [source, destination, weight]
	// edges = [][]int{
	// 	{9, 12, 9}, {12, 9, 9},{8, 12, 9}, {12, 8, 9},{7, 10, 8}, {10, 7, 8},{0, 3, 10}, {3, 0, 10},{6, 10, 8}, {10, 6, 8},
	// 	{1, 4, 7}, {4, 1, 7},{8, 4, 4}, {4, 8, 4},{10, 7, 6}, {7, 10, 6},{7, 4, 9}, {4, 7, 9},{10, 7, 4}, {7, 10, 4},{14, 3, 6},
	// 	{3, 14, 6},{9, 14, 8}, {14, 9, 8},{5, 0, 1}, {0, 5, 1},{10, 2, 5}, {2, 10, 5}, {0, 4, 9}, {4, 0, 9}, {3, 4, 12}, {4, 3, 12},
	// 	{14, 13, 1}, {13, 14, 1},{13, 0, 5}, {0, 13, 5},{4, 14, 3}, {14, 4, 3},{3, 12, 7}, {12, 3, 7},{13, 4, 1}, {4, 13, 1},
	// 	{6, 5, 5}, {5, 6, 5},{14, 8, 7}, {8, 14, 7},{1, 7, 5}, {7, 1, 5},{13, 3, 5}, {3, 13, 5},{1, 13, 5}, {13, 1, 5},
	// 	{12, 5, 10}, {5, 12, 10},{14, 6, 6}, {6, 14, 6},{14, 8, 2}, {8, 14, 2},{6, 9, 9}, {9, 6, 9},{6, 8, 2}, {8, 6, 2},
	// 	{10, 9, 3}, {9, 10, 3},{11, 10, 4}, {10, 11, 4},{3, 0, 7}, {0, 3, 7},{1, 6, 10}, {6, 1, 10},{6, 12, 9}, {12, 6, 9},
	// 	{12, 1, 8}, {1, 12, 8},{14, 0, 7}, {0, 14, 7},{4, 8, 10}, {8, 4, 10},{11, 2, 10}, {2, 11, 10},
	// }

	// Graph 2 (Used in slides):
	// vertices = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17}
	// edges = [][]int{
	// 	{0, 2, 7}, {2, 0, 7}, {1, 2, 7}, {2, 1, 7},{1, 3, 8}, {3, 1, 8},{1, 4, 9}, {4, 1, 9},{1, 12, 6}, {12, 1, 6},{2, 3, 5}, {3, 2, 5},
	// 	{2, 9, 7}, {9, 2, 7},{4, 5, 8}, {5, 4, 8},{4, 6, 9}, {6, 4, 9},{4, 7, 5}, {7, 4, 5},{4, 9, 6}, {9, 4, 6},{5, 6, 7}, {6, 5, 7},
	// 	{5, 7, 8}, {7, 5, 8},{6, 7, 6}, {7, 6, 6},{6, 17, 5}, {17, 6, 5},{7, 17, 7}, {17, 7, 7},{8, 9, 8}, {9, 8, 8},{8, 11, 6}, {11, 8, 6},
	// 	{9, 10, 9}, {10, 9, 9}, {10, 11, 7}, {11, 10, 7},{11, 12, 8}, {12, 11, 8},{12, 13, 6}, {13, 12, 6},{12, 15, 5}, {15, 12, 5},
	// 	{13, 15, 7}, {15, 13, 7},{13, 14, 9}, {14, 13, 9},{14, 15, 8}, {15, 14, 8},{14, 16, 5}, {16, 14, 5},{15, 16, 6}, {16, 15, 6},
	//   }

	// Graph 3 (Straight Line Graph):
	// vertices = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}
	// edges = [][]int{
	// 	{0, 1, 1}, {1, 0, 1},{1, 2, 1}, {2, 1, 1},{2, 3, 1}, {3, 2, 1},{3, 4, 1}, {4, 3, 1},{4, 5, 1}, {5, 4, 1},
	// 	{5, 6, 1}, {6, 5, 1},{6, 7, 1}, {7, 6, 1},{7, 8, 1}, {8, 7, 1},{8, 9, 1}, {9, 8, 1},{9, 10, 1}, {10, 9, 1},
	// 	{10, 11, 1}, {11, 10, 1},{11, 12, 1}, {12, 11, 1},{12, 13, 1}, {13, 12, 1},{13, 14, 1}, {14, 13, 1},{14, 15, 1}, {15, 14, 1},
	// 	{15, 16, 1}, {16, 15, 1},{16, 17, 1}, {17, 16, 1},{17, 18, 1}, {18, 17, 1},{18, 19, 1}, {19, 18, 1},
	// }

	// Graph 4 (US Road Network):
	vertices = []int{}
	edges    = [][]int{}

	// Number of levels in the MLP
	levels = 13
)

func TestMain(m *testing.M) {

	// Setup code: Initialize your graph here
	g = src.NewGraph()

	// // Read the graph from the text file: If using the US Road Network, use the following line
	err = sample_graphFromText(&vertices, &edges, "../../../datasets/master/data_100000.txt")
	fmt.Println("Reading the graph from the text file done. No of vertices: ", len(vertices), " No of edges: ", len(edges))

	//g.CreateSampleGraph(15, 30)
	g.CreateGraph(vertices, edges)

	// Visualize the Graph
	// g.GenerateGraphDot("graphs/graph.dot")

	// Run the tests
	code := m.Run()

	// Cleanup code if needed
	os.Exit(code)
}
