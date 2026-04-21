package ui

// ============================================================
// Main Menu Help
// ============================================================

var HelpMain = []string{
	"1) Daily Git Operations",
	"   - Push, pull, fetch, status (daily workflow)",
	"",
	"2) Branch / Remote",
	"   - Switch branches or change git remote",
	"",
	"3) Stash & Undo",
	"   - Temporarily save work or undo commits safely",
	"",
	"4) Tools",
	"   - Setup, GitHub repo linking, Doctor (health check)",
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
	"- First push will guide you if repo/remote is missing",
	"",
	"Pull",
	"- Pull latest changes from remote branch",
	"",
	"Smart Pull",
	"- Auto stashes local changes",
	"- Pulls latest code",
	"- Restores your changes safely",
	"",
	"Fetch",
	"- Downloads remote changes without merging",
	"",
	"Status",
	"- Shows modified, staged, and untracked files",
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
	"Create / Link GitHub Repository",
	"- Create repo on GitHub if missing",
	"- Link local project to GitHub",
	"- Works even if you only want to repair the local remote URL",
	"",
	"Change Project Directory",
	"- Switch to another project folder",
	"- Shows recent projects for faster switching",
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
