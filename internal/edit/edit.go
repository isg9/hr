// Package edit applies metadata changes (title, date) to an article and
// renames the file, moving its .meta.toml sidecar and raw HTML stash so
// the vault's naming convention stays consistent.
package edit

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/isdg/hr/internal/article"
	"github.com/isdg/hr/internal/meta"
	"github.com/isdg/hr/internal/textfmt"
	"github.com/isdg/hr/internal/vault"
)

// Options specifies the changes to apply; nil fields are left unchanged.
type Options struct {
	Title *string
	Date  *time.Time
}

// Article applies opts to the article at path. When the change alters
// the canonical filename, the .md file, its .meta.toml sidecar, and its
// raw HTML stash are renamed together — the stable id is preserved.
// Returns the resulting path.
func Article(v *vault.Vault, path string, opts Options) (string, error) {
	if opts.Title == nil && opts.Date == nil {
		return "", fmt.Errorf("nothing to edit: pass --title and/or --date")
	}
	if !strings.HasSuffix(path, ".md") {
		return "", fmt.Errorf("not a .md article: %s", path)
	}
	fm, body, err := article.ReadFile(path)
	if err != nil {
		return "", err
	}
	if opts.Title != nil {
		fm.Title = textfmt.Line(*opts.Title)
	}
	if opts.Date != nil {
		fm.Published = opts.Date.Format(time.RFC3339)
	}
	if err := article.Rewrite(path, fm, body); err != nil {
		return "", err
	}

	id, ok := article.IDFromName(filepath.Base(path))
	if !ok {
		return path, nil // non-canonical name: edited in place, no rename
	}
	newName := article.NameFor(fm.Title, article.ParseTime(fm.Published), id)
	if newName == filepath.Base(path) {
		return path, nil
	}
	newPath := filepath.Join(filepath.Dir(path), newName)
	if _, err := os.Stat(newPath); err == nil {
		return "", fmt.Errorf("target already exists: %s", newPath)
	}
	if err := os.Rename(path, newPath); err != nil {
		return "", err
	}
	if err := moveIfExists(meta.Path(path), meta.Path(newPath)); err != nil {
		return newPath, fmt.Errorf("renamed article but sidecar move failed: %w", err)
	}
	if v != nil {
		err := moveIfExists(
			v.RawPath(fm.Feed, filepath.Base(path)),
			v.RawPath(fm.Feed, newName),
		)
		if err != nil {
			return newPath, fmt.Errorf("renamed article but raw move failed: %w", err)
		}
	}
	return newPath, nil
}

// moveIfExists renames src→dst when src exists, creating dst's parent.
// A missing src is a no-op.
func moveIfExists(src, dst string) error {
	if _, err := os.Stat(src); errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	return os.Rename(src, dst)
}
