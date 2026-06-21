package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/isdg/hr/internal/listing"
)

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Inspect or rebuild the listing cache",
	Long: `Inspect or manage the listing cache (.hr/index.json).

The cache mirrors every article's metadata so 'hr list' is fast; it
refreshes automatically as files change. Use these commands when you
want to see its state or force a refresh.

With no subcommand, prints the cache's path, version, and entry count.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		v, _, err := openActiveVault()
		if err != nil {
			return err
		}
		s := listing.StatIndex(v)
		status := "up to date"
		switch {
		case !s.Exists:
			status = "missing (builds on next list)"
		case s.Version != s.WantVersion:
			status = fmt.Sprintf(
				"stale: on-disk v%d, want v%d (rebuilds on next list)",
				s.Version, s.WantVersion)
		}
		fmt.Printf("  path:    %s\n", s.Path)
		fmt.Printf("  version: %d\n", s.WantVersion)
		fmt.Printf("  entries: %d\n", s.Entries)
		fmt.Printf("  status:  %s\n", status)
		return nil
	},
}

var indexRebuildCmd = &cobra.Command{
	Use:   "rebuild",
	Short: "Discard the cache and re-parse every article",
	RunE: func(cmd *cobra.Command, args []string) error {
		v, _, err := openActiveVault()
		if err != nil {
			return err
		}
		n, err := listing.RebuildIndex(v)
		if err != nil {
			return err
		}
		fmt.Printf("rebuilt %d articles\n", n)
		return nil
	},
}

var indexClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Delete the cache (rebuilds on next list)",
	RunE: func(cmd *cobra.Command, args []string) error {
		v, _, err := openActiveVault()
		if err != nil {
			return err
		}
		if err := listing.ClearIndex(v); err != nil {
			return err
		}
		fmt.Println("removed index cache")
		return nil
	},
}

func init() {
	indexCmd.AddCommand(indexRebuildCmd)
	indexCmd.AddCommand(indexClearCmd)
	rootCmd.AddCommand(indexCmd)
}
