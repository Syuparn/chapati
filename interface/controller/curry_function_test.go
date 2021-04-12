package controller

import (
	"reflect"
	"testing"

	"github.com/syuparn/chapati/usecase"
)

func TestNewCurryFunctionController(t *testing.T) {
	port := newMockCurryFunctionInputPort()

	tests := []struct {
		name      string
		inputPort usecase.CurryFunctionInputPort
		expected  CurryFunctionController
	}{
		{
			"new controller",
			port,
			&curryFunctionController{
				inputPort: port,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := NewCurryFunctionController(tt.inputPort)

			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("wrong value: expected %#v, got %#v", tt.expected, actual)
			}
		})
	}
}

func TestCurryFunctionControllerHandle(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		expected *usecase.CurryFunctionInputData
	}{
		{
			"with one function",
			"simple.go",
			&usecase.CurryFunctionInputData{
				FuncName:        "PrintRepeat",
				CurriedFuncName: "CurriedPrintRepeat",
				Parameters: map[string]string{
					"msg": "string",
					"n":   "int",
				},
				ReturnTypes: []string{
					"error",
				},
				CurriedFunctionMetaData: usecase.CurriedFunctionMetaData{
					PackageName: "test",
				},
			},
		},
		{
			"multiple return values",
			"multi_return.go",
			&usecase.CurryFunctionInputData{
				FuncName:        "multi",
				CurriedFuncName: "CurriedMulti",
				Parameters:      map[string]string{},
				ReturnTypes: []string{
					"int",
					"bool",
				},
				CurriedFunctionMetaData: usecase.CurriedFunctionMetaData{
					PackageName: "test",
				},
			},
		},
		{
			"compound types",
			"compound.go",
			&usecase.CurryFunctionInputData{
				FuncName:        "handleCompound",
				CurriedFuncName: "CurriedHandleCompound",
				Parameters: map[string]string{
					"ptrArg":  "*string",
					"mapArg":  "map[string]interface{}",
					"arrArg":  "[]int",
					"funcArg": "func(int) string",
				},
				ReturnTypes: []string{},
				CurriedFunctionMetaData: usecase.CurriedFunctionMetaData{
					PackageName: "test",
				},
			},
		},
		{
			"defined type",
			"defined.go",
			&usecase.CurryFunctionInputData{
				FuncName:        "hello",
				CurriedFuncName: "CurriedHello",
				Parameters: map[string]string{
					"person": "test.Person",
				},
				ReturnTypes: []string{},
				CurriedFunctionMetaData: usecase.CurriedFunctionMetaData{
					PackageName: "test",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			port := newMockCurryFunctionInputPort()
			c := NewCurryFunctionController(port)

			if err := c.Handle("testdata/" + tt.fileName); err != nil {
				t.Fatalf("error must be nil: %v", err)
			}

			if !reflect.DeepEqual(port.in, tt.expected) {
				t.Errorf("wrong value: expected \n%#v\n, got \n%#v\n",
					tt.expected, port.in)
			}
		})
	}
}

func TestCurryFunctionControllerHandleFailed(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
	}{
		{
			"file not found",
			"notfound.go",
		},
		{
			"no functions",
			"no_funcs.go",
		},
		{
			"parse error",
			"i_am_not_go.py",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			port := newMockCurryFunctionInputPort()
			c := NewCurryFunctionController(port)

			if err := c.Handle("testdata/" + tt.fileName); err == nil {
				t.Fatalf("error must not be nil")
			}
		})
	}
}

func newMockCurryFunctionInputPort() *mockCurryFunctionInputPort {
	return &mockCurryFunctionInputPort{}
}

type mockCurryFunctionInputPort struct {
	in *usecase.CurryFunctionInputData
}

func (p *mockCurryFunctionInputPort) Exec(in *usecase.CurryFunctionInputData) error {
	p.in = in
	return nil
}
