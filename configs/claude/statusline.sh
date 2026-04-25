#!/bin/bash
input=$(cat)

# Extract core fields
MODEL=$(echo "$input" | jq -r '.model.display_name')
DIR=$(echo "$input" | jq -r '.workspace.current_dir')
PCT=$(echo "$input" | jq -r '.context_window.used_percentage // 0' | cut -d. -f1)

# Colors (optional, used for git status)
GREEN='\033[32m'
YELLOW='\033[33m'
RESET='\033[0m'

# -----------------------------
# Context usage bar
# -----------------------------
BAR_WIDTH=10
FILLED=$((PCT * BAR_WIDTH / 100))
EMPTY=$((BAR_WIDTH - FILLED))

BAR=""

if [ "$FILLED" -gt 0 ]; then
  printf -v FILL "%${FILLED}s"
  BAR="${FILL// /▓}"
fi

if [ "$EMPTY" -gt 0 ]; then
  printf -v PAD "%${EMPTY}s"
  BAR="${BAR}${PAD// /░}"
fi

# -----------------------------
# Git info
# -----------------------------
GIT_PART=""

if git -C "$DIR" rev-parse --git-dir > /dev/null 2>&1; then
  BRANCH=$(git -C "$DIR" branch --show-current 2>/dev/null)

  STAGED=$(git -C "$DIR" diff --cached --numstat 2>/dev/null | wc -l | tr -d ' ')
  MODIFIED=$(git -C "$DIR" diff --numstat 2>/dev/null | wc -l | tr -d ' ')

  GIT_STATUS=""

  if [ "$STAGED" -gt 0 ]; then
    GIT_STATUS="${GREEN}+${STAGED}${RESET}"
  fi

  if [ "$MODIFIED" -gt 0 ]; then
    GIT_STATUS="${GIT_STATUS}${YELLOW}~${MODIFIED}${RESET}"
  fi

  GIT_PART="| $BRANCH $GIT_STATUS"
fi

# -----------------------------
# Final output
# -----------------------------
echo "[$MODEL] 📁 ${DIR##*/} | $BAR $PCT% $GIT_PART"
