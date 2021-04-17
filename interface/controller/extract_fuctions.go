package controller

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"strings"

	"github.com/syuparn/chapati/usecase"
	"golang.org/x/xerrors"
)

// TODO: enable to set from config
const curriedFuncPrefix = "Curried"

type extracter struct{}

func (e extracter) extractFuncInfo(
	fileName string,
) (*usecase.CurryFunctionInputData, error) {
	fset := token.NewFileSet()
	conf := types.Config{
		// NOTE: use "source" to import directly from source
		// ("cg" cannot work if dependent package is not installed globally)
		Importer: importer.ForCompiler(fset, "source", nil),
	}

	//                                         src, mode
	f, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		return nil, xerrors.Errorf("failed to parse code: %w", err)
	}

	if f.Name == nil {
		return nil, xerrors.Errorf("failed to obtain package name")
	}
	packageName := f.Name.Name

	info := &types.Info{
		Defs: map[*ast.Ident]types.Object{},
	}

	_, err = conf.Check(packageName, fset, []*ast.File{f}, info)
	if err != nil {
		return nil, xerrors.Errorf("failed to extract info: %w", err)
	}

	for ident := range info.Defs {
		// TODO: return multiple functions
		funcType, ok := e.funcTypeOf(info, ident)
		if ok {
			return e.inputDataFrom(ident.Name, funcType, packageName), nil
		}
	}

	return nil, xerrors.Errorf("no functions found in soruce code")
}

func (e extracter) funcTypeOf(
	info *types.Info,
	ident *ast.Ident,
) (*types.Signature, bool) {
	f, ok := info.ObjectOf(ident).(*types.Func)
	if !ok {
		return nil, false
	}

	t, ok := f.Type().(*types.Signature)
	if !ok {
		return nil, false
	}

	// TODO: enable to use method
	if t.Recv() != nil {
		return nil, false
	}

	return t, true
}

func (e extracter) inputDataFrom(
	funcName string,
	t *types.Signature,
	packageName string,
) *usecase.CurryFunctionInputData {
	params := map[string]string{}
	returnTypes := make([]string, t.Results().Len())

	for i := 0; i < t.Params().Len(); i++ {
		p := t.Params().At(i)
		params[p.Name()] = p.Type().String()
	}

	for i := 0; i < t.Results().Len(); i++ {
		p := t.Results().At(i)
		returnTypes[i] = p.Type().String()
	}

	return &usecase.CurryFunctionInputData{
		FuncName:        funcName,
		CurriedFuncName: curriedFuncPrefix + strings.Title(funcName),
		Parameters:      params,
		ReturnTypes:     returnTypes,
		CurriedFunctionMetaData: usecase.CurriedFunctionMetaData{
			PackageName: packageName,
		},
	}
}
