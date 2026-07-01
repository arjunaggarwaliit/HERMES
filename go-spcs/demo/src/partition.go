package src

import (
	"math"
	"sync"
)

var lastPartitionID int = -1

type Partition struct {
	id       int
	nodes    map[int]*Node
	level    int
	parent   *Partition
	children []*Partition

	Border_Nodes map[int]*Node
	Border_Edges map[int]*Edge

	Shortcut_Nodes map[int]*Node
	Shortcut_Edges map[int]*Edge

	EdgesMap  map[int]*Edge
	EdgeNodeMap map[int]map[int]*Edge

	Apsp_Dist map[int]map[int]int
	Apsp_Pre  map[int]map[int]*Node

	Affected bool
	Affected_count int

	routingLock    sync.RWMutex 	// For routing queries (allows multiple read locks)
	exclusiveLock     sync.Mutex   	// For update queries (allows only one write lock)
	intentLock     sync.Mutex   	// For intent locking (controls access to update lock)
	intentLockAcquired bool         // Indicates if the intent lock is currently acquired
	exclusiveLockAcquired  bool     // Indicates if the write lock is currently acquired
}

// RLock acquires a read lock
func (p *Partition) RLock() bool {
	p.exclusiveLock.Lock()  	
	p.intentLock.Lock() 	
	result := false
	if !p.intentLockAcquired && !p.exclusiveLockAcquired {
		result = p.routingLock.TryRLock() // Acquire read lock after acquiring write lock
	} 
	p.intentLock.Unlock() 
	p.exclusiveLock.Unlock()
	return result
}

// RUnlock releases a read lock
func (p *Partition) RUnlock() {
	p.routingLock.RUnlock() // Release read lock
}

// IntentLock acquires the intent lock
func (p *Partition) IntentLock() bool {
	p.exclusiveLock.Lock()
	p.intentLock.Lock()
	result := false
	if !p.intentLockAcquired && !p.exclusiveLockAcquired {
		p.intentLockAcquired = true
		result = true
	} 
	p.intentLock.Unlock()
	p.exclusiveLock.Unlock()
	return result
}

// IntentUnlock releases the intent lock
func (p *Partition) IntentUnlock() {
	p.intentLock.Lock()
	p.intentLockAcquired = false
	p.intentLock.Unlock()
}

// Lock acquires a write lock
func (p *Partition) Lock() bool {
	p.exclusiveLock.Lock()
	p.intentLock.Lock()
	result := false
	if !p.intentLockAcquired && !p.exclusiveLockAcquired {
		result = p.routingLock.TryLock()
	}
	p.intentLock.Unlock()
	p.exclusiveLock.Unlock()
	return result
}

// Unlock releases the write lock
func (p *Partition) Unlock() {
	p.exclusiveLock.Lock()
	p.exclusiveLockAcquired = false
	p.routingLock.Unlock()
	p.exclusiveLock.Unlock()
}

type NodeDijkstra struct {
	ID       int
	EdgeID   int
	Distance float64
}

func NewPartition() *Partition {
	lastPartitionID++
	return &Partition{
		id:    lastPartitionID,
		nodes: make(map[int]*Node),
		parent: nil,
		EdgesMap: make(map[int]*Edge),
		EdgeNodeMap: make(map[int]map[int]*Edge),
	}
}

func (p *Partition) CreateNode(id int) error {
	if _, ok := p.nodes[id]; ok {
		return nil
	}
	p.nodes[id] = &Node{
		id: id,
	}
	return nil
}

func (p *Partition) AddNode(node *Node) {
	p.nodes[node.id] = node
}

func (p *Partition) DeleteNode(node *Node) {
	delete(p.nodes, node.id)
}

func (p *Partition) GetId() int {
	return p.id
}

// Return an map of IDs to nodes
func (p *Partition) GetNodes() map[int]*Node {
	return p.nodes
}

// Return an array of nodes
func (p *Partition) GetNodesPointer() []*Node {
	nodes := make([]*Node, 0, len(p.nodes))
	for _, node := range p.nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

func (p *Partition) GetLevel() int {
	return p.level
}

func (p *Partition) SetLevel(level int) {
	p.level = level
}

func (p *Partition) GetParent() *Partition {
	return p.parent
}

func (p *Partition) GetNeighbors(node *Node) []*NodeDijkstra {
	// Use Shortcut_Edges to get the neighbors
	neighbors := make([]*NodeDijkstra, 0)

	for _, e := range p.Shortcut_Edges {
		if e.src == node.id {
			NodeDijkstra := &NodeDijkstra{ ID: e.dest, Distance: float64(e.weight), EdgeID: e.id}
			neighbors = append(neighbors, NodeDijkstra)
		} else if e.dest == node.id {
			NodeDijkstra := &NodeDijkstra{ ID: e.src, Distance: float64(e.weight), EdgeID: e.id}
			neighbors = append(neighbors, NodeDijkstra)
		}
	}

	for _, e := range p.Border_Edges {
		if e.src == node.id {
			NodeDijkstra := &NodeDijkstra{ ID: e.dest, Distance: float64(e.weight), EdgeID: e.id}
			neighbors = append(neighbors, NodeDijkstra)
		} else if e.dest == node.id {
			NodeDijkstra := &NodeDijkstra{ ID: e.src, Distance: float64(e.weight), EdgeID: e.id}
			neighbors = append(neighbors, NodeDijkstra)
		}
	}

	// Make a unique list of neighbors. Make sure there is no repeated node or self node
	uniqueNeighbors := make([]*NodeDijkstra, 0)
	uniqueMap := make(map[int]bool)
	for _, n := range neighbors {
		if n.ID != node.id && uniqueMap[n.ID] == false {
			uniqueNeighbors = append(uniqueNeighbors, n)
			uniqueMap[n.ID] = true
		}
	}

	return uniqueNeighbors
}

func (p *Partition) SetParent(parent *Partition) {
	p.parent = parent
	parentLevel := parent.GetLevel()
	p.SetLevel(parentLevel + 1)
}

func (p *Partition) SetAPSP(dist map[int]map[int]int, pre map[int]map[int]*Node) {
	p.Apsp_Dist = dist
	p.Apsp_Pre = pre
}

func (p *Partition) GetChildren() []*Partition {
	return p.children
}

func (p *Partition) AddChild(child *Partition) {
	p.children = append(p.children, child)
}

func (p *Partition) contains(node *Node) bool {
	_, ok := p.nodes[node.id]
	return ok
}

func (p *Partition) size() int {
	return len(p.nodes)
}


// For a partition, shortcut network is constructed in the following manner:
// 1. MakeShortcutNodes: Add the border nodes of the child partitions to the shortcut nodes
// 2. MakeShortcutEdges:
// 						- Use Apsp_Dist from the child partition to add the shortcut edges
// 						- Use the border edges of the child partition to classify the edges as border or shortcut edges at the parent partition
//                      - This also sets the border edges at the current partition
// 3. SetBorderNodes: From the shortcut nodes, identify the border nodes at the current partition

// Get the shortcut nodes
func (p *Partition) MakeShortcutNodes() {
	// Shortcut nodes are the border nodes of the child partitions
	if p.Shortcut_Nodes == nil {
		p.Shortcut_Nodes = make(map[int]*Node)
		for _, child := range p.children {
			//fmt.Println("Child Border Nodes: ", child.Border_Nodes)
			for _, n := range child.Border_Nodes {
				p.Shortcut_Nodes[n.id] = n
			}
		}
	}
}

// Get the shortcut edges
func (p *Partition) MakeShortcutEdges() {
	p.Border_Edges = make(map[int]*Edge)
	p.Shortcut_Edges = make(map[int]*Edge)
	

	// Use Apsp_Dist from the child partition to add the shortcut edges
	for _, child := range p.children {
		for _, n := range child.Border_Nodes {
			for _, m := range child.Border_Nodes {
				if child.Apsp_Dist[n.id][m.id] == math.MaxInt32 {
					continue
				}

				new_edge := NewEdge(n.id, m.id, child.Apsp_Dist[n.id][m.id])
				p.Shortcut_Edges[new_edge.id] = new_edge
				p.EdgesMap[new_edge.id] = new_edge

				if p.EdgeNodeMap[n.id] == nil {
					p.EdgeNodeMap[n.id] = make(map[int]*Edge)
				}
				p.EdgeNodeMap[n.id][m.id] = new_edge
			}
		}
	}

	combined_child_Border_Edges := make(map[int]*Edge)
	for _, child := range p.children {
		for _, e := range child.Border_Edges {
			combined_child_Border_Edges[e.id] = e
		}
	}

	for _, e := range combined_child_Border_Edges {
		// Check if both the endpoints are in shortcut nodes, add the edge to the shortcut edges
		// Else add the edge to the border edges
		if _, ok := p.Shortcut_Nodes[e.src]; ok {
			if _, ok := p.Shortcut_Nodes[e.dest]; ok {
				p.Shortcut_Edges[e.id] = e
			} else {
				p.Border_Edges[e.id] = e
			}
		} else {
			p.Border_Edges[e.id] = e
		}
		p.EdgesMap[e.id] = e

		if p.EdgeNodeMap[e.src] == nil {
			p.EdgeNodeMap[e.src] = make(map[int]*Edge)
		}
		p.EdgeNodeMap[e.src][e.dest] = e
	}
}

// Using the border edges, identify the border nodes
func (p *Partition) SetBorderNodes() {
	if p.Border_Nodes == nil {
		p.Border_Nodes = make(map[int]*Node)
		// Use the border edges to identify the border nodes from the shortcut nodes
		for _, e := range p.Border_Edges {
			if p.Shortcut_Nodes[e.src] != nil {
				p.Border_Nodes[e.src] = p.Shortcut_Nodes[e.src]
			}
			if p.Shortcut_Nodes[e.dest] != nil {
				p.Border_Nodes[e.dest] = p.Shortcut_Nodes[e.dest]
			}
		}
	}
}

func (p *Partition) ReleaseLocks() {
	for _, e:= range p.EdgesMap {
		e.Unlock()
	}
}

// Construct the shortcut network at intermediate levels
func (p *Partition) ConstructShortcutNetwork() {
	p.MakeShortcutNodes()
	p.MakeShortcutEdges()
	p.SetBorderNodes()

	//fmt.Println("Number of Shortcut Nodes: ", len(p.Shortcut_Nodes), " Partition ID: ", p.GetId())

	AllPairsShortestPathNE(p.Shortcut_Nodes, p.Shortcut_Edges, p)
}

func (p *Partition) RebuildShortcutNetwork(){
	temp_dist := make(map[int]map[int]int)

	for _, child := range p.children {
		for _, n := range child.Border_Nodes {
			for _, m := range child.Border_Nodes {
				if child.Apsp_Dist[n.id][m.id] == math.MaxInt32 {
					continue
				}

				if temp_dist[n.id] == nil {
					temp_dist[n.id] = make(map[int]int)
				}

				if child.Apsp_Dist[n.id][m.id]!=p.Apsp_Dist[n.id][m.id]{
					temp_dist[n.id][m.id] = child.Apsp_Dist[n.id][m.id]
				}
			}
		}
	}

	// Make a temp map to the Shortcut Edges
	temp_shortcut_edges := make(map[int]map[int]*Edge)
	for e := range p.Shortcut_Edges {
		if temp_shortcut_edges[p.Shortcut_Edges[e].src] == nil {
			temp_shortcut_edges[p.Shortcut_Edges[e].src] = make(map[int]*Edge)
		}
		temp_shortcut_edges[p.Shortcut_Edges[e].src][p.Shortcut_Edges[e].dest] = p.Shortcut_Edges[e]
	}

	for src, dist := range temp_dist {
		for dest, val := range dist {
			temp_shortcut_edges[src][dest].SetWeight(val)
			//print("Src: ", src, " Dest: ", dest, " Val: ", val, " Updated: ", temp_shortcut_edges[src][dest].GetWeight(), "\n")
		}
	}

}

func (p *Partition) RecomputeShortcutNetwork(){
	// Check Shortcut Edges and their weight
	//print("Recomputing Shortcut Network for Partition: ", p.GetId(), "\n")
	AllPairsShortestPathNE(p.Shortcut_Nodes, p.Shortcut_Edges, p)
}

// Construct the shortcut network at the base level
func (p *Partition) ConstructShortcutNetworkBase() {
	p.Shortcut_Nodes = make(map[int]*Node)
	p.Shortcut_Edges = make(map[int]*Edge)
	p.Border_Edges = make(map[int]*Edge)

	for _, n := range p.nodes {
		p.Shortcut_Nodes[n.id] = n
	}

	for _, n := range p.Shortcut_Nodes {
		for _, edge := range n.edges {
			if p.Shortcut_Nodes[edge.GetOtherNode(n).id] != nil {
				p.Shortcut_Edges[edge.id] = edge
			} else {
				p.Border_Edges[edge.id] = edge
			}
			p.EdgesMap[edge.id] = edge

			if p.EdgeNodeMap[edge.src] == nil {
				p.EdgeNodeMap[edge.src] = make(map[int]*Edge)
			}
			p.EdgeNodeMap[edge.src][edge.dest] = edge
		}
	}

	p.SetBorderNodes()

	//fmt.Print("Number of Shortcut Nodes: ", len(p.Shortcut_Nodes), " Partition ID: ", p.GetId(), "\n")

	AllPairsShortestPathNE(p.Shortcut_Nodes, p.Shortcut_Edges, p)
}

func (p *Partition) GetPath(startID int, endID int) []int {
	// Get the shortest path from the Apsp_Pre
	path := make([]int, 0)
	at, ok := p.Apsp_Pre[startID][endID]

	for ok && at != nil && at.GetId() != startID {
		path = append([]int{at.GetId()}, path...)
		at, ok = p.Apsp_Pre[startID][at.GetId()]
	}

	if at != nil {
		path = append([]int{startID}, path...)
		path = append(path, endID)
	}
	return path
}

// Return the edges and for edges check if both the endpoints are in the partition
// func (p *Partition) GetPartitionEdges() []*Edge {
// 	edges := []*Edge{}
// 	for _, n := range p.nodes {
// 		for _, e := range n.edges {
// 			if p.contains(e.GetOtherNode(n)) {
// 				edges = append(edges, e)
// 			}
// 		}
// 	}
// 	return edges
// }

// func (p *Partition) GetAPSS() [][]int {

// 	p.apss = allPairsShortestPath(p.nodes)
// 	return p.apss
// }
