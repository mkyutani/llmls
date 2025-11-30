# llmls

List LLM models from various providers (OpenAI, Anthropic, etc.)

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

```bash
# List available models
go run main.go

# Build the binary
go build -o llmls
```

## Contributing

1. Create a feature branch from `main`
2. Make your changes following the guidelines in [CLAUDE.md](CLAUDE.md)
3. Submit a pull request

## License

[Add your license here]
