package main

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func panicif(cond bool) {
	if cond {
		panic("condition failed")
	}
}
