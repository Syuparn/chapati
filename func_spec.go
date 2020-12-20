package main

import "strings"

// FuncSpec is a specification of function name and signature.
type FuncSpec struct {
	name        string
	recv        *Param // may be nil
	params      []Param
	returnTypes []Type
}

// NewFuncSpec creates a new FuncSpec.
func NewFuncSpec(
	name string,
	recv *Param,
	params []Param,
	returnTypes []Type,
) *FuncSpec {
	return &FuncSpec{
		name:        name,
		recv:        recv,
		params:      params,
		returnTypes: returnTypes,
	}
}

// IsExported returns true if the function is exported.
func (s *FuncSpec) IsExported() bool {
	firstChar := s.name[:1]
	return strings.ToUpper(firstChar) == firstChar
}

// Name returns function name.
func (s *FuncSpec) Name() string {
	return s.name
}

// Params returns function parameters.
func (s *FuncSpec) Params() []Param {
	return s.params
}

// Recv returns function receiver.
func (s *FuncSpec) Recv() (Param, bool) {
	if s.recv == nil {
		return Param{}, false
	}
	return *s.recv, true
}

// ReturnTypes returns function return types.
func (s *FuncSpec) ReturnTypes() []Type {
	return s.returnTypes
}

// Param is a struct for parameter declaration.
type Param struct {
	Name string
	Type Type
}

// NewParam creates a new Param.
func NewParam(name string, t Type) Param {
	return Param{Name: name, Type: t}
}

// Type is a struct of type declaration.
type Type struct {
	Prefix string // combination of `*`, `[]`, `...` etc.
	Name   string
}

// NewType creates a new Type.
func NewType(name string, prefix string) Type {
	return Type{Prefix: prefix, Name: name}
}
