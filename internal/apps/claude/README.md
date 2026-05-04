# Claude Code App

Installs and configures [Claude Code](https://code.claude.com/docs/en) — the official CLI for Claude.

## After Installation

**Start Claude Code:**

```bash
cc
```

**Useful slash commands inside Claude Code:**

- `/help` — Get help with Claude Code features
- `/config` — Configure settings (model, theme, permissions)
- `/fast` — Toggle fast mode (Opus 4.6 with faster output)
- `/loop` — Run a prompt or command on a recurring interval
- `/skills` — List available skills
- `/init` — Initialize CLAUDE.md for your project
- `/review` — Review a pull request

**View configuration:**

```bash
# Settings, themes, skills, commands, and agents are installed in:
ls ~/.claude/
```

**Reset chat history:**

```bash
# message/input history
rm ~/.claude/history.jsonl
# chats
rm -rf ~/.claude/sessions
rm -rf ~/.claude/projects
rm -rf ~/.claude/tasks
```

Then restart Claude Code.
