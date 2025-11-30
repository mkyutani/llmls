# Claude Code Development Guidelines

This document provides guidelines for developing this project with Claude Code assistance.

## Environment Setup

### Dev Container

This project uses VS Code Dev Containers for consistent development environments.

**Container specifications:**
- Name: `llmls_dev` (repository name + "_dev")
- Base: Go 1.23 on Debian Bookworm
- Features:
  - Go development environment
  - GitHub CLI (`gh`)
  - Docker (outside-of-docker)
- Extensions:
  - Go language support
  - Claude Code

**Starting the environment:**
1. Ensure Docker is running on your host machine
2. Open the project in VS Code
3. When prompted, click "Reopen in Container" (or use Command Palette: "Dev Containers: Reopen in Container")
4. The container will build and start automatically

### Environment Variables

The project uses `.env.local` for environment-specific configuration.

**Setup:**
1. On first container startup, `.env.local.example` is automatically copied to `.env.local`
2. Edit `.env.local` and add your API keys:
   ```
   ANTHROPIC_API_KEY=your_actual_api_key_here
   ```
3. The `.env.local` file is automatically loaded into the container environment
4. This file is git-ignored for security

**Important:** Never commit `.env.local` or actual API keys to the repository.

## Git Workflow

### Commit Message Format

All commit messages must follow this format:

```
<prefix>: <description in 20 words or less>
```

**Requirements:**
- Language: English only
- Length: Maximum 20 words
- Format: Single line, single sentence
- Prefix: Use one of the standard prefixes below
- No footers: Do NOT include "Generated with Claude Code" or "Co-Authored-By" footers

**Standard Prefixes:**
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `style:` - Code style/formatting (no functional changes)
- `refactor:` - Code refactoring (no functional changes)
- `test:` - Adding or updating tests
- `chore:` - Maintenance tasks, dependencies, build config
- `perf:` - Performance improvements

**Examples:**
```
feat: add OpenAI model listing support
fix: resolve timeout error in Anthropic API calls
docs: update installation instructions for Windows
refactor: extract model formatting to separate function
```

### Branch Strategy

- `main` - Production-ready code
- Feature branches: `feature/description`
- Bug fixes: `fix/description`
- Always create pull requests for review before merging to `main`

### Workflow

1. Create a feature branch from `main`
2. Make your changes with properly formatted commits
3. Push your branch to the remote repository
4. Create a pull request to `main`
5. After review and approval, merge to `main`

## Working with Claude Code

### Project Context

When asking Claude Code to make changes:
- Provide clear, specific requirements
- Reference existing files when relevant
- Ask for explanations if you don't understand suggested changes

### Code Standards

- Follow Go best practices and idioms
- Write tests for new functionality
- Keep functions focused and single-purpose
- Use meaningful variable and function names
- Add comments for complex logic only (code should be self-documenting)

### Security

- Never commit sensitive data (API keys, credentials, secrets)
- Review all code changes before committing
- Use `.env.local` for local configuration
- Keep dependencies updated

## Additional Resources

- [Go Documentation](https://go.dev/doc/)
- [Claude Code Documentation](https://github.com/anthropics/claude-code)
- [Dev Containers Documentation](https://code.visualstudio.com/docs/devcontainers/containers)