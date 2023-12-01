package main

import (
	"fmt"
	"strings"
)

func main() {
	d, found := checkSpelledDigit("seven4xk")

	if found {
		fmt.Println(d)
	}
}

func checkSpelledDigit(line string) (int, bool) {
	numMap := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
		"four":  4,
		"five":  5,
		"six":   6,
		"seven": 7,
		"eight": 8,
		"nine":  9,
	}

	for k, v := range numMap {
		if strings.Contains(line, k) {
			if line[:len(k)] == k {
				return v, true
			}
		}
	}

	return 0, false
}
