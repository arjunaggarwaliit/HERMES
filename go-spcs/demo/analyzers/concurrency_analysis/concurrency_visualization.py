from matplotlib import pyplot as plt
import pandas as pd




def plot_timeline(df):
    # Filter out update and routing data
    update_df = df[df['query_type'] == 'UPDATE']
    routing_df = df[df['query_type'] == 'ROUTING']

    # Plot UPDATE queries
    for idx, row in update_df.iterrows():
        plt.plot([row['start_time'], row['end_time']], [row['duration'], row['duration']], color='blue', marker='o', markersize=5)

    # Plot ROUTING queries
    for idx, row in routing_df.iterrows():
        plt.plot([row['start_time'], row['end_time']], [row['duration'], row['duration']], color='red', marker='o', markersize=5)

    plt.xlabel('Timeline')
    plt.ylabel('Duration')
    plt.title('Timeline of Queries')
    plt.grid(True)
    # plt.xlabel('Query IDs')
    # plt.ylabel('Duration')
    # plt.title('Duration of Queries')
    # plt.legend()
    # plt.grid(True)
    plt.savefig('graphs/concurrency_timeline_v2.png')

if __name__ == "__main__":
    
    df = pd.read_csv('dfs/transformed_df_v2.csv')
    plot_timeline(df)