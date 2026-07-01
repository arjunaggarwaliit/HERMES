import os
import sys
import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
from sklearn.metrics import confusion_matrix


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
    
def group_data_into_batches(df, num_batches):
    """
    Group data into batches based on percentiles of the Query ID.
    """
    df['Batch'] = pd.qcut(df.index, num_batches, labels=False)
    return df

def plot_average_performance_time(dfs_map, output_dir):
    """
    Plot average performance time for each batch for all datasets.
    """
    plt.figure(figsize=(10, 6))

    for name, df in dfs_map.items():
        avg_performance_time = df.groupby('Batch')['RunTime'].mean()
        plt.plot(avg_performance_time.index, avg_performance_time.values, label=name)

    plt.xlabel('Batch')
    plt.ylabel('Average Performance Time (ms)')
    plt.title('Average Performance Time per Batch')
    plt.legend()
    plt.savefig(os.path.join(output_dir, 'average_performance_time.png'))
    plt.close()

def plot_performance_over_time(dfs_map, output_dir):
    """
    Plot performance over time for all datasets.
    """
    plt.figure(figsize=(10, 6))

    for name, df in dfs_map.items():
        plt.plot(df.index, df['RunTime'], label=name)

    plt.xlabel('Query ID')
    plt.ylabel('Run Time (ms)')
    plt.title('Performance Over Time')
    plt.legend()
    plt.savefig(os.path.join(output_dir, 'performance_over_time.png'))
    plt.close()


def plot_cost_vs_runtime(dfs_map, output_dir):
    """
    Plot cost vs. runtime for all datasets.
    """
    plt.figure(figsize=(10, 6))

    for name, df in dfs_map.items():
        plt.scatter(df['Cost'], df['RunTime'], label=name)

    plt.ylabel('Run Time (ms)')
    plt.xlabel('Cost')
    plt.title('Cost vs. Run Time')
    plt.legend()
    plt.savefig(os.path.join(output_dir, 'cost_vs_runtime.png'))
    plt.close()

def create_confusion_matrix(dirs_map, results_dir):
    """
    Create and save confusion matrix for cost comparison.
    """
    df1 = dirs_map['ch']
    df2 = dirs_map['cosmos']

    # Classify costs as match or mismatch
    expected_costs = df1['ExpectedCost'].values
    actual_costs_1 = (df1['Cost'].values)
    actual_costs_2 = df2['Cost'].values

    # Calculate confusion matrix
    cm1 = confusion_matrix(expected_costs, actual_costs_1, )
    cm2 = confusion_matrix(expected_costs, actual_costs_2)

    # Plot and save confusion matrix
    fig, axes = plt.subplots(1, 2, figsize=(12, 5))
    ax1, ax2 = axes.flatten()

    ax1.matshow(cm1, cmap='Oranges')
    ax1.set_title('ch')
    ax1.set_xlabel('Actual Cost')
    ax1.set_ylabel('Expected Cost')

    ax2.matshow(cm2, cmap='Blues')
    ax2.set_title('cosmos')
    ax2.set_xlabel('Actual Cost')
    ax2.set_ylabel('Expected Cost')

    plt.tight_layout()
    plt.savefig(os.path.join(results_dir, 'confusion_matrix.png'))
    plt.close()

# Modify the main function to pass dfs_map to plotting functions
def main(dirs_map, output_dir):
    """
    Main function to generate and save plots comparing all datasets.
    """
    os.makedirs(output_dir, exist_ok=True)

    dfs_map = {}

    for name, results_dir in dirs_map.items():
        # Load data for each dataset
        dfs_map[name] = load_data(results_dir)
        
        # Group data into batches based on percentiles
        num_batches = 10  # You can adjust the number of batches as needed
        dfs_map[name] = group_data_into_batches(dfs_map[name], num_batches)

    # Plot and save plots for all datasets
    plot_average_performance_time(dfs_map, output_dir)
    plot_performance_over_time(dfs_map, output_dir)
    plot_cost_vs_runtime(dfs_map, output_dir)
    create_confusion_matrix(dfs_map, output_dir)


if __name__ == "__main__":
    dirs_map = {
        'ch': "..\\..\\runnable\\results\\graph_CH\\routing",
        'cosmos': "..\\..\\runnable\\results\\graph1\\routing",
    }

    output = 'graphs/'
    main(dirs_map, output)
