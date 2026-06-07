# OpenCode App

Installs and configures [OpenCode](https://opencode.ai/docs) — a terminal-based AI code editor.

## After Installation

**Start OpenCode:**

```bash
opencode
```

**View configuration:**

```bash
ls ~/.config/opencode/
```

---

## Provider

Use **[OpenRouter](https://openrouter.ai)** — one API key, access to all models below.

```bash
export OPENROUTER_API_KEY="your-key-here"
```

---

## Models

| Model                     | Role                       |
| ------------------------- | -------------------------- |
| `anthropic/claude-opus-4` | Hardest reasoning          |
| `qwen/qwen3-coder`        | Daily coding (default)     |
| `moonshotai/kimi-k2`      | Agents + large repos       |
| `qwen/qwen3.5-coder-480b` | Deep review + architecture |
| `minimax/minimax-m1`      | Cheap bulk tasks           |

### When to switch

- **claude-opus-4** — when nothing else solves it
- **qwen3-coder** — everyday coding, bug fixes, PR reviews
- **kimi-k2** — large codebases, multi-file refactors, long agent workflows
- **qwen3.5-coder-480b** — architecture decisions, hard debugging, critical reviews
- **minimax-m1** — background or non-critical automation
