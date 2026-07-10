package cmd

import (
	"strings"
	"testing"

	"github.com/cjairm/devgita/internal/inventory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatValidate_EmptyItems(t *testing.T) {
	out, anyMissing, err := formatValidate(nil, "")
	require.NoError(t, err)
	assert.False(t, anyMissing)
	assert.Contains(t, out, "Nothing tracked yet")
}

func TestFormatValidate_TableHasStatusColumn(t *testing.T) {
	items := []inventory.Item{
		{Name: "git", Category: "packages", Source: "installed", State: inventory.StateOK},
		{Name: "tmux", Category: "packages", Source: "installed", State: inventory.StateMissing},
	}
	out, anyMissing, err := formatValidate(items, "")
	require.NoError(t, err)
	assert.True(t, anyMissing)
	assert.Contains(t, out, "STATUS")
	assert.Contains(t, out, "OK")
	assert.Contains(t, out, "MISSING")
	assert.Contains(t, out, "git")
	assert.Contains(t, out, "tmux")
}

func TestFormatValidate_UnknownDetailIsSurfaced(t *testing.T) {
	items := []inventory.Item{
		{
			Name:     "JetBrainsMono",
			Category: "fonts",
			Source:   "installed",
			State:    inventory.StateUnknown,
			Detail:   "fc-list: not found",
		},
	}
	out, _, err := formatValidate(items, "")
	require.NoError(t, err)
	assert.Contains(t, out, "DETAIL")
	assert.Contains(
		t,
		out,
		"fc-list: not found",
		"the check error should be visible so a user can diagnose an UNKNOWN result, not just see the bare status",
	)
}

func TestFormatValidate_UnknownDoesNotSetAnyMissing(t *testing.T) {
	items := []inventory.Item{
		{
			Name:     "JetBrainsMono",
			Category: "fonts",
			Source:   "installed",
			State:    inventory.StateUnknown,
		},
	}
	_, anyMissing, err := formatValidate(items, "")
	require.NoError(t, err)
	assert.False(t, anyMissing, "UNKNOWN must never fail dg validate's exit code")
}

func TestFormatValidate_CategoryFilter(t *testing.T) {
	items := []inventory.Item{
		{Name: "git", Category: "packages", Source: "installed", State: inventory.StateOK},
		{Name: "JetBrainsMono", Category: "fonts", Source: "installed", State: inventory.StateOK},
	}
	out, _, err := formatValidate(items, "fonts")
	require.NoError(t, err)
	assert.Contains(t, out, "JetBrainsMono")
	assert.NotContains(t, out, "git")
}

func TestFormatValidate_InvalidCategory(t *testing.T) {
	_, _, err := formatValidate(nil, "bogus")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid category")
}

func TestFormatValidate_GroupedByCategoryLabel(t *testing.T) {
	items := []inventory.Item{
		{Name: "git", Category: "packages", Source: "installed", State: inventory.StateOK},
	}
	out, _, err := formatValidate(items, "")
	require.NoError(t, err)
	assert.Contains(t, out, "Packages:")
	packagesIdx := strings.Index(out, "Packages:")
	require.NotEqual(t, -1, packagesIdx)
}
