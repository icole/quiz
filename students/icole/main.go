package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func getAnswer() string {
	input, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimRight(string(input), "\n")
}

var ErrTimeout = errors.New("timeout")

func answerQuizOrTimeout(lines [][]string) (int, error) {
	finish := make(chan bool)
	answers := make(chan bool)
	timeout := make(chan bool)
	defer close(finish)
	defer close(answers)
	defer close(timeout)
	correct := 0

	go func() {
		finish <- answerQuiz(answers, lines)
	}()

	go func() {
		time.Sleep(3 * time.Second)
		timeout <- true
	}()

	go func() {
		// input <- getAnswer()
		for answer := range answers {
			if answer {
				correct++
			}
		}
	}()

	select {
	case <-finish:
		return correct, nil
	case <-timeout:
		return correct, ErrTimeout
	}
}

func answerQuiz(answers chan bool, lines [][]string) bool {
	for _, rec := range lines {
		fmt.Printf("What is %s?\n", rec[0])
		answer := getAnswer()
		if answer == rec[1] {
			fmt.Println("Correct!")
			answers <- true
		} else {
			fmt.Println("Incorrect!")
			answers <- false
		}
	}
	return true
}

func main() {
	f, err := os.Open("./problems.csv")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	csvReader := csv.NewReader(f)
	// correct := 0
	lines, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	questions := len(lines)

	// TODO: Revisit how to do single in line reads.
	//for rec, err := csvReader.Read(); err == nil {

	correct, err := answerQuizOrTimeout(lines)
	if err != nil {
		if err == ErrTimeout {
			fmt.Println("Time's up!")
		} else {
			log.Fatal(err)
		}
	}

	fmt.Printf("You got %d out of %d correct\n", correct, questions)
}
