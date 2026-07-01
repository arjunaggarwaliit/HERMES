package src

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// MLP represents a Multi-Level Partitioning structure.
type MLP struct {
	levelNum    int
	level       map[int][]*Partition
	nodeMap     map[int]*Node
	edgeMap     map[int]*Edge
	AcquiredIDs []int
}

type EdgeSimplified struct {
	From   int `json:"from"`
	To     int `json:"to"`
	Weight int `json:"weight"`
}

type PartitionSimplified struct {
	Id              int                   `json:"id"`
	NodeIds         []int                 `json:"nodeIds"`
	Level           int                   `json:"level"`
	ParentID        int                   `json:"parentId"`
	Children        []int                 `json:"children"`
	BorderNodeIds   []int                 `json:"borderNodeIds"`
	BorderEdgeIds   []EdgeSimplified      `json:"borderEdgeIds"`
	ShortcutNodeIds []int                 `json:"shortcutNodeIds"`
	ShortcutEdgeIds []EdgeSimplified      `json:"shortcutEdgeIds"`
	APSPDist        map[int]map[int]int   `json:"apspDist"`
	APSPPre         map[int]map[int]*Node `json:"apspPre"`
}

type MLPSimplified struct {
	LevelNum    int                           `json:"levelNum"`
	Levels      map[int][]PartitionSimplified `json:"levels"`
	NodeMap     map[int]*Node                 `json:"nodeMap"`
	EdgeMap     map[int]*Edge                 `json:"edgeMap"`
	AcquiredIDs []int                         `json:"acquiredIDs"`
}

// NewMLP creates a new instance of MLP with the specified number of levels.
func NewMLP(levelNum int) *MLP {
	mlp := &MLP{
		levelNum: levelNum,
		level:    make(map[int][]*Partition),
		nodeMap:  make(map[int]*Node),
		edgeMap:  make(map[int]*Edge),
	}
	// Initialize partitions slice at each level
	for i := 0; i < levelNum; i++ {
		mlp.level[i] = []*Partition{}
	}
	return mlp
}

// Get Node by ID
func (mlp *MLP) GetNodeByID(id int) *Node {
	return mlp.nodeMap[id]
}

// Get Edge by ID
func (mlp *MLP) GetEdgeByID(id int) *Edge {
	return mlp.edgeMap[id]
}

// AddPartition adds a partition to the specified level.
func (mlp *MLP) AddPartition(level int, partition *Partition) {
	if level >= mlp.levelNum{
		fmt.Println("Level number exceeds the maximum level number")
		fmt.Println("Partition ID: ", partition.GetId())
	}
	partition.SetLevel(level)
	mlp.level[level] = append(mlp.level[level], partition)
}

// GetPartitions returns the partitions at the specified level.
func (mlp *MLP) GetPartitions(level int) []*Partition {
	return mlp.level[level]
}

// GetLevelNum returns the number of levels in the MLP.
func (mlp *MLP) GetLevelNum() int {
	return mlp.levelNum
}

// GetPartitionCountAtLevel returns the number of partitions at the specified level.
func (mlp *MLP) GetPartitionCountAtLevel(level int) int {
	return len(mlp.level[level])
}

// ModifyPartition modifies the partition at the specified level and index.
func (mlp *MLP) ModifyPartition(level, index int, partition *Partition) {
	mlp.level[level][index] = partition
}

func SerializeMLP(mlp *MLP) (string, error) {

	simplifiedMLP := MLPSimplified{
		LevelNum:    mlp.levelNum,
		Levels:      make(map[int][]PartitionSimplified),
		NodeMap:     make(map[int]*Node),
		EdgeMap:     make(map[int]*Edge),
		AcquiredIDs: mlp.AcquiredIDs,
	}

	simplifyAndSerializePartition := func(p *Partition) (PartitionSimplified, error) {
		simplified := PartitionSimplified{
			Id:              p.id,
			NodeIds:         make([]int, 0, len(p.nodes)),
			Level:           p.level,
			ParentID:        -1,
			Children:        make([]int, len(p.children)),
			BorderNodeIds:   make([]int, 0, len(p.Border_Nodes)),
			BorderEdgeIds:   make([]EdgeSimplified, 0, len(p.Border_Edges)),
			ShortcutNodeIds: make([]int, 0, len(p.Shortcut_Nodes)),
			ShortcutEdgeIds: make([]EdgeSimplified, 0, len(p.Shortcut_Edges)),
			APSPDist:        p.Apsp_Dist,
			APSPPre:         p.Apsp_Pre,
		}

		for id := range p.nodes {
			simplified.NodeIds = append(simplified.NodeIds, id)
		}

		for id := range p.Border_Nodes {
			simplified.BorderNodeIds = append(simplified.BorderNodeIds, id)
		}

		for _, edge := range p.Border_Edges {
			simplified.BorderEdgeIds = append(simplified.BorderEdgeIds, EdgeSimplified{From: edge.GetSrc(), To: edge.GetDest(), Weight: edge.weight})
		}

		for id := range p.Shortcut_Nodes {
			simplified.ShortcutNodeIds = append(simplified.ShortcutNodeIds, id)
		}

		for _, edge := range p.Shortcut_Edges {
			simplified.ShortcutEdgeIds = append(simplified.ShortcutEdgeIds, EdgeSimplified{From: edge.GetSrc(), To: edge.GetDest(), Weight: edge.weight})
		}

		for from, distMap := range p.Apsp_Dist {
			simplified.APSPDist[from] = make(map[int]int)
			for to, dist := range distMap {
				simplified.APSPDist[from][to] = dist
			}
		}

		for from, preMap := range p.Apsp_Pre {
			simplified.APSPPre[from] = make(map[int]*Node)
			for to, preNode := range preMap {
				simplified.APSPPre[from][to] = preNode
			}
		}

		if p.parent != nil {
			simplified.ParentID = p.parent.id
		}

		for i, child := range p.children {
			simplified.Children[i] = child.id
		}

		return simplified, nil
	}

	for level, partitions := range mlp.level {
		for _, partition := range partitions {
			simplifiedPartition, err := simplifyAndSerializePartition(partition)
			if err != nil {
				return "", err
			}
			simplifiedMLP.Levels[level] = append(simplifiedMLP.Levels[level], simplifiedPartition)
		}
	}

	data, err := json.Marshal(simplifiedMLP)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func SaveMLPToJsonFile(mlp *MLP, filename string) error {
	// Serialize the MLP to a JSON string
	serializedMLP, err := SerializeMLP(mlp)
	if err != nil {
		return err // Return serialization error
	}

	// Convert the JSON string to a byte slice for writing
	data := []byte(serializedMLP)

	// Write the serialized data to a file
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err // Return file writing error
	}

	return nil // Successfully saved MLP to file
}

func LoadMLPFromJSONFile(filename string) (*MLP, error) {
	// Define the simplified structures for deserialization with additional fields

	// Read the JSON file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var simplifiedMLP MLPSimplified

	// Deserialize the JSON data into the simplified MLP structure
	if err := json.Unmarshal(data, &simplifiedMLP); err != nil {
		return nil, err
	}

	// Begin reconstructing the original MLP structure from the simplified version
	mlp := &MLP{
		levelNum: simplifiedMLP.LevelNum,
		level:    make(map[int][]*Partition),
	}

	partitionMap := make(map[int]*Partition)

	for level, simplifiedPartitions := range simplifiedMLP.Levels {
		for _, sp := range simplifiedPartitions {
			partition := &Partition{
				id:       sp.Id,
				nodes:    make(map[int]*Node), // Reconstruct nodes based on NodeIds
				level:    sp.Level,
				parent:   nil, // To be linked later
				children: make([]*Partition, 0),
				// Initialize maps for Border_Nodes, Shortcut_Nodes, etc.
				Border_Nodes:   make(map[int]*Node),
				Border_Edges:   make(map[int]*Edge),
				Shortcut_Nodes: make(map[int]*Node),
				Shortcut_Edges: make(map[int]*Edge),
				Apsp_Dist:      sp.APSPDist,
				Apsp_Pre:       sp.APSPPre,
			}

			for _, nodeId := range sp.NodeIds {
				partition.CreateNode(nodeId)
			}

			partition.MakeShortcutNodes()
			partition.MakeShortcutEdges()
			partition.SetBorderNodes()

			partition.Apsp_Dist = sp.APSPDist
			partition.Apsp_Pre = sp.APSPPre

			partitionMap[sp.Id] = partition
			mlp.level[level] = append(mlp.level[level], partition)
		}
	}

	// Link parents and children
	for _, simplifiedPartitions := range simplifiedMLP.Levels {
		for _, sp := range simplifiedPartitions {
			partition := partitionMap[sp.Id]
			if sp.ParentID != -1 {
				partition.parent = partitionMap[sp.ParentID]
			}
			for _, childId := range sp.Children {
				childPartition := partitionMap[childId]
				partition.children = append(partition.children, childPartition)
			}
		}
	}

	return mlp, nil
}

func (mlp *MLP) MLPConstruction(graph *Graph) {
	// Construct the MLP
	// Have to think about the overhead of the MLP

	// Multi level partitioning
	nodes := graph.GetNodes()

	// Make the nodes map
	mlp.nodeMap = make(map[int]*Node)
	for _, node := range nodes {
		mlp.nodeMap[node.GetId()] = node
	}

	partition := NewPartition()
	for _, node := range nodes {
		partition.AddNode(node)
	}
	partition.SetLevel(0)
	mlp.AddPartition(0, partition)

	// Create the rest of the levels of the MLP
	type partitionResult struct {
		p1, p2 *Partition
	}

	for i := 1; i < mlp.levelNum; i++ {
		partitions := mlp.GetPartitions(i - 1)

		var wg sync.WaitGroup
		wg.Add(len(partitions))

		partitionCh := make(chan partitionResult, len(partitions))

		for _, p := range partitions {
			go func(p *Partition) {
				defer wg.Done()
				p1, p2 := Partition_algo(p)
				partitionCh <- partitionResult{p1, p2}
			}(p)
		}

		// Wait for all partitions at this level to complete
		wg.Wait()
		close(partitionCh)

		// Collect results and add partitions to the MLP
		for result := range partitionCh {
			mlp.AddPartition(i, result.p1)
			mlp.AddPartition(i, result.p2)
		}
		fmt.Printf("Partitioning done for level %d\n", i)
	}

	// At the base level, add the Partition to the node properties
	for _, p := range mlp.GetPartitions(mlp.levelNum - 1) {
		nodes := p.GetNodes()
		for _, node := range nodes {
			node.SetPartition(p)
		}
	}

	// Multi-level Shortcut Construction
	// Number of partitions at the last level
	numPartitions := mlp.GetPartitionCountAtLevel(mlp.levelNum - 1)

	var wg sync.WaitGroup
	wg.Add(numPartitions)

	for i := 0; i < numPartitions; i++ {
		// Construct the shortcuts for the last level
		go func(i int) {
			defer wg.Done()
			p := mlp.GetPartitionAtLevel(mlp.levelNum-1, i)
			p.ConstructShortcutNetworkBase()
		}(i)
	}

	wg.Wait()
	fmt.Printf("Shortcut Construction done at the base level\n")

	// Construct the shortcuts for the rest of the levels
	for i := mlp.levelNum - 2; i >= 0; i-- {
		partitionCount := mlp.GetPartitionCountAtLevel(i)
		var wg sync.WaitGroup
		wg.Add(partitionCount)
		for j := 0; j < partitionCount; j++ {
			go func(i, j int) {
				defer wg.Done()
				p := mlp.GetPartitionAtLevel(i, j)
				p.ConstructShortcutNetwork()
			}(i, j)
		}
		wg.Wait()
		fmt.Printf("Shortcut Construction done at level %d\n", i)
	}

	// Set the edge map for the MLP
	makeEdgeMap(mlp)
}

func makeEdgeMap(mlp *MLP) {
	mlp.edgeMap = make(map[int]*Edge)
	for i := mlp.levelNum - 1; i >= 0; i-- {
		for _, p := range mlp.GetPartitions(i) {
			// Add all the data from partiton edgeMap to the mlp edgeMap
			for id, edge := range p.EdgesMap {
				mlp.edgeMap[id] = edge
			}
		}
	}
}

func (mlp *MLP) GetPartitionAtLevel(level, index int) *Partition {
	return mlp.level[level][index]
}

func (mlp *MLP) GetNeighbors(nodeId int, level int) []*NodeDijkstra {
	// Get the node
	node := mlp.nodeMap[nodeId]

	// Get the partition
	partition := node.GetPartition()

	// Get to the desired level
	for i := 0; i < level; i++ {
		partition = partition.GetParent()
	}

	// Get the neighbors
	neighbors := partition.GetNeighbors(node)

	// Return the neighbors
	return neighbors
}

func (mlp *MLP) UnpackEdge(startID int, endID int, level int) [][]int {
	startNode := mlp.nodeMap[startID]
	endNode := mlp.nodeMap[endID]

	startPartition := startNode.GetPartition()
	endPartition := endNode.GetPartition()

	for i := 0; i < level-1; i++ {
		startPartition = startPartition.GetParent()
		endPartition = endPartition.GetParent()
	}

	if startPartition != endPartition {
		return [][]int{{startID, endID, 0}}
	} else {
		path := startPartition.GetPath(startID, endID)
		unpackedPath := make([][]int, 0)
		for i := 0; i < len(path)-1; i++ {
			unpackedPath = append(unpackedPath, []int{path[i], path[i+1], level - 1})
		}
		return unpackedPath
	}
}

func (mlp *MLP) LevelDijkstra(srcId int, destId int, nodeId int) int {
	// Get the node
	node := mlp.nodeMap[nodeId]
	srcNode := mlp.nodeMap[srcId]
	destNode := mlp.nodeMap[destId]

	levelU := unCommonLevel(srcNode, node)
	levelV := unCommonLevel(destNode, node)

	// Return minimum of the two levels
	if levelU < levelV {
		return levelU
	}
	return levelV
}

func unCommonLevel(srcNode *Node, node *Node) int {
	//fmt.Println("Finding unCommonLevel for ", srcNode.GetId(), " and ", node.GetId())
	partitionU := srcNode.GetPartition()
	partitionV := node.GetPartition()

	unCommonLevel := 0

	// Find the unCommonLevel
	for partitionU.id != partitionV.id {
		unCommonLevel++
		//fmt.Print("Partition ", partitionU.GetId(), " and ", partitionV.GetId(), "\n")
		partitionU = partitionU.GetParent()
		partitionV = partitionV.GetParent()
	}

	return unCommonLevel
}

// AcquireLocks acquires locks on partitions involved in routing query
func (mlp *MLP) AcquireLocks(startID int, endID int) {
	startPartition := mlp.GetNodeByID(startID).GetPartition()
	endPartition := mlp.GetNodeByID(endID).GetPartition()

	// Acquire locks on partitions along the traversal path
	for startPartition != endPartition {
		if startPartition != nil {
			startPartition.Lock()
			mlp.AcquiredIDs = append(mlp.AcquiredIDs, startPartition.GetId())
		}
		startPartition = startPartition.GetParent()
	}

	// Print accquired IDs
	fmt.Println("Acquired IDs: ", mlp.AcquiredIDs)
}

// AcquireLocksForUpdate acquires locks on partitions involved in update query
func (mlp *MLP) AcquireLocksForUpdate(updateQueries []UpdateQueryStruct) {
	for _, query := range updateQueries {
		startPartition := mlp.GetNodeByID(query.StartId).GetPartition()
		endPartition := mlp.GetNodeByID(query.EndId).GetPartition()

		// Acquire locks on partitions along the traversal path

		for startPartition != endPartition {
			fmt.Println("Acquiring lock for partition ", startPartition.GetId())
			if startPartition != nil {
				startPartition.Lock()
				mlp.AcquiredIDs = append(mlp.AcquiredIDs, startPartition.GetId())
			}
			startPartition = startPartition.GetParent()
		}

	}
}

func (mlp *MLP) AcquireLocksForUpdateEdge(updateQueries []UpdateQueryStruct) {
	for _, query := range updateQueries {

		startPartition := mlp.GetNodeByID(query.StartId).GetPartition()
		endPartition := mlp.GetNodeByID(query.EndId).GetPartition()

		for startPartition!= nil && startPartition != endPartition {
			startPartition = startPartition.GetParent()
			endPartition = endPartition.GetParent()
		}

		if startPartition == nil {
			continue
		}

		startEdgeId := query.StartId
		endEdgeId := query.EndId

		startPartition.EdgeNodeMap[startEdgeId][endEdgeId].Lock()
		borderNodes := startPartition.Border_Nodes
		startPartition = startPartition.GetParent()

		if startPartition!=nil && !startPartition.Affected {
			for startPartition != nil {
				for _, node := range borderNodes {
					for _, otherNode := range borderNodes {
						if startPartition.EdgeNodeMap[node.GetId()][otherNode.GetId()] != nil {
							startPartition.EdgeNodeMap[node.GetId()][otherNode.GetId()].Lock()
						}
					}
				}
				
				if startPartition.GetParent() != nil {
					//fmt.Println("Parent id ", startPartition.GetParent().GetId())
					startPartition = startPartition.GetParent()
				} else {
					break
				}
			}
			startPartition.Affected = true
		}

	}
}

// ReleaseLocks releases locks on all partitions
func (mlp *MLP) ReleaseLocks() {
	for _, id := range mlp.AcquiredIDs {
		node := mlp.GetNodeByID(id)
		if node != nil && node.GetPartition() != nil {
			node.GetPartition().Unlock()
		}
	}
}

func (mlp *MLP) PrintInfo() {
	print("This is all the info about the mlp\n")
	print("Level num: ", mlp.levelNum, "\n")
	for i := mlp.levelNum - 1; i >= 0; i-- {
		print("Level ", i, ": \n")
		for _, p := range mlp.GetPartitions(i) {

			print("Partition ", p.GetId(), ": \n")

			print("Shortcut Nodes: ")
			for _, n := range p.Shortcut_Nodes {
				print(n.GetId(), " ")
			}
			print("\n")

			for _, e := range p.Shortcut_Edges {
				print("Shortcut Edge ", e.GetId(), " from ", e.GetSrc(), " to ", e.GetDest(), " with weight ", e.GetWeight(), "\n")
			}

			print("Border Nodes: ")
			for _, n := range p.Border_Nodes {
				print(n.GetId(), " ")
			}
			print("\n")
			for _, e := range p.Border_Edges {
				print("Border Edge ", e.GetId(), " from ", e.GetSrc(), " to ", e.GetDest(), " with weight ", e.GetWeight(), "\n")
			}

			print("APSP Dist: \n")
			for from, distMap := range p.Apsp_Dist {
				for to, dist := range distMap {
					print("From ", from, " to ", to, " with distance ", dist, "\n")
				}
			}
		}
	}
}

// func main(){
// 	m := NewMLP(3)
// }
