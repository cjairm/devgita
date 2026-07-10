/*
* Copyright © 2025 Carlos Mendez <carlos@hadaelectronics.com> | https://cjairm.me/
 */
package cmd

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/internal/config"
	"github.com/cjairm/devgita/internal/inventory"
	tuiinventory "github.com/cjairm/devgita/internal/tui/inventory"
	"github.com/spf13/cobra"
)

var (
	validateCategoryFlag string
	validatePlainFlag    bool
)

// formatValidate renders items as one STATUS table per non-empty category (or
// just the matching one, if category is set). Returns the rendered text,
// whether any item is StateMissing (drives dg validate's exit code — UNKNOWN
// never does), and a category-validation error if any.
func formatValidate(items []inventory.Item, category string) (string, bool, error) {
	if category != "" && !isValidCategory(category) {
		return "", false, fmt.Errorf(
			"invalid category %q: valid categories are %s",
			category, strings.Join(validCategoryKeys(), ", "),
		)
	}

	byCategory := map[string][]inventory.Item{}
	for _, it := range items {
		if category != "" && it.Category != category {
			continue
		}
		byCategory[it.Category] = append(byCategory[it.Category], it)
	}

	var buf bytes.Buffer
	anyMissing := false
	wrote := false
	for _, cat := range inventory.Categories {
		rows := byCategory[cat.Key]
		if len(rows) == 0 {
			continue
		}
		sort.Slice(rows, func(i, j int) bool { return rows[i].Name < rows[j].Name })

		fmt.Fprintf(&buf, "%s:\n", cat.Label)
		w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tSTATUS\tSOURCE\tDETAIL")
		for _, it := range rows {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", it.Name, it.State.String(), it.Source, it.Detail)
			if it.State == inventory.StateMissing {
				anyMissing = true
			}
		}
		_ = w.Flush()
		fmt.Fprintln(&buf)
		wrote = true
	}

	if !wrote {
		return "Nothing tracked yet. Run `dg install` to get started.\n", false, nil
	}
	return buf.String(), anyMissing, nil
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Verify tracked installations are still present on the system",
	Long: `Verify tracked installations are still present on the system.

For every item devgita tracked (installed by devgita, or found pre-existing),
checks whether it's still actually present — catching drift between
global_config.yaml and system reality.

In a terminal, opens the shared inventory dashboard (same as 'dg list') pre-
filtered to problems only. Piped output, CI, or --plain get a plain STATUS
table and a non-zero exit code if anything tracked is missing.

Examples:
  dg validate                    # Interactive dashboard, problems only
  dg validate --plain            # Plain STATUS table, exits 1 if anything missing
  dg validate --category=fonts   # Limit to one category`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if validateCategoryFlag != "" && !isValidCategory(validateCategoryFlag) {
			return fmt.Errorf(
				"invalid category %q: valid categories are %s",
				validateCategoryFlag, strings.Join(validCategoryKeys(), ", "),
			)
		}

		gc := &config.GlobalConfig{}
		if err := gc.Load(); err != nil {
			return fmt.Errorf("failed to load global config: %w", err)
		}

		if !validatePlainFlag && isInteractiveTerminal() {
			return tuiinventory.Run(gc, tuiinventory.Options{
				ProblemsOnly: true,
				Category:     validateCategoryFlag,
			})
		}

		c := &inventory.Collector{Cmd: commands.NewCommand(), Base: commands.NewBaseCommand()}
		items := c.Collect(gc)

		out, anyMissing, err := formatValidate(items, validateCategoryFlag)
		if err != nil {
			return err
		}
		fmt.Print(out)
		if anyMissing {
			return fmt.Errorf("drift detected: one or more tracked items are missing")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	validateCmd.Flags().StringVar(
		&validateCategoryFlag,
		"category",
		"",
		fmt.Sprintf("Filter to a single category (%s)", strings.Join(validCategoryKeys(), ", ")),
	)
	validateCmd.Flags().BoolVar(
		&validatePlainFlag,
		"plain",
		false,
		"Force plain-text output even in a terminal",
	)
}
