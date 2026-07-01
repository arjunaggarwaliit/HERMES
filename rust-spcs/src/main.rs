use fast_paths::{InputGraph};
use std::fs::read_to_string;
use std::time::Instant;

fn add_all_edges(input_graph: &mut InputGraph, threshold: usize) {
    // Read the file "USA-road-d.NY.gr" and add all edges to the graph
    let file = read_to_string("../../datasets/USA-road-d.USA.gr").unwrap();

    let mut lines_count = 0;

    for line in file.lines() {
        let mut words = line.split_whitespace();
        let first_word = words.next().unwrap();
        if first_word == "a" {
            let node1 = words.next().unwrap().parse::<usize>().unwrap();
            let node2 = words.next().unwrap().parse::<usize>().unwrap();
            let weight = words.next().unwrap().parse::<usize>().unwrap();
            // println!("{} {} {}", node1, node2, weight);
            input_graph.add_edge(node1, node2, weight);
        }

        lines_count += 1;
        if lines_count > threshold {
            break;
        }
    }
}

fn main(){
    // begin with an empty graph
    let mut input_graph = InputGraph::new();

    // add an edge between nodes with ID 0 and 6, the weight of the edge is 12.
    // Note that the node IDs should be consecutive, if your graph has N nodes use 0...N-1 as node IDs,
    // otherwise performance will degrade.

    let mut threshold = String::new();

    println!("Enter the threshold: ");
    std::io::stdin().read_line(&mut threshold).expect("Failed to read line");

    let threshold: usize = threshold.trim().parse().expect("Please type a number!");

    let start = Instant::now();
    add_all_edges(&mut input_graph, threshold);
    // ... add many more edges here

    // freeze the graph before using it (you cannot add more edges afterwards, unless you call thaw() first)
    input_graph.freeze();

    // prepare the graph for fast shortest path calculations. note that you have to do this again if you want to change the
    // graph topology or any of the edge weights
    let fast_graph = fast_paths::prepare(&input_graph);

    let duration = start.elapsed();

    println!("Time elapsed in graph_processing is: {:?}", duration);

    while true{
        // Take the input from the user
        let mut source_id = String::new();
        let mut target_id = String::new();

        println!("Enter the source node id: ");
        std::io::stdin().read_line(&mut source_id).expect("Failed to read line");

        println!("Enter the target node id: ");
        std::io::stdin().read_line(&mut target_id).expect("Failed to read line");

        let source_id: usize = source_id.trim().parse().expect("Please type a number!");
        let target_id: usize = target_id.trim().parse().expect("Please type a number!");

        let start = Instant::now();
        // calculate the shortest path between nodes with ID 8 and 6 
        let shortest_path = fast_paths::calc_path(&fast_graph, source_id, target_id);

        
        match shortest_path {
            Some(p) => {
                // the weight of the shortest path
                let _weight = p.get_weight();
                
                // all nodes of the shortest path (including source and target)
                let _nodes = p.get_nodes();

                println!("Path: {:?}", _nodes);
                println!("Path length: {:?}", _nodes.len());
                println!("Weight: {:?}", _weight);
            },
            None => {
                // no path has been found (nodes are not connected in this graph)
            }
        }

        let duration = start.elapsed();
        println!("Time elapsed in computing shortest path is: {:?}", duration);

        let mut choice = String::new();
        println!("Do you want to continue? (y/n)");
        std::io::stdin().read_line(&mut choice).expect("Failed to read line");

        if choice.trim() == "n"{
            break;
        }
    }
}