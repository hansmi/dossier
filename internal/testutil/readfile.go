package testutil

import (
	"io/fs"
	"testing"
)

func MustReadFileString(t *testing.T, fsys fs.FS, path string) string {
	t.Helper()

	content, err := fs.ReadFile(fsys, path)
	if err != nil {
		t.Errorf("ReadFile(%q) failed: %v", path, err)
	}

	return string(content)
}
