package csps;

import java.io.BufferedReader;
import java.io.FileReader;
import java.io.IOException;
import java.util.HashSet;
import java.util.Set;

public class CSPS_test {
    private static final int INF = Integer.MAX_VALUE;

    public static void main(String[] args) {
        
        String PATH = "src/main/datasets/youtube_subset.txt";
        int NUM_OF_LABELS = 5;
        int NUM_OF_VERTICES = 2000;

        System.out.println("----Stage 1: initialization----");
        String input = PATH;
        int num_of_labels =  NUM_OF_LABELS;
        Index index = Algorithms.runAlgorithmOne(input, num_of_labels);
        System.out.println("Created index");

        System.out.println("----Stage 2: query processing----");
        Set<Integer> labels = new HashSet<>();
        labels.add(2);
        labels.add(0);
        int src = 74, dst = 4926;

        for (int i = 0; i < 5; ++i) {
            long t1 = System.currentTimeMillis();
            int cost = Algorithms.runAlgorithmTwo(index, src, dst, labels);
            long t2 = System.currentTimeMillis();

            if (cost == INF)
                System.out.println("Cannot find a route");
            else {
                System.out.println("Finished, cost is " + cost + " from " + src + " to " + dst + " with labels " +
                        labels.toString());
            }

            System.out.println("Time taken to run the query: " + (t2 - t1) + "ms");
        }
    }
}
