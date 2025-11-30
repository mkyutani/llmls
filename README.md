# llmls

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.23-blue.svg)](https://go.dev)

A command-line tool to list and explore LLM models from various providers via OpenRouter.

## Features

- **Browse Models** - List all available LLM models from OpenRouter
- **Filter by Provider** - Find models from specific providers (Anthropic, OpenAI, Google, etc.)
- **Filter by Model** - Search for specific models by name
- **Subcommands** - Quick access to providers and models list
- **Detailed Information** - View model IDs, providers, creation dates, and descriptions
- **Sorted Output** - Models are sorted by creation date (newest first)

## Installation

### Option 1: Download Pre-built Binary (Recommended)

Download the latest release from the [Releases page](https://github.com/yourusername/llmls/releases).

### Option 2: Install with Go

```bash
go install github.com/yourusername/llmls@latest
```

### Option 3: Build from Source

```bash
git clone https://github.com/yourusername/llmls.git
cd llmls
go build -o llmls
```

## Development Setup

### Prerequisites

- Docker Desktop (or Docker Engine)
- Visual Studio Code
- VS Code Dev Containers extension

### Getting Started

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd llmls
   ```

2. **Open in Dev Container:**
   - Open the project in VS Code
   - When prompted, click "Reopen in Container"
   - Or use Command Palette (Ctrl/Cmd+Shift+P): "Dev Containers: Reopen in Container"

3. **Configure environment variables:**
   - On first startup, `.env.local` is automatically created from `.env.local.example`
   - Edit `.env.local` and add your API keys:
     ```
     ANTHROPIC_API_KEY=your_actual_api_key_here
     ```

4. **Start developing:**
   - The dev container includes Go, GitHub CLI, Docker, and Claude Code extension
   - All dependencies are automatically installed

### Dev Container Features

The development environment includes:
- **Go 1.23** - Latest Go toolchain
- **GitHub CLI** - `gh` command for GitHub operations
- **Docker** - Docker-outside-of-docker for container operations
- **Claude Code** - AI-powered development assistance
- **Go Extensions** - Language server, linting, formatting tools

### Environment Variables

Configuration is managed through `.env.local`:
- Automatically created from `.env.local.example` on first run
- Git-ignored for security (never commit this file)
- Required variables:
  - `ANTHROPIC_API_KEY` - For Claude Code and Anthropic model listing

## Development Guidelines

See [CLAUDE.md](CLAUDE.md) for detailed development guidelines including:
- Git commit message format
- Branch strategy and workflow
- Code standards
- Security best practices

## Usage

### Basic Commands

List all available models from OpenRouter:

```bash
llmls
```

List all providers:

```bash
llmls providers
```

List all models with their providers:

```bash
llmls models
```

### Subcommands with Filters

Filter providers by name (partial match, case-insensitive):

```bash
llmls providers anthropic
llmls providers open
```

Filter models by name (partial match, case-insensitive):

```bash
llmls models claude
llmls models gpt-4
llmls models gemini
```

### Flag-based Filtering

Filter models by provider name:

```bash
llmls --provider anthropic
llmls --provider openai
llmls --provider google
```

Filter models by model name:

```bash
llmls --model gpt-4
llmls --model claude
llmls --model gemini
```

Combine filters:

```bash
llmls --provider google --model gemini
llmls --provider anthropic --model opus
```

### Output Format

Models are displayed with the following columns:
- **Model ID** - Full model identifier
- **Provider** - Provider name extracted from model ID
- **Created** - Creation date in YYYY-MM-DD format (local timezone)
- **Description** - Model description (truncated to 98 characters)

Results are sorted by creation date in descending order (newest first).

Example output:
```
anthropic/claude-opus-4.5      anthropic            2025-11-24  Claude Opus 4.5 is Anthropic's frontier reasoning model optimized for complex software engineeri..
openai/gpt-4.1                 openai               2025-04-14  GPT-4.1 is a flagship large language model optimized for advanced instruction following, real-worl..
google/gemini-3-pro-preview    google               2025-11-18  Gemini 3 Pro is Google's flagship frontier model for high-precision multimodal reasoning, combin..
```

### Help

Display usage information:

```bash
llmls --help
llmls providers --help
llmls models --help
```

## Troubleshooting

### Error: "Error fetching models"

This error typically occurs when:
- **Network issues** - Check your internet connection
- **OpenRouter API issues** - The OpenRouter API may be temporarily unavailable
- **Timeout** - The request may have timed out. Try again in a few moments.

### No output or empty results

If you see no output:
- **Check filters** - Your filter criteria may be too restrictive. Try without filters first.
- **API response** - The OpenRouter API may have returned no models. This is unusual but possible.

### Command not found

If you get "command not found":
- **Path issue** - Ensure the `llmls` binary is in your PATH or use `./llmls` if running from the current directory
- **Installation** - Verify the installation completed successfully with `which llmls`

### Permission denied

If you get "permission denied":
```bash
chmod +x llmls
```

## Contributing

1. Create a feature branch from `main`
2. Make your changes following the guidelines in [CLAUDE.md](CLAUDE.md)
3. Submit a pull request

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

Copyright (c) 2025 Miki Yutani
