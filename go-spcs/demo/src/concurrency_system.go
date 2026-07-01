package src

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

type RoutingQueryParams struct {
	ID 	int
	startId int
	endId   int
}

type ConcurrencySystem struct {
	mlp      *MLP
	interval int

	RoutingQueryDeque []RoutingQueryParams
	RoutingQueryMutex sync.Mutex
	UpdateQueryDeque []UpdateQueryStruct
	UpdateQueryMutex  sync.Mutex

	State int
	wg    sync.WaitGroup

	log       []string
	startTime time.Time 

}	

var csv_id int64 = 0
var csv_file string = "../../results/exec_times_v1.csv"

func NewConcurrencySystem(mlp *MLP, interval int) *ConcurrencySystem {
	
	outputFile, _ := os.Create("../../results/log_v1.txt")
	_ = createCSVFile(csv_file)

	log.SetOutput(outputFile)
	return &ConcurrencySystem{
		mlp:      mlp,
		interval: interval,
		startTime: time.Now(),
	}
}

func (cs *ConcurrencySystem) Start() {
	cs.State = 0
	cs.startTime = time.Now()
	go cs.Run()
}

func (cs *ConcurrencySystem) Stop() {
	cs.State = -1
}

func (cs *ConcurrencySystem) AddRoutingQuery(ID int, startId int, endId int) {
	cs.RoutingQueryMutex.Lock()
	cs.RoutingQueryDeque = append(cs.RoutingQueryDeque, RoutingQueryParams{ID, startId, endId})
	cs.RoutingQueryMutex.Unlock()
}

func (cs *ConcurrencySystem) AddUpdateQuery(ID int, startId int, endId int, weight int) {
	cs.UpdateQueryMutex.Lock()
	cs.UpdateQueryDeque = append(cs.UpdateQueryDeque, UpdateQueryStruct{ID, startId, endId, weight})
	cs.UpdateQueryMutex.Unlock()
}

func (cs *ConcurrencySystem) ExecuteUpdateQuery() {
	// Pop the first element from the update query stack

	startQueryTime := time.Now()
	updateQueries := make([]UpdateQueryStruct, len(cs.UpdateQueryDeque))

	
	for i := range cs.UpdateQueryDeque {
		// Pop the elements from the stack in reverse order
		cs.UpdateQueryMutex.Lock()
		updateQueries[i] = cs.UpdateQueryDeque[len(cs.UpdateQueryDeque)-1-i]
		log.Printf("[UPDATE QUERY START]  %v :  %v\n", updateQueries[i].ID, time.Since(cs.startTime) )
		cs.UpdateQueryMutex.Unlock()
	}

	// Empty out the UpdateQueryDeque
	cs.UpdateQueryDeque = nil

	BatchUpdateQuery(updateQueries, cs.mlp)
	endQueryTime := time.Now()

	for i := range updateQueries {
		log.Printf("[UPDATE QUERY END]  %v :  %v\n", updateQueries[i].ID, time.Since(cs.startTime))
	}

	// Time from the start of the system
	elapsed := time.Since(cs.startTime)

	writeToCSV("update", endQueryTime.Sub(startQueryTime))

	cs.log = append(cs.log, fmt.Sprintf("Update query executed. Time elapsed: %s, Query time: %s", elapsed, endQueryTime.Sub(startQueryTime)))
}

func (cs *ConcurrencySystem) ExecuteRoutingQuery() {
	// Pop the first element from the routing query stack

	cs.RoutingQueryMutex.Lock()

	if len(cs.RoutingQueryDeque) == 0 {
		cs.RoutingQueryMutex.Unlock()
		return
	}

	routingQueryParams := cs.RoutingQueryDeque[0]
	cs.RoutingQueryDeque = cs.RoutingQueryDeque[1:]
	cs.RoutingQueryMutex.Unlock()

	log.Printf("[ROUTING QUERY START]  %v :  %v\n", routingQueryParams.ID, time.Since(cs.startTime))
	
	_, dist, affected := RoutingQuery(routingQueryParams.startId, routingQueryParams.endId, cs.mlp)
	// Time from the start of the system
	startQueryTime := time.Now()
	elapsed := time.Since(cs.startTime)
	endQueryTime := time.Now()

	log.Printf("[ROUTING QUERY END]  %v :  %v\n", routingQueryParams.ID, time.Since(cs.startTime))   
	writeToCSV("routing", endQueryTime.Sub(startQueryTime))

	cs.log = append(cs.log, fmt.Sprintf("Routing query executed: %d %d, Answer: %f, Time elapsed: %s, Query time: %s", routingQueryParams.startId, routingQueryParams.endId, dist, elapsed, endQueryTime.Sub(startQueryTime)))

	if affected == 1 {
		cs.RoutingQueryDeque = append([]RoutingQueryParams{routingQueryParams}, cs.RoutingQueryDeque...)
	}
}

func (cs *ConcurrencySystem) ExecuteUpdateQuerys() {
	// TODO: Have to optimize this
	for len(cs.UpdateQueryDeque) > 0 {
		cs.ExecuteUpdateQuery()
	}
}

func (cs *ConcurrencySystem) ExecuteRoutingQuerys(channel chan int) {
	cs.wg.Add(1)
	timeout := time.NewTimer(time.Duration(cs.interval) * time.Second)
	defer timeout.Stop()

	fmt.Println("Routing query coroutine started!")

	for {
		select {
		case <-channel:
			cs.wg.Done()
			return
		case <-timeout.C:
			// Timeout reached
			cs.wg.Done()
			return
		default:
			if len(cs.RoutingQueryDeque) > 0 {
				cs.ExecuteRoutingQuery()
			}
		}
	}
}

func (cs *ConcurrencySystem) PrintLog() {
	for _, log := range cs.log {
		fmt.Println(log)
	}
}

func (cs *ConcurrencySystem) Run() {
	for cs.State != -1 {
		// Tell the current time
		if cs.State == 1 {
			fmt.Println(time.Now())
			quit := make(chan int)

			go cs.ExecuteRoutingQuerys(quit)

			time.Sleep(time.Duration(cs.interval) * time.Second)
			cs.wg.Wait()
			fmt.Println("Concurrency system - Routing Phase Ended!")
			if cs.State != -1 {
				cs.State = 0
			}
		} else if cs.State == 0 {
			fmt.Println("Concurrency system - Update Phase Started!")
			cs.ExecuteUpdateQuerys()
			if cs.State != -1 {
				cs.State = 1
			}
			fmt.Println("Concurrency system - Update Phase Ended!")
		}

		fmt.Println("State: ", cs.State)
		fmt.Println("Current Time Elapsed: ", time.Since(cs.startTime))
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

func writeToCSV(queryType string, timeTaken time.Duration) error {
    filePath := csv_file
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
    csv_id++

    // Write the observation row
    if err := writer.Write([]string{strconv.FormatInt(csv_id, 10), queryType, timeTaken.String()}); err != nil {
        return err
    }

    return nil
}
