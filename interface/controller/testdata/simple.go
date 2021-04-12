package test

import "fmt"

// PrintRepeat prints msg n times.
func PrintRepeat(msg string, n int) error {
	if n < 0 {
		return fmt.Errorf("n must be positive")
	}

	for i := 0; i < n; i++ {
		fmt.Println(msg)
	}

	return nil
}
