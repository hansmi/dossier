package mutool

import (
	"os"
)

func withTempdir() (string, func() error, error) {
	path, err := os.MkdirTemp("", "mutool*")
	if err != nil {
		return "", nil, err
	}

	return path, func() error {
		return os.RemoveAll(path)
	}, nil
}
