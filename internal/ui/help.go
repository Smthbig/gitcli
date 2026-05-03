package ui

// ============================================================
// Main Menu Help
// ============================================================

var HelpMain = []string{
	"1) Daily Git Operations",
	"   - Status, pull, smart pull, fetch, push",
	"",
	"2) Branch / Remote",
	"   - Switch existing branches, create branches, configure remotes",
	"",
	"3) Stash & Undo",
	"   - Temporarily save work or undo commits safely",
	"",
	"4) Tools",
	"   - Setup, repo linking, Git auth, project switching, Doctor",
	"",
	"h / help / ?",
	"   - Show this help screen",
	"",
	"6) Exit",
	"   - Quit Git Genius",
}

// ============================================================
// Daily Git Operations Help
// ============================================================

var HelpDaily = []string{
	"Push",
	"- Stages files, commits, and pushes to GitHub",
	"- If there are no file changes, Git Genius still attempts to push local commits",
	"- If the repo has no commits or no files yet, Git Genius explains the first-push next step",
	"- If push auth fails, configure a Git credential helper in Tools",
	"",
	"Pull",
	"- Pull latest changes from remote branch",
	"- Shows a short summary after pull completes",
	"- If local edits exist, Git Genius warns before a normal pull",
	"",
	"Smart Pull",
	"- Auto stashes local changes",
	"- Pulls latest code",
	"- Restores your changes safely",
	"- Recommended when you edited files locally and want the latest remote changes",
	"",
	"Fetch",
	"- Downloads remote changes without merging",
	"",
	"Status",
	"- Shows modified, staged, and untracked files",
	"",
	"Recommended order",
	"- Status",
	"- Smart Pull",
	"- edit files",
	"- Push changes",
}

// ============================================================
// Branch / Remote Help
// ============================================================

var HelpBranch = []string{
	"Switch to Existing Branch",
	"- Choose from local branches without resetting branch history",
	"- Automatically updates config branch",
	"",
	"Create New Branch",
	"- Creates a branch only when it does not already exist",
	"",
	"Configure Remote",
	"- Select an existing remote or safely add/update one",
	"- Existing remotes are updated with set-url instead of remove/add",
	"- Use this when push says the remote is missing or wrong",
}

// ============================================================
// Stash & Undo Help
// ============================================================

var HelpStash = []string{
	"Stash Changes",
	"- Save uncommitted work temporarily",
	"- Clean working directory",
	"",
	"Stash List",
	"- View all saved stashes",
	"",
	"Stash Pop",
	"- Restore last stashed changes",
	"- If conflicts appear, fix them manually and continue normally",
	"",
	"Undo Last Commit",
	"- Undo commit but KEEP file changes",
	"- Safe and reversible",
}

// ============================================================
// Tools Help
// ============================================================

var HelpTools = []string{
	"Setup / Reconfigure",
	"- Full guided setup (recommended first step)",
	"- Select project folder",
	"- Initialize git",
	"- Configure branch, remote, GitHub",
	"",
	"Switch Project",
	"- Switch to another project folder",
	"- Detects if the folder is a Git repo and shows context before switching",
	"- Loads that repo's saved branch and remote config when available",
	"- Can immediately switch branch or remote after changing directory",
	"",
	"Create / Link GitHub Repository",
	"- Create repo on GitHub if missing",
	"- Link local project to GitHub",
	"- Works even if you only want to repair the local remote URL",
	"",
	"Git Auth / Credential Helper",
	"- Configures Git's HTTPS credential helper",
	"- Can preload your current GitHub token into Git",
	"- Best fix when push keeps asking for username and token",
	"",
	"Doctor",
	"- Checks git, branch, remote, token, repo",
	"- Suggests fixes if something is wrong",
}

// ============================================================
// GitHub Help (NEW – very important for beginners)
// ============================================================

var HelpGitHub = []string{
	"What is GitHub?",
	"- Online platform to store and collaborate on code",
	"",
	"GitHub Token",
	"- Used for authentication (instead of password)",
	"- Create at: https://github.com/settings/tokens",
	"- Required scope: repo",
	"- Automation option: GIT_GENIUS_GITHUB_TOKEN",
	"",
	"GitHub Repository",
	"- Online copy of your project",
	"- Git Genius can create it automatically",
	"",
	"Remote",
	"- Link between local git and GitHub repo",
	"- Usually named: origin",
}

var HelpWorkflow = []string{
	"First-Time Setup",
	"- Run Tools -> Setup / Reconfigure",
	"- Pick project directory",
	"- Initialize git if needed",
	"- Configure branch, remote, GitHub repo, token, and auth helper",
	"- Brand-new repos are called out explicitly in the context panel",
	"",
	"Daily Workflow",
	"- Review context panel suggestions",
	"- Check ahead/behind status when Git Genius can compute it from local refs",
	"- Use Status to inspect local changes",
	"- Use Smart Pull before push when remote changes may exist",
	"- Push with a clear commit message",
	"",
	"Multi-Project Workflow",
	"- Use Tools -> Switch Project",
	"- Re-run Setup when switching to a brand-new repo",
}

var HelpTroubleshooting = []string{
	"Push asks for username/token repeatedly",
	"- Run Tools -> Git Auth / Credential Helper",
	"- Then preload the current GitHub token into Git",
	"",
	"Pull fails with local changes",
	"- Use Smart Pull",
	"- Or stash changes first",
	"",
	"Remote missing or wrong",
	"- Use Branch / Remote -> Configure remote",
	"- Or Tools -> Create / Link GitHub Repository",
	"",
	"Project feels misconfigured",
	"- Run Tools -> Doctor",
	"- Re-run Tools -> Setup / Reconfigure if needed",
	"",
	"Switching between repos should not require setup every time",
	"- Use Tools -> Switch Project",
	"- Each repo keeps its own local .git/.genius/config.json",
}
