import os
import argparse
import sys
import pandas as pd
import matplotlib.pyplot as plt

def plot_concurrency_results(input_dir):
    # Read CSV file
    csv_file = os.path.join(input_dir, "concurrency/exec_times.csv")
    data = pd.read_csv(csv_file)

    # Convert time_taken to microseconds for better visualization
    data['time_taken'] = pd.to_timedelta(data['time_taken']).dt.total_seconds() * 1e6

    # Separate data into two dataframes based on query_type
    routing_data = data[data['query_type'] == 'routing']
    update_data = data[data['query_type'] == 'update']

    # Plot line chart for routing data
    plt.figure(figsize=(10, 6))
    plt.scatter(routing_data.index, routing_data['time_taken'], label='Routing', color='blue')

    # Plot line chart for update data
    plt.scatter(update_data.index, update_data['time_taken'], label='Update', color='red')

    plt.xlabel('Execution Order')
    plt.ylabel('Time Taken (µs)')
    plt.title('Concurrency Results')
    plt.legend()
    plt.savefig('concurrency_analysis/graphs/concurrency_results.png')

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python main.py <results_directory_1>")
        sys.exit(1)

    results_directory1 = sys.argv[1]

    plot_concurrency_results(results_directory1)
