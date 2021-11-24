package fizzbuzz

import (
	"strconv"
	"strings"
)

func join(strs ...string) string {
	var sb strings.Builder

	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}

// FizzBuzz : return a fizzbuzz series accoridng to params
func FizzBuzz(x, y, limit int, fizz, buzz string) []string {
	res := make([]string, limit)

	for i := range res {
		z := i + 1
		switch {
		case z%(x*y) == 0:
			res[i] = join(fizz, buzz)
		case z%y == 0:
			res[i] = join(buzz)
		case z%x == 0:
			res[i] = join(fizz)
		default:
			res[i] = strconv.Itoa(z)
		}
	}
	return res
}
