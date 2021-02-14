package infrastructure

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/syuparn/chapati/domain"
)

func TestNewCurryService(t *testing.T) {
	actual := NewCurryService()
	if actual == nil {
		t.Fatalf("returned value must not be nil")
	}

	if _, ok := actual.(*curryService); !ok {
		t.Fatalf("expected type %T, got %T", &curryService{}, actual)
	}
}

func TestCurry(t *testing.T) {
	tests := []struct {
		fn       *domain.FunctionSignature
		name     string
		expected *domain.CurriedSignatureList
	}{
		{
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
			"curriedMyFunc",
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
		},
		{
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
			"curriedMyFunc",
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
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			svc := NewCurryService()
			actual, err := svc.Curry(tt.fn, tt.name)

			if err != nil {
				t.Fatalf("error must be nil. got=%s", err.Error())
			}

			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("wrong value: expected %#v, got %#v", tt.expected, actual)
			}
		})
	}
}

func TestCurryFailed(t *testing.T) {
	tests := []struct {
		fn       *domain.FunctionSignature
		name     string
		expected string
	}{
		{
			domain.NewFunctionSignature(
				"myFunc",
				[]domain.Parameter{},
				[]domain.Type{},
			),
			"curriedMyFunc",
			"no need to curry fn (arity=0)",
		},
		{
			domain.NewFunctionSignature(
				"myFunc",
				[]domain.Parameter{
					domain.NewParameter("arg0", domain.TermType("string")),
				},
				[]domain.Type{},
			),
			"curriedMyFunc",
			"no need to curry fn (arity=1)",
		},
		{
			nil,
			"curriedMyFunc",
			"fn must not be nil",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			svc := NewCurryService()
			_, err := svc.Curry(tt.fn, tt.name)

			if err == nil {
				t.Fatalf("error must not be nil")
			}

			if err.Error() != tt.expected {
				t.Errorf("got wrong message. expected `%s`, got `%s`",
					tt.expected, err.Error())
			}
		})
	}
}
