# Git Genius (gitcli)

Beginner-friendly interactive Git assistant for daily workflows, setup, and recovery.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](http://makeapullrequest.com)

## What It Solves

- Removes command memorization with guided terminal menus
- Handles common Git mistakes with safer defaults and confirmations
- Helps first-time setup for repo, branch, remote, and GitHub token
- Includes health diagnostics via built-in Doctor checks

## Features

- **Daily Git Operations:** push, pull, smart pull, fetch, status
- **Branch and Remote Management:** switch existing branches, create branches, and update remotes safely
- **Recovery Tools:** stash save/list/pop and undo last commit (keep changes)
- **Guided Setup:** configure work dir, repo, branch, remote, identity, token
- **GitHub Linking:** create or link remote repository
- **Automation-Friendly Auth:** supports `GIT_GENIUS_GITHUB_TOKEN`
- **Doctor:** validates git install, project state, token, remote, repo status

## Installation

### Quick install (recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/Smthbig/gitcli/main/install.sh | bash
```

Then run:

```bash
git-genius
```

### Fresh reinstall

```bash
curl -fsSL https://raw.githubusercontent.com/Smthbig/gitcli/main/uninstall.sh | bash
curl -fsSL https://raw.githubusercontent.com/Smthbig/gitcli/main/install.sh | bash
```

### Update to latest version

The installer is update-aware and exits when already up to date.

```bash
curl -fsSL https://raw.githubusercontent.com/Smthbig/gitcli/main/install.sh | bash
```

### Uninstall

```bash
curl -fsSL https://raw.githubusercontent.com/Smthbig/gitcli/main/uninstall.sh | bash
```

## Commands

- Primary command: `git-genius`
- Version command: `git-genius --version`
- Project/repo name: `gitcli` (this repository)

If you prefer `gitcli` as command, add an alias:

```bash
echo "alias gitcli='git-genius'" >> ~/.zshrc
source ~/.zshrc
```

## In-App "Update" Actions

Git Genius has two update paths inside the app:

- **Tools -> Setup / Reconfigure**
  - Re-runs complete setup and updates branch/remote/token/workdir settings
- **Tools -> Create / Link GitHub Repository**
  - Creates or relinks GitHub repository and updates remote URL

## Interface Map

```text
Main Menu
1) Daily Git Operations
2) Branch / Remote
3) Stash & Undo
4) Tools
5) Help / About
6) Exit
```

```text
Daily Git Operations
1) Push changes (commit + push)
2) Pull changes
3) Smart Pull (auto-stash + pull)
4) Fetch all remotes
5) Git status
6) Back
```

```text
Branch / Remote
1) Switch to existing branch
2) Create new branch
3) Configure remote
4) Back
```

```text
Tools
1) Setup / Reconfigure
2) Create / Link GitHub Repository
3) Change Project Directory
4) Doctor (health check)
5) Back
```

## Build From Source

Requirements:

- Go 1.21+
- Git installed

Build:

```bash
go fmt ./...
VERSION="$(cat VERSION)"
CGO_ENABLED=0 go build -ldflags "-X main.version=${VERSION}" -o git-genius ./cmd/genius
./git-genius
```

## Project Design (Codebase Overview)

Your project is modular and cleanly separated by responsibility:

- `cmd/genius/main.go`
  - Entry point, runtime safety env setup, app bootstrap
- `internal/menu`
  - Top-level interactive navigation and section routing
- `internal/setup`
  - First-time setup, reconfiguration, directory changes, GitHub repo linking
- `internal/gitops`
  - Daily operations: push/pull/fetch/smart pull/branch/remote/stash/undo
- `internal/doctor`
  - System and repository diagnostics
- `internal/system`
  - Git command execution, environment checks, runtime helpers
- `internal/github`
  - Token persistence and GitHub API interactions
- `internal/config`
  - Per-repo config at `<repo>/.git/.genius/config.json`
  - Global active-project state at `~/.git-genius/state.json` for multi-repo switching
- `internal/ui`
  - Terminal prompts, confirmations, rendering, help strings

Design strengths:

- Clear separation between UI flow and Git/system logic
- Safe defaults for first-run and restricted environments
- Config-driven operations across multiple working directories
- Good extensibility for future sections/features

## Roadmap Ideas

- Commit history viewer and better log UX
- Rich diff summaries for staged/unstaged changes
- Non-interactive flags for scripting/automation mode
- Optional shell autocompletion for command shortcuts

## License

MIT License.

## Note

Git Genius is a helper layer on top of Git. It improves flow and safety, but understanding core Git concepts is still important.
