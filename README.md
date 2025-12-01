# llmls

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.23-blue.svg)](https://go.dev)

A command-line tool to list and explore LLM models from various providers via OpenRouter and Ollama.

## Features

- **Multiple Providers** - List models from OpenRouter and Ollama in a unified view
- **Browse Models** - List all available LLM models from OpenRouter and local Ollama instances
- **Glob Pattern Search** - Search by model ID using `*` and `?` wildcards
- **Provider List** - Quick access to all provider names
- **Detailed Information** - View comprehensive model details with `--detail` flag
- **Ollama Support** - Automatically includes local Ollama models with customizable server URL
- **Sorted Output** - Models are sorted by creation date (newest first)
- **Pipe-Friendly** - Designed to work seamlessly with Unix tools like `grep`, `awk`, and `sort`
- **Version Display** - Show version information with `--version` or `-v` flag

## Installation

### Option 1: Download Pre-built Binary (Recommended)

Download the latest release from the [Releases page](https://github.com/mkyutani/llmls/releases).

Pre-built binaries include the correct version information and are available for:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

### Option 2: Build from Source

```bash
git clone https://github.com/mkyutani/llmls.git
cd llmls
go build -o llmls
```

**Note:** Building with `go build` without version flags will show version as "dev". This is fine for development, but for production use, download the pre-built binary from the releases page.

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

Display version information:

```bash
llmls --version  # or -v
```

List all available models (OpenRouter + Ollama):

```bash
llmls
```

Search models by ID (glob pattern):

```bash
llmls "anthropic/*"      # All Anthropic models
llmls "ollama/*"         # All Ollama local models
llmls "*gpt-4*"          # All GPT-4 models
llmls "*opus*"           # Models with "opus" in ID
```

View detailed information:

```bash
llmls --detail
llmls --detail "*claude*"
llmls --detail "ollama/*"  # See Ollama model details (size, quantization, etc.)
```

Use custom Ollama server:

```bash
llmls --ollama-host http://192.168.1.100:11434  # Remote server
llmls --ollama-host http://ollama:11434         # Docker container
export OLLAMA_HOST=http://remote:11434          # Set via environment variable
llmls
```

List all providers:

```bash
llmls providers
```

### Filtering with External Tools

Since `llmls` follows Unix philosophy, use standard tools for advanced filtering:

```bash
# Filter providers by pattern
llmls providers | grep open

# Filter by description
llmls | grep vision
llmls "*claude*" | grep reasoning

# Combine filters
llmls "anthropic/*" | grep -i vision
llmls "ollama/*" | grep llama

# Count models per provider
llmls | awk '{print $2}' | sort | uniq -c

# Get only model IDs
llmls | awk '{print $1}'

# List only local Ollama models
llmls "ollama/*"
```

### Glob Pattern Syntax

Patterns support standard glob wildcards:
- `*` - Match any sequence of characters
- `?` - Match any single character

**Important:** Quote patterns to prevent shell expansion:

```bash
llmls "anthropic/*"      # Correct
llmls anthropic/*        # May expand in shell
```

### Shell Configuration for Easier Usage

To avoid quoting glob patterns, configure your shell to disable glob expansion for `llmls`:

**For Bash (add to `~/.bashrc`):**
```bash
llmls() {
    set -f
    command llmls "$@"
    set +f
}
```

**For Zsh (add to `~/.zshrc`):**
```bash
llmls() {
    setopt local_options noglob
    command llmls "$@"
}
```

After adding this configuration and restarting your shell, you can use glob patterns without quotes:

```bash
llmls anthropic/*
llmls *gpt-4*
```

### Ollama Configuration

By default, `llmls` attempts to connect to Ollama at `http://localhost:11434`. If Ollama is not available, it silently continues with OpenRouter models only.

**Configuration Priority:**
1. `--ollama-host` command-line flag (highest priority)
2. `OLLAMA_HOST` environment variable
3. Default: `http://localhost:11434`

**Examples:**

```bash
# Use default (localhost:11434)
llmls

# Specify via command-line flag
llmls --ollama-host http://192.168.1.100:11434

# Set via environment variable
export OLLAMA_HOST=http://ollama:11434
llmls

# Ollama unavailable - shows OpenRouter models only (no error)
llmls  # When Ollama is not running
```

**Ollama Model Details:**

When using `--detail` with Ollama models, additional information is displayed:
- **Model Family** - Model architecture family (e.g., llama, mistral)
- **Parameter Size** - Model size (e.g., 7B, 13B, 70B)
- **Quantization** - Quantization level (e.g., Q4_0, Q8_0)
- **Format** - Model format (e.g., gguf)
- **Model Size** - Disk size in GB

### Output Format

Models are displayed with the following columns:
- **Model ID** - Full model identifier
- **Provider** - Provider name extracted from model ID
- **Created** - Creation date in YYYY-MM-DD format (local timezone)
- **Description** - Model description (truncated to 98 characters)

Results are sorted by creation date in descending order (newest first).

Example output:
```
anthropic/claude-opus-4.5      anthropic  2025-11-24  Claude Opus 4.5 is Anthropic's frontier reasoning model optimized for complex software engineeri..
openai/gpt-4.1                 openai     2025-04-14  GPT-4.1 is a flagship large language model optimized for advanced instruction following, real-worl..
ollama/llama3.2:latest         ollama     2024-12-01  llama 3B (Q4_0) - 2.0 GB
google/gemini-3-pro-preview    google     2025-11-18  Gemini 3 Pro is Google's flagship frontier model for high-precision multimodal reasoning, combin..
```

### Help

Display usage information:

```bash
llmls --help
llmls providers --help
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
