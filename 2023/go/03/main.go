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

var engineSchematic [][]string

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	input := getInput()
	inputSlice := strings.Split(input, "\n")
	inputSlice = inputSlice[:len(inputSlice)-1]

	for _, v := range inputSlice {
		splittedV := strings.Split(v, "")
		engineSchematic = append(engineSchematic, splittedV)
	}

	total := 0
	totalGearRatio := 0

	for i, line := range engineSchematic {
		out := getPartNumber(line, i)
		out2 := getGearRatio(line, i)
		total += out
		totalGearRatio += out2
	}

	fmt.Printf("Sum of All Part Number in the Engine Schematic is %d\n", total)
	fmt.Printf("Sum of All Gear Ratio in the Engine Schematic is %d\n", totalGearRatio)
}

func getInput() string {
	req, err := http.NewRequest("GET", "https://adventofcode.com/2023/day/3/input", nil)

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

func getPartNumber(line []string, i int) int {
	found := ""
	prevNum := false
	isValid := false
	validNum := 0

	for j, v := range line {
		if _, err := strconv.Atoi(v); err == nil {
			if !isValid {
				res := checkSurroundingSymbol(i, j)
				isValid = res
			}

			if !prevNum {
				prevNum = true
			}

			found += v
		} else {
			if prevNum {
				if isValid {
					num, err := strconv.Atoi(found)
					if err != nil {
						log.Fatal(err)
					}
					validNum += num
					isValid = false
				}
				prevNum = false
				found = ""
			}
		}
	}

	if prevNum && isValid {
		num, err := strconv.Atoi(found)
		if err != nil {
			log.Fatal(err)
		}
		validNum += num
	}

	return validNum
}

func getGearRatio(line []string, i int) int {
	total := 0
	for j, v := range line {
		if v == "*" {
			num := checkSurroundingNumber(i, j)

			if num > 0 {
				total += num
			}
		}
	}

	return total
}

func checkSurroundingSymbol(i, j int) bool {
	knownSymbol := []string{"*", "#", "=", "$", "%", "/", "@", "+", "-", "&"}
	surrounding := ""

	for x := i - 1; x <= i+1; x++ {
		if x < 0 || x >= len(engineSchematic) {
			continue
		}

		for y := j - 1; y <= j+1; y++ {
			if y < 0 || y >= len(engineSchematic[x]) || (x == i && y == j) {
				continue
			}

			surrounding += engineSchematic[x][y]
		}
	}

	for _, v := range knownSymbol {
		if strings.Contains(surrounding, v) {
			return true
		}
	}

	return false
}

func checkSurroundingNumber(i, j int) int {
	foundNum := make([]int, 2)
	idx := 0
	for x := i - 1; x <= i+1; x++ {
		if x < 0 || x >= len(engineSchematic) {
			continue
		}

		for y := j - 1; y <= j+1; y++ {
			if y < 0 || y >= len(engineSchematic[x]) || (x == i && y == j) {
				continue
			}

			if isNum(engineSchematic[x][y]) {
				num := getNumber(engineSchematic[x], y)
				foundNum[idx] = num
				idx++
			}
		}
	}

	return foundNum[0] * foundNum[1]
}

func getNumber(line []string, j int) int {
	found := line[j]

	x, y := j-1, j+1
	headEnd, tailEnd := false, false

	for {
		if headEnd && tailEnd {
			break
		}

		if x < 0 {
			headEnd = true
		}

		if y >= len(line) {
			tailEnd = true
		}

		if !headEnd {
			if isNum(line[x]) {
				found = line[x] + found
				line[x] = "V"
				x--
			} else {
				headEnd = true
			}
		}

		if !tailEnd {
			if isNum(line[y]) {
				found += line[y]
				line[y] = "V"
				y++
			} else {
				tailEnd = true
			}
		}
	}

	num, err := strconv.Atoi(found)
	if err != nil {
		log.Fatal(err)
	}
	return num
}

func isNum(str string) bool {
	if _, err := strconv.Atoi(str); err == nil {
		return true
	}

	return false
}
