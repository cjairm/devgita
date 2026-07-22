package task

import (
	"strings"
	"testing"
)

func TestParseActionsJobID(t *testing.T) {
	cases := []struct {
		name   string
		link   string
		wantID string
		wantOK bool
	}{
		{
			name:   "plain job link",
			link:   "https://github.com/octocat/hello/actions/runs/123456/job/789",
			wantID: "789",
			wantOK: true,
		},
		{
			name:   "job link with step fragment",
			link:   "https://github.com/octocat/hello/actions/runs/123456/job/789#step:5:12",
			wantID: "789",
			wantOK: true,
		},
		{
			name:   "job link with trailing slash",
			link:   "https://github.com/octocat/hello/actions/runs/123456/job/789/",
			wantID: "789",
			wantOK: true,
		},
		{
			name:   "job link with surrounding whitespace",
			link:   "  https://github.com/octocat/hello/actions/runs/123456/job/789  ",
			wantID: "789",
			wantOK: true,
		},
		{
			name:   "empty link (external check / commit status)",
			link:   "",
			wantID: "",
			wantOK: false,
		},
		{
			name:   "external status check link",
			link:   "https://ci.example.com/build/42",
			wantID: "",
			wantOK: false,
		},
		{
			name:   "run link without a job segment",
			link:   "https://github.com/octocat/hello/actions/runs/123456",
			wantID: "",
			wantOK: false,
		},
		{
			name:   "non-numeric job id never matches",
			link:   "https://github.com/octocat/hello/actions/runs/123456/job/abc",
			wantID: "",
			wantOK: false,
		},
		{
			name:   "github.com host but different path shape",
			link:   "https://github.com/octocat/hello/commit/deadbeef/checks",
			wantID: "",
			wantOK: false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotID, gotOK := parseActionsJobID(c.link)
			if gotID != c.wantID || gotOK != c.wantOK {
				t.Fatalf("parseActionsJobID(%q) = (%q, %v), want (%q, %v)",
					c.link, gotID, gotOK, c.wantID, c.wantOK)
			}
		})
	}
}

func TestDigestLogTail(t *testing.T) {
	t.Run("short log under the bound passes through unchanged (minus noise)", func(t *testing.T) {
		raw := "job\tstep\t2026-07-22T18:32:18.0000000Z line one\n" +
			"job\tstep\t2026-07-22T18:32:19.0000000Z line two\n" +
			"job\tstep\t2026-07-22T18:32:20.0000000Z line three\n"
		got := digestLogTail(raw, 10)
		want := "line one\nline two\nline three"
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		if strings.Contains(got, "omitted") {
			t.Fatalf("no truncation receipt expected, got: %q", got)
		}
	})

	t.Run("consecutive repeated lines collapse with a count marker", func(t *testing.T) {
		var b strings.Builder
		for i := 0; i < 12; i++ {
			b.WriteString("job\tstep\t2026-07-22T18:32:1" +
				string(rune('0'+i%10)) + ".0000000Z retrying connection\n")
		}
		b.WriteString("job\tstep\t2026-07-22T18:33:00.0000000Z connection failed: timeout\n")
		got := digestLogTail(b.String(), 10)
		want := "retrying connection  (×12)\nconnection failed: timeout"
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})

	t.Run("long log is truncated to the tail with a truncation receipt", func(t *testing.T) {
		var b strings.Builder
		for i := 1; i <= 100; i++ {
			b.WriteString("job\tstep\t2026-07-22T18:32:18.0000000Z log line ")
			b.WriteString(itoa(i))
			b.WriteString("\n")
		}
		got := digestLogTail(b.String(), 10)
		lines := strings.Split(got, "\n")
		if len(lines) != 10 {
			t.Fatalf("expected exactly 10 lines (receipt + 9 tail lines), got %d: %q",
				len(lines), got)
		}
		if lines[0] != "… 91 earlier lines omitted" {
			t.Fatalf("unexpected receipt: %q", lines[0])
		}
		if lines[1] != "log line 92" || lines[len(lines)-1] != "log line 100" {
			t.Fatalf("unexpected tail content: %q", got)
		}
	})

	t.Run("non-consecutive repeats are not collapsed", func(t *testing.T) {
		raw := "job\tstep\t2026-07-22T18:32:18.0000000Z boom\n" +
			"job\tstep\t2026-07-22T18:32:19.0000000Z something else\n" +
			"job\tstep\t2026-07-22T18:32:20.0000000Z boom\n"
		got := digestLogTail(raw, 10)
		want := "boom\nsomething else\nboom"
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})

	t.Run("embedded tab in message content passes through intact", func(t *testing.T) {
		raw := "job\tstep\t2026-07-22T18:32:18.0000000Z col1\tcol2\tcol3\n"
		got := digestLogTail(raw, 10)
		want := "col1\tcol2\tcol3"
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})

	t.Run("empty raw returns empty string", func(t *testing.T) {
		if got := digestLogTail("", 10); got != "" {
			t.Fatalf("expected empty string, got %q", got)
		}
	})

	t.Run("maxLines <= 0 still returns a receipt, never a crash", func(t *testing.T) {
		raw := "job\tstep\t2026-07-22T18:32:18.0000000Z a\n" +
			"job\tstep\t2026-07-22T18:32:19.0000000Z b\n"
		got := digestLogTail(raw, 0)
		if !strings.HasPrefix(got, "… ") || !strings.Contains(got, "omitted") {
			t.Fatalf("expected a truncation receipt, got %q", got)
		}
	})
}

// itoa avoids pulling in strconv just for this test file's synthetic fixtures.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	digits := []byte{}
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}
