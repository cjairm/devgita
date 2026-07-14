package commands

import (
	"testing"
	"time"
)

func TestCommandTimeoutContext(t *testing.T) {
	t.Run("zero timeout has no deadline", func(t *testing.T) {
		ctx, cancel := commandTimeoutContext(0)
		defer cancel()
		if _, ok := ctx.Deadline(); ok {
			t.Fatal("expected no deadline for zero timeout")
		}
		if err := ctx.Err(); err != nil {
			t.Fatalf("expected no error on a fresh context, got %v", err)
		}
	})

	t.Run("positive timeout sets a deadline within bound", func(t *testing.T) {
		ctx, cancel := commandTimeoutContext(50 * time.Millisecond)
		defer cancel()
		deadline, ok := ctx.Deadline()
		if !ok {
			t.Fatal("expected a deadline for a positive timeout")
		}
		if time.Until(deadline) > 50*time.Millisecond {
			t.Fatalf("deadline further out than the requested timeout: %v", deadline)
		}
	})

	t.Run("negative timeout treated as unbounded", func(t *testing.T) {
		ctx, cancel := commandTimeoutContext(-1 * time.Second)
		defer cancel()
		if _, ok := ctx.Deadline(); ok {
			t.Fatal("expected no deadline for a negative timeout")
		}
	})
}
