package infrastructure

import (
	"fmt"

	"github.com/syuparn/chapati/domain"
)

// NewCurryService generates a new CurryService.
func NewCurryService() domain.CurryService {
	return &curryService{}
}

type curryService struct{}

// Curry generates CurryiedSignatureList of FunctionSignature.
func (s *curryService) Curry(
	fn *domain.FunctionSignature,
	name string,
) (*domain.CurriedSignatureList, error) {
	if fn == nil {
		return nil, fmt.Errorf("fn must not be nil")
	}

	if fn.Arity() <= 1 {
		return nil, fmt.Errorf("no need to curry fn (arity=%d)", fn.Arity())
	}

	// [0]: func(x ArgN-1) (Ret0,...,RetM-1)
	// [1]: func(x ArgN-2) (func(ArgN-1) (Ret1,...,RetM-1)),
	// [2]: func(x ArgN-3) (func(ArgN-2) (func(ArgN-1) (Ret1,...,RetM-1))),
	// ...
	reversedPartiallyAppliedSignatures := make([]*domain.FunctionSignature, fn.Arity()-1)

	reversedPartiallyAppliedSignatures[0] = domain.NewFunctionSignature(
		fmt.Sprintf("%s%d", fn.Name(), fn.Arity()-1), // dummy name
		[]domain.Parameter{fn.Parameters()[fn.Arity()-1]},
		fn.ReturnTypes(),
	)

	for i := 1; i < fn.Arity()-1; i++ {
		reversedPartiallyAppliedSignatures[i] = domain.NewFunctionSignature(
			fmt.Sprintf("%s%d", fn.Name(), fn.Arity()-i-1), // dummy name
			[]domain.Parameter{fn.Parameters()[fn.Arity()-i-1]},
			[]domain.Type{reversedPartiallyAppliedSignatures[i-1].Type()},
		)
	}

	partiallyAppliedSignatures := make([]*domain.FunctionSignature, fn.Arity()-1)
	for i := 0; i < fn.Arity()-1; i++ {
		partiallyAppliedSignatures[i] = reversedPartiallyAppliedSignatures[fn.Arity()-2-i]
	}

	curriedSignature := domain.NewFunctionSignature(
		name,
		[]domain.Parameter{fn.Parameters()[0]},
		[]domain.Type{partiallyAppliedSignatures[0].Type()},
	)

	return domain.NewCurriedSignatureList(
		curriedSignature, partiallyAppliedSignatures), nil
}
