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
	inputStr := strings.Split(input, "\n")[0]
	inputSlice := strings.Split(inputStr, ",")

	fmt.Println(len(inputSlice))
	sum := 0
	for _, val := range inputSlice {
		hashVal := hashAlgorithm(val)
		sum += hashVal
	}

	fmt.Println("Sum of Initialization Result", sum)

	boxMap := mappingBox(inputSlice)

	totalPower := 0

	for k, v := range boxMap {
		power := calcFocusPower(v, k+1)
		totalPower += power
	}

	fmt.Println("Total Focusing Power", totalPower)
}

func getInput() string {
	req, err := http.NewRequest("GET", "https://adventofcode.com/2023/day/15/input", nil)

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

func hashAlgorithm(str string) int {
	cur := 0

	for _, c := range str {
		cur += int(c)
		cur *= 17
		cur %= 256
	}

	return cur
}

func mappingBox(seq []string) map[int][]string {
	mapBox := make(map[int][]string)
	hashMaps := make(map[string]int)

	for _, val := range seq {
		if val[len(val)-1] == '-' {
			code := strings.Split(val, "-")[0]
			if _, ok := hashMaps[code]; !ok {
				hashMaps[code] = hashAlgorithm(code)
			}
			boxVals, ok := mapBox[hashMaps[code]]
			if !ok {
				continue
			}

			mapBox[hashMaps[code]] = removeElement(boxVals, code)
		} else if val[len(val)-2] == '=' {
			code := strings.Split(val, "=")[0]
			if _, ok := hashMaps[code]; !ok {
				hashMaps[code] = hashAlgorithm(code)
			}
			boxVals, ok := mapBox[hashMaps[code]]
			if !ok {
				mapBox[hashMaps[code]] = []string{val}
				continue
			}

			mapBox[hashMaps[code]] = addOrUpdateElement(boxVals, val, code)
		}
	}

	return mapBox
}

func removeElement(boxVals []string, code string) []string {
	res := make([]string, 0, len(boxVals))

	for _, v := range boxVals {
		currCode := strings.Split(v, "=")[0]
		if currCode == code {
			continue
		}
		res = append(res, v)
	}

	return res
}

func addOrUpdateElement(boxVals []string, val string, code string) []string {
	found := false

	for idx, v := range boxVals {
		currCode := strings.Split(v, "=")[0]
		if currCode == code {
			boxVals[idx] = val
			found = true
			break
		}
	}

	if !found {
		boxVals = append(boxVals, val)
	}

	return boxVals
}

func calcFocusPower(boxVals []string, power int) int {
	res := 0

	for idx, v := range boxVals {
		focalLenStr := strings.Split(v, "=")[1]
		focalLen, err := strconv.Atoi(focalLenStr)
		if err != nil {
			log.Fatal(err)
		}
		res += power * (idx + 1) * focalLen
	}

	return res
}
