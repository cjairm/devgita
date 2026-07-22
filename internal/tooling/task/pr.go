package task

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cjairm/devgita/internal/tooling/terminal/dev_tools/githubcli"
	"github.com/cjairm/devgita/internal/tooling/terminal/dev_tools/jq"
)

// PRManager wires the githubcli primitives to the jq formatters so the dg task
// PR subcommands return compact, agent-friendly output. gh fetches/acts; jq
// renders. Methods return the text to print (markdown, a URL, or a one-line
// confirmation) plus an error.
type PRManager struct {
	Gh *githubcli.GithubCli
	Jq *jq.Jq
}

// NewPR creates a PRManager with real executors.
func NewPR() *PRManager {
	return &PRManager{Gh: githubcli.New(), Jq: jq.New()}
}

// resolvedPtrForState maps a --state value to the *bool the jq filter expects:
// "resolved" → &true, "unresolved"/"" → &false, "all" → nil.
func resolvedPtrForState(state string) (*bool, error) {
	t, f := true, false
	switch strings.ToLower(strings.TrimSpace(state)) {
	case "", "unresolved":
		return &f, nil
	case "resolved":
		return &t, nil
	case "all":
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown state %q (use unresolved, resolved, or all)", state)
	}
}

// resolveOwnerRepoPR fills in owner/name from the current repo and, when
// prNumber is empty, the PR number for the current branch.
func (p *PRManager) resolveOwnerRepoPR(prNumber string) (owner, name, pr string, err error) {
	repo, err := p.Gh.CurrentRepo()
	if err != nil {
		return "", "", "", err
	}
	parts := strings.SplitN(repo, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", "", fmt.Errorf("unexpected repo format %q (want owner/name)", repo)
	}

	pr = prNumber
	if pr == "" {
		pr, err = p.Gh.CurrentPRNumber()
		if err != nil {
			return "", "", "", err
		}
		if pr == "" {
			return "", "", "", fmt.Errorf("no pull request found for the current branch; pass --pr")
		}
	}
	return parts[0], parts[1], pr, nil
}

// ReviewThreads fetches a PR's inline review threads (filtered by state:
// "unresolved" default, "resolved", or "all") together with its review
// summary bodies and top-level conversation comments, and renders both as
// markdown. The discussion (summaries + conversation) is always included
// regardless of --state — reviews and conversation comments have no
// resolved/unresolved status of their own.
func (p *PRManager) ReviewThreads(prNumber, state string) (string, error) {
	resolved, err := resolvedPtrForState(state)
	if err != nil {
		return "", err
	}
	owner, name, pr, err := p.resolveOwnerRepoPR(prNumber)
	if err != nil {
		return "", err
	}

	rawThreads, err := p.Gh.FetchReviewThreads(owner, name, pr)
	if err != nil {
		return "", err
	}
	threads, err := p.Jq.FormatReviewThreads(rawThreads, resolved)
	if err != nil {
		return "", err
	}

	rawDiscussion, err := p.Gh.FetchPRDiscussion(owner, name, pr)
	if err != nil {
		return "", err
	}
	discussion, err := p.Jq.FormatPRDiscussion(rawDiscussion)
	if err != nil {
		return "", err
	}

	parts := make([]string, 0, 2)
	if strings.TrimSpace(threads) != "" {
		parts = append(parts, strings.TrimSpace(threads))
	}
	if strings.TrimSpace(discussion) != "" {
		parts = append(parts, strings.TrimSpace(discussion))
	}
	if len(parts) == 0 {
		return "No review threads or comments.", nil
	}
	return strings.Join(parts, "\n\n"), nil
}

// ResolveThread marks a review thread resolved.
func (p *PRManager) ResolveThread(threadID string) (string, error) {
	if _, err := p.Gh.ResolveReviewThread(threadID); err != nil {
		return "", err
	}
	return fmt.Sprintf("Resolved thread %s", threadID), nil
}

// UnresolveThread reopens a resolved review thread.
func (p *PRManager) UnresolveThread(threadID string) (string, error) {
	if _, err := p.Gh.UnresolveReviewThread(threadID); err != nil {
		return "", err
	}
	return fmt.Sprintf("Reopened thread %s", threadID), nil
}

// ReplyThread posts a reply on a review thread.
func (p *PRManager) ReplyThread(threadID, body string) (string, error) {
	if _, err := p.Gh.ReplyToReviewThread(threadID, body); err != nil {
		return "", err
	}
	return fmt.Sprintf("Replied to thread %s", threadID), nil
}

// reviewEventForVerdict maps a friendly verdict to the GitHub review event the
// reviews endpoint expects.
func reviewEventForVerdict(verdict string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(verdict)) {
	case "approve":
		return "APPROVE", nil
	case "request-changes", "request_changes":
		return "REQUEST_CHANGES", nil
	case "comment":
		return "COMMENT", nil
	default:
		return "", fmt.Errorf(
			"unknown review verdict %q (use approve, request-changes, or comment)",
			verdict,
		)
	}
}

// buildReviewPayload assembles the JSON request body for the reviews endpoint.
// commentsJSON, when non-empty, must be a JSON array of inline-comment objects
// ({path, line, body} plus optional start_line/side); it is embedded verbatim
// so callers keep full control of the GitHub comment shape. A body may be empty
// only for an APPROVE — REQUEST_CHANGES and COMMENT need a body or inline
// comments, matching GitHub's own requirement.
func buildReviewPayload(event, body, commentsJSON string) (string, error) {
	payload := struct {
		Body     string          `json:"body,omitempty"`
		Event    string          `json:"event"`
		Comments json.RawMessage `json:"comments,omitempty"`
	}{Event: event, Body: body}

	if strings.TrimSpace(commentsJSON) != "" {
		if !json.Valid([]byte(commentsJSON)) {
			return "", fmt.Errorf("comments must be valid JSON (an array of comment objects)")
		}
		payload.Comments = json.RawMessage(commentsJSON)
	}

	if strings.TrimSpace(body) == "" && len(payload.Comments) == 0 && event != "APPROVE" {
		return "", fmt.Errorf("a %s review requires a body or inline comments", event)
	}

	out, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to build review payload: %w", err)
	}
	return string(out), nil
}

// SubmitReview posts one PR review (approve, request-changes, or comment) in a
// single submission, optionally with inline comments anchored to diff lines.
// verdict is the friendly form; commentsJSON is an optional JSON array of inline
// comments. This is the line-anchored path the /review-pr skill uses, distinct
// from ApprovePR/RequestChangesPR/CommentPR which carry only a body.
func (p *PRManager) SubmitReview(prNumber, verdict, body, commentsJSON string) (string, error) {
	event, err := reviewEventForVerdict(verdict)
	if err != nil {
		return "", err
	}
	// Build (and validate) the payload before any gh call, so a malformed
	// request fails fast without touching the network.
	payload, err := buildReviewPayload(event, body, commentsJSON)
	if err != nil {
		return "", err
	}
	owner, name, pr, err := p.resolveOwnerRepoPR(prNumber)
	if err != nil {
		return "", err
	}
	if _, err := p.Gh.CreateReview(owner, name, pr, payload); err != nil {
		return "", err
	}
	return fmt.Sprintf("Submitted %s review on %s", verdictLabel(event), prLabel(prNumber)), nil
}

// verdictLabel renders the API event back as its friendly verdict for messages.
func verdictLabel(event string) string {
	switch event {
	case "APPROVE":
		return "approve"
	case "REQUEST_CHANGES":
		return "request-changes"
	default:
		return "comment"
	}
}

// CreatePR opens a PR and returns its URL.
func (p *PRManager) CreatePR(title, body, base string) (string, error) {
	return p.Gh.CreatePR(title, body, base)
}

// UpdatePRDescription replaces a PR's body.
func (p *PRManager) UpdatePRDescription(prNumber, body string) (string, error) {
	if err := p.Gh.UpdatePRDescription(prNumber, body); err != nil {
		return "", err
	}
	return "Updated PR description for " + prLabel(prNumber), nil
}

// ApprovePR approves a PR.
func (p *PRManager) ApprovePR(prNumber, body string) (string, error) {
	if err := p.Gh.ApprovePR(prNumber, body); err != nil {
		return "", err
	}
	return "Approved " + prLabel(prNumber), nil
}

// RequestChangesPR requests changes on a PR.
func (p *PRManager) RequestChangesPR(prNumber, body string) (string, error) {
	if err := p.Gh.RequestChangesPR(prNumber, body); err != nil {
		return "", err
	}
	return "Requested changes on " + prLabel(prNumber), nil
}

// RequestReviewPR re-requests review from the given reviewers by adding them
// back to the PR's requested-reviewers list.
func (p *PRManager) RequestReviewPR(prNumber string, reviewers []string) (string, error) {
	if err := p.Gh.RequestReviewPR(prNumber, reviewers); err != nil {
		return "", err
	}
	return fmt.Sprintf(
		"Requested review from %s on %s",
		strings.Join(reviewers, ", "),
		prLabel(prNumber),
	), nil
}

// CommentPR posts a top-level comment on a PR.
func (p *PRManager) CommentPR(prNumber, body string) (string, error) {
	if err := p.Gh.CommentPR(prNumber, body); err != nil {
		return "", err
	}
	return "Commented on " + prLabel(prNumber), nil
}

// MergePR merges a PR with the given strategy.
func (p *PRManager) MergePR(prNumber, method string) (string, error) {
	if err := p.Gh.MergePR(prNumber, method); err != nil {
		return "", err
	}
	return "Merged " + prLabel(prNumber), nil
}

// PRView returns a compact summary of a PR's metadata.
func (p *PRManager) PRView(prNumber string) (string, error) {
	raw, err := p.Gh.PRView(prNumber)
	if err != nil {
		return "", err
	}
	return p.Jq.FormatPRView(raw)
}

// PRChecks returns a compact, one-line-per-check CI status for a PR.
func (p *PRManager) PRChecks(prNumber string) (string, error) {
	raw, err := p.Gh.PRChecks(prNumber)
	if err != nil {
		return "", err
	}
	return p.Jq.FormatPRChecks(raw)
}

// CurrentPR returns the PR number for the current branch.
func (p *PRManager) CurrentPR() (string, error) {
	n, err := p.Gh.CurrentPRNumber()
	if err != nil {
		return "", err
	}
	if n == "" {
		return "No pull request found for the current branch.", nil
	}
	return n, nil
}

// CurrentRepo returns the current repository as "owner/name".
func (p *PRManager) CurrentRepo() (string, error) {
	return p.Gh.CurrentRepo()
}

// prLabel describes the PR target for confirmation messages.
func prLabel(prNumber string) string {
	if prNumber == "" {
		return "the current branch's PR"
	}
	return "PR #" + prNumber
}
