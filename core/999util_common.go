package core

func ifElse[T any](flag bool, val1 T, val2 T) T {
	if flag {
		return val1
	} else {
		return val2
	}
}