import csv
import random

def modify_and_save(input_file, output_file, n, multiplier):
    selected_lines = []
    with open(input_file, 'r') as infile:
        lines = infile.readlines()
        selected_lines = random.sample(lines, n)
    
    modified_data = []
    for line in selected_lines:
        fields = line.strip().split(';')
        first_two_fields = fields[:2]
        last_field = int(fields[-1])
        modified_last_field = last_field * multiplier
        modified_data.append(first_two_fields + [str(modified_last_field)])

    with open(output_file, 'w', newline='') as csvfile:
        writer = csv.writer(csvfile)
        writer.writerows(modified_data)

#Example usage:
n = 1000  # Number of lines to randomly select
input_file = '../../../../datasets/master/data_100000.txt'
output_file = f"update_queries_{n}.csv"
multiplier = 2  # Multiplier to apply to the last field


modify_and_save(input_file, output_file, n, multiplier)
