package usecase

import "github.com/syuparn/chapati/domain"

type curryFunctionInteractor struct {
	out          CurryFunctionOutputPort
	curryService domain.CurryService
}

func (p curryFunctionInteractor) Exec(in CurryFunctionInputData) {
	params := []domain.Parameter{}
	for n, t := range in.Parameters {
		param := domain.NewParameter(n, domain.TermType(t))
		params = append(params, param)
	}

	returnTypes := make([]domain.Type, 0, len(in.ReturnTypes))
	for i, t := range in.ReturnTypes {
		returnTypes[i] = domain.TermType(t)
	}

	funcSignature := domain.NewFunctionSignature(in.FuncName, params, returnTypes)

	curried, err := p.curryService.Curry(funcSignature, in.CurriedFuncName)

	out := CurryFunctionOutputData{
		OriginalSignatureList: funcSignature,
		CurriedSignatureList:  curried,
		Error:                 err,
	}

	p.out.Show(out)
}

// NewCurryFunctionInputPort creates a new CurryFunctionInputPort.
func NewCurryFunctionInputPort(
	out CurryFunctionOutputPort,
	curryService domain.CurryService,
) CurryFunctionInputPort {
	return &curryFunctionInteractor{
		out:          out,
		curryService: curryService,
	}
}
