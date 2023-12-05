package main

import (
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
	inputSlice := strings.Split(input, "\n\n")
	// inputSlice = inputSlice[:len(inputSlice)-1]

	seeds := strings.Split(inputSlice[0], " ")[1:]
	maps := make([][]string, 0, 7)

	for _, v := range inputSlice[1:] {
		m := strings.Split(v, "\n")[1:]

		if m[len(m)-1] == "" {
			m = m[:len(m)-1]
		}

		maps = append(maps, m)
	}

	locations := make([]int, 0, len(seeds))

	for _, seed := range seeds {
		// seedToSoil
		soil := getMapping(maps[0], seed)
		// soilToFertilizer
		fert := getMapping(maps[1], soil)
		// fertilizerToWater
		water := getMapping(maps[2], fert)
		// waterToLight
		light := getMapping(maps[3], water)
		// lightToTemp
		temp := getMapping(maps[4], light)
		// tempToHumidity
		humidity := getMapping(maps[5], temp)
		// humidityToLocation
		loc := getMapping(maps[6], humidity)

		intLoc, err := strconv.Atoi(loc)
		if err != nil {
			log.Fatal(err)
		}

		locations = append(locations, intLoc)
	}

	minLocation := slices.Min(locations)
	fmt.Printf("The Lowest Location Number is %d\n", minLocation)

	actualMinLocations := make([]int, 0)

	for x := 0; x < len(seeds); x += 2 {
		start, err := strconv.Atoi(seeds[x])
		if err != nil {
			log.Fatal(err)
		}
		rng, err := strconv.Atoi(seeds[x+1])
		if err != nil {
			log.Fatal(err)
		}

		seed := make([][2]int, 0)

		seed = append(seed, [2]int{start, start + rng - 1})
		// seedToSoil
		soil := getMapping2(maps[0], seed)
		// soilToFertilizer
		fert := getMapping2(maps[1], soil)
		// // fertilizerToWater
		water := getMapping2(maps[2], fert)
		// // waterToLight
		light := getMapping2(maps[3], water)
		// // lightToTemp
		temp := getMapping2(maps[4], light)
		// // tempToHumidity
		humidity := getMapping2(maps[5], temp)
		// // humidityToLocation
		loc := getMapping2(maps[6], humidity)

		minLocs := make([]int, 0, len(loc))

		for _, v := range loc {
			minLocs = append(minLocs, v[0])
		}

		minLoc := slices.Min(minLocs)
		fmt.Println(minLoc)

		actualMinLocations = append(actualMinLocations, minLoc)
	}

	actualMinLocation := slices.Min(actualMinLocations)
	fmt.Printf("The Lowest Location Number for Actual Seeds is %d", actualMinLocation)
}

func getInput() string {
	req, err := http.NewRequest("GET", "https://adventofcode.com/2023/day/5/input", nil)

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

func getMapping(mapping []string, src string) string {
	s, err := strconv.Atoi(src)
	if err != nil {
		log.Fatal(err)
	}

	dest := -1

	for _, v := range mapping {
		splittedV := strings.Split(v, " ")
		destStart, err := strconv.Atoi(splittedV[0])
		if err != nil {
			log.Fatal(err)
		}
		srcStart, err := strconv.Atoi(splittedV[1])
		if err != nil {
			log.Fatal(err)
		}
		rng, err := strconv.Atoi(splittedV[2])
		if err != nil {
			log.Fatal(err)
		}

		if s >= srcStart && s < (srcStart+rng) {
			dest = destStart + (s - srcStart)
			break
		}
	}

	if dest == -1 {
		dest = s
	}

	return fmt.Sprint(dest)
}

func getMapping2(mapping []string, rgs [][2]int) [][2]int {
	coll := make([][2]int, 0)
	rngColl := make([][2]int, 0)

	rngColl = append(rngColl, rgs...)

	for _, v := range mapping {
		splittedV := strings.Split(v, " ")
		destStart, err := strconv.Atoi(splittedV[0])
		if err != nil {
			log.Fatal(err)
		}
		srcStart, err := strconv.Atoi(splittedV[1])
		if err != nil {
			log.Fatal(err)
		}
		rng, err := strconv.Atoi(splittedV[2])
		if err != nil {
			log.Fatal(err)
		}

		srcEnd := srcStart + rng - 1

		tempColl := make([][2]int, 0)

		for _, rg := range rngColl {
			if rg[0] < srcStart && rg[0] < srcEnd && rg[1] > srcStart && rg[1] > srcEnd {
				tempColl = append(tempColl, [2]int{rg[0], srcStart - 1})
				coll = append(coll, [2]int{destStart, destStart + rng - 1})
				tempColl = append(tempColl, [2]int{srcEnd + 1, rg[1]})
			} else if rg[0] < srcStart && rg[0] < srcEnd && rg[1] > srcStart && rg[1] < srcEnd {
				tempColl = append(tempColl, [2]int{rg[0], srcStart - 1})
				coll = append(coll, [2]int{destStart, destStart + rg[1] - srcStart})
			} else if rg[0] > srcStart && rg[0] < srcEnd && rg[1] > srcStart && rg[1] > srcEnd {
				coll = append(coll, [2]int{destStart + rg[0] - srcStart, destStart + rng - 1})
				tempColl = append(tempColl, [2]int{srcEnd + 1, rg[1]})
			} else if rg[0] > srcStart && rg[0] < srcEnd && rg[1] > srcStart && rg[1] < srcEnd {
				coll = append(coll, [2]int{destStart + rg[0] - srcStart, destStart + rg[1] - srcStart})
			} else if rg[0] == srcStart && rg[1] > srcStart && rg[1] > srcEnd {
				coll = append(coll, [2]int{destStart, destStart + rng - 1})
				tempColl = append(tempColl, [2]int{rg[0] + rng, rg[1]})
			} else if rg[0] == srcStart && rg[1] > srcStart && rg[1] < srcEnd {
				coll = append(coll, [2]int{destStart, destStart + rg[1] - rg[0]})
			} else if rg[0] > srcStart && rg[0] < srcEnd && rg[1] == srcEnd {
				coll = append(coll, [2]int{destStart + rg[0] - srcStart, destStart + rng - 1})
			} else if rg[0] < srcStart && rg[0] < srcEnd && rg[1] == srcEnd {
				coll = append(coll, [2]int{destStart, destStart + rng - 1})
				tempColl = append(tempColl, [2]int{rg[0], srcStart -1})
			} else {
				tempColl = append(tempColl, rg)
			}
		}

		copy(rngColl, tempColl)
	}

	return append(coll, rngColl...)
}
