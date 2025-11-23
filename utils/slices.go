package utils

func MapSlice[T any, R any](f func (T) R, s []T) []R {
	var res []R = make([]R, len(s))

	for i, val := range s {
		res[i] = f(val)
	}

	return res
}

