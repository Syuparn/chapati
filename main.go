package main

import (
	"fmt"
	"os"

	"github.com/syuparn/chapati/di"
	"github.com/syuparn/chapati/interface/controller"
)

func main() {
	container := di.NewContainer(os.Stdout)
	derr := container.Invoke(func(c controller.CurryFunctionController) {
		if err := c.Handle("example/example.go"); err != nil {
			fmt.Fprintf(os.Stderr, "%+v\n", err)
			os.Exit(1)
		}
	})

	if derr != nil {
		fmt.Fprintf(os.Stderr, "di failed: %+v\n", derr)
		os.Exit(1)
	}
}
