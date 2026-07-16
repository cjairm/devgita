package tuicomponents

import "strings"

// FuzzyRank scores how well a query matches a candidate string. Higher ranks
// are better matches; FuzzyNoMatch (zero value) means the query does not
// match at all. Values are ordered so a plain int comparison sorts matches
// best-first: exact prefix > substring > bare subsequence.
type FuzzyRank int

const (
	FuzzyNoMatch FuzzyRank = iota
	FuzzySubsequence
	FuzzySubstring
	FuzzyExactPrefix
)

// FuzzyMatch case-insensitively scores query against candidate. A subsequence
// match requires every rune of query to appear in candidate in order (not
// necessarily contiguous); substring and prefix matches are strictly stronger
// forms of subsequence so they're checked first. An empty query matches every
// candidate at the weakest rank (FuzzySubsequence), so an unfiltered list
// keeps its original order (see refilter's stable sort).
func FuzzyMatch(query, candidate string) FuzzyRank {
	q := strings.ToLower(query)
	c := strings.ToLower(candidate)

	if q == "" {
		return FuzzySubsequence
	}
	if strings.HasPrefix(c, q) {
		return FuzzyExactPrefix
	}
	if strings.Contains(c, q) {
		return FuzzySubstring
	}
	if isSubsequence(q, c) {
		return FuzzySubsequence
	}
	return FuzzyNoMatch
}

// isSubsequence reports whether every rune of q appears in c in order.
// q and c are assumed already lowercased by the caller.
func isSubsequence(q, c string) bool {
	qr := []rune(q)
	qi := 0
	for _, r := range c {
		if qi >= len(qr) {
			break
		}
		if r == qr[qi] {
			qi++
		}
	}
	return qi == len(qr)
}
