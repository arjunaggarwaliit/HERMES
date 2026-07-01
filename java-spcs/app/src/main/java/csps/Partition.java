package csps;

import java.util.HashMap;
import java.util.Map;
import java.util.List;
import java.util.ArrayList;

public class Partition {
    private Map<Integer, Vertex> vertices;
    private Map<Long, Integer> costs;

    public Partition() {
        this.vertices = new HashMap<>();
        this.costs = new HashMap<>();
    }

    public boolean contains(int vertexId) {
        return vertices.containsKey(vertexId);
    }

    public void addVertex(int id, boolean isBridge, List<Integer> otherHosts) {
        vertices.put(id, new Vertex(id, isBridge, new ArrayList<>(otherHosts)));
    }

    public void addEdge(int src, int dst, int weight) {
        vertices.get(src).addEdge(dst, weight);
    }

    public Vertex getVertex(int vertexId) {
        return vertices.get(vertexId);
    }

    public boolean isBridge(int vertexId) {
        return vertices.get(vertexId).isBridge();
    }

    public List<Edge> getEdges(int vertexId) {
        return vertices.get(vertexId).getEdges();
    }

    public boolean containsCost(int src, int dst) {
        long key = ((long) src << 32) + dst;
        return costs.containsKey(key);
    }

    public void addCost(int src, int dst, int distance) {
        long key = ((long) src << 32) + dst;
        costs.put(key, distance);
    }

    public int getCost(int src, int dst) {
        long key = ((long) src << 32) + dst;
        return costs.get(key);
    }

    public int getCostSize() {
        return costs.size();
    }
}
