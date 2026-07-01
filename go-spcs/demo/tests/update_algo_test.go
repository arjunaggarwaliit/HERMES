package tests

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"src"
	"strconv"
	"strings"
	"testing"
)

type UpdateQuery struct {
	ID        int
	StartID   int
	EndID     int
	NewWeight float64
}


func ReadUpdateQueriesFromCSV(filePath string) ([]UpdateQuery, error) {
	var testData []UpdateQuery
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

		testData = append(testData, UpdateQuery{id, startId, endId, expectedCost})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return testData, nil
}

func TestUpdateAlgo(t *testing.T){
	m := src.NewMLP(levels)
	m.MLPConstruction(g)
	//m.PrintInfo()
	
	batch := []src.UpdateQueryStruct{
        {StartId: 67472 , EndId: 67554 , Weight: 10},
        {StartId: 8, EndId: 7, Weight: 10},
    }
	src.BatchUpdateQuery(batch, m)
	// m.PrintInfo()

	startRouteId := 10
	endRouteId := 12

	path, dist, _ := src.RoutingQuery(startRouteId, endRouteId, m)
	fmt.Printf("Shortest path from %d to %d: %v\n", startRouteId, endRouteId, path)
	fmt.Printf("Shortest distance from %d to %d: %v\n", startRouteId, endRouteId, dist)
}

func TestUpdateAlgoTwo(t *testing.T){
	
	// Create or open the file for writing
	outputFile, _ := os.Create("test_output.txt")

	log.SetOutput(outputFile)

	m := src.NewMLP(levels)
	m.MLPConstruction(g)
	//m.PrintInfo()
	var updateBatch []src.UpdateQueryStruct
	
	batch1 , _ := ReadUpdateQueriesFromCSV("update_queries_1000.csv")
	for _, query := range batch1 {
		updateBatch = append(updateBatch, src.UpdateQueryStruct{
			StartId: query.StartID,
			EndId:   query.EndID,
			Weight:  int(query.NewWeight),
		})
		// time.Sleep(100 * time.Millisecond) // Adjust sleep time as needed
	}

	src.BatchUpdateQuery(updateBatch, m)
	fmt.Print("Batch 1 processed\n")
}