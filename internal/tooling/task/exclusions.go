package task

import "path/filepath"

// defaultExclusionPatterns are lockfiles and generated/minified assets that
// are unreviewable diff noise. Matched against a file's basename so they
// exclude at any depth (e.g. packages/app/package-lock.json in a monorepo).
// Deliberately non-exhaustive — Pipfile.lock, mix.lock, Podfile.lock,
// packages.lock.json, etc. are absent by design; the --file bypass on
// BranchDiff and raw `git diff` permissions cover the gaps. Shared by
// ReviewScope (partitions an already-fetched file list) and BranchDiff
// (excludes at the git-diff level via exclusionPathspecs).
var defaultExclusionPatterns = []string{
	"package-lock.json",
	"yarn.lock",
	"pnpm-lock.yaml",
	"bun.lockb",
	"go.sum",
	"Cargo.lock",
	"Gemfile.lock",
	"composer.lock",
	"poetry.lock",
	"uv.lock",
	"*.min.js",
	"*.min.css",
}

// isExcludedPath reports whether path matches one of the default exclusion
// patterns, checked against the basename so nested paths are covered too.
func isExcludedPath(path string) bool {
	base := filepath.Base(path)
	for _, pattern := range defaultExclusionPatterns {
		if ok, _ := filepath.Match(pattern, base); ok {
			return true
		}
	}
	return false
}

// exclusionPathspecs renders defaultExclusionPatterns as git pathspecs that
// exclude matches at any depth, for use as extra `git diff` arguments.
func exclusionPathspecs() []string {
	pathspecs := make([]string, len(defaultExclusionPatterns))
	for i, p := range defaultExclusionPatterns {
		pathspecs[i] = ":(exclude,glob)**/" + p
	}
	return pathspecs
}
