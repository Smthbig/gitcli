# Workflow Guide

## First-Time Setup

Use this when the current directory is a new project or when the repo has not been configured in Git Genius yet.

1. Launch `git-genius`
2. Open `Tools -> Setup / Reconfigure`
3. Select the project directory
4. Initialize Git if needed
5. Configure:
   - branch
   - remote
   - GitHub owner and repo
   - GitHub token
   - Git credential helper
6. Optionally create the GitHub repository
7. Optionally do the first commit and push

Brand-new repos are now called out directly in the context panel so the first run is harder to miss.

## Recommended Daily Flow

1. Review the context panel
2. Check the ahead/behind line when it is shown
3. Run `Daily Git Operations -> Git status`
4. Run `Smart Pull` if you may be behind the remote and have local edits
5. Edit files outside Git Genius
6. Run `Push changes`

Why this order:

- `Status` helps you inspect what changed
- the context panel can show whether the current branch is ahead of or behind the fetched remote ref
- `Smart Pull` reduces pull failures caused by local edits
- `Push changes` stages, commits, and pushes in one guided step

## Branch Workflow

Use `Branch / Remote` for branch work:

- `Switch to existing branch`
  - selects from real local branches
  - does not force-reset branch history
- `Create new branch`
  - creates only when the branch does not already exist
  - switches to it immediately

## Remote Workflow

Use `Branch / Remote -> Configure remote` when:

- the repo points to the wrong GitHub URL
- the configured remote name is missing
- you need to add a new remote

Git Genius updates existing remotes with `set-url` instead of deleting and recreating them.

## Multi-Project Workflow

1. Open `Tools -> Change Project Directory`
2. Pick a recent project or enter a new path
3. If the new path is not a Git repo, initialize it or re-run setup

Git Genius stores:

- repo-local config in `.git/.genius/config.json`
- global active project and recent directories in `~/.git-genius/state.json`

## Recovery Workflow

Use `Stash & Undo` when something feels risky or broken:

- `Stash changes`
  - temporarily saves local edits
- `List stashes`
  - shows saved work
- `Apply last stash (pop)`
  - restores the most recent stash
- `Undo last commit`
  - removes the last commit while keeping file changes

If the repo still feels wrong after recovery, run `Tools -> Doctor`.
