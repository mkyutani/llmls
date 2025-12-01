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

## Git Communication Language

**All git-related communications must be in English:**
- Commit messages
- Pull request titles and descriptions
- Issue titles and descriptions
- Code comments and documentation
- Release notes

This ensures consistency and accessibility for all contributors and users in version control history.

**Note:** Other communications (chat, discussion, etc.) can be in any language.

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

**For solo development:**
- You may commit directly to `main` for small changes
- Use feature branches for larger changes to maintain clean history
- PRs are optional but can be useful for documenting major changes

**For team development:**
- Always create pull requests for review before merging to `main`

### GitHub Issue Workflow

Follow this workflow when working on GitHub issues:

#### 1. Issue Selection and Assignment
```bash
# List all open issues
gh issue list

# View specific issue details
gh issue view <number>

# Assign issue to yourself
gh issue edit <number> --add-assignee @me
```

#### 2. Design Discussion (for non-trivial issues)
- Discuss design approach in GitHub issue comments
- For complex issues, document:
  - Proposed approach
  - Key architectural decisions
  - Alternative options considered
  - Testing strategy
- Get feedback/approval in issue comments before implementation

#### 3. Branch Creation
- **Naming convention**:
  - Features: `feature/issue-<number>-brief-description`
  - Bug fixes: `fix/issue-<number>-brief-description`
  - Example: `feature/issue-1-cli-interface`
```bash
git checkout -b feature/issue-<number>-description
```

#### 4. Implementation
- Follow the design discussed in issue comments
- Make atomic commits with proper format
- Reference issue in commits when relevant: `feat: add CLI interface for issue #1`
- Follow code standards (see below)

#### 5. Testing
- Perform manual testing
- Add unit tests if applicable
- Verify issue requirements are met

#### 6. Merge to Main

**Option A: Direct merge (solo development)**
```bash
# Switch to main
git checkout main

# Merge feature branch
git merge feature/issue-<number>-description

# Delete feature branch
git branch -d feature/issue-<number>-description

# Push to remote (requires user confirmation)
git push

# Close issue (requires user confirmation)
gh issue close <number> --comment "Implemented in commits [commit-hash]"
```

**Option B: Pull Request (team development or major changes)**
```bash
# Push branch
git push -u origin feature/issue-<number>-description

# Create PR
gh pr create --title "feat: description for issue #<number>" \
  --body "Closes #<number>

## Summary
Brief description of changes.

## Changes
- Change 1
- Change 2

## Testing
- [x] Manual testing completed
- [x] Issue requirements verified"

# After review and merge
git checkout main
git pull
git branch -d feature/issue-<number>-description
```

### Standard Workflow Summary

1. Create a feature branch from `main`
2. Make your changes with properly formatted commits
3. Push your branch to the remote repository
4. Create a pull request to `main`
5. After review and approval, merge to `main`

### Release Workflow

This project uses GitHub Actions for automated release workflow.

#### Automated Release (Recommended)

The release workflow is automated via GitHub Actions. Simply push a tag to trigger:

```bash
# 1. Determine version number
# Review commits since last release
git log v<last-version>..HEAD --oneline

# Use semantic versioning (MAJOR.MINOR.PATCH):
# - MAJOR: Breaking changes
# - MINOR: New features (backwards compatible)
# - PATCH: Bug fixes (backwards compatible)

# 2. Create and push tag
git tag v<version>
git push origin v<version>
```

The GitHub Actions workflow (`.github/workflows/release.yml`) will automatically:
- Build binaries for all supported platforms (linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64)
- Create a GitHub release with the tag
- Generate release notes from commits
- Upload all built binaries to the release

#### Manual Release (Alternative)

If you need to create a release manually:

**1. Determine Version Number**
- Review commits since last release: `git log v<last-version>..HEAD --oneline`
- Use semantic versioning (MAJOR.MINOR.PATCH)

**2. Create Release with Notes**
```bash
# Create release with release notes
gh release create v<version> \
  --title "v<version>" \
  --notes "## New Features
- Feature 1
- Feature 2

## Bug Fixes
- Fix 1

## Documentation
- Doc update 1"
```

**3. Build and Upload Binaries**
```bash
# Create dist directory
mkdir -p dist

# Build for all platforms
GOOS=linux GOARCH=amd64 go build -o dist/llmls-linux-amd64 .
GOOS=linux GOARCH=arm64 go build -o dist/llmls-linux-arm64 .
GOOS=darwin GOARCH=amd64 go build -o dist/llmls-darwin-amd64 .
GOOS=darwin GOARCH=arm64 go build -o dist/llmls-darwin-arm64 .
GOOS=windows GOARCH=amd64 go build -o dist/llmls-windows-amd64.exe .

# Upload binaries to release
gh release upload v<version> \
  dist/llmls-linux-amd64 \
  dist/llmls-linux-arm64 \
  dist/llmls-darwin-amd64 \
  dist/llmls-darwin-arm64 \
  dist/llmls-windows-amd64.exe
```

**4. Verify Release**
- Check release page: `gh release view v<version> --web`
- Verify all binaries are attached
- Test download and execution of at least one binary

## Working with Claude Code

### Project Context

When asking Claude Code to make changes:
- Provide clear, specific requirements
- Reference existing files when relevant
- Ask for explanations if you don't understand suggested changes

### User Confirmation Requirements

Certain operations require explicit user confirmation before execution:
- **`git push`** - Pushing changes to remote repository is irreversible and affects shared code
- **`gh issue close`** - Closing issues should be reviewed to ensure implementation is complete and satisfactory

These restrictions are enforced in `.claude/settings.local.json` and help prevent unintended modifications to the remote repository and project tracking.

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