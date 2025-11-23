package utils

func MapSlice[T any, R any](f func(T) R, s []T) []R {
	var res []R = make([]R, len(s))

	for i, val := range s {
		res[i] = f(val)
	}

	return res
}

func FilterSlice[T any](f func(T) bool, s []T) []T {
	var res []T

	for _, val := range s {
		if f(val) {
			res = append(res, val)
		}
	}

	return res
}
