package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {

	csvFile := flag.String("csv", "problems.csv", "a csv file in format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")

	flag.Parse()

	file, err := os.Open(*csvFile)

	if err != nil {
		exit(fmt.Sprintf("Failed to open csv file: %s", *csvFile))

	}

	reader := csv.NewReader(file)
	lines, _ := reader.ReadAll()

	if err != nil {
		exit("Failed to read the provided CSV file")
	}

	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0

	for i, problem := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, problem.question)
		answerCh := make(chan string)

		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("you scored %d out of %d \n", correct, len(problems))
			return
		case answer := <-answerCh:
			if answer == problem.answer {
				correct++
			}
		}
	}

	fmt.Printf("you scored %d out of %d \n", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))

	for i, line := range lines {
		ret[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}

	return ret
}

type problem struct {
	question string
	answer   string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
