import csv
import matplotlib.pyplot as plt

# Initialize empty lists to store data
avg_time_taken = []
total_time_taken = 0

# Read data from CSV file
with open('results_algo_two.csv', 'r') as file:
    reader = csv.reader(file)
    for row in reader:
        total_time_taken += int(row[4])/1000
        avg_time_taken.append(total_time_taken/(len(avg_time_taken)+1))

# Skip inital 40 values for average time taken
time_taken = avg_time_taken[40:]

avg_time_taken_2 = []
total_time_taken_2 = 0

# Read data from CSV file
with open('results_algo_ch.csv', 'r') as file:
    reader = csv.reader(file)
    for row in reader:
        total_time_taken_2 += int(row[4])/1000
        avg_time_taken_2.append(total_time_taken_2/(len(avg_time_taken_2)+1))

# Skip inital 40 values for average time taken
time_taken_2 = avg_time_taken_2[40:]

# Plotting the data
plt.figure(figsize=(10, 5))
plt.plot(time_taken, label='Our Algorithm')
plt.plot(time_taken_2, label='Algorithm CH')

# Adding labels and title
plt.xlabel('Query')
plt.ylabel('Time Taken (in ms)')
plt.title('Time Taken by Queries')

# Adding legend
plt.legend()

# Rotating x-axis labels for better readability
plt.xticks(rotation=45)

# Displaying the plot
plt.tight_layout()  # Adjust layout to prevent clipping of labels
plt.savefig('results_comparsion.png')
