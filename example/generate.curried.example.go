package example

func CurriedAdd(i1 int) func(int) int {
	return func(i2 int) int {
		return Add(i1, i2)
	}
}
