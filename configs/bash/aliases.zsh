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
