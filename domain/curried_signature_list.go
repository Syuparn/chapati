package domain

// CurriedSignatureList represents a list of curried function signatures.
type CurriedSignatureList struct {
	CurriedSignature           *FunctionSignature
	PartiallyAppliedSignatures []*FunctionSignature
}

// NewCurriedSignatureList returns a new CurriedSignatureList.
func NewCurriedSignatureList(
	curriedSignature *FunctionSignature,
	partiallyAppliedSignatures []*FunctionSignature,
) *CurriedSignatureList {
	return &CurriedSignatureList{
		CurriedSignature:           curriedSignature,
		PartiallyAppliedSignatures: partiallyAppliedSignatures,
	}
}
