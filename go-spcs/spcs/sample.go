package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"

	//"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"utils" // Assuming this package exists and is correctly imported
)

const (
	eps1 = 0.0001
)

type TestRoutingQuery struct {
	startId int
	endId int
	expectedCost float64
}

func ReadTestData(filePath string) ([]TestRoutingQuery, error) {
	var testData []TestRoutingQuery

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ",")
		if len(fields) != 3 {
			return nil, fmt.Errorf("invalid line format: %s", line)
		}

		startId, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, err
		}

		endId, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, err
		}

		expectedCost, err := strconv.ParseFloat(fields[2], 64)
		if err != nil {
			return nil, err
		}

		testData = append(testData, TestRoutingQuery{startId, endId, expectedCost})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return testData, nil
}

func main() {
	g := utils.Graph{}

	nodes := "100000"
	err := sample_graphFromText(&g, fmt.Sprintf("C:/MasterFolder/General/CODING/New_Projects/shortest-path-concurrent-system/datasets/master/data_%s.txt", nodes))
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1) // Exit the program because of an error
	}

	fmt.Println("Please wait until contraction hierarchy is prepared")
	timech := time.Now()
	g.PrepareContractionHierarchies()
	elapsedch := time.Since(timech)
	fmt.Println("Time for CH:", elapsedch)

	// Find the shortest path
	fmt.Println("ShortestPath Query is starting...")

	// u := 0
	// v := 17

	// timequery := time.Now()
	// ans, path := g.ShortestPath(int64(u), int64(v))
	// elapsedquery := time.Since(timequery)

	// fmt.Println("Time for query:", elapsedquery)
	// fmt.Printf("Shortest path from %d to %d is: %v\n", u, v, path)
	// fmt.Printf("Cost of path: %f\n", ans)

	numPairs := 10000
	nodePairs := make([][]int64,0)
	// Convert the number of nodes to an integer
	vertices, _ := strconv.Atoi(nodes)

	// Store this data in a CSV file
	file, err := os.OpenFile("results.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	

	for i := 0; i < numPairs; i++ {
		// Randomly select a pair of nodes
		u := rand.Int63n(int64(vertices))
		v := rand.Int63n(int64(vertices))

		for u == v {
			v = rand.Int63n(int64(vertices))
		}
	
		timequery := time.Now()
		ans, path := g.ShortestPath(u, v)
		elapsedquery := time.Since(timequery)

		// Check if the node pair is already in the list
		if ans < 0 {
			i--
			continue
		}

		if len(nodePairs) > 0 {
			for _, pair := range nodePairs {
				//print(pair)
				if pair[0] == u && pair[1] == v {
					i--
					continue
				}
			}
		}

		nodePairs = append(nodePairs, []int64{u, v})
		fmt.Println("Time for query:", elapsedquery)
		fmt.Printf("Shortest path from %d to %d is: %v\n", u, v, path)
		fmt.Printf("Cost of path: %f\n", ans)

		

		err = writer.Write([]string{strconv.FormatInt(u, 10), strconv.FormatInt(v, 10), strconv.FormatFloat(ans, 'f', -1, 64)})

		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	}

	// Read the CSV file
	testData, err := ReadTestData("results.csv")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Write the results to a file
	file_new, err := os.OpenFile("results_algo_ch.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer file_new.Close()

	writer_new := csv.NewWriter(file)
	defer writer_new.Flush()

	for _, test := range testData {
		start := time.Now()
		ans, _ := g.ShortestPath(int64(test.startId), int64(test.endId))
		elapsed := time.Since(start)

		if ans < 0 {
			continue
		}

		if ans != test.expectedCost {
			print("Error")
		}

		// Save the results to a file
		record := []string{strconv.Itoa(test.startId), strconv.Itoa(test.endId), strconv.FormatFloat(ans, 'f', -1, 64), strconv.FormatFloat(test.expectedCost, 'f', -1, 64), strconv.FormatInt(elapsed.Microseconds(), 10)}
		err := writer.Write(record)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	}

}

func sample_graphFromCSV(graph *utils.Graph, fname string) error {
	file, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := csv.NewReader(bufio.NewReader(file))

	reader.Comma = ';'
	// reader.LazyQuotes = true

	// Read header
	_, err = reader.Read()
	if err != nil {
		return err
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		source, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			return err
		}
		target, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			return err
		}

		oneway := record[2]
		weight, err := strconv.ParseFloat(record[3], 64)
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
	return nil
}

func sample_graphFromText(graph *utils.Graph, fname string) error {
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
