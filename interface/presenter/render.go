package presenter

import (
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/syuparn/chapati/domain"
)

func renderParam(p domain.Parameter) jen.Code {
	ident := jen.Id(p.Name)
	ident.Add(renderType(p.Type))

	return ident
}

func renderParams(params []domain.Parameter) []jen.Code {
	rendered := make([]jen.Code, len(params))
	for _, p := range params {
		rendered = append(rendered, renderParam(p))
	}

	return rendered
}

func renderParamValue(p domain.Parameter) jen.Code {
	return jen.Id(p.Name)
}

func renderParamValues(params []domain.Parameter) []jen.Code {
	rendered := make([]jen.Code, len(params))
	for _, p := range params {
		rendered = append(rendered, renderParamValue(p))
	}

	return rendered
}

func renderTypes(types []domain.Type) []jen.Code {
	rendered := make([]jen.Code, len(types))
	for _, t := range types {
		rendered = append(rendered, renderType(t))
	}

	return rendered
}

func renderType(t domain.Type) jen.Code {
	if t.IsFuncType() {
		return renderFuncType(t)
	}
	return renderTermType(t)
}

func renderTermType(t domain.Type) jen.Code {
	modulePath, typeName := splitType(t.(domain.TermType))
	if modulePath == "" {
		return jen.Id(typeName)
	}

	return jen.Qual(modulePath, typeName)
}

func renderFuncType(t domain.Type) jen.Code {
	fn := jen.Func()

	ft := t.(domain.FuncType)

	// function params
	fn.Params(renderTypes(ft.ParamTypes())...)

	// function return types
	if len(ft.ReturnTypes()) > 0 {
		fn.Params(renderTypes(ft.ReturnTypes())...)
	}

	return fn
}

func splitType(t domain.TermType) (modulePath, typeName string) {
	whole := string(t)
	iSep := strings.LastIndex(whole, ".")

	if iSep == -1 {
		modulePath, typeName = "", whole
		return
	}

	modulePath, typeName = whole[:iSep], whole[iSep+1:]
	return
}
