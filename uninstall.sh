#!/usr/bin/env bash
set -e

APP_NAME="git-genius"
BIN_NAME="git-genius"

GREEN="\033[0;32m"
RED="\033[0;31m"
YELLOW="\033[1;33m"
NC="\033[0m"

log()  { echo -e "${GREEN}✔ $1${NC}"; }
warn() { echo -e "${YELLOW}⚠ $1${NC}"; }
err()  { echo -e "${RED}✖ $1${NC}"; }

echo "======================================="
echo " Git Genius Uninstaller"
echo "======================================="
echo

# --------------------------------------------------
# Detect environment
# --------------------------------------------------

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
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

log "Environment: $ENV"

# --------------------------------------------------
# Remove binary
# --------------------------------------------------

REMOVED=false

remove_bin() {
  if [[ -f "$1/$BIN_NAME" ]]; then
    rm -f "$1/$BIN_NAME"
    log "Removed binary: $1/$BIN_NAME"
    REMOVED=true
  fi
}

remove_bin "$HOME/bin"
remove_bin "/usr/local/bin"
remove_bin "/usr/bin"

# --------------------------------------------------
# Remove source directory
# --------------------------------------------------

SRC_DIR="$HOME/.git-genius-src"
if [[ -d "$SRC_DIR" ]]; then
  rm -rf "$SRC_DIR"
  log "Removed source directory: $SRC_DIR"
fi

# --------------------------------------------------
# Remove config cache (safe)
# --------------------------------------------------

CONFIG_DIR="$HOME/.config/git-genius"
if [[ -d "$CONFIG_DIR" ]]; then
  rm -rf "$CONFIG_DIR"
  log "Removed config directory: $CONFIG_DIR"
fi

# --------------------------------------------------
# Remove PATH entries (non-destructive)
# --------------------------------------------------

clean_path_file() {
  FILE="$1"
  [[ -f "$FILE" ]] || return

  if grep -q "git-genius" "$FILE"; then
    sed -i.bak '/git-genius/d' "$FILE"
    log "Cleaned PATH entry from $FILE (backup created)"
  fi
}

clean_path_file "$HOME/.bashrc"
clean_path_file "$HOME/.profile"
clean_path_file "$HOME/.zshrc"

# --------------------------------------------------
# Summary
# --------------------------------------------------

echo
if [[ "$REMOVED" == true ]]; then
  log "Git Genius completely uninstalled ✅"
else
  warn "Git Genius binary not found (already removed?)"
fi

echo
echo "You can now:"
echo "  • Reinstall a fresh version"
echo "  • Test installer safely"
echo
echo "======================================="