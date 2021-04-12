package controller

import (
	"golang.org/x/xerrors"

	"github.com/syuparn/chapati/usecase"
)

type CurryFunctionController interface {
	Handle(src string) error
}

type curryFunctionController struct {
	inputPort usecase.CurryFunctionInputPort
	extracter
}

// NewCurryFunctionController creates a new CurryFunctionController.
func NewCurryFunctionController(
	inputPort usecase.CurryFunctionInputPort,
) CurryFunctionController {
	return &curryFunctionController{
		inputPort: inputPort,
	}
}

// Handle generates curried function from source code.
func (c *curryFunctionController) Handle(fileName string) error {
	in, err := c.extractFuncInfo(fileName)
	if err != nil {
		return xerrors.Errorf("failed to extract function from source code: %w", err)
	}

	if err := c.inputPort.Exec(in); err != nil {
		return err
	}

	return nil
}
