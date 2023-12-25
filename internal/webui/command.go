package webui

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"runtime"

	"github.com/google/subcommands"
)

type Command struct {
	listenAddr string
	serverOpts serverOptions
}

func (*Command) Name() string {
	return "web"
}

func (*Command) Synopsis() string {
	return `Start an HTTP server to render a sketch.`
}

func (c *Command) Usage() string {
	return `Arguments: ` + c.Name() + ` <document_file> <sketch_file>

Flags:
`
}

func (c *Command) SetFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.listenAddr, "listen", "[::1]:8080",
		"TCP network address for listening.")
	fs.IntVar(&c.serverOpts.maxConcurrent, "max_concurrent", runtime.NumCPU(),
		"Maximum number of concurrently processed requests.")
	fs.IntVar(&c.serverOpts.maxPages, "max_pages", 10,
		"Maximum number of pages to parse.")
}

func (c *Command) execute(ctx context.Context) error {
	ln, err := net.Listen("tcp", c.listenAddr)
	if err != nil {
		return err
	}

	log.Printf("HTTP server listening on http://%s", ln.Addr())

	s, err := newServer(c.serverOpts)
	if err != nil {
		return err
	}

	return http.Serve(ln, s.makeRouter())
}

func (c *Command) Execute(ctx context.Context, fs *flag.FlagSet, args ...any) subcommands.ExitStatus {
	if fs.NArg() != 2 {
		fs.Usage()
		return subcommands.ExitUsageError
	}

	c.serverOpts.documentPath = fs.Arg(0)
	c.serverOpts.sketchPath = fs.Arg(1)

	if err := c.execute(ctx); err != nil {
		log.Printf("Error: %v", err)
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
