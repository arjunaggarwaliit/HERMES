package csps;

import java.util.HashMap;
import java.util.Map;
import java.util.Set;
import java.util.List;

public class Index {
    private final Map<Integer, Partition> partitions = new HashMap<>();
    private static final int INF = Integer.MAX_VALUE;

    public void createPartition(int label) {
        partitions.put(label, new Partition());
    }

    public int getMinPr(int src, Set<Integer> labels) {
        int min = INF;
        int selectedLabel = INF;

        for (Map.Entry<Integer, Partition> entry : partitions.entrySet()) {
            int currentLabel = entry.getKey();
            Partition par = entry.getValue();

            if (par.contains(src) && labels.contains(currentLabel)) {
                for (Edge edge : par.getEdges(src)) {
                    if (edge.getWeight() < min) {
                        min = edge.getWeight();
                        selectedLabel = currentLabel;
                    }
                }
            }
        }

        return selectedLabel;
    }

    public boolean isBridge(int label, int vertexId) {
        return getPartition(label).isBridge(vertexId);
    }

    public List<Integer> getOtherHosts(int label, int vertexId) {
        return getPartition(label).getVertex(vertexId).getOtherHosts();
    }

    public Partition getPartition(int label) {
        return partitions.get(label);
    }

    public int getCost(int label, int src, int dst) {
        return partitions.get(label).getCost(src, dst);
    }

    public boolean containsCost(int label, int src, int dst) {
        return partitions.get(label).containsCost(src, dst);
    }

    public List<Integer> getBridgeEdges(int label, int vertexId) {
        return partitions.get(label).getVertex(vertexId).getBridgeEdges();
    }
}
