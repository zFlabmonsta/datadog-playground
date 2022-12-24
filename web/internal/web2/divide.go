package web2

import "errors"

func Divide(a, b float32) (float32, error) {
	if b == 0 {
		return 0, errors.New("divide(): cannot divide by zero")
	}

	return a / b, nil
}
