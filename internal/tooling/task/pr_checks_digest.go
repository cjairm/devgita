package task

import (
	"fmt"
	"regexp"
	"strings"
)

// actionsJobLinkRE matches a GitHub Actions workflow-run check's link field:
// https://github.com/<owner>/<repo>/actions/runs/<run-id>/job/<job-id>
// optionally followed by a "#step:N:M" fragment (the UI deep-link to a
// specific step) or a trailing slash. Anything else — external checks,
// commit statuses, or an empty link — does not match, and per the cycle
// plan that must never be guessed at: the caller falls back to a one-line
// "log unavailable: external check" note instead.
var actionsJobLinkRE = regexp.MustCompile(
	`^https://github\.com/[^/]+/[^/]+/actions/runs/\d+/job/(\d+)/?(?:#.*)?$`,
)

// parseActionsJobID extracts the job id from a GitHub Actions check link. It
// returns ("", false) for any link that doesn't match the workflow-run job
// shape exactly.
func parseActionsJobID(link string) (string, bool) {
	m := actionsJobLinkRE.FindStringSubmatch(strings.TrimSpace(link))
	if m == nil {
		return "", false
	}
	return m[1], true
}

// logLineTimestampRE strips the leading RFC3339-ish timestamp gh prints on
// every log line (e.g. "2026-07-22T18:32:18.0019976Z ").
var logLineTimestampRE = regexp.MustCompile(
	`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?Z\s*`,
)

// splitLogLines turns a raw `gh run view --log-failed` payload into cleaned
// content lines: one gh log line is "<job>\t<step>\t<timestamp> <message>"
// (gh prepends the job and step name so multi-job/multi-step output can be
// told apart); since the caller already knows which failing check this log
// belongs to, the job/step columns are pure repetition and are dropped by
// splitting on gh's first two tabs (its own column separators) and keeping
// everything after them — not on the last tab, which a tab embedded in the
// step's own message (TSV output, tab-formatted diagnostics, etc.) could
// otherwise shadow, silently swallowing real message content. A line that
// doesn't have gh's two-tab column shape is left unchanged rather than
// guessed at. The per-line timestamp is stripped too — see dedupConsecutive's
// doc comment for why that matters. A trailing empty element from a final
// newline is dropped; no other filtering happens (blank lines and GitHub's
// "##[group]" markers are left as-is).
func splitLogLines(raw string) []string {
	raw = strings.ReplaceAll(raw, "\r\n", "\n")
	rawLines := strings.Split(raw, "\n")
	if len(rawLines) > 0 && rawLines[len(rawLines)-1] == "" {
		rawLines = rawLines[:len(rawLines)-1]
	}
	lines := make([]string, 0, len(rawLines))
	for _, l := range rawLines {
		if parts := strings.SplitN(l, "\t", 3); len(parts) == 3 {
			l = parts[2]
		}
		l = logLineTimestampRE.ReplaceAllString(l, "")
		lines = append(lines, l)
	}
	return lines
}

// dedupConsecutive collapses runs of consecutive identical lines into one,
// suffixed with "(×N)" when a run is longer than one line.
//
// Dedup semantics (deliberately narrow): only CONSECUTIVE repeats collapse.
// This is the shape CI noise actually takes — a retry/poll loop prints the
// identical line dozens of times in a row — and it needs the timestamp
// already stripped from each line (splitLogLines does this), since gh
// timestamps every line and a naive compare would never see two "identical"
// lines. Non-consecutive repeats (the same text appearing far apart in the
// log) are intentionally NOT collapsed: that would need unbounded lookback
// and risks silently merging two distinct failures that happen to share a
// line of text.
func dedupConsecutive(lines []string) []string {
	out := make([]string, 0, len(lines))
	i := 0
	for i < len(lines) {
		j := i + 1
		for j < len(lines) && lines[j] == lines[i] {
			j++
		}
		run := j - i
		if run > 1 {
			out = append(out, fmt.Sprintf("%s  (×%d)", lines[i], run))
		} else {
			out = append(out, lines[i])
		}
		i = j
	}
	return out
}

// digestLogTail reduces a raw failing-job log payload to a bounded,
// deduplicated tail of plain failure text, suitable for embedding under a
// failing check's one-line summary.
//
// Pipeline: strip per-line noise (job/step columns, timestamps) → collapse
// consecutive-repeat lines → keep only the last maxLines lines (the tail,
// since gh's failed-step log starts with setup/checkout noise and the
// actual failure is almost always at the end).
//
// When lines are dropped by the tail cut, a truncation receipt ("… N earlier
// lines omitted") is prepended so the cut is never silent — and is NEVER
// emitted when nothing was cut. The receipt line itself counts toward
// maxLines, so the total returned line count never exceeds maxLines (when
// maxLines >= 1).
//
// Returns "" only when raw contains no lines at all (the caller treats that
// as "no log content", distinct from a log that exists but was truncated).
func digestLogTail(raw string, maxLines int) string {
	lines := splitLogLines(raw)
	if len(lines) == 0 {
		return ""
	}

	deduped := dedupConsecutive(lines)

	if maxLines < 1 {
		maxLines = 1
	}
	if len(deduped) <= maxLines {
		return strings.Join(deduped, "\n")
	}

	keep := maxLines - 1
	omitted := len(deduped) - keep
	tail := deduped[len(deduped)-keep:]
	receipt := fmt.Sprintf("… %d earlier lines omitted", omitted)
	return strings.Join(append([]string{receipt}, tail...), "\n")
}
