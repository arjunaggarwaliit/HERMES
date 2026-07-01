package src

import (
	"container/heap"
	"fmt"
	"math"
)

// Node represents a node in the graph
type NodeContainer struct {
    ID       int
    Distance float64
}



// MinHeap implements heap.Interface for []*Node based on Node.Distance
type MinHeap []*NodeContainer

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i].Distance < h[j].Distance }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x interface{}) {
    *h = append(*h, x.(*NodeContainer))
}

func (h *MinHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[0 : n-1]
    return x
}

// VanillaDijkstra finds the shortest path using Dijkstra's algorithm
func VanillaDijkstra(mlp *MLP, startID int, endID int) ([]int, float64) {
    // Initialize the distance and parent maps
    dist := make(map[int]float64)
    parent := make(map[int]int)

    // Initialize the queue
    pq := make(MinHeap, 0)
    heap.Init(&pq)

    // Set the distance of startID to 0 and push it to the priority queue
    dist[startID] = 0
    parent[startID] = -1
    heap.Push(&pq, &NodeContainer{ID: startID, Distance: 0})
	nodesChecked := 0

    for len(pq) > 0 {
		nodesChecked++
		fmt.Print("Nodes Checked: ", nodesChecked, "\n")

        // Pop the node with the smallest distance
        u := heap.Pop(&pq).(*NodeContainer)

        // Break if we reached the end node
        if u.ID == endID {
            break
        }

        // Get the level of the current node
		level := mlp.LevelDijkstra(startID, endID, u.ID)

        // Print the current node and its distance
        //fmt.Printf("Current node: %d, Distance: %f, Level: %d\n", u.ID, dist[u.ID], level)

		// Iterate over the neighbors of the current node
		neighbors := mlp.GetNeighbors(u.ID, level)

		for _, neighbor := range neighbors {
			// If the dist[neighbor.ID] is not set, set it to infinity
            if _, ok := dist[neighbor.ID]; !ok {
                dist[neighbor.ID] = math.Inf(1)
            }

            // Calculate the new distance
            alt := dist[u.ID] + neighbor.Distance

            // If the new distance is less than the current distance, update it
            if alt < dist[neighbor.ID] {
                dist[neighbor.ID] = alt
                parent[neighbor.ID] = u.ID

                // Print the updated distance
                //fmt.Printf("Updated distance to node %d: %f\n", neighbor.ID, alt)

                heap.Push(&pq, &NodeContainer{ID: neighbor.ID, Distance: alt})
            }
        }
    }

	fmt.Println("Dijkstra done")

    // Reconstruct the shortest path
    path := make([]int, 0)
	fmt.Println("End ID: ", endID)
    for u := endID; u != -1; u = parent[u] {
		fmt.Println("Parent: ", parent[u])
		if parent[u] == 0 {
			break
		}
        path = append(path, u)
    }

	fmt.Println("Path reconstruction done")

    // Reverse the path
    for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
        path[i], path[j] = path[j], path[i]
    }

	path = pathUnpacker(mlp, startID, endID, path)
	fmt.Println("Path unpacking done")

    // Return the shortest path and its distance
    return path, dist[endID]
}

// BidirectionalDijkstra finds the shortest path using bidirectional Dijkstra's algorithm
func BidirectionalDijkstra(mlp *MLP, startID int, endID int) ([]int, float64,int) {
	// Check if the start ID is there in the graph
	if _, ok := mlp.nodeMap[startID]; !ok {
		return []int{}, 0, 1
	}

	// Check if the end ID is there in the graph
	if _, ok := mlp.nodeMap[endID]; !ok {
		return []int{}, 0, 1
	}

	// Initialize the forward and backward distances and parents
	forwardDist := make(map[int]float64)
	backwardDist := make(map[int]float64)
	forwardParent := make(map[int]int)
	backwardParent := make(map[int]int)
	allVisited := []int{}

	// Initialize the forward and backward priority queues
	forwardPQ := make(MinHeap, 0)
	backwardPQ := make(MinHeap, 0)
	heap.Init(&forwardPQ)
	heap.Init(&backwardPQ)

	// Set the distances of startID and endID to 0 and push them to their respective priority queues
	forwardDist[startID] = 0
	backwardDist[endID] = 0
	forwardParent[startID] = -1
	backwardParent[endID] = -1
	heap.Push(&forwardPQ, &NodeContainer{ID: startID, Distance: 0})
	heap.Push(&backwardPQ, &NodeContainer{ID: endID, Distance: 0})

	// LocksAcquired is an array of edge id that are locked in order
	locksAcquired := make([]int, 0)
	queryStatus := true
	// Initialize variables to track common meeting point and shortest distance
	commonNode := -1
	shortestDistance := math.Inf(1)

	// For logging purposes
	nodesChecked := 0

	for len(forwardPQ) > 0 && len(backwardPQ) > 0 {
		nodesChecked++

		// Perform forward search
		forwardNodeContainer := heap.Pop(&forwardPQ).(*NodeContainer)
		forwardNodeContainerID := forwardNodeContainer.ID
		forwardNodeContainerNode := mlp.GetNodeByID(forwardNodeContainerID)
		forwardNodeContainerBasePartition := forwardNodeContainerNode.basePartition
		forwardNodeContainerAffected := forwardNodeContainerBasePartition.Affected
		
		if forwardNodeContainerAffected {
			queryStatus = false
			break
		}

		// Perform backward search
		backwardNodeContainer := heap.Pop(&backwardPQ).(*NodeContainer)
		backwardNodeContainerID := backwardNodeContainer.ID
		backwardNodeContainerNode := mlp.GetNodeByID(backwardNodeContainerID)
		backwardNodeContainerBasePartition := backwardNodeContainerNode.basePartition
		backwardNodeContainerAffected := backwardNodeContainerBasePartition.Affected

		if backwardNodeContainerAffected  {
			queryStatus = false
			break
		}

		// Update forward distances and parents
		updateNeighbors(mlp, startID, endID, forwardNodeContainer, forwardDist, forwardParent, &forwardPQ, &queryStatus, &locksAcquired, &allVisited)

		if !queryStatus{
			break
		}


		// Update backward distances and parents
		updateNeighbors(mlp, startID, endID, backwardNodeContainer, backwardDist, backwardParent, &backwardPQ, &queryStatus, &locksAcquired, &allVisited)

		if !queryStatus {
			break
		}

		// Check for common meeting point
		for _, node := range allVisited {
			// Check if the node is in both forward and backward visited nodes
			if _, ok := forwardDist[node]; !ok {
				continue
			}
			if _, ok := backwardDist[node]; !ok {
				continue
			}
			if forwardDist[node] + backwardDist[node] < shortestDistance {
				commonNode = node
				shortestDistance = forwardDist[node] + backwardDist[node]
			}
		}

		if commonNode != -1 {
			break
		}
	}

	if !queryStatus {
		// Release the locks
		for i := len(locksAcquired) - 1; i >= 0; i-- {
			mlp.edgeMap[locksAcquired[i]].RUnlock()
		}
		return []int{},0,1
	}

	// Reconstruct the shortest path
	path := make([]int, 0)
	for u := commonNode; u != -1; u = forwardParent[u] {
		path = append(path, u)
	}

    // Reverse the path
    for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
        path[i], path[j] = path[j], path[i]
    }

	// Add the backward path
	backwardPath := make([]int, 0)
	for u := backwardParent[commonNode]; u != -1; u = backwardParent[u] {
		backwardPath = append(backwardPath, u)
	}
    for i := 0 ; i < len(backwardPath); i++ {
        path = append(path, backwardPath[i])
    }

	path = pathUnpacker(mlp, startID, endID, path)

	// Release the locks
	for i := len(locksAcquired) - 1; i >= 0; i-- {
		mlp.edgeMap[locksAcquired[i]].RUnlock()
	}
	
	// Return the shortest path and its distance
	return path, shortestDistance , 0
}

func updateNeighbors(mlp *MLP, startID int, endID int, node *NodeContainer, dist map[int]float64, parent map[int]int, pq *MinHeap, queryStatus *bool, locksAcquired *[]int, allVisited *[]int) {
	level := mlp.LevelDijkstra(startID, endID, node.ID)
	neighbors := mlp.GetNeighbors(node.ID, level)
	for _, neighbor := range neighbors {
		// Lock the edge
		edgeId:= neighbor.EdgeID
		edge := mlp.GetEdgeByID(edgeId)
		// fmt.Println("Locking edge ", edgeId)
		result := edge.RLock()

		if !result {
			*queryStatus = false
			break
		} else {
			*locksAcquired = append(*locksAcquired, edgeId)
		}

		if _, ok := dist[neighbor.ID]; !ok {
			dist[neighbor.ID] = math.Inf(1)
		}
		alt := dist[node.ID] + neighbor.Distance
		if alt < dist[neighbor.ID] {
			dist[neighbor.ID] = alt
			parent[neighbor.ID] = node.ID
			heap.Push(pq, &NodeContainer{ID: neighbor.ID, Distance: alt})
			*allVisited = append(*allVisited, neighbor.ID)
		}
	}
}

func pathUnpacker(mlp *MLP, startID int, endID int, path []int) []int {
	// Path should be of the form [{startID, level}, {nodeID, level}, {nodeID, level}, ..., {endID, level}]
	if len(path) < 2 {
		return path
	}

	newPath := make([][] int, 0)
	firstLevelDijkstra := 0
	secondLevelDijkstra := mlp.LevelDijkstra(startID, endID, path[0])

	for i := 0; i < len(path)-1; i++ {
		firstLevelDijkstra = secondLevelDijkstra
		secondLevelDijkstra = mlp.LevelDijkstra(startID, endID, path[i+1])

		if firstLevelDijkstra != secondLevelDijkstra {
			newPath = append(newPath, []int{path[i], path[i+1], 0})
		} else {
			newPath = append(newPath, []int{path[i], path[i+1], firstLevelDijkstra})
		}
	}

	// Unpack the path
	unpackedPath := pathUnpackerRecursive(mlp, newPath)

	finalPath := make([]int, 0)
	for i := 0; i < len(unpackedPath); i++ {
		finalPath = append(finalPath, unpackedPath[i][0])
	}

	// Add the last node to the final path
	finalPath = append(finalPath, unpackedPath[len(unpackedPath)-1][1])

	// Return the unpacked path
	return finalPath
}


func pathUnpackerRecursive(mlp *MLP, path [][]int) [][]int {
	// Base case: if the path has only one node, return it
	if len(path) == 1 {
		return path
	}

	var isUnpacked bool = false

	newPath := make([][]int, 0)
	for i := 0; i < len(path); i++ {
		// If the levels of the two nodes are different, add the current node to the new path
		if path[i][2] != 0 {
			isUnpacked = true
			unPackEdge := mlp.UnpackEdge(path[i][0], path[i][1], path[i][2])
			for j := 0; j < len(unPackEdge); j++ {
				newPath = append(newPath, unPackEdge[j])
			}

		} else {
			newPath = append(newPath, []int{path[i][0], path[i][1], path[i][2]})
		}
	}

	// Return the unpacked path
	if isUnpacked {
		return pathUnpackerRecursive(mlp, newPath)
	} else {
		return newPath
	}
}


func RoutingQuery(startID int, endID int, mlp *MLP) ([]int, float64,int){
	// Get the shortest path from startID to endID
	//path, distance := VanillaDijkstra(mlp, startID, endID)

	// defer func() {
    //     if r := recover(); r != nil {
    //         fmt.Println("No route found")
    //     }
    // }()
	
    bipath, bidistance, affected := BidirectionalDijkstra(mlp, startID, endID)
	//bipath, bidistance := VanillaDijkstra(mlp, startID, endID)


    // // Print the output from the two algorithms
    // fmt.Printf("Vanilla Dijkstra: Shortest path from %d to %d: %v, Distance: %f\n", startID, endID, path, distance)
    // fmt.Printf("Bidirectional Dijkstra: Shortest path from %d to %d: %v, Distance: %f\n", startID, endID, bipath, bidistance)

    // // Check if the two algorithms return the same path and distance
    // if distance != bidistance || len(path) != len(bipath) {
    //     fmt.Println("The two algorithms return different results!")
    // } else {
    //     fmt.Println("The two algorithms return the same results!")
    // }

    return bipath, bidistance, affected
}

