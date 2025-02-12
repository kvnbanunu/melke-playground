package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/kvnbanunu/melke-playground/cli/internal/codegen"
)

func main() {
	// Define command line flags
	configFile := flag.String("config", "config.yaml", "Path to YAML configuration file")
	flag.Parse()

	// Read and parse the configuration file
	cfg, err := codegen.ParseConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	// Create the generator instance
	generator := codegen.NewGenerator(cfg)

	// Generate the code
	if err := generator.Generate(); err != nil {
		log.Fatalf("Failed to generate code: %v", err)
	}

	fmt.Println("Code generation completed successfully!")
}
