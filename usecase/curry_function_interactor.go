package usecase

import (
	"github.com/syuparn/chapati/domain"
	"golang.org/x/xerrors"
)

type curryFunctionInteractor struct {
	out          CurryFunctionOutputPort
	curryService domain.CurryService
}

func (p curryFunctionInteractor) Exec(in *CurryFunctionInputData) error {
	params := []domain.Parameter{}
	for n, t := range in.Parameters {
		param := domain.NewParameter(n, domain.TermType(t))
		params = append(params, param)
	}

	returnTypes := make([]domain.Type, len(in.ReturnTypes))
	for i, t := range in.ReturnTypes {
		returnTypes[i] = domain.TermType(t)
	}

	funcSignature := domain.NewFunctionSignature(in.FuncName, params, returnTypes)
	if funcSignature.Arity() <= 1 {
		return xerrors.Errorf("no need to curry fn (arity=%d)", funcSignature.Arity())
	}

	curried, err := p.curryService.Curry(funcSignature, in.CurriedFuncName)
	if err != nil {
		return xerrors.Errorf("failed to curry: %w", err)
	}

	out := &CurryFunctionOutputData{
		OriginalSignatureList:   funcSignature,
		CurriedSignatureList:    curried,
		CurriedFunctionMetaData: in.CurriedFunctionMetaData,
	}

	if err := p.out.Show(out); err != nil {
		return xerrors.Errorf("failed to present outputdata: %w", err)
	}

	return nil
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
