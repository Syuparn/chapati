package di

import (
	"io"

	"go.uber.org/dig"

	"github.com/syuparn/chapati/infrastructure"
	"github.com/syuparn/chapati/interface/controller"
	"github.com/syuparn/chapati/interface/presenter"
	"github.com/syuparn/chapati/usecase"
)

// NewContainer creates a new DI container.
func NewContainer(w io.Writer) *dig.Container {
	c := dig.New()

	// domain
	c.Provide(infrastructure.NewCurryService)

	// usecase
	c.Provide(usecase.NewCurryFunctionInputPort)
	c.Provide(presenter.NewCurryFunctionPresenter)

	// interface
	c.Provide(controller.NewCurryFunctionController)

	// writer
	c.Provide(func() io.Writer { return w })

	return c
}
