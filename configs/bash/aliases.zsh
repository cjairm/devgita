# File system
alias ls='eza -lh --group-directories-first --icons'
alias lsa='ls -a'
alias lt='eza --tree --level=2 --long --icons --git'
alias lta='lt -a'
alias ff="fzf --preview 'bat --style=numbers --color=always {}'"

# quick edits
alias edit-zsh="nvim ~/.zshrc"
alias edit-ohmyzsh="nvim ~/.oh-my-zsh"
# TODO: Add aerospace config
alias edit-tmux="nvim ~/.tmux.conf"
# TODO: Add config per app (?)
alias edit-nvim="nvim ~/.config/nvim/init.lua"
alias edit-alias="nvim ~/.config/devgita/aliases.zsh"

# Directories
alias ..='cd ..'
alias ...='cd ../..'
alias ....='cd ../../..'

# Tools
alias lzg='lazygit'
alias lzd='lazydocker'

# Neovim
alias n='nvim'
alias ns="fd --type f --hidden --exclude .git | fzf-tmux -p  --reverse | xargs nvim"

# ---- Zoxide (better cd) ----
alias cd="z"

# ---- Syntax highlighting (better cat) ----
alias cat="bat"

# Git
alias g='git'
alias gcm='git commit -m'
alias gcam='git commit -a -m'
alias gcad='git commit -a --amend'

# Tmux
alias tml="tmux list-sessions"
alias tmk="tmux kill-session -t"
alias tms="tmux switch-client -t"
alias tmkall="tmux kill-server"
alias tmclear="clear && tmux clear-history && clear"
tmn() {
    local session_name="$1"
    local path="${2:-$HOME}"
    /usr/local/bin/tmux new-session -d -s "$session_name" -c "$path"
    /usr/local/bin/tmux switch-client -t "$session_name"
}

# Compression
compress() { tar -czf "${1%/}.tar.gz" "${1%/}"; }
alias decompress="tar -xzf"

# Mise
alias mx="mise x --"

# Custom
dg() {
  local allowed=(
    # git branch management
    "clean-branch"
    "refresh-branch"
    "reset-branch"
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
    clean-branch)
      # Default to "main" if not provided
      target="${target:-main}"

      git checkout "$target" \
        && git reset --hard origin/"$target"
      ;;

    refresh-branch)
      # Default to "main" if not provided
      target="${target:-main}"

      git checkout "$target" \
        && git pull origin "$target" \
        && git checkout - \
        && git merge "$target"
      ;;

    reset-branch)
      # Default to "main" if not provided
      target="${target:-main}"

      git checkout "$target" \
        && git fetch origin \
        && git pull origin "$target" \
        && git branch | fzf-tmux -p 50% --reverse | xargs git branch -D
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
