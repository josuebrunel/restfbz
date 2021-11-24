package fizzbuzz_test

import (
	fbz "restfbz/pkg/fizzbuzz"
	"testing"
)

func count(str string, strs []string) int {
	ct := 0
	for i := range strs {
		if strs[i] == str {
			ct++
		}
	}
	return ct
}

func TestFizzBuzz(t *testing.T) {
	res := fbz.FizzBuzz(3, 5, 25, "fizz", "buzz")
	if count("fizz", res) != 7 {
		t.Fatal("Invalid number of fizz")
	}
	if count("buzz", res) != 4 {
		t.Fatal("Invalid number of buzz")
	}
	if count("fizzbuzz", res) != 1 {
		t.Fatal("Invalid number of FizzBuzz")
	}
}
