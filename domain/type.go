package domain

// Type represents a type of each parameter or returned value in function signature.
type Type interface {
	// dummy method
	Type()
	// whether this is Function type or not
	IsFuncType() bool
}

// TermType represents a type other than Function type.
type TermType string

func (t TermType) String() string {
	return t.String()
}

// Type is a dummy method of Type interface.
func (t TermType) Type() {}

// IsFuncType returns whether this is Function type or not.
func (t TermType) IsFuncType() bool { return false }

// NewTermType creates a new TermType.
func NewTermType(s string) TermType { return TermType(s) }

// FuncType represents a type of a function.
type FuncType struct {
	paramTypes  []Type
	returnTypes []Type
}

// Type is a dummy method of Type interface.
func (t FuncType) Type() {}

// IsFuncType returns whether this is Function type or not.
func (t FuncType) IsFuncType() bool { return true }

// ParamTypes returns the parameter types of the signature.
func (t FuncType) ParamTypes() []Type { return t.paramTypes }

// ReturnTypes returns the return types of the signature.
func (t FuncType) ReturnTypes() []Type { return t.returnTypes }

// NewFuncType creates a new FuncType.
func NewFuncType(paramTypes []Type, returnTypes []Type) FuncType {
	return FuncType{
		paramTypes:  paramTypes,
		returnTypes: returnTypes,
	}
}
