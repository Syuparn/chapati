package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"golang.org/x/xerrors"
)

// ReadFileAST returns the AST of the source file.
func ReadFileAST(filename string) (*ast.File, *token.FileSet, error) {
	fset := token.NewFileSet()
	//                                         src, parse all
	f, err := parser.ParseFile(fset, filename, nil, parser.Mode(0))
	if err != nil {
		return nil, nil, xerrors.Errorf("failed to parse %s: %v", filename, err)
	}

	return f, fset, nil
}

// ExtractFuncAst returns nodes of function definitions in AST.
func ExtractFuncAst(f *ast.File, fset *token.FileSet) ([]*FuncDecl, error) {
	decls := []*FuncDecl{}
	for _, d := range f.Decls {
		if f, ok := d.(*ast.FuncDecl); ok {
			decls = append(decls, NewFuncDecl(f, fset))
		}
	}

	return decls, nil
}

// FuncDecl is a wrapper of ast.FuncDecl
type FuncDecl struct {
	fset *token.FileSet
	ast.FuncDecl
}

// NewFuncDecl creates a new FuncDecl.
func NewFuncDecl(d *ast.FuncDecl, fset *token.FileSet) *FuncDecl {
	return &FuncDecl{fset: fset, FuncDecl: *d}
}

// Spec returns the specification info of the function.
func (d *FuncDecl) Spec() (*FuncSpec, error) {
	name := d.name()
	params, err := d.parameters()
	if err != nil {
		return nil, xerrors.Errorf("failed to get parameters: %v", err)
	}

	returnTypes, err := d.returnTypes()
	if err != nil {
		return nil, xerrors.Errorf("failed to get return types: %v", err)
	}

	recv, ok := d.recv()
	if !ok {
		return NewFuncSpec(name, nil, params, returnTypes), nil
	}
	return NewFuncSpec(name, &recv, params, returnTypes), nil
}

// IsExported returns true if the function is exported.
func (d *FuncDecl) IsExported() bool {
	firstChar := d.name()[:1]
	return strings.ToUpper(firstChar) == firstChar
}

func (d *FuncDecl) name() string {
	return d.FuncDecl.Name.String()
}

func (d *FuncDecl) recv() (Param, bool) {
	recv := d.FuncDecl.Recv
	if recv == nil || recv.NumFields() == 0 {
		return Param{}, false
	}

	param, err := extractParam(recv.List[0])
	if err != nil {
		return Param{}, false
	}

	return param, true
}

func (d *FuncDecl) parameters() ([]Param, error) {
	return extractParams(d.FuncDecl.Type.Params)
}

func (d *FuncDecl) returnTypes() ([]Type, error) {
	return extractTypes(d.FuncDecl.Type.Results)
}

func extractParams(fieldList *ast.FieldList) ([]Param, error) {
	if fieldList == nil || fieldList.NumFields() == 0 {
		return []Param{}, nil
	}

	params := []Param{}
	for i, field := range fieldList.List {
		p, err := extractParam(field)
		if err != nil {
			return []Param{}, xerrors.Errorf(
				"failed to parse param List[%d] (%+V): %v", i, field, err)
		}

		params = append(params, p)
	}

	return params, nil
}

func extractParam(field *ast.Field) (Param, error) {
	t, err := extractType(field.Type)
	if err != nil {
		return Param{}, xerrors.Errorf(
			"failed to parse type %v: %v", field.Type, err)
	}

	if len(field.Names) == 0 {
		return Param{}, xerrors.Errorf(
			"failed to parse name of type %v: %v", field.Type, err)
	}

	return Param{
		Name: field.Names[0].String(),
		Type: t,
	}, nil
}

func extractTypes(fieldList *ast.FieldList) ([]Type, error) {
	if fieldList == nil || fieldList.NumFields() == 0 {
		return []Type{}, nil
	}

	types := []Type{}
	for i, field := range fieldList.List {
		t, err := extractType(field.Type)
		if err != nil {
			return []Type{}, xerrors.Errorf(
				"failed to parse param List[%d] (%+V): %v", i, field, err)
		}

		types = append(types, t)
	}

	return types, nil
}

func extractType(t ast.Expr) (Type, error) {
	// TODO: handle map
	switch t := t.(type) {
	case *ast.Ident:
		return Type{
			Prefix: "",
			Name:   t.Name,
		}, nil
	case *ast.StarExpr:
		inner, err := extractType(t.X)
		if err != nil {
			return Type{}, xerrors.Errorf("type %T is not supported to parse", t.X)
		}

		return Type{
			Prefix: "*" + inner.Prefix,
			Name:   inner.Name,
		}, nil
	case *ast.SelectorExpr:
		inner, err := extractType(t.X)
		if err != nil {
			return Type{}, xerrors.Errorf("type %T is not supported to parse", t.X)
		}

		return Type{
			Prefix: inner.Prefix,
			Name:   fmt.Sprintf("%s.%s", inner.Name, t.Sel.Name),
		}, nil
	case *ast.ArrayType:
		inner, err := extractType(t.Elt)
		if err != nil {
			return Type{}, xerrors.Errorf("type %T is not supported to parse", t.Elt)
		}

		return Type{
			Prefix: "[]" + inner.Prefix,
			Name:   inner.Name,
		}, nil
	case *ast.InterfaceType:
		return Type{
			Prefix: "",
			Name:   "interface{}",
		}, nil
	}
	// TODO: ... and map
	return Type{}, xerrors.Errorf("type %T is not supported to parse", t)
}

// TODO: control import modules
