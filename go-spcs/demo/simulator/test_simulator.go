package simulator

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"src"
	"time"
	"utils"

	"gopkg.in/yaml.v2"
)

// Config struct to hold the configuration from config.yaml
type Config struct {
	Name       string       `yaml:"name"`
	Type       string       `yaml:"type"`
	ResultsDir  string    `yaml:"results_path"`
	MLPConfig  MLPConfig    `yaml:"mlp"`
	Include    []string     `yaml:"include"`
	Tests []TestConfig `yaml:"tests"`
}

// TestConfig struct to hold the configuration for each test
type Tests struct {
	Configs []Config `yaml:"configs"`
}

type TestConfig struct {
	Name 	     string    `yaml:"name"`
	RoutingQueriesPath  string    `yaml:"routing_queires_path"`
	UpdateQueriesPath   string    `yaml:"update_queries_path"`
	
}

type MLPConfig struct {
	GraphFile string `yaml:"graphfile"`
	Levels    int    `yaml:"levels"`
}

var (
	g   *src.Graph // Global graph variable
	err error      // Error variable
	testFuncsSPCS  = map[string]func(*src.MLP, []TestConfig, string){
		"routing_test":                runRoutingTestSPCS,
		"concurrency_test":            runConcurrencyTestSPCS,
		// "conflict_serialization_test": runConflictSerializationTest,
	}
	testFuncsDCH  = map[string]func(utils.Graph,[]TestConfig, string){
		"routing_test":                runRoutingTestDCH,
		// "stress_test":                 runStressTest,
		// "conflict_serialization_test": runConflictSerializationTest,
	}
	logger     *log.Logger
	logFile, _ = os.OpenFile("spcs_simulator_1.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	vertices = []int{}
	edges    = [][]int{}
)

func init() {
	logger = log.New(logFile, "", log.Ldate|log.Ltime)
}

func Simulate() {
	
	logger.Println("Starting simulation")

	dir, _ := os.Getwd()
	
	configPath := filepath.Join(dir, "../config/config.yaml")
	tests, _ := loadConfig(configPath)

	for _, config := range tests.Configs {	
		
		
		fmt.Println("Running simulation for config: ", config.Name)	
		resultsDir := filepath.Join(dir, config.ResultsDir)
		routingDir := filepath.Join(resultsDir, "routing")
		concurrencyDir := filepath.Join(resultsDir, "concurrency")
		err = os.MkdirAll(routingDir, 0755)
		err = os.MkdirAll(concurrencyDir, 0755)
		

		if config.Type == "spcs" {
			simulateSPCS(resultsDir, config)
		} else if config.Type == "dch"{
			simulateDCH(resultsDir, config)
		} else {
			fmt.Println("Unknown type of simulation")			
		}		
	}
}

func simulateDCH(resultsDir string, config Config){
	// Create results directory if it doesn't exist
	
	err := os.MkdirAll(resultsDir, 0755)
	if err != nil {
		logger.Fatalf("Failed to create results directory: %v", err)
	}

	// Save test config to test.config file in results directory
	testConfigFile := filepath.Join(resultsDir, "test.config")
	testConfigData, err := yaml.Marshal(config)
	if err != nil {
		logger.Fatalf("Failed to marshal test config: %v", err)
	}
	err = os.WriteFile(testConfigFile, testConfigData, 0644)
	if err != nil {
		logger.Fatalf("Failed to write test config to file: %v", err)
	}

	g := utils.Graph{}

	// Read the graph from the text file: If using the US Road Network, use the following line
	err = ReadGraphDCH(&g, config.MLPConfig.GraphFile)
	fmt.Println("Reading the graph from the text file done. No of vertices: ", len(vertices), " No of edges: ", len(edges))
	
	fmt.Println("Please wait until contraction hierarchy is prepared")
	timech := time.Now()
	g.PrepareContractionHierarchies()
	elapsedch := time.Since(timech)
	log.Println("Time for CH:", elapsedch)

	// Run tests specified in the config
	for _, testName := range config.Include {
		fmt.Println(testName)
		if testFunc, ok := testFuncsDCH[testName]; ok {
			logger.Printf("Running test: %s\n", testName)
			testFunc(g, config.Tests, resultsDir)
		} else {
			logger.Printf("Unknown test: %s\n", testName)
		}
	}
}

func simulateSPCS(resultsDir string, config Config){
	err := os.MkdirAll(resultsDir, 0755)
	if err != nil {
		logger.Fatalf("Failed to create results directory: %v", err)
	}

	// Save test config to test.config file in results directory
	testConfigFile := filepath.Join(resultsDir, "test.config")
	testConfigData, err := yaml.Marshal(config)
	if err != nil {
		logger.Fatalf("Failed to marshal test config: %v", err)
	}
	err = os.WriteFile(testConfigFile, testConfigData, 0644)
	if err != nil {
		logger.Fatalf("Failed to write test config to file: %v", err)
	}

	g = src.NewGraph()
	// Read the graph from the text file: If using the US Road Network, use the following line
	err = ReadGraphSPCS(&vertices, &edges, config.MLPConfig.GraphFile)
	fmt.Println("Reading the graph from the text file done. No of vertices: ", len(vertices), " No of edges: ", len(edges))

	//g.CreateSampleGraph(15, 30)
	g.CreateGraph(vertices, edges)

	mlp := src.NewMLP(config.MLPConfig.Levels)
	mlp.MLPConstruction(g)

	// Run tests specified in the config
	for _, testName := range config.Include {
		fmt.Println(testName)
		if testFunc, ok := testFuncsSPCS[testName]; ok {
			logger.Printf("Running test: %s\n", testName)
			testFunc(mlp, config.Tests, resultsDir)
		} else {
			logger.Printf("Unknown test: %s\n", testName)
		}
	}
}

func loadConfig(filename string) (Tests, error) {
	var tests Tests
	data, err := os.ReadFile(filename)
	if err != nil {
		logger.Fatalf("Failed to read config file: %v", err)
	}

	err = yaml.Unmarshal(data, &tests)
	if err != nil {
		logger.Fatalf("Failed to unmarshal config: %v", err)
	}

	fmt.Printf("Config: %+v\n", tests)
	return tests, nil
}



