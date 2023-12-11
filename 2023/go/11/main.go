package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
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

	images := make([][]string, 0, len(inputSlice))

	for _, v := range inputSlice {
		slc := strings.Split(v, "")
		images = append(images, slc)
	}

	emptyRows := getEmptyRows(images)
	emptyCols := getEmptyColumns(images)

	totalGalaxies := 0
	totalPairs := 0
	galaxiesCoords := make([][2]int, 0, len(images))

	for y, row := range images {
		for x, v := range row {
			if v == "#" {
				totalGalaxies++
				galaxiesCoords = append(galaxiesCoords, [2]int{x, y})
			}
		}
	}

	for i := totalGalaxies - 1; i > 0; i-- {
		totalPairs += i
	}

	fmt.Println("Total Galaxies:", totalGalaxies)
	fmt.Println("Total Pairs:", totalPairs)

	distMatrix := createDistMatrix(galaxiesCoords, emptyRows, emptyCols, 2)

	total := 0

	for y := 0; y < len(distMatrix); y++ {
		row := distMatrix[y]
		for x := y + 1; x < len(row); x++ {
			total += row[x]
		}
	}

	fmt.Println("Sum of Shortest Path between Galaxies where Empty Column and Row Expand 2 Times Larger are", total)

	distMatrix2 := createDistMatrix(galaxiesCoords, emptyRows, emptyCols, 1_000_000)

	total2 := 0

	for y := 0; y < len(distMatrix2); y++ {
		row := distMatrix2[y]
		for x := y + 1; x < len(row); x++ {
			total2 += row[x]
		}
	}

	fmt.Println("Sum of Shortest Path between Galaxies where Empty Column and Row Expand 1 Million Times Larger are", total2)
}

func getInput() string {
	req, err := http.NewRequest("GET", "https://adventofcode.com/2023/day/11/input", nil)

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

func getEmptyColumns(images [][]string) []int {
	emptyColumns := make([]int, 0)
	for x := 0; x < len(images[0]); x++ {
		empty := true

		for y := 0; y < len(images); y++ {
			if images[y][x] == "#" {
				empty = false
				break
			}
		}

		if empty {
			emptyColumns = append(emptyColumns, x)
		}
	}

	return emptyColumns
}

func getEmptyRows(images [][]string) []int {
	emptyRows := make([]int, 0)
	for y := 0; y < len(images); y++ {
		row := images[y]
		empty := true

		for x := 0; x < len(row); x++ {
			if images[y][x] == "#" {
				empty = false
				break
			}
		}

		if empty {
			emptyRows = append(emptyRows, y)
		}
	}

	return emptyRows
}

func createDistMatrix(coords [][2]int, emptyRows []int, emptyCols []int, expansionRate int) [][]int {
	dists := make([][]int, len(coords))
	inf := math.Pow(2, 31)

	for start := 0; start < len(coords); start++ {
		for end := 0; end < len(coords); end++ {
			if start == end {
				dists[start] = append(dists[start], 0)
				continue
			}

			if end < start {
				dists[start] = append(dists[start], int(inf))
				continue
			}
			x, y := 0, 0

			if coords[end][0] > coords[start][0] {
				exp := checkExpansion(coords[start][0], coords[end][0], emptyCols, expansionRate)
				x = coords[end][0] - coords[start][0]
				x += exp
			} else {
				exp := checkExpansion(coords[end][0], coords[start][0], emptyCols, expansionRate)
				x = coords[start][0] - coords[end][0]
				x += exp
			}

			if coords[end][1] > coords[start][1] {
				exp := checkExpansion(coords[start][1], coords[end][1], emptyRows, expansionRate)
				y = coords[end][1] - coords[start][1]
				y += exp
			} else {
				exp := checkExpansion(coords[end][1], coords[start][1], emptyRows, expansionRate)
				y = coords[start][1] - coords[end][1]
				y += exp
			}

			dist := x + y
			dists[start] = append(dists[start], dist)
		}
	}

	return dists
}

func checkExpansion(start int, end int, emptySlc []int, expRate int) int {
	expansion := 0
	empty := 0
	for _, v := range emptySlc {
		if v > start && v < end {
			expansion += expRate
			empty++
		}
	}
	return expansion - empty
}
