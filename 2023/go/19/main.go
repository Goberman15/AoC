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

type flow struct {
	part       string
	op         string
	threeshold int
	next       string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
// 
	input := getInput()
	// input := "px{a<2006:qkq,m>2090:A,rfg}\npv{a>1716:R,A}\nlnx{m>1548:A,A}\nrfg{s<537:gd,x>2440:R,A}\nqs{s>3448:A,lnx}\nqkq{x<1416:A,crn}\ncrn{x>2662:A,R}\nin{s<1351:px,qqz}\nqqz{s>2770:qs,m<1801:hdj,R}\ngd{a>3333:R,R}\nhdj{m>838:A,pv}\n\n{x=787,m=2655,a=1222,s=2876}\n{x=1679,m=44,a=2067,s=496}\n{x=2036,m=264,a=79,s=2244}\n{x=2461,m=1339,a=466,s=291}\n{x=2127,m=1623,a=2188,s=1013}"
	inputSlice := strings.Split(input, "\n\n")

	rawWorkFlows := strings.Split(inputSlice[0], "\n")
	rawParts := strings.Split(inputSlice[1], "\n")
	rawParts = rawParts[:len(rawParts)-1]

	workflow := parseWorkflow(rawWorkFlows)
	parts := ParseParts(rawParts)

	fmt.Printf("%+v\n", workflow)

	total := 0

	for _, part := range parts {
		
		if isAccepted := evaluatepart(part, workflow); isAccepted {
			total += part["total"]
		}
	}

	fmt.Println("Total Part Rating", total)
}

func getInput() string {
	req, err := http.NewRequest("GET", "https://adventofcode.com/2023/day/19/input", nil)

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

func parseWorkflow(raw []string) map[string][]flow {
	workflow := make(map[string][]flow)

	for _, line := range raw {
		splitted1 := strings.Split(line, "{")
		flowName := splitted1[0]

		splitted2 := strings.Split(splitted1[1][:len(splitted1[1])-1], ",")

		flowSlc := make([]flow, 0, len(splitted2))

		for idx, rule := range splitted2 {
			if idx == len(splitted2)-1 {
				flowSlc = append(flowSlc, flow{
					next: rule,
				})
			} else {
				splittedRule := strings.Split(rule, ":")

				flowSlc = append(flowSlc, flow{
					part:       string(splittedRule[0][0]),
					op:         string(splittedRule[0][1]),
					threeshold: stringToInt(splittedRule[0][2:]),
					next:       splittedRule[1],
				})
			}
		}

		workflow[flowName] = flowSlc
	}

	return workflow
}

func ParseParts(raw []string) []map[string]int {
	parts := make([]map[string]int, 0, len(raw))

	for _, line := range raw {
		partMap := make(map[string]int)

		splitted1 := strings.Split(line[1:len(line)-1], ",")

		total := 0
		for _, part := range splitted1 {
			splitted2 := strings.Split(part, "=")
			partLabel := splitted2[0]
			partRating := stringToInt(splitted2[1])

			partMap[partLabel] = partRating
			total += partRating
		}

		partMap["total"] = total

		parts = append(parts, partMap)
	}

	return parts
}

func evaluatepart(part map[string]int, workflows map[string][]flow) bool {
	cur := "in"

	for {
		flows := workflows[cur]

		for idx, flow := range flows {
			if idx != len(flows)-1 {
				testedRating := part[flow.part]

				if flow.op == ">" {
					if testedRating <= flow.threeshold {
						continue
					}
				} else {
					if testedRating >= flow.threeshold {
						continue
					}
				}
			}

			if flow.next == "R" {
				return false
			} else if flow.next == "A" {
				return true
			} else {
				cur = flow.next
				break
			}
		}
	}
}

func stringToInt(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		log.Fatal(err)
	}
	return num
}
