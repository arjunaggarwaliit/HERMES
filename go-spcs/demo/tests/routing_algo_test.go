package tests

import (
	"bufio"
	"fmt"
	"math/rand"
	"encoding/csv"
	"os"
	"src"
	"strconv"
	"strings"
	"testing"
	"time"
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

func TestRoutingAlgo(t *testing.T){
	m := src.NewMLP(levels)
	m.MLPConstruction(g)

	// Save the MLP to a file
	// err := src.SaveMLPToJsonFile(m, "mlp_1000.json")
	// if err != nil {
	// 	t.Errorf("Error saving MLP to file: %v", err)
	// }

	//m, err := src.LoadMLPFromJSONFile("mlp_1000.json")

	//m.PrintInfo()

	// if err != nil {
	// 	t.Errorf("Error loading MLP from file: %v", err)
	// }

	// endId := 0
	// startId := 17
	
	// startTime := time.Now()
	// path, cost, _ := src.RoutingQuery(startId, endId, m)
	// elapsedTime := time.Since(startTime)
	// fmt.Printf("Routing cost from %d to %d is %f, Path: %d. Time elapsed: %v\n", startId, endId, cost, path, elapsedTime)


	

}
// TestRoutingAlgoTwo tests the routing algorithm with a set of test data from a CSV file
func TestRoutingAlgoTwo(t *testing.T){
	m := src.NewMLP(levels)
	m.MLPConstruction(g)

	// Read test data
	testData, err := ReadTestData("results.csv")
	if err != nil {
		t.Errorf("Error reading test data: %v", err)
	}

	// Run the tests
	testResult := true
	correctResults := 0
	totalTime := 0.0

	// Save the results to a file
	file, err := os.OpenFile("results/results_algo_two.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, test := range testData {
		start := time.Now()
		path, cost, _ := src.RoutingQuery(test.startId, test.endId, m)
		elapsed := time.Since(start)
		totalTime += elapsed.Seconds()

		if cost != test.expectedCost {
			t.Errorf("Routing cost from %d to %d is %f, expected %f", test.startId, test.endId, cost, test.expectedCost)
			testResult = false
		} else {
			correctResults++
		}
		
		fmt.Printf("Routing cost from %d to %d is %f, Path: %d , expected %f. Time elapsed: %v\n", test.startId, test.endId, cost, path, test.expectedCost, elapsed)
		record := []string{strconv.Itoa(test.startId), strconv.Itoa(test.endId), strconv.FormatFloat(cost, 'f', -1, 64), strconv.FormatFloat(test.expectedCost, 'f', -1, 64), strconv.FormatInt(elapsed.Microseconds(), 10)}
		err := writer.Write(record)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	}

	averageTime := totalTime / float64(len(testData))
	fmt.Printf("Average time elapsed for %d queries: %f seconds\n", len(testData), averageTime)
	fmt.Printf("Correct results: %d\n", correctResults)
	fmt.Printf("Routing tests result: %v\n", testResult)
}

// TestRoutingAlgoThree tests the routing algorithm with a set of test data generated randomly
func TestRoutingAlgoThree(t *testing.T){
	m := src.NewMLP(levels)
	m.MLPConstruction(g)

	// Store this data in a CSV file
	file, err := os.OpenFile("results_algo.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Randomly generate queries
	queriesSize := 1000
	verticesNum := 100000
	for i := 0; i < queriesSize; i++ {
		startId := rand.Intn(verticesNum)
		endId := rand.Intn(verticesNum)
		for startId == endId {
			endId = rand.Intn(verticesNum)
		}

		start := time.Now()
		_, cost, _ := src.RoutingQuery(startId, endId, m)
		elapsed := time.Since(start)

		// Write the result to the file
		record := []string{strconv.Itoa(startId), strconv.Itoa(endId), strconv.FormatFloat(cost, 'f', -1, 64), strconv.FormatInt(elapsed.Milliseconds(), 10)}
		err := writer.Write(record)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	}
}	

