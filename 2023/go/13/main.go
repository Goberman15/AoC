package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	input := getInput()
	inputSlice := strings.Split(input, "\n\n")

	total := 0

	for idx, pattern := range inputSlice {
		parsed := parsePattern(pattern)
		if idx == 99 {
			parsed = parsed[:len(parsed)-1]
		}
		num, isFound := findReflection(parsed)

		if isFound {
			total += 100 * num
			continue
		}

		transposed := transpose(parsed)
		num, isFound = findReflection(transposed)

		if isFound {
			total += num
			continue
		}
	}

	fmt.Println("Total Reflection", total)

	total2 := 0

	for idx, pattern := range inputSlice {
		parsed := parsePattern(pattern)
		if idx == 99 {
			parsed = parsed[:len(parsed)-1]
		}
		num, isFound := findReflection2(parsed)

		if isFound {
			total2 += 100 * num
			continue
		}

		transposed := transpose(parsed)
		num, isFound = findReflection2(transposed)

		if isFound {
			total2 += num
			continue
		}
	}

	fmt.Println("Total Reflection", total2)
}

func getInput() string {
	req, err := http.NewRequest("GET", "https://adventofcode.com/2023/day/13/input", nil)

	if err != nil {
		panic(err)
	}

	cookie := os.Getenv("AOC_COOKIE")

	req.Header.Set("cookie", cookie)

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		panic("Fail to Retrieve Input")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading response body: %v\n", err)
	}

	return string(body)
}

func findReflection(pattern [][]string) (int, bool) {
	for i := 0; i < len(pattern)-1; i++ {
		isFound := false
		if slices.Equal(pattern[i], pattern[i+1]) {
			x := i
			y := i + 1

			for {
				x--
				y++

				if x < 0 || y >= len(pattern) {
					isFound = true
					break
				}

				if !slices.Equal(pattern[x], pattern[y]) {
					break
				}
			}

			if isFound {
				return i + 1, isFound
			}
		}
	}
	return 0, false
}

func findReflection2(pattern [][]string) (int, bool) {
	for i := 0; i < len(pattern)-1; i++ {
		totalSmudge := 0
		isFound := false
		eq, smudge := equal(pattern[i], pattern[i+1])
		totalSmudge += smudge
		if eq {
			x := i
			y := i + 1

			for {
				x--
				y++

				if x < 0 || y >= len(pattern) {
					if totalSmudge == 1 {
						isFound = true
					}
					break
				}

				eq, smudge = equal(pattern[x], pattern[y])
				totalSmudge += smudge

				if !eq {
					break
				}

				if totalSmudge > 1 {
					break
				}
			}

			if isFound {
				return i + 1, isFound
			}
		}
	}
	return 0, false
}

func parsePattern(pattern string) [][]string {
	parsed := make([][]string, 0, len(pattern))

	splitPattern := strings.Split(pattern, "\n")

	for _, v := range splitPattern {
		splitted := strings.Split(v, "")
		parsed = append(parsed, splitted)
	}

	return parsed
}

func transpose(pattern [][]string) [][]string {
	transposed := make([][]string, 0, len(pattern[0]))

	for x := 0; x < len(pattern[0]); x++ {
		line := make([]string, 0, len(pattern))
		for y := 0; y < len(pattern); y++ {
			line = append(line, pattern[y][x])
		}
		transposed = append(transposed, line)
	}

	return transposed
}

func equal(p1, p2 []string) (bool, int) {
	smudge := 0
	for i, v := range p1 {
		if v != p2[i] {
			smudge++

			if smudge > 1 {
				return false, smudge
			}
		}
	}

	return true, smudge
}
