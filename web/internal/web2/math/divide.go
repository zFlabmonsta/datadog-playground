package math

func Divide(a, b float32) (float32, error) {
	if b == 0 {
		return 0, ErrorDivisibleZero
	}

	return a / b, nil
}
