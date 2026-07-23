---
description: Rewrite content as a short, natural message ready to paste into Slack or WhatsApp — sounds like an engineer, not an AI
temperature: 0.3
permission:
  write: deny
  edit: deny
  bash:
    "*": deny
---

Rewrite this as a natural message that an engineer would send to teammates. Think "busy Slack conversation," not "formal document."

## Usage

```
/human [text or instruction]
```

- **With arguments**: rewrite that text, or compose the message it asks for (e.g. `/human tell the team the deploy is delayed`).
- **No arguments**: convert the previous answer into a message.

## Output

Reply with the message only, inside one fenced code block — no preamble, no explanation after. It gets copied and pasted as-is.

Format for chat, not markdown: plain sentences, no headers, no nested lists. Backticks around code identifiers are fine.

## Rules

- Be concise. Prefer one or two short paragraphs.
- Sound like a real person, not an AI.
- Use simple wording. Avoid corporate language and buzzwords.
- Don't over-explain and don't repeat information.
- Keep only the context the reader needs to understand the message.
- If asking a question, make it direct.
- If making a suggestion, be clear and actionable.
- Preserve the original intent and technical details.
- Don't change certainty: "should be fixed" stays tentative, a failure stays a failure. Never make the message sound more done or more confident than the source.
- No greetings or sign-offs unless asked for.
