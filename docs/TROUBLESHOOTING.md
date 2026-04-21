# Troubleshooting Guide

## Push keeps asking for username and token

Likely cause:

- Git has no credential helper configured

Fix:

1. Open `Tools -> Git Auth / Credential Helper`
2. Choose `Memory cache` or `Persistent store`
3. Preload the current GitHub token into Git

Git Genius now repeats this guidance automatically after failed HTTPS pushes.

If the issue continues:

- run `Tools -> Doctor`

## Push fails even though the token is saved

Remember:

- GitHub API auth and Git HTTPS auth are separate

Git Genius may be able to create/check repos with the token, while `git push` still fails if Git itself has no working credential helper.

## Pull fails because of local changes

Fix:

- use `Daily Git Operations -> Smart Pull`

Why:

- Smart Pull stashes local changes
- pulls the latest remote changes
- restores the local changes after pull

## Remote not found or wrong remote URL

Fix options:

- `Branch / Remote -> Configure remote`
- `Tools -> Create / Link GitHub Repository`

Use the first option when you only need to repair a local remote.

## Branch feels wrong after switching repos

Fix:

1. Run `Tools -> Doctor`
2. Use `Tools -> Switch Project / Repo` so Git Genius loads the target repo's own config
3. If branch mismatch is reported, re-run `Tools -> Setup / Reconfigure`

Git Genius keeps the configured branch and the real Git branch in sync, but older repos may need one cleanup pass.

## Switching repos asks for setup every time

Use `Tools -> Switch Project / Repo`.

That flow now:

- changes the active directory
- loads the selected repo's own `.git/.genius/config.json` when it exists
- avoids copying branch/remote settings from the previous repo
- can immediately switch branch or remote in the new repo

## No Git repository found

Fix:

- re-run `Tools -> Setup / Reconfigure`
- or `Tools -> Switch Project / Repo`

Git Genius can initialize the directory as a Git repo during setup.

## First push does not happen

Possible causes:

- the repo still has no commit
- there are no files to include in the first commit
- the branch exists locally but has never been published

Fix:

1. Create at least one file if the repo is empty
2. Run `Daily Git Operations -> Push changes`
3. If auth fails, configure `Tools -> Git Auth / Credential Helper`

## GitHub repo checks fail

Possible causes:

- missing token
- expired token
- insufficient GitHub permissions
- temporary network failure

Fix:

1. Check `Tools -> Doctor`
2. Re-run setup and validate the token
3. If needed, create the GitHub repo manually and then relink it

## Limited Mode

If Git is unavailable in a restricted environment, Git Genius enters limited mode.

Available in limited mode:

- setup
- change project directory
- Doctor
- Help

Unavailable until Git is installed:

- push
- pull
- branch switching
- stash and undo
