#!/usr/bin/env bash
input=$(cat)

eval "$(printf '%s' "$input" | jq -r '
  @sh "MODEL=\(.model.display_name // "?")",
  @sh "DIR=\(.workspace.current_dir // "")",
  @sh "WORKTREE=\(.workspace.git_worktree // "")",
  @sh "AGENT=\(.agent.name // "")",
  @sh "SESSION_ID=\(.session_id // "default")",
  @sh "DUR_MS=\(.cost.total_duration_ms // 0)",
  @sh "ADDED=\(.cost.total_lines_added // 0)",
  @sh "REMOVED=\(.cost.total_lines_removed // 0)",
  @sh "PCT=\((.context_window.used_percentage // 0) | floor)",
  @sh "FIVE_H=\((.rate_limits.five_hour.used_percentage // -1) | floor)",
  @sh "SEVEN_D=\((.rate_limits.seven_day.used_percentage // -1) | floor)"')"

EFFORT=$(jq -r '.effortLevel // empty' "$HOME/.claude/settings.json" 2>/dev/null)

: "${MODEL:=?}" "${DIR:=}" "${PCT:=0}" "${FIVE_H:=-1}" "${SEVEN_D:=-1}"
: "${ADDED:=0}" "${REMOVED:=0}" "${DUR_MS:=0}"

DIM=$'\033[2m'; BOLD=$'\033[1m'; RESET=$'\033[0m'
CYAN=$'\033[36m'; GREEN=$'\033[32m'; YELLOW=$'\033[33m'; RED=$'\033[31m'
MAGENTA=$'\033[35m'; BLUE=$'\033[34m'

TMP_DIR="${TMPDIR:-/tmp}"
DIR_HASH=$(printf '%s' "$DIR" | cksum | awk '{print $1}')
CACHE_FILE="${TMP_DIR%/}/cc-statusline-${SESSION_ID}-${DIR_HASH}"
CACHE_MAX_AGE=5

cache_stale() {
    [ ! -f "$CACHE_FILE" ] && return 0
    local mtime
    mtime=$(stat -f %m "$CACHE_FILE" 2>/dev/null || stat -c %Y "$CACHE_FILE" 2>/dev/null || echo 0)
    [ $(($(date +%s) - mtime)) -gt "$CACHE_MAX_AGE" ]
}

if cache_stale; then
    if [ -n "$DIR" ] && git -C "$DIR" rev-parse --git-dir >/dev/null 2>&1; then
        GB=$(git -C "$DIR" branch --show-current 2>/dev/null)
        GS=$(git -C "$DIR" diff --cached --numstat 2>/dev/null | wc -l | tr -d ' ')
        GM=$(git -C "$DIR" diff --numstat 2>/dev/null | wc -l | tr -d ' ')
        GU=$(git -C "$DIR" ls-files --others --exclude-standard 2>/dev/null | wc -l | tr -d ' ')
        printf '%s|%s|%s|%s\n' "$GB" "$GS" "$GM" "$GU" > "$CACHE_FILE"
    else
        printf '|||\n' > "$CACHE_FILE"
    fi
fi
IFS='|' read -r GB GS GM GU < "$CACHE_FILE"

BAR_WIDTH=12
FILLED=$((PCT * BAR_WIDTH / 100))
[ "$FILLED" -gt "$BAR_WIDTH" ] && FILLED=$BAR_WIDTH
EMPTY=$((BAR_WIDTH - FILLED))
BAR_FILL=""; BAR_EMPTY=""
[ "$FILLED" -gt 0 ] && printf -v BAR_FILL  "%${FILLED}s" "" && BAR_FILL="${BAR_FILL// /█}"
[ "$EMPTY"  -gt 0 ] && printf -v BAR_EMPTY "%${EMPTY}s"  "" && BAR_EMPTY="${BAR_EMPTY// /░}"

if   [ "$PCT" -ge 90 ]; then BAR_COLOR="$RED"
elif [ "$PCT" -ge 70 ]; then BAR_COLOR="$YELLOW"
else                         BAR_COLOR="$GREEN"
fi

color_for_limit() {
    local p="$1"
    if   [ "$p" -ge 90 ]; then printf '%s' "$RED"
    elif [ "$p" -ge 75 ]; then printf '%s' "$YELLOW"
    else                       printf '%s' "$GREEN"
    fi
}

MINS=$((DUR_MS / 60000))
SECS=$(((DUR_MS % 60000) / 1000))
SEP=" ${DIM}|${RESET} "

if [ -n "$HOME" ] && [ "$DIR" = "$HOME" ]; then
    DISPLAY_DIR="~"
elif [ -n "$HOME" ] && [ "${DIR#"$HOME"/}" != "$DIR" ]; then
    DISPLAY_DIR="~/${DIR#"$HOME"/}"
else
    DISPLAY_DIR="$DIR"
fi
L="${DIM}${DISPLAY_DIR}${RESET}"
if [ -n "$GB" ]; then
    L+="${SEP}${BLUE}${GB}${RESET}"
    [ "${GS:-0}" -gt 0 ] && L+=" ${GREEN}+${GS}${RESET}"
    [ "${GM:-0}" -gt 0 ] && L+=" ${YELLOW}~${GM}${RESET}"
    [ "${GU:-0}" -gt 0 ] && L+=" ${CYAN}?${GU}${RESET}"
fi
if [ "${ADDED:-0}" -gt 0 ] || [ "${REMOVED:-0}" -gt 0 ]; then
    L+="${SEP}${GREEN}+${ADDED}${RESET}${DIM}/${RESET}${RED}-${REMOVED}${RESET}"
fi
SESSION=""
[ -n "$AGENT" ]    && SESSION+="${MAGENTA}@${AGENT}${RESET}"
[ -n "$WORKTREE" ] && SESSION+="${SESSION:+ }${DIM}wt:${WORKTREE}${RESET}"
[ -n "$SESSION" ]  && L+="${SEP}${SESSION}"
L+="${SEP}${BAR_COLOR}${BAR_FILL}${DIM}${BAR_EMPTY}${RESET} ${PCT}% ctx"
if [ "$FIVE_H" -ge 0 ] || [ "$SEVEN_D" -ge 0 ]; then
    L+="${SEP}"
    [ "$FIVE_H"  -ge 0 ] && L+="$(color_for_limit "$FIVE_H")5h:${FIVE_H}%${RESET}"
    [ "$SEVEN_D" -ge 0 ] && L+=" $(color_for_limit "$SEVEN_D")7d:${SEVEN_D}%${RESET}"
fi
L+="${SEP}${MINS}m${SECS}s"
L+="${SEP}${CYAN}${BOLD}${MODEL}${RESET}"
[ -n "$EFFORT" ] && L+=" ${DIM}·${RESET} ${CYAN}${EFFORT}${RESET}"

echo "$L"
