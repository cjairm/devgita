// Auto-generation of a standalone-session name when the user leaves the name
// prompt blank. Names are "devgita-<character>" for a Dragon Ball character —
// a nod to devgita itself, whose name comes from Vegeta — and are always
// checked against the live tmux sessions so a blank-name create can never
// collide with an existing session (tmux new-session -s fails on a duplicate).

package tuiworktree

import (
	"fmt"
	"math/rand"
)

// sessionNamePrefix prefixes every auto-generated session name. Kept separate
// from the character pool so both the primary "devgita-goku" form and the
// "devgita-goku-2" collision fallback share one source of truth.
const sessionNamePrefix = "devgita-"

// dragonBallNames is the pool the auto-namer draws from. Every entry is
// already lowercase and free of spaces, dots, and colons, so "devgita-<name>"
// is a valid tmux session name with no sanitizing needed (tmux treats "." and
// ":" specially in target names). Kept short and recognizable on purpose — a
// generated name should read as a friendly label, not a random string.
var dragonBallNames = []string{
	"goku", "vegeta", "gohan", "piccolo", "trunks", "goten",
	"krillin", "bulma", "roshi", "beerus", "whis", "frieza",
	"cell", "buu", "raditz", "nappa", "tien", "yamcha",
	"videl", "bardock", "broly", "zamasu", "gotenks", "shenron",
}

// nextFreeSessionName returns a session name that is not present in taken,
// preferring "devgita-<character>" for a Dragon Ball character (tried in the
// order given, so callers control randomness/determinism), and falling back to
// "devgita-<character>-<n>" with an incrementing n only if every bare
// character name is already taken. It is guaranteed to return a name absent
// from taken, so the caller's blank-name create can always proceed without
// clashing with a live session.
//
// order indexes into dragonBallNames; callers pass a shuffled order for a
// random pick (see randomSessionNameOrder) or a fixed one in tests. It must be
// non-empty — dragonBallNames is a non-empty package var, so its randomSessionNameOrder
// output always is too.
func nextFreeSessionName(taken map[string]bool, order []int) string {
	for _, i := range order {
		if name := sessionNamePrefix + dragonBallNames[i]; !taken[name] {
			return name
		}
	}
	// Every bare character is already taken (the user has >= len(pool) live
	// sessions named after the pool): disambiguate the first-ordered character
	// with an incrementing numeric suffix. Bounded by taken's size, so it
	// always terminates on a free name rather than looping forever.
	base := sessionNamePrefix + dragonBallNames[order[0]]
	for n := 2; ; n++ {
		if name := fmt.Sprintf("%s-%d", base, n); !taken[name] {
			return name
		}
	}
}

// randomSessionNameOrder returns the indices of dragonBallNames in a random
// order, so successive blank-name creates tend to pick different characters
// rather than always starting from "goku". Go auto-seeds the global rand
// source (1.20+), so this needs no manual seeding. Extracted from
// nextFreeSessionName so tests can supply a deterministic order instead.
func randomSessionNameOrder() []int {
	return rand.Perm(len(dragonBallNames))
}
