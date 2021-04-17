package test

import "io"

func write(w io.Writer) {
	w.Write([]byte("hello"))
}
