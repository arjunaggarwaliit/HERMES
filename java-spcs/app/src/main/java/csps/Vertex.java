package csps;

import java.util.ArrayList;
import java.util.List;

public class Vertex {
    private int id;
    private boolean isBridge;
    private List<Integer> otherHosts;
    private List<Integer> bridgeEdges;
    private List<Edge> edges;

    public Vertex(int id, boolean isBridge, List<Integer> otherHosts) {
        this.id = id;
        this.isBridge = isBridge;
        this.otherHosts = new ArrayList<>(otherHosts);
        this.bridgeEdges = new ArrayList<>();
        this.edges = new ArrayList<>();
    }

    public int getId() {
        return id;
    }

    public boolean isBridge() {
        return isBridge;
    }

    public List<Integer> getOtherHosts() {
        return new ArrayList<>(otherHosts);
    }

    public void addEdge(int dst, int weight) {
        edges.add(new Edge(dst, weight));
    }

    public List<Edge> getEdges() {
        return new ArrayList<>(edges);
    }

    public List<Integer> getBridgeEdges() {
        return new ArrayList<>(bridgeEdges);
    }

    public void addBridgeEdge(int dst) {
        bridgeEdges.add(dst);
    }
}

