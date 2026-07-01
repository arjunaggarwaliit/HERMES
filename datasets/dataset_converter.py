def dataset_convert(data_path, output_file_name, nodes_threshold = 1000):
    
    # Read the file
    with open(data_path, "r") as file:
        lines = file.readlines()

    egdes_count = 0

    # Open the file to write
    with open(output_file_name, "w") as file:

        # Add the first line
        # file.write("source_id;target_id;one_way;length\n")

        for i in range(len(lines)):
            lines[i] = lines[i].strip()
            if lines[i][0] == 'a':

                # Split the line
                line = lines[i].split(" ")

                if int(line[1]) < nodes_threshold and int(line[2]) < nodes_threshold:
                    file.write(f"{line[1]};{line[2]};TF;{line[3]}\n")
                    egdes_count += 1
                else:
                    continue
    
    print(f"Edges count: {egdes_count}")
    print("Conversion completed")
                

if __name__ == "__main__":
    data_path = "master/USA-road-d.USA.gr"

    nodes_threshold = 100
    output_file_name = f"master/data_{nodes_threshold}.txt"
    dataset_convert(data_path, output_file_name, nodes_threshold)