package src

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
)

var lastEdgeID int = -1

type Edge struct {
	id                int
	src, dest, weight int

	routingLock    sync.RWMutex 	// For routing queries (allows multiple read locks)
	exclusiveLock     sync.Mutex   	// For update queries (allows only one write lock)
	intentLock     sync.Mutex   	// For intent locking (controls access to update lock)
	intentLockAcquired bool         // Indicates if the intent lock is currently acquired
	exclusiveLockAcquired  bool     // Indicates if the write lock is currently acquired
}

type Node struct {
	id                           int
	incomingEdges, outgoingEdges []*Edge
	edges                        []*Edge
	basePartition                *Partition

	NodeEdgeMap map[int]*Edge
}

type Graph struct {
	nodes map[int]*Node
	pre   map[int]map[int]*Node
}

func NewEdge(src, dest, weight int) *Edge {
	lastEdgeID++
	return &Edge{
		id:     lastEdgeID,
		src:    src,
		dest:   dest,
		weight: weight,
	}
}

// RLock acquires a read lock
func (e *Edge) RLock() bool {
	e.exclusiveLock.Lock()  	
	e.intentLock.Lock() 	
	result := false
	if !e.intentLockAcquired && !e.exclusiveLockAcquired {
		result = e.routingLock.TryRLock() // Acquire read lock after acquiring write lock
	} 
	e.intentLock.Unlock() 
	e.exclusiveLock.Unlock()
	return result
}

// RUnlock releases a read lock
func (e *Edge) RUnlock() {
	e.routingLock.RUnlock() // Release read lock
}

// IntentLock acquires the intent lock
func (e *Edge) IntentLock() bool {
	e.exclusiveLock.Lock()
	e.intentLock.Lock()
	result := false
	if !e.intentLockAcquired && !e.exclusiveLockAcquired {
		e.intentLockAcquired = true
		result = true
	} 
	e.intentLock.Unlock()
	e.exclusiveLock.Unlock()
	return result
}

// IntentUnlock releases the intent lock
func (e *Edge) IntentUnlock() {
	e.intentLock.Lock()
	e.intentLockAcquired = false
	e.intentLock.Unlock()
}

// Lock acquires a write lock
func (e *Edge) Lock() bool {
	e.exclusiveLock.Lock()
	e.intentLock.Lock()
	result := false
	if !e.intentLockAcquired && !e.exclusiveLockAcquired {
		result = e.routingLock.TryLock()
	}
	if result{
		e.exclusiveLockAcquired = true
	}
	e.intentLock.Unlock()
	e.exclusiveLock.Unlock()
	return result
}

// Unlock releases the write lock
func (e *Edge) Unlock() {
	e.exclusiveLock.Lock()
	if ! e.exclusiveLockAcquired {
		e.exclusiveLock.Unlock()
		return
	}
	e.exclusiveLockAcquired = false
	e.routingLock.Unlock()
	e.exclusiveLock.Unlock()
}

func NewGraph() *Graph {
	return &Graph{
		nodes: make(map[int]*Node),
	}
}

func (g *Graph) GetPre() map[int]map[int]*Node {
	return g.pre
}

func (g *Graph) SetPre(pre map[int]map[int]*Node) {
	g.pre = pre
}

func (g *Graph) AddNode(id int) error {
	if _, ok := g.nodes[id]; ok {
		return nil
	}
	g.nodes[id] = &Node{
		id: id,
	}
	return nil
}

func (g *Graph) AddEdge(src, dest, weight int) error {
	if _, ok := g.nodes[src]; !ok {
		return fmt.Errorf("source node %d not found", src)
	}
	if _, ok := g.nodes[dest]; !ok {
		return fmt.Errorf("destination node %d not found", dest)
	}
	var edge *Edge = NewEdge(src, dest, weight)
	g.nodes[src].outgoingEdges = append(g.nodes[src].outgoingEdges, edge)
	g.nodes[dest].incomingEdges = append(g.nodes[dest].incomingEdges, edge)
	g.nodes[src].edges = append(g.nodes[src].edges, edge)
	g.nodes[dest].edges = append(g.nodes[dest].edges, edge)

	if g.nodes[src].NodeEdgeMap == nil {
		g.nodes[src].NodeEdgeMap = make(map[int]*Edge)
	}
	if g.nodes[dest].NodeEdgeMap == nil {
		g.nodes[dest].NodeEdgeMap = make(map[int]*Edge)
	}
	g.nodes[src].NodeEdgeMap[dest] = edge
	return nil
}

func (g *Graph) GetNodes() map[int]*Node {
	return g.nodes
}

func (g *Graph) GetEdges() map[int]*Edge {
	edges := make(map[int]*Edge)
	for _, node := range g.GetNodes() {
		for _, edge := range node.GetEdges() {
			edges[edge.id] = edge
		}
	}
	return edges
}

func (n *Node) GetId() int {
	return n.id
}

func (n *Node) GetEdges() []*Edge {
	return n.edges
}

func (n *Node) SetPartition(p *Partition) {
	n.basePartition = p
}

func (n *Node) GetPartition() *Partition {
	return n.basePartition
}

func (e *Edge) GetId() int {
	return e.id
}

func (e *Edge) GetSrc() int {
	return e.src
}

func (e *Edge) GetDest() int {
	return e.dest
}

func (e *Edge) GetWeight() int {
	return e.weight
}

func (e *Edge) SetWeight(weight int) {
	e.weight = weight
}

func (e *Edge) GetOtherNode(n *Node) *Node {
	if e.src == n.id {
		return &Node{id: e.dest}
	}
	return &Node{id: e.src}
}

func (g *Graph) CreateGraph(nodes []int, edges [][]int) {
	// Add nodes
	for _, nodeID := range nodes {
		g.AddNode(nodeID)
	}

	// Add edges
	for _, edge := range edges {
		if len(edge) == 3 {
			g.AddEdge(edge[0], edge[1], edge[2])
		}
	}
}

func (g *Graph) CreateSampleGraph(nodes int, edges int) {
	// Add nodes
	for i := 0; i < nodes; i++ {
		g.AddNode(i)
	}

	// A simple check to avoid infinite loops in case edges > nodes*(nodes-1) for a directed graph
	maxEdges := nodes * (nodes - 1)
	if edges > maxEdges {
		edges = maxEdges
	}

	// Add edges with a check to avoid multiple edges between the same two nodes
	addedEdges := make(map[[2]int]bool)
	for i := 0; i < edges; i++ {
		src := i % nodes
		dest := (i + 1) % nodes
		weight := i + 1 // Assign a weight, for example

		edgeKey := [2]int{src, dest}
		if !addedEdges[edgeKey] {
			g.AddEdge(src, dest, weight)
			addedEdges[edgeKey] = true
		} else {
			// If an edge between src and dest already exists, find a new dest to ensure uniqueness
			for j := 0; j < nodes; j++ {
				newDest := (dest + j + 1) % nodes
				edgeKey = [2]int{src, newDest}
				if !addedEdges[edgeKey] {
					g.AddEdge(src, newDest, weight)
					addedEdges[edgeKey] = true
					break // Break after adding a unique edge
				}
			}
		}
	}

	// // Print all edges
	// for _, node := range g.GetNodes() {
	// 	for _, edge := range node.GetEdges() {
	// 		fmt.Printf("Edge from %d to %d with weight %d\n", edge.src, edge.dest, edge.weight)
	// 	}
	// }
}

func (g *Graph) CreateGraphFromText(fname string) error {

	vertices := []int{}
	edges := [][]int{}

	file, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a buffered scanner to read the file
	scanner := bufio.NewScanner(file)

	// Create channels for communication
	lineCh := make(chan string)
	errCh := make(chan error)

	// Goroutine to read lines from the file and send them to the channel
	go func() {
		for scanner.Scan() {
			lineCh <- scanner.Text()
		}
		close(lineCh)
	}()

	// Goroutine to process lines concurrently
	var wg sync.WaitGroup
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for line := range lineCh {
				err := processLine(line, &vertices, &edges)
				if err != nil {
					errCh <- err
					return
				}
			}
		}()
	}

	// Wait for all processing goroutines to finish
	go func() {
		wg.Wait()
		close(errCh)
	}()

	// Check for any errors
	for err := range errCh {
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func processLine(line string, vertices *[]int, edges *[][]int) error {
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

	// Use a mutex to synchronize access to shared data
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()

	// Check if the source and target vertices are in the vertices array
	if !slices.Contains(*vertices, int(source)) {
		*vertices = append(*vertices, int(source))
	}

	if !slices.Contains(*vertices, int(target)) {
		*vertices = append(*vertices, int(target))
	}

	// Add the edge to the edges array
	*edges = append(*edges, []int{int(source), int(target), int(weight)})

	if oneway == "B" {
		// Add the reverse edge if it's bidirectional
		*edges = append(*edges, []int{int(target), int(source), int(weight)})
	}

	return nil
}

func (g *Graph) GenerateGraphDot(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintln(file, "digraph G {")
	if err != nil {
		return err
	}

	for _, node := range g.GetNodes() {
		for _, edge := range node.GetEdges() {
			_, err = fmt.Fprintf(file, "    %d -> %d [label=\"%d\"];\n", edge.src, edge.dest, edge.weight)
			if err != nil {
				return err
			}
		}
	}

	_, err = fmt.Fprintln(file, "}")

	// Run the subprocess command to generate the PNG image
	cmd := exec.Command("dot", "-Tpng", filePath, "-o", strings.Replace(filePath, ".dot", ".png", 1))
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error running dot command:", err)
	}
	fmt.Println("MLP DOT and PNG files generated successfully!")

	return err
}

func (m *MLP) GenerateMLPDOT(filePath string) {
	var b strings.Builder
	b.WriteString("digraph G {\n")

	// Level 0: All nodes and edges
	b.WriteString("\t// Level 0: All nodes and edges\n")

	// Subsequent Levels: Partitions with their nodes and edges
	for i := 1; i < m.GetLevelNum(); i++ {
		b.WriteString(fmt.Sprintf("\t// Level %d: Partitions\n", i))
		for _, p := range m.GetPartitions(i) {
			b.WriteString(fmt.Sprintf("\tsubgraph cluster_%d_%d {\n\t\tlabel=\"Level %d: Partition %d\";\n", i, p.GetId(), i, p.GetId()))
			nodes := p.GetNodes()
			nodeIDs := make([]int, 0, len(nodes))
			for id := range nodes {
				nodeIDs = append(nodeIDs, id)
				// Assuming you have a way to get edges for a node: node.GetEdges()
				for _, edge := range nodes[id].GetEdges() {
					if _, exists := nodes[edge.dest]; exists {
						// Internal edge
						b.WriteString(fmt.Sprintf("\t\t%d -> %d [label=\"%d\"];\n", edge.src, edge.dest, edge.weight))
					}
				}
			}
			b.WriteString("\t}\n")
		}
	}

	b.WriteString("}")

	// Create the directory path first
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Println("Error creating directory path:", err)
		return
	}

	// Write the DOT string to the specified file path
	err := os.WriteFile(filePath, []byte(b.String()), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	// Run the subprocess command to generate the PNG image
	cmd := exec.Command("dot", "-Tpng", filePath, "-o", strings.Replace(filePath, ".dot", ".png", 1))
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error running dot command:", err)
	}
	fmt.Println("MLP DOT and PNG files generated successfully!")

}

func main() {
	// Create a new graph
	g := NewGraph()

	// Add vertices
	g.AddNode(0)
	g.AddNode(1)
	g.AddNode(2)

	// Add edges
	g.AddEdge(0, 1, 5)
	g.AddEdge(0, 2, 10)

	// Get nodes
	nodes := g.GetNodes()

	// Print nodes
	fmt.Println("Nodes:")
	for _, n := range nodes {
		fmt.Println(n.GetId())
	}
}
