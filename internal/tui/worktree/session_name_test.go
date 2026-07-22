package tuiworktree

import (
	"strings"
	"testing"
)

func TestNextFreeSessionNamePicksFirstInOrder(t *testing.T) {
	// Order [0,1,2,...] with an empty taken-set must return the first
	// character in dragonBallNames, prefixed.
	order := seqOrder(len(dragonBallNames))
	got := nextFreeSessionName(map[string]bool{}, order)
	want := sessionNamePrefix + dragonBallNames[0]
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestNextFreeSessionNameSkipsTaken(t *testing.T) {
	order := seqOrder(len(dragonBallNames))
	taken := map[string]bool{
		sessionNamePrefix + dragonBallNames[0]: true,
		sessionNamePrefix + dragonBallNames[1]: true,
	}
	got := nextFreeSessionName(taken, order)
	want := sessionNamePrefix + dragonBallNames[2]
	if got != want {
		t.Errorf("expected the first free character %q, got %q", want, got)
	}
}

func TestNextFreeSessionNameNeverReturnsTaken(t *testing.T) {
	order := seqOrder(len(dragonBallNames))
	// Every bare character is taken — the numeric-suffix fallback must kick in
	// and still return something absent from taken.
	taken := map[string]bool{}
	for _, n := range dragonBallNames {
		taken[sessionNamePrefix+n] = true
	}
	got := nextFreeSessionName(taken, order)
	if taken[got] {
		t.Errorf("fallback returned a taken name %q", got)
	}
	if !strings.HasPrefix(got, sessionNamePrefix+dragonBallNames[0]+"-") {
		t.Errorf("expected a numeric-suffixed fallback on the first character, got %q", got)
	}
}

func TestNextFreeSessionNameFallbackIncrements(t *testing.T) {
	order := seqOrder(len(dragonBallNames))
	base := sessionNamePrefix + dragonBallNames[0]
	taken := map[string]bool{}
	for _, n := range dragonBallNames {
		taken[sessionNamePrefix+n] = true
	}
	// -2 also taken → must land on -3.
	taken[base+"-2"] = true
	got := nextFreeSessionName(taken, order)
	if got != base+"-3" {
		t.Errorf("expected %q, got %q", base+"-3", got)
	}
}

func TestRandomSessionNameOrderIsAPermutation(t *testing.T) {
	order := randomSessionNameOrder()
	if len(order) != len(dragonBallNames) {
		t.Fatalf(
			"expected an order covering all %d names, got %d",
			len(dragonBallNames),
			len(order),
		)
	}
	seen := make(map[int]bool, len(order))
	for _, i := range order {
		if i < 0 || i >= len(dragonBallNames) {
			t.Fatalf("index %d out of range", i)
		}
		if seen[i] {
			t.Fatalf("index %d repeated — not a permutation", i)
		}
		seen[i] = true
	}
}

// seqOrder returns [0,1,2,...,n-1], a deterministic order for tests.
func seqOrder(n int) []int {
	order := make([]int, n)
	for i := range order {
		order[i] = i
	}
	return order
}
