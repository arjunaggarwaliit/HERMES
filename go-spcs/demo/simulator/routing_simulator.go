package simulator

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"src"
	"strconv"
	"strings"
	"time"
	"utils"
)

// RoutingQuery represents a single query with an ID, startID, endID, and ExpectedCost

type RoutingTestResult struct {
	QueryID      int
	StartID      int
	EndID        int
	ExpectedCost float64
	RunTime      time.Duration
	Cost         float64 // Cost of the routing
}

type RoutingQuery struct {
	ID           int
	StartID      int
	EndID        int
	ExpectedCost float64
}

func ReadRoutingQueriesFromCSV(filePath string) ([]RoutingQuery, error) {
	var testData []RoutingQuery
	id := 0

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fmt.Println("Reading queries from file")

	// Read the file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		id += 1
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

		testData = append(testData, RoutingQuery{id, startId, endId, expectedCost})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return testData, nil
}

// runRoutingTest executes routing test queries and stores the results in a CSV file
func runRoutingTestSPCS(mlp *src.MLP, testConfigs []TestConfig, resultsDir string) {

	var results []RoutingTestResult

	for _, config := range testConfigs {
		if config.Name == "routing_test" {
			fmt.Println("Running routing test")
			// Read queries from the file
			queries, _ := ReadRoutingQueriesFromCSV(config.RoutingQueriesPath)

			// Execute routing test for each query
			for _, query := range queries {
				startTime := time.Now()

				fmt.Printf("Running query %d: StartID=%d, EndID=%d\n", query.ID, query.StartID, query.EndID)

				_, cost, _ := src.RoutingQuery(query.StartID, query.EndID, mlp)

				runTime := time.Since(startTime)

				// Store the result
				results = append(results, RoutingTestResult{
					QueryID:      query.ID,
					StartID:      query.StartID,
					EndID:        query.EndID,
					ExpectedCost: query.ExpectedCost,
					RunTime:      runTime,
					Cost:         float64(cost), // Convert cost to float64
				})
			}

			// Assuming config.ResultsDir contains the directory path
			resultsFile := filepath.Join(resultsDir, "routing/routing_results.csv")

			// Write results to CSV file incrementally
			file, _ := os.OpenFile(resultsFile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)

			defer file.Close()

			writer := csv.NewWriter(file)
			defer writer.Flush()

			// Write header if file is empty
			fileInfo, _ := file.Stat()
			if fileInfo.Size() == 0 {
				writer.Write([]string{"QueryID", "StartID", "EndID", "ExpectedCost", "RunTime", "Cost"})
			}

			// Write results to CSV
			for _, result := range results {
				// Convert path slice to a comma-separated string

				record := []string{
					strconv.Itoa(result.QueryID),
					strconv.Itoa(result.StartID),
					strconv.Itoa(result.EndID),
					strconv.Itoa(int(result.ExpectedCost)),
					result.RunTime.String(),
					strconv.FormatFloat(result.Cost, 'f', -1, 64), // Convert float64 to string
				}
				if err := writer.Write(record); err != nil {
					log.Println("Error writing record to CSV:", err)
				}
			}

		}
	}

}

func runRoutingTestDCH(g utils.Graph, testConfigs []TestConfig, resultsDir string) {

	var results []RoutingTestResult

	for _, config := range testConfigs {
		if config.Name == "routing_test" {
			fmt.Println("Running routing test DCH")
			// Read queries from the file
			queries, _ := ReadRoutingQueriesFromCSV(config.RoutingQueriesPath)

			// Execute routing test for each query
			for _, query := range queries {
				startTime := time.Now()

				fmt.Printf("Running query %d: StartID=%d, EndID=%d\n", query.ID, query.StartID, query.EndID)

				cost, _ := g.ShortestPath(int64(query.StartID), int64(query.EndID))

				runTime := time.Since(startTime)

				// Store the result
				results = append(results, RoutingTestResult{
					QueryID:      query.ID,
					StartID:      query.StartID,
					EndID:        query.EndID,
					ExpectedCost: query.ExpectedCost,
					RunTime:      runTime,
					Cost:         float64(cost), // Convert cost to float64
				})
			}

			// Assuming config.ResultsDir contains the directory path
			resultsFile := filepath.Join(resultsDir, "routing/routing_results.csv")

			// Write results to CSV file incrementally
			file, _ := os.OpenFile(resultsFile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)

			defer file.Close()

			writer := csv.NewWriter(file)
			defer writer.Flush()

			// Write header if file is empty
			fileInfo, _ := file.Stat()
			if fileInfo.Size() == 0 {
				writer.Write([]string{"QueryID", "StartID", "EndID", "ExpectedCost", "RunTime", "Cost"})
			}

			// Write results to CSV
			for _, result := range results {
				// Convert path slice to a comma-separated string

				record := []string{
					strconv.Itoa(result.QueryID),
					strconv.Itoa(result.StartID),
					strconv.Itoa(result.EndID),
					strconv.Itoa(int(result.ExpectedCost)),
					result.RunTime.String(),
					strconv.FormatFloat(result.Cost, 'f', -1, 64), // Convert float64 to string
				}
				if err := writer.Write(record); err != nil {
					log.Println("Error writing record to CSV:", err)
				}
			}

		}
	}

}
