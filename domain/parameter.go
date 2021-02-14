package domain

// Parameter represents a parameter in a function signature.
type Parameter struct {
	Name string
	Type Type
}

// NewParameter creates a new Parameter.
func NewParameter(name string, t Type) Parameter {
	return Parameter{Name: name, Type: t}
}
