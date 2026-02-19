# 🧠 Git Genius

**Git without the "Grit."** Git Genius is a beginner-friendly, interactive CLI tool that helps you manage your repositories **without memorizing commands.**

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](http://makeapullrequest.com)

---

## 🚀 Why Git Genius?

Git is powerful, but the command-line workflow can be a maze of cryptic flags and accidental "detached HEAD" nightmares. Git Genius acts as your **intelligent navigator**, providing a menu-driven interface that handles the heavy lifting while you focus on your code.

* **Safe for Beginners:** Built-in guardrails and explicit confirmations.
* **Smart Workflows:** Automated stashing and recovery during pulls.
* **Built-in Doctor:** One-click health checks for your local environment.

---

## 🛠 Installation

Get up and running in seconds. This one-liner ensures a fresh installation by clearing previous versions.

```bash
# Uninstall old version & Install fresh
curl -fsSL https://raw.githubusercontent.com/Smthbig/gitcli/main/uninstall.sh | bash
curl -fsSL https://raw.githubusercontent.com/Smthbig/gitcli/main/install.sh | bash

# Launch the tool
git-genius
```

---

## ✨ Core Features

### 📂 Project & Repo Management
* **Smart Init:** Initialize Git if a repository doesn't exist.
* **Multi-Project Support:** Switch between project directories easily.
* **Safety First:** Destructive actions always require a confirmation.

### 🔄 Daily Operations (Simplified)
* **Push/Pull:** Handle remotes with clear, guided prompts.
* **Smart Pull:** Automatically stash -> pull -> pop to prevent merge conflicts on uncommitted work.
* **Branch/Remote Switching:** No more typing long branch names; just select from a list.

### 🩺 Git Doctor (Health Check)
The Doctor diagnostic tool checks:
- [x] Git installation & Internet connectivity.
- [x] Valid user.name and user.email configuration.
- [x] GitHub Token validation & API status.
- [x] Repository integrity and error log detection.

### ⏪ Recovery Features
* **Stash Manager:** Visually save, list, and restore stashes.
* **Undo Last Commit:** Safely revert your last commit while **keeping your changes** in the workspace.

---

## 🖥 Interface Preview

Git Genius transforms your terminal into a guided command center:

```text
Main Menu:
  [1] 📋 Status & Changes
  [2] 🚀 Push to Remote
  [3] 📥 Smart Pull (Auto-stash)
  [4] 🌿 Branch Manager
  [5] 🩺 Run Git Doctor (Health Check)
  [6] ⚙️  Setup/Configuration
  [Q] Exit
```

---

## 🗺 Roadmap (Future-Ready)

Git Genius is built with a modular architecture, ready for these upcoming features:
- [ ] History Viewer: A scrollable log of your recent commits.
- [ ] Diff Visualizer: See exactly what changed line-by-line.
- [ ] GitHub API Integration: Create repositories directly from the CLI.
- [ ] CI/CD Helper: Quick-start templates for GitHub Actions.

---

## ⚖️ Disclaimer & License

**Use Responsibly.** Git Genius is a helper tool, not a replacement for Git knowledge. It wraps Git commands to make them approachable, but you are still the captain of your code. Always review actions before confirming.

Distributed under the **MIT License**.

---

### 💡 Final Note
If you understand Git better after using this tool, **it has done its job.**
