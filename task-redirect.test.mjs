// Behavioral test for configs/opencode/plugin/task-redirect.js — the
// OpenCode-side mirror of configs/claude/task-redirect.sh (see
// task_redirect_test.go, package main, for the bash script's equivalent
// end-to-end test).
//
// This file is deliberately NOT under configs/opencode/plugin/: that
// directory is copied byte-for-byte to the user's OpenCode plugin dir
// (internal/apps/opencode/opencode.go's ForceConfigure step, via
// files.CopyDir) — a test file living there would ship to every installed
// machine. Keeping this test at the repo root (sibling to
// task_redirect_test.go) mirrors the Go test's placement without polluting
// the deployed plugin directory.
//
// Run with: node --test task-redirect.test.mjs
// Uses only Node's built-in test runner and assert module (available since
// Node 18) — no new dependency for a single test file, per CLAUDE.md's
// "prefer the standard library" rule applied to this repo's one JS file.

import { test } from "node:test";
import assert from "node:assert/strict";
import { mkdtempSync, writeFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";
import {
  TaskRedirect,
  isDevgitaRepo,
} from "./configs/opencode/plugin/task-redirect.js";

// This test file lives at the repo root, so its own directory is a devgita
// repo (repo-root go.mod has module github.com/cjairm/devgita). Used as the
// devgita-cwd for the release-gating tests below.
const DEVGITA_DIR = dirname(fileURLToPath(import.meta.url));

// runHook drives the plugin exactly as OpenCode would: build the plugin,
// invoke its tool.execute.before hook with a bash tool call, and report
// whether it denied (threw) or allowed (returned normally). `ctx` is the
// OpenCode plugin context (defaults to the devgita repo dir so the release
// rules are exercised in their firing state unless a test overrides it).
async function runHook(command, env = {}, ctx = { directory: DEVGITA_DIR }) {
  const previous = {};
  for (const [key, value] of Object.entries(env)) {
    previous[key] = process.env[key];
    if (value === undefined) {
      delete process.env[key];
    } else {
      process.env[key] = value;
    }
  }
  try {
    const plugin = await TaskRedirect(ctx);
    const hook = plugin["tool.execute.before"];
    try {
      await hook({ tool: "bash" }, { args: { command } });
      return { denied: false, message: null };
    } catch (err) {
      return { denied: true, message: err.message };
    }
  } finally {
    for (const [key, value] of Object.entries(previous)) {
      if (value === undefined) {
        delete process.env[key];
      } else {
        process.env[key] = value;
      }
    }
  }
}

test("allows legitimate single commands", async () => {
  const allowed = [
    "git diff",
    "git diff HEAD~1",
    "git diff --stat",
    "git log",
    "git log -5",
    "git log --oneline",
    "git tag",
    "git tag v1.0.0",
    "git tag -l",
    "git reset --soft HEAD",
    "git worktree list",
    "git worktree prune",
    "git status",
    "git push origin main",
    'git commit -m "fix: something"',
    // Compound commands where no segment matches any rule.
    "cd some/dir && git status",
    "git fetch && git log -5",
    // A commit message mentioning a trigger word must not itself trigger a
    // rule — "worktree" here is message text, not a `git worktree` call.
    'git commit -m "fix: worktree stuff"',
    // A commit message containing separator-like characters (';', '&&')
    // must not be split apart: the quoted span is one segment, and neither
    // half looks like a git invocation on its own.
    'git commit -m "fix: a; b"',
    // Pathological case for quote-aware splitting: this commit message
    // literally contains "&& git worktree" as text. A naive
    // (non-quote-aware) splitter would slice this into a segment starting
    // with "git worktree" and falsely deny it. Still a single `git commit`
    // command — must be allowed.
    'git commit -m "notes && git worktree stuff"',
    // gh commands that must NOT match the new gh rules: `pr view` is a
    // different subcommand from `pr review`/`pr checks`; `pr status`/`pr list`
    // are neither; a graphql query without reviewThreads and a bare `gh api`
    // are not the review-threads fetch.
    "gh pr view",
    "gh pr view --json title",
    "gh pr status",
    "gh pr list",
    "gh api graphql -f query='{ viewer { login } }'",
    "gh api repos/cjairm/devgita",
  ];
  for (const command of allowed) {
    const result = await runHook(command);
    assert.equal(
      result.denied,
      false,
      `expected allow for ${JSON.stringify(command)}, got deny: ${result.message}`,
    );
  }
});

test("denies narrow patterns, including compound commands and env-var prefixes", async () => {
  const cases = [
    ["git diff main..feature", "devgita task review-package"],
    ["git diff v1.2.0..v1.3.0", "devgita task review-package"],
    ["git diff --stat A..B", "devgita task review-package"],
    ["git log --oneline base..head", "devgita task review-package"],
    ["git worktree add ../wt -b feature-x", "devgita task worktree-start"],
    ["git worktree remove ../wt", "devgita task worktree-finish"],
    // New gh rules — all GLOBAL, so they deny regardless of cwd. (runHook
    // defaults ctx to the devgita dir, but these rules never consult it; the
    // release-gating tests below cover the cwd-dependent rules.)
    ["gh pr checks", "devgita task pr-checks"],
    ["gh pr checks --watch", "devgita task pr-checks"],
    ["gh pr review --approve", "devgita task submit-review"],
    ["gh pr review --request-changes -b bad", "devgita task submit-review"],
    [
      "gh api graphql --paginate -f query='{ repository { pullRequest { reviewThreads { nodes { id } } } } }'",
      "devgita task review-threads",
    ],
    // Compound commands: a matching segment anywhere in the chain must
    // deny, not just a bare command at position 0.
    [
      "cd some/dir && git worktree add ../wt -b x",
      "devgita task worktree-start",
    ],
    ["git status; git worktree remove ../wt", "devgita task worktree-finish"],
    ["git fetch && git diff main..feature", "devgita task review-package"],
    ["gh pr view && gh pr checks", "devgita task pr-checks"],
    // git diff a..b | less: the LHS of the pipe is still a git invocation
    // itself, so this must deny too.
    ["git diff main..feature | less", "devgita task review-package"],
    // Env-var-prefix case (no separator character before `git` at all):
    // deliberately handled now that GIT_PREFIX is being reworked anyway — a
    // simple `NAME=value` prefix, single or repeated, in front of `git`
    // still denies.
    ["GIT_PAGER=cat git diff main..feature", "devgita task review-package"],
    [
      "FOO=bar BAZ=qux git worktree add ../wt -b x",
      "devgita task worktree-start",
    ],
  ];
  for (const [command, wantReplacement] of cases) {
    const result = await runHook(command);
    assert.equal(
      result.denied,
      true,
      `expected deny for ${JSON.stringify(command)}`,
    );
    assert.ok(
      result.message.includes(wantReplacement),
      `expected deny reason for ${JSON.stringify(command)} to mention ${JSON.stringify(wantReplacement)}, got ${JSON.stringify(result.message)}`,
    );
    assert.ok(
      result.message.includes("DEVGITA_SKIP_TASK_REDIRECT"),
      `expected deny reason for ${JSON.stringify(command)} to state the bypass escape hatch, got ${JSON.stringify(result.message)}`,
    );
  }
});

// isDevgitaRepo is unit-tested directly with injected paths so the release
// gate's core decision is verified independently of the plugin plumbing.
test("isDevgitaRepo detects the devgita go.mod by walking up", () => {
  // The repo root itself is devgita.
  assert.equal(isDevgitaRepo(DEVGITA_DIR), true);
  // A nested subdirectory still resolves upward to the devgita go.mod.
  assert.equal(isDevgitaRepo(join(DEVGITA_DIR, "cmd")), true);

  // A dir with no go.mod anywhere up the (temp) tree: not devgita.
  const noGoMod = mkdtempSync(join(tmpdir(), "no-gomod-"));
  assert.equal(isDevgitaRepo(noGoMod), false);

  // A dir whose go.mod is a different module: not devgita.
  const otherModule = mkdtempSync(join(tmpdir(), "other-mod-"));
  writeFileSync(
    join(otherModule, "go.mod"),
    "module github.com/other/thing\n\ngo 1.25\n",
  );
  assert.equal(isDevgitaRepo(otherModule), false);

  // Indeterminate inputs fail toward false (release rules do not fire).
  assert.equal(isDevgitaRepo(undefined), false);
  assert.equal(isDevgitaRepo(""), false);
  assert.equal(isDevgitaRepo(123), false);
});

test("release rules deny inside devgita, allow everywhere else", async () => {
  const releaseCommands = [
    "git reset --soft HEAD~1",
    "git reset --soft HEAD~3",
    "git tag -a v0.12.0 -m release",
    "git tag -a -m release v0.12.0",
    "cd wt && git reset --soft HEAD~2",
    "git status && git tag -a v1.0.0 -m release",
  ];

  const noGoMod = mkdtempSync(join(tmpdir(), "no-gomod-"));
  const otherModule = mkdtempSync(join(tmpdir(), "other-mod-"));
  writeFileSync(
    join(otherModule, "go.mod"),
    "module github.com/other/thing\n\ngo 1.25\n",
  );

  for (const command of releaseCommands) {
    // Inside devgita: deny.
    const inside = await runHook(command, {}, { directory: DEVGITA_DIR });
    assert.equal(
      inside.denied,
      true,
      `expected deny inside devgita for ${JSON.stringify(command)}`,
    );
    assert.ok(inside.message.includes("devgita task release"));

    // Outside devgita (no go.mod, and a different module): allow.
    for (const dir of [noGoMod, otherModule]) {
      const outside = await runHook(command, {}, { directory: dir });
      assert.equal(
        outside.denied,
        false,
        `expected allow outside devgita (${dir}) for ${JSON.stringify(command)}, got: ${outside.message}`,
      );
    }
  }
});

test("release gate uses worktree fallback and fails toward not firing", async () => {
  // No `directory` in ctx: falls back to `worktree`.
  const viaWorktree = await runHook(
    "git reset --soft HEAD~1",
    {},
    { worktree: DEVGITA_DIR },
  );
  assert.equal(viaWorktree.denied, true);

  // Empty ctx and a non-devgita process.cwd() would allow; here we assert the
  // explicit non-devgita directory path allows (the fail-toward-not-firing
  // posture), which is the safety-critical direction.
  const noGoMod = mkdtempSync(join(tmpdir(), "no-gomod-"));
  const failOpen = await runHook(
    "git tag -a v1.2.3 -m r",
    {},
    { directory: noGoMod },
  );
  assert.equal(failOpen.denied, false);
});

test("bypass env var allows everything", async () => {
  const result = await runHook("git worktree add ../wt -b x", {
    DEVGITA_SKIP_TASK_REDIRECT: "1",
  });
  assert.equal(result.denied, false);
});

test("ignores non-bash tool calls and non-string commands", async () => {
  const plugin = await TaskRedirect();
  const hook = plugin["tool.execute.before"];

  // Non-bash tool: never even inspects the command.
  await assert.doesNotReject(() =>
    hook(
      { tool: "edit" },
      { args: { command: "git worktree add ../wt -b x" } },
    ),
  );

  // Missing/non-string command: falls through without throwing.
  await assert.doesNotReject(() => hook({ tool: "bash" }, { args: {} }));
  await assert.doesNotReject(() =>
    hook({ tool: "bash" }, { args: { command: 123 } }),
  );
});
