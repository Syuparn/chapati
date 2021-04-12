package usecase

import (
	"github.com/syuparn/chapati/domain"
)

// CurryFunctionInputPort executes currying function.
type CurryFunctionInputPort interface {
	Exec(in *CurryFunctionInputData) error
}

// CurryFunctionInputData is a DTO for CurryFunctionInputPort.
type CurryFunctionInputData struct {
	FuncName        string
	CurriedFuncName string
	Parameters      map[string]string
	ReturnTypes     []string
}

// CurryFunctionOutputPort presents the result of currying function.
type CurryFunctionOutputPort interface {
	Show(out *CurryFunctionOutputData) error
}

// CurryFunctionOutputData is a DTO for CurryFunctionOutputPort.
type CurryFunctionOutputData struct {
	OriginalSignatureList *domain.FunctionSignature
	CurriedSignatureList  *domain.CurriedSignatureList
}
