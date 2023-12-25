package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"
	"github.com/hansmi/dossier/internal/clianalyzesketch"
	"github.com/hansmi/dossier/internal/webui"
)

func main() {
	for _, cmd := range []subcommands.Command{
		subcommands.HelpCommand(),
		subcommands.FlagsCommand(),
		subcommands.CommandsCommand(),
		&clianalyzesketch.Command{},
		&webui.Command{},
	} {
		subcommands.Register(cmd, "")
	}

	flag.Parse()

	os.Exit(int(subcommands.Execute(context.Background(), nil)))
}
