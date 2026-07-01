package csps;
import java.util.HashSet;
import java.util.Set;
import java.util.Random;
import java.util.concurrent.ThreadLocalRandom;
import java.util.stream.Collectors;

public class CSPS_main {
    private static final int INF = Integer.MAX_VALUE;

    public static void main(String[] args) {

        String PATH = "src/main/datasets/youtube_subset.txt";
        int NUM_OF_LABELS = 5;
        int NUM_OF_VERTICES = 2000;
        
        System.err.println("----Stage 1: initialization----");

        String input = PATH;
        int num_of_labels = NUM_OF_LABELS;
        Index index = Algorithms.runAlgorithmOne(input, num_of_labels);
        System.err.println("Created index");

        System.err.println("----Stage 2: query processing----");

        int num_of_vertices = NUM_OF_VERTICES;
        Random rng = new Random();
        int num_of_queries = 1000;

        for (int iter = 1; iter <= num_of_queries; ++iter) {
            // generate query labels
            Set<Integer> labels = new HashSet<>();
            while (labels.size() < num_of_labels / 2)
                labels.add(rng.nextInt(num_of_labels));

            for (int label : labels)
                System.err.println("label: " + label);

            // run query
            System.err.println("query num: " + iter);
            int src;
            do {
                src = rng.nextInt(num_of_vertices);
            } while (index.getMinPr(src, labels) == INF);
            int dst = rng.nextInt(num_of_vertices);
            System.err.println("src: " + src + ", dst: " + dst);

            long t1 = System.currentTimeMillis();
            int cost = Algorithms.runAlgorithmTwo(index, src, dst, labels);
            long t2 = System.currentTimeMillis();

            if (cost == INF)
                System.err.println("Cannot find a route");
            else {
                System.err.println("Finished, cost is " + cost + " from " + src + " to " + dst +
                        " with labels " + labels.stream().map(Object::toString).collect(Collectors.joining(" ")));
            }

            int total_entries = 0;
            for (int l = 0; l < num_of_labels; ++l)
                total_entries += index.getPartition(l).getCostSize();
            System.err.println("Now there are " + total_entries + " entries in index");
            System.out.println(iter + " " + total_entries);

            System.err.println("Time taken to run the query: " + (t2 - t1) + "ms");
            System.err.println("========================================");
        }
    }
}
