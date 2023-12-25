package testutil

import (
	"os"
	"testing"
)

func MustWriteFile(t *testing.T, path string, content []byte) string {
	t.Helper()

	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Errorf("WriteFile(%q) failed: %v", path, err)
	}

	return path
}

func MustWriteFileString(t *testing.T, path string, content string) string {
	t.Helper()

	return MustWriteFile(t, path, []byte(content))
}
