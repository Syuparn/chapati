package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"golang.org/x/xerrors"
)

var (
	outputFile = flag.String("o", "", "output file name (default: 'generate.curried.{input file name}.go')")
)

type CmdArgs struct {
	InputFile  string
	OutputFile string
}

func parseArgs() (*CmdArgs, error) {
	prependUsage("chapati [options] <inputfile>\n\n")

	flag.Parse()

	if flag.Arg(0) == "" {
		return nil, xerrors.Errorf("input file name must not be empty")
	}

	in := flag.Arg(0)
	out := filepath.Join(filepath.Dir(in), DefaultOutputFilePrefix+filepath.Base(in))

	if *outputFile != "" {
		out = *outputFile
	}

	return &CmdArgs{
		InputFile:  in,
		OutputFile: out,
	}, nil
}

func prependUsage(msg string) {
	origUsage := flag.Usage
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), msg)
		origUsage()
	}
}
