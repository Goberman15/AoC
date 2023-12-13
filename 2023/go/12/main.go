package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"slices"
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
	// input := "???.### 1,1,3\n.??..??...?##. 1,1,3\n?#?#?#?#?#?#?#? 1,3,1,6\n????.#...#... 4,1,1\n????.######..#####. 1,6,5\n?###???????? 3,2,1"
	inputSlice := strings.Split(input, "\n")
	inputSlice = inputSlice[:len(inputSlice)-1]

	total := 0

	for idx, v := range inputSlice {
		springs, sizes := parseSpringLine(v)

		arrangements := countArrangements(springs, sizes)
		fmt.Println(idx, arrangements)
		total += arrangements
		fmt.Println()
	}

	fmt.Println("Total Arrangements", total)

}

func getInput() string {
	req, err := http.NewRequest("GET", "https://adventofcode.com/2023/day/12/input", nil)

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

func parseSpringLine(line string) ([]string, []int) {
	splitted := strings.Split(line, " ")
	springs := strings.Split(splitted[0], "")
	sizes := strings.Split(splitted[1], ",")

	return springs, toIntSlice(sizes)
}

func countArrangements(springs []string, sizes []int) int {
	fmt.Println(springs)
	fmt.Println(sizes)
	count := 0
	arrangements := 0
	blankIdx := make([]int, 0)
	for idx, s := range springs {
		if s == "?" {
			count++
			blankIdx = append(blankIdx, idx)
		}
	}

	cache := make(map[string]int)
	patterns := make(map[string]int)

	for x := 0; x < int(math.Pow(2, float64(count))); x++ {
		arr := make([]string, 0, count)
		checker := ""
		isNotWorked := false
		for y := 0; y < count; y++ {
			if (x & (1 << y)) > 0 {
				arr = append(arr, ".")
				checker += "."
			} else {
				arr = append(arr, "#")
				checker += "#"
			}

			if cache[checker] == -1 {
				break
			}
		}

		patterns[checker]++

		if isNotWorked {
			continue
		}

		if len(blankIdx) != count && len(arr) != count {
			log.Fatal("Something went Wrong")
		}

		row := make([]string, len(springs))
		copy(row, springs)

		for idx, v := range arr {
			row[blankIdx[idx]] = v
		}

		isValid := validator(row, sizes, cache, blankIdx)

		if isValid {
			arrangements++
		}
	}

	fmt.Println(math.Pow(2, float64(count)), len(patterns))

	return arrangements
}

func validator(springs []string, sizes []int, cache map[string]int, blankIdx []int) (isValid bool) {
	group := 0
	cp := make([]int, len(sizes))
	strChecker := ""
	prevBroken := false
	copy(cp, sizes)

	defer func() {
		if !isValid {
			cache[strChecker] = -1
		}
	}()

	for idx, v := range springs {
		if slices.Contains(blankIdx, idx) {
			strChecker += v
		}
		if v == "." {
			if prevBroken {
				if cp[group] > 0 {
					isValid = false
					return
				} else {
					prevBroken = false
					group++
				}
			}
		} else if v == "#" {
			if group >= len(sizes) {
				isValid = false
				return
			}

			cp[group]--

			if cp[group] < 0 {
				isValid = false
				return
			}

			prevBroken = true
		}
	}

	if group < len(sizes)-1 {
		isValid = false
		return
	}

	if group < len(sizes) && cp[group] != 0 {
		isValid = false
		return
	}

	// if !checkAllZero(cp) {
	// 	isValid = false
	// 	return
	// }

	isValid = true
	return
}

func toInt(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		panic("NaN")
	}

	return num
}

func toIntSlice(slc []string) []int {
	numSlc := make([]int, 0, len(slc))

	for _, v := range slc {
		numSlc = append(numSlc, toInt(v))
	}

	return numSlc
}
