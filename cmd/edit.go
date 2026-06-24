package cmd

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/isdg/hr/internal/edit"
)

var (
	editTitle string
	editDate  string
)

var editCmd = &cobra.Command{
	Use:   "edit <article.md>",
	Short: "Edit an article's title and/or date (renames the file to match)",
	Long: `Edit an article's frontmatter title and/or published date.

The .md file is renamed to its canonical "<date>-<slug>-<id>.md" form,
and its .meta.toml sidecar and raw HTML stash move with it. The stable
id is preserved, so dedup still works on the next sync.

Dates accept YYYY-MM-DD, -YYYY-MM-DD (BC), or full RFC3339.

  hr edit feeds/aristotle/...md --title "Poetics" --date -0350-01-01`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		v, _, err := openActiveVault()
		if err != nil {
			return err
		}
		var opts edit.Options
		if cmd.Flags().Changed("title") {
			opts.Title = &editTitle
		}
		if cmd.Flags().Changed("date") {
			d, err := parseDateFlag(editDate)
			if err != nil {
				return err
			}
			opts.Date = &d
		}
		newPath, err := edit.Article(v, args[0], opts)
		if err != nil {
			return err
		}
		if newPath != args[0] {
			fmt.Printf("%s → %s\n",
				filepath.Base(args[0]), filepath.Base(newPath))
		} else {
			fmt.Printf("updated %s\n", filepath.Base(newPath))
		}
		return nil
	},
}

// parseDateFlag parses YYYY-MM-DD, -YYYY-MM-DD (BC), or RFC3339 into a
// UTC time.
func parseDateFlag(s string) (time.Time, error) {
	neg := strings.HasPrefix(s, "-")
	body := strings.TrimPrefix(s, "-")
	for _, layout := range []string{"2006-01-02", time.RFC3339} {
		if t, err := time.Parse(layout, body); err == nil {
			if neg {
				t = time.Date(-t.Year(), t.Month(), t.Day(),
					t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.UTC)
			}
			return t.UTC(), nil
		}
	}
	return time.Time{}, fmt.Errorf(
		"invalid date %q (use YYYY-MM-DD, -YYYY-MM-DD, or RFC3339)", s)
}

func init() {
	editCmd.Flags().StringVar(&editTitle, "title", "", "new title")
	editCmd.Flags().StringVar(&editDate, "date", "",
		"new published date (YYYY-MM-DD, -YYYY-MM-DD for BC, or RFC3339)")
	rootCmd.AddCommand(editCmd)
}
