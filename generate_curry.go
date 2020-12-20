package main

import (
	"fmt"
	"strings"

	"github.com/dave/jennifer/jen"
)

// GenCurry returns curried function soruce code of decl.
func GenCurry(fs *FuncSpec) (string, error) {
	// TODO: currifying
	// TODO: pass path info
	// TODO: pass import info
	f := jen.NewFilePath("a.b/c")
	fn := f.Func()

	// function recv if exists
	recv, ok := fs.Recv()
	if ok {
		fn.Params(renderParam(recv))
	}

	// function name
	fn.Id(toExportedName(fs.Name()))

	// function params
	fn.Params(renderParams(fs.Params())...)

	// function return types
	returnTypes := fs.ReturnTypes()
	if len(returnTypes) > 0 {
		fn.Params(renderTypes(returnTypes)...)
	}

	// function body
	fn.Block(
		jen.Qual("a.b/c", fs.Name()).Call(renderParamValues(fs.Params())...),
	)

	return fmt.Sprintf("%#v", f), nil
}

func renderParam(p Param) jen.Code {
	ident := jen.Id(p.Name)

	if p.Type.Prefix != "" {
		ident.Op(p.Type.Prefix)
	}

	ident.Id(p.Type.Name)

	return ident
}

func renderParams(params []Param) []jen.Code {
	rendered := make([]jen.Code, len(params))
	for _, p := range params {
		rendered = append(rendered, renderParam(p))
	}

	return rendered
}

func renderParamValue(p Param) jen.Code {
	return jen.Id(p.Name)
}

func renderParamValues(params []Param) []jen.Code {
	rendered := make([]jen.Code, len(params))
	for _, p := range params {
		rendered = append(rendered, renderParamValue(p))
	}

	return rendered
}

func renderTypes(types []Type) []jen.Code {
	rendered := make([]jen.Code, len(types))
	for _, t := range types {
		rendered = append(rendered, renderType(t))
	}

	return rendered
}

func renderType(t Type) jen.Code {
	if t.Prefix == "" {
		return jen.Id(t.Name)

	}

	return jen.Op(t.Prefix).Id(t.Name)
}

func toExportedName(name string) string {
	return strings.ToUpper(name[:1]) + name[1:]
}
