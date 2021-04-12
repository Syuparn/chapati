package presenter

import (
	"io"

	"golang.org/x/xerrors"

	"github.com/dave/jennifer/jen"

	"github.com/syuparn/chapati/domain"
	"github.com/syuparn/chapati/usecase"
)

type curryFunctionPresenter struct {
	writer      io.Writer
	packageName string
}

// NewCurryFunctionPresenter creates a new CurryFunctionPresenter.
func NewCurryFunctionPresenter(
	writer io.Writer,
	packageName string,
) usecase.CurryFunctionOutputPort {
	return &curryFunctionPresenter{
		writer:      writer,
		packageName: packageName,
	}
}

// Show writes source code of curried function to p.writer.
func (p *curryFunctionPresenter) Show(out usecase.CurryFunctionOutputData) error {
	if out.Error != nil {
		return out.Error
	}

	f := jen.NewFilePath(p.packageName)

	curryCode, err := p.curryCode(out.CurriedSignatureList, out.OriginalSignatureList)
	if err != nil {
		return xerrors.Errorf("failed to generate code: %w", err)
	}
	f.Add(curryCode)

	if err := f.Render(p.writer); err != nil {
		return xerrors.Errorf("failed to write code: %w", err)
	}

	return nil
}

func (p *curryFunctionPresenter) curryCode(
	currySig *domain.CurriedSignatureList,
	origSig *domain.FunctionSignature,
) (jen.Code, error) {
	if len(currySig.PartiallyAppliedSignatures) == 0 {
		return nil, xerrors.Errorf("PartiallyAppliedSignatures must not be zero")
	}

	reversedSigs := make([]*domain.FunctionSignature, len(currySig.PartiallyAppliedSignatures))
	for i, sig := range currySig.PartiallyAppliedSignatures {
		reversedSigs[len(currySig.PartiallyAppliedSignatures)-i-1] = sig
	}

	// inner most function
	code := p.curryCoreCode(reversedSigs[0], origSig)

	// inner functions from inner to outer
	for _, sig := range reversedSigs[1:] {
		code = p.curryMiddleCode(sig, code)
	}

	// outer function
	code = p.curryOuterCode(currySig.CurriedSignature, code)

	return code, nil
}

func (p *curryFunctionPresenter) curryOuterCode(
	sig *domain.FunctionSignature,
	inner jen.Code,
) jen.Code {
	fn := jen.Func()

	// function name
	fn.Id(sig.Name())

	// function params
	fn.Params(renderParams(sig.Parameters())...)

	// function return types
	if len(sig.ReturnTypes()) > 0 {
		fn.Params(renderTypes(sig.ReturnTypes())...)
	}

	fn.Block(
		jen.Return(inner),
	)

	return fn
}

func (p *curryFunctionPresenter) curryMiddleCode(
	sig *domain.FunctionSignature,
	inner jen.Code,
) jen.Code {
	fn := jen.Func()

	// function params
	fn.Params(renderParams(sig.Parameters())...)

	// function return types
	if len(sig.ReturnTypes()) > 0 {
		fn.Params(renderTypes(sig.ReturnTypes())...)
	}

	fn.Block(
		jen.Return(inner),
	)

	return fn
}

func (p *curryFunctionPresenter) curryCoreCode(
	sig *domain.FunctionSignature,
	origSig *domain.FunctionSignature,
) jen.Code {
	fn := jen.Func()

	// function params
	fn.Params(renderParams(sig.Parameters())...)

	// function return types
	if len(sig.ReturnTypes()) > 0 {
		fn.Params(renderTypes(sig.ReturnTypes())...)
	}

	fn.Block(
		jen.Return(
			jen.Id(origSig.Name()).
				Call(renderParamValues(origSig.Parameters())...),
		),
	)

	return fn
}
