# Authentication Guide

Git Genius uses two authentication layers. Keeping them separate is intentional.

## 1. GitHub API Authentication

Used for:

- checking whether a GitHub repository exists
- creating a repository on GitHub
- validating the token and resolving the authenticated username

Sources:

- local Git Genius token storage at `~/.git-genius/token`
- environment variable `GIT_GENIUS_GITHUB_TOKEN`

Environment-variable usage:

```bash
export GIT_GENIUS_GITHUB_TOKEN=your_token_here
git-genius
```

Recommended token scope:

- `repo`

## 2. Git HTTPS Authentication

Used for:

- `git push`
- `git pull`
- any Git operation that talks to GitHub over HTTPS

This is separate from the GitHub API token storage above. Git itself still needs a credential helper if you want to avoid repeated prompts.

Configure it from:

- `Tools -> Git Auth / Credential Helper`

Available modes:

- Memory cache
  - recommended default
  - keeps credentials for a limited time in memory
  - does not write plaintext credentials to disk
- Persistent store
  - saves credentials to `~/.git-credentials`
  - survives restarts
  - less secure because credentials are stored in plain text

## Preloading Credentials

After configuring a credential helper, Git Genius can preload the current GitHub token into Git.

This helps future HTTPS pushes stop prompting for:

- GitHub username
- GitHub token/password

## Recommended Setup

For most local machines:

1. Configure a GitHub token in Setup
2. Open `Tools -> Git Auth / Credential Helper`
3. Choose `Memory cache`
4. Preload the GitHub credentials into Git

For automation-heavy environments:

1. Export `GIT_GENIUS_GITHUB_TOKEN`
2. Run Git Genius
3. Optionally save the env token locally if you want interactive reuse

## Security Notes

- Environment variables are good for automation and ephemeral sessions
- Memory cache is safer than persistent plain-text storage
- Persistent store is convenient but writes credentials to disk
- Git Genius never embeds the token in the remote URL
