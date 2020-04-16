package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	filename := flag.String("csv", "problems", "input csv file name")
	timeLimit := flag.Int("limit", 30, "time allowed for the quiz before it terminated")
	flag.Parse()
	csvFile, err := os.Open(fmt.Sprintf("%v.csv", *filename))
	if err != nil {
		exit(fmt.Sprintf("%v, %v\n", err, csvFile))
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	data, err := reader.ReadAll()
	if err != nil {
		exit(fmt.Sprintln(err))
	}

	aCh := make(chan string)
	go func() {
		var answer string
		fmt.Scanf("%s\n", &answer)
		aCh <- answer
	}()
	correct := 0
	quizs := parseQuiz(data)
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	quizLoop:
	for i, q := range quizs {
		fmt.Printf("Quiz #%d: %s = ?\n", i+1, q.q)
		select {
		case <-timer.C:
			break quizLoop
		case answer := <- aCh:	
			if answer == q.a {
				correct++
			}
		}
	}
	
	fmt.Printf("Got %d out of %d correct\n", correct, len(quizs))
}

func parseQuiz(qs [][]string) []quiz {
	ret := make([]quiz, len(qs))

	for i, q := range qs {
		ret[i] = quiz {
			q: q[0],
			a: strings.TrimSpace(q[1]),
		}
	}
	return ret
}

type quiz struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}