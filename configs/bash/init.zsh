if command -v mise &> /dev/null; then
  eval "$(~/.local/bin/mise activate zsh)"
fi

if command -v zoxide &> /dev/null; then
  eval "$(zoxide init bash)"
fi
