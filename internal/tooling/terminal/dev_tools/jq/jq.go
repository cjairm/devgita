// jq JSON processor with devgita integration
//
// jq is a lightweight and flexible command-line JSON processor. It lets you
// slice, filter, map, and transform structured JSON data with ease.
//
// References:
// - jq Documentation: https://jqlang.github.io/jq/
// - jq Manual: https://jqlang.github.io/jq/manual/
//
// Common jq commands available through ExecuteCommand():
//   - jq '.' file.json          - Pretty-print JSON
//   - jq '.key' file.json       - Extract a field
//   - jq '.[] | .field' file.json - Iterate array, extract field
//   - jq -r '.key' file.json    - Raw string output (no quotes)
//   - jq -c '.' file.json       - Compact output

package jq

import (
	"fmt"
	"os"
	"strconv"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
)

// reviewThreadsFilter turns a GitHub GraphQL pull-request reviewThreads
// response into compact markdown intended to be fed to an LLM: one block per
// thread, headed by "file:line (thread <id>)", an optional diff hunk for
// context, then one line per comment as "**author** (<id>): body".
//
// The layout is deliberately terse — no "## Location" / "### Comments" scaffold
// — so it reads clearly while spending few tokens. Thread and comment ids are
// kept because follow-up tasks (resolve, reply) act on them.
//
// The $resolved argument controls filtering:
//   - true  → only resolved threads
//   - false → only unresolved threads
//   - null  → all threads (no filtering)
//
// The GraphQL query feeding this MUST select, per reviewThread: id, isResolved,
// path, line, originalLine; per comment: id, author.login, body; and the diff
// hunk via a "firstComment: comments(first: 1) { nodes { diffHunk } }" alias
// (kept separate so diffHunk isn't refetched for every comment). Missing fields
// render as null/empty.
//
// Triple backticks are concatenated in because Go raw string literals cannot
// contain backticks; everything else stays raw so jq sees literal "\n" / "\(...)".
const reviewThreadsFilter = `.data.repository.pullRequest.reviewThreads.nodes
| map(select(if $resolved == null then true else .isResolved == $resolved end))
| .[]
| "## \(.path):\(.line // .originalLine // "?") (thread \(.id))\n\n"
  + (
      .firstComment.nodes[0].diffHunk
      | if . then "` + "```" + `diff\n\(.)\n` + "```" + `\n\n" else "" end
    )
  + (
      .comments.nodes
      | map("**\(.author.login)** (\(.id)): \(.body)")
      | join("\n\n")
    )
  + "\n\n---\n"`

// prViewFilter renders `gh pr view --json ...` output (a single object) into a
// compact three-line summary. Assumes the default field set selected by
// githubcli.PRView; missing fields fall back so it never errors on absence.
const prViewFilter = `"PR #\(.number): \(.title)\n"
  + "state: \(.state)  mergeable: \(.mergeable // "?")  review: \(.reviewDecision // "none")\n"
  + "branch: \(.headRefName // "?") -> \(.baseRefName // "?")"`

// prChecksFilter renders `gh pr checks --json ...` output (an array) into one
// line per check: "<state><TAB><name>" plus the link when present. Empty input
// renders a short "No checks." note.
const prChecksFilter = `if length == 0 then "No checks." else
  ( .[] | "\(.state)\t\(.name)" + (if (.link // "") != "" then "  \(.link)" else "" end) )
end`

type Jq struct {
	Cmd  cmd.Command
	Base cmd.BaseCommandExecutor
}

func New() *Jq {
	return &Jq{
		Cmd:  cmd.NewCommand(),
		Base: cmd.NewBaseCommand(),
	}
}

func (j *Jq) Install() error {
	return j.Cmd.InstallPackage(constants.Jq)
}

func (j *Jq) SoftInstall() error {
	return j.Cmd.MaybeInstallPackage(constants.Jq)
}

func (j *Jq) ForceInstall() error {
	if err := j.Uninstall(); err != nil {
		return fmt.Errorf("failed to uninstall jq: %w", err)
	}
	return j.Install()
}

func (j *Jq) Uninstall() error {
	return fmt.Errorf("jq uninstall not supported through devgita")
}

func (j *Jq) ForceConfigure() error {
	return nil
}

func (j *Jq) SoftConfigure() error {
	return nil
}

func (j *Jq) ExecuteCommand(args ...string) error {
	_, _, err := j.Base.ExecCommand(cmd.CommandParams{
		IsSudo:  false,
		Command: constants.Jq,
		Args:    args,
	})
	if err != nil {
		return fmt.Errorf("failed to run jq command: %w", err)
	}
	return nil
}

func (j *Jq) Update() error {
	return fmt.Errorf("jq update not implemented through devgita")
}

// runFilter is the standard entry point every Format* method goes through: it
// writes the JSON payload to a temp file, runs `jq -r [extraArgs] <filter>
// <tmpfile>`, and returns stdout. Passing the payload as a file argument
// (rather than stdin) survives arbitrarily large inputs and keeps execution
// routed through Base.ExecCommand for testability.
//
// Standard for adding a new formatter: define a documented filter const that
// lists the input fields it expects, then write a thin FormatX(json, ...) that
// calls runFilter. Keep output compact and LLM-oriented.
func (j *Jq) runFilter(payload, filter string, extraArgs ...string) (string, error) {
	tmp, err := os.CreateTemp("", "devgita-jq-*.json")
	if err != nil {
		return "", fmt.Errorf("jq: failed to create temp file: %w", err)
	}
	tmpName := tmp.Name()
	defer func() { _ = os.Remove(tmpName) }()

	if _, err := tmp.WriteString(payload); err != nil {
		_ = tmp.Close()
		return "", fmt.Errorf("jq: failed to write temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return "", fmt.Errorf("jq: failed to close temp file: %w", err)
	}

	args := append([]string{"-r"}, extraArgs...)
	args = append(args, filter, tmpName)

	stdout, stderr, err := j.Base.ExecCommand(cmd.CommandParams{
		Command: constants.Jq,
		Args:    args,
	})
	if err != nil {
		if stderr != "" {
			return "", fmt.Errorf("jq: %s", stderr)
		}
		return "", fmt.Errorf("jq: failed to run filter: %w", err)
	}
	return stdout, nil
}

// FormatReviewThreads runs reviewThreadsFilter over a GitHub GraphQL
// reviewThreads JSON payload and returns the formatted markdown.
//
// resolved filters by thread resolution status: pass &true for resolved
// threads only, &false for unresolved only, or nil to include all threads.
func (j *Jq) FormatReviewThreads(ghJSON string, resolved *bool) (string, error) {
	// jq's --argjson parses the value as JSON, so "true"/"false"/"null" all
	// arrive as their native JSON types — exactly what the filter compares against.
	resolvedArg := "null"
	if resolved != nil {
		resolvedArg = strconv.FormatBool(*resolved)
	}
	return j.runFilter(ghJSON, reviewThreadsFilter, "--argjson", "resolved", resolvedArg)
}

// FormatPRView renders `gh pr view --json ...` output into a compact summary.
func (j *Jq) FormatPRView(ghJSON string) (string, error) {
	return j.runFilter(ghJSON, prViewFilter)
}

// FormatPRChecks renders `gh pr checks --json ...` output into one line per check.
func (j *Jq) FormatPRChecks(ghJSON string) (string, error) {
	return j.runFilter(ghJSON, prChecksFilter)
}
