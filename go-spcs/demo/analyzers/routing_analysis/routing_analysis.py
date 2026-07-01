import os
import sys
import pandas as pd
import numpy as np
import matplotlib.pyplot as plt



def load_data(results_dir):
    """
    Load data from CSV file in the results directory.
    """
    csv_file = os.path.join(results_dir, 'routing_results.csv')
    converters = {'RunTime': lambda x: convert_runtime(x)}
    df = pd.read_csv(csv_file, converters=converters)
    return df

def convert_runtime(time_str):
    """
    Convert runtime string to milliseconds.
    """
    if time_str.endswith('ms'):
        return float(time_str.replace('ms', ''))
    elif time_str.endswith('µs'):
        return float(time_str.replace('µs', '')) / 1000
    elif time_str.endswith('s'):
        return float(time_str.replace('s', '')) * 1000
    else:
        return float(time_str)

def calculate_percentiles(dfs_map):
    """
    Calculate percentiles of time taken for both datasets.
    """
    percentiles = [50, 75, 90, 95, 99]

    percentile_data = {
        'percentile': percentiles,
    }

    for name, df in dfs_map.items():
        percentile_data[name] = np.percentile(df['RunTime'], percentiles)

    percentile_df = pd.DataFrame(percentile_data)
    return percentile_df

def plot_histogram(dfs_map, output_dir):
    """
    Plot histogram of time taken for both datasets.
    """
    plt.figure(figsize=(10, 6))

    colors = ['blue', 'red']    
    for name, df in dfs_map.items():
        plt.hist(df['RunTime'], bins=100, alpha=0.5, label=name, color=colors.pop(0))
    
    plt.xlabel('Time Taken (ms)')
    plt.ylabel('Frequency')
    plt.title('Distribution of Time Taken')
    plt.legend()
    plt.savefig(os.path.join(output_dir, 'time_taken_distribution.png'))

def main(dirs_map, output_dir):
    """
    Main function to generate and save plots comparing both datasets.
    """
    os.makedirs(output_dir, exist_ok=True)

    dfs_map = {}

    for name , path in dirs_map.items():
        dfs_map[name] = load_data(path)

    # Calculate percentiles of time taken
    percentile_df = calculate_percentiles(dfs_map)
    percentile_df.to_csv(os.path.join(output_dir, 'percentile_table.csv'))

    # Plot histogram
    plot_histogram(dfs_map, output_dir)

if __name__ == "__main__":
    
    dirs_map = {
        'ch': "..\\..\\runnable\\results\\graph_CH\\routing",
        'cosmos': "..\\..\\runnable\\results\\graph1\\routing",
    }

    output = 'graphs/'
    main(dirs_map, output)
