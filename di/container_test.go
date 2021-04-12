package di

import (
	"os"
	"testing"

	"github.com/syuparn/chapati/interface/controller"
)

func TestDI(t *testing.T) {
	container := NewContainer(os.Stdout)

	err := container.Invoke(func(c controller.CurryFunctionController) {
		// noop
	})
	if err != nil {
		t.Errorf("failed to invoke controller: %v", err)
	}
}
