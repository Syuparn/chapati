package presenter

import (
	"io"

	"github.com/syuparn/chapati/usecase"
)

type curryFunctionPresenter struct {
	writer io.Writer
}

func (p *curryFunctionPresenter) Show(out usecase.CurryFunctionOutputData) {
	// TODO: implement
}

// NewCurryFunctionPresenter creates a new CurryFunctionPresenter.
func NewCurryFunctionPresenter(writer io.Writer) usecase.CurryFunctionOutputPort {
	return &curryFunctionPresenter{
		writer: writer,
	}
}
