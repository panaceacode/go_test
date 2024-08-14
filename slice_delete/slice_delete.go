package slice_delete

import "errors"

func DeleteAt[T any](source []T, index int) ([]T, error) {

	length := len(source)
	if index < 0 || index >= length {
		return source, errors.New("index out of range")
	}

	if index >= length/2 {
		for i := index + 1; i < length; i++ {
			source[i-1] = source[i]
		}
		return source[:length-1], nil
	} else {
		for i := index; i > 0; i-- {
			source[i] = source[i-1]
		}
		return source[1:], nil
	}
}

func Shrink[T any](source []T) ([]T, error) {
	length, capacity := len(source), cap(source)
	len, changed := getCapacity(length, capacity)
	if changed {
		res := make([]T, 0, len)
		res = append(res, source...)
		return res, nil
	}
	return source, nil
}

func getCapacity(len, cap int) (int, bool) {
	if cap <= 64 {
		return cap, false
	}

	if cap >= 2048 && cap/len >= 2 {
		factor := 0.625
		return int(float32(cap) * float32(factor)), true
	}

	if cap < 2048 && cap/len >= 4 {
		return cap / 2, true
	}

	return cap, false
}
