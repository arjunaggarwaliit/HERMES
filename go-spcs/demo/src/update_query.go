package src

import "fmt"

type UpdateQueryStruct struct {
	ID 	int
	StartId int
	EndId   int
	Weight  int
}

func BatchUpdateQuery(updateQueries []UpdateQueryStruct, mlp *MLP) {

	fmt.Println("Processing batch update queries")
	
	// mlp.AcquireLocksForUpdate(updateQueries)
	// fmt.Println("Acquired locks")
    // defer mlp.ReleaseLocks() // Release locks when the function returns

	mlp.AcquireLocksForUpdateEdge(updateQueries)
	
	// Map to store affected partitions and their corresponding level numbers
	affectedPartitions := make(map[*Partition]int)

	// Process each update query
	for _, query := range updateQueries {
		// Retrieve nodes and their partitions
		startNode := mlp.nodeMap[query.StartId]
		endNode := mlp.nodeMap[query.EndId]
		startPartition := startNode.GetPartition()
		endPartition := endNode.GetPartition()

		// Find the common ancestor partition and calculate the level number
		levelNum := mlp.GetLevelNum()
		for startPartition.GetId() != endPartition.GetId() {
			startPartition = startPartition.GetParent()
			endPartition = endPartition.GetParent()
			levelNum--
		}

		// Update edge weights
		updateEdgeWeights(startNode, endNode, query.Weight)

		// Mark affected partitions
		markAffectedPartitions(startNode.GetPartition(), levelNum, affectedPartitions)
	}

	// Print affected partitions
	// fmt.Println("Affected partitions:")
	// for partition, levelNum := range affectedPartitions {
	// 	fmt.Printf("Partition ID: %d, Level Number: %d \n", partition.GetId(), levelNum)
	// }

	fmt.Println("Recomputing shortcut networks")

	// Recompute shortcut networks for affected partitions
	recomputeShortcutNetworks(affectedPartitions, mlp.GetLevelNum())
}

// Function to update edge weights
func updateEdgeWeights(startNode, endNode *Node, weight int) {
	startID, endID := startNode.GetId(), endNode.GetId()
	startNode.NodeEdgeMap[endID].SetWeight(weight)
	endNode.NodeEdgeMap[startID].SetWeight(weight)
}

// Function to mark affected partitions and their level numbers
func markAffectedPartitions(partition *Partition, levelNum int, affectedPartitions map[*Partition]int) {
	for i := 0; i < levelNum; i++ {
		partition.Affected = true
		affectedPartitions[partition]++
		partition = partition.GetParent()
	}
}

// Function to recompute shortcut networks for affected partitions
func recomputeShortcutNetworks(affectedPartitions map[*Partition]int, maxLevel int) {
	// Group affected partitions by their levels
	partitionsByLevel := make(map[int][]*Partition)
	for partition, levelNum := range affectedPartitions {
		partitionsByLevel[levelNum] = append(partitionsByLevel[levelNum], partition)
	}

	// Iterate over levels in reverse order
	for level := maxLevel; level >= 0; level-- {
		partitions, exists := partitionsByLevel[level]
		if exists {
			// Recompute shortcut networks for partitions at the current level
			for _, partition := range partitions {
				partition.RebuildShortcutNetwork()
				partition.RecomputeShortcutNetwork()
				partition.Affected = false
				partition.ReleaseLocks()
			}
		}
	}
}
