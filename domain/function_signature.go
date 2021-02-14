package domain

// FunctionSignature represents a signature format of a function.
type FunctionSignature struct {
	name        string
	params      []Parameter
	returnTypes []Type
}

// Name returns the name of the function.
func (s *FunctionSignature) Name() string { return s.name }

// Parameters returns the parameters of the signature.
func (s *FunctionSignature) Parameters() []Parameter {
	copied := make([]Parameter, len(s.params))
	copy(copied, s.params)
	return copied
}

// ReturnTypes returns the return types of the signature.
func (s *FunctionSignature) ReturnTypes() []Type {
	copied := make([]Type, len(s.returnTypes))
	copy(copied, s.returnTypes)
	return copied
}

// Arity returns the number of parameters of the signature.
func (s *FunctionSignature) Arity() int { return len(s.params) }

// Type returns the type of the signature.
func (s *FunctionSignature) Type() Type {
	paramTypes := make([]Type, len(s.params))
	for i, param := range s.params {
		paramTypes[i] = param.Type
	}

	return FuncType{
		paramTypes:  paramTypes,
		returnTypes: s.returnTypes,
	}
}

// NewFunctionSignature creates a new FunctionSignature.
func NewFunctionSignature(
	name string,
	params []Parameter,
	returnTypes []Type,
) *FunctionSignature {
	return &FunctionSignature{
		name:        name,
		params:      params,
		returnTypes: returnTypes,
	}
}
