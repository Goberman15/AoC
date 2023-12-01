package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	input := getInput()
	inputSlice := strings.Split(input, "\n")
	inputSlice = inputSlice[:len(inputSlice)-1]
	sum := 0

	for _, line := range inputSlice {
		coord := extractCalibrationValue2(line)
		// fmt.Println(line, coord)
		sum += coord
	}

	fmt.Printf("The Sum of All the Calibration Value is %d\n", sum)
}

func getInput() string {
	req, err := http.NewRequest("GET", "https://adventofcode.com/2023/day/1/input", nil)

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

func extractCalibrationValue(line string) int {
	headIdx, tailIdx := 0, len(line)-1
	headFound, tailFound := false, false
	head, tail := 0, 0

	for {
		if !headFound {
			if headIdx == tailIdx && tailFound {
				head = tail
				headFound = true
			} else {
				i, err := strconv.Atoi(string(line[headIdx]))

				if err != nil {
					headIdx++
				} else {
					headFound = true
					head = i
				}
			}
		}

		if !tailFound {
			if tailIdx == headIdx && headFound {
				tail = head
				tailFound = true
			} else {
				i, err := strconv.Atoi(string(line[tailIdx]))

				if err != nil {
					tailIdx--
				} else {
					tailFound = true
					tail = i
				}
			}
		}

		if headFound && tailFound {
			break
		}
	}

	return (10 * head) + tail
}

func extractCalibrationValue2(line string) int {
	foundSlice := make([]int, 0, len(line))

	ptr := 0

	for {
		if ptr >= len(line) {
			break
		}

		i, err := strconv.Atoi(string(line[ptr]))

		if err != nil {
			d, found := checkSpelledDigit(line[ptr:])

			if found {
				foundSlice = append(foundSlice, d)
			}
			ptr++
			continue
		}

		foundSlice = append(foundSlice, i)
		ptr++
	}

	return (10 * foundSlice[0]) + foundSlice[len(foundSlice)-1]
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
