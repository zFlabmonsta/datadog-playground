package main

import "errors"

func divide(a, b float32) (float32, error) {
	if b == 0 {
		return 0, errors.New("divide(): cannot divide by zero")
	}

	return a / b, nil
}
