if command -v mise &> /dev/null; then
  eval "$(mise activate zsh)"
fi

if command -v zoxide &> /dev/null; then
  eval "$(zoxide init bash)"
fi

# plugins=(
#   zsh-autosuggestions
#   zsh-syntax-highlighting
# )

# Extended capabilities
dg() {
  local allowed=(
    # git branch management
    "delete-branch"
    "read-branches"
    "refresh-branch"
    "reset-main-branch"
    # npm dependency management
    "reinstall-libraries"
    "reinstall-library"
  )

  local cmd="$1"
  local target="$2"

  # Validate command
  local is_valid=false
  for a in "${allowed[@]}"; do
    if [[ "$cmd" == "$a" ]]; then
      is_valid=true
      break
    fi
  done

  if [[ "$is_valid" == false ]]; then
    echo "valid commands are: ${allowed[*]}"
    return 1
  fi

  case "$cmd" in
    # -------------------------------------
    # Git Branch Management 
    # -------------------------------------
    delete-branch)
      # Default to "main" if not provided
      target="${target:-main}"

      git checkout "$target" \
        && git fetch origin \
        && git pull origin "$target" \
        && git branch | fzf-tmux -p 50% --reverse | xargs git branch -D
      ;;

    read-branches)
      selected_branch=$(git branch | sort -u | fzf-tmux -p 50% --reverse)
      if [[ -n "$selected_branch" ]]; then
        # Copy to clipboard
        if command -v pbcopy &>/dev/null; then
          echo -n "$selected_branch" | pbcopy
        elif command -v xclip &>/dev/null; then
          echo -n "$selected_branch" | xclip -selection clipboard
        fi
        echo "Branch '$selected_branch' copied to clipboard!"
      fi
      ;;

    refresh-branch)
      # Default to "main" if not provided
      target="${target:-main}"

      git checkout "$target" \
        && git pull origin "$target" \
        && git checkout - \
        && git merge "$target"
      ;;

    reset-main-branch)
      git checkout main && git reset --hard origin/main
      ;;

    # -------------------------------------
    # NPM Dependency Management
    # -------------------------------------
    reinstall-libraries)
      git clean -Xdf \
        && rm -rf node_modules/ \
        && npm install \
        && rm -f tsconfig.tsbuildinfo
      ;;

    reinstall-library)
      # Must provide target
      if [[ -z "$target" ]]; then
        echo "update-library requires a library name"
        return 1
      fi

      rm -rf "node_modules/$target" && npm install
      ;;

  esac
}
