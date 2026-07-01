package simulator

import (
	"bufio"
	"fmt"
	"os"

	//"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"utils"
)


func ReadGraphDCH(graph *utils.Graph, fname string) error {
	file, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
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

		err = graph.CreateVertex(source)
		if err != nil {
			return err
		}
		err = graph.CreateVertex(target)
		if err != nil {
			return err
		}

		// Assuming the distance is used as the weight for the edge
		err = graph.AddEdge(source, target, weight)
		if err != nil {
			return err
		}

		if oneway == "B" {
			err = graph.AddEdge(target, source, weight)
			if err != nil {
				return err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}


func ReadGraphSPCS(vertices *[]int, edges *[][]int, fname string) error {
	// Write sequential code to read the graph from the text file
	file, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		err := processLineSPCS(line, vertices, edges)
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func processLineSPCS(line string, vertices *[]int, edges *[][]int) error {
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
