package mutool

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"go.uber.org/multierr"
	"golang.org/x/exp/slices"
)

func makeCommand(ctx context.Context, args []string) *exec.Cmd {
	return exec.CommandContext(ctx, args[0], args[1:]...)
}

type showArgs struct {
	input string
}

func (a *showArgs) build() []string {
	return []string{"show", "-g", "--", a.input}
}

type drawArgs struct {
	input     string
	pageRange string

	output string
	stdout io.Writer

	format string
	width  int
	height int
}

func (a drawArgs) build() []string {
	args := []string{"draw", "-N", "-a",
		"-F", a.format,
		"-o", a.output,
	}

	for _, i := range []struct {
		flag  string
		value int
	}{
		{"-w", a.width},
		{"-h", a.height},
	} {
		if i.value != 0 {
			args = append(args, i.flag, strconv.Itoa(i.value))
		}
	}

	return append(args, "--", a.input, a.pageRange)
}

type mutoolInvoker interface {
	CheckCommand(context.Context) error
	Show(context.Context, showArgs) error
	Draw(context.Context, drawArgs) error
}

type mutoolCommand struct {
	args []string
}

var _ mutoolInvoker = (*mutoolCommand)(nil)

func (c *mutoolCommand) makeArgs(args ...string) []string {
	return append(slices.Clone(c.args), args...)
}

func (c *mutoolCommand) CheckCommand(ctx context.Context) (err error) {
	tmpdir, tmpdirCleanup, err := withTempdir()
	if err != nil {
		return err
	}

	defer multierr.AppendFunc(&err, tmpdirCleanup)

	cmd := makeCommand(ctx, c.makeArgs("create",
		"-o", filepath.Join(tmpdir, "out.pdf"),
		os.DevNull,
	))

	if _, err := cmd.Output(); err != nil {
		return fmt.Errorf("mutool: %w", err)
	}

	return nil
}

func (c *mutoolCommand) Show(ctx context.Context, a showArgs) error {
	cmd := makeCommand(ctx, c.makeArgs(a.build()...))

	if _, err := cmd.Output(); err != nil {
		return fmt.Errorf("mutool: %w", err)
	}

	return nil
}

func (c *mutoolCommand) Draw(ctx context.Context, a drawArgs) error {
	cmd := makeCommand(ctx, c.makeArgs(a.build()...))

	var err error

	if a.stdout == nil {
		_, err = cmd.Output()
	} else {
		cmd.Stdout = a.stdout

		err = cmd.Run()
	}

	if err != nil {
		return fmt.Errorf("mutool: %w", err)
	}

	return nil
}
