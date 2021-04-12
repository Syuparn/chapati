package presenter

import (
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
	if !t.IsFuncType() {
		return jen.Id(string(t.(domain.TermType)))
	}

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
