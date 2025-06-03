package main

import (
	"errors"
	"flag"
	"log/slog"
	"os"
	"github.com/genshen/cmds"
	_ "github.com/genshen/wssocks/cmd/client"
	_ "github.com/genshen/wssocks/cmd/server"
	_ "github.com/genshen/wssocks/version"
)

func init() {
	// Configure slog with debug level
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewTextHandler(os.Stdout, opts)
	slog.SetDefault(slog.New(handler))
}

func main() {
	cmds.SetProgramName("wssocks")
	if err := cmds.Parse(); err != nil {
		if !errors.Is(err, flag.ErrHelp) && !errors.Is(err, &cmds.SubCommandParseError{}) {
			slog.Error("parsing error", "error", err)
			os.Exit(1)
		}
	}
}
