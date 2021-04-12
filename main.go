package main

import (
	"fmt"
	"os"

	"github.com/syuparn/chapati/di"
	"github.com/syuparn/chapati/interface/controller"
)

func main() {
	args, err := parseArgs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	f, err := os.Create(args.OutputFile)
	defer f.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create output file %s %+v\n",
			args.OutputFile, err)
		os.Exit(1)
	}

	container := di.NewContainer(f)
	derr := container.Invoke(func(c controller.CurryFunctionController) {
		if err := c.Handle(args.InputFile); err != nil {
			fmt.Fprintf(os.Stderr, "%+v\n", err)
			os.Exit(1)
		}
	})

	if derr != nil {
		fmt.Fprintf(os.Stderr, "di failed: %+v\n", derr)
		os.Exit(1)
	}
}
