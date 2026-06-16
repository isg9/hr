// Package textfmt normalizes strings that need to live on a single
// line: frontmatter titles, sidecar aliases, list/TSV columns, nvim
// buffer rows.
package textfmt

import (
	"strings"
	"unicode"
)

// Line collapses s into a single-line, whitespace-normalized string.
// Newlines, tabs, and control characters become a single space; runs
// of whitespace are collapsed; leading and trailing whitespace is
// trimmed. Non-ASCII letters and punctuation are preserved.
func Line(s string) string {
	if s == "" {
		return ""
	}
	var b strings.Builder
	b.Grow(len(s))
	inSpace := true
	for _, r := range s {
		if r == '\t' || r == '\n' || r == '\r' ||
			(unicode.IsControl(r) && !unicode.IsSpace(r)) {
			if !inSpace {
				b.WriteByte(' ')
				inSpace = true
			}
			continue
		}
		if unicode.IsSpace(r) {
			if !inSpace {
				b.WriteByte(' ')
				inSpace = true
			}
			continue
		}
		b.WriteRune(r)
		inSpace = false
	}
	return strings.TrimRight(b.String(), " ")
}
