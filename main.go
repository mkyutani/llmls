package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		// Default behavior: list all models
		listModelsCommand(nil)
		return
	}

	subcommand := os.Args[1]

	switch subcommand {
	case "providers":
		providersCommand()
	case "models":
		modelsCommand()
	default:
		// If not a subcommand, treat as old flag-based behavior
		listModelsCommand(os.Args[1:])
	}
}

func listModelsCommand(args []string) {
	fs := flag.NewFlagSet("llmls", flag.ExitOnError)
	providerFilter := fs.String("provider", "", "Filter models by provider name (glob pattern: * and ?)")
	modelFilter := fs.String("model", "", "Filter models by model name (glob pattern: * and ?)")
	descriptionFilter := fs.String("description", "", "Filter models by description text (glob pattern: * and ?)")
	detail := fs.Bool("detail", false, "Display detailed model information")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "llmls - List and manage LLM models\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  llmls [flags] [search-term]\n")
		fmt.Fprintf(os.Stderr, "  llmls providers [filter]\n")
		fmt.Fprintf(os.Stderr, "  llmls models [filter]\n\n")
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		fmt.Fprintf(os.Stderr, "  search-term  Search across model ID, provider, and description using glob pattern\n")
		fmt.Fprintf(os.Stderr, "               Supports * (any sequence) and ? (single character)\n")
		fmt.Fprintf(os.Stderr, "               Ignored if flags are used\n\n")
		fmt.Fprintf(os.Stderr, "Subcommands:\n")
		fmt.Fprintf(os.Stderr, "  providers [filter]  List provider names (optionally filtered with glob pattern)\n")
		fmt.Fprintf(os.Stderr, "  models [filter]     List models with provider names (optionally filtered with glob pattern)\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		fs.PrintDefaults()
	}

	if args != nil {
		fs.Parse(args)
	} else {
		fs.Parse(os.Args[1:])
	}

	// Get search term from positional argument
	searchTerm := ""
	if fs.NArg() > 0 {
		searchTerm = fs.Arg(0)
	}

	// Fetch models from OpenRouter
	models, err := FetchModels()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Filter models
	models = FilterModels(models, *providerFilter, *modelFilter, *descriptionFilter, searchTerm)

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
		fmt.Fprintf(os.Stderr, "Usage: llmls providers [filter]\n\n")
		fmt.Fprintf(os.Stderr, "List provider names, optionally filtered using glob pattern.\n")
		fmt.Fprintf(os.Stderr, "Supports * (any sequence) and ? (single character).\n")
		fmt.Fprintf(os.Stderr, "Example: llmls providers \"open*\"\n")
	}

	fs.Parse(os.Args[2:])

	filter := ""
	if fs.NArg() > 0 {
		filter = fs.Arg(0)
	}

	// Fetch models from OpenRouter
	models, err := FetchModels()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Display providers
	DisplayProviders(models, filter)
}

func modelsCommand() {
	fs := flag.NewFlagSet("models", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: llmls models [filter]\n\n")
		fmt.Fprintf(os.Stderr, "List models with provider names, optionally filtered using glob pattern.\n")
		fmt.Fprintf(os.Stderr, "Supports * (any sequence) and ? (single character).\n")
		fmt.Fprintf(os.Stderr, "Example: llmls models \"*gpt-4*\"\n")
	}

	fs.Parse(os.Args[2:])

	filter := ""
	if fs.NArg() > 0 {
		filter = fs.Arg(0)
	}

	// Fetch models from OpenRouter
	models, err := FetchModels()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Filter and display models
	DisplayModelsFiltered(models, filter)
}
