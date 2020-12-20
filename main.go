package main

import (
	"fmt"
	"os"
)

func main() {
	fileAST, fset, err := ReadFileAST("extract_func.go")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
	}

	funcNodes, err := ExtractFuncAst(fileAST, fset)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
	}

	for _, f := range funcNodes {
		if f.IsExported() {
			continue
		}

		fSpec, err := f.Spec()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%+v\n", err)
		}

		src, err := GenCurry(fSpec)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%+v\n", err)
		}

		fmt.Println(src)
	}
}
