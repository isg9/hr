package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var vaultFlag string

var rootCmd = &cobra.Command{
	Use:   "hr",
	Short: "File-first RSS / blog reader",
	Long: `hr syncs feeds into a directory of plain markdown files you can read
in any editor. Articles are immutable; per-article state (read, favorite,
tags) lives in sidecar .meta.toml files. The whole vault is git-syncable.

Running 'hr' with no subcommand opens the reading panel in nvim
(override the editor with $HR_EDITOR).`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return openInEditor()
	},
}

// openInEditor launches the configured editor with the panel auto-
// opened. nvim is the default; the editor command must accept
// `-c 'lua require("hr").open()'` (nvim-style). Override with
// $HR_EDITOR for a different binary or alias.
func openInEditor() error {
	v, _, err := openActiveVault()
	if err != nil {
		return err
	}
	editor := os.Getenv("HR_EDITOR")
	if editor == "" {
		editor = "nvim"
	}
	c := exec.Command(
		editor, "-c", `lua require("hr").start()`)
	c.Dir = v.Root
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "hr:", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&vaultFlag, "vault", "C", "",
		"vault directory (default: $HR_VAULT or ~/blogs)")
}
