package domain

// CurryService curries the function signature.
type CurryService interface {
	Curry(fn *FunctionSignature, name string) (*CurriedSignatureList, error)
}
