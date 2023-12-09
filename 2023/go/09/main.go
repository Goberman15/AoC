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

	totalForward, totalBackward := 0, 0

	for _, x := range inputSlice {
		strSlice := strings.Split(x, " ")
		seq := toIntSlice(strSlice)

		fwd := getNextSequence(seq)
		totalForward += fwd
		
		bwd := getPrevSequence(seq)
		totalBackward += bwd
	}
	

	fmt.Printf("Sum of Extrapolated Forward Value is %d\n", totalForward)
	fmt.Printf("Sum of Extrapolated Backward Value is %d\n", totalBackward)

}

func getInput() string {
	req, err := http.NewRequest("GET", "https://adventofcode.com/2023/day/9/input", nil)

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

func getNextSequence(seq []int) int {
	if containsAll(seq, seq[0]) {
		return seq[0]
	}

	diffs := make([]int, 0, len(seq)-1)

	for x := 1; x < len(seq); x++ {
		diff := seq[x] - seq[x-1]
		diffs = append(diffs, diff)
	}

	return seq[len(seq)-1] + getNextSequence(diffs)
}

func getPrevSequence(seq []int) int {
	if containsAll(seq, seq[0]) {
		return seq[0]
	}

	diffs := make([]int, 0, len(seq)-1)

	for x := 1; x < len(seq); x++ {
		diff := seq[x] - seq[x-1]
		diffs = append(diffs, diff)
	}

	return seq[0] - getPrevSequence(diffs)
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

func containsAll[S ~[]E, E comparable](slc S, x E) bool {
	for _, v := range slc {
		if v != x {
			return false
		}

	}
	return true
}
