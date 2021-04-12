package presenter

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/lithammer/dedent"

	"github.com/syuparn/chapati/domain"
	"github.com/syuparn/chapati/usecase"
)

func TestNewCurryFunctionPresenter(t *testing.T) {
	tests := []struct {
		name        string
		writer      io.Writer
		packageName string
		expected    usecase.CurryFunctionOutputPort
	}{
		{
			"new presenter",
			os.Stdout,
			"myPackage",
			&curryFunctionPresenter{
				writer:      os.Stdout,
				packageName: "myPackage",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := NewCurryFunctionPresenter(tt.writer, tt.packageName)

			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("wrong value: expected %#v, got %#v", tt.expected, actual)
			}
		})
	}
}

func TestCurryFunctionPresenterShow(t *testing.T) {
	tests := []struct {
		name        string
		packageName string
		origSig     *domain.FunctionSignature
		currySig    *domain.CurriedSignatureList
		expected    string
	}{
		{
			"arity 2 currying",
			"mypackage",
			domain.NewFunctionSignature(
				"myFunc",
				[]domain.Parameter{
					domain.NewParameter("arg0", domain.TermType("string")),
					domain.NewParameter("arg1", domain.TermType("int")),
				},
				[]domain.Type{
					domain.TermType("error"),
				},
			),
			domain.NewCurriedSignatureList(
				domain.NewFunctionSignature(
					"curriedMyFunc",
					[]domain.Parameter{
						domain.NewParameter("arg0", domain.TermType("string")),
					},
					[]domain.Type{
						domain.NewFuncType(
							[]domain.Type{
								domain.TermType("int"),
							},
							[]domain.Type{
								domain.TermType("error"),
							},
						),
					},
				),
				[]*domain.FunctionSignature{
					domain.NewFunctionSignature(
						"myFunc1",
						[]domain.Parameter{
							domain.NewParameter("arg1", domain.TermType("int")),
						},
						[]domain.Type{
							domain.TermType("error"),
						},
					),
				},
			),
			`
			package mypackage

			func curriedMyFunc(arg0 string) func(int) error {
				return func(arg1 int) error {
					return myFunc(arg0, arg1)
				}
			}
			`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			p := NewCurryFunctionPresenter(&buf, tt.packageName)

			err := p.Show(&usecase.CurryFunctionOutputData{
				OriginalSignatureList: tt.origSig,
				CurriedSignatureList:  tt.currySig,
			})

			if err != nil {
				t.Fatalf("error must be nil: %v", err)
			}

			actual := buf.String()
			expected := strings.TrimPrefix(dedent.Dedent(tt.expected), "\n")

			if actual != expected {
				t.Errorf("wrong value: expected ```\n%s\n```, got ```\n%s\n```", expected, actual)
			}
		})
	}
}

func TestCurryFunctionPresenterShowFailed(t *testing.T) {
	tests := []struct {
		name        string
		packageName string
		out         *usecase.CurryFunctionOutputData
	}{
		{
			"curriedSignatureList is not curry func",
			"mypackage",
			&usecase.CurryFunctionOutputData{
				OriginalSignatureList: domain.NewFunctionSignature(
					"myFunc",
					[]domain.Parameter{
						domain.NewParameter("arg0", domain.TermType("string")),
					},
					[]domain.Type{
						domain.TermType("error"),
					},
				),
				CurriedSignatureList: domain.NewCurriedSignatureList(
					domain.NewFunctionSignature(
						"nonCurriedMyFunc",
						[]domain.Parameter{
							domain.NewParameter("arg0", domain.TermType("string")),
						},
						[]domain.Type{
							domain.TermType("error"),
						},
					),
					[]*domain.FunctionSignature{},
				),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			p := NewCurryFunctionPresenter(&buf, tt.packageName)

			err := p.Show(tt.out)

			if err == nil {
				t.Fatalf("error must not be nil")
			}
		})
	}
}

func TestCurryFunctionPresenterCurryCode(t *testing.T) {
	tests := []struct {
		name     string
		origSig  *domain.FunctionSignature
		currySig *domain.CurriedSignatureList
		expected string
	}{
		{
			"arity 2 currying",
			domain.NewFunctionSignature(
				"myFunc",
				[]domain.Parameter{
					domain.NewParameter("arg0", domain.TermType("string")),
					domain.NewParameter("arg1", domain.TermType("int")),
				},
				[]domain.Type{
					domain.TermType("error"),
				},
			),
			domain.NewCurriedSignatureList(
				domain.NewFunctionSignature(
					"curriedMyFunc",
					[]domain.Parameter{
						domain.NewParameter("arg0", domain.TermType("string")),
					},
					[]domain.Type{
						domain.NewFuncType(
							[]domain.Type{
								domain.TermType("int"),
							},
							[]domain.Type{
								domain.TermType("error"),
							},
						),
					},
				),
				[]*domain.FunctionSignature{
					domain.NewFunctionSignature(
						"myFunc1",
						[]domain.Parameter{
							domain.NewParameter("arg1", domain.TermType("int")),
						},
						[]domain.Type{
							domain.TermType("error"),
						},
					),
				},
			),
			`
			func curriedMyFunc(arg0 string) func(int) error {
				return func(arg1 int) error {
					return myFunc(arg0, arg1)
				}
			}`,
		},
		{
			"arity 3 currying",
			domain.NewFunctionSignature(
				"myFunc",
				[]domain.Parameter{
					domain.NewParameter("arg0", domain.TermType("Arg0")),
					domain.NewParameter("arg1", domain.TermType("Arg1")),
					domain.NewParameter("arg2", domain.TermType("Arg2")),
				},
				[]domain.Type{
					domain.TermType("Ret0"),
				},
			),
			domain.NewCurriedSignatureList(
				domain.NewFunctionSignature(
					"curriedMyFunc",
					[]domain.Parameter{
						domain.NewParameter("arg0", domain.TermType("Arg0")),
					},
					[]domain.Type{
						domain.NewFuncType(
							[]domain.Type{
								domain.TermType("Arg1"),
							},
							[]domain.Type{
								domain.NewFuncType(
									[]domain.Type{
										domain.TermType("Arg2"),
									},
									[]domain.Type{
										domain.TermType("Ret0"),
									},
								),
							},
						),
					},
				),
				[]*domain.FunctionSignature{
					domain.NewFunctionSignature(
						"myFunc1",
						[]domain.Parameter{
							domain.NewParameter("arg1", domain.TermType("Arg1")),
						},
						[]domain.Type{
							domain.NewFuncType(
								[]domain.Type{
									domain.TermType("Arg2"),
								},
								[]domain.Type{
									domain.TermType("Ret0"),
								},
							),
						},
					),
					domain.NewFunctionSignature(
						"myFunc2",
						[]domain.Parameter{
							domain.NewParameter("arg2", domain.TermType("Arg2")),
						},
						[]domain.Type{
							domain.TermType("Ret0"),
						},
					),
				},
			),
			`
			func curriedMyFunc(arg0 Arg0) func(Arg1) func(Arg2) Ret0 {
				return func(arg1 Arg1) func(Arg2) Ret0 {
					return func(arg2 Arg2) Ret0 {
						return myFunc(arg0, arg1, arg2)
					}
				}
			}`,
		},
		{
			"first arg is func",
			domain.NewFunctionSignature(
				"myFunc",
				[]domain.Parameter{
					domain.NewParameter(
						"arg0",
						domain.NewFuncType(
							[]domain.Type{domain.TermType("int")},
							[]domain.Type{domain.TermType("string")},
						),
					),
					domain.NewParameter("arg1", domain.TermType("bool")),
				},
				[]domain.Type{
					domain.TermType("error"),
				},
			),
			domain.NewCurriedSignatureList(
				domain.NewFunctionSignature(
					"curriedMyFunc",
					[]domain.Parameter{
						domain.NewParameter(
							"arg0",
							domain.NewFuncType(
								[]domain.Type{domain.TermType("int")},
								[]domain.Type{domain.TermType("string")},
							),
						),
					},
					[]domain.Type{
						domain.NewFuncType(
							[]domain.Type{
								domain.TermType("bool"),
							},
							[]domain.Type{
								domain.TermType("error"),
							},
						),
					},
				),
				[]*domain.FunctionSignature{
					domain.NewFunctionSignature(
						"myFunc1",
						[]domain.Parameter{
							domain.NewParameter("arg1", domain.TermType("bool")),
						},
						[]domain.Type{
							domain.TermType("error"),
						},
					),
				},
			),
			`
			func curriedMyFunc(arg0 func(int) string) func(bool) error {
				return func(arg1 bool) error {
					return myFunc(arg0, arg1)
				}
			}`,
		},
		{
			"second arg is func",
			domain.NewFunctionSignature(
				"myFunc",
				[]domain.Parameter{
					domain.NewParameter("arg0", domain.TermType("string")),
					domain.NewParameter(
						"arg1",
						domain.NewFuncType(
							[]domain.Type{domain.TermType("int")},
							[]domain.Type{domain.TermType("bool")},
						),
					),
				},
				[]domain.Type{
					domain.TermType("error"),
				},
			),
			domain.NewCurriedSignatureList(
				domain.NewFunctionSignature(
					"curriedMyFunc",
					[]domain.Parameter{
						domain.NewParameter("arg0", domain.TermType("string")),
					},
					[]domain.Type{
						domain.NewFuncType(
							[]domain.Type{
								domain.NewFuncType(
									[]domain.Type{domain.TermType("int")},
									[]domain.Type{domain.TermType("bool")},
								),
							},
							[]domain.Type{
								domain.TermType("error"),
							},
						),
					},
				),
				[]*domain.FunctionSignature{
					domain.NewFunctionSignature(
						"myFunc1",
						[]domain.Parameter{
							domain.NewParameter(
								"arg1",
								domain.NewFuncType(
									[]domain.Type{domain.TermType("int")},
									[]domain.Type{domain.TermType("bool")},
								),
							),
						},
						[]domain.Type{
							domain.TermType("error"),
						},
					),
				},
			),
			`
			func curriedMyFunc(arg0 string) func(func(int) bool) error {
				return func(arg1 func(int) bool) error {
					return myFunc(arg0, arg1)
				}
			}`,
		},
		{
			"return type is func",
			domain.NewFunctionSignature(
				"myFunc",
				[]domain.Parameter{
					domain.NewParameter("arg0", domain.TermType("string")),
					domain.NewParameter("arg1", domain.TermType("int")),
				},
				[]domain.Type{
					domain.NewFuncType(
						[]domain.Type{domain.TermType("bool")},
						[]domain.Type{domain.TermType("error")},
					),
				},
			),
			domain.NewCurriedSignatureList(
				domain.NewFunctionSignature(
					"curriedMyFunc",
					[]domain.Parameter{
						domain.NewParameter("arg0", domain.TermType("string")),
					},
					[]domain.Type{
						domain.NewFuncType(
							[]domain.Type{
								domain.TermType("int"),
							},
							[]domain.Type{
								domain.NewFuncType(
									[]domain.Type{domain.TermType("bool")},
									[]domain.Type{domain.TermType("error")},
								),
							},
						),
					},
				),
				[]*domain.FunctionSignature{
					domain.NewFunctionSignature(
						"myFunc1",
						[]domain.Parameter{
							domain.NewParameter("arg1", domain.TermType("int")),
						},
						[]domain.Type{
							domain.NewFuncType(
								[]domain.Type{domain.TermType("bool")},
								[]domain.Type{domain.TermType("error")},
							),
						},
					),
				},
			),
			`
			func curriedMyFunc(arg0 string) func(int) func(bool) error {
				return func(arg1 int) func(bool) error {
					return myFunc(arg0, arg1)
				}
			}`,
		},
		{
			"return 2 values",
			domain.NewFunctionSignature(
				"myFunc",
				[]domain.Parameter{
					domain.NewParameter("arg0", domain.TermType("string")),
					domain.NewParameter("arg1", domain.TermType("int")),
				},
				[]domain.Type{
					domain.TermType("bool"),
					domain.TermType("error"),
				},
			),
			domain.NewCurriedSignatureList(
				domain.NewFunctionSignature(
					"curriedMyFunc",
					[]domain.Parameter{
						domain.NewParameter("arg0", domain.TermType("string")),
					},
					[]domain.Type{
						domain.NewFuncType(
							[]domain.Type{
								domain.TermType("int"),
							},
							[]domain.Type{
								domain.TermType("bool"),
								domain.TermType("error"),
							},
						),
					},
				),
				[]*domain.FunctionSignature{
					domain.NewFunctionSignature(
						"myFunc1",
						[]domain.Parameter{
							domain.NewParameter("arg1", domain.TermType("int")),
						},
						[]domain.Type{
							domain.TermType("bool"),
							domain.TermType("error"),
						},
					),
				},
			),
			`
			func curriedMyFunc(arg0 string) func(int) (bool, error) {
				return func(arg1 int) (bool, error) {
					return myFunc(arg0, arg1)
				}
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &curryFunctionPresenter{
				writer:      os.Stdout,
				packageName: "myPackage",
			}
			code, err := p.curryCode(tt.currySig, tt.origSig)

			if err != nil {
				t.Fatalf("error must be nil: got=%v", err)
			}

			actual := fmt.Sprintf("%#v", code)
			expected := strings.TrimPrefix(dedent.Dedent(tt.expected), "\n")
			if actual != expected {
				t.Errorf("wrong value: expected ```\n%s\n```, got ```\n%s\n```", expected, actual)
			}
		})
	}
}
