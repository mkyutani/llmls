package main

import (
	"flag"
	"fmt"
	"os"
)

const version = "1.0.0"

func main() {
	if len(os.Args) < 2 {
		// Default behavior: list all models
		listModelsCommand(nil)
		return
	}

	subcommand := os.Args[1]

	switch subcommand {
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

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "llmls - List and manage LLM models\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  llmls [flags] [pattern]\n")
		fmt.Fprintf(os.Stderr, "  llmls providers\n\n")
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		fmt.Fprintf(os.Stderr, "  pattern  Search by model ID using glob pattern or provider name\n")
		fmt.Fprintf(os.Stderr, "           Glob pattern: * (any sequence) and ? (single character)\n")
		fmt.Fprintf(os.Stderr, "           Provider name: exact match (case-insensitive)\n")
		fmt.Fprintf(os.Stderr, "           Examples: \"anthropic/*\", \"*gpt-4*\", \"cohere\"\n\n")
		fmt.Fprintf(os.Stderr, "Subcommands:\n")
		fmt.Fprintf(os.Stderr, "  providers  List all provider names\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		fmt.Fprintf(os.Stderr, "  --detail       Display detailed model information\n")
		fmt.Fprintf(os.Stderr, "  --ollama-host  Ollama server URL (default: $OLLAMA_HOST or http://localhost:11434)\n")
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  llmls                                List all models (OpenRouter + Ollama)\n")
		fmt.Fprintf(os.Stderr, "  llmls cohere                         List all Cohere models (provider exact match)\n")
		fmt.Fprintf(os.Stderr, "  llmls \"anthropic/*\"                   List Anthropic models (glob pattern)\n")
		fmt.Fprintf(os.Stderr, "  llmls \"ollama/*\"                      List Ollama models only\n")
		fmt.Fprintf(os.Stderr, "  llmls \"*gpt-4*\"                       Search for GPT-4 models\n")
		fmt.Fprintf(os.Stderr, "  llmls --detail \"*opus*\"               Detailed view of Opus models\n")
		fmt.Fprintf(os.Stderr, "  llmls --ollama-host http://remote:11434  Use remote Ollama server\n")
		fmt.Fprintf(os.Stderr, "  llmls providers                      List all providers\n")
		fmt.Fprintf(os.Stderr, "  llmls | grep vision                  Filter by description\n")
	}

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
		fmt.Fprintf(os.Stderr, "Use external tools like grep to filter:\n")
		fmt.Fprintf(os.Stderr, "  llmls providers | grep open\n")
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

