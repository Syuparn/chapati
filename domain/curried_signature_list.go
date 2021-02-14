package domain

// CurriedSignatureList represents a list of curried function signatures.
type CurriedSignatureList struct {
	CurriedSignature           *FunctionSignature
	PartiallyAppliedSignatures []*FunctionSignature
}
