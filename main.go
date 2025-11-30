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
	providerFilter := fs.String("provider", "", "Filter models by provider name (partial match)")
	modelFilter := fs.String("model", "", "Filter models by model name (partial match)")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "llmls - List and manage LLM models\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  llmls [flags]\n")
		fmt.Fprintf(os.Stderr, "  llmls providers [filter]\n")
		fmt.Fprintf(os.Stderr, "  llmls models [filter]\n\n")
		fmt.Fprintf(os.Stderr, "Subcommands:\n")
		fmt.Fprintf(os.Stderr, "  providers [filter]  List provider names (optionally filtered)\n")
		fmt.Fprintf(os.Stderr, "  models [filter]     List models with provider names (optionally filtered)\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		fs.PrintDefaults()
	}

	if args != nil {
		fs.Parse(args)
	} else {
		fs.Parse(os.Args[1:])
	}

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

func providersCommand() {
	fs := flag.NewFlagSet("providers", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: llmls providers [filter]\n\n")
		fmt.Fprintf(os.Stderr, "List provider names, optionally filtered by a search string.\n")
		fmt.Fprintf(os.Stderr, "Filter performs case-insensitive partial matching.\n")
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
		fmt.Fprintf(os.Stderr, "List models with provider names, optionally filtered by a search string.\n")
		fmt.Fprintf(os.Stderr, "Filter performs case-insensitive partial matching on model names.\n")
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
