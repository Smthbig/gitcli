# Git Genius (gitcli)

Beginner-friendly interactive Git assistant for daily workflows, setup, recovery, and GitHub linking.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](http://makeapullrequest.com)

## What It Solves

- Removes Git command memorization with guided terminal menus
- Gives safer defaults for common branch, remote, and recovery actions
- Helps first-time setup for local repos and GitHub remotes
- Supports multi-project switching with recent-directory history
- Includes built-in diagnostics with Doctor checks
- Shows cheap local ahead/behind context when remote refs are already available
- Gives clearer first-run, empty-repo, first-push, and push-auth guidance

## Core Features

- Daily Git operations: `status`, `pull`, `smart pull`, `fetch`, `push`
- Safe branch and remote management
- Stash and undo flows for recovery
- Guided setup and reconfiguration
- Context panel with branch, remote, ahead/behind, and first-push cues
- Fast project/repo switching that loads the selected repo's saved config
- GitHub repository create/link flow
- GitHub token support from local storage or `GIT_GENIUS_GITHUB_TOKEN`
- Git credential-helper configuration to reduce repeated HTTPS prompts
- Per-repo config plus global active-project state

## Install

### Quick install

```bash
curl -fsSL https://raw.githubusercontent.com/Smthbig/gitcli/main/install.sh | bash
```

Then run:

```bash
git-genius
```

### Reinstall

```bash
curl -fsSL https://raw.githubusercontent.com/Smthbig/gitcli/main/uninstall.sh | bash
curl -fsSL https://raw.githubusercontent.com/Smthbig/gitcli/main/install.sh | bash
```

### Update

```bash
curl -fsSL https://raw.githubusercontent.com/Smthbig/gitcli/main/install.sh | bash
```

### Uninstall

```bash
curl -fsSL https://raw.githubusercontent.com/Smthbig/gitcli/main/uninstall.sh | bash
```

## Commands

- Main command: `git-genius`
- Version: `git-genius --version`
- Project/repo name: `gitcli`

If you prefer `gitcli` as the command name:

```bash
echo "alias gitcli='git-genius'" >> ~/.zshrc
source ~/.zshrc
```

## Docs

- [Workflow Guide](docs/WORKFLOWS.md)
- [Authentication Guide](docs/AUTH.md)
- [Troubleshooting Guide](docs/TROUBLESHOOTING.md)
- [Roadmap](Roadmap.txt)

## Quick Start

### First-time setup

1. Run `git-genius`
2. Open `Tools -> Setup / Reconfigure`
3. Choose the project directory
4. Initialize Git if the folder is not already a repo
5. Configure branch, remote, GitHub owner/repo, token, and auth helper
6. Optionally create/link the GitHub repository
7. Optionally do the first commit and push

### Daily workflow

1. Review the context panel
2. Run `Daily Git Operations -> Git status`
3. Run `Smart Pull` before pushing when the remote may have changed
4. Edit files in your normal editor
5. Run `Push changes`

The context panel now calls out:

- first runs in brand-new repos
- `ahead / behind` counts when matching remote refs already exist locally
- first-push status when a branch has not been published yet
- direct credential-helper guidance after failed HTTPS pushes

### Multi-project workflow

1. Open `Tools -> Switch Project / Repo`
2. Pick a recent project or enter a path manually
3. Git Genius switches the active directory and loads that repo's saved branch/remote config
4. Optionally switch branch or remote immediately in the same flow
5. Re-run setup only for a brand-new repo that has never been configured

## Authentication Model

Git Genius uses two separate layers:

- GitHub API authentication
  - Used for repo existence checks and repo creation
  - Reads the token from local Git Genius storage or `GIT_GENIUS_GITHUB_TOKEN`
- Git HTTPS authentication
  - Used by `git push` and `git pull`
  - Managed through `Tools -> Git Auth / Credential Helper`

Recommended setup:

- Use `GIT_GENIUS_GITHUB_TOKEN` for automation or CI-like environments
- Configure the Git credential helper from the Tools menu
- Preload the current GitHub token into Git so HTTPS pushes stop prompting

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
3) Git Auth / Credential Helper
4) Switch Project / Repo
5) Doctor (health check)
6) Back
```

## Limited Mode

If Git is not available in a restricted environment, Git Genius enters limited mode. In that mode it keeps setup, project switching, Doctor, and Help available while hiding normal daily Git operations.

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

Test:

```bash
go test ./...
```

## Codebase Overview

- `cmd/genius/main.go`
  - app bootstrap and version entrypoint
- `internal/menu`
  - top-level navigation and section routing
- `internal/setup`
  - setup, reconfiguration, auth helper, repo linking, and project/repo switching
- `internal/gitops`
  - daily Git operations and safe branch/remote flows
- `internal/doctor`
  - diagnostics and health checks
- `internal/system`
  - Git command helpers and shared system behavior
- `internal/github`
  - token storage and GitHub API client logic
- `internal/config`
  - repo-local config plus global active-project state and history
- `internal/ui`
  - prompts, rendering, and help text

## Current Direction

The current priority is product hardening:

- smoother auth and credential-helper behavior
- stronger tests around interactive workflows
- richer help and docs
- clearer recovery and troubleshooting guidance
- future automation and release improvements

See [Roadmap.txt](Roadmap.txt) for the forward plan.

## License

MIT License.

## Note

Git Genius improves safety and usability, but it is still a helper layer on top of Git. Understanding basic Git concepts remains important.
