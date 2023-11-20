package main

import (
	"cmp"
	"fmt"
	"io"
	"log"
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
	slc := sliceInput(input)
	cal := groupCalByElf(slc)
	elfsCal := sumCalByElf(cal)
	max := findMaxCals(elfsCal)
	fmt.Println("Calories Own by Elf with Most Calories: ", max)

	sortElfsByCal(elfsCal)
	top3 := sumTop3ElfsCal(elfsCal)
	fmt.Println("Calories Own by Top 3 Elves with Most Calories: ", top3)
}

func getInput() string {
	req, err := http.NewRequest("GET", "https://adventofcode.com/2022/day/1/input", nil)

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

func sliceInput(i string) []string {
	slc := strings.Split(i, "\n")
	return slc
}

func groupCalByElf(cals []string) [][]string {
	out := make([][]string, 0)

	in := make([]string, 0)

	for _, cal := range cals {
		if cal == "" {
			out = append(out, in)
			in = []string{}
			continue
		}

		in = append(in, cal)
	}

	return out
}

func sumCalByElf(cals [][]string) []int {
	elfsCal := make([]int, 0)

	for _, v := range cals {
		t := 0
		for _, cal := range v {
			i, err := strconv.Atoi(cal)

			if err != nil {
				panic(err)
			}

			t += i
		}
		elfsCal = append(elfsCal, t)
	}

	return elfsCal
}

func findMaxCals(elfsCal []int) int {
	return slices.Max[[]int](elfsCal)
}

func sortElfsByCal(elfsCal []int) {
	slices.SortFunc[[]int](elfsCal, func(a, b int) int {
		return cmp.Compare(b, a)
	})
}

func sumTop3ElfsCal(elfsCal []int) int {
	c := 3
	t := 0

	for i := 0; i < c; i++ {
		t += elfsCal[i]
	}

	return t
}
