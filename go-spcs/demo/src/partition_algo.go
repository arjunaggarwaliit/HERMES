package src

func Partition_algo(Partition *Partition) (*Partition, *Partition) {
	// Divides the partition into two parts and updates them directly

	// Create a new partition
	partition1 := NewPartition()
	partition2 := NewPartition()

	// Initialize the level of the partitions
	partition1.SetLevel(Partition.GetLevel() + 1)
	partition2.SetLevel(Partition.GetLevel() + 1)
	partition1.SetParent(Partition)
	partition2.SetParent(Partition)

	// Add the partitions to the parent partition
	Partition.AddChild(partition1)
	Partition.AddChild(partition2)

	// Get nodes
	nodes := Partition.GetNodes()

	for _, n := range nodes {
		var scoreP1 = 0
		var scoreP2 = 0

		for _, e := range n.GetEdges() {
			if partition1.contains(e.GetOtherNode(n)) {
				//scoreP1 += e.GetWeight()
				scoreP1 += 1
			} else if partition2.contains(e.GetOtherNode(n)) {
				//scoreP2 += e.GetWeight()
				scoreP2 += 1
			}
		}

		if scoreP1 > scoreP2 || (scoreP1 == scoreP2 && partition1.size() < partition2.size()) {
			partition1.AddNode(n)
		} else {
			partition2.AddNode(n)
		}
	}

	// Run the process 20 times
	for i := 0; i < 20; i++ {
		for _, n := range nodes {
			var scoreP1 = 0
			var scoreP2 = 0

			for _, e := range n.GetEdges() {
				if partition1.contains(e.GetOtherNode(n)) {
					//scoreP1 += e.GetWeight()
					scoreP1 += 1
				} else if partition2.contains(e.GetOtherNode(n)) {
					//scoreP2 += e.GetWeight()
					scoreP2 += 1
				}
			}

			if scoreP1 > scoreP2 || (scoreP1 == scoreP2 && partition1.size() < partition2.size()) {
				if !partition1.contains(n) {
					partition2.DeleteNode(n)
					partition1.AddNode(n)
				}
			} else {
				if !partition2.contains(n) {
					partition1.DeleteNode(n)
					partition2.AddNode(n)
				}
			}
		}
	}

	// Return the updated partitions
	return partition1, partition2
}
