package main

import (
	"flag"
	"fmt"
	"os"
)

var version = "dev"

func showHelp() {
	fmt.Fprintf(os.Stderr, "Usage: llmls [options] [pattern]\n\n")
	fmt.Fprintf(os.Stderr, "List LLM models from OpenRouter and Ollama.\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	fmt.Fprintf(os.Stderr, "  --detail         Show detailed model information\n")
	fmt.Fprintf(os.Stderr, "  --ollama-host    Ollama server URL (default: $OLLAMA_HOST or http://localhost:11434)\n")
	fmt.Fprintf(os.Stderr, "  -h, --help       Show this help message\n")
	fmt.Fprintf(os.Stderr, "  -v, --version    Show version information\n\n")
	fmt.Fprintf(os.Stderr, "Subcommands:\n")
	fmt.Fprintf(os.Stderr, "  providers        List all provider names\n\n")
	fmt.Fprintf(os.Stderr, "Examples:\n")
	fmt.Fprintf(os.Stderr, "  llmls                   List all models\n")
	fmt.Fprintf(os.Stderr, "  llmls \"anthropic/*\"     List Anthropic models\n")
	fmt.Fprintf(os.Stderr, "  llmls --detail cohere   Show detailed Cohere models\n")
	fmt.Fprintf(os.Stderr, "  llmls providers         List all providers\n")
}

func main() {
	if len(os.Args) < 2 {
		// Default behavior: list all models
		listModelsCommand(nil)
		return
	}

	subcommand := os.Args[1]

	switch subcommand {
	case "--help", "-h":
		showHelp()
		return
	case "--version", "-v":
		fmt.Printf("llmls version %s\n", version)
		return
	case "providers":
		providersCommand()
	default:
		// If not a subcommand, treat as search pattern
		listModelsCommand(os.Args[1:])
	}
}

func listModelsCommand(args []string) {
	fs := flag.NewFlagSet("llmls", flag.ExitOnError)
	detail := fs.Bool("detail", false, "Display detailed model information")
	ollamaHost := fs.String("ollama-host", "", "Ollama server URL (default: $OLLAMA_HOST or http://localhost:11434)")

	fs.Usage = showHelp

	if args != nil {
		fs.Parse(args)
	} else {
		fs.Parse(os.Args[1:])
	}

	// Get search pattern from positional argument
	pattern := ""
	if fs.NArg() > 0 {
		pattern = fs.Arg(0)
	}

	// Fetch models from OpenRouter
	models, err := FetchModels()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Fetch models from Ollama and merge
	ollamaModels := FetchOllamaModels(GetOllamaHost(*ollamaHost))
	models = append(models, ollamaModels...)

	// Filter models by pattern
	models = FilterModels(models, pattern)

	// Sort by creation date descending
	SortModelsByCreatedDesc(models)

	// Display models
	if *detail {
		DisplayModelsDetailed(models)
	} else {
		DisplayModels(models)
	}
}

func providersCommand() {
	fs := flag.NewFlagSet("providers", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: llmls providers\n\n")
		fmt.Fprintf(os.Stderr, "List all provider names.\n")
	}

	fs.Parse(os.Args[2:])

	// providers subcommand does not accept arguments
	if fs.NArg() > 0 {
		fmt.Fprintf(os.Stderr, "Error: providers subcommand does not accept arguments\n")
		fmt.Fprintf(os.Stderr, "Use 'llmls providers | grep pattern' to filter\n\n")
		fs.Usage()
		os.Exit(1)
	}

	// Fetch models from OpenRouter
	models, err := FetchModels()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Display all providers
	DisplayProviders(models)
}

