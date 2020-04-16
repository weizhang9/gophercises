package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	filename := flag.String("csv", "problems", "input csv file name")
	t := flag.Int("timer", 30, "time allowed for the quiz before it terminated")
	flag.Parse()
	count := 0
	csvFile, err := os.Open(fmt.Sprintf("%v.csv", *filename))
	if err != nil {
		log.Fatalln(err)
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	data, err := reader.ReadAll()
	if err != nil {
		log.Fatalln(err)
	}
	
	for _, line := range data {
		timer := time.AfterFunc(time.Second * time.Duration(*t), func() {
			fmt.Printf("\nGot %v out of %v correct\n", count, len(data))
			os.Exit(0)
		})
		defer timer.Stop()
		fmt.Print(line[0], " = ? ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		answer := scanner.Text()
		
		if answer == line[1] {
			count++
		}
	}
	
	fmt.Printf("Got %v out of %v correct\n", count, len(data))
}

type secondsTimer struct {
	timer *time.Timer
	end time.Time
}

func newSecondsTimer(t time.Duration) *secondsTimer {
	return &secondsTimer{time.NewTimer(t), time.Now().Add(t)}
}

func (s *secondsTimer) timeRemaining() time.Duration {
	return s.end.Sub(time.Now())
}

func (s *secondsTimer) stop() {
	s.timer.Stop()
}