package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// Define flags
	providerFilter := flag.String("provider", "", "Filter models by provider name (partial match)")
	modelFilter := flag.String("model", "", "Filter models by model name (partial match)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "llmls - List and manage LLM models\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  llmls [flags]\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	// Fetch models from OpenRouter
	models, err := FetchModels()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Filter models
	models = FilterModels(models, *providerFilter, *modelFilter)

	// Sort by creation date descending
	SortModelsByCreatedDesc(models)

	// Display models
	DisplayModels(models)
}
