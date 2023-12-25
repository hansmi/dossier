package clianalyzesketch

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/subcommands"
	"github.com/hansmi/aurum"
	"github.com/hansmi/dossier"
	"github.com/hansmi/dossier/internal/cliutil"
	"github.com/hansmi/dossier/pkg/geometry"
	"github.com/hansmi/dossier/pkg/pagerange"
	"github.com/hansmi/dossier/pkg/sketch"
)

type Command struct {
	maxPages        int
	textProtoFormat bool
	lengthUnit      geometry.LengthUnit

	documentPath string
	sketchPath   string
}

func (*Command) Name() string {
	return "analyze-sketch"
}

func (*Command) Synopsis() string {
	return `Analyze a document using a sketch and print the results.`
}

func (c *Command) Usage() string {
	return `Arguments: ` + c.Name() + ` <document_file> <sketch_file>

Flags:
`
}

func (c *Command) SetFlags(fs *flag.FlagSet) {
	fs.IntVar(&c.maxPages, "max_pages", 0,
		"Maximum number of pages to analyze.")
	fs.BoolVar(&c.textProtoFormat, "textproto", false,
		"Write output using the Protocol Buffer text format instead of JSON.")

	lu := cliutil.NewLengthUnitVar(&c.lengthUnit, geometry.Millimeter)
	fs.Var(lu, "unit", lu.Usage("Length unit for output."))
}

func (c *Command) execute(ctx context.Context) error {
	doc := dossier.NewDocument(c.documentPath)

	if err := doc.Validate(ctx); err != nil {
		return fmt.Errorf("document validation: %w", err)
	}

	sketchBytes, err := os.ReadFile(c.sketchPath)
	if err != nil {
		return fmt.Errorf("reading sketch file: %w", err)
	}

	s, err := sketch.CompileFromTextproto(sketchBytes)
	if err != nil {
		return fmt.Errorf("parsing sketch: %w", err)
	}

	r := pagerange.All

	if c.maxPages != 0 {
		if r, err = pagerange.New(1, c.maxPages); err != nil {
			return err
		}
	}

	report, err := s.AnalyzeDocument(ctx, doc, r)
	if err != nil {
		return fmt.Errorf("analyzing document: %w", err)
	}

	var codec aurum.Codec

	if c.textProtoFormat {
		codec = &aurum.TextProtoCodec{}
	} else {
		jc := &aurum.JSONCodec{}
		jc.ProtoMarshalOptions.EmitUnpopulated = true

		codec = jc
	}

	buf, err := codec.Marshal(report.AsProto(c.lengthUnit))
	if err != nil {
		return fmt.Errorf("marshalling report: %w", err)
	}

	_, err = os.Stdout.WriteString(string(buf))

	return err
}

func (c *Command) Execute(ctx context.Context, fs *flag.FlagSet, args ...any) subcommands.ExitStatus {
	if fs.NArg() != 2 {
		fs.Usage()
		return subcommands.ExitUsageError
	}

	c.documentPath = fs.Arg(0)
	c.sketchPath = fs.Arg(1)

	if err := c.execute(ctx); err != nil {
		log.Printf("Error: %v", err)
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
