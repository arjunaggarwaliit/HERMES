package tests

import (
	"src"
	"testing"
	"time"
)

func TestConcurrencySystem(t *testing.T){
	m := src.NewMLP(levels)
	m.MLPConstruction(g)

	system := src.NewConcurrencySystem(m, 2)
	system.Start()

	// After 10 seconds, stop the system
	go func(){
		time.Sleep(10 * time.Second)
		system.Stop()
	}()

	// Update Query
	ID := 1
	startUpdateId := 7
	endUpdateId := 8
	weight := 10
	system.AddUpdateQuery(ID, startUpdateId, endUpdateId, weight)

	// Routing Query
	ID = 2
	startRouteId := 0
	endRouteId := 17
	system.AddRoutingQuery(ID, startRouteId, endRouteId)
	
	count := 0

	for system.State != -1 {
		for i := 0; i < 3; i++ {
			ID++
			startRouteId = 0
			endRouteId = i + 3
			system.AddRoutingQuery(ID, startRouteId, endRouteId)
		}
		count++
		time.Sleep(1 * time.Second)

		if count == 3 {
			if weight == 10 {
				weight = 1
			} else {
				weight = 10
			}
			ID++
			system.AddUpdateQuery(ID, startUpdateId, endUpdateId, weight)
			count = 0
		}

		//fmt.Println("State: ", system.State)
	}

	// Print the log
	system.PrintLog()
}