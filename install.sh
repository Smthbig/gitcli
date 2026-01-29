#!/usr/bin/env bash
set -e

APP_NAME="git-genius"
BIN_NAME="git-genius"
REPO_URL="https://github.com/Smthbig/gitcli.git"
PREBUILT_BASE="https://raw.githubusercontent.com/Smthbig/gitcli/main"
VERSION_FILE="${PREBUILT_BASE}/VERSION"

GREEN="\033[0;32m"
RED="\033[0;31m"
YELLOW="\033[1;33m"
NC="\033[0m"

log()  { echo -e "${GREEN}✔ $1${NC}"; }
warn() { echo -e "${YELLOW}⚠ $1${NC}"; }
err()  { echo -e "${RED}✖ $1${NC}"; exit 1; }

echo "======================================="
echo " Git Genius Installer"
echo "======================================="
echo

# --------------------------------------------------
# Safety (AndroidIDE / broken cwd fix)
# --------------------------------------------------
cd "$HOME" || err "HOME directory not accessible"

# --------------------------------------------------
# Detect OS / ARCH
# --------------------------------------------------
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

log "Detected OS: $OS"
log "Detected ARCH: $ARCH"

# --------------------------------------------------
# Detect environment
# --------------------------------------------------
ENV="unknown"

if [[ -n "$PREFIX" && "$PREFIX" == *"com.termux"* ]]; then
  ENV="termux"
elif [[ "$HOME" == *"androidide"* ]]; then
  ENV="androidide"
elif [[ "$OS" == "darwin" ]]; then
  ENV="macos"
elif [[ "$OS" == "linux" ]]; then
  if command -v sudo >/dev/null 2>&1; then
    ENV="linux"
  else
    ENV="restricted-linux"
  fi
fi

if [[ "$ENV" == "unknown" ]]; then
  warn "Unable to auto-detect environment"
  echo "1) Linux (sudo)"
  echo "2) Restricted Linux"
  echo "3) Termux"
  echo "4) AndroidIDE"
  echo "5) macOS"
  read -rp "Select [1-5]: " choice
  case "$choice" in
    1) ENV="linux" ;;
    2) ENV="restricted-linux" ;;
    3) ENV="termux" ;;
    4) ENV="androidide" ;;
    5) ENV="macos" ;;
    *) err "Invalid selection" ;;
  esac
fi

log "Environment: $ENV"

# --------------------------------------------------
# Install directory
# --------------------------------------------------
BIN_DIR="$HOME/bin"
[[ "$ENV" == "linux" || "$ENV" == "macos" ]] && BIN_DIR="/usr/local/bin"
mkdir -p "$BIN_DIR"

BIN_PATH="$BIN_DIR/$BIN_NAME"

# --------------------------------------------------
# Version check (🔥 KEY FEATURE 🔥)
# --------------------------------------------------
REMOTE_VERSION="$(curl -fsSL "$VERSION_FILE" 2>/dev/null || true)"

if [[ -x "$BIN_PATH" ]]; then
  INSTALLED_VERSION="$("$BIN_PATH" --version 2>/dev/null || true)"

  if [[ -n "$REMOTE_VERSION" && "$INSTALLED_VERSION" == *"$REMOTE_VERSION"* ]]; then
    log "Git Genius already installed ($INSTALLED_VERSION)"
    log "No update required"
    exit 0
  fi

  warn "Installed version differs, upgrading…"
fi

# --------------------------------------------------
# Ensure Git
# --------------------------------------------------
if ! command -v git >/dev/null 2>&1; then
  warn "Git not found"
  case "$ENV" in
    termux) pkg install -y git ;;
    linux) sudo apt update && sudo apt install -y git ;;
    macos) brew install git || xcode-select --install ;;
    *) err "Please install git manually" ;;
  esac
else
  log "Git already installed"
fi

# --------------------------------------------------
# Prebuilt binary (preferred)
# --------------------------------------------------
USE_PREBUILT=false

if [[ "$OS" == "linux" && ( "$ARCH" == "aarch64" || "$ARCH" == "arm64" || "$ARCH" == "x86_64" ) ]]; then
  USE_PREBUILT=true
fi

if [[ "$USE_PREBUILT" == true ]]; then
  PREBUILT_URL="${PREBUILT_BASE}/${BIN_NAME}"

  log "Installing prebuilt Git Genius binary"
  curl -fsSL "$PREBUILT_URL" -o "$BIN_PATH" || err "Download failed"
  chmod +x "$BIN_PATH"

  if ! echo "$PATH" | grep -q "$BIN_DIR"; then
    warn "Adding $BIN_DIR to PATH"
    echo "export PATH=\"$BIN_DIR:\$PATH\"" >> "$HOME/.profile"
    echo "export PATH=\"$BIN_DIR:\$PATH\"" >> "$HOME/.bashrc" 2>/dev/null || true
  fi

  log "Git Genius installed successfully 🎉"
  echo "Run:"
  echo "  git-genius"
  echo "Next:"
  echo "  Tools → Setup / Reconfigure"
  exit 0
fi

# --------------------------------------------------
# Fallback: source build (rare)
# --------------------------------------------------
warn "No prebuilt binary available, building from source"

if ! command -v go >/dev/null 2>&1; then
  err "Go not found. Please install Go manually."
fi

SRC_DIR="$HOME/.git-genius-src"

if [[ -d "$SRC_DIR/.git" ]]; then
  cd "$SRC_DIR"
  git pull --rebase
else
  rm -rf "$SRC_DIR"
  git clone "$REPO_URL" "$SRC_DIR"
  cd "$SRC_DIR"
fi

go build -o "$BIN_PATH" ./cmd/genius || err "Build failed"
chmod +x "$BIN_PATH"

log "Git Genius installed successfully 🎉"
echo "Run:"
echo "  git-genius"