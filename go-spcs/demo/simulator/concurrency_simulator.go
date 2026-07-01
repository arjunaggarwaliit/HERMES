package simulator

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"src"
	"strconv"
	"strings"
	"sync"
	"time"
)

// RoutingQuery represents a single query with an ID, startID, endID, and ExpectedCost

type UpdateQueryTestResult struct {
	QueryID int
	StartID int
	EndID   int
	RunTime time.Duration
}

type UpdateQuery struct {
	ID        int
	StartID   int
	EndID     int
	NewWeight float64
}

// var routingCounter int
// var updateCounter int
// var n int = 5
var id int64 = 0
var globalStartTime = time.Now()

func ReadUpdateQueriesFromCSV(filePath string) ([]UpdateQuery, error) {
	var testData []UpdateQuery
	id := 0

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	log.Println("Reading queries from file")

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
func runConcurrencyTestSPCS(mlp *src.MLP, testConfigs []TestConfig, resultsDir string) {
    // Create or open the file for writing
    outputFile, err := os.Create(filepath.Join(resultsDir, "concurrency/concurrency_results_1.txt"))
    if err != nil {
        log.Println("Error creating/opening output file:", err)
        return
    }
    defer outputFile.Close() // Make sure to close the file when done

    // Redirecting log output to the file
    log.SetOutput(outputFile)

    // Create the CSV file for recording observations
    if err := createCSVFile(filepath.Join(resultsDir, "concurrency/exec_times.csv")); err != nil {
        log.Println("Error creating CSV file:", err)
        return
    }

    for _, config := range testConfigs {
        if config.Name == "concurrency_test" {
            fmt.Println("Running concurrency test")

            routingQueries, err := ReadRoutingQueriesFromCSV(config.RoutingQueriesPath)
            if err != nil {
                log.Println("Error reading routing queries:", err)
                return
            }

            updateQueries, err := ReadUpdateQueriesFromCSV(config.UpdateQueriesPath)
            if err != nil {
                log.Println("Error reading update queries:", err)
                return
            }

            var wg sync.WaitGroup
            wg.Add(2)

            rand.Seed(time.Now().UnixNano())

            system := src.NewConcurrencySystem(mlp, 3)
            system.Start()

            go func() {
                defer wg.Done()
                for _, query := range routingQueries {
                    system.AddRoutingQuery(query.ID, query.StartID, query.EndID)                    
                    time.Sleep(randomSleepDuration()) // Sleep for a random duration before firing the next update query

                }
            }()

            go func() {
                defer wg.Done()
            
             
                for _, query := range updateQueries {
                    
                    system.AddUpdateQuery(query.ID, query.StartID, query.EndID, int(query.NewWeight))
            
                    time.Sleep(randomSleepDuration()) // Sleep for a random duration before firing the next routing query
                }
            }()

            wg.Wait()
        }
    }
}

func createCSVFile(filePath string) error {
    // Create the file and write the header row
    file, err := os.Create(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    // Write the header row
    header := []string{"id", "query_type", "time_taken"}
    writer := csv.NewWriter(file)
    defer writer.Flush()
    if err := writer.Write(header); err != nil {
        return err
    }

    return nil
}

func writeToCSV(resultsDir string ,queryType string, timeTaken time.Duration) error {
    filePath := filepath.Join(resultsDir, "concurrency/exec_times.csv")
    var file *os.File

    // Open the file in append mode
    file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    // Write the observation to the file
    writer := csv.NewWriter(file)
    defer writer.Flush()

    // Generate a unique ID for the record
    id++

    // Write the observation row
    if err := writer.Write([]string{strconv.FormatInt(id, 10), queryType, timeTaken.String()}); err != nil {
        return err
    }

    return nil
}


// Function to generate a random sleep duration between 1 and 100 milliseconds
func randomSleepDuration() time.Duration {
	sleepTime := rand.Intn(100) + 1
	return time.Duration(sleepTime) * time.Millisecond
}
