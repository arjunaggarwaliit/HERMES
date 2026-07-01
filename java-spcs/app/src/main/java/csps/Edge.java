package csps;

public class Edge {
    private int dst;
    private int weight;

    public Edge(int dst, int weight) {
        this.dst = dst;
        this.weight = weight;
    }

    public int getDst() {
        return dst;
    }

    public int getWeight() {
        return weight;
    }
}
