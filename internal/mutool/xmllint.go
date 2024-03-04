package mutool

import (
	"context"
	"fmt"
	"slices"
)

type recoverArgs struct {
	input  string
	output string
}

type xmllintInvoker interface {
	CheckCommand(context.Context) error
	Recover(context.Context, recoverArgs) error
}

type xmllintCommand struct {
	args []string
}

var _ xmllintInvoker = (*xmllintCommand)(nil)

func (c *xmllintCommand) makeArgs(args ...string) []string {
	return append(slices.Clone(c.args), args...)
}

func (c *xmllintCommand) CheckCommand(ctx context.Context) (err error) {
	cmd := makeCommand(ctx, c.makeArgs("--version"))

	if _, err := cmd.Output(); err != nil {
		return fmt.Errorf("xmllint: %w", err)
	}

	return nil
}

func (c *xmllintCommand) Recover(ctx context.Context, a recoverArgs) error {
	cmd := makeCommand(ctx, c.makeArgs("--recover", "--output", a.output, a.input))

	_, err := cmd.Output()

	if err != nil {
		err = fmt.Errorf("xmllint: %w", err)
	}

	return err
}
