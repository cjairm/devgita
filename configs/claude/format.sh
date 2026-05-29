#!/usr/bin/env bash
# PostToolUse hook: format the single file Claude just edited or wrote.
#
# Claude Code delivers the hook payload as JSON on stdin; the edited file path
# lives at `.tool_input.file_path`. (The old `$CLAUDE_FILE_PATHS` env var no
# longer exists in Claude Code 2.x, which is why directory-wide formatting
# silently stopped working.) We parse the path once and route to the relevant
# formatters by extension so we only run tools that apply to the file.
input=$(cat)

FILE=$(printf '%s' "$input" | jq -r '.tool_input.file_path // empty' 2>/dev/null)
[ -z "$FILE" ] && exit 0
[ -f "$FILE" ] || exit 0

BIN="$HOME/.local/share/nvim/mason/bin"

# run <tool> [args...] — invoke a Mason-installed tool only if it's executable,
# never failing the hook (formatting is best-effort).
run() {
    local tool="$1"
    shift
    if [ -x "$tool" ]; then
        "$tool" "$@" || true
    fi
}

case "$FILE" in
*.js | *.jsx | *.ts | *.tsx | *.mjs | *.cjs)
    run "$BIN/eslint_d" "$FILE" --fix
    run "$BIN/prettier" "$FILE" --write
    ;;
*.json | *.css | *.scss | *.less | *.html | *.md | *.markdown | *.yaml | *.yml)
    run "$BIN/prettier" "$FILE" --write
    ;;
*.py)
    run "$BIN/isort" "$FILE"
    run "$BIN/black" "$FILE"
    run "$BIN/flake8" "$FILE"
    ;;
*.go)
    run "$BIN/goimports" -w "$FILE"
    run "$BIN/gofumpt" -w "$FILE"
    run "$BIN/golines" -w "$FILE"
    ;;
*.lua)
    run "$BIN/stylua" "$FILE"
    ;;
*.sh | *.bash)
    run "$BIN/shfmt" -w "$FILE"
    ;;
esac

exit 0
