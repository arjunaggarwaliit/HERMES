package csps;

import java.io.BufferedReader;
import java.io.FileReader;
import java.io.IOException;
import java.util.*;
import java.util.stream.Collectors;

public class Algorithms {
    private static final int INF = Integer.MAX_VALUE;

    static class PQElement {
        int label, dst, cost;

        PQElement(int label, int dst, int cost) {
            this.label = label;
            this.dst = dst;
            this.cost = cost;
        }
    }

    static class PQCompare implements Comparator<PQElement> {
        @Override
        public int compare(PQElement obj1, PQElement obj2) {
            return Integer.compare(obj1.cost, obj2.cost);
        }
    }

    private static final Map<Long, Integer> globalDistance = new HashMap<>();

    public static Index runAlgorithmOne(String inputFilename, int numOfLabels) {
        Index index = new Index();

        for (int i = 0; i < numOfLabels; ++i) {

            index.createPartition(i);
        }

        try (BufferedReader input = new BufferedReader(new FileReader(inputFilename))) {
            String line;
            Map<Integer, Set<Integer>> otherHostsMap = new HashMap<>();

            while ((line = input.readLine()) != null) {
                String[] parts = line.split(" ");
                int src = Integer.parseInt(parts[0]);
                int dst = Integer.parseInt(parts[1]);
                int weight = Integer.parseInt(parts[2]);
                int label = Integer.parseInt(parts[3]);
                
                otherHostsMap.computeIfAbsent(src, k -> new HashSet<>()).add(label);
            }

            System.err.println("Finish scanning");

            try (BufferedReader input2 = new BufferedReader(new FileReader(inputFilename))) {
                while ((line = input2.readLine()) != null) {
                    String[] parts = line.split(" ");
                    int src = Integer.parseInt(parts[0]);
                    int dst = Integer.parseInt(parts[1]);
                    int weight = Integer.parseInt(parts[2]);
                    int label = Integer.parseInt(parts[3]);
                    
                    Partition par = index.getPartition(label);

                    if (!par.contains(src)) {
                        boolean isBridge = false;
                        List<Integer> otherHosts = otherHostsMap.getOrDefault(src, Collections.emptySet())
                                .stream()
                                .filter(l -> l != label)
                                .collect(Collectors.toList());

                        if (!otherHosts.isEmpty()) {
                            isBridge = true;
                        }

                        par.addVertex(src, isBridge, new ArrayList<>(otherHosts));
                    }

                    if (!par.contains(dst)) {
                        boolean isBridge = false;
                        List<Integer> otherHosts = otherHostsMap.getOrDefault(dst, Collections.emptySet())
                                .stream()
                                .filter(l -> l != label)
                                .collect(Collectors.toList());

                        if (!otherHosts.isEmpty()) {
                            isBridge = true;
                        }

                        par.addVertex(dst, isBridge, new ArrayList<>(otherHosts));
                    }

                    par.addEdge(src, dst, weight);
                }
            }

        } catch (IOException e) {
            System.err.println("Input file cannot be opened!");
            System.exit(-1);
        }

        return index;
    }

    private static long getKey(int a, int b) {
        return ((long) a << 32) + b;
    }

    private static void insertIfRelaxed(PriorityQueue<PQElement> q, int label, int dst, int distance) {
        long key = getKey(label, dst);
        Integer globalDist = globalDistance.get(key);

        if (globalDist == null) {
            q.add(new PQElement(label, dst, distance));
            globalDistance.put(key, distance);
        } else if (globalDist > distance) {
            q.add(new PQElement(label, dst, distance));
            globalDistance.put(key, distance);
        }
    }

    public static int runAlgorithmTwo(Index index, int src, int dst, Set<Integer> labels) {
        PriorityQueue<PQElement> q = new PriorityQueue<>(new PQCompare());
        globalDistance.clear();

        int minPar = index.getMinPr(src, labels);
        if (minPar == INF) {
            return INF;
        }

        insertIfRelaxed(q, minPar, src, 0);

        while (!q.isEmpty()) {
            PQElement currentVertex = q.poll();
            long key = getKey(currentVertex.label, currentVertex.dst);

            if (currentVertex.cost > globalDistance.getOrDefault(key, 0)) {
                continue;
            }

            System.out.println("Now visiting " + currentVertex.dst);

            if (currentVertex.dst == dst) {
                return currentVertex.cost;
            }

            if (!index.containsCost(currentVertex.label, currentVertex.dst, dst)) {
                System.out.println("Running Dijkstra");
                Partition par = index.getPartition(currentVertex.label);
                Map<Integer, Integer> distances = new HashMap<>();
                distances.put(currentVertex.dst, 0);

                PriorityQueue<Map.Entry<Integer, Integer>> djQ = new PriorityQueue<>(Comparator.comparingInt(Map.Entry::getValue));
                djQ.add(new AbstractMap.SimpleEntry<>(0, currentVertex.dst));

                while (!djQ.isEmpty()) {
                    Map.Entry<Integer, Integer> v = djQ.poll();
                    System.out.println("Internal: now visiting " + v.getValue());

                    if (distances.containsKey(v.getValue()) && v.getKey() == distances.get(v.getValue())) {
                        if (par.isBridge(v.getValue()) && !v.getValue().equals(currentVertex.dst)) {
                            par.addCost(currentVertex.dst, v.getValue(), v.getKey());
                            par.getVertex(currentVertex.dst).addBridgeEdge(v.getValue());
                        } else if (v.getValue() == dst) {
                            par.addCost(currentVertex.dst, v.getValue(), v.getKey());
                        }

                        for (Edge edge : par.getEdges(v.getValue())) {
                            int newWeight = v.getKey() + edge.getWeight();
                            if (!distances.containsKey(edge.getDst()) || distances.get(edge.getDst()) > newWeight) {
                                distances.put(edge.getDst(), newWeight);
                                djQ.add(new AbstractMap.SimpleEntry<>(newWeight, edge.getDst()));
                            }
                        }
                    }
                }

                if (!distances.containsKey(dst)) {
                    par.addCost(currentVertex.dst, dst, INF);
                }
            }

            System.out.println("Retrieving from index");

            int costToDst = index.getCost(currentVertex.label, currentVertex.dst, dst);
            if (costToDst != INF) {
                insertIfRelaxed(q, currentVertex.label, dst, currentVertex.cost + costToDst);
            }

            for (int bridgeV : index.getBridgeEdges(currentVertex.label, currentVertex.dst)) {
                int costToBridge = index.getCost(currentVertex.label, currentVertex.dst, bridgeV);

                for (int otherLabel : index.getOtherHosts(currentVertex.label, bridgeV)) {
                    if (labels.contains(otherLabel)) {
                        insertIfRelaxed(q, otherLabel, bridgeV, currentVertex.cost + costToBridge);
                    }
                }
            }

            if (index.isBridge(currentVertex.label, currentVertex.dst)) {
                for (int otherLabel : index.getOtherHosts(currentVertex.label, currentVertex.dst)) {
                    if (labels.contains(otherLabel)) {
                        insertIfRelaxed(q, otherLabel, currentVertex.dst, currentVertex.cost);
                    }
                }
            }
        }

        return INF;
    }
}
