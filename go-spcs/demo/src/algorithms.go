package src

import (
	"container/heap"
	"math"
	"sync"
)

func AllPairsShortestPath(g *Graph) {

	nodeIDs := []int{}
	for id := range g.GetNodes() {
		nodeIDs = append(nodeIDs, id)
	}

	dist := make(map[int]map[int]int)
	pre := make(map[int]map[int]*Node)

	// Initialize distances and predecessors
	for _, id := range nodeIDs {
		dist[id] = make(map[int]int)
		pre[id] = make(map[int]*Node)
		for _, otherID := range nodeIDs {
			if id == otherID {
				dist[id][otherID] = 0 // Distance to self is 0
			} else {
				dist[id][otherID] = math.MaxInt32 // Set to infinity
			}
			pre[id][otherID] = nil // Initialize predecessors to nil (no predecessor)
		}
	}

	// Set initial distances based on direct edges
	for _, id := range nodeIDs {
		node := g.nodes[id]
		for _, edge := range node.outgoingEdges {
			dist[edge.src][edge.dest] = edge.weight
			pre[edge.src][edge.dest] = g.nodes[edge.src]
		}
	}

	// Floyd-Warshall algorithm
	for _, k := range nodeIDs {
		for _, i := range nodeIDs {
			for _, j := range nodeIDs {
				if dist[i][k] == math.MaxInt32 || dist[k][j] == math.MaxInt32 {
					continue
				}

				if dist[i][k]+dist[k][j] < dist[i][j] {
					dist[i][j] = dist[i][k] + dist[k][j]
					pre[i][j] = pre[k][j]
				}
			}
		}
	}

	g.SetPre(pre)
}

func AllPairsShortestPathNE(nodes map[int]*Node, edges map[int]*Edge, partition *Partition) {
	nodeIDs := make([]int, 0, len(nodes))
	for id := range nodes {
		nodeIDs = append(nodeIDs, id)
	}

	dist := make(map[int]map[int]int)
	pre := make(map[int]map[int]*Node)
	var wg sync.WaitGroup
	var mutex sync.Mutex

	for _, sourceID := range nodeIDs {
		wg.Add(1)
		go func(sourceID int) {
			defer wg.Done()
			distFromSource, preFromSource := Dijkstra(nodes, edges, sourceID)
			mutex.Lock()
			defer mutex.Unlock()
			dist[sourceID] = distFromSource
			pre[sourceID] = preFromSource
		}(sourceID)
	}

	wg.Wait()

	partition.SetAPSP(dist, pre)
}

func Dijkstra(nodes map[int]*Node, edges map[int]*Edge, sourceID int) (map[int]int, map[int]*Node) {
	dist := make(map[int]int)
	pre := make(map[int]*Node)
	for id := range nodes {
		dist[id] = math.MaxInt32
		pre[id] = nil
	}
	dist[sourceID] = 0

	priorityQueue := make(PriorityQueue, 0)
	heap.Init(&priorityQueue)
	heap.Push(&priorityQueue, &Item{nodeID: sourceID, priority: 0})

	for priorityQueue.Len() > 0 {
		item := heap.Pop(&priorityQueue).(*Item)
		u := item.nodeID

		for _, edge := range edges{
			if edge.src != u {
				continue
			}
			v := edge.dest
			weight := edge.weight
			alt := dist[u] + weight
			if alt < dist[v] {
				dist[v] = alt
				pre[v] = nodes[u]
				heap.Push(&priorityQueue, &Item{nodeID: v, priority: alt})
			}
		}
	}

	return dist, pre
}

type Item struct {
	nodeID   int
	priority int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Item)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}


func GetShortestPath(g *Graph, startID, endID int) []int {
	pre := g.GetPre()
	path := make([]int, 0)
	at, ok := pre[startID][endID]

	for ok && at != nil && at.GetId() != startID {
		path = append([]int{at.GetId()}, path...)
		at, ok = pre[startID][at.GetId()]
	}

	if at != nil {
		path = append([]int{startID}, path...)
		path = append(path, endID)
	}
	return path
}

func GetShortestPathPre(pre map[int]map[int]*Node, startID, endID int) []int {
	path := make([]int, 0)
	at, ok := pre[startID][endID]

	for ok && at != nil && at.GetId() != startID {
		path = append([]int{at.GetId()}, path...)
		at, ok = pre[startID][at.GetId()]
	}

	if at != nil {
		path = append([]int{startID}, path...)
		path = append(path, endID)
	}
	return path
}