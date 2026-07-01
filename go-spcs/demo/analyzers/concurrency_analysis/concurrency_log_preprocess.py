import pandas as pd


# Initialize lists to store data
timestamps = []
query_types = []
query_ids = []
durations = []

def extract_data_from_log(log_file_path):
    timestamps = []
    query_types = []
    query_ids = []
    durations = []

    with open(log_file_path, 'r') as file:
        for line in file:
            parts = line.split()
            if len(parts) >= 7:
                timestamp = ' '.join(parts[:2])
                query_type = ' '.join(parts[2:5]).strip('[]')  # Get the whole query type

                # Strip the [ from the query type
                
                query_id = parts[5]

                duration_str = parts[-1]

                if duration_str.endswith('µs'):
                    duration = float(duration_str[:-2]) / 1000000 
                elif duration_str.endswith('ms'):
                    duration = float(duration_str[:-2]) / 1000  # Duration in milliseconds, convert to seconds
                elif duration_str.endswith('s'):
                    duration = float(duration_str[:-1])  # Duration in seconds
                else:
                    print("Unknown duration format:", duration_str)
                    continue

                timestamps.append(timestamp)
                query_types.append(query_type)
                query_ids.append(query_id)
                durations.append(duration)
            else:
                print("Skipping line:", line)

    df = pd.DataFrame({
        'Timestamp': timestamps,
        'Type': query_types,
        'ID': query_ids,
        'Duration': durations
    })
    return df



def transform_data(original_df):
    # Initialize an empty DataFrame to store the new data
    new_data = []

    for query_id in original_df["ID"].unique():
        # Get the rows corresponding to the current query ID
        query_rows = original_df[original_df["ID"] == query_id]

        # Initialize variables to store start and end times for each query type
        start_times = {}
        end_times = {}
        
        # Check if both "UPDATE" and "ROUTING" events exist for the current query ID
        if "UPDATE QUERY START" in query_rows["Type"].values and "ROUTING QUERY START" in query_rows["Type"].values:
            # Iterate through each query type
            for query_type in ["UPDATE QUERY", "ROUTING QUERY"]:
                # Extract the start and end events for the current query type
                start_event = f"{query_type} START"
                end_event = f"{query_type} END"
                
                
                    # Get the start and end times for the current query type
                start_time = query_rows[query_rows["Type"] == start_event]["Duration"].iloc[0]
                end_time = query_rows[query_rows["Type"] == end_event]["Duration"].iloc[0]
                
                # Store start and end times for the current query type
                start_times[query_type] = start_time
                end_times[query_type] = end_time
            
            # Append the data to the list
            new_data.append({
                "query_id": query_id,
                "query_type": "UPDATE",
                "start_time": start_times.get("UPDATE QUERY", None),
                "end_time": end_times.get("UPDATE QUERY", None)
            })
            
            new_data.append({
                "query_id": query_id,
                "query_type": "ROUTING",
                "start_time": start_times.get("ROUTING QUERY", None),
                "end_time": end_times.get("ROUTING QUERY", None)
            })
        elif "UPDATE QUERY START" in query_rows["Type"].values and "ROUTING QUERY START" not in query_rows["Type"].values:
            # Only "UPDATE" events exist for the current query ID
            start_time = query_rows[query_rows["Type"] == "UPDATE QUERY START"]["Duration"].iloc[0]
            end_time = query_rows[query_rows["Type"] == "UPDATE QUERY END"]["Duration"].iloc[0]
            
            new_data.append({
                "query_id": query_id,
                "query_type": "UPDATE",
                "start_time": start_time,
                "end_time": end_time
            })
        elif "ROUTING QUERY START" in query_rows["Type"].values and "UPDATE QUERY START" not in query_rows["Type"].values:
            # Only "ROUTING" events exist for the current query ID
            start_time = query_rows[query_rows["Type"] == "ROUTING QUERY START"]["Duration"].iloc[0]
            end_time = query_rows[query_rows["Type"] == "ROUTING QUERY END"]["Duration"].iloc[0]
            
            new_data.append({
                "query_id": query_id,
                "query_type": "ROUTING",
                "start_time": start_time,
                "end_time": end_time
            })


    # Create a DataFrame from the list
    new_df = pd.DataFrame(new_data)
    new_df['duration'] = new_df['end_time'] - new_df['start_time']
    return new_df


if __name__ == "__main__":

    version = 2
    # log_file_path = 'dfs/log_backup.txt'
    log_file_path = '../../../results/log_v1.txt'
    original_df = extract_data_from_log(log_file_path)
    original_df.to_csv(f"dfs/concurrency_results_v{version}.csv", index=False)

    original_df = pd.read_csv(f"dfs/concurrency_results_v{version}.csv")
    results = transform_data(original_df)
    results.to_csv(f"dfs/transformed_df_v{version}.csv", index=False)